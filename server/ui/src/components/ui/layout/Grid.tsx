import { forwardRef, type CSSProperties, type HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

/**
 * Auto-fit responsive grid. Columns = `auto-fit, minmax(min(min, 100%), 1fr)`.
 * No breakpoints — children automatically reflow when their min-width fits
 * one more or one fewer time.
 */
export interface GridProps extends HTMLAttributes<HTMLDivElement> {
  /** Minimum child width before columns wrap. Default `16rem`. */
  min?: string;
  gap?: string;
}

export const Grid = forwardRef<HTMLDivElement, GridProps>(
  ({ min, gap, className, style, ...rest }, ref) => {
    const cssVars: CSSProperties = {};
    if (min) (cssVars as Record<string, string>)["--l-grid-min"] = min;
    if (gap) (cssVars as Record<string, string>)["--l-grid-gap"] = gap;
    return (
      <div
        ref={ref}
        className={cn("l-grid", className)}
        style={{ ...cssVars, ...style }}
        {...rest}
      />
    );
  },
);
Grid.displayName = "Grid";
