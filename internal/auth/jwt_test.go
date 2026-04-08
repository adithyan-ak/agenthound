package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	secret := "test-secret-key"
	token, expiresAt, err := GenerateJWT("user-1", "admin", "admin", secret, 24*time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateJWT() returned empty token")
	}
	if expiresAt.Before(time.Now()) {
		t.Error("expiresAt should be in the future")
	}

	claims, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("UserID = %q, want user-1", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("Username = %q, want admin", claims.Username)
	}
	if claims.Role != "admin" {
		t.Errorf("Role = %q, want admin", claims.Role)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	token, _, err := GenerateJWT("user-1", "admin", "admin", "secret-1", time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Error("ValidateJWT() with wrong secret should fail")
	}
}

func TestValidateJWTExpired(t *testing.T) {
	token, _, err := GenerateJWT("user-1", "admin", "admin", "secret", -time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	_, err = ValidateJWT(token, "secret")
	if err == nil {
		t.Error("ValidateJWT() with expired token should fail")
	}
}

func TestValidateJWTInvalid(t *testing.T) {
	_, err := ValidateJWT("not-a-jwt", "secret")
	if err == nil {
		t.Error("ValidateJWT() with invalid token should fail")
	}
}
