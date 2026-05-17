package embeddinginvert

import (
	"context"
	"testing"

	"github.com/adithyan-ak/agenthound/sdk/action"
)

func TestExtract_DetectsOutliers(t *testing.T) {
	e := &Extractor{}
	res, err := e.Extract(context.Background(), action.Target{
		Kind:    "node",
		Address: "sha256:test-model-id",
	}, action.ExtractOptions{
		SourceNodeID: "sha256:test-model-id",
		ArtifactPath: fixturePath(),
		EngagementID: "TEST-001",
		DryRun:       false,
		Extras:       map[string]any{"confidence-threshold": 1.5},
	})
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	if res.IngestData == nil {
		t.Fatal("IngestData nil")
	}
	if res.Summary.ArtifactsProduced < 2 {
		t.Errorf("expected at least 2 outliers (rows 8+9), got %d", res.Summary.ArtifactsProduced)
	}

	var foundSecret, foundTool bool
	for _, n := range res.IngestData.Graph.Nodes {
		tok, _ := n.Properties["token_string"].(string)
		switch tok {
		case "[fine_tune_secret]":
			foundSecret = true
		case "[internal_tool_xyz]":
			foundTool = true
		}
	}
	if !foundSecret {
		t.Error("outlier token [fine_tune_secret] not detected")
	}
	if !foundTool {
		t.Error("outlier token [internal_tool_xyz] not detected")
	}

	if len(res.IngestData.Graph.Edges) != res.Summary.ArtifactsProduced {
		t.Errorf("edges (%d) != artifacts (%d)", len(res.IngestData.Graph.Edges), res.Summary.ArtifactsProduced)
	}
	for _, e := range res.IngestData.Graph.Edges {
		if e.Kind != "EXTRACTED_FROM" {
			t.Errorf("edge kind = %q, want EXTRACTED_FROM", e.Kind)
		}
		if e.SourceKind != "AIModel" || e.TargetKind != "ExtractedTrainingSignal" {
			t.Errorf("edge endpoints: %s -> %s", e.SourceKind, e.TargetKind)
		}
	}
}

func TestExtract_DryRunEmitsNoData(t *testing.T) {
	e := &Extractor{}
	res, err := e.Extract(context.Background(), action.Target{
		Kind:    "node",
		Address: "sha256:test-model-id",
	}, action.ExtractOptions{
		SourceNodeID: "sha256:test-model-id",
		ArtifactPath: fixturePath(),
		EngagementID: "TEST-002",
		DryRun:       true,
		Extras:       map[string]any{"confidence-threshold": 1.5},
	})
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	if res.IngestData != nil {
		t.Error("DryRun should not produce IngestData")
	}
	if !res.Summary.DryRun {
		t.Error("Summary.DryRun should be true")
	}
	if res.Summary.ArtifactsProduced < 2 {
		t.Errorf("dry-run should still count outliers: got %d", res.Summary.ArtifactsProduced)
	}
}

func TestExtract_MissingArtifact(t *testing.T) {
	e := &Extractor{}
	_, err := e.Extract(context.Background(), action.Target{}, action.ExtractOptions{
		ArtifactPath: "/nonexistent.gguf",
		EngagementID: "X",
	})
	if err == nil {
		t.Error("expected error on missing artifact")
	}
}

func TestExtract_RequiresArtifactPath(t *testing.T) {
	e := &Extractor{}
	_, err := e.Extract(context.Background(), action.Target{}, action.ExtractOptions{
		EngagementID: "X",
	})
	if err == nil {
		t.Error("expected error when --artifact not provided")
	}
}
