import { useMemo } from "react";
import { TrendingUp } from "lucide-react";
import { useDashboardScans } from "@/hooks/useDashboardData";
import { Skeleton } from "@/components/ui/skeleton";
import { WidgetCard, AreaTrend } from "./kit";
import type { TrendSeries } from "./kit";
import { NODE_KIND_COLORS } from "@/theme/tokens";
import { shortDate } from "@/lib/format";

const INFO =
  "Nodes and edges discovered per scan over time, drawn from the scan history — your attack surface's growth.";

const SERIES: TrendSeries[] = [
  { key: "nodes", label: "Nodes", color: NODE_KIND_COLORS.AgentInstance ?? "#06B6D4" },
  { key: "edges", label: "Edges", color: NODE_KIND_COLORS.A2AAgent ?? "#A855F7" },
];

export function InventoryTrend() {
  const { data: scans, isLoading } = useDashboardScans();

  const data = useMemo(() => {
    const usable = (scans ?? []).filter((s) => s.status === "completed");
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
    <WidgetCard title="Surface Growth" info={INFO} icon={TrendingUp} accent="#06B6D4">
      {isLoading ? (
        <Skeleton className="h-44 w-full" />
      ) : data.length === 0 ? (
        <div className="flex h-44 items-center justify-center text-sm text-muted-foreground">
          No completed scans yet
        </div>
      ) : (
        <AreaTrend data={data} series={SERIES} xKey="t" height={176} />
      )}
    </WidgetCard>
  );
}
