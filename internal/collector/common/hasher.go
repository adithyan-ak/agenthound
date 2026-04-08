package common

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func HashSHA256(input string) string {
	h := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", h)
}

func CanonicalJSONHash(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("canonical json marshal: %w", err)
	}
	return HashSHA256(string(data)), nil
}

func DescriptionHash(name, description string, schema any) string {
	m := map[string]any{
		"name":        name,
		"description": description,
	}
	if schema != nil {
		m["schema"] = schema
	}

	hash, err := CanonicalJSONHash(m)
	if err != nil {
		return HashSHA256(name + ":" + description)
	}
	return hash
}
