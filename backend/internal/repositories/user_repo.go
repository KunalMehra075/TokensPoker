package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"freetokenspoker/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrNotFound is returned when a document does not exist.
var ErrNotFound = errors.New("not found")

// UserRepository persists users keyed by email.
type UserRepository struct {
	col *mongo.Collection
}

// UpsertByEmail finds-or-creates a user by email and refreshes name + lastLogin.
// Email is normalized to lowercase so it stays a stable identity key.
func (r *UserRepository) UpsertByEmail(ctx context.Context, email, name string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	now := time.Now().UTC()

	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{"name": strings.TrimSpace(name), "lastLogin": now},
		"$setOnInsert": bson.M{
			"email":     email,
			"createdAt": now,
		},
	}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var user models.User
	if err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID loads a user by id.
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
