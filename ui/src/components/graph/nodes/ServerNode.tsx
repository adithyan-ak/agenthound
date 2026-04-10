import { Handle, Position, type Node, type NodeProps } from "@xyflow/react";
import { Server } from "lucide-react";
import { cn } from "@/lib/utils";

type ServerNodeData = {
  label: string;
  kind: string;
  color: string;
  riskScore: number;
  sharedWith?: string[];
  properties: Record<string, unknown>;
};

type ServerNodeType = Node<ServerNodeData, "server">;

const AUTH_DOTS: Record<string, { color: string; label: string }> = {
  oauth: { color: "#50C878", label: "oauth" },
  mtls: { color: "#50C878", label: "mTLS" },
  apiKey: { color: "#F5A623", label: "apiKey" },
  bearer: { color: "#F5A623", label: "bearer" },
  none: { color: "#FF6B6B", label: "none" },
};

export function ServerNode({
  data,
  selected,
}: NodeProps<ServerNodeType>) {
  const authMethod = String(data.properties.auth_method ?? "none");
  const dot = AUTH_DOTS[authMethod] ?? AUTH_DOTS["none"]!;

  return (
    <div
      className={cn(
        "rounded-lg border px-2.5 py-1.5 shadow-sm transition-all",
        "bg-[#1a1f2e] border-[#2a2f3e]",
        selected && "ring-2 ring-offset-1 ring-offset-[#0a0f1e]",
      )}
      style={{
        width: 180,
        borderLeftWidth: 4,
        borderLeftColor: "#50C878",
      }}
    >
      <Handle
        type="target"
        position={Position.Left}
        className="!bg-[#4a4f5e] !w-2 !h-2 !border-0"
      />

      <div className="flex items-center gap-1.5">
        <Server size={14} className="text-[#50C878] flex-shrink-0" />
        <span
          className="text-[11px] font-bold text-white truncate flex-1 min-w-0"
          title={data.label}
        >
          {data.label}
        </span>
        <span
          title={`Auth: ${dot.label}`}
          className="flex-shrink-0"
          style={{
            width: 6,
            height: 6,
            borderRadius: "50%",
            backgroundColor: dot.color,
            display: "inline-block",
          }}
        />
      </div>

      <Handle
        type="source"
        position={Position.Right}
        className="!bg-[#4a4f5e] !w-2 !h-2 !border-0"
      />
    </div>
  );
}
