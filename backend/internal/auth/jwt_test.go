package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndVerifyRoundTrip(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)

	token, err := m.Generate("user-123", "dev@example.com", "Dev User")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if token == "" {
		t.Fatal("Generate returned an empty token")
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user-123")
	}
	if claims.Email != "dev@example.com" {
		t.Errorf("Email = %q, want %q", claims.Email, "dev@example.com")
	}
	if claims.Name != "Dev User" {
		t.Errorf("Name = %q, want %q", claims.Name, "Dev User")
	}
	// Subject mirrors UserID so standard JWT tooling can read it.
	if claims.Subject != "user-123" {
		t.Errorf("Subject = %q, want %q", claims.Subject, "user-123")
	}
}

func TestVerifyEmptyToken(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)
	if _, err := m.Verify(""); err == nil {
		t.Fatal("Verify(\"\") should return an error")
	}
}

func TestVerifyRejectsWrongSecret(t *testing.T) {
	signer := NewJWTManager("real-secret", time.Hour)
	attacker := NewJWTManager("other-secret", time.Hour)

	token, err := signer.Generate("u1", "a@b.com", "A")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := attacker.Verify(token); err == nil {
		t.Fatal("token signed with a different secret should not verify")
	}
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	// Negative lifetime makes the token expire the instant it is issued.
	m := NewJWTManager("test-secret", -time.Minute)
	token, err := m.Generate("u1", "a@b.com", "A")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := m.Verify(token); err == nil {
		t.Fatal("expired token should not verify")
	}
}

func TestVerifyRejectsTamperedToken(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)
	token, err := m.Generate("u1", "a@b.com", "A")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	// Flip the final character of the signature.
	tampered := token[:len(token)-1]
	if token[len(token)-1] == 'a' {
		tampered += "b"
	} else {
		tampered += "a"
	}
	if _, err := m.Verify(tampered); err == nil {
		t.Fatal("tampered token should not verify")
	}
}

func TestVerifyRejectsNoneAlgorithm(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)
	// A token signed with the "none" algorithm must be rejected: Verify pins
	// HMAC as the only accepted signing method.
	claims := Claims{UserID: "u1", RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}}
	unsigned := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := unsigned.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("signing none token: %v", err)
	}
	if _, err := m.Verify(tokenString); err == nil {
		t.Fatal("token with alg=none should not verify")
	}
}
