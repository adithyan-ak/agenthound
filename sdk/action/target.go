package action

// Target is the input shape for every per-target action (Fingerprint,
// Enumerate, Loot, Extract, Poison, Implant). v0 keeps it deliberately flat
// — Kind discriminates and Meta carries discovery hints. Typed sub-structs
// (HostTarget, URLTarget, ConfigTarget, etc.) may land at v1 if real
// consumers demand stronger typing.
type Target struct {
	Kind    string            // "host", "url", "config_path", "cidr_member", "local"
	Address string            // "10.0.0.42:11434", "https://api.example.com", ""
	Meta    map[string]string // discovery hints; conventions documented per Kind
}
