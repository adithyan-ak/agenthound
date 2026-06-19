import type { APIEdge, APINode } from "@entities/graph/dto";
import type { BuildOptions, BuildResult } from "./types";
import { buildFindingIndex, buildLogicalEdges } from "./build-edges";
import { buildLogicalNodes } from "./build-nodes";
import { buildOrphanClusters } from "./clustering";
import { computeMetrics } from "./metrics";
import {
  logicalClusterNodeToReactFlow,
  logicalEdgeToReactFlow,
  logicalHexNodeToReactFlow,
} from "./to-react-flow";

/**
 * Pure orchestrator: transform raw API data + active lens into React Flow nodes
 * and edges ready for rendering. Composes the edge filter/bundling, node build,
 * orphan clustering, metrics, and the React-Flow adapter. Output is identical
 * to the previous monolithic builder: orphan-cluster nodes are prepended to the
 * hex nodes, edges follow bundle order.
 */
export function buildExplorerGraph(
  raw: { nodes: APINode[]; edges: APIEdge[] },
  opts: BuildOptions,
): BuildResult {
  const findingIndex = buildFindingIndex(opts.findings);

  // Map node kinds by ID for fast source/target lookup.
  const nodeById = new Map<string, APINode>();
  for (const n of raw.nodes) nodeById.set(n.id, n);

  const { edges: logicalEdges, touchedNodeIds } = buildLogicalEdges(
    raw,
    opts,
    findingIndex,
    nodeById,
  );

  const { hexNodes, orphanByKind, orphanCount } = buildLogicalNodes(
    raw,
    opts,
    logicalEdges,
    touchedNodeIds,
  );

  const clusterNodes = buildOrphanClusters(orphanByKind, opts, orphanCount);

  const metrics = computeMetrics(
    hexNodes,
    logicalEdges,
    orphanCount,
    orphanByKind,
  );

  // Cluster nodes are prepended so ELK places them above connected nodes.
  const nodes = [
    ...clusterNodes.map(logicalClusterNodeToReactFlow),
    ...hexNodes.map(logicalHexNodeToReactFlow),
  ];
  const edges = logicalEdges.map(logicalEdgeToReactFlow);

  return { nodes, edges, metrics };
}
