package embeddinginvert

import (
	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

func init() {
	module.Register(&Extractor{})
}

func (*Extractor) ID() string            { return "embedding.extract" }
func (*Extractor) Action() action.Action { return action.Extract }
func (*Extractor) Target() string        { return "embedding-invert" }
func (*Extractor) Description() string {
	return "Detect fine-tune training signals via embedding-layer statistical outlier analysis (GGUF)"
}
func (*Extractor) Version() string     { return "0.5.0-dev" }
func (*Extractor) IsDestructive() bool { return false }
