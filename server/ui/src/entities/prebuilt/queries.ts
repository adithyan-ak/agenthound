import { useMutation, useQuery } from "@tanstack/react-query";
import { qk } from "@shared/api/query-keys";
import { fetchPreBuiltQueries, runPreBuiltQuery } from "./api";

export function usePreBuiltQueries() {
  return useQuery({
    queryKey: qk.prebuiltQueries(),
    queryFn: fetchPreBuiltQueries,
  });
}

// Cached result of a named pre-built query (dashboard cross-protocol /
// chokepoint widgets read this).
export function usePreBuiltResult(id: string) {
  return useQuery({
    queryKey: qk.prebuiltResult(id),
    queryFn: () => runPreBuiltQuery(id),
  });
}

// On-demand query execution (Query Library). No cache invalidation — it only
// produces rows for display, matching today's behavior.
export function useRunPreBuiltQuery() {
  return useMutation({
    mutationFn: (id: string) => runPreBuiltQuery(id),
  });
}
