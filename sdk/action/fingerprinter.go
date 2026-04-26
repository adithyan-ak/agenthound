package action

import "context"

// Fingerprinter probes a single Target to identify the running service
// (kind, version, auth posture). Output feeds enumeration and module
// selection downstream.
//
// Implementations also implement sdk/module.Module.
type Fingerprinter interface {
	Fingerprint(ctx context.Context, t Target) (*FingerprintResult, error)
}

// FingerprintResult is a v0 stub. The concrete shape lands with the first
// Fingerprinter implementation; treat as opaque until then.
type FingerprintResult struct{}
