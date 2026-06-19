import { forwardRef, type CSSProperties, type HTMLAttributes } from "react";
import { cn } from "@shared/lib/utils";

/**
 * Switcher: row → column when container falls below `threshold`. No media
 * queries — pure flex-basis arithmetic. Use for content panels that should
 * sit side-by-side on wider containers and stack on narrow ones.
 */
export interface SwitcherProps extends HTMLAttributes<HTMLDivElement> {
  /** Width below which children stack. Default `30rem`. */
  threshold?: string;
  gap?: string;
}

export const Switcher = forwardRef<HTMLDivElement, SwitcherProps>(
  ({ threshold, gap, className, style, ...rest }, ref) => {
    const cssVars: CSSProperties = {};
    if (threshold) (cssVars as Record<string, string>)["--l-switcher-threshold"] = threshold;
    if (gap) (cssVars as Record<string, string>)["--l-switcher-gap"] = gap;
    return (
      <div
        ref={ref}
        className={cn("l-switcher", className)}
        style={{ ...cssVars, ...style }}
        {...rest}
      />
    );
  },
);
Switcher.displayName = "Switcher";
