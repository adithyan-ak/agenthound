import {
  forwardRef,
  type CSSProperties,
  type HTMLAttributes,
  type ReactNode,
} from "react";
import { cn } from "@shared/lib/utils";

/**
 * Sidebar layout: a fixed-ish sidebar paired with intrinsically growing
 * main content. Wraps to a column when main's `contentMin` no longer fits.
 *
 * Pass `side` and `main` as props rather than children — the layout uses
 * specific class names on the two slots and accepting children would let
 * callers accidentally inject extra elements that break the wrapping logic.
 */
export interface SidebarProps extends Omit<HTMLAttributes<HTMLDivElement>, "children"> {
  side: ReactNode;
  main: ReactNode;
  /** Ideal sidebar width. Default `18rem`. */
  sideWidth?: string;
  /** Min main width before wrap (% of total). Default `60%`. */
  contentMin?: string;
  gap?: string;
  /** Place sidebar on the right instead of left. */
  sidePosition?: "left" | "right";
}

export const Sidebar = forwardRef<HTMLDivElement, SidebarProps>(
  ({ side, main, sideWidth, contentMin, gap, sidePosition = "left", className, style, ...rest }, ref) => {
    const cssVars: CSSProperties = {};
    if (sideWidth) (cssVars as Record<string, string>)["--l-sidebar-side-width"] = sideWidth;
    if (contentMin) (cssVars as Record<string, string>)["--l-sidebar-content-min"] = contentMin;
    if (gap) (cssVars as Record<string, string>)["--l-sidebar-gap"] = gap;
    return (
      <div
        ref={ref}
        className={cn("l-sidebar", className)}
        style={{ ...cssVars, ...style }}
        {...rest}
      >
        {sidePosition === "left" ? (
          <>
            <div className="l-sidebar__side">{side}</div>
            <div className="l-sidebar__main">{main}</div>
          </>
        ) : (
          <>
            <div className="l-sidebar__main">{main}</div>
            <div className="l-sidebar__side">{side}</div>
          </>
        )}
      </div>
    );
  },
);
Sidebar.displayName = "Sidebar";
