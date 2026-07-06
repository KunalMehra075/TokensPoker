package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Member is an embedded snapshot of a user inside a room. Denormalized so the
// room view does not need a join for names.
type Member struct {
	UserID   primitive.ObjectID `bson:"userId" json:"userId"`
	Email    string             `bson:"email" json:"email"`
	Name     string             `bson:"name" json:"name"`
	JoinedAt time.Time          `bson:"joinedAt" json:"joinedAt"`
}

// Room represents a collaborative estimation session.
type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomCode  string             `bson:"roomCode" json:"roomCode"`
	Name      string             `bson:"name" json:"name"`
	OwnerID   primitive.ObjectID `bson:"ownerId" json:"ownerId"`
	Members   []Member           `bson:"members" json:"members"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
