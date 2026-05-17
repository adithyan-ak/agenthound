package litellmfp

import (
	"log/slog"

	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

func init() {
	f, err := New()
	if err != nil {
		slog.Warn("litellm fingerprinter init failed; service will not be detected",
			"error", err)
		module.Register(&disabledFingerprinter{})
		return
	}
	module.Register(f)
}

func (*Fingerprinter) ID() string            { return "litellm.fingerprint" }
func (*Fingerprinter) Action() action.Action { return action.Fingerprint }
func (*Fingerprinter) Target() string        { return "litellm" }
func (*Fingerprinter) Description() string {
	return "Identify LiteLLM proxy/gateway servers by GET /health/liveliness"
}
func (*Fingerprinter) Version() string     { return "0.2.0-dev" }
func (*Fingerprinter) IsDestructive() bool { return false }

type disabledFingerprinter struct{}

func (*disabledFingerprinter) ID() string            { return "litellm.fingerprint" }
func (*disabledFingerprinter) Action() action.Action { return action.Fingerprint }
func (*disabledFingerprinter) Target() string        { return "litellm" }
func (*disabledFingerprinter) Description() string {
	return "LiteLLM fingerprinter (disabled — rule failed to load)"
}
func (*disabledFingerprinter) Version() string     { return "0.2.0-dev" }
func (*disabledFingerprinter) IsDestructive() bool { return false }
