import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { AlertOctagon, ArrowRight, ExternalLink } from "lucide-react";
import { useExplorerGraph } from "@features/explorer/model/useExplorerGraph";
import { useExplorerStore } from "@features/explorer/model/store";
import {
  extractCriticalChains,
  type CriticalChain,
} from "@features/explorer/model/critical-chains";
import { SEVERITY } from "@shared/theme/tokens";
import { cn } from "@shared/lib/utils";

export function ChainRibbon() {
  const navigate = useNavigate();
  const activeLens = useExplorerStore((s) => s.activeLens);
  const selectNode = useExplorerStore((s) => s.selectNode);
  const selectedNodeId = useExplorerStore((s) => s.selectedNodeId);
  const openDrawer = useExplorerStore((s) => s.openDrawer);
  const { data } = useExplorerGraph();

  const chains = useMemo(() => {
    if (!data) return [];
    return extractCriticalChains(data.findings);
  }, [data]);

  if (activeLens !== "critical") return null;

  // A clean critical scan is rendered by the canvas' uniform empty-state
  // overlay; the ribbon only appears when there are chains to list.
  if (chains.length === 0) return null;

  return (
    <div
      className={cn(
        "pointer-events-auto absolute bottom-9 left-1/2 z-20 -translate-x-1/2",
        "max-w-[calc(100vw-48px)] overflow-hidden rounded-md border border-border bg-card/95 p-3 backdrop-blur-md elev-2",
      )}
      style={{ boxShadow: `inset 2px 0 0 0 ${SEVERITY.critical.solid}` }}
    >
      <span aria-hidden className="pointer-events-none absolute inset-x-0 top-0 h-px bg-white/[0.05]" />
      <div className="mb-2 flex items-center gap-2 px-1">
        <AlertOctagon className="h-3.5 w-3.5" style={{ color: SEVERITY.critical.solid }} strokeWidth={2.5} />
        <div
          className="font-mono text-[10px] font-semibold uppercase tracking-[0.14em]"
          style={{ color: SEVERITY.critical.text }}
        >
          {chains.length} critical attack path{chains.length === 1 ? "" : "s"}
        </div>
        <div className="font-mono text-[10px] uppercase tracking-[0.06em] text-muted-foreground">
          · click a card to focus
        </div>
      </div>
      <div className="flex gap-2 overflow-x-auto pb-1">
        {chains.slice(0, 8).map((chain) => (
          <ChainCard
            key={chain.id}
            chain={chain}
            selected={selectedNodeId === chain.sourceId}
            onSelect={() => {
              selectNode(chain.sourceId);
              openDrawer();
            }}
            onOpenFinding={() => navigate(`/findings/${chain.findingId}`)}
          />
        ))}
      </div>
    </div>
  );
}

function ChainCard({
  chain,
  selected,
  onSelect,
  onOpenFinding,
}: {
  chain: CriticalChain;
  selected: boolean;
  onSelect: () => void;
  onOpenFinding: () => void;
}) {
  return (
    <div
      role="button"
      tabIndex={0}
      onClick={onSelect}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onSelect();
        }
      }}
      className={cn(
        "relative w-[280px] flex-shrink-0 cursor-pointer overflow-hidden rounded-[3px] border p-3 text-left outline-none",
        "transition-[border-color,background-color] duration-150 ease-out focus-visible:border-primary/60",
        selected
          ? "border-destructive/60 bg-destructive/10"
          : "border-border bg-black/30 hover:border-mauve-7 hover:bg-white/[0.03]",
      )}
      style={{ boxShadow: `inset 2px 0 0 0 ${SEVERITY.critical.solid}` }}
    >
      <div className="mb-1.5 flex items-center justify-between">
        <div className="flex items-center gap-1.5">
          <span
            className="h-1.5 w-1.5 animate-led-pulse rounded-[1px]"
            style={{ background: SEVERITY.critical.solid }}
          />
          <span
            className="font-mono text-[9px] font-semibold uppercase tracking-[0.12em]"
            style={{ color: SEVERITY.critical.text }}
          >
            Critical
          </span>
        </div>
        <div className="flex items-center gap-1.5">
          <span className="font-mono text-[9px] tabular-nums text-muted-foreground">
            {(chain.confidence * 100).toFixed(0)}% conf
          </span>
          <button
            onClick={(e) => {
              e.stopPropagation();
              onOpenFinding();
            }}
            className="flex h-5 w-5 items-center justify-center rounded-[2px] text-muted-foreground transition-colors hover:bg-white/[0.08] hover:text-primary"
            aria-label="Open finding detail"
            title="Open finding detail"
          >
            <ExternalLink className="h-3 w-3" />
          </button>
        </div>
      </div>
      <div className="mb-1.5 line-clamp-2 text-xs font-semibold text-foreground">{chain.title}</div>
      <div className="flex items-center gap-1.5 font-mono text-[10px] text-muted-foreground">
        <span className="max-w-[80px] truncate">{chain.sourceName}</span>
        <ArrowRight className="h-2.5 w-2.5 flex-shrink-0" style={{ color: SEVERITY.critical.solid }} />
        <span className="max-w-[80px] truncate">{chain.targetName}</span>
        <span className="ml-auto text-[9px] uppercase">{chain.edgeKind.replace(/_/g, " ")}</span>
      </div>
    </div>
  );
}
