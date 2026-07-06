package services

import (
	"context"

	"freetokenspoker/internal/apperr"
	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/models"
	"freetokenspoker/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HistoryService assembles a user's archived estimations across their rooms.
type HistoryService struct {
	rooms  *repositories.RoomRepository
	tasks  *repositories.TaskRepository
	votes  *repositories.VoteRepository
	finals *repositories.FinalRepository
}

// NewHistoryService wires the history service.
func NewHistoryService(
	rooms *repositories.RoomRepository,
	tasks *repositories.TaskRepository,
	votes *repositories.VoteRepository,
	finals *repositories.FinalRepository,
) *HistoryService {
	return &HistoryService{rooms: rooms, tasks: tasks, votes: votes, finals: finals}
}

// List returns closed tasks (with final values and votes) from every room the
// user belongs to, newest first.
func (s *HistoryService) List(ctx context.Context, userID string) ([]dto.HistoryItem, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperr.BadRequest("invalid user id")
	}
	rooms, err := s.rooms.ListByMember(ctx, uid)
	if err != nil {
		return nil, apperr.Internal("could not load rooms")
	}
	if len(rooms) == 0 {
		return []dto.HistoryItem{}, nil
	}

	roomIDs := make([]primitive.ObjectID, 0, len(rooms))
	roomMeta := map[string]models.Room{}
	for _, r := range rooms {
		roomIDs = append(roomIDs, r.ID)
		roomMeta[r.ID.Hex()] = r
	}

	tasks, err := s.tasks.ListClosedByRooms(ctx, roomIDs)
	if err != nil {
		return nil, apperr.Internal("could not load tasks")
	}
	if len(tasks) == 0 {
		return []dto.HistoryItem{}, nil
	}

	taskIDs := make([]primitive.ObjectID, 0, len(tasks))
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}
	finals, err := s.finals.FindByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, apperr.Internal("could not load final decisions")
	}
	votes, err := s.votes.ListByTasks(ctx, taskIDs)
	if err != nil {
		return nil, apperr.Internal("could not load votes")
	}

	items := make([]dto.HistoryItem, 0, len(tasks))
	for _, t := range tasks {
		room := roomMeta[t.RoomID.Hex()]
		item := dto.HistoryItem{
			Task:     t,
			RoomName: room.Name,
			RoomCode: room.RoomCode,
			Votes:    votes[t.ID.Hex()],
		}
		if f, ok := finals[t.ID.Hex()]; ok {
			fc := f
			item.Final = &fc
		}
		if item.Votes == nil {
			item.Votes = []models.Vote{}
		}
		items = append(items, item)
	}
	return items, nil
}
