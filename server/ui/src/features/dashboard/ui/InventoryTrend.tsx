import { useMemo } from "react";
import { TrendingUp } from "lucide-react";
import { useScans, isUsableScan } from "@entities/scan";
import { Skeleton } from "@shared/ui/primitives/skeleton";
import { AsyncBoundary } from "@shared/ui/feedback";
import { WidgetCard, AreaTrend } from "@shared/ui/widgets";
import type { TrendSeries } from "@shared/ui/widgets";
import { ACCENT, INSTRUMENT } from "@shared/theme/tokens";
import { shortDate } from "@shared/lib/format";

const INFO =
  "Nodes and edges discovered per scan over time, drawn from the scan history — your attack surface's growth.";

const SERIES: TrendSeries[] = [
  { key: "nodes", label: "Nodes", color: ACCENT },
  { key: "edges", label: "Edges", color: INSTRUMENT.grayMuted },
];

export function InventoryTrend() {
  const { data: scans, isLoading } = useScans(20);

  const data = useMemo(() => {
    // Include completed_with_errors: the graph was populated, so the
    // node/edge counts are real and part of the surface growth trend.
    const usable = (scans ?? []).filter(isUsableScan);
    return usable
      .slice()
      .reverse()
      .map((s) => ({
        t: shortDate(s.started_at),
        nodes: s.node_count,
        edges: s.edge_count,
      }));
  }, [scans]);

  return (
    <WidgetCard
      title="Surface Growth"
      info={INFO}
      icon={TrendingUp}
      accent={ACCENT}
      action={
        <div className="flex items-center gap-3">
          {SERIES.map((s) => (
            <span
              key={s.key}
              className="flex items-center gap-1.5 font-mono text-[9px] uppercase tracking-[0.1em] text-muted-foreground"
            >
              <span className="h-2 w-2 rounded-[1px]" style={{ backgroundColor: s.color }} />
              {s.label}
            </span>
          ))}
        </div>
      }
    >
      <AsyncBoundary
        isLoading={isLoading}
        isEmpty={data.length === 0}
        loading={<Skeleton className="h-44 w-full" />}
        empty={
          <div className="flex h-44 items-center justify-center font-mono text-xs uppercase tracking-wider text-muted-foreground">
            No completed scans yet
          </div>
        }
      >
        <AreaTrend data={data} series={SERIES} xKey="t" height={176} />
      </AsyncBoundary>
    </WidgetCard>
  );
}
