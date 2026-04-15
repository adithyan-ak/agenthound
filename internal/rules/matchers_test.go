package rules

import (
	"strings"
	"testing"
)

func TestRegexMatcher(t *testing.T) {
	tests := []struct {
		name    string
		spec    MatcherSpec
		input   string
		wantN   int
		wantErr bool
	}{
		{
			name:  "simple match",
			spec:  MatcherSpec{Type: "regex", Pattern: `</?IMPORTANT>`},
			input: "text <IMPORTANT>override</IMPORTANT> end",
			wantN: 2,
		},
		{
			name:  "case insensitive",
			spec:  MatcherSpec{Type: "regex", Pattern: `<important>`, CaseInsensitive: true},
			input: "<IMPORTANT>",
			wantN: 1,
		},
		{
			name:  "no match",
			spec:  MatcherSpec{Type: "regex", Pattern: `foobar`},
			input: "clean text",
			wantN: 0,
		},
		{
			name:    "invalid regex",
			spec:    MatcherSpec{Type: "regex", Pattern: `[invalid`},
			wantErr: true,
		},
		{
			name:  "already has case flag",
			spec:  MatcherSpec{Type: "regex", Pattern: `(?i)test`, CaseInsensitive: true},
			input: "TEST",
			wantN: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := compileMatcher(tc.spec)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			results := m.Match(tc.input)
			if len(results) != tc.wantN {
				t.Errorf("got %d results, want %d", len(results), tc.wantN)
			}
		})
	}
}

func TestRegexMatcherTextTruncation(t *testing.T) {
	spec := MatcherSpec{Type: "regex", Pattern: `A+`}
	m, err := compileMatcher(spec)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	input := strings.Repeat("A", 200)
	results := m.Match(input)
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if len(results[0].Text) != 100 {
		t.Errorf("text length = %d, want 100", len(results[0].Text))
	}
}

func TestKeywordMatcher(t *testing.T) {
	tests := []struct {
		name  string
		spec  MatcherSpec
		input string
		wantN int
	}{
		{
			name:  "any mode first match",
			spec:  MatcherSpec{Type: "keyword", Keywords: []string{"shell", "bash"}, CaseInsensitive: true, MatchMode: "any"},
			input: "run shell bash command",
			wantN: 1,
		},
		{
			name:  "all mode all present",
			spec:  MatcherSpec{Type: "keyword", Keywords: []string{"shell", "bash"}, CaseInsensitive: true, MatchMode: "all"},
			input: "shell and bash",
			wantN: 2,
		},
		{
			name:  "all mode missing one",
			spec:  MatcherSpec{Type: "keyword", Keywords: []string{"shell", "bash"}, CaseInsensitive: true, MatchMode: "all"},
			input: "only shell here",
			wantN: 0,
		},
		{
			name:  "case sensitive no match",
			spec:  MatcherSpec{Type: "keyword", Keywords: []string{"Shell"}, CaseInsensitive: false},
			input: "shell command",
			wantN: 0,
		},
		{
			name:  "case sensitive match",
			spec:  MatcherSpec{Type: "keyword", Keywords: []string{"Shell"}, CaseInsensitive: false},
			input: "Shell command",
			wantN: 1,
		},
		{
			name:  "no keywords no match",
			spec:  MatcherSpec{Type: "keyword", Keywords: []string{}, CaseInsensitive: true},
			input: "anything",
			wantN: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := compileMatcher(tc.spec)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			results := m.Match(tc.input)
			if len(results) != tc.wantN {
				t.Errorf("got %d results, want %d", len(results), tc.wantN)
			}
		})
	}
}

func TestCompoundMatcher(t *testing.T) {
	tests := []struct {
		name  string
		spec  MatcherSpec
		input string
		wantN int
	}{
		{
			name: "and both match",
			spec: MatcherSpec{
				Type:     "compound",
				Operator: "and",
				Matchers: []MatcherSpec{
					{Type: "keyword", Keywords: []string{"postgres"}, CaseInsensitive: true},
					{Type: "keyword", Keywords: []string{"prod"}, CaseInsensitive: true},
				},
			},
			input: "postgres://prod-db",
			wantN: 2,
		},
		{
			name: "and one missing",
			spec: MatcherSpec{
				Type:     "compound",
				Operator: "and",
				Matchers: []MatcherSpec{
					{Type: "keyword", Keywords: []string{"postgres"}, CaseInsensitive: true},
					{Type: "keyword", Keywords: []string{"prod"}, CaseInsensitive: true},
				},
			},
			input: "postgres://dev-db",
			wantN: 0,
		},
		{
			name: "or first matches",
			spec: MatcherSpec{
				Type:     "compound",
				Operator: "or",
				Matchers: []MatcherSpec{
					{Type: "keyword", Keywords: []string{"alpha"}, CaseInsensitive: true},
					{Type: "keyword", Keywords: []string{"beta"}, CaseInsensitive: true},
				},
			},
			input: "alpha version",
			wantN: 1,
		},
		{
			name: "or none matches",
			spec: MatcherSpec{
				Type:     "compound",
				Operator: "or",
				Matchers: []MatcherSpec{
					{Type: "keyword", Keywords: []string{"alpha"}, CaseInsensitive: true},
					{Type: "keyword", Keywords: []string{"beta"}, CaseInsensitive: true},
				},
			},
			input: "gamma version",
			wantN: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := compileMatcher(tc.spec)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			results := m.Match(tc.input)
			if len(results) != tc.wantN {
				t.Errorf("got %d results, want %d", len(results), tc.wantN)
			}
		})
	}
}

func TestEntropyMatcher(t *testing.T) {
	tests := []struct {
		name  string
		spec  MatcherSpec
		input string
		wantN int
	}{
		{
			name:  "high entropy base64",
			spec:  MatcherSpec{Type: "entropy", Charset: "base64", Threshold: 4.5, MinLength: 8},
			input: "sk+ant+abc123XYZdefGHIjklMNOpqrSTUvwx",
			wantN: 1,
		},
		{
			name:  "low entropy base64",
			spec:  MatcherSpec{Type: "entropy", Charset: "base64", Threshold: 4.5, MinLength: 8},
			input: "aaaaaaaa",
			wantN: 0,
		},
		{
			name:  "too short",
			spec:  MatcherSpec{Type: "entropy", Charset: "base64", Threshold: 4.5, MinLength: 8},
			input: "abc",
			wantN: 0,
		},
		{
			name:  "non base64 chars rejected",
			spec:  MatcherSpec{Type: "entropy", Charset: "base64", Threshold: 4.5, MinLength: 8},
			input: "hello world with spaces",
			wantN: 0,
		},
		{
			name:  "high entropy hex",
			spec:  MatcherSpec{Type: "entropy", Charset: "hex", Threshold: 3.0, MinLength: 8},
			input: "a1b2c3d4e5f6a7b8",
			wantN: 1,
		},
		{
			name:  "hex with non-hex chars rejected",
			spec:  MatcherSpec{Type: "entropy", Charset: "hex", Threshold: 3.0, MinLength: 8},
			input: "a1b2c3g4",
			wantN: 0,
		},
		{
			name:  "empty string",
			spec:  MatcherSpec{Type: "entropy", Charset: "base64", Threshold: 4.5, MinLength: 8},
			input: "",
			wantN: 0,
		},
		{
			name:  "base64 with padding",
			spec:  MatcherSpec{Type: "entropy", Charset: "base64", Threshold: 3.0, MinLength: 8},
			input: "dGVzdDE2Y2hhcg==",
			wantN: 1,
		},
		{
			name:  "hex charset rejects base64 plus sign",
			spec:  MatcherSpec{Type: "entropy", Charset: "hex", Threshold: 3.0, MinLength: 8},
			input: "a1b2c3d4+5f6",
			wantN: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := compileMatcher(tc.spec)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			results := m.Match(tc.input)
			if len(results) != tc.wantN {
				t.Errorf("got %d results, want %d", len(results), tc.wantN)
			}
		})
	}
}

func TestPrefixMatcher(t *testing.T) {
	tests := []struct {
		name  string
		spec  MatcherSpec
		input string
		wantN int
	}{
		{
			name:  "matching prefix",
			spec:  MatcherSpec{Type: "prefix", Prefixes: []string{"sk-ant-"}},
			input: "sk-ant-abc123",
			wantN: 1,
		},
		{
			name:  "no matching prefix",
			spec:  MatcherSpec{Type: "prefix", Prefixes: []string{"sk-ant-"}},
			input: "sk-abc123",
			wantN: 0,
		},
		{
			name:  "case insensitive",
			spec:  MatcherSpec{Type: "prefix", Prefixes: []string{"vault://"}, CaseInsensitive: true},
			input: "VAULT://secret/path",
			wantN: 1,
		},
		{
			name:  "case sensitive no match",
			spec:  MatcherSpec{Type: "prefix", Prefixes: []string{"sk-ant-"}, CaseInsensitive: false},
			input: "SK-ANT-abc123",
			wantN: 0,
		},
		{
			name:  "multiple prefixes first wins",
			spec:  MatcherSpec{Type: "prefix", Prefixes: []string{"ghp_", "gho_", "sk-"}},
			input: "gho_abcdef",
			wantN: 1,
		},
		{
			name:  "empty input",
			spec:  MatcherSpec{Type: "prefix", Prefixes: []string{"sk-"}},
			input: "",
			wantN: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := compileMatcher(tc.spec)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			results := m.Match(tc.input)
			if len(results) != tc.wantN {
				t.Errorf("got %d results, want %d", len(results), tc.wantN)
			}
			if tc.wantN > 0 && results[0].Offset != 0 {
				t.Errorf("prefix offset = %d, want 0", results[0].Offset)
			}
		})
	}
}

func TestUnknownMatcherType(t *testing.T) {
	_, err := compileMatcher(MatcherSpec{Type: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown matcher type")
	}
}
