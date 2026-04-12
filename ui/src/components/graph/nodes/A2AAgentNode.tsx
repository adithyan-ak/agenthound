import { Handle, Position, type Node, type NodeProps } from "@xyflow/react";
import { Bot, AlertTriangle, Target, Crown } from "lucide-react";
import { cn } from "@/lib/utils";
import { useGraphStore } from "@/store/graph";

type A2AAgentNodeData = {
  label: string;
  kind: string;
  color: string;
  riskScore: number;
  properties: Record<string, unknown>;
};

type A2AAgentNodeType = Node<A2AAgentNodeData, "a2aAgent">;

const AUTH_DOTS: Record<string, { color: string; label: string }> = {
  oauth: { color: "#50C878", label: "oauth" },
  oidc: { color: "#50C878", label: "oidc" },
  mtls: { color: "#50C878", label: "mTLS" },
  apiKey: { color: "#F5A623", label: "apiKey" },
  bearer: { color: "#F5A623", label: "bearer" },
  none: { color: "#FF6B6B", label: "none" },
};

export function A2AAgentNode({
  id,
  data,
  selected,
}: NodeProps<A2AAgentNodeType>) {
  const authMethod = String(data.properties.auth_method ?? "none");
  const dot = AUTH_DOTS[authMethod] ?? AUTH_DOTS["none"]!;
  const skillCount = Number(data.properties.skill_count ?? 0);
  const isSigned = data.properties.is_signed;
  const unsigned = isSigned === false;
  const isOwned = useGraphStore((s) => s.ownedNodeIds.includes(id));
  const isHighValue = useGraphStore((s) => s.highValueNodeIds.includes(id));

  return (
    <div
      className={cn(
        "relative rounded-lg border px-2.5 py-1.5 shadow-sm transition-all",
        "bg-[#1a1f2e] border-[#2a2f3e]",
        selected && "ring-2 ring-offset-1 ring-offset-[#0a0f1e]",
        isOwned && "ring-2 ring-red-500 shadow-red-900/50 shadow-lg",
        isHighValue && "ring-2 ring-yellow-400 shadow-yellow-900/50 shadow-lg",
      )}
      style={{
        width: 180,
        borderLeftWidth: 4,
        borderLeftColor: "#7B68EE",
      }}
    >
      {isOwned && (
        <div className="absolute -top-1.5 -right-1.5 h-4 w-4 rounded-full bg-red-600 flex items-center justify-center shadow-md">
          <Target className="h-2.5 w-2.5 text-white" />
        </div>
      )}
      {isHighValue && !isOwned && (
        <div className="absolute -top-1.5 -right-1.5 h-4 w-4 rounded-full bg-yellow-500 flex items-center justify-center shadow-md">
          <Crown className="h-2.5 w-2.5 text-black" />
        </div>
      )}
      <Handle
        type="target"
        position={Position.Left}
        className="!bg-[#4a4f5e] !w-2 !h-2 !border-0"
      />

      <div className="flex items-center gap-1.5">
        <Bot size={14} className="text-[#7B68EE] flex-shrink-0" />
        <span
          className="text-[11px] font-bold text-white truncate flex-1 min-w-0"
          title={data.label}
        >
          {data.label}
          {skillCount > 0 && (
            <span className="font-normal text-gray-400">
              {" "}({skillCount})
            </span>
          )}
        </span>
        {unsigned && (
          <span title="Unsigned agent card" className="flex-shrink-0">
            <AlertTriangle size={12} className="text-amber-400" />
          </span>
        )}
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
