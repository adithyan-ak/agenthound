import * as React from "react";
import type { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { InfoTip } from "@/components/dashboard/InfoTip";

interface WidgetCardProps {
  title?: React.ReactNode;
  /** Tooltip text rendered next to the title. */
  info?: string;
  /** Optional leading icon for the header. */
  icon?: LucideIcon;
  /** Color for the leading icon + top accent hairline. */
  accent?: string;
  /** Right-aligned header slot (e.g. a "View all" link or a badge). */
  action?: React.ReactNode;
  className?: string;
  contentClassName?: string;
  /** Drop the default content padding (for edge-to-edge tables). */
  flush?: boolean;
  children: React.ReactNode;
}

/**
 * Standard dashboard panel: an elevated, glowing surface with a consistent
 * header (icon + title + info tip + action slot). Every widget on the home
 * page is wrapped in one so the grid reads as a single, cohesive system.
 */
export function WidgetCard({
  title,
  info,
  icon: Icon,
  accent,
  action,
  className,
  contentClassName,
  flush = false,
  children,
}: WidgetCardProps) {
  return (
    <section
      className={cn(
        "card-elevated group relative flex h-full flex-col overflow-hidden rounded-xl",
        className,
      )}
    >
      {accent && (
        <span
          aria-hidden
          className="pointer-events-none absolute inset-x-0 top-0 h-px opacity-60"
          style={{ background: `linear-gradient(90deg, transparent, ${accent}, transparent)` }}
        />
      )}
      {(title || action) && (
        <header className="flex items-center justify-between gap-3 px-5 pb-3 pt-4">
          <div className="flex min-w-0 items-center gap-2">
            {Icon && (
              <Icon
                className="h-4 w-4 shrink-0"
                style={accent ? { color: accent } : undefined}
              />
            )}
            {title && (
              <h3 className="truncate text-[13px] font-semibold tracking-tight text-foreground">
                {title}
              </h3>
            )}
            {info && <InfoTip text={info} />}
          </div>
          {action && <div className="shrink-0">{action}</div>}
        </header>
      )}
      <div className={cn(!flush && "px-5 pb-5", flush && "pb-0", contentClassName)}>
        {children}
      </div>
    </section>
  );
}
