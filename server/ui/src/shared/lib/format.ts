/** Compact relative time, e.g. "12s ago", "4m ago", "3h ago", "2d ago". */
export function timeAgo(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const seconds = Math.floor(diff / 1000);
  if (seconds < 60) return `${Math.max(seconds, 0)}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}

/** Short clock label for chart ticks, e.g. "Jun 18 14:05". */
export function shortDateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

/** Short date label, e.g. "Jun 18". */
export function shortDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric" });
}

/**
 * Human-friendly label for a scan status. Statuses render uppercased in the
 * UI, so the multi-word value needs spaces instead of the raw underscored
 * enum (e.g. "completed_with_errors" -> "Completed with errors").
 */
export function scanStatusLabel(status: string): string {
  if (status === "completed_with_errors") return "Completed with errors";
  return status;
}
