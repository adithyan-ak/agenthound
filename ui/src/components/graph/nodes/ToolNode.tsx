import { Handle, Position, type Node, type NodeProps } from "@xyflow/react";
import {
  Zap,
  FileText,
  Globe,
  Database,
  Code,
  Mail,
} from "lucide-react";
import type { ElementType } from "react";
import { cn } from "@/lib/utils";

type ToolNodeData = {
  label: string;
  kind: string;
  color: string;
  riskScore: number;
  isOverflow?: boolean;
  overflowCount?: number;
  properties: Record<string, unknown>;
};

type ToolNodeType = Node<ToolNodeData, "tool">;

const CAPABILITY_ICONS: Record<string, ElementType> = {
  shell_access: Zap,
  file_read: FileText,
  file_write: FileText,
  network_outbound: Globe,
  database_access: Database,
  code_execution: Code,
  email_send: Mail,
};

const CAPABILITY_DANGER_ORDER = [
  "shell_access",
  "code_execution",
  "network_outbound",
  "database_access",
  "email_send",
  "file_write",
  "file_read",
];

export function ToolNode({
  data,
  selected,
}: NodeProps<ToolNodeType>) {
  if (data.isOverflow) {
    return (
      <div
        className={cn(
          "rounded-full border border-dashed px-2 py-0.5 shadow-sm transition-all cursor-pointer",
          "bg-[#1a1f2e]/60 border-[#2a2f3e]",
          "flex items-center justify-center",
          selected && "ring-1 ring-offset-1 ring-offset-[#0a0f1e]",
        )}
        style={{ width: 110, height: 26 }}
      >
        <Handle
          type="target"
          position={Position.Left}
          className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
        />
        <span className="text-[10px] text-gray-500 italic truncate">
          +{data.overflowCount ?? 0} more
        </span>
        <Handle
          type="source"
          position={Position.Right}
          className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
        />
      </div>
    );
  }

  const capabilities = Array.isArray(data.properties.capability_surface)
    ? (data.properties.capability_surface as string[])
    : [];

  const topCap = CAPABILITY_DANGER_ORDER.find((c) => capabilities.includes(c));
  const TopIcon = topCap ? CAPABILITY_ICONS[topCap] : null;

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
          backgroundColor: "#F5A623",
          display: "inline-block",
        }}
      />
      <span
        className="text-[10px] text-white truncate flex-1 min-w-0"
        title={data.label}
      >
        {data.label}
      </span>
      {TopIcon && (
        <span title={topCap} className="flex-shrink-0">
          <TopIcon size={10} className="text-[#F5A623]/70" />
        </span>
      )}

      <Handle
        type="source"
        position={Position.Right}
        className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
      />
    </div>
  );
}
