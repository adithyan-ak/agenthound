import type { APIEdge, APINode } from "@entities/graph/dto";

export interface RemediationItem {
  severity: "critical" | "high" | "medium" | "low";
  title: string;
  body: string;
}

/**
 * Derives the node-level remediation checklist from a node's properties and the
 * composite edges it participates in. Extracted verbatim from RemediationTab so
 * the rules live in the security domain rather than the presentation layer.
 */
export function deriveRemediations(
  node: APINode,
  kind: string,
  edges: APIEdge[],
): RemediationItem[] {
  const items: RemediationItem[] = [];
  const props = node.properties ?? {};

  // Unpinned package
  if (kind === "MCPServer" && props.is_pinned === false) {
    items.push({
      severity: "medium",
      title: "Pin this server's package version",
      body: "The server was launched without a pinned version (e.g. `npx -y @pkg` without `@x.y.z`). A malicious update could ship new tool descriptions or behavior without warning. Pin the version in the client config.",
    });
  }

  // No auth
  if ((kind === "MCPServer" || kind === "A2AAgent") && props.auth_method === "none") {
    items.push({
      severity: "high",
      title: "Add an authentication method",
      body: "This endpoint accepts requests without any authentication. Configure at minimum a bearer token or API key, and prefer OAuth or mTLS for anything reaching sensitive resources.",
    });
  }

  // Exposed credential
  if (kind === "Credential" && props.is_exposed === true) {
    items.push({
      severity: "critical",
      title: "Rotate this credential",
      body: "This credential was found inlined in a config file or environment variable in plaintext. Rotate it immediately, move it to a secret manager, and reference it by env var alone.",
    });
  }

  // High entropy secret
  if (kind === "Credential" && props.high_entropy === true && props.is_exposed !== true) {
    items.push({
      severity: "medium",
      title: "Review this high-entropy value",
      body: "The value has Shannon entropy high enough to suggest it may be a raw secret. Confirm it is referenced via environment variable or vault, not inlined.",
    });
  }

  // Poisoned tool
  if (kind === "MCPTool" && props.has_injection_patterns === true) {
    items.push({
      severity: "high",
      title: "Remove or re-review this tool",
      body: "This tool's description contains patterns consistent with prompt injection (`<IMPORTANT>` tags, 'ignore previous instructions', hidden Unicode). Agents that read the description as planning context will treat it as trusted instructions.",
    });
  }

  // Poisoned instruction file
  if (kind === "InstructionFile" && props.is_suspicious === true) {
    items.push({
      severity: "high",
      title: "Inspect and sanitize this instruction file",
      body: "This file contains suspicious directives (imperative overrides, outbound curl/wget, or encoded payloads). Review the file, remove any injected directives, and add it to your repo's suspicious-path audit list.",
    });
  }

  // Composite edge participation
  const hasCanExfiltrate = edges.some(
    (e) => e.kind === "CAN_EXFILTRATE_VIA" && (e.source === node.id || e.target === node.id),
  );
  if (hasCanExfiltrate) {
    items.push({
      severity: "critical",
      title: "Break the exfiltration path",
      body: "This node participates in a computed CAN_EXFILTRATE_VIA path. Either remove the sensitive resource reach, or remove the outbound channel from the same agent's trust scope. Both legs must remain disabled to fully close the exfil path.",
    });
  }

  const hasCriticalCanReach = edges.some(
    (e) =>
      e.kind === "CAN_REACH" &&
      (e.source === node.id || e.target === node.id) &&
      e.properties?.cross_protocol === true,
  );
  if (hasCriticalCanReach) {
    items.push({
      severity: "critical",
      title: "Close the cross-protocol pivot",
      body: "This node is on a cross-protocol CAN_REACH path (A2A agent pivoting through an MCP server). Enforce auth on the delegated agent AND remove the shared host correlation between the two protocols.",
    });
  }

  return items;
}
