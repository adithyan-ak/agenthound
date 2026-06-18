import { useQuery } from "@tanstack/react-query";
import { fetchNodes } from "@/api/graph";
import { fetchFindings } from "@/api/analysis";
import { fetchScans } from "@/api/scans";

const STALE = 30_000;

/**
 * Shared dashboard queries. Co-locating the query keys here lets every widget
 * derive from the same cached datasets (one node fetch, one findings fetch,
 * one scans fetch) instead of each panel re-requesting overlapping data.
 */

export function useAllNodes() {
  return useQuery({
    queryKey: ["dashboard", "all-nodes"],
    queryFn: () => fetchNodes(undefined, 10000),
    staleTime: STALE,
  });
}

export function useDashboardFindings() {
  return useQuery({
    queryKey: ["dashboard", "findings"],
    queryFn: () => fetchFindings(),
    staleTime: STALE,
  });
}

export function useDashboardScans(limit = 20) {
  return useQuery({
    queryKey: ["dashboard", "scans", limit],
    queryFn: () => fetchScans(limit, 0),
    staleTime: STALE,
  });
}
