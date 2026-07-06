package handlers

import (
	"net/http"

	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/services"

	"github.com/gin-gonic/gin"
)

// TaskHandler exposes task lifecycle endpoints.
type TaskHandler struct {
	svc *services.TaskService
}

// NewTaskHandler builds the task handler.
func NewTaskHandler(svc *services.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// Create opens a new estimation task (owner only).
func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, http.StatusBadRequest, dto.CodeValidation, "title, room and mode are required")
		return
	}
	detail, err := h.svc.Create(c.Request.Context(), currentUser(c).UserID, req)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusCreated, detail)
}

// Get returns full task state.
func (h *TaskHandler) Get(c *gin.Context) {
	detail, err := h.svc.Get(c.Request.Context(), c.Param("id"), currentUser(c).UserID)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, detail)
}

// ListByRoom returns all tasks in a room.
func (h *TaskHandler) ListByRoom(c *gin.Context) {
	details, err := h.svc.ListByRoom(c.Request.Context(), c.Param("id"), currentUser(c).UserID)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, details)
}

// Reveal makes all votes visible (owner only).
func (h *TaskHandler) Reveal(c *gin.Context) {
	detail, err := h.svc.Reveal(c.Request.Context(), c.Param("id"), currentUser(c).UserID)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, detail)
}

// Final commits a final value and closes the task (owner only).
func (h *TaskHandler) Final(c *gin.Context) {
	var req dto.FinalDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, http.StatusBadRequest, dto.CodeValidation, "a final value is required")
		return
	}
	detail, err := h.svc.Final(c.Request.Context(), c.Param("id"), currentUser(c).UserID, req.FinalValue)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, detail)
}
