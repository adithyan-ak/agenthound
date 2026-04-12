import ELK, { type ElkNode, type ElkExtendedEdge } from "elkjs/lib/elk-api.js";
import ElkWorker from "elkjs/lib/elk-worker.min.js?worker";
import type { Node, Edge } from "@xyflow/react";
import { getHexConfig, HEX_NODE_WIDTH, HEX_TOTAL_HEIGHT } from "./hex-config";

const elk = new ELK({
  workerFactory: () => new ElkWorker() as unknown as Worker,
});

const COLUMN_SPACING = 260;
const ROW_SPACING = 48;

/**
 * ELK options for strict left-to-right column layout. Uses the `layered`
 * algorithm with partitioning so each node kind sits in its designated
 * column index (set per-node via `layoutOptions.partitioning.partition`).
 */
const ELK_OPTIONS: Record<string, string> = {
  "elk.algorithm": "layered",
  "elk.direction": "RIGHT",
  "elk.partitioning.activate": "true",
  "elk.layered.layering.strategy": "NETWORK_SIMPLEX",
  "elk.layered.nodePlacement.strategy": "NETWORK_SIMPLEX",
  "elk.layered.crossingMinimization.strategy": "LAYER_SWEEP",
  "elk.layered.crossingMinimization.forceNodeModelOrder": "true",
  "elk.spacing.nodeNode": String(ROW_SPACING),
  "elk.layered.spacing.nodeNodeBetweenLayers": String(COLUMN_SPACING),
  "elk.layered.spacing.edgeNodeBetweenLayers": "36",
  "elk.spacing.edgeNode": "20",
  "elk.layered.considerModelOrder.strategy": "NODES_AND_EDGES",
};

export interface LayoutResult<T extends Node = Node> {
  nodes: T[];
  bounds: { width: number; height: number };
}

/**
 * Compute left-to-right column layout for hex nodes. Each node is tagged
 * with its target column (0..4) based on its data.kind (via HEX_CONFIG),
 * then ELK's partitioning layering places it exactly on that column.
 * Within each column, ELK minimizes edge crossings. Accepts both regular
 * hex nodes and orphan-cluster nodes — both carry `data.kind`.
 */
export async function computeExplorerLayout<T extends Node = Node>(
  nodes: T[],
  edges: Edge[],
): Promise<LayoutResult<T>> {
  if (nodes.length === 0) return { nodes, bounds: { width: 0, height: 0 } };

  const elkChildren: ElkNode[] = nodes.map((n) => {
    const data = n.data as Record<string, unknown> | undefined;
    const kind = typeof data?.kind === "string" ? data.kind : "";
    const config = getHexConfig(kind);
    return {
      id: n.id,
      width: HEX_NODE_WIDTH,
      height: HEX_TOTAL_HEIGHT,
      layoutOptions: {
        "elk.partitioning.partition": String(config.column),
      },
    };
  });

  const nodeIds = new Set(nodes.map((n) => n.id));
  const elkEdges: ElkExtendedEdge[] = [];
  const seen = new Set<string>();

  for (const e of edges) {
    if (!nodeIds.has(e.source) || !nodeIds.has(e.target)) continue;
    const key = `${e.source}->${e.target}`;
    if (seen.has(key)) continue;
    seen.add(key);
    elkEdges.push({
      id: e.id,
      sources: [e.source],
      targets: [e.target],
    });
  }

  const result = await elk.layout({
    id: "explorer-root",
    layoutOptions: ELK_OPTIONS,
    children: elkChildren,
    edges: elkEdges,
  });

  const positions = new Map<string, { x: number; y: number }>();
  let maxX = 0;
  let maxY = 0;
  for (const child of result.children ?? []) {
    const x = child.x ?? 0;
    const y = child.y ?? 0;
    positions.set(child.id, { x, y });
    if (x + HEX_NODE_WIDTH > maxX) maxX = x + HEX_NODE_WIDTH;
    if (y + HEX_TOTAL_HEIGHT > maxY) maxY = y + HEX_TOTAL_HEIGHT;
  }

  const positioned = nodes.map((n) => {
    const pos = positions.get(n.id);
    if (!pos) return n;
    return { ...n, position: pos };
  });

  return { nodes: positioned, bounds: { width: maxX, height: maxY } };
}
