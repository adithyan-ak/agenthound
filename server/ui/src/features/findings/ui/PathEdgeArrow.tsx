import { memo } from "react";
import { ChevronRight } from "lucide-react";
import { getEdgeCategory } from "@entities/edge";
import { EDGE_COLORS } from "@shared/theme/tokens";

interface PathEdgeArrowProps {
  kind: string;
}

function PathEdgeArrowComponent({ kind }: PathEdgeArrowProps) {
  const category = getEdgeCategory(kind);
  const color = EDGE_COLORS[category as keyof typeof EDGE_COLORS] ?? EDGE_COLORS.structure;

  return (
    <div className="flex min-w-[64px] flex-shrink-0 flex-col items-center justify-center px-1">
      <div
        className="mb-1 whitespace-nowrap rounded-[2px] border bg-black/50 px-1.5 py-0.5 font-mono text-[8px] font-semibold uppercase tracking-[0.08em]"
        style={{ color, borderColor: `${color}55` }}
      >
        {kind.replace(/_/g, " ")}
      </div>
      <div className="flex w-full items-center">
        <div className="flex-1 border-t border-dashed" style={{ borderColor: `${color}99` }} />
        <ChevronRight className="-ml-0.5 h-3 w-3 flex-shrink-0" style={{ color }} />
      </div>
    </div>
  );
}

export const PathEdgeArrow = memo(PathEdgeArrowComponent);
