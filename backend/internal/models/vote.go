package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vote is a single participant's private selection for a task. It stays hidden
// until the task is revealed.
type Vote struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TaskID       primitive.ObjectID `bson:"taskId" json:"taskId"`
	RoomID       primitive.ObjectID `bson:"roomId" json:"roomId"`
	UserID       primitive.ObjectID `bson:"userId" json:"userId"`
	UserName     string             `bson:"userName" json:"userName"`
	SelectedCard string             `bson:"selectedCard" json:"selectedCard"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}
