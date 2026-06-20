package common

import "testing"

func TestRedact(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"sk-1234567890abcdef", "sk-12345..."},
		{"123456789", "12345678..."}, // 9 chars (> 8) → prefix + ...
		{"sk-ABC", "***"},            // 6 chars ≤ 8 → fully redacted
		{"", "***"},                  // empty → fully redacted
		{"12345678", "***"},          // exactly 8 → fully redacted
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Redact(tt.input); got != tt.want {
				t.Errorf("Redact(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
