import { useMemo } from "react";
import { Layers } from "lucide-react";
import { useDashboardFindings } from "@/hooks/useDashboardData";
import { Skeleton } from "@/components/ui/skeleton";
import { WidgetCard } from "./kit";
import { SEVERITY, SEVERITY_ORDER, severityColor } from "@/theme/tokens";

const INFO =
  "Security findings grouped by category. Bar length reflects total findings; the colored segments show the severity mix within each category.";

const MAX_ROWS = 8;

interface CategoryRow {
  name: string;
  total: number;
  counts: Record<string, number>;
}

export function CategoryBreakdown() {
  const { data: findings, isLoading } = useDashboardFindings();

  const rows = useMemo<CategoryRow[]>(() => {
    const map = new Map<string, CategoryRow>();
    for (const f of findings ?? []) {
      const row = map.get(f.category) ?? { name: f.category, total: 0, counts: {} };
      row.total += 1;
      row.counts[f.severity] = (row.counts[f.severity] ?? 0) + 1;
      map.set(f.category, row);
    }
    return Array.from(map.values())
      .sort((a, b) => b.total - a.total)
      .slice(0, MAX_ROWS);
  }, [findings]);

  const maxTotal = rows.reduce((m, r) => Math.max(m, r.total), 0);

  return (
    <WidgetCard title="Findings by Category" info={INFO} icon={Layers}>
      {isLoading ? (
        <Skeleton className="h-56 w-full" />
      ) : rows.length === 0 ? (
        <div className="flex h-56 items-center justify-center text-sm text-muted-foreground">
          No findings yet
        </div>
      ) : (
        <div className="space-y-3.5">
          <div className="flex flex-wrap items-center gap-x-4 gap-y-1">
            {SEVERITY_ORDER.map((sev) => (
              <span key={sev} className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
                <span
                  className="h-2 w-2 rounded-full"
                  style={{ backgroundColor: severityColor(sev) }}
                />
                {SEVERITY[sev]?.label ?? sev}
              </span>
            ))}
          </div>

          <div className="space-y-2.5">
            {rows.map((row) => (
              <div key={row.name} className="space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="truncate pr-2 text-foreground">{row.name}</span>
                  <span className="font-mono tabular-nums text-muted-foreground">{row.total}</span>
                </div>
                <div className="h-2.5 w-full overflow-hidden rounded-full bg-white/[0.04]">
                  <div
                    className="flex h-full overflow-hidden rounded-full"
                    style={{ width: `${maxTotal > 0 ? (row.total / maxTotal) * 100 : 0}%` }}
                  >
                    {SEVERITY_ORDER.map((sev) => {
                      const count = row.counts[sev] ?? 0;
                      if (count === 0) return null;
                      return (
                        <span
                          key={sev}
                          title={`${SEVERITY[sev]?.label ?? sev}: ${count}`}
                          style={{
                            width: `${(count / row.total) * 100}%`,
                            backgroundColor: severityColor(sev),
                          }}
                        />
                      );
                    })}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </WidgetCard>
  );
}
