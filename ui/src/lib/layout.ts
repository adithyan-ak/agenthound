import ELK, { type ElkNode, type ElkExtendedEdge } from "elkjs/lib/elk-api.js";
import ElkWorker from "elkjs/lib/elk-worker.min.js?worker";
import type { Node, Edge } from "@xyflow/react";

const elk = new ELK({
  workerFactory: () => new ElkWorker() as unknown as Worker,
});

const TIER1_TYPES = new Set(["server", "a2aAgent"]);
const TIER3_TYPES = new Set(["infra"]);

function getNodeDims(type: string): { w: number; h: number } {
  if (TIER1_TYPES.has(type)) return { w: 180, h: 44 };
  if (TIER3_TYPES.has(type)) return { w: 28, h: 28 };
  return { w: 110, h: 26 };
}

const ELK_OPTIONS: Record<string, string> = {
  "elk.algorithm": "layered",
  "elk.direction": "RIGHT",
  "elk.layered.layering.strategy": "NETWORK_SIMPLEX",
  "elk.layered.nodePlacement.strategy": "NETWORK_SIMPLEX",
  "elk.layered.crossingMinimization.strategy": "LAYER_SWEEP",
  "elk.spacing.nodeNode": "12",
  "elk.spacing.edgeNode": "10",
  "elk.layered.spacing.nodeNodeBetweenLayers": "150",
  "elk.layered.spacing.edgeNodeBetweenLayers": "30",
};

export async function computeLayout(
  nodes: Node[],
  edges: Edge[],
): Promise<Node[]> {
  if (nodes.length === 0) return nodes;

  const elkChildren: ElkNode[] = nodes.map((n) => {
    const dims = getNodeDims(n.type ?? "tool");
    return { id: n.id, width: dims.w, height: dims.h };
  });

  const nodeIds = new Set(nodes.map((n) => n.id));
  const elkEdges: ElkExtendedEdge[] = [];
  const seen = new Set<string>();

  for (const e of edges) {
    if (!nodeIds.has(e.source) || !nodeIds.has(e.target)) continue;
    const key = `${e.source}->${e.target}`;
    if (seen.has(key)) continue;
    seen.add(key);
    elkEdges.push({ id: e.id, sources: [e.source], targets: [e.target] });
  }

  const result = await elk.layout({
    id: "root",
    layoutOptions: ELK_OPTIONS,
    children: elkChildren,
    edges: elkEdges,
  });

  const positions = new Map<string, { x: number; y: number }>();
  for (const child of result.children ?? []) {
    positions.set(child.id, { x: child.x ?? 0, y: child.y ?? 0 });
  }

  return nodes.map((node) => {
    const pos = positions.get(node.id);
    if (pos) return { ...node, position: pos };
    return node;
  });
}
