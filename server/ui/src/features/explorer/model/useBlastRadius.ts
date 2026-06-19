import { useQuery } from "@tanstack/react-query";
import { qk } from "@shared/api/query-keys";
import { fetchBlastRadius, type BlastRadiusResponse } from "@entities/node/api";
import type { BlastDirection } from "./store";

export interface BlastRadiusData extends BlastRadiusResponse {
  nodeIdSet: Set<string>;
  edgeKeySet: Set<string>;
}

/**
 * Fetch the blast radius subgraph for a given source node, augmented with
 * Set<string> indices for O(1) membership checks during rendering.
 * staleTime is the global 30s default.
 */
export function useBlastRadius(
  sourceId: string | null,
  direction: BlastDirection,
  maxHops: number,
) {
  return useQuery({
    queryKey: qk.blastRadius(sourceId, direction, maxHops),
    queryFn: async (): Promise<BlastRadiusData> => {
      if (!sourceId) throw new Error("no source id");
      const raw = await fetchBlastRadius(sourceId, {
        direction,
        maxHops,
      });
      const nodeIdSet = new Set(raw.nodes.map((n) => n.id));
      const edgeKeySet = new Set(
        raw.edges.map((e) => `${e.source}|${e.target}|${e.kind}`),
      );
      return { ...raw, nodeIdSet, edgeKeySet };
    },
    enabled: sourceId !== null,
  });
}
