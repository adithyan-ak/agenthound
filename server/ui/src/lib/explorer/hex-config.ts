import type { LucideIcon } from "lucide-react";
import {
  User,
  Server,
  Wrench,
  FileLock2,
  MessageSquareQuote,
  Bot,
  Sparkles,
  Monitor,
  Shield,
  KeyRound,
  FileCode2,
  FileText,
  Layers,
  ShieldCheck,
  Hexagon,
  Brain,
  Rocket,
  Database,
  FlaskConical,
  GitFork,
  Notebook,
  Link2,
  MessageSquare,
  Boxes,
} from "lucide-react";
import { SEVERITY } from "@/theme/tokens";
import type { SeverityLevel } from "./lens-config";

export interface HexKindConfig {
  /** Stroke color applied to the hex outline */
  strokeColor: string;
  /** Hex fill color (very dark, almost transparent) */
  fillColor: string;
  /** Icon rendered centered inside the hex */
  icon: LucideIcon;
  /** Uppercase tag shown below the label ("MCP SERVER", "AGENT INSTANCE", …) */
  kindTag: string;
  /** Layout column index (0 = leftmost entry, 4 = rightmost resources) */
  column: number;
  /** Display label for grouping UI */
  groupLabel: string;
}

export const HEX_CONFIG: Record<string, HexKindConfig> = {
  AgentInstance: {
    strokeColor: "#06B6D4",
    fillColor: "#0B1220",
    icon: User,
    kindTag: "AGENT",
    column: 0,
    groupLabel: "Agents",
  },
  A2AAgent: {
    strokeColor: "#A855F7",
    fillColor: "#0B1220",
    icon: Bot,
    kindTag: "A2A AGENT",
    column: 0,
    groupLabel: "Agents",
  },
  MCPServer: {
    strokeColor: "#10B981",
    fillColor: "#0B1220",
    icon: Server,
    kindTag: "MCP SERVER",
    column: 1,
    groupLabel: "Servers",
  },
  MCPTool: {
    strokeColor: "#F59E0B",
    fillColor: "#0B1220",
    icon: Wrench,
    kindTag: "MCP TOOL",
    column: 2,
    groupLabel: "Tools & Skills",
  },
  A2ASkill: {
    strokeColor: "#C084FC",
    fillColor: "#0B1220",
    icon: Sparkles,
    kindTag: "A2A SKILL",
    column: 2,
    groupLabel: "Tools & Skills",
  },
  MCPPrompt: {
    strokeColor: "#FB923C",
    fillColor: "#0B1220",
    icon: MessageSquareQuote,
    kindTag: "MCP PROMPT",
    column: 2,
    groupLabel: "Tools & Skills",
  },
  Host: {
    strokeColor: "#475569",
    fillColor: "#0B1220",
    icon: Monitor,
    kindTag: "HOST",
    column: 3,
    groupLabel: "Infra",
  },
  Identity: {
    strokeColor: "#94A3B8",
    fillColor: "#0B1220",
    icon: Shield,
    kindTag: "IDENTITY",
    column: 3,
    groupLabel: "Infra",
  },
  Credential: {
    strokeColor: "#EC4899",
    fillColor: "#0B1220",
    icon: KeyRound,
    kindTag: "CREDENTIAL",
    column: 3,
    groupLabel: "Infra",
  },
  MCPResource: {
    strokeColor: "#EF4444",
    fillColor: "#0B1220",
    icon: FileLock2,
    kindTag: "RESOURCE",
    column: 4,
    groupLabel: "Resources",
  },
  ConfigFile: {
    strokeColor: "#D97706",
    fillColor: "#0B1220",
    icon: FileCode2,
    kindTag: "CONFIG FILE",
    column: 4,
    groupLabel: "Resources",
  },
  InstructionFile: {
    strokeColor: "#EAB308",
    fillColor: "#0B1220",
    icon: FileText,
    kindTag: "INSTRUCTION",
    column: 4,
    groupLabel: "Resources",
  },
  ResourceGroup: {
    strokeColor: "#64748B",
    fillColor: "#0B1220",
    icon: Layers,
    kindTag: "RESOURCE GROUP",
    column: 4,
    groupLabel: "Resources",
  },
  TrustZone: {
    strokeColor: "#22D3EE",
    fillColor: "#0B1220",
    icon: ShieldCheck,
    kindTag: "TRUST ZONE",
    column: 3,
    groupLabel: "Infra",
  },

  // v0.2 AI services. All sit in column 2 ("Tools & Skills") so they
  // place between :MCPServer and :MCPResource in the layered explorer
  // view — visually consistent with how a LiteLLM gateway sits between
  // an MCP server (which knows the master key) and the upstream
  // provider credentials it exposes. Per-kind label is kinds[0]; the
  // umbrella :AIService is rendered via the AIService entry only when
  // a node has NO known per-kind label (defensive fallback for
  // forward-compat with v0.3+ kinds the UI hasn't been taught yet).
  OllamaInstance: {
    strokeColor: "#FF7043",
    fillColor: "#0B1220",
    icon: Brain,
    kindTag: "OLLAMA",
    column: 2,
    groupLabel: "AI Services",
  },
  LiteLLMGateway: {
    strokeColor: "#EC407A",
    fillColor: "#0B1220",
    icon: GitFork,
    kindTag: "LITELLM",
    column: 2,
    groupLabel: "AI Services",
  },
  AIService: {
    strokeColor: "#7E57C2",
    fillColor: "#0B1220",
    icon: Hexagon,
    kindTag: "AI SERVICE",
    column: 2,
    groupLabel: "AI Services",
  },
  // v0.3 / v0.4 fingerprinters. Final colors locked in to be visually
  // distinct from each other AND from the v0.1 palette above. Stroke
  // colors must match theme/tokens.ts NODE_KIND_COLORS exactly.
  VLLMInstance: {
    strokeColor: "#26A69A",
    fillColor: "#0B1220",
    icon: Rocket,
    kindTag: "VLLM",
    column: 2,
    groupLabel: "AI Services",
  },
  QdrantInstance: {
    strokeColor: "#5C6BC0",
    fillColor: "#0B1220",
    icon: Database,
    kindTag: "QDRANT",
    column: 2,
    groupLabel: "AI Services",
  },
  MLflowServer: {
    strokeColor: "#42A5F5",
    fillColor: "#0B1220",
    icon: FlaskConical,
    kindTag: "MLFLOW",
    column: 2,
    groupLabel: "AI Services",
  },
  JupyterServer: {
    strokeColor: "#F57C00",
    fillColor: "#0B1220",
    icon: Notebook,
    kindTag: "JUPYTER",
    column: 2,
    groupLabel: "AI Services",
  },
  LangServeApp: {
    strokeColor: "#9CCC65",
    fillColor: "#0B1220",
    icon: Link2,
    kindTag: "LANGSERVE",
    column: 2,
    groupLabel: "AI Services",
  },
  OpenWebUIInstance: {
    strokeColor: "#66BB6A",
    fillColor: "#0B1220",
    icon: MessageSquare,
    kindTag: "OPEN WEBUI",
    column: 2,
    groupLabel: "AI Services",
  },
  // v0.3 model-artifact node. Sits one column right of OllamaInstance so the
  // PROVIDES_MODEL edge reads OllamaInstance(col 2) -> AIModel(col 3) — model
  // artifacts visually live "downstream" of the service that hosts them.
  // Distinct from MCPResource (col 4) which is a remote-resource pointer, not
  // a stored artifact.
  AIModel: {
    strokeColor: "#6A1B9A",
    fillColor: "#0B1220",
    icon: Boxes,
    kindTag: "AI MODEL",
    column: 3,
    groupLabel: "AI Models",
  },
};

const FALLBACK: HexKindConfig = {
  strokeColor: "#94A3B8",
  fillColor: "#0B1220",
  icon: Hexagon,
  kindTag: "NODE",
  column: 3,
  groupLabel: "Other",
};

export function getHexConfig(kind: string): HexKindConfig {
  return HEX_CONFIG[kind] ?? FALLBACK;
}

export const SEVERITY_HALO: Record<Exclude<SeverityLevel, "info">, string> = {
  critical: "drop-shadow(0 0 10px rgba(239,68,68,0.85)) drop-shadow(0 0 18px rgba(239,68,68,0.45))",
  high: "drop-shadow(0 0 8px rgba(249,115,22,0.75)) drop-shadow(0 0 16px rgba(249,115,22,0.35))",
  medium: "drop-shadow(0 0 6px rgba(234,179,8,0.65)) drop-shadow(0 0 12px rgba(234,179,8,0.3))",
  low: "drop-shadow(0 0 4px rgba(148,163,184,0.5))",
};

export const SEVERITY_STROKE_COLOR: Record<Exclude<SeverityLevel, "info">, string> = {
  critical: SEVERITY.critical!.solid,
  high: SEVERITY.high!.solid,
  medium: SEVERITY.medium!.solid,
  low: SEVERITY.low!.solid,
};

/**
 * Hex node visual constants. All hex nodes share the same outer viewport
 * size so handles and port dots stay at consistent relative positions,
 * regardless of the kind-specific stroke or halo.
 */
export const HEX_NODE_WIDTH = 84;
export const HEX_NODE_HEIGHT = 96;
export const HEX_LABEL_HEIGHT = 32;
export const HEX_TOTAL_HEIGHT = HEX_NODE_HEIGHT + HEX_LABEL_HEIGHT;

/**
 * Six vertex positions for the hexagon (pointy-top orientation), used for
 * both the decorative port dots and the invisible React Flow handles.
 * Coordinates are in the node's local pixel space.
 */
export const HEX_VERTICES: Array<{
  id: string;
  x: number;
  y: number;
  side: "top" | "right" | "bottom" | "left";
}> = [
  { id: "h-top-left", x: 8, y: 22, side: "left" },
  { id: "h-top-right", x: 76, y: 22, side: "right" },
  { id: "h-right", x: 82, y: 48, side: "right" },
  { id: "h-bottom-right", x: 76, y: 74, side: "right" },
  { id: "h-bottom-left", x: 8, y: 74, side: "left" },
  { id: "h-left", x: 2, y: 48, side: "left" },
];

/**
 * SVG polygon points for the hex, rendered inside an 84x96 viewBox.
 * Pointy-top hexagon oriented to match the Trickest aesthetic.
 */
export const HEX_POLYGON_POINTS = "42,4 78,22 78,74 42,92 6,74 6,22";
