package action

import "context"

// Extractor pulls a specific resource by reference (a known file path,
// memory region, table, etc.) from a Target. Distinct from Looter, which
// performs broader untargeted collection.
//
// Implementations also implement sdk/module.Module.
type Extractor interface {
	Extract(ctx context.Context, t Target, opts ExtractOptions) (*ExtractResult, error)
}

// ExtractOptions is a v0 stub.
type ExtractOptions struct{}

// ExtractResult is a v0 stub.
type ExtractResult struct{}
