package networkscan

import (
	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

func init() {
	module.Register(&Scanner{})
}

// Scanner doubles as both the sdk/action.Scanner implementation AND the
// sdk/module.Module identity carrier — the registry indexes on the
// six identity methods below.

func (*Scanner) ID() string            { return "network.scan" }
func (*Scanner) Action() action.Action { return action.Scan }
func (*Scanner) Target() string        { return "network" }
func (*Scanner) Description() string {
	return "Scan a CIDR / host / file-of-targets for AI/ML services on standard ports"
}
func (*Scanner) Version() string     { return "0.2.0-dev" }
func (*Scanner) IsDestructive() bool { return false }
