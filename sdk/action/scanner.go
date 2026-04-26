package action

import "context"

// Scanner expands a CIDR / range / discovery seed into concrete Targets.
// Distinct from Fingerprinter, which probes a single Target.
//
// Implementations also implement sdk/module.Module.
type Scanner interface {
	Scan(ctx context.Context, cidr string) ([]Target, error)
}
