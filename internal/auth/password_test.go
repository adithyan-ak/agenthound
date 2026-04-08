package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("test-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
		t.Errorf("hash should be bcrypt format, got %q", hash[:10])
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if err := CheckPassword(hash, "correct-password"); err != nil {
		t.Errorf("CheckPassword() with correct password: %v", err)
	}
	if err := CheckPassword(hash, "wrong-password"); err == nil {
		t.Error("CheckPassword() with wrong password should fail")
	}
}

func TestHashPasswordUnique(t *testing.T) {
	h1, _ := HashPassword("same-password")
	h2, _ := HashPassword("same-password")
	if h1 == h2 {
		t.Error("bcrypt hashes of same password should differ (different salts)")
	}
}
