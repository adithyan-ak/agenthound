import { Handle, Position, type Node, type NodeProps } from "@xyflow/react";
import {
  Key,
  Lock,
  Server,
  FileText,
  FileWarning,
} from "lucide-react";
import type { ElementType } from "react";
import { cn } from "@/lib/utils";

type InfraNodeData = {
  label: string;
  kind: string;
  color: string;
  riskScore: number;
  properties: Record<string, unknown>;
};

type InfraNodeType = Node<InfraNodeData, "infra">;

const KIND_ICONS: Record<string, ElementType> = {
  Identity: Key,
  Credential: Lock,
  Host: Server,
  ConfigFile: FileText,
  InstructionFile: FileWarning,
};

const KIND_COLORS: Record<string, string> = {
  Identity: "#8E8E93",
  Credential: "#FF6B6B",
  Host: "#2C3E50",
  ConfigFile: "#95A5A6",
  InstructionFile: "#BDC3C7",
};

export function InfraNode({
  data,
  selected,
}: NodeProps<InfraNodeType>) {
  const kind = data.kind ?? "";
  const Icon = KIND_ICONS[kind] ?? FileText;
  const accentColor = KIND_COLORS[kind] ?? "#8E8E93";

  const isExposedCredential =
    kind === "Credential" && data.properties.is_exposed === true;

  const borderColor = isExposedCredential ? "#FF6B6B" : accentColor;

  return (
    <div
      className={cn(
        "rounded-full shadow-sm transition-all flex items-center justify-center",
        "bg-[#1a1f2e]",
        isExposedCredential && "shadow-red-500/40 shadow-md",
        selected && "ring-1 ring-offset-1 ring-offset-[#0a0f1e]",
      )}
      style={{
        width: 28,
        height: 28,
        border: `2px solid ${borderColor}`,
      }}
      title={data.label}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
      />

      <Icon size={12} style={{ color: borderColor }} />

      <Handle
        type="source"
        position={Position.Right}
        className="!bg-[#4a4f5e] !w-1.5 !h-1.5 !border-0"
      />
    </div>
  );
}
