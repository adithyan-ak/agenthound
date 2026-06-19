import { forwardRef, type CSSProperties, type HTMLAttributes } from "react";
import { cn } from "@shared/lib/utils";

export interface StackProps extends HTMLAttributes<HTMLDivElement> {
  /** Spacing token between children. Tailwind-aware lengths only ('1rem', '0.5rem'…). */
  gap?: string;
}

export const Stack = forwardRef<HTMLDivElement, StackProps>(
  ({ gap, className, style, ...rest }, ref) => {
    const cssVars = gap
      ? ({ "--l-stack-gap": gap } as CSSProperties)
      : undefined;
    return (
      <div
        ref={ref}
        className={cn("l-stack", className)}
        style={{ ...cssVars, ...style }}
        {...rest}
      />
    );
  },
);
Stack.displayName = "Stack";
