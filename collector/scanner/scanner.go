package scanner

import (
	"context"
	"errors"

	"github.com/adithyan-ak/agenthound/sdk/action"
)

// ErrNotImplemented is returned by stub Scan calls until the real implementation lands.
var ErrNotImplemented = errors.New("scanner: not yet implemented — see docs/future-modules.md")

// Stub is a placeholder Scanner for use until the real network discovery engine is built.
// It satisfies sdk/action.Scanner so call sites can compile against the contract today.
type Stub struct{}

// Scan returns ErrNotImplemented unconditionally.
func (Stub) Scan(_ context.Context, _ string) ([]action.Target, error) {
	return nil, ErrNotImplemented
}
