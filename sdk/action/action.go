package action

// Action identifies a phase of an engagement. The string values are stable
// identifiers used in CLI flags, registry lookups, and module metadata.
type Action string

const (
	Scan        Action = "scan"
	Fingerprint Action = "fingerprint"
	Enumerate   Action = "enumerate"
	Loot        Action = "loot"
	Extract     Action = "extract"
	Poison      Action = "poison"
	Implant     Action = "implant"
)
