package repositories

import (
	"context"
	"errors"

	"freetokenspoker/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoomRepository persists rooms and their embedded members.
type RoomRepository struct {
	col *mongo.Collection
}

// Create inserts a new room.
func (r *RoomRepository) Create(ctx context.Context, room *models.Room) error {
	res, err := r.col.InsertOne(ctx, room)
	if err != nil {
		return err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID loads a room by its object id.
func (r *RoomRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Room, error) {
	return r.findOne(ctx, bson.M{"_id": id})
}

// FindByCode loads a room by its short join code.
func (r *RoomRepository) FindByCode(ctx context.Context, code string) (*models.Room, error) {
	return r.findOne(ctx, bson.M{"roomCode": code})
}

func (r *RoomRepository) findOne(ctx context.Context, filter bson.M) (*models.Room, error) {
	var room models.Room
	err := r.col.FindOne(ctx, filter).Decode(&room)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// CodeExists reports whether a room code is already taken.
func (r *RoomRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	n, err := r.col.CountDocuments(ctx, bson.M{"roomCode": code})
	return n > 0, err
}

// AddMember appends a member if not already present (idempotent on userId).
func (r *RoomRepository) AddMember(ctx context.Context, roomID primitive.ObjectID, m models.Member) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"_id": roomID, "members.userId": bson.M{"$ne": m.UserID}},
		bson.M{"$push": bson.M{"members": m}},
	)
	return err
}

// Delete removes a room.
func (r *RoomRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ListByMember returns every room the user belongs to, newest first.
func (r *RoomRepository) ListByMember(ctx context.Context, userID primitive.ObjectID) ([]models.Room, error) {
	cur, err := r.col.Find(ctx, bson.M{"members.userId": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	rooms := []models.Room{}
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
