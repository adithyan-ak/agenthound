import { useQuery } from "@tanstack/react-query";
import { qk } from "@shared/api/query-keys";
import { fetchNodes } from "@entities/node/api";
import { fetchEdges } from "@entities/edge/api";
import { fetchAllFindings } from "@entities/finding/api";
import type { APIEdge, APINode } from "@entities/graph/dto";
import type { Finding } from "@entities/finding/model";

export interface ExplorerRawData {
  nodes: APINode[];
  edges: APIEdge[];
  findings: Finding[];
}

/**
 * Fetches the full graph (all nodes + all edges + all findings) in one call.
 * Lens switching filters this data client-side with no extra round-trips.
 * staleTime is the global 30s default.
 */
export function useExplorerGraph() {
  return useQuery({
    queryKey: qk.explorerGraph(),
    queryFn: async (): Promise<ExplorerRawData> => {
      const [nodes, edges, findings] = await Promise.all([
        fetchNodes(undefined, 10000),
        fetchEdges(undefined, 100000),
        fetchAllFindings(),
      ]);
      return { nodes, edges, findings };
    },
  });
}
