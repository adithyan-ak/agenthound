import { AlertOctagon, AlertTriangle, AlertCircle, Info } from "lucide-react";
import type { LucideIcon } from "lucide-react";

// Canonical node kind colors — matches hex-config palette
export const NODE_KIND_COLORS: Record<string, string> = {
  AgentInstance: "#06B6D4",    // cyan-500
  A2AAgent: "#A855F7",        // purple-500
  MCPServer: "#10B981",       // emerald-500
  MCPTool: "#F59E0B",         // amber-500
  MCPResource: "#EF4444",     // red-500
  MCPPrompt: "#FB923C",       // orange-400
  A2ASkill: "#C084FC",        // purple-400
  Host: "#475569",            // slate-600
  Identity: "#94A3B8",        // slate-400
  Credential: "#EC4899",      // pink-500
  ConfigFile: "#D97706",      // amber-600
  InstructionFile: "#EAB308", // yellow-500
  ResourceGroup: "#64748B",   // slate-500
  TrustZone: "#22D3EE",       // cyan-400
  // AI service kinds. Per-kind label is the dispatch key (per
  // node-styles.ts kinds[0] semantics); the umbrella :AIService stays
  // on the node as a multi-label companion so unified queries can match
  // (s:AIService) without enumerating per-kind labels. Every per-service
  // color is checked against the v0.1 palette above for visual distinctness
  // — the umbrella :AIService color is the fallback only when the
  // per-kind dispatch returns nothing.
  OllamaInstance: "#FF7043",     // orange-red — distinct from MCPTool amber
  LiteLLMGateway: "#EC407A",     // pink — sibling of Credential pink (the gateway IS the credential exposure)
  AIService: "#7E57C2",          // muted purple — generic umbrella fallback
  VLLMInstance: "#26A69A",       // teal — vLLM's "rocket" identity
  QdrantInstance: "#5C6BC0",     // indigo — vector-search identity, distinct from AgentInstance cyan
  MLflowServer: "#42A5F5",       // blue — MLflow brand-aligned blue
  JupyterServer: "#F57C00",      // deep orange — distinct from MCPPrompt #FB923C (orange-400)
  LangServeApp: "#9CCC65",       // chartreuse — distinct from AIService purple AND OpenWebUI green
  OpenWebUIInstance: "#66BB6A",  // green — chat-UI identity
  // v0.3 model-artifact node — emitted by the Ollama Looter (one per /api/tags
  // entry, properties from /api/show). Deep purple chosen because it sits next
  // to AIService #7E57C2 and A2AAgent #A855F7 in hue but is materially darker,
  // so the explorer renders model artifacts as the "weight" beneath their AI
  // service. Plan suggested #F44336 red but it collides with MCPResource
  // #EF4444 — readers can't distinguish a sensitive resource from a model.
  AIModel: "#6A1B9A",            // deep purple — distinct from A2AAgent / AIService / QdrantInstance
};

// Severity system — solid color, muted bg, text, border, icon
export interface SeverityStyle {
  solid: string;
  bg: string;
  text: string;
  border: string;
  icon: LucideIcon;
  label: string;
  /** Tailwind classes for badge: bg + text + border */
  badgeClass: string;
  /** Tailwind class for left-border accent */
  borderLeftClass: string;
  /** Tailwind class for solid dot/indicator */
  dotClass: string;
}

export const SEVERITY: Record<string, SeverityStyle> = {
  critical: {
    solid: "#EF4444",
    bg: "rgba(239,68,68,0.12)",
    text: "#F87171",
    border: "rgba(239,68,68,0.30)",
    icon: AlertOctagon,
    label: "Critical",
    badgeClass: "bg-red-500/10 text-red-400 border-red-500/30",
    borderLeftClass: "border-l-red-500",
    dotClass: "bg-red-500",
  },
  high: {
    solid: "#F97316",
    bg: "rgba(249,115,22,0.12)",
    text: "#FB923C",
    border: "rgba(249,115,22,0.30)",
    icon: AlertTriangle,
    label: "High",
    badgeClass: "bg-orange-500/10 text-orange-400 border-orange-500/30",
    borderLeftClass: "border-l-orange-500",
    dotClass: "bg-orange-500",
  },
  medium: {
    solid: "#EAB308",
    bg: "rgba(234,179,8,0.12)",
    text: "#FACC15",
    border: "rgba(234,179,8,0.30)",
    icon: AlertCircle,
    label: "Medium",
    badgeClass: "bg-yellow-500/10 text-yellow-400 border-yellow-500/30",
    borderLeftClass: "border-l-yellow-500",
    dotClass: "bg-yellow-500",
  },
  low: {
    solid: "#64748B",
    bg: "rgba(100,120,143,0.12)",
    text: "#94A3B8",
    border: "rgba(100,120,143,0.20)",
    icon: Info,
    label: "Low",
    badgeClass: "bg-slate-500/10 text-slate-400 border-slate-500/30",
    borderLeftClass: "border-l-slate-500",
    dotClass: "bg-slate-500",
  },
  info: {
    solid: "#3B82F6",
    bg: "rgba(59,130,246,0.12)",
    text: "#60A5FA",
    border: "rgba(59,130,246,0.30)",
    icon: Info,
    label: "Info",
    badgeClass: "bg-blue-500/10 text-blue-400 border-blue-500/30",
    borderLeftClass: "border-l-blue-500",
    dotClass: "bg-blue-500",
  },
};

// Feedback colors for form validation / system messages
export const FEEDBACK = {
  success: { solid: "#22C55E", bg: "rgba(34,197,94,0.12)", text: "#4ADE80", border: "rgba(34,197,94,0.30)" },
  warning: { solid: "#F59E0B", bg: "rgba(245,158,11,0.12)", text: "#FBBF24", border: "rgba(245,158,11,0.30)" },
  error: { solid: "#EF4444", bg: "rgba(239,68,68,0.12)", text: "#F87171", border: "rgba(239,68,68,0.30)" },
  info: { solid: "#3B82F6", bg: "rgba(59,130,246,0.12)", text: "#60A5FA", border: "rgba(59,130,246,0.30)" },
} as const;

// Recharts theme constants
export const CHART_THEME = {
  tooltip: {
    bg: "#111B2E",
    border: "#1A2540",
    text: "#EDF0F3",
  },
  grid: "rgba(255,255,255,0.05)",
  axis: "#64788F",
  series: ["#06B6D4", "#A855F7", "#10B981", "#F59E0B", "#EC4899", "#3B82F6", "#EF4444", "#22D3EE"],
} as const;

// Risk score -> color utility
export function riskColor(score: number): string {
  if (score >= 75) return "#EF4444";
  if (score >= 50) return "#F97316";
  if (score >= 25) return "#EAB308";
  return "#22C55E";
}

// Canonical severity ordering (worst first) for dashboards and legends.
export const SEVERITY_ORDER = ["critical", "high", "medium", "low"] as const;

// Severity key -> solid color, falling back to the "low" slate.
export function severityColor(sev: string): string {
  return SEVERITY[sev]?.solid ?? SEVERITY.low!.solid;
}

// Risk score -> Tailwind bg class
export function riskBgClass(score: number): string {
  if (score >= 75) return "bg-red-500";
  if (score >= 50) return "bg-orange-500";
  if (score >= 25) return "bg-yellow-500";
  return "bg-green-500";
}

// Risk score -> Tailwind text class
export function riskTextClass(score: number): string {
  if (score >= 75) return "text-red-400";
  if (score >= 50) return "text-orange-400";
  if (score >= 25) return "text-amber-400";
  return "text-green-400";
}

// Edge category colors
export const EDGE_COLORS = {
  attack: "#EF4444",
  trust: "#4A90D9",
  structure: "#475569",
} as const;

// Explorer canvas background
export const EXPLORER_CANVAS_BG = "#050B18";
export const EXPLORER_HEX_FILL = "#0B1220";
