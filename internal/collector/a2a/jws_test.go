package a2a

import "testing"

func TestVerifySignatures_NoField(t *testing.T) {
	raw := map[string]any{"name": "test"}
	signed, valid := VerifySignatures(nil, raw)
	if signed || valid {
		t.Errorf("expected signed=false, valid=false for missing signatures; got signed=%v, valid=%v", signed, valid)
	}
}

func TestVerifySignatures_EmptyArray(t *testing.T) {
	raw := map[string]any{"signatures": []any{}}
	signed, valid := VerifySignatures(nil, raw)
	if signed || valid {
		t.Errorf("expected signed=false, valid=false for empty signatures; got signed=%v, valid=%v", signed, valid)
	}
}

func TestVerifySignatures_Present(t *testing.T) {
	raw := map[string]any{
		"signatures": []any{
			map[string]any{
				"protected": "eyJhbGciOiJSUzI1NiJ9",
				"signature": "dGVzdA",
			},
		},
	}
	signed, valid := VerifySignatures(nil, raw)
	if !signed {
		t.Error("expected signed=true when signatures present")
	}
	if valid {
		t.Error("expected valid=false in Phase 2 MVP")
	}
}

func TestVerifySignatures_WrongType(t *testing.T) {
	raw := map[string]any{"signatures": "not-an-array"}
	signed, valid := VerifySignatures(nil, raw)
	if signed || valid {
		t.Errorf("expected signed=false, valid=false for non-array signatures; got signed=%v, valid=%v", signed, valid)
	}
}
