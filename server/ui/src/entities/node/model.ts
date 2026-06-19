import type { APINode } from "@entities/graph/dto";

// Typed accessors over the raw `properties: Record<string, unknown>` bag, plus
// the small set of derived values widgets re-coerced inline everywhere. Each
// helper reproduces the exact coercion it replaces, so wiring a call site to it
// is behavior-preserving.

export function nodeString(
  node: APINode,
  key: string,
  fallback = "",
): string {
  const value = node.properties[key];
  return value == null ? fallback : String(value);
}

export function nodeNumber(node: APINode, key: string, fallback = 0): number {
  return Number(node.properties[key] ?? fallback);
}

export function nodeBool(node: APINode, key: string): boolean {
  return node.properties[key] === true;
}

/** Best available human label, falling back through common identity props. */
export function displayName(node: APINode): string {
  return String(
    node.properties.name ??
      node.properties.uri ??
      node.properties.path ??
      node.properties.hostname ??
      node.id.slice(0, 12),
  );
}

/** Declared auth method, defaulting to "none" when unset. */
export function authMethod(node: APINode): string {
  return String(node.properties.auth_method ?? "none");
}

/** True when the node advertises no authentication. */
export function isUnauth(node: APINode): boolean {
  return authMethod(node) === "none";
}

/** Computed risk score (0 when unset). */
export function riskScore(node: APINode): number {
  return Number(node.properties.risk_score ?? 0);
}
