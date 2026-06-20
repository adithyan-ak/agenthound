package processors

import (
	"context"
	"errors"
	"testing"

	"github.com/adithyan-ak/agenthound/server/internal/graph"
)

func TestConfusedDeputy_Name(t *testing.T) {
	p := &ConfusedDeputy{}
	if p.Name() != "confused_deputy" {
		t.Errorf("Name() = %q, want confused_deputy", p.Name())
	}
}

func TestConfusedDeputy_Dependencies(t *testing.T) {
	p := &ConfusedDeputy{}
	deps := p.Dependencies()
	if len(deps) != 2 {
		t.Fatalf("Dependencies() = %v, want 2 entries", deps)
	}
	wantSet := map[string]bool{"auth_strength": true, "can_reach": true}
	for _, d := range deps {
		if !wantSet[d] {
			t.Errorf("unexpected dependency %q", d)
		}
		delete(wantSet, d)
	}
	if len(wantSet) > 0 {
		t.Errorf("missing dependencies: %v", wantSet)
	}
}

func TestConfusedDeputy_ProcessSuccess(t *testing.T) {
	mock := &graph.MockGraphDB{ExecuteWriteResult: 1}

	p := &ConfusedDeputy{}
	stats, err := p.Process(context.Background(), mock, "scan-1")
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if stats.ProcessorName != "confused_deputy" {
		t.Errorf("ProcessorName = %q", stats.ProcessorName)
	}
	if stats.EdgesCreated != 1 {
		t.Errorf("EdgesCreated = %d, want 1", stats.EdgesCreated)
	}

	calls := mock.CallsTo("ExecuteWrite")
	if len(calls) != 1 {
		t.Fatalf("ExecuteWrite called %d times, want 1", len(calls))
	}
	params, _ := calls[0].Args[1].(map[string]any)
	if params["scan_id"] != "scan-1" {
		t.Errorf("scan_id = %v, want scan-1", params["scan_id"])
	}

	cypher, _ := calls[0].Args[0].(string)
	// The auth-strength gradient literals are the boundary the detector
	// fires on: low (weak) >= 80, high (strong) <= 30. Assert the exact
	// predicate so a future threshold drift is caught.
	for _, want := range []string{"CONFUSED_DEPUTY", "low.auth_strength >= 80", "high.auth_strength <= 30", "source_collector = 'a2a'"} {
		if !contains(cypher, want) {
			t.Errorf("Cypher missing load-bearing predicate %q; query:\n%s", want, cypher)
		}
	}
}

func TestConfusedDeputy_ProcessError(t *testing.T) {
	mock := &graph.MockGraphDB{ExecuteWriteError: errors.New("db error")}

	p := &ConfusedDeputy{}
	_, err := p.Process(context.Background(), mock, "scan-1")
	if err == nil {
		t.Fatal("expected error")
	}
}
