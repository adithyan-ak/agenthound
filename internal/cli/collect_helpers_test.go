package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestWriteCollectorOutput_File(t *testing.T) {
	data := &model.IngestData{
		Meta: model.IngestMeta{
			Version:   1,
			Type:      "agenthound-ingest",
			Collector: "test",
		},
		Graph: model.GraphData{
			Nodes: []model.Node{
				{ID: "n1", Kinds: []string{"MCPServer"}, Properties: map[string]any{"name": "srv"}},
			},
		},
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.json")

	if err := writeCollectorOutput(data, outPath); err != nil {
		t.Fatalf("writeCollectorOutput: %v", err)
	}

	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	var got model.IngestData
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Meta.Collector != "test" {
		t.Errorf("collector = %q, want %q", got.Meta.Collector, "test")
	}
	if len(got.Graph.Nodes) != 1 {
		t.Errorf("nodes = %d, want 1", len(got.Graph.Nodes))
	}
}

func TestWriteCollectorOutput_Stdout(t *testing.T) {
	data := &model.IngestData{
		Meta: model.IngestMeta{
			Version:   1,
			Type:      "agenthound-ingest",
			Collector: "stdout-test",
		},
	}

	out := captureStdout(t, func() {
		if err := writeCollectorOutput(data, ""); err != nil {
			t.Fatalf("writeCollectorOutput: %v", err)
		}
	})

	var got model.IngestData
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nraw: %q", err, out)
	}
	if got.Meta.Collector != "stdout-test" {
		t.Errorf("collector = %q, want %q", got.Meta.Collector, "stdout-test")
	}
}
