import type { APINode } from "@/api/types";
import { NODE_KIND_COLORS } from "@/theme/tokens";

export const NODE_COLORS: Record<string, string> = NODE_KIND_COLORS;

const BASE_SIZE = 6;
const MAX_SIZE = 20;

export function getNodeColor(kinds: string[]): string {
  for (const kind of kinds) {
    if (kind in NODE_COLORS) return NODE_COLORS[kind]!;
  }
  return "#999999";
}

export function getNodeSize(node: APINode): number {
  const kind = node.kinds[0] ?? "";
  const props = node.properties;

  switch (kind) {
    case "AgentInstance": {
      const score = Number(props.risk_score ?? 0);
      return BASE_SIZE + (score / 100) * (MAX_SIZE - BASE_SIZE);
    }
    case "MCPServer": {
      const tools = Number(props.tool_count ?? 0);
      return Math.min(BASE_SIZE + tools * 1.5, MAX_SIZE);
    }
    case "MCPTool": {
      const caps = Array.isArray(props.capability_surface)
        ? props.capability_surface.length
        : 0;
      return Math.min(BASE_SIZE + caps * 2, MAX_SIZE);
    }
    case "MCPResource": {
      const sensitivity = String(props.sensitivity ?? "low");
      if (sensitivity === "critical") return 14;
      if (sensitivity === "high") return 11;
      if (sensitivity === "medium") return 9;
      return BASE_SIZE;
    }
    case "A2AAgent": {
      const skills = Number(props.skill_count ?? 0);
      return Math.min(BASE_SIZE + skills * 1.5, MAX_SIZE);
    }
    case "Credential": {
      const exposed = Boolean(props.is_exposed);
      return exposed ? 12 : BASE_SIZE;
    }
    case "ConfigFile": {
      const servers = Number(props.server_count ?? 0);
      return Math.min(BASE_SIZE + servers, MAX_SIZE);
    }
    // v0.2 — LiteLLM gateways carry the highest blast radius among AI
    // services because a single master key fans out into every
    // upstream provider key. We size them by EXPOSES_CREDENTIAL
    // out-degree so a gateway exposing 5 provider keys looks
    // distinctly bigger than one exposing 1. The property
    // exposes_credential_count is populated by the v0.2 scan envelope
    // post-process; we fall back to BASE_SIZE when missing (the
    // edge-count summary lands in v0.3 once the post-processor
    // back-fills it).
    case "LiteLLMGateway": {
      const fanout = Number(props.exposes_credential_count ?? 0);
      return Math.min(BASE_SIZE + fanout * 1.5, MAX_SIZE);
    }
    case "OllamaInstance": {
      // Ollama nodes scale with the number of models loaded on the
      // host (a fine-tune-heavy box is more interesting than a stock
      // install). The looter that populates loaded_model_count is in
      // v0.3+; until then this falls back to BASE_SIZE.
      const models = Number(props.loaded_model_count ?? 0);
      return Math.min(BASE_SIZE + models, MAX_SIZE);
    }
    case "AIService": {
      // Generic umbrella fallback when a node carries only the
      // umbrella label (shouldn't happen with v0.2 emitters, but
      // defensive against forward-compat issues with v0.3+ kinds).
      return BASE_SIZE + 2;
    }
    default:
      return BASE_SIZE;
  }
}

export function getNodeLabel(node: APINode): string {
  const p = node.properties;
  return String(
    p.name ?? p.uri ?? p.path ?? p.hostname ?? p.id ?? node.id.slice(0, 12),
  );
}
