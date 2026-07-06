package handlers

import (
	"net/http"

	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/models"

	"github.com/gin-gonic/gin"
)

// MetaHandler exposes non-resource endpoints: health and the mode catalog.
type MetaHandler struct{}

// NewMetaHandler builds the meta handler.
func NewMetaHandler() *MetaHandler { return &MetaHandler{} }

// Health is a simple liveness probe.
func (h *MetaHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "freetokenspoker-api"})
}

// Modes returns the estimation mode catalog so the UI never hardcodes cards.
func (h *MetaHandler) Modes(c *gin.Context) {
	dto.OK(c, http.StatusOK, models.EstimationModes)
}
