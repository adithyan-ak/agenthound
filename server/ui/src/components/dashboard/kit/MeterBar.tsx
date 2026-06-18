import { cn } from "@/lib/utils";

interface MeterBarProps {
  /** Current value. */
  value: number;
  max?: number;
  color: string;
  className?: string;
  /** Bar thickness in px. */
  height?: number;
}

/** Thin gradient progress/meter bar with a glowing fill. */
export function MeterBar({ value, max = 100, color, className, height = 6 }: MeterBarProps) {
  const pct = max > 0 ? Math.max(0, Math.min(100, (value / max) * 100)) : 0;

  return (
    <div
      className={cn("w-full overflow-hidden rounded-full bg-white/[0.06]", className)}
      style={{ height }}
    >
      <div
        className="h-full rounded-full"
        style={{
          width: `${pct}%`,
          background: `linear-gradient(90deg, ${color}99, ${color})`,
          boxShadow: `0 0 12px -2px ${color}99`,
          transition: "width 700ms cubic-bezier(0.22,1,0.36,1)",
        }}
      />
    </div>
  );
}
