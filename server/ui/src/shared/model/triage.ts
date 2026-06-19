import { SEVERITY, ACCENT, SIGNAL_OK, TRIAGE_NEUTRAL } from "@shared/theme/tokens";

/**
 * Finding triage workflow types + view-model constants (shared kernel).
 *
 * Triage state is now server-backed (Postgres finding_triage table, keyed by
 * the finding fingerprint) and read/written via TanStack Query — see
 * `@entities/finding` useTriage / useSetTriage. This module keeps only the
 * status enum and its display metadata so consumers (the register, the
 * dossier header) stay decoupled from the transport.
 *
 * "new" is the implicit default for a finding with no recorded decision.
 */
export type TriageStatus =
  | "new"
  | "triaging"
  | "confirmed"
  | "accepted-risk"
  | "false-positive";

export const TRIAGE_ORDER: TriageStatus[] = [
  "new",
  "triaging",
  "confirmed",
  "accepted-risk",
  "false-positive",
];

export interface TriageMeta {
  label: string;
  short: string;
  color: string;
}

export const TRIAGE_META: Record<TriageStatus, TriageMeta> = {
  new: { label: "New", short: "NEW", color: TRIAGE_NEUTRAL },
  triaging: { label: "Triaging", short: "TRIAGE", color: ACCENT },
  confirmed: { label: "Confirmed", short: "CONFIRMED", color: SEVERITY.critical.solid },
  "accepted-risk": { label: "Accepted risk", short: "ACCEPTED", color: SEVERITY.info.solid },
  "false-positive": { label: "False positive", short: "FALSE+", color: SIGNAL_OK },
};
