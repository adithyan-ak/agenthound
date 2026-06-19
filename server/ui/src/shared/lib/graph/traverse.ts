/**
 * Generic, domain-free graph traversal. Builds an adjacency index once and
 * runs a hop-bounded BFS over it, replacing the per-node "scan every edge on
 * every hop" loops that were duplicated across the explorer (O(N·E) -> O(E)).
 *
 * Edges are opaque to this module: callers supply `edgeKey` (and optionally
 * `includeEdge`) closures so any edge shape with `source`/`target` works.
 */

export type TraversalDirection = "in" | "out" | "both";

export interface MinimalEdge {
  source: string;
  target: string;
}

export interface AdjacencyIndex<E extends MinimalEdge> {
  /** Edges keyed by their `source` node id. */
  outgoing: Map<string, E[]>;
  /** Edges keyed by their `target` node id. */
  incoming: Map<string, E[]>;
}

export function buildAdjacencyIndex<E extends MinimalEdge>(
  edges: readonly E[],
): AdjacencyIndex<E> {
  const outgoing = new Map<string, E[]>();
  const incoming = new Map<string, E[]>();
  for (const e of edges) {
    let outList = outgoing.get(e.source);
    if (!outList) {
      outList = [];
      outgoing.set(e.source, outList);
    }
    outList.push(e);

    let inList = incoming.get(e.target);
    if (!inList) {
      inList = [];
      incoming.set(e.target, inList);
    }
    inList.push(e);
  }
  return { outgoing, incoming };
}

export interface BfsOptions<E extends MinimalEdge> {
  /** Maximum number of hops to expand from the start node. */
  maxHops: number;
  /** Which edge directions to follow. */
  direction: TraversalDirection;
  /** Stable key for an edge; collected into `edgeKeys`. */
  edgeKey: (edge: E) => string;
  /** Optional per-edge filter; edges returning false are skipped. */
  includeEdge?: (edge: E) => boolean;
}

export interface BfsResult {
  nodeIds: Set<string>;
  edgeKeys: Set<string>;
}

/**
 * Hop-bounded BFS from `startId`. Mirrors the frontier-per-hop structure of
 * the original explorer traversals: a node is "processed" only while it sits
 * in a frontier (hops 0..maxHops-1), and every traversed edge incident to a
 * processed node (in an allowed direction) is recorded. Results are returned
 * as sets — order is not significant for any consumer.
 */
export function bfsFrom<E extends MinimalEdge>(
  startId: string,
  index: AdjacencyIndex<E>,
  options: BfsOptions<E>,
): BfsResult {
  const { maxHops, direction, edgeKey, includeEdge } = options;
  const nodeIds = new Set<string>([startId]);
  const edgeKeys = new Set<string>();
  let frontier: string[] = [startId];

  for (let hop = 0; hop < maxHops && frontier.length > 0; hop++) {
    const next: string[] = [];
    for (const current of frontier) {
      if (direction !== "in") {
        for (const e of index.outgoing.get(current) ?? []) {
          if (includeEdge && !includeEdge(e)) continue;
          if (!nodeIds.has(e.target)) {
            nodeIds.add(e.target);
            next.push(e.target);
          }
          edgeKeys.add(edgeKey(e));
        }
      }
      if (direction !== "out") {
        for (const e of index.incoming.get(current) ?? []) {
          if (includeEdge && !includeEdge(e)) continue;
          if (!nodeIds.has(e.source)) {
            nodeIds.add(e.source);
            next.push(e.source);
          }
          edgeKeys.add(edgeKey(e));
        }
      }
    }
    frontier = next;
  }

  return { nodeIds, edgeKeys };
}
