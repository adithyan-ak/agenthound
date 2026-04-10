import ELK, { type ElkNode, type ElkExtendedEdge } from "elkjs/lib/elk.bundled.js";
import type { Node, Edge } from "@xyflow/react";

const elk = new ELK();

const NODE_W = 200;
const NODE_H = 50;

const ELK_OPTIONS: Record<string, string> = {
  "elk.algorithm": "layered",
  "elk.direction": "RIGHT",
  "elk.layered.layering.strategy": "NETWORK_SIMPLEX",
  "elk.layered.nodePlacement.strategy": "NETWORK_SIMPLEX",
  "elk.layered.crossingMinimization.strategy": "LAYER_SWEEP",
  "elk.spacing.nodeNode": "35",
  "elk.spacing.edgeNode": "25",
  "elk.layered.spacing.nodeNodeBetweenLayers": "120",
  "elk.layered.spacing.edgeNodeBetweenLayers": "40",
};

export async function computeLayout(
  nodes: Node[],
  edges: Edge[],
): Promise<Node[]> {
  if (nodes.length === 0) return nodes;

  const elkChildren: ElkNode[] = nodes.map((n) => ({
    id: n.id,
    width: NODE_W,
    height: NODE_H,
  }));

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
