// Package model is a temporary shim re-exporting symbols moved to sdk/ingest.
// Deleted in Step 4 after callers are mass-rewritten.
package model

import "github.com/adithyan-ak/agenthound/sdk/ingest"

// Type aliases (work for type names).
type Node = ingest.Node
type Edge = ingest.Edge
type EdgeEndpoints = ingest.EdgeEndpoints
type IngestData = ingest.IngestData
type IngestMeta = ingest.IngestMeta
type GraphData = ingest.GraphData
type IngestResult = ingest.IngestResult
type PostProcessingStat = ingest.PostProcessingStat

// Function values (Go has no `func X = pkg.X` syntax).
var (
	ComputeNodeID        = ingest.ComputeNodeID
	ComputeMCPServerID   = ingest.ComputeMCPServerID
	ResolveEdgeEndpoints = ingest.ResolveEdgeEndpoints
)

// Map and slice value re-exports — share the underlying data by reference.
// All these are immutable globals in practice (verified: no caller mutates them).
var (
	AllowedNodeKinds  = ingest.AllowedNodeKinds
	AllNodeLabels     = ingest.AllNodeLabels
	RawEdgeKinds      = ingest.RawEdgeKinds
	AllowedEdgeKinds  = ingest.AllowedEdgeKinds
	AllowedCollectors = ingest.AllowedCollectors
	EdgeKindEndpoints = ingest.EdgeKindEndpoints
)
