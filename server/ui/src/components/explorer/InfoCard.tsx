import { useMemo } from "react";
import { EyeOff, Eye } from "lucide-react";
import { useExplorerStore } from "@/store/explorer";
import { useExplorerGraph } from "@/hooks/useExplorerGraph";
import { useBlastRadius } from "@/hooks/useBlastRadius";
import { getLens } from "@/lib/explorer/lens-config";
import { buildExplorerGraph } from "@/lib/explorer/graph-builder";
import { SEVERITY } from "@/theme/tokens";
import { cn } from "@/lib/utils";

export function InfoCard() {
  const { data } = useExplorerGraph();
  const activeLens = useExplorerStore((s) => s.activeLens);
  const subPresets = useExplorerStore((s) => s.subPresets[activeLens] ?? []);
  const showOrphans = useExplorerStore((s) => s.showOrphans);
  const toggleShowOrphans = useExplorerStore((s) => s.toggleShowOrphans);
  const blastRadiusSourceId = useExplorerStore((s) => s.blastRadiusSourceId);
  const blastDirection = useExplorerStore((s) => s.blastRadiusDirection);
  const blastMaxHops = useExplorerStore((s) => s.blastRadiusMaxHops);

  const { data: blastData } = useBlastRadius(
    activeLens === "blast-radius" ? blastRadiusSourceId : null,
    blastDirection,
    blastMaxHops,
  );

  const metrics = useMemo(() => {
    if (!data) return null;
    const lens = getLens(activeLens);

    const blastRadius =
      activeLens === "blast-radius" && blastData && blastRadiusSourceId
        ? {
            sourceId: blastRadiusSourceId,
            nodeIds: blastData.nodeIdSet,
            edgeKeys: blastData.edgeKeySet,
          }
        : undefined;

    const built = buildExplorerGraph(
      { nodes: data.nodes, edges: data.edges },
      {
        lens,
        activeLensId: activeLens,
        subPresets,
        findings: data.findings,
        blastRadius,
        showOrphans,
      },
    );
    return built.metrics;
  }, [data, activeLens, subPresets, showOrphans, blastData, blastRadiusSourceId]);

  const lens = getLens(activeLens);

  if (!metrics) return null;

  return (
    <div
      className={cn(
        "pointer-events-auto absolute left-4 top-20 z-20 w-[280px] overflow-hidden rounded-md",
        "border border-border bg-card/95 p-3.5 backdrop-blur-md elev-2",
      )}
    >
      <span aria-hidden className="pointer-events-none absolute inset-x-0 top-0 h-px bg-white/[0.05]" />
      <span
        aria-hidden
        className="pointer-events-none absolute left-0 top-0 h-px w-12"
        style={{ background: lens.activeTint, opacity: 0.9 }}
      />
      <div className="flex items-center justify-between">
        <div className="font-mono text-[10px] font-semibold uppercase tracking-[0.16em] text-muted-foreground">
          {lens.label} lens
        </div>
        <div className="flex items-center gap-1 rounded-[2px] border border-border bg-black/30 px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-[0.1em] text-muted-foreground">
          <span
            className="h-1.5 w-1.5 animate-led-pulse rounded-[1px]"
            style={{ background: lens.activeTint }}
          />
          active
        </div>
      </div>

      <div className="mt-2.5 flex items-baseline gap-1.5">
        <div className="font-mono text-[28px] font-bold leading-none tabular-nums text-foreground">
          {metrics.visibleEdgeCount}
        </div>
        <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-muted-foreground">
          visible edges
        </div>
      </div>
      <div className="mt-1 font-mono text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
        across {metrics.visibleNodeCount} nodes
      </div>

      <div className="mt-3 space-y-1">
        {metrics.criticalCount > 0 && (
          <MetricRow
            label="Critical"
            value={metrics.criticalCount}
            color={SEVERITY.critical.solid}
          />
        )}
        {metrics.highCount > 0 && (
          <MetricRow
            label="High"
            value={metrics.highCount}
            color={SEVERITY.high.solid}
          />
        )}
        {metrics.mediumCount > 0 && (
          <MetricRow
            label="Medium"
            value={metrics.mediumCount}
            color={SEVERITY.medium.solid}
          />
        )}
        {metrics.lowCount > 0 && (
          <MetricRow
            label="Low"
            value={metrics.lowCount}
            color={SEVERITY.low.solid}
          />
        )}
      </div>

      {metrics.orphanCount > 0 && !lens.dimOthers && (
        <button
          onClick={toggleShowOrphans}
          className={cn(
            "mt-3 flex w-full items-center justify-between rounded-[3px] border px-2.5 py-1.5 font-mono text-[10px] uppercase tracking-[0.06em] transition-colors",
            showOrphans
              ? "border-primary/40 bg-primary/10 text-primary hover:bg-primary/15"
              : "border-border bg-black/30 text-muted-foreground hover:border-mauve-7 hover:text-foreground",
          )}
          aria-label={
            showOrphans
              ? "Hide unconnected node clusters"
              : "Show unconnected node clusters"
          }
        >
          <span className="flex items-center gap-1.5">
            {showOrphans ? (
              <Eye className="h-3 w-3" strokeWidth={2.25} />
            ) : (
              <EyeOff className="h-3 w-3" strokeWidth={2.25} />
            )}
            <span className="tabular-nums">{metrics.orphanCount}</span>
            <span>unconnected</span>
          </span>
          <span className="font-semibold text-[9px] uppercase tracking-widest">
            {showOrphans ? "hide" : "show clusters"}
          </span>
        </button>
      )}

      <div className="mt-3 border-t border-border/70 pt-2 text-[10px] leading-relaxed text-muted-foreground">
        {lens.description}
      </div>
    </div>
  );
}

function MetricRow({
  label,
  value,
  color,
}: {
  label: string;
  value: number;
  color: string;
}) {
  return (
    <div className="flex items-center gap-2">
      <span className="h-1.5 w-1.5 rounded-[1px]" style={{ background: color }} />
      <div className="flex-1 font-mono text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
        {label}
      </div>
      <div className="font-mono text-[11px] font-semibold tabular-nums" style={{ color }}>
        {value}
      </div>
    </div>
  );
}
