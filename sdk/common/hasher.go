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

// HashCredentialValue is the cross-collector merge primitive on
// :Credential nodes. The Config Collector emits Credential nodes with
// objectid derived from (config-path, env-var-name); the LiteLLM Looter
// emits Credential nodes with objectid derived from (litellm-endpoint,
// key-name). Those two objectid derivations are completely different
// — but if both collectors observed the same secret value, the
// resulting nodes carry the same value_hash, and the
// cross_service_credential_chain post-processor can join them.
//
// Always populated, regardless of --include-credential-values. The raw
// Credential.value is omitted in the default safety mode; value_hash
// stays so the chain still fires.
//
// One swap point if the algorithm needs to change later (e.g. salted
// HMAC for cross-deployment privacy).
func HashCredentialValue(value string) string {
	return HashSHA256(value)
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
