package vllmfp

import (
	"log/slog"

	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

func init() {
	f, err := New()
	if err != nil {
		slog.Warn("vllm fingerprinter init failed; service will not be detected",
			"error", err)
		module.Register(&disabledFingerprinter{})
		return
	}
	module.Register(f)
}

func (*Fingerprinter) ID() string            { return "vllm.fingerprint" }
func (*Fingerprinter) Action() action.Action { return action.Fingerprint }
func (*Fingerprinter) Target() string        { return "vllm" }
func (*Fingerprinter) Description() string {
	return "Identify vLLM inference servers by GET /v1/models"
}
func (*Fingerprinter) Version() string     { return "0.3.0-dev" }
func (*Fingerprinter) IsDestructive() bool { return false }

type disabledFingerprinter struct{}

func (*disabledFingerprinter) ID() string            { return "vllm.fingerprint" }
func (*disabledFingerprinter) Action() action.Action { return action.Fingerprint }
func (*disabledFingerprinter) Target() string        { return "vllm" }
func (*disabledFingerprinter) Description() string {
	return "vLLM fingerprinter (disabled — rule failed to load)"
}
func (*disabledFingerprinter) Version() string     { return "0.3.0-dev" }
func (*disabledFingerprinter) IsDestructive() bool { return false }
