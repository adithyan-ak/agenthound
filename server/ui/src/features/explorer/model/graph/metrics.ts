import type { APINode } from "@entities/graph/dto";
import type { LensMetrics, LogicalEdge, LogicalHexNode } from "./types";

/**
 * Derive the lens metrics from the built hex nodes + edges. Visible counts use
 * the non-dimmed subset; cluster nodes are intentionally excluded (only hexes
 * are counted), matching the original behavior.
 */
export function computeMetrics(
  hexNodes: LogicalHexNode[],
  edges: LogicalEdge[],
  orphanCount: number,
  orphanByKind: Record<string, APINode[]>,
): LensMetrics {
  const orphanByKindCounts: Record<string, number> = {};
  for (const [kind, members] of Object.entries(orphanByKind)) {
    orphanByKindCounts[kind] = members.length;
  }

  const metrics: LensMetrics = {
    visibleNodeCount: hexNodes.filter((n) => !n.data.dim).length,
    visibleEdgeCount: edges.filter((e) => !e.data.dim).length,
    criticalCount: 0,
    highCount: 0,
    mediumCount: 0,
    lowCount: 0,
    orphanCount,
    orphanByKind: orphanByKindCounts,
  };
  for (const e of edges) {
    const sev = e.data.severity;
    if (sev === "critical") metrics.criticalCount++;
    else if (sev === "high") metrics.highCount++;
    else if (sev === "medium") metrics.mediumCount++;
    else if (sev === "low") metrics.lowCount++;
  }
  return metrics;
}
