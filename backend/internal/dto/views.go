package dto

import "freetokenspoker/internal/models"

// MemberVoteView is a per-member voting status. The card stays hidden until the
// task is revealed so the reveal stays simultaneous and unbiased.
type MemberVoteView struct {
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	HasVoted bool   `json:"hasVoted"`
	Card     string `json:"card,omitempty"`
}

// TaskDetail is the full state of a task for the room UI.
type TaskDetail struct {
	Task        models.Task           `json:"task"`
	Votes       []MemberVoteView      `json:"votes"`
	Final       *models.FinalDecision `json:"final,omitempty"`
	VoteCount   int                   `json:"voteCount"`
	MemberCount int                   `json:"memberCount"`
}

// HistoryItem is one archived estimation in a user's history.
type HistoryItem struct {
	Task     models.Task           `json:"task"`
	RoomName string                `json:"roomName"`
	RoomCode string                `json:"roomCode"`
	Final    *models.FinalDecision `json:"final,omitempty"`
	Votes    []models.Vote         `json:"votes"`
}
