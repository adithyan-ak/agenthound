package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/ingest"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

// setupSentinels creates all sentinel files in a temp dir and sets HOME
// so the AUTHORIZED prompts are bypassed. Also resets rootCmd args.
func setupSentinels(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	ahDir := filepath.Join(home, ".agenthound")
	_ = os.MkdirAll(ahDir, 0o700)
	for _, name := range []string{"loot-acknowledged", "poison-acknowledged", "extract-acknowledged"} {
		_ = os.WriteFile(filepath.Join(ahDir, name), []byte(`{}`), 0o600)
	}
	// Reset rootCmd state between tests to avoid arg pollution.
	rootCmd.SetArgs(nil)
	rootCmd.SetOut(nil)
	rootCmd.SetErr(nil)
	return home
}

// --- Mock Looter ---

type mockLooter struct {
	called bool
	result *action.LootResult
	err    error
}

func (m *mockLooter) ID() string            { return "mock.loot" }
func (m *mockLooter) Action() action.Action { return action.Loot }
func (m *mockLooter) Target() string        { return "mock-svc" }
func (m *mockLooter) Description() string   { return "mock" }
func (m *mockLooter) Version() string       { return "0.0.0" }
func (m *mockLooter) IsDestructive() bool   { return false }
func (m *mockLooter) Loot(ctx context.Context, t action.Target, opts action.LootOptions) (*action.LootResult, error) {
	m.called = true
	if m.err != nil {
		return nil, m.err
	}
	if m.result != nil {
		return m.result, nil
	}
	return &action.LootResult{
		IngestData: &ingest.IngestData{},
		Summary:    action.LootSummary{CredentialsFound: 3},
	}, nil
}

// --- Mock Extractor ---

type mockExtractor struct {
	called bool
}

func (m *mockExtractor) ID() string            { return "mock.extract" }
func (m *mockExtractor) Action() action.Action { return action.Extract }
func (m *mockExtractor) Target() string        { return "mock-extract" }
func (m *mockExtractor) Description() string   { return "mock" }
func (m *mockExtractor) Version() string       { return "0.0.0" }
func (m *mockExtractor) IsDestructive() bool   { return false }
func (m *mockExtractor) Extract(ctx context.Context, t action.Target, opts action.ExtractOptions) (*action.ExtractResult, error) {
	m.called = true
	if opts.DryRun {
		return &action.ExtractResult{
			Summary: action.ExtractSummary{ArtifactsProduced: 5, DryRun: true},
		}, nil
	}
	return &action.ExtractResult{
		IngestData: &ingest.IngestData{
			Graph: ingest.GraphData{
				Nodes: []ingest.Node{{ID: "sha256:test", Kinds: []string{"ExtractedTrainingSignal"}}},
			},
		},
		Summary: action.ExtractSummary{ArtifactsProduced: 5},
	}, nil
}

// --- Loot tests ---

func TestRunLoot_HappyPath(t *testing.T) {
	setupSentinels(t)
	mock := &mockLooter{}
	module.Register(mock)
	defer deregisterModule(t, "mock.loot")

	out := &bytes.Buffer{}
	lootCmd.SetOut(out)
	lootCmd.SetErr(out)
	_ = lootCmd.Flags().Set("type", "mock-svc")
	_ = lootCmd.Flags().Set("engagement-id", "TEST")
	// Use root --output to stdout so writeCollectorOutput goes to buffer.
	_ = rootCmd.PersistentFlags().Set("output", "-")
	defer func() { _ = rootCmd.PersistentFlags().Set("output", "") }()
	err := runLoot(lootCmd, []string{"10.0.0.1:4000"})
	if err != nil {
		t.Fatalf("runLoot: %v", err)
	}
	if !mock.called {
		t.Error("mock Looter was not called")
	}
	if !strings.Contains(out.String(), "credentials_found=3") {
		t.Errorf("output missing credential count: %s", out.String())
	}
}

func TestRunLoot_NoModule(t *testing.T) {
	setupSentinels(t)
	out := &bytes.Buffer{}
	lootCmd.SetOut(out)
	lootCmd.SetErr(out)
	_ = lootCmd.Flags().Set("type", "definitely-not-registered")
	_ = lootCmd.Flags().Set("engagement-id", "X")
	err := runLoot(lootCmd, []string{"10.0.0.1:4000"})
	if err == nil || !strings.Contains(err.Error(), "no looter registered") {
		t.Errorf("expected 'no looter registered' error, got: %v", err)
	}
}

// --- Extract tests ---

func TestRunExtract_DryRun(t *testing.T) {
	setupSentinels(t)
	mock := &mockExtractor{}
	module.Register(mock)
	defer deregisterModule(t, "mock.extract")

	tmpArtifact := filepath.Join(t.TempDir(), "test.gguf")
	_ = os.WriteFile(tmpArtifact, []byte("fake"), 0o600)

	out := &bytes.Buffer{}
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{"extract", "sha256:node-id", "--type", "mock-extract",
		"--artifact", tmpArtifact, "--engagement-id", "TEST"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !mock.called {
		t.Error("mock Extractor was not called")
	}
	if !strings.Contains(out.String(), "DRY-RUN") {
		t.Errorf("expected DRY-RUN in output: %s", out.String())
	}
}

func TestRunExtract_NoModule(t *testing.T) {
	setupSentinels(t)
	out := &bytes.Buffer{}
	extractCmd.SetOut(out)
	extractCmd.SetErr(out)
	_ = extractCmd.Flags().Set("type", "not-registered-extractor")
	_ = extractCmd.Flags().Set("engagement-id", "TEST")
	_ = extractCmd.Flags().Set("artifact", "/tmp/fake.gguf")
	err := runExtract(extractCmd, []string{"sha256:node-id"})
	if err == nil || !strings.Contains(err.Error(), "no extractor registered") {
		t.Errorf("expected 'no extractor registered' error, got: %v", err)
	}
}

// deregisterModule removes a module from the registry after test.
// This is a test-only hack — the registry doesn't expose a Remove API,
// so we reach in via the package-internal access (same package).
func deregisterModule(t *testing.T, id string) {
	t.Helper()
	// Since we're in package cli (not package module), we can't directly
	// access the registry map. Instead we just accept that re-registration
	// in the same process will panic. Tests using mocks should use
	// unique IDs or subtests with separate processes.
	// For now, this is a no-op placeholder — the mock IDs don't collide
	// with real modules because they use "mock.*" prefix.
}
