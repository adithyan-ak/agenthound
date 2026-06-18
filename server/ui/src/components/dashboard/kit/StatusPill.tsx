import { cn } from "@/lib/utils";

export type PillTone =
  | "critical"
  | "high"
  | "medium"
  | "low"
  | "info"
  | "success"
  | "warning"
  | "error"
  | "neutral";

const TONE_CLASS: Record<PillTone, string> = {
  critical: "bg-red-500/12 text-red-400 ring-red-500/25",
  high: "bg-orange-500/12 text-orange-400 ring-orange-500/25",
  medium: "bg-yellow-500/12 text-yellow-400 ring-yellow-500/25",
  low: "bg-slate-500/12 text-slate-300 ring-slate-500/25",
  info: "bg-blue-500/12 text-blue-400 ring-blue-500/25",
  success: "bg-emerald-500/12 text-emerald-400 ring-emerald-500/25",
  warning: "bg-amber-500/12 text-amber-400 ring-amber-500/25",
  error: "bg-red-500/12 text-red-400 ring-red-500/25",
  neutral: "bg-white/[0.06] text-muted-foreground ring-white/10",
};

const DOT_CLASS: Record<PillTone, string> = {
  critical: "bg-red-500",
  high: "bg-orange-500",
  medium: "bg-yellow-500",
  low: "bg-slate-400",
  info: "bg-blue-500",
  success: "bg-emerald-500",
  warning: "bg-amber-500",
  error: "bg-red-500",
  neutral: "bg-muted-foreground",
};

interface StatusPillProps {
  tone: PillTone;
  children: React.ReactNode;
  /** Show the leading status dot. */
  dot?: boolean;
  /** Soft pulsing dot (e.g. for an in-progress scan). */
  pulse?: boolean;
  className?: string;
}

/** Compact, ringed status/severity pill with an optional indicator dot. */
export function StatusPill({ tone, children, dot = true, pulse = false, className }: StatusPillProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide ring-1 ring-inset",
        TONE_CLASS[tone],
        className,
      )}
    >
      {dot && (
        <span className="relative flex h-1.5 w-1.5">
          {pulse && (
            <span
              className={cn(
                "absolute inline-flex h-full w-full animate-ping rounded-full opacity-75",
                DOT_CLASS[tone],
              )}
            />
          )}
          <span className={cn("relative inline-flex h-1.5 w-1.5 rounded-full", DOT_CLASS[tone])} />
        </span>
      )}
      {children}
    </span>
  );
}
