import { useEffect, useMemo, useRef } from "react";
import {
  Crosshair,
  ArrowDown,
  ArrowUp,
  Radar,
  FileLock2,
  Shield,
  Target,
  Crown,
  Copy,
  FileText,
  Focus,
  type LucideIcon,
} from "lucide-react";
import { useExplorerStore } from "@features/explorer/model/store";
import { useExplorerGraph } from "@features/explorer/model/useExplorerGraph";
import { formatNodeAsMarkdown } from "@features/explorer/lib/export-markdown";
import { useMarksStore } from "@shared/model/marks";
import { getHexConfig } from "@shared/lib/hex-config";
import { cn } from "@shared/lib/utils";
import { useEscapeKey } from "@shared/lib/useEscapeKey";
import {
  buildAdjacencyIndex,
  bfsFrom,
  type TraversalDirection,
} from "@shared/lib/graph/traverse";

const MENU_WIDTH = 260;
const MENU_HEIGHT_ESTIMATE = 520;

interface MenuItem {
  type: "item";
  icon: LucideIcon;
  label: string;
  shortcut?: string;
  onSelect: () => void;
  destructive?: boolean;
  active?: boolean;
  disabled?: boolean;
}

interface MenuDivider {
  type: "divider";
}

interface MenuSection {
  type: "section";
  label: string;
}

type MenuEntry = MenuItem | MenuDivider | MenuSection;

export function ExplorerNodeContextMenu() {
  const contextMenu = useExplorerStore((s) => s.contextMenu);
  const closeContextMenu = useExplorerStore((s) => s.closeContextMenu);
  const activeLens = useExplorerStore((s) => s.activeLens);
  const blastSourceId = useExplorerStore((s) => s.blastRadiusSourceId);
  const setActiveLens = useExplorerStore((s) => s.setActiveLens);
  const setBlastRadiusSource = useExplorerStore((s) => s.setBlastRadiusSource);
  const selectNode = useExplorerStore((s) => s.selectNode);
  const openDrawer = useExplorerStore((s) => s.openDrawer);
  const setHighlight = useExplorerStore((s) => s.setHighlight);

  const toggleOwned = useMarksStore((s) => s.toggleOwned);
  const toggleHighValue = useMarksStore((s) => s.toggleHighValue);
  const isOwned = useMarksStore((s) => s.isOwned);
  const isHighValue = useMarksStore((s) => s.isHighValue);

  const { data } = useExplorerGraph();

  const menuRef = useRef<HTMLDivElement | null>(null);

  useEscapeKey(closeContextMenu, {
    enabled: !!contextMenu,
    ignoreInputs: false,
  });

  useEffect(() => {
    if (!contextMenu) return;
    function onClick(e: MouseEvent) {
      if (menuRef.current?.contains(e.target as Node)) return;
      closeContextMenu();
    }
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [contextMenu, closeContextMenu]);

  const node = useMemo(() => {
    if (!contextMenu || !data) return null;
    return data.nodes.find((n) => n.id === contextMenu.nodeId) ?? null;
  }, [contextMenu, data]);

  // Adjacency index for the focus/reach actions — built once per data load,
  // then reused by all three traversals (replaces the per-call edge scan).
  const adjacency = useMemo(
    () => (data ? buildAdjacencyIndex(data.edges) : null),
    [data],
  );

  if (!contextMenu || !node) return null;

  const kind = node.kinds[0] ?? "Unknown";
  const config = getHexConfig(kind);
  const Icon = config.icon;
  const name = String(
    node.properties?.name ??
      node.properties?.uri ??
      node.properties?.path ??
      node.id.slice(0, 12),
  );

  // Viewport-aware positioning: if the click is near the right edge, open
  // leftwards. If near the bottom, open upwards.
  const vw = typeof window !== "undefined" ? window.innerWidth : 1920;
  const vh = typeof window !== "undefined" ? window.innerHeight : 1080;
  const left =
    contextMenu.x + MENU_WIDTH > vw - 16
      ? Math.max(8, contextMenu.x - MENU_WIDTH)
      : contextMenu.x;
  const top =
    contextMenu.y + MENU_HEIGHT_ESTIMATE > vh - 16
      ? Math.max(8, contextMenu.y - MENU_HEIGHT_ESTIMATE)
      : contextMenu.y;

  // --- action handlers ---

  function applyNeighborhood(
    maxHops: number,
    direction: TraversalDirection,
    title: string,
  ) {
    if (!adjacency) return;
    const { nodeIds, edgeKeys } = bfsFrom(node!.id, adjacency, {
      maxHops,
      direction,
      edgeKey: (e) => `${e.source}|${e.target}|${e.kind}`,
    });
    setHighlight({
      nodeIds: Array.from(nodeIds),
      edgeIds: Array.from(edgeKeys),
      title,
    });
    selectNode(node!.id);
    closeContextMenu();
  }

  function handleSetBlastSource() {
    setActiveLens("blast-radius");
    setBlastRadiusSource(node!.id);
    selectNode(node!.id);
    closeContextMenu();
  }

  function handleFocusNeighborhood() {
    applyNeighborhood(2, "both", `2-hop neighborhood · ${name}`);
  }

  function handleShowInbound() {
    applyNeighborhood(6, "in", `Inbound reach · ${name}`);
  }

  function handleShowOutbound() {
    applyNeighborhood(6, "out", `Outbound reach · ${name}`);
  }

  function handleViewEvidence() {
    selectNode(node!.id);
    openDrawer("evidence");
    closeContextMenu();
  }

  function handleViewRemediation() {
    selectNode(node!.id);
    openDrawer("remediation");
    closeContextMenu();
  }

  function handleToggleOwned() {
    toggleOwned(node!.id);
    closeContextMenu();
  }

  function handleToggleHighValue() {
    toggleHighValue(node!.id);
    closeContextMenu();
  }

  function handleCopyObjectId() {
    void navigator.clipboard.writeText(node!.id);
    closeContextMenu();
  }

  function handleCopyAsMarkdown() {
    const md = formatNodeAsMarkdown(node!, config.kindTag, owned, highValue);
    void navigator.clipboard.writeText(md);
    closeContextMenu();
  }

  // Contextual label for the blast radius action.
  let blastRadiusLabel = "Set as blast radius source";
  if (activeLens === "blast-radius") {
    if (blastSourceId === node.id) {
      blastRadiusLabel = "Refocus here";
    } else {
      blastRadiusLabel = "Use as blast radius source";
    }
  }

  const owned = isOwned(node.id);
  const highValue = isHighValue(node.id);

  const entries: MenuEntry[] = [
    { type: "section", label: "Focus" },
    {
      type: "item",
      icon: Radar,
      label: blastRadiusLabel,
      onSelect: handleSetBlastSource,
    },
    {
      type: "item",
      icon: Crosshair,
      label: "Focus 2-hop neighborhood",
      onSelect: handleFocusNeighborhood,
    },
    {
      type: "item",
      icon: ArrowDown,
      label: "Show what reaches this",
      onSelect: handleShowInbound,
    },
    {
      type: "item",
      icon: ArrowUp,
      label: "Show what this reaches",
      onSelect: handleShowOutbound,
    },
    { type: "divider" },
    { type: "section", label: "Inspect" },
    {
      type: "item",
      icon: FileLock2,
      label: "View evidence",
      onSelect: handleViewEvidence,
    },
    {
      type: "item",
      icon: Shield,
      label: "View remediation",
      onSelect: handleViewRemediation,
    },
    { type: "divider" },
    { type: "section", label: "Mark" },
    {
      type: "item",
      icon: Target,
      label: owned ? "Unmark Owned" : "Mark as Owned",
      onSelect: handleToggleOwned,
      active: owned,
      destructive: !owned,
    },
    {
      type: "item",
      icon: Crown,
      label: highValue ? "Unmark High Value" : "Mark as High Value",
      onSelect: handleToggleHighValue,
      active: highValue,
    },
    { type: "divider" },
    { type: "section", label: "Copy" },
    {
      type: "item",
      icon: Copy,
      label: "Copy objectid",
      onSelect: handleCopyObjectId,
    },
    {
      type: "item",
      icon: FileText,
      label: "Copy as Markdown",
      onSelect: handleCopyAsMarkdown,
    },
  ];

  return (
    <div
      ref={menuRef}
      role="menu"
      className={cn(
        "fixed z-[70] overflow-hidden rounded-md border border-border bg-card/95 backdrop-blur-md elev-3",
        "animate-in fade-in zoom-in-95 duration-100",
      )}
      style={{ left, top, width: MENU_WIDTH }}
      onContextMenu={(e) => e.preventDefault()}
    >
      {/* Header */}
      <div className="relative flex items-center gap-2.5 border-b border-border px-3 py-2.5">
        <span
          aria-hidden
          className="pointer-events-none absolute inset-x-0 top-0 h-px"
          style={{ background: config.strokeColor, opacity: 0.7 }}
        />
        <div
          className="flex h-7 w-7 items-center justify-center rounded-[3px] border"
          style={{
            borderColor: config.strokeColor,
            background: `${config.strokeColor}15`,
          }}
        >
          <Icon
            className="h-3.5 w-3.5"
            style={{ color: config.strokeColor }}
            strokeWidth={2.25}
          />
        </div>
        <div className="flex min-w-0 flex-col">
          <div
            className="font-mono text-[10px] uppercase tracking-[0.14em]"
            style={{ color: config.strokeColor }}
          >
            {config.kindTag}
          </div>
          <div className="truncate text-xs font-semibold text-foreground">
            {name}
          </div>
        </div>
        {owned && (
          <Target
            className="ml-auto h-3 w-3 flex-shrink-0 text-red-400"
            strokeWidth={2.75}
          />
        )}
        {highValue && (
          <Crown
            className="ml-auto h-3 w-3 flex-shrink-0 text-yellow-400"
            strokeWidth={2.5}
            fill="currentColor"
            fillOpacity={0.25}
          />
        )}
      </div>

      {/* Items */}
      <div className="py-1">
        {entries.map((entry, i) => {
          if (entry.type === "divider") {
            return (
              <div
                key={`div-${i}`}
                className="my-1 border-t border-border"
              />
            );
          }
          if (entry.type === "section") {
            return (
              <div
                key={`sec-${i}`}
                className="px-3 pb-1 pt-2 font-mono text-[9px] font-semibold uppercase tracking-[0.14em] text-muted-foreground/70"
              >
                {entry.label}
              </div>
            );
          }
          const ItemIcon = entry.icon;
          return (
            <button
              key={`item-${i}-${entry.label}`}
              role="menuitem"
              onClick={(e) => {
                e.stopPropagation();
                entry.onSelect();
              }}
              disabled={entry.disabled}
              className={cn(
                "group/mi flex w-full items-center gap-2.5 px-3 py-1.5 text-left text-xs",
                "transition-colors",
                !entry.disabled && !entry.active && "text-foreground hover:bg-white/[0.05]",
                entry.active && "bg-white/[0.05] text-foreground",
                entry.disabled && "cursor-not-allowed text-muted-foreground/70",
              )}
            >
              <ItemIcon
                className={cn(
                  "h-3.5 w-3.5 flex-shrink-0",
                  entry.active && "text-red-400",
                )}
                strokeWidth={2.25}
              />
              <span className="flex-1">{entry.label}</span>
            </button>
          );
        })}
      </div>

      {/* Footer hint */}
      <div className="flex items-center gap-2 border-t border-border bg-black/20 px-3 py-1.5 font-mono text-[9px] uppercase tracking-[0.14em] text-muted-foreground/70">
        <Focus className="h-2.5 w-2.5" strokeWidth={2.5} />
        <span>Esc to close</span>
      </div>
    </div>
  );
}
