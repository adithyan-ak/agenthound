package processors

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/adithyan-ak/agenthound/server/internal/graph"
)

func TestShadows_Name(t *testing.T) {
	p := &Shadows{}
	if p.Name() != "shadows" {
		t.Errorf("Name() = %q, want shadows", p.Name())
	}
}

func TestShadows_Dependencies(t *testing.T) {
	p := &Shadows{}
	if deps := p.Dependencies(); deps != nil {
		t.Errorf("Dependencies() = %v, want nil", deps)
	}
}

func TestShadows_ProcessSuccess(t *testing.T) {
	mock := &graph.MockGraphDB{ExecuteWriteResult: 2}

	p := &Shadows{}
	stats, err := p.Process(context.Background(), mock, "scan-1")
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if stats.ProcessorName != "shadows" {
		t.Errorf("ProcessorName = %q", stats.ProcessorName)
	}
	// Two ExecuteWrite passes now run: SHADOWS, then POISONS_CONTEXT. The
	// mock returns 2 for each, so EdgesCreated sums to 4.
	if stats.EdgesCreated != 4 {
		t.Errorf("EdgesCreated = %d, want 4 (SHADOWS + POISONS_CONTEXT)", stats.EdgesCreated)
	}

	calls := mock.CallsTo("ExecuteWrite")
	if len(calls) != 2 {
		t.Fatalf("ExecuteWrite called %d times, want 2 (SHADOWS + POISONS_CONTEXT)", len(calls))
	}
	params, _ := calls[0].Args[1].(map[string]any)
	if params["scan_id"] != "scan-1" {
		t.Errorf("scan_id = %v, want scan-1", params["scan_id"])
	}

	shadowsCypher, _ := calls[0].Args[0].(string)
	if !contains(shadowsCypher, "SHADOWS") {
		t.Errorf("first pass should emit SHADOWS, got:\n%s", shadowsCypher)
	}
	poisonsCypher, _ := calls[1].Args[0].(string)
	if !contains(poisonsCypher, "POISONS_CONTEXT") {
		t.Errorf("second pass should emit POISONS_CONTEXT, got:\n%s", poisonsCypher)
	}
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }

func TestShadows_ProcessError(t *testing.T) {
	mock := &graph.MockGraphDB{ExecuteWriteError: errors.New("db error")}

	p := &Shadows{}
	_, err := p.Process(context.Background(), mock, "scan-1")
	if err == nil {
		t.Fatal("expected error")
	}
}
