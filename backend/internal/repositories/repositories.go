package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repositories bundles every collection repository for dependency injection.
type Repositories struct {
	Users  *UserRepository
	Rooms  *RoomRepository
	Tasks  *TaskRepository
	Votes  *VoteRepository
	Finals *FinalRepository
}

// New wires repositories against a database handle.
func New(db *mongo.Database) *Repositories {
	return &Repositories{
		Users:  &UserRepository{col: db.Collection("users")},
		Rooms:  &RoomRepository{col: db.Collection("rooms")},
		Tasks:  &TaskRepository{col: db.Collection("tasks")},
		Votes:  &VoteRepository{col: db.Collection("votes")},
		Finals: &FinalRepository{col: db.Collection("final_decisions")},
	}
}

// EnsureIndexes creates the indexes called out in the architecture doc.
func (r *Repositories) EnsureIndexes(ctx context.Context) error {
	unique := func(keys bson.D) mongo.IndexModel {
		return mongo.IndexModel{Keys: keys, Options: options.Index().SetUnique(true)}
	}
	plain := func(keys bson.D) mongo.IndexModel {
		return mongo.IndexModel{Keys: keys}
	}

	if _, err := r.Users.col.Indexes().CreateOne(ctx, unique(bson.D{{Key: "email", Value: 1}})); err != nil {
		return err
	}
	if _, err := r.Rooms.col.Indexes().CreateOne(ctx, unique(bson.D{{Key: "roomCode", Value: 1}})); err != nil {
		return err
	}
	if _, err := r.Tasks.col.Indexes().CreateOne(ctx, plain(bson.D{{Key: "roomId", Value: 1}})); err != nil {
		return err
	}
	// One vote per user per task; updates upsert on this key.
	if _, err := r.Votes.col.Indexes().CreateOne(ctx, unique(bson.D{{Key: "taskId", Value: 1}, {Key: "userId", Value: 1}})); err != nil {
		return err
	}
	if _, err := r.Finals.col.Indexes().CreateOne(ctx, unique(bson.D{{Key: "taskId", Value: 1}})); err != nil {
		return err
	}
	return nil
}
