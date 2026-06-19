import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { qk } from "@shared/api/query-keys";
import { deleteScan, fetchScans, uploadScan } from "./api";

// Page size for the scan-manager list. The dashboard requests its own smaller
// page (20) under a distinct cache key, so writes here never disturb it.
export const SCANS_LIST_LIMIT = 50;

export function useScans(limit = SCANS_LIST_LIMIT) {
  return useQuery({
    queryKey: qk.scans(limit),
    queryFn: () => fetchScans(limit, 0),
  });
}

// Mutations replicate today's post-write refetch target ONLY: the scan-manager
// list. They deliberately do NOT invalidate the dashboard scan widgets (a
// separately-listed, out-of-scope behavioral change).
export function useDeleteScan() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteScan(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: qk.scans(SCANS_LIST_LIMIT) });
    },
  });
}

export function useUploadScan() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (file: File) => uploadScan(file),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: qk.scans(SCANS_LIST_LIMIT) });
    },
  });
}
