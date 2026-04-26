package common

import (
	"math"
	"strings"
	"testing"
)

func TestShannonEntropy(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantMin float64
		wantMax float64
	}{
		{
			name:    "empty string",
			input:   "",
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "single char",
			input:   "a",
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "repeated char",
			input:   "aaaaaaa",
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "two equal chars",
			input:   "ab",
			wantMin: 0.99,
			wantMax: 1.01,
		},
		{
			name:    "high entropy random-looking",
			input:   "aB3$xY7!mN9@pQ2&",
			wantMin: 4.0,
			wantMax: 8.0,
		},
		{
			name:    "low entropy repeated pattern",
			input:   "abababababababab",
			wantMin: 0.99,
			wantMax: 1.01,
		},
		{
			name:    "all 256 byte values produce max entropy",
			input:   allBytes(),
			wantMin: 7.99,
			wantMax: 8.01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShannonEntropy(tt.input)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("ShannonEntropy(%q) = %f, want [%f, %f]", tt.input, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func allBytes() string {
	var b [256]byte
	for i := range b {
		b[i] = byte(i)
	}
	return string(b[:])
}

func TestShannonEntropyNonNegative(t *testing.T) {
	inputs := []string{"", "a", "ab", "abc", "password123", strings.Repeat("x", 1000)}
	for _, s := range inputs {
		e := ShannonEntropy(s)
		if e < 0 || math.IsNaN(e) {
			t.Errorf("ShannonEntropy(%q) = %f, should be non-negative", s, e)
		}
	}
}

func TestIsBase64Charset(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", false},
		{"valid base64", "ABCDabcd0123+/==", true},
		{"with space", "ABC DEF", false},
		{"with hyphen", "abc-def", false},
		{"pure alpha", "abcXYZ", true},
		{"pure digits", "0123456789", true},
		{"padding only", "====", true},
		{"unicode", "éè", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBase64Charset(tt.input)
			if got != tt.want {
				t.Errorf("IsBase64Charset(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsHexCharset(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", false},
		{"valid hex lower", "0123456789abcdef", true},
		{"valid hex upper", "0123456789ABCDEF", true},
		{"valid hex mixed", "aAbBcC00FF", true},
		{"with g", "abcdefg", false},
		{"with space", "ab cd", false},
		{"unicode", "é", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsHexCharset(tt.input)
			if got != tt.want {
				t.Errorf("IsHexCharset(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsLikelySecret(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", false},
		{"short", "abc", false},
		{"7 chars", "abcdefg", false},
		{"low entropy base64", "aaaaaaaa", false},
		{"high entropy base64", "aB3xY7mN9pQ2kL5wRjT8vU6hZeF4gI0s", true},
		{"high entropy hex", "4a8f2c1d9e3b7a0f5c6d8e2b1a4f7c3d", true},
		{"low entropy hex", "00000000", false},
		{"not base64 or hex", "hello world!!", false},
		{"real-looking api key", "sk7KjR3mNpQ2xY8bL4wT9vU6hA1cE5fG", true},
		{"real-looking hex token", "9f3a7b2e8d1c4f06a5e7b9d2c8f1a3e6", true},
		{"simple password", "password", false},
		{"boundary hex at 3.0 not secret", "a1b2c3d4", false},
		{"boundary base64 at 4.0 not secret", "aB3xY7mN9pQ2kL5w", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsLikelySecret(tt.input)
			if got != tt.want {
				t.Errorf("IsLikelySecret(%q) = %v, want %v (entropy=%f, base64=%v, hex=%v)",
					tt.input, got, tt.want,
					ShannonEntropy(tt.input), IsBase64Charset(tt.input), IsHexCharset(tt.input))
			}
		})
	}
}
