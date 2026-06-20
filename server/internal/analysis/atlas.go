package analysis

// ATLASTechniques is the single source of truth for the MITRE ATLAS techniques
// AgentHound maps findings to. IDs and titles are verbatim from MITRE ATLAS
// v5.5.0 (2026-03-30) — the release that introduced AML.T0110, the newest
// technique referenced here (AML.T0086 arrived in v5.0.0; the agentic-AI
// techniques did not exist in any v4.x release). MITRE renamed several
// "ML"→"AI"; titles below are the current names.
var ATLASTechniques = map[string]string{
	"AML.T0051": "LLM Prompt Injection",
	"AML.T0110": "AI Agent Tool Poisoning",
	"AML.T0086": "Exfiltration via AI Agent Tool Invocation",
	"AML.T0024": "Exfiltration via AI Inference API",
	"AML.T0057": "LLM Data Leakage",
	"AML.T0020": "Poison Training Data",
	"AML.T0010": "AI Supply Chain Compromise",
}
