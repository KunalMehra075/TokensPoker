package handlers

import (
	"net/http"

	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/services"

	"github.com/gin-gonic/gin"
)

// RoomHandler exposes room CRUD and join.
type RoomHandler struct {
	svc *services.RoomService
}

// NewRoomHandler builds the room handler.
func NewRoomHandler(svc *services.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

// Create makes a new room owned by the caller.
func (h *RoomHandler) Create(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, http.StatusBadRequest, dto.CodeValidation, "a room name is required")
		return
	}
	u := currentUser(c)
	room, err := h.svc.Create(c.Request.Context(), u.UserID, u.Email, u.Name, req.Name)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusCreated, room)
}

// Preview returns the public, minimal view of a room for an invite link. No
// authentication required.
func (h *RoomHandler) Preview(c *gin.Context) {
	preview, err := h.svc.Preview(c.Request.Context(), c.Param("code"))
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, preview)
}

// Get returns a room the caller belongs to.
func (h *RoomHandler) Get(c *gin.Context) {
	room, err := h.svc.Get(c.Request.Context(), c.Param("id"), currentUser(c).UserID)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, room)
}

// Join adds the caller to a room by code.
func (h *RoomHandler) Join(c *gin.Context) {
	var req dto.JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, http.StatusBadRequest, dto.CodeValidation, "a room code is required")
		return
	}
	u := currentUser(c)
	room, err := h.svc.Join(c.Request.Context(), req.RoomCode, u.UserID, u.Email, u.Name)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, room)
}

// Delete removes a room (owner only).
func (h *RoomHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id"), currentUser(c).UserID); err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, gin.H{"deleted": true})
}
