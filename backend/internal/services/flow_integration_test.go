package services

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/models"
	"freetokenspoker/internal/realtime"
	"freetokenspoker/internal/repositories"

	"go.mongodb.org/mongo-driver/mongo"
)

// recordingBroadcaster captures every emitted realtime event so tests can
// assert the service layer broadcasts the right things at the right time.
type recordingBroadcaster struct {
	mu     sync.Mutex
	events []recordedEvent
}

type recordedEvent struct {
	RoomID  string
	Event   string
	Payload any
}

func (r *recordingBroadcaster) EmitToRoom(roomID, event string, payload any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, recordedEvent{roomID, event, payload})
}

func (r *recordingBroadcaster) count(event string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, e := range r.events {
		if e.Event == event {
			n++
		}
	}
	return n
}

func (r *recordingBroadcaster) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = nil
}

// testDB connects to a throwaway database, or skips the test when no Mongo is
// configured. Set TEST_MONGODB_URI (or MONGODB_URI) to run these against a real
// MongoDB / Atlas instance; without it the suite skips instead of failing so CI
// without a database stays green.
func testDB(t *testing.T) (*repositories.Repositories, func()) {
	t.Helper()

	uri := os.Getenv("TEST_MONGODB_URI")
	if uri == "" {
		uri = os.Getenv("MONGODB_URI")
	}
	if uri == "" {
		t.Skip("no TEST_MONGODB_URI or MONGODB_URI set; skipping Mongo integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := repositories.Connect(ctx, uri)
	if err != nil {
		t.Skipf("could not reach MongoDB at %s: %v; skipping integration test", uri, err)
	}

	dbName := fmt.Sprintf("ftp_test_%d", time.Now().UnixNano())
	db := client.Database(dbName)
	repos := repositories.New(db)
	if err := repos.EnsureIndexes(ctx); err != nil {
		t.Fatalf("EnsureIndexes: %v", err)
	}

	cleanup := func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = db.Drop(c)
		_ = client.Disconnect(c)
	}
	return repos, cleanup
}

func mustUser(t *testing.T, repos *repositories.Repositories, email, name string) *models.User {
	t.Helper()
	u, err := repos.Users.UpsertByEmail(context.Background(), email, name)
	if err != nil {
		t.Fatalf("UpsertByEmail(%s): %v", email, err)
	}
	return u
}

// TestFullEstimationFlow drives the real service stack against MongoDB:
// create room -> join -> create task -> vote (hidden) -> reveal -> final.
func TestFullEstimationFlow(t *testing.T) {
	repos, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()
	rt := &recordingBroadcaster{}
	rooms := NewRoomService(repos.Rooms, rt)
	tasks := NewTaskService(repos.Tasks, repos.Rooms, repos.Votes, repos.Finals, rt)
	votes := NewVoteService(repos.Votes, repos.Tasks, repos.Rooms, rt)

	owner := mustUser(t, repos, "owner@example.com", "Olivia Owner")
	member := mustUser(t, repos, "member@example.com", "Manny Member")

	// --- Create room -------------------------------------------------------
	room, err := rooms.Create(ctx, owner.ID.Hex(), owner.Email, owner.Name, "Sprint 42")
	if err != nil {
		t.Fatalf("Create room: %v", err)
	}
	if len(room.Members) != 1 || room.Members[0].UserID != owner.ID {
		t.Fatalf("new room should have the owner as its only member, got %+v", room.Members)
	}
	if len(room.RoomCode) == 0 {
		t.Fatal("room should have a generated code")
	}

	// --- Preview (public, unauthenticated) ---------------------------------
	preview, err := rooms.Preview(ctx, room.RoomCode)
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	if preview.Name != "Sprint 42" || preview.MemberCount != 1 {
		t.Errorf("preview = %+v, want name=Sprint 42 memberCount=1", preview)
	}

	// --- Join --------------------------------------------------------------
	rt.reset()
	if _, err := rooms.Join(ctx, room.RoomCode, member.ID.Hex(), member.Email, member.Name); err != nil {
		t.Fatalf("Join: %v", err)
	}
	if rt.count(realtime.EventMemberJoined) != 1 {
		t.Errorf("member_joined emitted %d times, want 1", rt.count(realtime.EventMemberJoined))
	}
	// Re-join by owner is idempotent and emits nothing new.
	rt.reset()
	if _, err := rooms.Join(ctx, room.RoomCode, owner.ID.Hex(), owner.Email, owner.Name); err != nil {
		t.Fatalf("owner re-join: %v", err)
	}
	if rt.count(realtime.EventMemberJoined) != 0 {
		t.Error("re-joining an existing member should not emit member_joined")
	}

	reloaded, err := rooms.Get(ctx, room.ID.Hex(), owner.ID.Hex())
	if err != nil {
		t.Fatalf("Get room: %v", err)
	}
	if len(reloaded.Members) != 2 {
		t.Fatalf("room should have 2 members after join, got %d", len(reloaded.Members))
	}

	// --- Only the owner can create a task ----------------------------------
	if _, err := tasks.Create(ctx, member.ID.Hex(), dto.CreateTaskRequest{
		RoomID: room.ID.Hex(), Title: "X", Mode: models.ModeTokens,
	}); err == nil {
		t.Error("a non-owner should not be able to create a task")
	}

	// --- Create task -------------------------------------------------------
	rt.reset()
	detail, err := tasks.Create(ctx, owner.ID.Hex(), dto.CreateTaskRequest{
		RoomID:      room.ID.Hex(),
		Title:       "Add semantic search",
		Description: "Vector index + query rewrite",
		Mode:        models.ModeTokens,
	})
	if err != nil {
		t.Fatalf("Create task: %v", err)
	}
	if detail.Task.Status != models.TaskStatusActive {
		t.Errorf("new task status = %q, want ACTIVE", detail.Task.Status)
	}
	if rt.count(realtime.EventTaskCreated) != 1 {
		t.Error("creating a task should emit task_created once")
	}
	taskID := detail.Task.ID.Hex()

	// Only one active task at a time.
	if _, err := tasks.Create(ctx, owner.ID.Hex(), dto.CreateTaskRequest{
		RoomID: room.ID.Hex(), Title: "Second", Mode: models.ModeDays,
	}); err == nil {
		t.Error("a second active task should be rejected")
	}

	// --- Votes (hidden until reveal) ---------------------------------------
	if _, err := votes.Submit(ctx, owner.ID.Hex(), owner.Name, dto.SubmitVoteRequest{
		TaskID: taskID, SelectedCard: "5M",
	}); err != nil {
		t.Fatalf("owner vote: %v", err)
	}
	if _, err := votes.Submit(ctx, member.ID.Hex(), member.Name, dto.SubmitVoteRequest{
		TaskID: taskID, SelectedCard: "1M",
	}); err != nil {
		t.Fatalf("member vote: %v", err)
	}
	// Invalid card for this mode is rejected.
	if _, err := votes.Submit(ctx, owner.ID.Hex(), owner.Name, dto.SubmitVoteRequest{
		TaskID: taskID, SelectedCard: "$5",
	}); err == nil {
		t.Error("a COST card should be invalid for a TOKENS task")
	}
	// Changing a vote is an upsert, not a new vote.
	if _, err := votes.Submit(ctx, owner.ID.Hex(), owner.Name, dto.SubmitVoteRequest{
		TaskID: taskID, SelectedCard: "10M",
	}); err != nil {
		t.Fatalf("owner change vote: %v", err)
	}

	// Before reveal, cards must be hidden even though HasVoted is true.
	pre, err := tasks.Get(ctx, taskID, member.ID.Hex())
	if err != nil {
		t.Fatalf("Get pre-reveal: %v", err)
	}
	if pre.VoteCount != 2 {
		t.Errorf("VoteCount = %d, want 2 (upsert, not duplicate)", pre.VoteCount)
	}
	for _, v := range pre.Votes {
		if !v.HasVoted {
			t.Errorf("member %s should be marked as voted", v.Name)
		}
		if v.Card != "" {
			t.Errorf("card for %s leaked before reveal: %q", v.Name, v.Card)
		}
	}

	// --- Reveal (owner only) ----------------------------------------------
	if _, err := tasks.Reveal(ctx, taskID, member.ID.Hex()); err == nil {
		t.Error("a non-owner should not be able to reveal")
	}
	rt.reset()
	revealed, err := tasks.Reveal(ctx, taskID, owner.ID.Hex())
	if err != nil {
		t.Fatalf("Reveal: %v", err)
	}
	if !revealed.Task.Revealed || revealed.Task.Status != models.TaskStatusRevealed {
		t.Errorf("after reveal status=%q revealed=%v", revealed.Task.Status, revealed.Task.Revealed)
	}
	if rt.count(realtime.EventVotesRevealed) != 1 {
		t.Error("reveal should emit votes_revealed once")
	}
	cards := map[string]string{}
	for _, v := range revealed.Votes {
		cards[v.Name] = v.Card
	}
	if cards[owner.Name] != "10M" {
		t.Errorf("owner card after reveal = %q, want 10M (the changed vote)", cards[owner.Name])
	}
	if cards[member.Name] != "1M" {
		t.Errorf("member card after reveal = %q, want 1M", cards[member.Name])
	}

	// --- Final decision (owner only) --------------------------------------
	rt.reset()
	final, err := tasks.Final(ctx, taskID, owner.ID.Hex(), "8M")
	if err != nil {
		t.Fatalf("Final: %v", err)
	}
	if final.Task.Status != models.TaskStatusClosed {
		t.Errorf("after final status = %q, want CLOSED", final.Task.Status)
	}
	if final.Final == nil || final.Final.FinalValue != "8M" {
		t.Errorf("final decision = %+v, want FinalValue=8M", final.Final)
	}
	if rt.count(realtime.EventFinalDecision) != 1 || rt.count(realtime.EventTaskClosed) != 1 {
		t.Error("final should emit both final_decision and task_closed")
	}

	// Voting is closed once the task is closed.
	if _, err := votes.Submit(ctx, member.ID.Hex(), member.Name, dto.SubmitVoteRequest{
		TaskID: taskID, SelectedCard: "2M",
	}); err == nil {
		t.Error("voting on a closed task should be rejected")
	}

	// A new active task can be started now that the previous one is closed.
	if _, err := tasks.Create(ctx, owner.ID.Hex(), dto.CreateTaskRequest{
		RoomID: room.ID.Hex(), Title: "Next", Mode: models.ModeModel,
	}); err != nil {
		t.Errorf("should be able to start a new task after closing: %v", err)
	}
}

// TestJoinUnknownRoom verifies a clear not-found path.
func TestJoinUnknownRoom(t *testing.T) {
	repos, cleanup := testDB(t)
	defer cleanup()

	rt := &recordingBroadcaster{}
	rooms := NewRoomService(repos.Rooms, rt)
	u := mustUser(t, repos, "lost@example.com", "Lost")

	if _, err := rooms.Join(context.Background(), "ZZZZZZ", u.ID.Hex(), u.Email, u.Name); err == nil {
		t.Error("joining a non-existent code should return an error")
	}
	if _, err := rooms.Preview(context.Background(), "ZZZZZZ"); err == nil {
		t.Error("previewing a non-existent code should return an error")
	}
}

var _ = mongo.ErrNoDocuments
