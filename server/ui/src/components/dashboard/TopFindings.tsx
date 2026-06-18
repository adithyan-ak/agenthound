import { useNavigate, Link } from "react-router-dom";
import { Siren, ArrowRight } from "lucide-react";
import { useDashboardFindings } from "@/hooks/useDashboardData";
import { Skeleton } from "@/components/ui/skeleton";
import { WidgetCard, StatusPill } from "./kit";
import type { PillTone } from "./kit";
import { severityColor } from "@/theme/tokens";

const INFO =
  "Most critical security findings sorted by severity and confidence. Click any alert to see the full attack path, evidence, and remediation.";

const SEVERITY_RANK: Record<string, number> = { critical: 4, high: 3, medium: 2, low: 1 };

export function TopFindings() {
  const navigate = useNavigate();
  const { data: findings, isLoading } = useDashboardFindings();

  const top = (findings ?? [])
    .filter((f) => f.severity === "critical" || f.severity === "high")
    .sort(
      (a, b) =>
        (SEVERITY_RANK[b.severity] ?? 0) - (SEVERITY_RANK[a.severity] ?? 0) ||
        (b.confidence ?? 0) - (a.confidence ?? 0),
    )
    .slice(0, 7);

  return (
    <WidgetCard
      title="Critical Alerts"
      info={INFO}
      icon={Siren}
      accent="#EF4444"
      action={
        <Link
          to="/findings"
          className="flex items-center gap-1 text-xs text-primary transition-colors hover:text-primary/80"
        >
          View all <ArrowRight className="h-3 w-3" />
        </Link>
      }
    >
      {isLoading ? (
        <Skeleton className="h-56 w-full" />
      ) : top.length === 0 ? (
        <div className="flex h-56 flex-col items-center justify-center gap-1 text-center">
          <p className="text-sm font-medium text-foreground">No critical alerts</p>
          <p className="text-xs text-muted-foreground">No critical or high-severity findings detected.</p>
        </div>
      ) : (
        <ul className="space-y-2">
          {top.map((f) => (
            <li key={f.id}>
              <button
                onClick={() => navigate(`/findings/${f.id}`)}
                className="group/row flex w-full items-center gap-3 overflow-hidden rounded-lg border border-white/[0.05] bg-white/[0.02] px-3 py-2.5 text-left transition-colors hover:border-primary/30 hover:bg-white/[0.04]"
              >
                <span
                  aria-hidden
                  className="h-9 w-1 shrink-0 rounded-full"
                  style={{ backgroundColor: severityColor(f.severity) }}
                />
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium text-foreground">{f.title}</p>
                  <p className="truncate text-xs text-muted-foreground">
                    {f.source_name} <span className="text-muted-foreground/60">&rarr;</span> {f.target_name}
                  </p>
                </div>
                {f.owasp_map && f.owasp_map.length > 0 && (
                  <span className="hidden shrink-0 rounded bg-white/[0.05] px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground sm:inline">
                    {f.owasp_map[0]}
                  </span>
                )}
                <StatusPill tone={f.severity as PillTone} dot={false}>
                  {f.severity}
                </StatusPill>
              </button>
            </li>
          ))}
        </ul>
      )}
    </WidgetCard>
  );
}
