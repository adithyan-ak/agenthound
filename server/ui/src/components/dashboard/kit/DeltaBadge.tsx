import { ArrowDownRight, ArrowUpRight, Minus } from "lucide-react";
import { cn } from "@/lib/utils";

interface DeltaBadgeProps {
  /** Change vs the comparison point. Null/undefined renders nothing. */
  value: number | null | undefined;
  /** When true, an increase is "bad" (red) — e.g. risk going up. */
  invert?: boolean;
  /** Trailing context, e.g. "vs prev scan". */
  suffix?: string;
  className?: string;
}

/** Small directional delta chip with honest up/down semantics. */
export function DeltaBadge({ value, invert = false, suffix, className }: DeltaBadgeProps) {
  if (value === null || value === undefined) return null;

  const isFlat = value === 0;
  const isBad = invert ? value > 0 : value < 0;
  const Icon = isFlat ? Minus : value > 0 ? ArrowUpRight : ArrowDownRight;

  const tone = isFlat
    ? "text-muted-foreground"
    : isBad
      ? "text-red-400"
      : "text-emerald-400";

  return (
    <span className={cn("inline-flex items-center gap-1 text-xs font-medium", tone, className)}>
      <Icon className="h-3.5 w-3.5" />
      <span className="font-mono tabular-nums">
        {value > 0 ? "+" : ""}
        {value.toLocaleString()}
      </span>
      {suffix && <span className="text-muted-foreground">{suffix}</span>}
    </span>
  );
}
