import { useQuery } from "@tanstack/react-query";
import { qk } from "@shared/api/query-keys";
import { fetchNode, fetchNodes } from "./api";

export function useNodes(kind?: string, limit = 10000) {
  return useQuery({
    queryKey: qk.nodes(kind, limit),
    queryFn: () => fetchNodes(kind, limit),
  });
}

export function useNode(id: string | null) {
  return useQuery({
    queryKey: qk.node(id ?? ""),
    queryFn: () => fetchNode(id!),
    enabled: id !== null,
  });
}
