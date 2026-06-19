import { useQuery } from "@tanstack/react-query";
import { api } from "@shared/api/client";
import { qk } from "@shared/api/query-keys";

export interface GraphStats {
  node_counts: Record<string, number>;
  edge_counts: Record<string, number>;
  total_nodes: number;
  total_edges: number;
}

export async function fetchGraphStats(): Promise<GraphStats> {
  return api.get("graph/stats").json<GraphStats>();
}

export function useGraphStats() {
  return useQuery({
    queryKey: qk.graphStats(),
    queryFn: fetchGraphStats,
  });
}
