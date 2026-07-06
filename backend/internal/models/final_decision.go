package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FinalDecision records the value the room owner committed to after discussion.
type FinalDecision struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TaskID     primitive.ObjectID `bson:"taskId" json:"taskId"`
	RoomID     primitive.ObjectID `bson:"roomId" json:"roomId"`
	FinalValue string             `bson:"finalValue" json:"finalValue"`
	SelectedBy primitive.ObjectID `bson:"selectedBy" json:"selectedBy"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}
