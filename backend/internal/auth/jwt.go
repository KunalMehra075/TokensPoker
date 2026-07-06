package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload identifying a user. Email is an identity label, not
// a secret, so the token simply attributes actions to a user.
type Claims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// JWTManager signs and verifies identity tokens.
type JWTManager struct {
	secret []byte
	expiry time.Duration
}

// NewJWTManager builds a manager with a signing secret and token lifetime.
func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), expiry: expiry}
}

// Generate issues a signed token for a user.
func (m *JWTManager) Generate(userID, email, name string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Verify parses and validates a token, returning its claims.
func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
