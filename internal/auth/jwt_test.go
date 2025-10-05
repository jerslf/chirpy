package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTCreateAndValidate(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()
	expiresIn := time.Hour

	// Create JWT
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Validate JWT
	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if validatedID != userID {
		t.Errorf("Expected userID %v, got %v", userID, validatedID)
	}
}

func TestJWTExpired(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()

	// Expired token (1 second ago)
	token, err := MakeJWT(userID, secret, -1*time.Second)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("Expected error for expired token, got none")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	secret := "mysecret"
	wrongSecret := "wrongsecret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Errorf("Expected error for invalid secret, got none")
	}
}
