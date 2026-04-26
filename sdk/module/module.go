// Package module declares the AgentHound module super-interface and a
// process-global registry. Modules self-register at init() time; the CLI and
// API discover them by ID, by Action, or by Target kind.
//
// Stability: v0 — UNSTABLE. The Module interface (in particular
// IsDestructive() bool) may be refined to a richer side-effect class enum
// at v1 if real UX gates emerge.
package module

import "github.com/adithyan-ak/agenthound/sdk/action"

// Module is the super-interface every registered module satisfies in
// addition to one or more action interfaces.
type Module interface {
	ID() string            // dotted lowercase, e.g. "mcp.enumerate"
	Action() action.Action // which action this module performs
	Target() string        // service kind targeted, e.g. "mcp", "a2a", "config"
	Description() string   // one-line human-readable summary
	Version() string       // semver of the module
	IsDestructive() bool   // v0 binary; refined at v1 if needed
}
