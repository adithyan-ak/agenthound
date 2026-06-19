// Explorer graph transform barrel. Public surface = the same names the former
// monolithic graph-builder.ts exported, plus the protocol helpers.
export { buildExplorerGraph } from "./build-graph";
export { buildFindingIndex, severityRank } from "./build-edges";
export {
  isCrossProtocolEdge,
  protocolDomain,
  MCP_NODE_KINDS,
  A2A_NODE_KINDS,
} from "./protocol";
export type {
  HexNodeData,
  OrphanClusterData,
  LensEdgeData,
  BundledEdge,
  BuildOptions,
  BuildResult,
  LensMetrics,
} from "./types";
