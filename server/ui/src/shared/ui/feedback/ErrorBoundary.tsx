import { Component, type ErrorInfo, type ReactNode } from "react";

interface ErrorBoundaryProps {
  children: ReactNode;
  /** Custom fallback. Receives the caught error when it is a render function. */
  fallback?: ReactNode | ((error: Error) => ReactNode);
  /** Side-effect hook for logging/telemetry. */
  onError?: (error: Error, info: ErrorInfo) => void;
}

interface ErrorBoundaryState {
  error: Error | null;
}

/**
 * Catches render-time errors in its subtree and shows a fallback instead of
 * unmounting the whole React tree. Use a root instance in app/providers plus
 * per-feature instances so one widget crashing does not blank the page.
 */
export class ErrorBoundary extends Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  state: ErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    this.props.onError?.(error, info);
  }

  render(): ReactNode {
    const { error } = this.state;
    if (error) {
      const { fallback } = this.props;
      if (typeof fallback === "function") return fallback(error);
      if (fallback !== undefined) return fallback;
      return (
        <div
          role="alert"
          className="flex h-full min-h-[120px] flex-col items-center justify-center gap-1 p-4 text-center"
        >
          <p className="font-mono text-sm font-medium text-foreground">
            Something went wrong
          </p>
          <p className="font-mono text-[11px] uppercase tracking-wide text-muted-foreground">
            {error.message || "Unexpected error"}
          </p>
        </div>
      );
    }
    return this.props.children;
  }
}
