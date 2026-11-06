package auth

import (
	"strings"
	"testing"
)

func TestPasswordHashAndCheck(t *testing.T) {
	hash, err := HashPassword("Secret123!")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("expected bcrypt hash, got %s", hash)
	}
	if !CheckPassword(hash, "Secret123!") {
		t.Fatal("expected password to match")
	}
	if CheckPassword(hash, "WrongPassword!") {
		t.Fatal("expected password mismatch")
	}
}

func TestRefreshTokenHashing(t *testing.T) {
	raw, hash, err := newRefreshToken()
	if err != nil {
		t.Fatalf("new refresh token: %v", err)
	}
	if raw == "" || hash == "" {
		t.Fatal("expected non-empty token and hash")
	}
	if len(hash) != 64 {
		t.Fatalf("hash length = %d, want 64", len(hash))
	}
	if hashRefreshToken(raw) != hash {
		t.Fatal("expected deterministic hash")
	}

	otherRaw, otherHash, err := newRefreshToken()
	if err != nil {
		t.Fatalf("new refresh token: %v", err)
	}
	if raw == otherRaw || hash == otherHash {
		t.Fatal("expected unique refresh tokens")
	}
}
