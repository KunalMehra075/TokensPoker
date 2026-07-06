package repositories

import (
	"context"
	"time"

	"freetokenspoker/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VoteRepository persists votes, one per user per task.
type VoteRepository struct {
	col *mongo.Collection
}

// Upsert casts or changes a user's vote for a task, returning the stored vote.
func (r *VoteRepository) Upsert(ctx context.Context, v *models.Vote) (*models.Vote, error) {
	now := time.Now().UTC()
	filter := bson.M{"taskId": v.TaskID, "userId": v.UserID}
	update := bson.M{
		"$set": bson.M{
			"selectedCard": v.SelectedCard,
			"userName":     v.UserName,
			"updatedAt":    now,
		},
		"$setOnInsert": bson.M{
			"taskId":    v.TaskID,
			"roomId":    v.RoomID,
			"userId":    v.UserID,
			"createdAt": now,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var stored models.Vote
	if err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&stored); err != nil {
		return nil, err
	}
	return &stored, nil
}

// ListByTask returns all votes for a task.
func (r *VoteRepository) ListByTask(ctx context.Context, taskID primitive.ObjectID) ([]models.Vote, error) {
	cur, err := r.col.Find(ctx, bson.M{"taskId": taskID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	votes := []models.Vote{}
	if err := cur.All(ctx, &votes); err != nil {
		return nil, err
	}
	return votes, nil
}

// CountByTask returns how many votes a task has received.
func (r *VoteRepository) CountByTask(ctx context.Context, taskID primitive.ObjectID) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{"taskId": taskID})
}

// ListByTasks returns votes for several tasks grouped by taskId hex.
func (r *VoteRepository) ListByTasks(ctx context.Context, taskIDs []primitive.ObjectID) (map[string][]models.Vote, error) {
	out := map[string][]models.Vote{}
	if len(taskIDs) == 0 {
		return out, nil
	}
	cur, err := r.col.Find(ctx, bson.M{"taskId": bson.M{"$in": taskIDs}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var votes []models.Vote
	if err := cur.All(ctx, &votes); err != nil {
		return nil, err
	}
	for _, v := range votes {
		out[v.TaskID.Hex()] = append(out[v.TaskID.Hex()], v)
	}
	return out, nil
}

var _ = mongo.ErrNoDocuments
