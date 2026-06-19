import { api } from "@shared/api/client";
import type { Finding, FindingDetail, TriageState } from "./model";

export async function fetchFindings(
  severity?: string,
  includeSuppressed?: boolean,
): Promise<Finding[]> {
  const params: Record<string, string> = {};
  if (severity) params["severity"] = severity;
  if (includeSuppressed) params["include_suppressed"] = "true";
  return api
    .get("analysis/findings", { searchParams: params })
    .json<Finding[]>();
}

export async function fetchFindingDetail(id: string): Promise<FindingDetail> {
  return api.get(`analysis/findings/${id}`).json<FindingDetail>();
}

export async function getTriage(fingerprint: string): Promise<TriageState> {
  return api.get(`findings/triage/${fingerprint}`).json<TriageState>();
}

export async function setTriage(
  fingerprint: string,
  status: string,
  note: string,
): Promise<TriageState> {
  return api
    .put(`findings/triage/${fingerprint}`, { json: { status, note } })
    .json<TriageState>();
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
        .json<Finding[]>(),
    ),
  );
  return results.flat();
}
