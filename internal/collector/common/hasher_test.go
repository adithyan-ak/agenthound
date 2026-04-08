package common

import (
	"testing"
)

func TestHashSHA256(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:  "hello",
			input: "hello",
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "unicode",
			input: "\u00e9\u00e8\u00ea",
			want:  HashSHA256("\u00e9\u00e8\u00ea"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashSHA256(tt.input)
			if got != tt.want {
				t.Errorf("HashSHA256(%q) = %q, want %q", tt.input, got, tt.want)
			}
			if len(got) != 64 {
				t.Errorf("HashSHA256(%q) length = %d, want 64", tt.input, len(got))
			}
		})
	}
}

func TestCanonicalJSONHash(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:  "simple map",
			input: map[string]string{"b": "2", "a": "1"},
		},
		{
			name:  "nested structure",
			input: map[string]any{"name": "test", "props": map[string]int{"x": 1}},
		},
		{
			name:  "nil input",
			input: nil,
		},
		{
			name:  "empty string",
			input: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CanonicalJSONHash(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CanonicalJSONHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(got) != 64 {
				t.Errorf("CanonicalJSONHash() length = %d, want 64", len(got))
			}
		})
	}

	t.Run("sorted keys determinism", func(t *testing.T) {
		a := map[string]string{"z": "1", "a": "2", "m": "3"}
		b := map[string]string{"a": "2", "m": "3", "z": "1"}
		ha, _ := CanonicalJSONHash(a)
		hb, _ := CanonicalJSONHash(b)
		if ha != hb {
			t.Errorf("same map different order produced different hashes: %s vs %s", ha, hb)
		}
	})

	t.Run("unmarshallable input", func(t *testing.T) {
		_, err := CanonicalJSONHash(make(chan int))
		if err == nil {
			t.Error("expected error for channel input")
		}
	})
}

func TestDescriptionHash(t *testing.T) {
	t.Run("with schema", func(t *testing.T) {
		schema := map[string]any{"type": "object"}
		h1 := DescriptionHash("tool1", "does stuff", schema)
		h2 := DescriptionHash("tool1", "does stuff", schema)
		if h1 != h2 {
			t.Errorf("same inputs produced different hashes")
		}
		if len(h1) != 64 {
			t.Errorf("hash length = %d, want 64", len(h1))
		}
	})

	t.Run("nil schema omitted", func(t *testing.T) {
		withNil := DescriptionHash("tool1", "does stuff", nil)
		withSchema := DescriptionHash("tool1", "does stuff", map[string]any{"type": "object"})
		if withNil == withSchema {
			t.Error("nil schema should produce different hash than non-nil schema")
		}
	})

	t.Run("different names different hashes", func(t *testing.T) {
		h1 := DescriptionHash("tool1", "desc", nil)
		h2 := DescriptionHash("tool2", "desc", nil)
		if h1 == h2 {
			t.Error("different names should produce different hashes")
		}
	})

	t.Run("empty strings", func(t *testing.T) {
		h := DescriptionHash("", "", nil)
		if len(h) != 64 {
			t.Errorf("hash length = %d, want 64", len(h))
		}
	})
}
