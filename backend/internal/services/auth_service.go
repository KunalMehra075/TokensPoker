package services

import (
	"context"
	"net/http"

	"freetokenspoker/internal/apperr"
	"freetokenspoker/internal/auth"
	"freetokenspoker/internal/models"
	"freetokenspoker/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthService handles email-identity login. No OTP, no passwords.
type AuthService struct {
	users *repositories.UserRepository
	jwt   *auth.JWTManager
}

// NewAuthService wires the auth service.
func NewAuthService(users *repositories.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

// Login upserts a user by email and returns a signed identity token.
func (s *AuthService) Login(ctx context.Context, name, email string) (string, *models.User, error) {
	user, err := s.users.UpsertByEmail(ctx, email, name)
	if err != nil {
		return "", nil, apperr.Internal("could not save user")
	}
	token, err := s.jwt.Generate(user.ID.Hex(), user.Email, user.Name)
	if err != nil {
		return "", nil, apperr.Internal("could not issue token")
	}
	return token, user, nil
}

// Me returns the current user from their id.
func (s *AuthService) Me(ctx context.Context, userID string) (*models.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperr.New(http.StatusBadRequest, "VALIDATION_ERROR", "invalid user id")
	}
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.NotFound("user not found")
	}
	return user, nil
}
