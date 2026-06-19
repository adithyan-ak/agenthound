import type { ReactNode } from "react";

interface AsyncBoundaryProps {
  /** Query/request in flight. */
  isLoading: boolean;
  /** Request failed. */
  isError?: boolean;
  /** Request succeeded but produced no rows to render. */
  isEmpty?: boolean;
  /** Rendered while loading (e.g. a sized <Skeleton/>). */
  loading: ReactNode;
  /** Rendered on error. Falls back to the empty slot, then nothing. */
  error?: ReactNode;
  /** Rendered when empty. */
  empty?: ReactNode;
  children: ReactNode;
}

/**
 * Declarative loading / error / empty / ready switch.
 *
 * Consolidates the hand-rolled `isLoading ? <Skeleton/> : isError ? ... : ...`
 * ladders scattered across the app. It renders the EXACT node passed for each
 * state (no opinion on markup), so adopting it is pixel-identical to the
 * branch it replaces. Precedence: loading > error > empty > children.
 */
export function AsyncBoundary({
  isLoading,
  isError = false,
  isEmpty = false,
  loading,
  error,
  empty,
  children,
}: AsyncBoundaryProps) {
  if (isLoading) return <>{loading}</>;
  if (isError) return <>{error ?? empty ?? null}</>;
  if (isEmpty) return <>{empty ?? null}</>;
  return <>{children}</>;
}
