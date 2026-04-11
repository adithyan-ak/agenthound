import { useMemo } from "react";
import { useExplorerStore } from "@/store/explorer";
import { useExplorerGraph } from "@/hooks/useExplorerGraph";
import { getLens } from "@/lib/explorer/lens-config";
import { buildExplorerGraph } from "@/lib/explorer/graph-builder";
import { cn } from "@/lib/utils";

export function InfoCard() {
  const { data } = useExplorerGraph();
  const activeLens = useExplorerStore((s) => s.activeLens);
  const subPresets = useExplorerStore((s) => s.subPresets[activeLens] ?? []);

  const metrics = useMemo(() => {
    if (!data) return null;
    const lens = getLens(activeLens);
    const built = buildExplorerGraph(
      { nodes: data.nodes, edges: data.edges },
      {
        lens,
        activeLensId: activeLens,
        subPresets,
        findings: data.findings,
      },
    );
    return built.metrics;
  }, [data, activeLens, subPresets]);

  const lens = getLens(activeLens);

  if (!metrics) return null;

  return (
    <div
      className={cn(
        "pointer-events-auto absolute left-6 top-24 z-20 w-[280px] rounded-xl",
        "border border-slate-800/90 bg-slate-950/92 p-4 shadow-2xl backdrop-blur-md",
      )}
      style={{
        borderTopColor: lens.activeTint,
        borderTopWidth: 3,
      }}
    >
      <div className="flex items-center justify-between">
        <div className="text-[10px] font-semibold uppercase tracking-widest text-slate-500">
          {lens.label} lens
        </div>
        <div className="flex h-5 items-center gap-1 rounded-full bg-slate-900 px-2 text-[9px] text-slate-400">
          <div
            className="h-1.5 w-1.5 rounded-full animate-pulse"
            style={{ background: lens.activeTint }}
          />
          active
        </div>
      </div>

      <div className="mt-2 flex items-baseline gap-1.5">
        <div className="text-2xl font-bold text-white tabular-nums">
          {metrics.visibleEdgeCount}
        </div>
        <div className="text-xs text-slate-400">visible edges</div>
      </div>
      <div className="mt-0.5 text-[10px] text-slate-500">
        across {metrics.visibleNodeCount} nodes
      </div>

      <div className="mt-3 space-y-1">
        {metrics.criticalCount > 0 && (
          <MetricRow
            label="Critical"
            value={metrics.criticalCount}
            color="#EF4444"
          />
        )}
        {metrics.highCount > 0 && (
          <MetricRow
            label="High"
            value={metrics.highCount}
            color="#F97316"
          />
        )}
        {metrics.mediumCount > 0 && (
          <MetricRow
            label="Medium"
            value={metrics.mediumCount}
            color="#EAB308"
          />
        )}
        {metrics.lowCount > 0 && (
          <MetricRow
            label="Low"
            value={metrics.lowCount}
            color="#94A3B8"
          />
        )}
      </div>

      <div className="mt-3 border-t border-slate-800/80 pt-2 text-[10px] leading-relaxed text-slate-400">
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
      <div
        className="h-1.5 w-1.5 rounded-full"
        style={{ background: color }}
      />
      <div className="flex-1 text-[11px] text-slate-300">{label}</div>
      <div className="text-[11px] font-semibold tabular-nums text-white">
        {value}
      </div>
    </div>
  );
}
