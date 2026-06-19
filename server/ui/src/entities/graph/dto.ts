// Shared graph wire types — the raw shapes returned by the graph API and
// consumed across multiple entities (node, edge, finding, explorer). These are
// the canonical DTOs; richer per-entity view-models build on top of them.

export type NodeKind =
  | "MCPServer"
  | "MCPTool"
  | "MCPResource"
  | "MCPPrompt"
  | "A2AAgent"
  | "A2ASkill"
  | "AgentInstance"
  | "Identity"
  | "Credential"
  | "Host"
  | "ConfigFile"
  | "InstructionFile"
  | "OllamaInstance"
  | "VLLMInstance"
  | "QdrantInstance"
  | "MLflowServer"
  | "LiteLLMGateway"
  | "JupyterServer"
  | "LangServeApp"
  | "OpenWebUIInstance"
  | "AIService"
  | "AIModel"
  | "ExtractedTrainingSignal"
  | "ResourceGroup"
  | "TrustZone";

export type EdgeKind =
  | "TRUSTS_SERVER"
  | "PROVIDES_TOOL"
  | "PROVIDES_RESOURCE"
  | "PROVIDES_PROMPT"
  | "ADVERTISES_SKILL"
  | "DELEGATES_TO"
  | "AUTHENTICATES_WITH"
  | "USES_CREDENTIAL"
  | "RUNS_ON"
  | "CONFIGURED_IN"
  | "HAS_ENV_VAR"
  | "LOADS_INSTRUCTIONS"
  | "SAME_AUTH_DOMAIN"
  | "EXPOSES"
  | "EXPOSES_CREDENTIAL"
  | "PROVIDES_MODEL"
  | "EXTRACTED_FROM"
  | "HAS_ACCESS_TO"
  | "CAN_EXECUTE"
  | "SHADOWS"
  | "POISONED_DESCRIPTION"
  | "CAN_REACH"
  | "CAN_EXFILTRATE_VIA"
  | "CAN_IMPERSONATE"
  | "POISONED_INSTRUCTIONS";

export interface APINode {
  id: string;
  kinds: string[];
  properties: Record<string, unknown>;
}

export interface APIEdge {
  source: string;
  target: string;
  kind: string;
  source_kind?: string;
  target_kind?: string;
  properties: Record<string, unknown>;
}
