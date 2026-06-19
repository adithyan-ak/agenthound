import { api } from "@shared/api/client";
import type { Finding, FindingDetail } from "./model";

export async function fetchFindings(severity?: string): Promise<Finding[]> {
  const params: Record<string, string> = {};
  if (severity) params["severity"] = severity;
  return api
    .get("analysis/findings", { searchParams: params })
    .json<Finding[]>();
}

export async function fetchFindingDetail(id: string): Promise<FindingDetail> {
  return api.get(`analysis/findings/${id}`).json<FindingDetail>();
}

/**
 * Fetch findings across all severities in a single call by fanning out
 * parallel requests (the backend only filters one severity at a time) and
 * flattening. Used by the explorer's bundled graph fetch.
 */
export async function fetchAllFindings(): Promise<Finding[]> {
  const severities = ["critical", "high", "medium", "low"];
  const results = await Promise.all(
    severities.map((sev) =>
      api
        .get("analysis/findings", { searchParams: { severity: sev } })
        .json<Finding[]>()
        .catch(() => [] as Finding[]),
    ),
  );
  return results.flat();
}
