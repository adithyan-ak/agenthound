// Scan domain types. The `isUsableScan` rule and other view-model logic are
// added in WS5; this is the type home split out of the old api/types.ts.

export type ScanStatus =
  | "pending"
  | "running"
  | "completed"
  // Collection succeeded (real node/edge counts) but analysis
  // post-processing failed; the `error` field carries the detail.
  | "completed_with_errors"
  | "failed";

export interface Scan {
  id: string;
  collector: string;
  status: ScanStatus;
  started_at: string;
  completed_at?: string;
  node_count: number;
  edge_count: number;
  error?: string;
  metadata?: Record<string, unknown>;
}

/**
 * Whether a scan actually populated the graph. `completed_with_errors` means
 * collection succeeded (real node/edge counts) even though analysis
 * post-processing failed, so it counts as usable — the single home for the
 * `completed | completed_with_errors` rule that was duplicated across widgets.
 */
export function isUsableScan(scan: Scan): boolean {
  return scan.status === "completed" || scan.status === "completed_with_errors";
}
