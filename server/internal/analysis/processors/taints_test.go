package processors

import (
	"context"
	"errors"
	"testing"

	"github.com/adithyan-ak/agenthound/server/internal/graph"
)

func TestTaints_Name(t *testing.T) {
	p := &Taints{}
	if p.Name() != "taints" {
		t.Errorf("Name() = %q, want taints", p.Name())
	}
}

func TestTaints_Dependencies(t *testing.T) {
	p := &Taints{}
	if deps := p.Dependencies(); deps != nil {
		t.Errorf("Dependencies() = %v, want nil", deps)
	}
}

func TestTaints_ProcessSuccess(t *testing.T) {
	mock := &graph.MockGraphDB{ExecuteWriteResult: 2}

	p := &Taints{}
	stats, err := p.Process(context.Background(), mock, "scan-1")
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if stats.ProcessorName != "taints" {
		t.Errorf("ProcessorName = %q", stats.ProcessorName)
	}
	if stats.EdgesCreated != 2 {
		t.Errorf("EdgesCreated = %d, want 2", stats.EdgesCreated)
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
	for _, want := range []string{"TAINTS", ">= 2", "source_trust", "INGESTS_UNTRUSTED", "source_collector = 'mcp'"} {
		if !contains(cypher, want) {
			t.Errorf("Cypher missing load-bearing predicate %q; query:\n%s", want, cypher)
		}
	}
}

func TestTaints_ProcessError(t *testing.T) {
	mock := &graph.MockGraphDB{ExecuteWriteError: errors.New("db error")}

	p := &Taints{}
	_, err := p.Process(context.Background(), mock, "scan-1")
	if err == nil {
		t.Fatal("expected error")
	}
}
