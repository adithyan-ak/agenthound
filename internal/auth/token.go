package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	tokenPrefix   = "ah_"
	tokenRawBytes = 32
	tokenAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GenerateAPIToken() (raw, hash string, err error) {
	b := make([]byte, tokenRawBytes)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate token: %w", err)
	}

	encoded := make([]byte, tokenRawBytes)
	for i, v := range b {
		encoded[i] = tokenAlphabet[int(v)%len(tokenAlphabet)]
	}

	raw = tokenPrefix + string(encoded)
	hash = HashToken(raw)
	return raw, hash, nil
}

func HashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
