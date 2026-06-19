import { forwardRef, type CSSProperties, type HTMLAttributes } from "react";
import { cn } from "@shared/lib/utils";

export interface ClusterProps extends HTMLAttributes<HTMLDivElement> {
  gap?: string;
  justify?: "flex-start" | "flex-end" | "center" | "space-between";
  align?: "flex-start" | "flex-end" | "center" | "baseline" | "stretch";
}

export const Cluster = forwardRef<HTMLDivElement, ClusterProps>(
  ({ gap, justify, align, className, style, ...rest }, ref) => {
    const cssVars: CSSProperties = {};
    if (gap) (cssVars as Record<string, string>)["--l-cluster-gap"] = gap;
    if (justify) (cssVars as Record<string, string>)["--l-cluster-justify"] = justify;
    if (align) (cssVars as Record<string, string>)["--l-cluster-align"] = align;
    return (
      <div
        ref={ref}
        className={cn("l-cluster", className)}
        style={{ ...cssVars, ...style }}
        {...rest}
      />
    );
  },
);
Cluster.displayName = "Cluster";
