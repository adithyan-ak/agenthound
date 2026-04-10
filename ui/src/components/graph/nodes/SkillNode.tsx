import { Handle, Position, type Node, type NodeProps } from "@xyflow/react";
import { cn } from "@/lib/utils";

type SkillNodeData = {
  label: string;
  kind: string;
  color: string;
  riskScore: number;
  properties: Record<string, unknown>;
};

type SkillNodeType = Node<SkillNodeData, "skill">;

export function SkillNode({
  data,
  selected,
}: NodeProps<SkillNodeType>) {
  return (
    <div
      className={cn(
        "rounded-full border px-2 py-0.5 shadow-sm transition-all",
        "bg-[#1a1f2e] border-[#2a2f3e]",
        "flex items-center gap-1",
        selected && "ring-1 ring-offset-1 ring-offset-[#0a0f1e]",
      )}
      style={{ width: 110, height: 26 }}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
      />

      <span
        className="flex-shrink-0"
        style={{
          width: 4,
          height: 4,
          borderRadius: "50%",
          backgroundColor: "#9B59B6",
          display: "inline-block",
        }}
      />
      <span
        className="text-[10px] text-white truncate flex-1 min-w-0"
        title={data.label}
      >
        {data.label}
      </span>

      <Handle
        type="source"
        position={Position.Right}
        className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
      />
    </div>
  );
}
