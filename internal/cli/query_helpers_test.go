package cli

import (
	"sort"
	"testing"
)

func TestParseNodeRef(t *testing.T) {
	tests := []struct {
		name    string
		ref     string
		kind    string
		nname   string
		wantErr string
	}{
		{"valid MCPServer", "MCPServer:my-server", "MCPServer", "my-server", ""},
		{"valid A2AAgent", "A2AAgent:agent1", "A2AAgent", "agent1", ""},
		{"valid with colon in name", "MCPResource:postgres://prod:5432", "MCPResource", "postgres://prod:5432", ""},
		{"unknown kind", "FakeKind:x", "", "", "unknown node kind"},
		{"no colon", "MCPServer", "", "", "invalid format"},
		{"empty name", "MCPServer:", "", "", "invalid format"},
		{"empty kind", ":name", "", "", "invalid format"},
		{"just colon", ":", "", "", "invalid format"},
		{"empty string", "", "", "", "invalid format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, name, err := parseNodeRef(tt.ref)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !containsSubstring(err.Error(), tt.wantErr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if kind != tt.kind {
				t.Errorf("kind = %q, want %q", kind, tt.kind)
			}
			if name != tt.nname {
				t.Errorf("name = %q, want %q", name, tt.nname)
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{"nil", nil, "<null>"},
		{"string", "hello", "hello"},
		{"integer float", float64(42.0), "42"},
		{"fractional float", float64(3.14159), "3.1416"},
		{"int64", int64(99), "99"},
		{"true", true, "true"},
		{"false", false, "false"},
		{"slice", []any{"a", "b"}, "[a, b]"},
		{"empty slice", []any{}, "[]"},
		{"nested nil in slice", []any{nil, "x"}, "[<null>, x]"},
		{"default type", struct{ X int }{42}, "{42}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatValue(tt.val)
			if got != tt.want {
				t.Errorf("formatValue(%v) = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

func TestOrderedColumns(t *testing.T) {
	t.Run("returns all keys", func(t *testing.T) {
		row := map[string]any{"name": "x", "kind": "MCPServer", "score": 42}
		cols := orderedColumns(row)

		if len(cols) != len(row) {
			t.Fatalf("len = %d, want %d", len(cols), len(row))
		}

		sort.Strings(cols)
		want := []string{"kind", "name", "score"}
		for i, c := range cols {
			if c != want[i] {
				t.Errorf("cols[%d] = %q, want %q", i, c, want[i])
			}
		}
	})

	t.Run("empty map", func(t *testing.T) {
		cols := orderedColumns(map[string]any{})
		if len(cols) != 0 {
			t.Fatalf("expected empty slice, got %v", cols)
		}
	})
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
