package model

// AllowedNodeKinds are the 12 collector-produced node kinds accepted in ingest input.
var AllowedNodeKinds = map[string]bool{
	"MCPServer":       true,
	"MCPTool":         true,
	"MCPResource":     true,
	"MCPPrompt":       true,
	"A2AAgent":        true,
	"A2ASkill":        true,
	"AgentInstance":   true,
	"Identity":        true,
	"Credential":      true,
	"Host":            true,
	"ConfigFile":      true,
	"InstructionFile": true,
}

// AllNodeLabels includes all 14 node labels (12 collector + 2 synthetic) for Neo4j schema constraints.
var AllNodeLabels = []string{
	"MCPServer", "MCPTool", "MCPResource", "MCPPrompt",
	"A2AAgent", "A2ASkill", "AgentInstance",
	"Identity", "Credential", "Host",
	"ConfigFile", "InstructionFile",
	"ResourceGroup", "TrustZone",
}

// RawEdgeKinds are the 13 collector-produced edge kinds accepted in ingest input.
var RawEdgeKinds = map[string]bool{
	"TRUSTS_SERVER":      true,
	"PROVIDES_TOOL":      true,
	"PROVIDES_RESOURCE":  true,
	"PROVIDES_PROMPT":    true,
	"ADVERTISES_SKILL":   true,
	"DELEGATES_TO":       true,
	"AUTHENTICATES_WITH": true,
	"USES_CREDENTIAL":    true,
	"RUNS_ON":            true,
	"CONFIGURED_IN":      true,
	"HAS_ENV_VAR":        true,
	"LOADS_INSTRUCTIONS": true,
	"SAME_AUTH_DOMAIN":   true,
}

// AllowedEdgeKinds includes all 21 edge kinds (13 raw + 8 composite) for Neo4j writer dispatch.
var AllowedEdgeKinds = map[string]bool{
	// Raw (collector-produced)
	"TRUSTS_SERVER":      true,
	"PROVIDES_TOOL":      true,
	"PROVIDES_RESOURCE":  true,
	"PROVIDES_PROMPT":    true,
	"ADVERTISES_SKILL":   true,
	"DELEGATES_TO":       true,
	"AUTHENTICATES_WITH": true,
	"USES_CREDENTIAL":    true,
	"RUNS_ON":            true,
	"CONFIGURED_IN":      true,
	"HAS_ENV_VAR":        true,
	"LOADS_INSTRUCTIONS": true,
	"SAME_AUTH_DOMAIN":   true,
	// Composite (post-processor produced)
	"HAS_ACCESS_TO":         true,
	"CAN_EXECUTE":           true,
	"CAN_REACH":             true,
	"CAN_EXFILTRATE_VIA":    true,
	"SHADOWS":               true,
	"POISONED_DESCRIPTION":  true,
	"CAN_IMPERSONATE":       true,
	"POISONED_INSTRUCTIONS": true,
}

// AllowedCollectors are the valid collector identifiers in ingest meta.
var AllowedCollectors = map[string]bool{
	"mcp":    true,
	"a2a":    true,
	"config": true,
}

// EdgeEndpoints defines the expected source and target node kinds for an edge kind.
type EdgeEndpoints struct {
	SourceKinds []string
	TargetKinds []string
}

// EdgeKindEndpoints maps each edge kind to its expected source/target node labels.
var EdgeKindEndpoints = map[string]EdgeEndpoints{
	"TRUSTS_SERVER":         {SourceKinds: []string{"AgentInstance"}, TargetKinds: []string{"MCPServer"}},
	"PROVIDES_TOOL":         {SourceKinds: []string{"MCPServer"}, TargetKinds: []string{"MCPTool"}},
	"PROVIDES_RESOURCE":     {SourceKinds: []string{"MCPServer"}, TargetKinds: []string{"MCPResource"}},
	"PROVIDES_PROMPT":       {SourceKinds: []string{"MCPServer"}, TargetKinds: []string{"MCPPrompt"}},
	"ADVERTISES_SKILL":      {SourceKinds: []string{"A2AAgent"}, TargetKinds: []string{"A2ASkill"}},
	"DELEGATES_TO":          {SourceKinds: []string{"A2AAgent"}, TargetKinds: []string{"A2AAgent"}},
	"AUTHENTICATES_WITH":    {SourceKinds: []string{"MCPServer", "A2AAgent"}, TargetKinds: []string{"Identity"}},
	"USES_CREDENTIAL":       {SourceKinds: []string{"Identity"}, TargetKinds: []string{"Credential"}},
	"RUNS_ON":               {SourceKinds: []string{"MCPServer", "A2AAgent"}, TargetKinds: []string{"Host"}},
	"CONFIGURED_IN":         {SourceKinds: []string{"MCPServer"}, TargetKinds: []string{"ConfigFile"}},
	"HAS_ENV_VAR":           {SourceKinds: []string{"MCPServer"}, TargetKinds: []string{"Credential"}},
	"LOADS_INSTRUCTIONS":    {SourceKinds: []string{"AgentInstance"}, TargetKinds: []string{"InstructionFile"}},
	"SAME_AUTH_DOMAIN":      {SourceKinds: []string{"A2AAgent"}, TargetKinds: []string{"A2AAgent"}},
	"HAS_ACCESS_TO":         {SourceKinds: []string{"MCPTool"}, TargetKinds: []string{"MCPResource"}},
	"CAN_EXECUTE":           {SourceKinds: []string{"MCPTool"}, TargetKinds: []string{"Host"}},
	"CAN_REACH":             {SourceKinds: []string{"AgentInstance", "A2AAgent"}, TargetKinds: []string{"MCPResource"}},
	"CAN_EXFILTRATE_VIA":    {SourceKinds: []string{"AgentInstance"}, TargetKinds: []string{"MCPTool"}},
	"SHADOWS":               {SourceKinds: []string{"MCPTool"}, TargetKinds: []string{"MCPTool"}},
	"POISONED_DESCRIPTION":  {SourceKinds: []string{"MCPTool"}, TargetKinds: []string{"MCPTool"}},
	"CAN_IMPERSONATE":       {SourceKinds: []string{"A2AAgent"}, TargetKinds: []string{"A2AAgent"}},
	"POISONED_INSTRUCTIONS": {SourceKinds: []string{"InstructionFile"}, TargetKinds: []string{"InstructionFile"}},
}

// ResolveEdgeEndpoints returns the source and target node kinds for an edge,
// using explicit values when set, falling back to the EdgeKindEndpoints registry.
func ResolveEdgeEndpoints(kind, sourceKind, targetKind string) (string, string) {
	if sourceKind != "" && targetKind != "" {
		return sourceKind, targetKind
	}
	ep, ok := EdgeKindEndpoints[kind]
	if !ok {
		return sourceKind, targetKind
	}
	if sourceKind == "" && len(ep.SourceKinds) > 0 {
		sourceKind = ep.SourceKinds[0]
	}
	if targetKind == "" && len(ep.TargetKinds) > 0 {
		targetKind = ep.TargetKinds[0]
	}
	return sourceKind, targetKind
}
