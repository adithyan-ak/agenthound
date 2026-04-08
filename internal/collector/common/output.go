package common

import (
	"fmt"
	"time"

	"github.com/adithyan-ak/agenthound/internal/model"
)

const CollectorVersion = "0.1.0"

func NewIngestData(collector, scanID string) *model.IngestData {
	if scanID == "" {
		scanID = GenerateScanID(collector)
	}
	return &model.IngestData{
		Meta: model.IngestMeta{
			Version:          1,
			Type:             "agenthound-ingest",
			Collector:        collector,
			CollectorVersion: CollectorVersion,
			Timestamp:        time.Now().UTC().Format(time.RFC3339),
			ScanID:           scanID,
		},
		Graph: model.GraphData{
			Nodes: []model.Node{},
			Edges: []model.Edge{},
		},
	}
}

func GenerateScanID(collector string) string {
	return fmt.Sprintf("scan-%s-%d", collector, time.Now().UnixMilli())
}

func NewEdgeProps(scanID string, confidence, riskWeight float64) map[string]any {
	return map[string]any{
		"scan_id":      scanID,
		"last_seen":    time.Now().UTC().Format(time.RFC3339),
		"confidence":   confidence,
		"risk_weight":  riskWeight,
		"is_composite": false,
	}
}

func DefaultEdgeProps(scanID string) map[string]any {
	return NewEdgeProps(scanID, 1.0, 0.0)
}

func NewNode(id string, kinds []string, props map[string]any) model.Node {
	if props == nil {
		props = make(map[string]any)
	}
	props["objectid"] = id
	return model.Node{
		ID:         id,
		Kinds:      kinds,
		Properties: props,
	}
}

func NewEdge(source, target, kind string, props map[string]any) model.Edge {
	if props == nil {
		props = make(map[string]any)
	}
	return model.Edge{
		Source:     source,
		Target:     target,
		Kind:       kind,
		Properties: props,
	}
}
