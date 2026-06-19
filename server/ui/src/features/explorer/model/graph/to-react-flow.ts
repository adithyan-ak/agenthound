import type { Node, Edge } from "@xyflow/react";
import type {
  LensEdgeData,
  LogicalClusterNode,
  LogicalEdge,
  LogicalHexNode,
} from "./types";

/**
 * React-Flow adapter. The pure transform stays free of RF concerns; this
 * module layers on the node `type`/`position` and edge `type`/handle ids that
 * the canvas needs. Edge `type` is derived (self-loop > cross-protocol > lens)
 * exactly as the original inline switch did.
 */
export function logicalEdgeToReactFlow(e: LogicalEdge): Edge<LensEdgeData> {
  const isSelfLoop = e.source === e.target;
  const type = isSelfLoop
    ? "self-loop"
    : e.data.isCrossProtocol
      ? "lens-cross"
      : "lens";
  return {
    id: e.id,
    source: e.source,
    target: e.target,
    type,
    sourceHandle: "h-bottom",
    targetHandle: "h-top",
    data: e.data,
  };
}

export function logicalHexNodeToReactFlow(n: LogicalHexNode): Node {
  return {
    id: n.id,
    type: "hex",
    position: { x: 0, y: 0 },
    data: n.data,
  };
}

export function logicalClusterNodeToReactFlow(n: LogicalClusterNode): Node {
  return {
    id: n.id,
    type: "orphan-cluster",
    position: { x: 0, y: 0 },
    data: n.data,
  };
}
