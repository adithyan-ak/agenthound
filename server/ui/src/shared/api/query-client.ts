import { QueryClient } from "@tanstack/react-query";

// Central TanStack Query client. The 30s staleTime is the global default for
// every query, so entity hooks must NOT repeat `staleTime: 30_000`; they only
// set genuine per-query overrides (e.g. the health poll's refetchInterval).
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});
