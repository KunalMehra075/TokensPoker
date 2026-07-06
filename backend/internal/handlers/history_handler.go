package handlers

import (
	"net/http"

	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/services"

	"github.com/gin-gonic/gin"
)

// HistoryHandler exposes a user's archived estimations.
type HistoryHandler struct {
	svc *services.HistoryService
}

// NewHistoryHandler builds the history handler.
func NewHistoryHandler(svc *services.HistoryService) *HistoryHandler {
	return &HistoryHandler{svc: svc}
}

// List returns the caller's estimation history.
func (h *HistoryHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context(), currentUser(c).UserID)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, items)
}
