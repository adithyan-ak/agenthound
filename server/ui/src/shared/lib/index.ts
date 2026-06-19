// shared/lib barrel — generic, domain-free utilities.
export { cn } from "./utils";
export { useCountUp } from "./useCountUp";
export { useEscapeKey, type UseEscapeKeyOptions } from "./useEscapeKey";
export {
  buildAdjacencyIndex,
  bfsFrom,
  type AdjacencyIndex,
  type BfsOptions,
  type BfsResult,
  type TraversalDirection,
  type MinimalEdge,
} from "./graph/traverse";
