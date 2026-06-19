import { Handle, Position } from "@xyflow/react";
import type { CSSProperties } from "react";

// Invisible 1px handles. Edges connect via the bottom (source) / top (target)
// pair; the left/right handles exist for completeness. Shared verbatim by
// HexNode and OrphanClusterNode so the connection geometry stays identical.
const BASE: CSSProperties = {
  position: "absolute",
  width: 1,
  height: 1,
  background: "transparent",
  border: "none",
  pointerEvents: "none",
};

/**
 * The four cardinal React Flow handles (top/bottom/left/right) shared by every
 * explorer hex. Previously duplicated byte-for-byte across HexNode and
 * OrphanClusterNode.
 */
export function NodeHandles() {
  return (
    <>
      <Handle
        id="h-top"
        type="target"
        position={Position.Top}
        style={{ ...BASE, left: 42, top: -6 }}
        isConnectable={false}
      />
      <Handle
        id="h-bottom"
        type="source"
        position={Position.Bottom}
        style={{ ...BASE, left: 42, top: 134 }}
        isConnectable={false}
      />
      <Handle
        id="h-left"
        type="target"
        position={Position.Left}
        style={{ ...BASE, left: 2, top: 48 }}
        isConnectable={false}
      />
      <Handle
        id="h-right"
        type="source"
        position={Position.Right}
        style={{ ...BASE, left: 82, top: 48 }}
        isConnectable={false}
      />
    </>
  );
}
