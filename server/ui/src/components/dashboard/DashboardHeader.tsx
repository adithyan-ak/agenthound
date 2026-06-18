import { useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Activity, Download, RefreshCw, Radar } from "lucide-react";
import { api } from "@/api/client";
import { fetchScans } from "@/api/scans";
import { fetchGraphStats } from "@/api/graph";
import { fetchFindings } from "@/api/analysis";
import type { HealthResponse } from "@/api/types";
import { useDashboardScans } from "@/hooks/useDashboardData";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { timeAgo } from "@/lib/format";
import { SEVERITY_ORDER } from "@/theme/tokens";

function greeting(): string {
  const h = new Date().getHours();
  if (h < 12) return "Good morning";
  if (h < 18) return "Good afternoon";
  return "Good evening";
}

const today = new Date().toLocaleDateString(undefined, {
  weekday: "long",
  month: "short",
  day: "numeric",
});

function downloadJSON(name: string, data: unknown) {
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = name;
  a.click();
  URL.revokeObjectURL(url);
}

export function DashboardHeader() {
  const queryClient = useQueryClient();
  const [refreshing, setRefreshing] = useState(false);
  const [exporting, setExporting] = useState(false);

  const { data: health } = useQuery({
    queryKey: ["health"],
    queryFn: () => api.get("health").json<HealthResponse>(),
    refetchInterval: 30_000,
  });

  const { data: scans } = useDashboardScans();

  const isHealthy = health?.status === "healthy";
  const lastCompleted = (scans ?? []).find((s) => s.status === "completed");
  const running = (scans ?? []).some((s) => s.status === "running");

  async function refresh() {
    setRefreshing(true);
    try {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
        queryClient.invalidateQueries({ queryKey: ["graph"] }),
        queryClient.invalidateQueries({ queryKey: ["health"] }),
      ]);
    } finally {
      setRefreshing(false);
    }
  }

  async function exportSnapshot() {
    setExporting(true);
    try {
      const [stats, findings, scanList] = await Promise.all([
        fetchGraphStats(),
        fetchFindings(),
        fetchScans(20, 0),
      ]);
      const bySeverity: Record<string, number> = {};
      for (const sev of SEVERITY_ORDER) bySeverity[sev] = 0;
      for (const f of findings) bySeverity[f.severity] = (bySeverity[f.severity] ?? 0) + 1;
      downloadJSON(`agenthound-attack-surface-${new Date().toISOString().slice(0, 10)}.json`, {
        generated_at: new Date().toISOString(),
        totals: {
          nodes: stats.total_nodes,
          edges: stats.total_edges,
          findings: findings.length,
        },
        node_counts: stats.node_counts,
        edge_counts: stats.edge_counts,
        findings_by_severity: bySeverity,
        recent_scans: scanList,
      });
    } finally {
      setExporting(false);
    }
  }

  return (
    <header className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
      <div className="min-w-0">
        <p className="text-[11px] font-semibold uppercase tracking-[0.14em] text-muted-foreground">
          {greeting()} &middot; {today}
        </p>
        <h1 className="mt-1 flex items-center gap-2.5 text-2xl font-bold tracking-tight text-foreground sm:text-[28px]">
          <Radar className="h-6 w-6 text-primary" />
          Attack Surface Command
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Live security posture across your agent, MCP, and A2A infrastructure.
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <div className="flex items-center gap-2 rounded-lg border border-white/[0.07] bg-white/[0.03] px-3 py-1.5">
          <Activity className="h-3.5 w-3.5 text-muted-foreground" />
          <span className="text-xs text-muted-foreground">
            {running ? (
              <span className="text-amber-400">Scan in progress…</span>
            ) : lastCompleted ? (
              <>Last scan {timeAgo(lastCompleted.started_at)}</>
            ) : (
              "No scans yet"
            )}
          </span>
        </div>

        <div className="flex items-center gap-2 rounded-lg border border-white/[0.07] bg-white/[0.03] px-3 py-1.5">
          <span
            className={cn(
              "h-2 w-2 rounded-full",
              isHealthy ? "bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.8)]" : "bg-destructive",
            )}
          />
          <span className="text-xs text-muted-foreground">
            {isHealthy ? "Operational" : "Degraded"}
          </span>
        </div>

        <Button variant="outline" size="sm" onClick={refresh} disabled={refreshing} className="h-9">
          <RefreshCw className={cn("h-4 w-4", refreshing && "animate-spin")} />
          Refresh
        </Button>
        <Button variant="outline" size="sm" onClick={exportSnapshot} disabled={exporting} className="h-9">
          <Download className="h-4 w-4" />
          Export
        </Button>
      </div>
    </header>
  );
}
