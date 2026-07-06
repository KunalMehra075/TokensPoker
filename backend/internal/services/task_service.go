package services

import (
	"context"
	"errors"
	"sort"

	"freetokenspoker/internal/apperr"
	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/models"
	"freetokenspoker/internal/realtime"
	"freetokenspoker/internal/repositories"
	"freetokenspoker/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskService owns the estimation task lifecycle and reveal/final flow.
type TaskService struct {
	tasks  *repositories.TaskRepository
	rooms  *repositories.RoomRepository
	votes  *repositories.VoteRepository
	finals *repositories.FinalRepository
	rt     realtime.Broadcaster
}

// NewTaskService wires the task service.
func NewTaskService(
	tasks *repositories.TaskRepository,
	rooms *repositories.RoomRepository,
	votes *repositories.VoteRepository,
	finals *repositories.FinalRepository,
	rt realtime.Broadcaster,
) *TaskService {
	return &TaskService{tasks: tasks, rooms: rooms, votes: votes, finals: finals, rt: rt}
}

// Create opens a new task. Only the room owner may create tasks, and only one
// task can be active at a time.
func (s *TaskService) Create(ctx context.Context, userID string, req dto.CreateTaskRequest) (*dto.TaskDetail, error) {
	room, err := s.ownerRoom(ctx, req.RoomID, userID)
	if err != nil {
		return nil, err
	}
	if _, ok := models.ModeByName(req.Mode); !ok {
		return nil, apperr.BadRequest("unknown estimation mode")
	}
	active, err := s.tasks.ActiveExists(ctx, room.ID)
	if err != nil {
		return nil, apperr.Internal("could not check active tasks")
	}
	if active {
		return nil, apperr.Conflict("finish the current task before starting a new one")
	}

	creatorID, _ := primitive.ObjectIDFromHex(userID)
	task := &models.Task{
		RoomID:      room.ID,
		Title:       req.Title,
		Description: req.Description,
		Mode:        req.Mode,
		Status:      models.TaskStatusActive,
		Revealed:    false,
		CreatedBy:   creatorID,
		CreatedAt:   utils.Now(),
	}
	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, apperr.Internal("could not create task")
	}

	detail := &dto.TaskDetail{Task: *task, Votes: []dto.MemberVoteView{}, MemberCount: len(room.Members)}
	s.rt.EmitToRoom(room.ID.Hex(), realtime.EventTaskCreated, detail)
	return detail, nil
}

// Get returns full task state for a room member.
func (s *TaskService) Get(ctx context.Context, taskID, userID string) (*dto.TaskDetail, error) {
	task, room, err := s.taskWithRoom(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if !isMember(room, userID) {
		return nil, apperr.Forbidden("you are not a member of this room")
	}
	return s.buildDetail(ctx, task, room)
}

// ListByRoom returns every task in a room with summarized state.
func (s *TaskService) ListByRoom(ctx context.Context, roomID, userID string) ([]dto.TaskDetail, error) {
	rid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, apperr.BadRequest("invalid room id")
	}
	room, err := s.rooms.FindByID(ctx, rid)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, apperr.NotFound("room not found")
	}
	if err != nil {
		return nil, apperr.Internal("could not load room")
	}
	if !isMember(room, userID) {
		return nil, apperr.Forbidden("you are not a member of this room")
	}
	tasks, err := s.tasks.ListByRoom(ctx, rid)
	if err != nil {
		return nil, apperr.Internal("could not list tasks")
	}
	out := make([]dto.TaskDetail, 0, len(tasks))
	for i := range tasks {
		detail, err := s.buildDetail(ctx, &tasks[i], room)
		if err != nil {
			return nil, err
		}
		out = append(out, *detail)
	}
	return out, nil
}

// Reveal makes all votes visible. Owner-only.
func (s *TaskService) Reveal(ctx context.Context, taskID, userID string) (*dto.TaskDetail, error) {
	task, room, err := s.taskWithRoom(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if room.OwnerID.Hex() != userID {
		return nil, apperr.Forbidden("only the room owner can reveal")
	}
	if task.Status == models.TaskStatusClosed {
		return nil, apperr.Conflict("task is already closed")
	}
	if err := s.tasks.Reveal(ctx, task.ID); err != nil {
		return nil, apperr.Internal("could not reveal task")
	}
	task.Revealed = true
	task.Status = models.TaskStatusRevealed

	detail, err := s.buildDetail(ctx, task, room)
	if err != nil {
		return nil, err
	}
	s.rt.EmitToRoom(room.ID.Hex(), realtime.EventVotesRevealed, detail)
	return detail, nil
}

// Final commits a value and closes the task. Owner-only.
func (s *TaskService) Final(ctx context.Context, taskID, userID, finalValue string) (*dto.TaskDetail, error) {
	task, room, err := s.taskWithRoom(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if room.OwnerID.Hex() != userID {
		return nil, apperr.Forbidden("only the room owner can set the final decision")
	}
	selectedBy, _ := primitive.ObjectIDFromHex(userID)
	final := &models.FinalDecision{
		TaskID:     task.ID,
		RoomID:     room.ID,
		FinalValue: finalValue,
		SelectedBy: selectedBy,
	}
	if _, err := s.finals.Upsert(ctx, final); err != nil {
		return nil, apperr.Internal("could not save final decision")
	}
	if err := s.tasks.Close(ctx, task.ID); err != nil {
		return nil, apperr.Internal("could not close task")
	}
	task.Status = models.TaskStatusClosed

	detail, err := s.buildDetail(ctx, task, room)
	if err != nil {
		return nil, err
	}
	s.rt.EmitToRoom(room.ID.Hex(), realtime.EventFinalDecision, detail)
	s.rt.EmitToRoom(room.ID.Hex(), realtime.EventTaskClosed, detail)
	return detail, nil
}

// buildDetail assembles a TaskDetail, hiding card values until revealed.
func (s *TaskService) buildDetail(ctx context.Context, task *models.Task, room *models.Room) (*dto.TaskDetail, error) {
	votes, err := s.votes.ListByTask(ctx, task.ID)
	if err != nil {
		return nil, apperr.Internal("could not load votes")
	}
	voted := map[string]models.Vote{}
	for _, v := range votes {
		voted[v.UserID.Hex()] = v
	}

	views := make([]dto.MemberVoteView, 0, len(room.Members))
	for _, m := range room.Members {
		uid := m.UserID.Hex()
		v, has := voted[uid]
		view := dto.MemberVoteView{UserID: uid, Name: m.Name, HasVoted: has}
		if has && task.Revealed {
			view.Card = v.SelectedCard
		}
		views = append(views, view)
	}
	sort.Slice(views, func(i, j int) bool { return views[i].Name < views[j].Name })

	detail := &dto.TaskDetail{
		Task:        *task,
		Votes:       views,
		VoteCount:   len(votes),
		MemberCount: len(room.Members),
	}
	if task.Status == models.TaskStatusClosed || task.Revealed {
		if final, err := s.finals.FindByTask(ctx, task.ID); err == nil {
			detail.Final = final
		}
	}
	return detail, nil
}

func (s *TaskService) ownerRoom(ctx context.Context, roomID, userID string) (*models.Room, error) {
	rid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, apperr.BadRequest("invalid room id")
	}
	room, err := s.rooms.FindByID(ctx, rid)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, apperr.NotFound("room not found")
	}
	if err != nil {
		return nil, apperr.Internal("could not load room")
	}
	if room.OwnerID.Hex() != userID {
		return nil, apperr.Forbidden("only the room owner can do that")
	}
	return room, nil
}

func (s *TaskService) taskWithRoom(ctx context.Context, taskID string) (*models.Task, *models.Room, error) {
	tid, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, nil, apperr.BadRequest("invalid task id")
	}
	task, err := s.tasks.FindByID(ctx, tid)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, nil, apperr.NotFound("task not found")
	}
	if err != nil {
		return nil, nil, apperr.Internal("could not load task")
	}
	room, err := s.rooms.FindByID(ctx, task.RoomID)
	if err != nil {
		return nil, nil, apperr.Internal("could not load room")
	}
	return task, room, nil
}
