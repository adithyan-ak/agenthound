package openwebuifp

import (
	"log/slog"

	"github.com/adithyan-ak/agenthound/sdk/action"
	"github.com/adithyan-ak/agenthound/sdk/module"
)

func init() {
	f, err := New()
	if err != nil {
		slog.Warn("openwebui fingerprinter init failed; service will not be detected",
			"error", err)
		module.Register(&disabledFingerprinter{})
		return
	}
	module.Register(f)
}

func (*Fingerprinter) ID() string            { return "openwebui.fingerprint" }
func (*Fingerprinter) Action() action.Action { return action.Fingerprint }
func (*Fingerprinter) Target() string        { return "openwebui" }
func (*Fingerprinter) Description() string {
	return "Identify Open WebUI by GET /api/version; capture configured Ollama backend"
}
func (*Fingerprinter) Version() string     { return "0.3.0-dev" }
func (*Fingerprinter) IsDestructive() bool { return false }

type disabledFingerprinter struct{}

func (*disabledFingerprinter) ID() string            { return "openwebui.fingerprint" }
func (*disabledFingerprinter) Action() action.Action { return action.Fingerprint }
func (*disabledFingerprinter) Target() string        { return "openwebui" }
func (*disabledFingerprinter) Description() string {
	return "Open WebUI fingerprinter (disabled — rule failed to load)"
}
func (*disabledFingerprinter) Version() string     { return "0.3.0-dev" }
func (*disabledFingerprinter) IsDestructive() bool { return false }
