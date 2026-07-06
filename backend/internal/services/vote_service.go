package services

import (
	"context"
	"errors"

	"freetokenspoker/internal/apperr"
	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/models"
	"freetokenspoker/internal/realtime"
	"freetokenspoker/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VoteService handles casting and changing votes.
type VoteService struct {
	votes *repositories.VoteRepository
	tasks *repositories.TaskRepository
	rooms *repositories.RoomRepository
	rt    realtime.Broadcaster
}

// NewVoteService wires the vote service.
func NewVoteService(
	votes *repositories.VoteRepository,
	tasks *repositories.TaskRepository,
	rooms *repositories.RoomRepository,
	rt realtime.Broadcaster,
) *VoteService {
	return &VoteService{votes: votes, tasks: tasks, rooms: rooms, rt: rt}
}

// Submit casts or updates the caller's vote for an active, unrevealed task.
func (s *VoteService) Submit(ctx context.Context, userID, name string, req dto.SubmitVoteRequest) (*models.Vote, error) {
	tid, err := primitive.ObjectIDFromHex(req.TaskID)
	if err != nil {
		return nil, apperr.BadRequest("invalid task id")
	}
	task, err := s.tasks.FindByID(ctx, tid)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, apperr.NotFound("task not found")
	}
	if err != nil {
		return nil, apperr.Internal("could not load task")
	}
	if task.Revealed || task.Status != models.TaskStatusActive {
		return nil, apperr.Conflict("voting is closed for this task")
	}

	room, err := s.rooms.FindByID(ctx, task.RoomID)
	if err != nil {
		return nil, apperr.Internal("could not load room")
	}
	if !isMember(room, userID) {
		return nil, apperr.Forbidden("you are not a member of this room")
	}

	def, _ := models.ModeByName(task.Mode)
	if !def.IsValidCard(req.SelectedCard) {
		return nil, apperr.BadRequest("card is not valid for this estimation mode")
	}

	uid, _ := primitive.ObjectIDFromHex(userID)
	stored, err := s.votes.Upsert(ctx, &models.Vote{
		TaskID:       task.ID,
		RoomID:       room.ID,
		UserID:       uid,
		UserName:     name,
		SelectedCard: req.SelectedCard,
	})
	if err != nil {
		return nil, apperr.Internal("could not save vote")
	}

	count, _ := s.votes.CountByTask(ctx, task.ID)
	// Pre-reveal we only announce that this member voted, never the card value.
	s.rt.EmitToRoom(room.ID.Hex(), realtime.EventVoteReceived, map[string]any{
		"taskId":      task.ID.Hex(),
		"userId":      userID,
		"name":        name,
		"hasVoted":    true,
		"voteCount":   count,
		"memberCount": len(room.Members),
	})
	return stored, nil
}
