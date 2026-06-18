import { Gauge } from "lucide-react";
import { useAllNodes, useDashboardFindings, useDashboardScans } from "@/hooks/useDashboardData";
import { Skeleton } from "@/components/ui/skeleton";
import { WidgetCard, RadialGauge, DeltaBadge, StatusPill } from "./kit";
import type { PillTone } from "./kit";
import { riskColor } from "@/theme/tokens";

const INFO =
  "Composite exposure index = (critical findings x8) + (high findings x3) + (unauthenticated servers x5), capped at 100. Higher is worse.";

function level(score: number): { label: string; tone: PillTone } {
  if (score >= 75) return { label: "Critical", tone: "critical" };
  if (score >= 50) return { label: "Elevated", tone: "high" };
  if (score >= 25) return { label: "Guarded", tone: "medium" };
  return { label: "Low", tone: "success" };
}

interface ContributorProps {
  value: number;
  label: string;
  color: string;
}

function Contributor({ value, label, color }: ContributorProps) {
  return (
    <div className="rounded-lg border border-white/[0.06] bg-white/[0.02] px-3 py-2.5 text-center">
      <p className="font-mono text-xl font-bold tabular-nums" style={{ color: value > 0 ? color : undefined }}>
        {value}
      </p>
      <p className="mt-0.5 text-[10px] font-medium uppercase tracking-wide text-muted-foreground">
        {label}
      </p>
    </div>
  );
}

export function ExposureGauge() {
  const { data: findings, isLoading: loadingFindings } = useDashboardFindings();
  const { data: nodes, isLoading: loadingNodes } = useAllNodes();
  const { data: scans } = useDashboardScans();

  if (loadingFindings || loadingNodes) {
    return (
      <WidgetCard title="Exposure Index" info={INFO} icon={Gauge}>
        <Skeleton className="mx-auto h-44 w-full max-w-[240px]" />
        <div className="mt-4 grid grid-cols-3 gap-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-16 w-full rounded-lg" />
          ))}
        </div>
      </WidgetCard>
    );
  }

  const critical = (findings ?? []).filter((f) => f.severity === "critical").length;
  const high = (findings ?? []).filter((f) => f.severity === "high").length;
  const unauthServers = (nodes ?? []).filter(
    (n) => n.kinds.includes("MCPServer") && String(n.properties.auth_method ?? "none") === "none",
  ).length;

  const score = Math.min(100, critical * 8 + high * 3 + unauthServers * 5);
  const color = riskColor(score);
  const { label, tone } = level(score);

  const completed = (scans ?? []).filter((s) => s.status === "completed");
  const delta =
    completed.length >= 2
      ? (completed[0]?.node_count ?? 0) - (completed[1]?.node_count ?? 0)
      : null;

  return (
    <WidgetCard
      title="Exposure Index"
      info={INFO}
      icon={Gauge}
      accent={color}
      action={<DeltaBadge value={delta} invert suffix="entities" />}
    >
      <div className="flex flex-col items-center">
        <RadialGauge value={score} valueColor={color} caption="of 100" />
        <StatusPill tone={tone} className="-mt-1">
          {label} risk
        </StatusPill>
      </div>
      <div className="mt-4 grid grid-cols-3 gap-2">
        <Contributor value={critical} label="Critical" color="#F87171" />
        <Contributor value={high} label="High" color="#FB923C" />
        <Contributor value={unauthServers} label="Unauth Srv" color="#FACC15" />
      </div>
    </WidgetCard>
  );
}
