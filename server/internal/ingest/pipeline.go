package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	sdkingest "github.com/adithyan-ak/agenthound/sdk/ingest"
	"github.com/adithyan-ak/agenthound/server/internal/analysis"
	"github.com/adithyan-ak/agenthound/server/internal/appdb"
	"github.com/adithyan-ak/agenthound/server/internal/graph"
	"github.com/adithyan-ak/agenthound/server/model"
)

// nodeEdgeWriter is the subset of *graph.Writer that Pipeline needs.
// Defining it as an interface lets tests substitute a recording mock.
// *graph.Writer satisfies it implicitly.
type nodeEdgeWriter interface {
	WriteNodes(ctx context.Context, nodes []sdkingest.Node, scanID string) (int, error)
	WriteEdges(ctx context.Context, edges []sdkingest.Edge, scanID string) (int, error)
}

// scanRecorder is the subset of *appdb.ScanStore that Pipeline needs.
// *appdb.ScanStore satisfies it implicitly.
type scanRecorder interface {
	CreateScan(ctx context.Context, scan *model.Scan) error
	UpdateScan(ctx context.Context, id, status string, nodeCount, edgeCount int, scanErr string) error
}

// postProcessFunc runs the analysis post-processors. Defaulted to
// analysis.RunPostProcessors; replaceable in tests to assert behavior
// without exercising every concrete processor.
type postProcessFunc func(ctx context.Context, db graph.GraphDB, scanID string, collectors []string) ([]graph.ProcessingStats, error)

// Pipeline serializes ingests through a single mutex.
//
// The post-processing stage's stale-edge cleanup deletes composite
// edges where scan_id != current AND source_collector IN current. Two
// concurrent mcp ingests would both run that DELETE, each treating the
// other's freshly-written edges as stale, producing edge-flapping or
// duplicate work. For a single-user server, serializing ingests is the
// correct trade-off: ingest is rare (operator-driven), the lock holds
// for the duration of a single scan, and concurrent UI uploads are not
// a target scenario.
type Pipeline struct {
	mu         sync.Mutex
	validator  *Validator
	normalizer *Normalizer
	writer     nodeEdgeWriter
	graphDB    graph.GraphDB
	scanStore  scanRecorder
	runPP      postProcessFunc
}

func NewPipeline(writer *graph.Writer, graphDB graph.GraphDB, scanStore *appdb.ScanStore) *Pipeline {
	p := &Pipeline{
		validator:  NewValidator(),
		normalizer: NewNormalizer(),
		graphDB:    graphDB,
		runPP:      analysis.RunPostProcessors,
	}
	// Avoid the typed-nil-into-interface trap: assign only if the
	// concrete pointer is non-nil so `p.writer != nil` checks behave
	// the same as before the interface seam was introduced.
	if writer != nil {
		p.writer = writer
	}
	if scanStore != nil {
		p.scanStore = scanStore
	}
	return p
}

func (p *Pipeline) Ingest(ctx context.Context, data *sdkingest.IngestData) (*sdkingest.IngestResult, error) {
	// Serialize ingests; see Pipeline doc-comment for why.
	p.mu.Lock()
	defer p.mu.Unlock()

	start := time.Now()
	result := &sdkingest.IngestResult{
		ScanID: data.Meta.ScanID,
	}

	// Stage 1: Validate
	if err := p.validator.Validate(data); err != nil {
		return nil, err
	}
	slog.Info("validation passed", "nodes", len(data.Graph.Nodes), "edges", len(data.Graph.Edges))

	// Stage 2: Normalize
	warnings := p.normalizer.Normalize(data)
	result.Warnings = warnings
	if len(warnings) > 0 {
		slog.Info("normalization warnings", "count", len(warnings))
	}

	// Stage 3: Record scan start
	if p.scanStore != nil {
		scan := &model.Scan{
			ID:        data.Meta.ScanID,
			Collector: data.Meta.Collector,
			Status:    model.ScanStatusRunning,
			StartedAt: time.Now().UTC(),
		}
		if err := p.scanStore.CreateScan(ctx, scan); err != nil {
			slog.Warn("failed to create scan record", "error", err)
		}
	}

	// Stage 4: Write nodes
	nodesWritten, err := p.writer.WriteNodes(ctx, data.Graph.Nodes, data.Meta.ScanID)
	if err != nil {
		p.failScan(ctx, data.Meta.ScanID, err)
		return nil, fmt.Errorf("write nodes: %w", err)
	}
	result.NodesWritten = nodesWritten
	slog.Info("nodes written", "count", nodesWritten)

	// Stage 5: Write edges
	edgesWritten, err := p.writer.WriteEdges(ctx, data.Graph.Edges, data.Meta.ScanID)
	if err != nil {
		p.failScan(ctx, data.Meta.ScanID, err)
		return nil, fmt.Errorf("write edges: %w", err)
	}
	result.EdgesWritten = edgesWritten
	slog.Info("edges written", "count", edgesWritten)

	// Stage 6: Post-processing (non-fatal)
	var ppErr error
	if p.graphDB != nil && p.runPP != nil {
		var ppStats []graph.ProcessingStats
		ppStats, ppErr = p.runPP(ctx, p.graphDB, data.Meta.ScanID, []string{data.Meta.Collector})
		if ppErr != nil {
			slog.Error("post-processing failed", "error", ppErr)
		}
		for _, s := range ppStats {
			result.PostProcessingStats = append(result.PostProcessingStats, sdkingest.PostProcessingStat{
				ProcessorName: s.ProcessorName,
				EdgesCreated:  s.EdgesCreated,
				NodesUpdated:  s.NodesUpdated,
				Duration:      s.Duration,
				Error:         s.Error,
			})
		}
	}

	// Stage 7: Record completion
	if p.scanStore != nil {
		status := model.ScanStatusCompleted
		scanErr := ""
		if ppErr != nil {
			// Collection (node/edge writes) succeeded; only analysis
			// failed. Record the real counts plus the error under a
			// dedicated status so this is distinguishable from a hard
			// collection failure (failScan -> ScanStatusFailed, 0/0).
			status = model.ScanStatusCompletedWithErrors
			scanErr = fmt.Sprintf("post-processing: %v", ppErr)
		}
		if err := p.scanStore.UpdateScan(ctx, data.Meta.ScanID, status, nodesWritten, edgesWritten, scanErr); err != nil {
			slog.Warn("failed to update scan record", "error", err)
		}
	}

	result.Duration = time.Since(start)
	slog.Info("ingest complete", "scan_id", data.Meta.ScanID, "nodes", nodesWritten, "edges", edgesWritten, "duration", result.Duration)
	return result, nil
}

func (p *Pipeline) failScan(ctx context.Context, scanID string, scanErr error) {
	if p.scanStore != nil {
		if err := p.scanStore.UpdateScan(ctx, scanID, model.ScanStatusFailed, 0, 0, scanErr.Error()); err != nil {
			slog.Warn("failed to record scan failure", "error", err)
		}
	}
}
