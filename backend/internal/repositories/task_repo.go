package repositories

import (
	"context"
	"errors"
	"time"

	"freetokenspoker/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TaskRepository persists estimation tasks.
type TaskRepository struct {
	col *mongo.Collection
}

// Create inserts a new task.
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	res, err := r.col.InsertOne(ctx, task)
	if err != nil {
		return err
	}
	task.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID loads a task by id.
func (r *TaskRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Task, error) {
	var task models.Task
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// ListByRoom returns a room's tasks, newest first.
func (r *TaskRepository) ListByRoom(ctx context.Context, roomID primitive.ObjectID) ([]models.Task, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.col.Find(ctx, bson.M{"roomId": roomID}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	tasks := []models.Task{}
	if err := cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// ListClosedByRooms returns closed tasks across several rooms, newest first.
func (r *TaskRepository) ListClosedByRooms(ctx context.Context, roomIDs []primitive.ObjectID) ([]models.Task, error) {
	if len(roomIDs) == 0 {
		return []models.Task{}, nil
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(200)
	cur, err := r.col.Find(ctx, bson.M{
		"roomId": bson.M{"$in": roomIDs},
		"status": models.TaskStatusClosed,
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	tasks := []models.Task{}
	if err := cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// ActiveExists reports whether the room already has a non-closed task.
func (r *TaskRepository) ActiveExists(ctx context.Context, roomID primitive.ObjectID) (bool, error) {
	n, err := r.col.CountDocuments(ctx, bson.M{
		"roomId": roomID,
		"status": bson.M{"$ne": models.TaskStatusClosed},
	})
	return n > 0, err
}

// Reveal flips a task to revealed state.
func (r *TaskRepository) Reveal(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now().UTC()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"revealed":   true,
		"status":     models.TaskStatusRevealed,
		"revealedAt": now,
	}})
	return err
}

// Close marks a task closed after a final decision.
func (r *TaskRepository) Close(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now().UTC()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"status":   models.TaskStatusClosed,
		"closedAt": now,
	}})
	return err
}
