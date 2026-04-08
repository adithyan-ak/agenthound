package analysis

import (
	"context"
	"time"

	"github.com/adithyan-ak/agenthound/internal/graph"
)

type PostProcessor interface {
	Name() string
	Dependencies() []string
	Process(ctx context.Context, db graph.GraphDB, scanID string) (ProcessingStats, error)
}

type ProcessingStats struct {
	ProcessorName string        `json:"processor_name"`
	EdgesCreated  int           `json:"edges_created"`
	NodesUpdated  int           `json:"nodes_updated"`
	Duration      time.Duration `json:"duration"`
	Error         string        `json:"error,omitempty"`
}
