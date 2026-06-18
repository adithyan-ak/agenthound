import { ShieldAlert } from "lucide-react";
import { useDashboardFindings } from "@/hooks/useDashboardData";
import { Skeleton } from "@/components/ui/skeleton";
import { WidgetCard, MiniRing } from "./kit";
import { SEVERITY, SEVERITY_ORDER, severityColor } from "@/theme/tokens";

const INFO =
  "Open security findings grouped by severity. Each ring shows that severity's share of all findings.";

export function SeverityRings() {
  const { data: findings, isLoading } = useDashboardFindings();

  const counts: Record<string, number> = {};
  for (const sev of SEVERITY_ORDER) counts[sev] = 0;
  for (const f of findings ?? []) counts[f.severity] = (counts[f.severity] ?? 0) + 1;
  const total = (findings ?? []).length;

  return (
    <WidgetCard title="Threat Severity" info={INFO} icon={ShieldAlert}>
      {isLoading ? (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="mx-auto h-[92px] w-[92px] rounded-full" />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          {SEVERITY_ORDER.map((sev) => {
            const count = counts[sev] ?? 0;
            return (
              <div key={sev} className="flex flex-col items-center gap-2">
                <MiniRing
                  fraction={total > 0 ? count / total : 0}
                  value={count}
                  unit="total"
                  color={severityColor(sev)}
                />
                <span className="text-xs font-medium text-muted-foreground">
                  {SEVERITY[sev]?.label ?? sev}
                </span>
              </div>
            );
          })}
        </div>
      )}
    </WidgetCard>
  );
}
