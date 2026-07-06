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

// FinalRepository persists one final decision per task.
type FinalRepository struct {
	col *mongo.Collection
}

// Upsert records (or replaces) the final decision for a task.
func (r *FinalRepository) Upsert(ctx context.Context, f *models.FinalDecision) (*models.FinalDecision, error) {
	now := time.Now().UTC()
	filter := bson.M{"taskId": f.TaskID}
	update := bson.M{
		"$set": bson.M{
			"finalValue": f.FinalValue,
			"selectedBy": f.SelectedBy,
			"createdAt":  now,
		},
		"$setOnInsert": bson.M{
			"taskId": f.TaskID,
			"roomId": f.RoomID,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var stored models.FinalDecision
	if err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&stored); err != nil {
		return nil, err
	}
	return &stored, nil
}

// FindByTask returns the final decision for a task, if any.
func (r *FinalRepository) FindByTask(ctx context.Context, taskID primitive.ObjectID) (*models.FinalDecision, error) {
	var f models.FinalDecision
	err := r.col.FindOne(ctx, bson.M{"taskId": taskID}).Decode(&f)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// FindByTaskIDs returns final decisions for a set of tasks, keyed by taskId hex.
func (r *FinalRepository) FindByTaskIDs(ctx context.Context, ids []primitive.ObjectID) (map[string]models.FinalDecision, error) {
	out := map[string]models.FinalDecision{}
	if len(ids) == 0 {
		return out, nil
	}
	cur, err := r.col.Find(ctx, bson.M{"taskId": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var finals []models.FinalDecision
	if err := cur.All(ctx, &finals); err != nil {
		return nil, err
	}
	for _, f := range finals {
		out[f.TaskID.Hex()] = f
	}
	return out, nil
}
