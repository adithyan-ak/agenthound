import { useQuery } from "@tanstack/react-query";
import { qk } from "@shared/api/query-keys";
import { fetchFindingDetail, fetchFindings } from "./api";

// One findings cache for every surface that needs the full set (dashboard,
// findings list, node findings, references, navigation).
export function useFindings() {
  return useQuery({
    queryKey: qk.findings(),
    queryFn: () => fetchFindings(),
  });
}

export function useFindingDetail(findingId: string | undefined) {
  return useQuery({
    queryKey: qk.findingDetail(findingId ?? ""),
    queryFn: () => fetchFindingDetail(findingId!),
    enabled: !!findingId,
  });
}
