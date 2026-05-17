package module_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

// flagModule satisfies both Module and FlagsModule.
type flagModule struct {
	called bool
}

func (m *flagModule) ID() string            { return "test.flags" }
func (m *flagModule) Action() action.Action { return action.Loot }
func (m *flagModule) Target() string        { return "test" }
func (m *flagModule) Description() string   { return "test flags module" }
func (m *flagModule) Version() string       { return "0.0.0" }
func (m *flagModule) IsDestructive() bool   { return false }
func (m *flagModule) RegisterFlags(*pflag.FlagSet) {
	m.called = true
}

// noFlagModule satisfies Module but NOT FlagsModule.
type noFlagModule struct{}

func (*noFlagModule) ID() string            { return "test.noflags" }
func (*noFlagModule) Action() action.Action { return action.Loot }
func (*noFlagModule) Target() string        { return "test" }
func (*noFlagModule) Description() string   { return "module without flags" }
func (*noFlagModule) Version() string       { return "0.0.0" }
func (*noFlagModule) IsDestructive() bool   { return false }

func TestRegisterFlagsFor_AssertsAndCalls(t *testing.T) {
	m := &flagModule{}
	cmd := &cobra.Command{Use: "x"}
	module.RegisterFlagsFor(cmd, m)
	if !m.called {
		t.Fatal("RegisterFlags was not called on a FlagsModule")
	}
}

func TestRegisterFlagsFor_NoOpWhenModuleLacksInterface(t *testing.T) {
	cmd := &cobra.Command{Use: "x"}
	module.RegisterFlagsFor(cmd, &noFlagModule{}) // must not panic
}

func TestRegisterFlagsFor_NoOpWhenNil(t *testing.T) {
	// Both nils must be tolerated — the dispatcher uses this in code paths
	// where a target may not resolve to any module.
	module.RegisterFlagsFor(nil, &noFlagModule{})
	module.RegisterFlagsFor(&cobra.Command{Use: "x"}, nil)
}
