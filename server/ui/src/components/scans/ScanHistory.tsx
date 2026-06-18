import { useState, type ReactNode } from "react";
import { Trash2 } from "lucide-react";
import type { Scan } from "@/api/types";
import { deleteScan } from "@/api/scans";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { scanStatusLabel } from "@/lib/format";

interface ScanHistoryProps {
  scans: Scan[];
  onDeleted?: () => void;
}

const STATUS_COLOR: Record<string, string> = {
  completed: "#3FB950",
  completed_with_errors: "#F59E0B",
  running: "#F5A623",
  pending: "#7A828E",
  failed: "#EF4444",
};

const COLLECTOR_COLOR: Record<string, string> = {
  mcp: "#10B981",
  a2a: "#A855F7",
  config: "#D97706",
};

function formatDate(dateStr: string | undefined): string {
  if (!dateStr) return "\u2014";
  return new Date(dateStr).toLocaleString();
}

function Th({ children, className }: { children?: ReactNode; className?: string }) {
  return (
    <th
      className={cn(
        "px-3 py-2 font-mono text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground",
        className,
      )}
    >
      {children}
    </th>
  );
}

export function ScanHistory({ scans, onDeleted }: ScanHistoryProps) {
  const [confirmScan, setConfirmScan] = useState<Scan | null>(null);
  const [deleting, setDeleting] = useState(false);

  async function handleConfirmedDelete() {
    if (!confirmScan) return;
    setDeleting(true);
    try {
      await deleteScan(confirmScan.id);
      setConfirmScan(null);
      onDeleted?.();
    } finally {
      setDeleting(false);
    }
  }

  if (scans.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 font-mono text-xs uppercase tracking-[0.12em] text-muted-foreground">
        No scans recorded yet
      </div>
    );
  }

  return (
    <>
      <div className="overflow-x-auto">
        <table className="w-full border-collapse text-left">
          <thead>
            <tr className="border-b border-border bg-black/20">
              <Th className="w-10 pr-2 text-right">#</Th>
              <Th>ID</Th>
              <Th>Collector</Th>
              <Th>Status</Th>
              <Th>Started</Th>
              <Th>Completed</Th>
              <Th className="text-right">Nodes</Th>
              <Th className="text-right">Edges</Th>
              <Th className="w-10" />
            </tr>
          </thead>
          <tbody>
            {scans.map((scan, i) => {
              const statusColor = STATUS_COLOR[scan.status] ?? "#7A828E";
              const collectorColor = COLLECTOR_COLOR[scan.collector] ?? "#7A828E";
              const running = scan.status === "running";
              return (
                <tr
                  key={`${scan.id}-${scan.collector}`}
                  className="border-b border-border/60 transition-colors last:border-0 hover:bg-white/[0.03]"
                >
                  <td
                    className="px-3 py-2.5 text-right align-middle font-mono text-[10px] tabular-nums text-muted-foreground/60"
                    style={{ boxShadow: `inset 2px 0 0 0 ${statusColor}` }}
                  >
                    {String(i + 1).padStart(2, "0")}
                  </td>
                  <td className="px-3 py-2.5 align-middle font-mono text-[11px] text-foreground/80">
                    {scan.id.slice(0, 8)}
                  </td>
                  <td className="px-3 py-2.5 align-middle">
                    <span className="inline-flex items-center gap-1.5 rounded-[2px] border border-border bg-black/40 px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-[0.06em] text-muted-foreground">
                      <span className="h-1.5 w-1.5 rounded-[1px]" style={{ backgroundColor: collectorColor }} />
                      {scan.collector}
                    </span>
                  </td>
                  <td className="px-3 py-2.5 align-middle">
                    <span
                      className={cn("inline-flex items-center gap-1.5", scan.error && "cursor-help")}
                      title={scan.error || undefined}
                    >
                      <span
                        className={cn("h-2 w-2 rounded-[1px]", running && "animate-led-pulse")}
                        style={{ backgroundColor: statusColor, boxShadow: `0 0 6px -1px ${statusColor}` }}
                      />
                      <span
                        className="font-mono text-[10px] font-semibold uppercase tracking-[0.08em]"
                        style={{ color: statusColor }}
                      >
                        {scanStatusLabel(scan.status)}
                      </span>
                    </span>
                  </td>
                  <td className="px-3 py-2.5 align-middle font-mono text-[11px] text-muted-foreground">
                    {formatDate(scan.started_at)}
                  </td>
                  <td className="px-3 py-2.5 align-middle font-mono text-[11px] text-muted-foreground">
                    {formatDate(scan.completed_at)}
                  </td>
                  <td className="px-3 py-2.5 text-right align-middle font-mono text-[11px] tabular-nums text-foreground">
                    {scan.node_count}
                  </td>
                  <td className="px-3 py-2.5 text-right align-middle font-mono text-[11px] tabular-nums text-foreground">
                    {scan.edge_count}
                  </td>
                  <td className="px-3 py-2.5 text-right align-middle">
                    <button
                      onClick={() => setConfirmScan(scan)}
                      title="Delete scan"
                      className="inline-flex h-7 w-7 items-center justify-center rounded-[3px] text-muted-foreground transition-colors hover:bg-white/[0.05] hover:text-destructive"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <Dialog open={!!confirmScan} onOpenChange={(v) => !v && setConfirmScan(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 font-mono uppercase tracking-[0.04em]">
              <span className="h-2 w-2 rounded-[1px] bg-destructive" />
              Delete scan?
            </DialogTitle>
            <DialogDescription>
              This will permanently delete the scan record and remove all
              {confirmScan ? ` ${confirmScan.node_count} nodes and ${confirmScan.edge_count} edges` : ""}{" "}
              it contributed to the graph. Nodes shared with other scans will be preserved.
            </DialogDescription>
          </DialogHeader>
          <div className="mt-2 flex justify-end gap-2">
            <button
              onClick={() => setConfirmScan(null)}
              disabled={deleting}
              className="inline-flex h-8 items-center rounded-[3px] border border-border bg-black/30 px-3 font-mono text-[11px] uppercase tracking-[0.08em] text-foreground/80 transition-colors hover:border-mauve-7 hover:text-foreground disabled:opacity-40"
            >
              Cancel
            </button>
            <button
              onClick={handleConfirmedDelete}
              disabled={deleting}
              className="inline-flex h-8 items-center rounded-[3px] bg-destructive px-3 font-mono text-[11px] font-semibold uppercase tracking-[0.08em] text-destructive-foreground transition-colors hover:bg-destructive/90 disabled:opacity-40"
            >
              {deleting ? "Deleting…" : "Delete scan"}
            </button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
