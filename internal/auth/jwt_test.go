package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTGenerateAndParse(t *testing.T) {
	manager := NewJWTManager("01234567890123456789012345678901", "test-issuer", time.Hour)
	userID := uuid.New()
	sessionID := uuid.New()

	token, expiresAt, err := manager.Generate(userID, sessionID, 2, "admin")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if time.Until(expiresAt) < 55*time.Minute {
		t.Fatalf("unexpected expiry: %v", expiresAt)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Subject != userID.String() {
		t.Fatalf("subject = %s, want %s", claims.Subject, userID)
	}
	if claims.SessionID != sessionID {
		t.Fatalf("session id = %s, want %s", claims.SessionID, sessionID)
	}
	if claims.RoleLevel != 2 {
		t.Fatalf("role level = %d, want 2", claims.RoleLevel)
	}
	if claims.Username != "admin" {
		t.Fatalf("username = %s, want admin", claims.Username)
	}
}

func TestJWTRejectsWrongIssuer(t *testing.T) {
	manager := NewJWTManager("01234567890123456789012345678901", "issuer-a", time.Hour)
	token, _, err := manager.Generate(uuid.New(), uuid.New(), 0, "developer")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	other := NewJWTManager("01234567890123456789012345678901", "issuer-b", time.Hour)
	if _, err := other.Parse(token); err == nil {
		t.Fatal("expected issuer validation error")
	}
}

func TestJWTRejectsWrongSecret(t *testing.T) {
	manager := NewJWTManager("01234567890123456789012345678901", "issuer", time.Hour)
	token, _, err := manager.Generate(uuid.New(), uuid.New(), 0, "developer")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	other := NewJWTManager("abcdefghijklmnopqrstuvwxyz123456", "issuer", time.Hour)
	if _, err := other.Parse(token); err == nil {
		t.Fatal("expected signature validation error")
	}
}
