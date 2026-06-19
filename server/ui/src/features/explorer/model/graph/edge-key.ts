import type { APIEdge } from "@entities/graph/dto";

/** Unique per-edge key: `${source}|${target}|${kind}`. */
export function edgeKey(e: APIEdge): string {
  return `${e.source}|${e.target}|${e.kind}`;
}

/** Bundling key collapsing parallel edges between the same pair of nodes. */
export function bundleKey(e: APIEdge): string {
  return `${e.source}|${e.target}`;
}
