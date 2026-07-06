package handlers

import (
	"net/http"

	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/services"

	"github.com/gin-gonic/gin"
)

// VoteHandler exposes vote submission and updates.
type VoteHandler struct {
	svc *services.VoteService
}

// NewVoteHandler builds the vote handler.
func NewVoteHandler(svc *services.VoteService) *VoteHandler {
	return &VoteHandler{svc: svc}
}

// Submit casts or updates the caller's vote. Both POST and PATCH map here since
// a vote is idempotent per user per task.
func (h *VoteHandler) Submit(c *gin.Context) {
	var req dto.SubmitVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, http.StatusBadRequest, dto.CodeValidation, "a task id and selected card are required")
		return
	}
	u := currentUser(c)
	vote, err := h.svc.Submit(c.Request.Context(), u.UserID, u.Name, req)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, vote)
}
