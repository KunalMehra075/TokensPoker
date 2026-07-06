package realtime

import (
	"encoding/json"
	"testing"
)

func TestNoopBroadcasterDoesNotPanic(t *testing.T) {
	var b Broadcaster = NoopBroadcaster{}
	// Should be a safe no-op for any input.
	b.EmitToRoom("room-1", EventVoteReceived, map[string]any{"x": 1})
	b.EmitToRoom("", "", nil)
}

func TestEnvelopeJSONShape(t *testing.T) {
	env := Envelope{
		Event:   EventVotesRevealed,
		RoomID:  "abc123",
		Payload: map[string]any{"count": 3},
	}
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// The wire contract the frontend parses is {event, roomId, payload}.
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded["event"] != EventVotesRevealed {
		t.Errorf("event = %v, want %q", decoded["event"], EventVotesRevealed)
	}
	if decoded["roomId"] != "abc123" {
		t.Errorf("roomId = %v, want abc123", decoded["roomId"])
	}
	if _, ok := decoded["payload"]; !ok {
		t.Error("payload key missing from envelope JSON")
	}
}

func TestEventNamesMatchWireContract(t *testing.T) {
	// These string values must stay in lockstep with the frontend's
	// SOCKET_EVENTS map in src/constants/index.ts. Changing one side without
	// the other silently breaks realtime updates, so pin them here.
	want := map[string]string{
		"EventPresence":      "presence",
		"EventMemberJoined":  "member_joined",
		"EventMemberLeft":    "member_left",
		"EventTaskCreated":   "task_created",
		"EventVoteReceived":  "vote_received",
		"EventVotesRevealed": "votes_revealed",
		"EventTaskClosed":    "task_closed",
		"EventFinalDecision": "final_decision",
	}
	got := map[string]string{
		"EventPresence":      EventPresence,
		"EventMemberJoined":  EventMemberJoined,
		"EventMemberLeft":    EventMemberLeft,
		"EventTaskCreated":   EventTaskCreated,
		"EventVoteReceived":  EventVoteReceived,
		"EventVotesRevealed": EventVotesRevealed,
		"EventTaskClosed":    EventTaskClosed,
		"EventFinalDecision": EventFinalDecision,
	}
	for name, wantVal := range want {
		if got[name] != wantVal {
			t.Errorf("%s = %q, want %q", name, got[name], wantVal)
		}
	}
}
