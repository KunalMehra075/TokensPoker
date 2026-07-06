package dto

import "freetokenspoker/internal/models"

// LoginRequest carries the email identity (name + email) used to upsert a user.
type LoginRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=80"`
	Email string `json:"email" binding:"required,email,max=160"`
}

// AuthResponse is returned after a successful login.
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// CreateRoomRequest creates a new room owned by the caller.
type CreateRoomRequest struct {
	Name string `json:"name" binding:"required,min=1,max=80"`
}

// JoinRoomRequest joins an existing room by its short code.
type JoinRoomRequest struct {
	RoomCode string `json:"roomCode" binding:"required,min=4,max=12"`
}

// RoomPreview is the minimal, public view of a room behind an invite link. It
// lets an unauthenticated invitee see what they are joining without exposing
// members or the internal room id.
type RoomPreview struct {
	Name        string `json:"name"`
	RoomCode    string `json:"roomCode"`
	MemberCount int    `json:"memberCount"`
}

// CreateTaskRequest opens a new estimation round in a room.
type CreateTaskRequest struct {
	RoomID      string                `json:"roomId" binding:"required"`
	Title       string                `json:"title" binding:"required,min=1,max=160"`
	Description string                `json:"description" binding:"max=2000"`
	Mode        models.EstimationMode `json:"mode" binding:"required"`
}

// FinalDecisionRequest records the committed estimate for a task.
type FinalDecisionRequest struct {
	FinalValue string `json:"finalValue" binding:"required,max=80"`
}

// SubmitVoteRequest casts or updates a vote for a task.
type SubmitVoteRequest struct {
	TaskID       string `json:"taskId" binding:"required"`
	SelectedCard string `json:"selectedCard" binding:"required,max=40"`
}
