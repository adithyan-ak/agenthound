import { StrictMode, type ReactNode } from "react";
import { BrowserRouter } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@shared/api/query-client";
import { ErrorBoundary } from "@shared/ui/feedback";

/**
 * Root infrastructure wiring: the TanStack Query cache, the router, and a
 * top-level error boundary so a render-time crash shows a fallback instead of
 * blanking the page. Composition-only — no domain logic lives here.
 */
export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <StrictMode>
      <ErrorBoundary>
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>{children}</BrowserRouter>
        </QueryClientProvider>
      </ErrorBoundary>
    </StrictMode>
  );
}
