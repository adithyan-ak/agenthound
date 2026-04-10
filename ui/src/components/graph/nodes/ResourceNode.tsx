import { Handle, Position, type Node, type NodeProps } from "@xyflow/react";
import { cn } from "@/lib/utils";

type ResourceNodeData = {
  label: string;
  kind: string;
  color: string;
  riskScore: number;
  properties: Record<string, unknown>;
};

type ResourceNodeType = Node<ResourceNodeData, "resource">;

const SENSITIVITY_CONFIG: Record<
  string,
  { accent: string; bg: string; text: string }
> = {
  critical: { accent: "#D0021B", bg: "#D0021B20", text: "#FF4D4D" },
  high: { accent: "#FF8C00", bg: "#FF8C0020", text: "#FFA940" },
  medium: { accent: "#F5A623", bg: "#F5A62320", text: "#F5C451" },
  low: { accent: "#8E8E93", bg: "#8E8E9320", text: "#A0A0A5" },
};

export function ResourceNode({
  data,
  selected,
}: NodeProps<ResourceNodeType>) {
  const sensitivity = String(data.properties.sensitivity ?? "low");
  const config = SENSITIVITY_CONFIG[sensitivity] ?? SENSITIVITY_CONFIG["low"]!;

  return (
    <div
      className={cn(
        "rounded-full border px-2 py-0.5 shadow-sm transition-all",
        "bg-[#1a1f2e] border-[#2a2f3e]",
        "flex items-center gap-1",
        sensitivity === "critical" && "shadow-red-900/30 shadow-md",
        selected && "ring-1 ring-offset-1 ring-offset-[#0a0f1e]",
      )}
      style={{
        width: 110,
        height: 26,
        borderColor: sensitivity === "critical" ? "#D0021B40" : undefined,
      }}
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
          backgroundColor: config.accent,
          display: "inline-block",
        }}
      />
      <span
        className="text-[10px] text-white truncate flex-1 min-w-0"
        title={data.label}
      >
        {data.label}
      </span>
      <span
        title={sensitivity.toUpperCase()}
        className="flex-shrink-0"
        style={{
          width: 5,
          height: 5,
          borderRadius: "50%",
          backgroundColor: config.accent,
          display: "inline-block",
        }}
      />

      <Handle
        type="source"
        position={Position.Right}
        className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
      />
    </div>
  );
}
