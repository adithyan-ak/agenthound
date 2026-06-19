import { api } from "@shared/api/client";
import type { APIEdge, APINode } from "@entities/graph/dto";

export async function fetchNodes(
  kind?: string,
  limit = 10000,
): Promise<APINode[]> {
  const params: Record<string, string> = { limit: String(limit) };
  if (kind) params["kind"] = kind;
  const result = await api
    .get("graph/nodes", { searchParams: params })
    .json<APINode[] | null>();
  return result ?? [];
}

export async function fetchNode(
  id: string,
): Promise<{ node: APINode; edges: APIEdge[] }> {
  return api
    .get(`graph/nodes/${encodeURIComponent(id)}`)
    .json<{ node: APINode; edges: APIEdge[] }>();
}

export interface BlastRadiusResponse {
  nodes: APINode[];
  edges: APIEdge[];
  rings: Record<string, string[]>;
  direction: "out" | "in" | "both";
  max_hops: number;
}

export interface BlastRadiusOptions {
  direction?: "out" | "in" | "both";
  maxHops?: number;
}

export async function fetchBlastRadius(
  nodeId: string,
  opts: BlastRadiusOptions = {},
): Promise<BlastRadiusResponse> {
  const params: Record<string, string> = {};
  if (opts.direction) params["direction"] = opts.direction;
  if (opts.maxHops) params["max_hops"] = String(opts.maxHops);
  return api
    .get(`graph/nodes/${encodeURIComponent(nodeId)}/blast-radius`, {
      searchParams: params,
    })
    .json<BlastRadiusResponse>();
}
