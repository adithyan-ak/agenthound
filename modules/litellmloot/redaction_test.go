package litellmloot

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adithyan-ak/agenthound/sdk/action"
)

// TestLoot_MasterKeyNeverInLogs is a hard regression guard: the FULL
// master key MUST NEVER appear in slog output across any code path of
// the looter. The redact() helper trims to an 8-char prefix; this test
// captures every slog write into a buffer and greps for the full
// secret.
//
// Per docs/plans/sprint3-offensive-primitives.md 9.5 — operators
// running the looter against a real engagement will pipe slog output
// to a log file that may end up in a SIEM. A leaked master key in a
// log line is the same as a leaked master key in a network capture.
func TestLoot_MasterKeyNeverInLogs(t *testing.T) {
	const sneakyMasterKey = "sk-NEVER-LOG-THIS-VALUE-NOT-REAL-XYZ"

	var buf bytes.Buffer
	captureHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	originalLogger := slog.Default()
	slog.SetDefault(slog.New(captureHandler))
	defer slog.SetDefault(originalLogger)

	// Stub server that fails BOTH probes so we exercise the slog.Warn
	// paths in the looter's partial-failure branches — the most likely
	// place a careless implementer might log the key.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer srv.Close()

	l := &Looter{}
	_, _ = l.Loot(context.Background(), action.Target{
		Address: strings.TrimPrefix(srv.URL, "http://"),
	}, action.LootOptions{
		Credentials:  map[string]string{"master_key": sneakyMasterKey},
		EngagementID: "REDACTION-TEST",
	})

	logs := buf.String()
	if strings.Contains(logs, sneakyMasterKey) {
		t.Errorf("FULL master key leaked into slog output:\n%s", logs)
	}
	// Sanity: the redacted prefix should appear, proving redact() ran.
	wantPrefix := sneakyMasterKey[:8] + "..."
	if !strings.Contains(logs, wantPrefix) {
		t.Errorf("expected redacted prefix %q in logs; got:\n%s", wantPrefix, logs)
	}
}

func TestRedact(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"sk-1234567890abcdef", "sk-12345..."},
		{"sk-ABC", "***"},   // 6 chars ≤ 8 → fully redacted
		{"", "***"},         // empty → fully redacted
		{"12345678", "***"}, // exactly 8 → fully redacted
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := redact(tt.input); got != tt.want {
				t.Errorf("redact(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
