import { riskColor } from "@shared/theme/tokens";

export interface ExposureInput {
  critical: number;
  high: number;
  unauthServers: number;
}

/**
 * Composite exposure index = (critical x8) + (high x3) + (unauth servers x5),
 * capped at 100. Byte-identical to the formula previously inlined in both the
 * dashboard header strip and the exposure gauge.
 */
export function exposureScore({
  critical,
  high,
  unauthServers,
}: ExposureInput): number {
  return Math.min(100, critical * 8 + high * 3 + unauthServers * 5);
}

export type ExposureBand = "critical" | "elevated" | "guarded" | "low";

/**
 * Threshold -> band. Each call site maps the band to its own label phrasing
 * ("Critical" in the header strip vs "Critical Risk" in the gauge), so only the
 * band id lives here — not the display string.
 */
export function exposureBand(score: number): ExposureBand {
  if (score >= 75) return "critical";
  if (score >= 50) return "elevated";
  if (score >= 25) return "guarded";
  return "low";
}

/**
 * Band color. Tracks riskColor() exactly (same thresholds), imported from the
 * theme so no hex leaves tokens.ts.
 */
export function exposureColor(score: number): string {
  return riskColor(score);
}
