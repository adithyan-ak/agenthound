import { useQuery } from "@tanstack/react-query";
import { api } from "@shared/api/client";
import { qk } from "@shared/api/query-keys";

export interface HealthResponse {
  status: string;
  neo4j: string;
  postgres: string;
}

export async function fetchHealth(): Promise<HealthResponse> {
  return api.get("health").json<HealthResponse>();
}

// 30s poll is a genuine override (kept from the original inline queries);
// staleTime uses the global 30s default.
export function useHealth() {
  return useQuery({
    queryKey: qk.health(),
    queryFn: fetchHealth,
    refetchInterval: 30_000,
  });
}
