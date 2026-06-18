import type { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";

interface SubMetric {
  /** Numeric value; when > 0 the chip turns into a risk highlight. */
  value: number;
  /** Trailing label, e.g. "unauth" or "exec/shell". */
  label: string;
}

interface StatTileProps {
  icon: LucideIcon;
  label: string;
  value: number;
  /** Node-kind color (hex). Drives the icon chip + accents. */
  color: string;
  sub?: SubMetric;
  onClick?: () => void;
  className?: string;
}

/**
 * KPI tile: a kind-colored icon chip, a big mono count, a label, and an
 * optional *real* risk sub-metric chip (red when non-zero, muted "clear"
 * otherwise) so each tile carries signal, not just a number.
 */
export function StatTile({ icon: Icon, label, value, color, sub, onClick, className }: StatTileProps) {
  const danger = sub ? sub.value > 0 : false;
  const Comp = onClick ? "button" : "div";

  return (
    <Comp
      onClick={onClick}
      className={cn(
        "card-elevated group h-full rounded-xl p-4 text-left",
        onClick && "cursor-pointer transition-transform hover:-translate-y-0.5",
        className,
      )}
    >
      <div className="flex items-start justify-between">
        <span
          className="flex h-9 w-9 items-center justify-center rounded-lg ring-1 ring-inset ring-white/10"
          style={{ backgroundColor: `${color}1f` }}
        >
          <Icon className="h-[18px] w-[18px]" style={{ color }} />
        </span>
        {sub && (
          <span
            className={cn(
              "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide",
              danger
                ? "bg-red-500/12 text-red-400 ring-1 ring-inset ring-red-500/25"
                : "bg-emerald-500/10 text-emerald-400 ring-1 ring-inset ring-emerald-500/20",
            )}
          >
            {danger ? (
              <>
                <span className="h-1.5 w-1.5 rounded-full bg-red-500" />
                {sub.value} {sub.label}
              </>
            ) : (
              "clear"
            )}
          </span>
        )}
      </div>
      <p className="mt-3 font-mono text-[28px] font-semibold leading-none tabular-nums text-foreground">
        {value.toLocaleString()}
      </p>
      <p className="mt-1.5 text-xs text-muted-foreground">{label}</p>
    </Comp>
  );
}
