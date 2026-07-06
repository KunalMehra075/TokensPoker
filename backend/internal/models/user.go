package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is identified solely by email. Email is a data-association key, not an
// authentication secret. There are no passwords and no OTP in V1.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	LastLogin time.Time          `bson:"lastLogin" json:"lastLogin"`
}
