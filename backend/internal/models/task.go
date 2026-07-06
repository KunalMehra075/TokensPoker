package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskStatus is the lifecycle state of an estimation task.
type TaskStatus string

const (
	TaskStatusActive   TaskStatus = "ACTIVE"
	TaskStatusRevealed TaskStatus = "REVEALED"
	TaskStatusClosed   TaskStatus = "CLOSED"
)

// Task is a single estimation round inside a room.
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID      primitive.ObjectID `bson:"roomId" json:"roomId"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Mode        EstimationMode     `bson:"mode" json:"mode"`
	Status      TaskStatus         `bson:"status" json:"status"`
	Revealed    bool               `bson:"revealed" json:"revealed"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	RevealedAt  *time.Time         `bson:"revealedAt,omitempty" json:"revealedAt,omitempty"`
	ClosedAt    *time.Time         `bson:"closedAt,omitempty" json:"closedAt,omitempty"`
}
