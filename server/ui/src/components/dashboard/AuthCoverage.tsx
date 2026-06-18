import { useMemo } from "react";
import { Lock } from "lucide-react";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";
import { useAllNodes } from "@/hooks/useDashboardData";
import { Skeleton } from "@/components/ui/skeleton";
import { WidgetCard } from "./kit";
import { CHART_THEME } from "@/theme/tokens";

const INFO =
  "Authentication method distribution across MCP servers and A2A agents. 'none' (red) means no authentication — your highest-priority targets.";

const AUTH_COLORS: Record<string, string> = {
  none: "#EF4444",
  apiKey: CHART_THEME.series[3] ?? "#F59E0B",
  bearer: CHART_THEME.series[5] ?? "#3B82F6",
  oauth: CHART_THEME.series[2] ?? "#10B981",
  mTLS: CHART_THEME.series[0] ?? "#06B6D4",
};
const FALLBACK = CHART_THEME.axis;

const AUTH_LABELS: Record<string, string> = {
  none: "None",
  apiKey: "API Key",
  bearer: "Bearer",
  oauth: "OAuth",
  mTLS: "mTLS",
};

export function AuthCoverage() {
  const { data: nodes, isLoading } = useAllNodes();

  const { chartData, total } = useMemo(() => {
    const grouped: Record<string, number> = {};
    for (const node of nodes ?? []) {
      if (!node.kinds.includes("MCPServer") && !node.kinds.includes("A2AAgent")) continue;
      const method = String(node.properties.auth_method ?? "none");
      grouped[method] = (grouped[method] ?? 0) + 1;
    }
    const data = Object.entries(grouped)
      .map(([name, value]) => ({ name, value }))
      .sort((a, b) => b.value - a.value);
    return { chartData: data, total: data.reduce((s, d) => s + d.value, 0) };
  }, [nodes]);

  return (
    <WidgetCard title="Auth Coverage" info={INFO} icon={Lock}>
      {isLoading ? (
        <Skeleton className="h-48 w-full" />
      ) : total === 0 ? (
        <div className="flex h-48 items-center justify-center text-sm text-muted-foreground">
          No endpoints to assess
        </div>
      ) : (
        <div className="flex items-center gap-4">
          <div className="relative h-40 w-40 shrink-0">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={chartData}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  innerRadius={52}
                  outerRadius={72}
                  paddingAngle={2}
                  strokeWidth={0}
                  startAngle={90}
                  endAngle={-270}
                  isAnimationActive={false}
                >
                  {chartData.map((entry) => (
                    <Cell key={entry.name} fill={AUTH_COLORS[entry.name] ?? FALLBACK} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{
                    backgroundColor: CHART_THEME.tooltip.bg,
                    border: `1px solid ${CHART_THEME.tooltip.border}`,
                    borderRadius: 8,
                    fontSize: 12,
                  }}
                  itemStyle={{ color: CHART_THEME.tooltip.text }}
                  formatter={(value: number, name: string) => [value, AUTH_LABELS[name] ?? name]}
                />
              </PieChart>
            </ResponsiveContainer>
            <div className="pointer-events-none absolute inset-0 flex flex-col items-center justify-center">
              <span className="font-mono text-2xl font-bold tabular-nums text-foreground">{total}</span>
              <span className="text-[10px] uppercase tracking-wide text-muted-foreground">endpoints</span>
            </div>
          </div>

          <ul className="min-w-0 flex-1 space-y-1.5">
            {chartData.map((d) => {
              const pct = total > 0 ? Math.round((d.value / total) * 100) : 0;
              return (
                <li key={d.name} className="flex items-center gap-2 text-sm">
                  <span
                    className="h-2.5 w-2.5 shrink-0 rounded-sm"
                    style={{ backgroundColor: AUTH_COLORS[d.name] ?? FALLBACK }}
                  />
                  <span className="truncate text-muted-foreground">{AUTH_LABELS[d.name] ?? d.name}</span>
                  <span className="ml-auto font-mono tabular-nums text-foreground">{d.value}</span>
                  <span className="w-9 text-right font-mono text-xs tabular-nums text-muted-foreground">
                    {pct}%
                  </span>
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </WidgetCard>
  );
}
