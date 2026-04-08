package auth

import (
	"strings"
	"testing"
)

func TestGenerateAPIToken(t *testing.T) {
	raw, hash, err := GenerateAPIToken()
	if err != nil {
		t.Fatalf("GenerateAPIToken() error = %v", err)
	}

	if !strings.HasPrefix(raw, "ah_") {
		t.Errorf("raw token should start with ah_, got %q", raw[:6])
	}
	if len(raw) != 3+tokenRawBytes { // "ah_" + 32 chars
		t.Errorf("raw token len = %d, want %d", len(raw), 3+tokenRawBytes)
	}
	if len(hash) != 64 { // SHA-256 hex
		t.Errorf("hash len = %d, want 64", len(hash))
	}
}

func TestGenerateAPITokenUnique(t *testing.T) {
	r1, h1, _ := GenerateAPIToken()
	r2, h2, _ := GenerateAPIToken()
	if r1 == r2 {
		t.Error("two generated tokens should be different")
	}
	if h1 == h2 {
		t.Error("two generated hashes should be different")
	}
}

func TestHashTokenDeterministic(t *testing.T) {
	h1 := HashToken("ah_test123")
	h2 := HashToken("ah_test123")
	if h1 != h2 {
		t.Error("HashToken should be deterministic")
	}
}

func TestHashTokenDifferentInputs(t *testing.T) {
	h1 := HashToken("ah_token1")
	h2 := HashToken("ah_token2")
	if h1 == h2 {
		t.Error("different inputs should produce different hashes")
	}
}
