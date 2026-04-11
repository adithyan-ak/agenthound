import { useEffect, useRef } from "react";
import {
  Crosshair,
  Eye,
  EyeOff,
  Target,
  Crown,
  Copy,
  GitBranch,
  ArrowRight,
  ArrowLeft,
} from "lucide-react";
import { useGraphStore } from "@/store/graph";

export interface ContextMenuState {
  nodeId: string;
  nodeLabel: string;
  nodeKind: string;
  x: number;
  y: number;
}

interface NodeContextMenuProps {
  state: ContextMenuState | null;
  onClose: () => void;
  onFocus: (id: string) => void;
  onShowBlastRadius: (id: string) => void;
  onShowInbound: (id: string) => void;
}

export function NodeContextMenu({
  state,
  onClose,
  onFocus,
  onShowBlastRadius,
  onShowInbound,
}: NodeContextMenuProps) {
  const menuRef = useRef<HTMLDivElement>(null);
  const toggleOwned = useGraphStore((s) => s.toggleOwned);
  const toggleHighValue = useGraphStore((s) => s.toggleHighValue);
  const ownedNodeIds = useGraphStore((s) => s.ownedNodeIds);
  const highValueNodeIds = useGraphStore((s) => s.highValueNodeIds);

  useEffect(() => {
    if (!state) return;
    function handleClickOutside(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        onClose();
      }
    }
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [state, onClose]);

  if (!state) return null;

  const isOwned = ownedNodeIds.includes(state.nodeId);
  const isHighValue = highValueNodeIds.includes(state.nodeId);

  async function copyValue(value: string) {
    try {
      await navigator.clipboard.writeText(value);
    } catch {
      // ignore
    }
  }

  const items: Array<{
    label: string;
    icon: React.ReactNode;
    onClick: () => void;
    divider?: boolean;
    danger?: boolean;
  }> = [
    {
      label: "Focus here",
      icon: <Crosshair className="h-3.5 w-3.5" />,
      onClick: () => {
        onFocus(state.nodeId);
        onClose();
      },
    },
    {
      label: "Show what this can reach",
      icon: <ArrowRight className="h-3.5 w-3.5" />,
      onClick: () => {
        onShowBlastRadius(state.nodeId);
        onClose();
      },
    },
    {
      label: "Show what reaches this",
      icon: <ArrowLeft className="h-3.5 w-3.5" />,
      onClick: () => {
        onShowInbound(state.nodeId);
        onClose();
      },
      divider: true,
    },
    {
      label: isOwned ? "Unmark Owned" : "Mark as Owned",
      icon: isOwned ? (
        <EyeOff className="h-3.5 w-3.5 text-red-400" />
      ) : (
        <Target className="h-3.5 w-3.5 text-red-400" />
      ),
      onClick: () => {
        toggleOwned(state.nodeId);
        onClose();
      },
    },
    {
      label: isHighValue ? "Unmark High Value" : "Mark as High Value",
      icon: isHighValue ? (
        <Eye className="h-3.5 w-3.5 text-yellow-400" />
      ) : (
        <Crown className="h-3.5 w-3.5 text-yellow-400" />
      ),
      onClick: () => {
        toggleHighValue(state.nodeId);
        onClose();
      },
      divider: true,
    },
    {
      label: "Copy name",
      icon: <Copy className="h-3.5 w-3.5" />,
      onClick: () => {
        void copyValue(state.nodeLabel);
        onClose();
      },
    },
    {
      label: "Copy ID",
      icon: <Copy className="h-3.5 w-3.5" />,
      onClick: () => {
        void copyValue(state.nodeId);
        onClose();
      },
    },
    {
      label: "Copy Cypher MATCH",
      icon: <GitBranch className="h-3.5 w-3.5" />,
      onClick: () => {
        void copyValue(
          `MATCH (n:${state.nodeKind} {objectid: '${state.nodeId}'}) RETURN n`,
        );
        onClose();
      },
    },
  ];

  return (
    <div
      ref={menuRef}
      className="fixed z-50 min-w-[220px] rounded-md bg-card border border-border shadow-xl py-1 text-sm"
      style={{ left: state.x, top: state.y }}
    >
      <div className="px-3 py-1.5 text-[10px] uppercase tracking-wide text-muted-foreground border-b border-border mb-1 truncate">
        {state.nodeLabel}
      </div>
      {items.map((item, idx) => (
        <div key={idx}>
          <button
            onClick={item.onClick}
            className="w-full flex items-center gap-2 px-3 py-1.5 text-left hover:bg-accent text-foreground"
          >
            {item.icon}
            <span className="flex-1">{item.label}</span>
          </button>
          {item.divider && <div className="h-px bg-border my-1" />}
        </div>
      ))}
    </div>
  );
}
