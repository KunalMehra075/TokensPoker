package handlers

import (
	"net/http"

	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler exposes login and current-user endpoints.
type AuthHandler struct {
	svc *services.AuthService
}

// NewAuthHandler builds the auth handler.
func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login upserts a user by email and returns a JWT. No OTP, no password.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, http.StatusBadRequest, dto.CodeValidation, "a valid name and email are required")
		return
	}
	token, user, err := h.svc.Login(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, dto.AuthResponse{Token: token, User: *user})
}

// Me returns the authenticated user.
func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.svc.Me(c.Request.Context(), currentUser(c).UserID)
	if err != nil {
		respondErr(c, err)
		return
	}
	dto.OK(c, http.StatusOK, user)
}
