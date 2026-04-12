import {
  getBezierPath,
  BaseEdge,
  EdgeLabelRenderer,
  type Edge,
  type EdgeProps,
} from "@xyflow/react";

type AttackEdgeData = {
  kind: string;
  animated?: boolean;
  isCrossProtocol?: boolean;
  sourceKind?: string;
  targetKind?: string;
};

type AttackEdgeType = Edge<AttackEdgeData, "attack">;

export function AttackEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
  selected,
  markerEnd,
}: EdgeProps<AttackEdgeType>) {
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
    sourcePosition,
    targetPosition,
  });

  const showAnimation = selected || data?.animated;
  const kind = data?.kind ?? "";
  const isCrossProtocol = data?.isCrossProtocol ?? false;

  const stroke = isCrossProtocol ? "#FF2D2D" : "#FF2D2D";
  const strokeWidth = isCrossProtocol ? 4 : 2.5;
  const glowFilter = isCrossProtocol ? `url(#ah-cross-glow)` : undefined;

  return (
    <>
      {isCrossProtocol && (
        <defs>
          <filter id="ah-cross-glow" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur stdDeviation="3" result="coloredBlur" />
            <feMerge>
              <feMergeNode in="coloredBlur" />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>
      )}
      <BaseEdge
        id={id}
        path={edgePath}
        markerEnd={markerEnd}
        style={{
          stroke,
          strokeWidth,
          strokeDasharray: "8 4",
          filter: glowFilter,
        }}
      />

      {showAnimation && (
        <circle r={isCrossProtocol ? "4" : "3"} fill="#FFB800">
          <animateMotion dur="2s" repeatCount="indefinite" path={edgePath} />
        </circle>
      )}

      {kind && (
        <EdgeLabelRenderer>
          <div
            className="nodrag nopan pointer-events-auto"
            style={{
              position: "absolute",
              transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            }}
          >
            {isCrossProtocol ? (
              <div className="flex items-center gap-1">
                <span className="text-[9px] font-bold px-1.5 py-0.5 rounded bg-red-600 text-white whitespace-nowrap shadow-lg shadow-red-900/50 border border-red-300/30">
                  {kind}
                </span>
                <span className="text-[8px] font-bold px-1 py-0.5 rounded bg-amber-500 text-black whitespace-nowrap uppercase tracking-wide">
                  cross-protocol
                </span>
              </div>
            ) : (
              <span className="text-[9px] font-semibold px-1.5 py-0.5 rounded bg-red-600/90 text-white whitespace-nowrap">
                {kind}
              </span>
            )}
          </div>
        </EdgeLabelRenderer>
      )}
    </>
  );
}
