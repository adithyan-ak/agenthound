package action

import (
	"context"

	"github.com/adithyan-ak/agenthound/sdk/ingest"
)

// Enumerator inspects a single Target and produces a graph patch.
// Supersedes today's Collector interface; the existing mcp / a2a / config
// collectors will adapt to this shape in a later step.
//
// Implementations also implement sdk/module.Module.
type Enumerator interface {
	Enumerate(ctx context.Context, t Target, opts EnumerateOptions) (*ingest.IngestData, error)
}

// EnumerateOptions is a v0 stub. Concrete fields land per-module
// (timeouts, transport overrides, redaction toggles) as needs surface.
type EnumerateOptions struct{}
