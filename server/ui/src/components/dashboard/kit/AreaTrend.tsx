import { useId } from "react";
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
} from "recharts";
import { CHART_THEME } from "@/theme/tokens";

export interface TrendSeries {
  key: string;
  label: string;
  color: string;
}

interface AreaTrendProps {
  data: Array<Record<string, number | string>>;
  series: TrendSeries[];
  xKey: string;
  height?: number;
}

/** Gradient-filled multi-series area chart, themed for the dark dashboard. */
export function AreaTrend({ data, series, xKey, height = 180 }: AreaTrendProps) {
  const idBase = useId().replace(/:/g, "");

  return (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data} margin={{ top: 8, right: 8, bottom: 0, left: -8 }}>
        <defs>
          {series.map((s) => (
            <linearGradient key={s.key} id={`${idBase}-${s.key}`} x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor={s.color} stopOpacity={0.45} />
              <stop offset="100%" stopColor={s.color} stopOpacity={0} />
            </linearGradient>
          ))}
        </defs>
        <CartesianGrid vertical={false} stroke={CHART_THEME.grid} />
        <XAxis
          dataKey={xKey}
          tick={{ fill: CHART_THEME.axis, fontSize: 10 }}
          axisLine={false}
          tickLine={false}
          minTickGap={16}
        />
        <YAxis
          tick={{ fill: CHART_THEME.axis, fontSize: 10 }}
          axisLine={false}
          tickLine={false}
          width={44}
          allowDecimals={false}
        />
        <Tooltip
          cursor={{ stroke: "rgba(255,255,255,0.12)", strokeWidth: 1 }}
          contentStyle={{
            backgroundColor: CHART_THEME.tooltip.bg,
            border: `1px solid ${CHART_THEME.tooltip.border}`,
            borderRadius: 8,
            fontSize: 12,
          }}
          labelStyle={{ color: CHART_THEME.tooltip.text, fontWeight: 600 }}
          itemStyle={{ color: CHART_THEME.tooltip.text }}
        />
        {series.map((s) => (
          <Area
            key={s.key}
            type="monotone"
            dataKey={s.key}
            name={s.label}
            stroke={s.color}
            strokeWidth={2}
            fill={`url(#${idBase}-${s.key})`}
            dot={false}
            activeDot={{ r: 4, strokeWidth: 0 }}
            isAnimationActive={false}
          />
        ))}
      </AreaChart>
    </ResponsiveContainer>
  );
}
