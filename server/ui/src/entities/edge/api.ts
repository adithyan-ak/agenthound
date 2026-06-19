import { useQuery } from "@tanstack/react-query";
import { api } from "@shared/api/client";
import { qk } from "@shared/api/query-keys";
import type { APIEdge } from "@entities/graph/dto";

export async function fetchEdges(
  kind?: string,
  limit = 100000,
): Promise<APIEdge[]> {
  const params: Record<string, string> = { limit: String(limit) };
  if (kind) params["kind"] = kind;
  const result = await api
    .get("graph/edges", { searchParams: params })
    .json<APIEdge[] | null>();
  return result ?? [];
}

// Single "all edges" cache (the inspector pulls the full set and filters
// client-side). `enabled` gates the fetch when there is nothing to inspect.
export function useEdges(enabled = true) {
  return useQuery({
    queryKey: qk.edges(),
    queryFn: () => fetchEdges(undefined, 100000),
    enabled,
  });
}
