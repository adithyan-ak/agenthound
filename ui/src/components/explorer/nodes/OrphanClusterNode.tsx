import { memo, useEffect, useMemo, useRef, useState } from "react";
import { Handle, Position, type NodeProps } from "@xyflow/react";
import { Search, ChevronRight } from "lucide-react";
import type { OrphanClusterData } from "@/lib/explorer/graph-builder";
import {
  getHexConfig,
  HEX_POLYGON_POINTS,
  HEX_NODE_WIDTH,
  HEX_NODE_HEIGHT,
} from "@/lib/explorer/hex-config";
import { useExplorerStore } from "@/store/explorer";
import { cn } from "@/lib/utils";

function OrphanClusterNodeComponent({ data }: NodeProps) {
  const d = data as OrphanClusterData;
  const config = getHexConfig(d.kind);
  const Icon = config.icon;

  const selectNode = useExplorerStore((s) => s.selectNode);
  const openDrawer = useExplorerStore((s) => s.openDrawer);
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const popoverRef = useRef<HTMLDivElement | null>(null);
  const triggerRef = useRef<HTMLButtonElement | null>(null);

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return d.orphanNodes;
    return d.orphanNodes.filter((n) => n.name.toLowerCase().includes(q));
  }, [d.orphanNodes, search]);

  useEffect(() => {
    if (!open) return;
    function onDocClick(e: MouseEvent) {
      const t = e.target as Node;
      if (popoverRef.current?.contains(t)) return;
      if (triggerRef.current?.contains(t)) return;
      setOpen(false);
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") setOpen(false);
    }
    document.addEventListener("mousedown", onDocClick);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onDocClick);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  const strokeColor = config.strokeColor;

  return (
    <div
      className="relative flex flex-col items-center select-none"
      style={{ width: HEX_NODE_WIDTH }}
    >
      <button
        ref={triggerRef}
        onClick={(e) => {
          e.stopPropagation();
          setOpen((v) => !v);
        }}
        className="relative cursor-pointer transition-transform duration-150 hover:scale-[1.08] focus:outline-none"
        style={{ width: HEX_NODE_WIDTH, height: HEX_NODE_HEIGHT }}
        aria-label={`${d.count} unconnected ${d.kindTag}`}
      >
        <svg
          width={HEX_NODE_WIDTH}
          height={HEX_NODE_HEIGHT}
          viewBox={`0 0 ${HEX_NODE_WIDTH} ${HEX_NODE_HEIGHT}`}
          className="absolute inset-0 pointer-events-none"
        >
          <g opacity={0.25}>
            <polygon
              points={HEX_POLYGON_POINTS}
              fill="none"
              stroke={strokeColor}
              strokeWidth={2}
              strokeLinejoin="round"
              transform="translate(-6, -3)"
            />
          </g>
          <g opacity={0.45}>
            <polygon
              points={HEX_POLYGON_POINTS}
              fill="none"
              stroke={strokeColor}
              strokeWidth={2}
              strokeLinejoin="round"
              transform="translate(-3, -1.5)"
            />
          </g>
          <polygon
            points={HEX_POLYGON_POINTS}
            fill="#0B1220"
            stroke={strokeColor}
            strokeWidth={2.5}
            strokeLinejoin="round"
            strokeDasharray="4 2.5"
          />
        </svg>

        <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
          <Icon
            className="h-4 w-4 mb-0.5"
            style={{ color: strokeColor }}
            strokeWidth={2}
          />
          <div
            className="text-[15px] font-bold leading-none tabular-nums"
            style={{ color: strokeColor }}
          >
            {d.count}
          </div>
        </div>
      </button>

      {open && (
        <div
          ref={popoverRef}
          className={cn(
            "absolute left-full top-0 ml-3 w-[300px] rounded-lg border border-slate-700/80 bg-slate-950/98 shadow-2xl backdrop-blur-md",
            "z-[60] overflow-hidden",
            "animate-in fade-in zoom-in-95 duration-150",
          )}
          onClick={(e) => e.stopPropagation()}
        >
          <div
            className="flex items-center gap-2 border-b border-slate-800 px-3 py-2.5"
            style={{ borderTopColor: strokeColor, borderTopWidth: 2 }}
          >
            <Icon
              className="h-4 w-4"
              style={{ color: strokeColor }}
              strokeWidth={2.25}
            />
            <div className="flex flex-col">
              <div className="text-[10px] uppercase tracking-widest text-slate-500">
                Unconnected
              </div>
              <div className="text-sm font-semibold text-white">
                {d.count} {d.kindTag.toLowerCase()}
              </div>
            </div>
          </div>

          {d.orphanNodes.length > 8 && (
            <div className="flex items-center gap-2 border-b border-slate-800 bg-slate-900/50 px-3 py-2">
              <Search className="h-3 w-3 text-slate-500" strokeWidth={2.5} />
              <input
                autoFocus
                placeholder="Filter…"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full bg-transparent text-xs text-white placeholder:text-slate-600 focus:outline-none"
              />
            </div>
          )}

          <div className="max-h-[320px] overflow-y-auto py-1">
            {filtered.length === 0 ? (
              <div className="px-3 py-4 text-center text-xs text-slate-500">
                No matches.
              </div>
            ) : (
              filtered.map((n) => (
                <button
                  key={n.id}
                  onClick={(e) => {
                    e.stopPropagation();
                    selectNode(n.id);
                    openDrawer();
                    setOpen(false);
                  }}
                  className={cn(
                    "group/row flex w-full items-center gap-2 px-3 py-1.5 text-left text-xs",
                    "transition-colors hover:bg-slate-800/60",
                  )}
                >
                  <div
                    className="h-1.5 w-1.5 flex-shrink-0 rounded-full"
                    style={{ background: strokeColor }}
                  />
                  <span className="flex-1 truncate text-slate-200 group-hover/row:text-white">
                    {n.name}
                  </span>
                  <ChevronRight className="h-3 w-3 flex-shrink-0 text-slate-600 group-hover/row:text-slate-400" />
                </button>
              ))
            )}
          </div>

          <div className="border-t border-slate-800 bg-slate-900/40 px-3 py-1.5 text-[9px] uppercase tracking-widest text-slate-600">
            Click a row to inspect
          </div>
        </div>
      )}

      <div className="mt-1 flex flex-col items-center text-center">
        <div
          className="text-[10px] font-semibold whitespace-nowrap max-w-[180px] overflow-hidden text-ellipsis"
          style={{
            color: strokeColor,
            textShadow: "0 1px 4px rgba(0,0,0,0.9)",
          }}
        >
          {d.count} UNCONNECTED
        </div>
        <div className="text-[8px] tracking-[0.12em] text-slate-500 font-medium mt-0.5">
          {d.kindTag}
        </div>
      </div>

      <Handle
        id="h-left"
        type="target"
        position={Position.Left}
        style={{
          position: "absolute",
          left: 2,
          top: 48,
          width: 1,
          height: 1,
          background: "transparent",
          border: "none",
          pointerEvents: "none",
        }}
        isConnectable={false}
      />
      <Handle
        id="h-right"
        type="source"
        position={Position.Right}
        style={{
          position: "absolute",
          left: 82,
          top: 48,
          width: 1,
          height: 1,
          background: "transparent",
          border: "none",
          pointerEvents: "none",
        }}
        isConnectable={false}
      />
    </div>
  );
}

export const OrphanClusterNode = memo(OrphanClusterNodeComponent);
