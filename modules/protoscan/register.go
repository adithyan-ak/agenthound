package protoscan

import (
	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

// We register two modules: mcp.discover and a2a.discover. Both use the
// new action.Discover constant. The CLI dispatches by --type or --mcp /
// --a2a flags. Both modules wrap a *Scanner with the right Mode.

type mcpDiscoverer struct{ *Scanner }
type a2aDiscoverer struct{ *Scanner }

func init() {
	module.Register(&mcpDiscoverer{Scanner: &Scanner{Mode: ModeMCP}})
	module.Register(&a2aDiscoverer{Scanner: &Scanner{Mode: ModeA2A}})
}

func (*mcpDiscoverer) ID() string            { return "mcp.discover" }
func (*mcpDiscoverer) Action() action.Action { return action.Discover }
func (*mcpDiscoverer) Target() string        { return "mcp" }
func (*mcpDiscoverer) Description() string {
	return "Discover MCP servers on a network via JSON-RPC initialize probe"
}
func (*mcpDiscoverer) Version() string     { return "0.3.0-dev" }
func (*mcpDiscoverer) IsDestructive() bool { return false }

func (*a2aDiscoverer) ID() string            { return "a2a.discover" }
func (*a2aDiscoverer) Action() action.Action { return action.Discover }
func (*a2aDiscoverer) Target() string        { return "a2a" }
func (*a2aDiscoverer) Description() string {
	return "Discover A2A agents on a network via /.well-known/agent-card.json probe"
}
func (*a2aDiscoverer) Version() string     { return "0.3.0-dev" }
func (*a2aDiscoverer) IsDestructive() bool { return false }
