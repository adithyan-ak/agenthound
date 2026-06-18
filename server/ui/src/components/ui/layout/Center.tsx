import { forwardRef, type CSSProperties, type HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

/**
 * Centers a block horizontally and caps its measure (line length). Uses
 * content-box so padding doesn't fight the max-inline-size.
 */
export interface CenterProps extends HTMLAttributes<HTMLDivElement> {
  /** Max measure (e.g. '60ch', '72rem'). Default `60ch`. */
  measure?: string;
}

export const Center = forwardRef<HTMLDivElement, CenterProps>(
  ({ measure, className, style, ...rest }, ref) => {
    const cssVars: CSSProperties = {};
    if (measure) (cssVars as Record<string, string>)["--l-center-measure"] = measure;
    return (
      <div
        ref={ref}
        className={cn("l-center", className)}
        style={{ ...cssVars, ...style }}
        {...rest}
      />
    );
  },
);
Center.displayName = "Center";
