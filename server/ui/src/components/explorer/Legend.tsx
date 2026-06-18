import { useExplorerStore } from "@/store/explorer";
import { getLens } from "@/lib/explorer/lens-config";
import { cn } from "@/lib/utils";
import { NODE_KIND_COLORS, SEVERITY, DIMMED } from "@/theme/tokens";
import { EDGE_CATEGORY_COLORS } from "@/lib/edge-styles";

interface LegendItem {
  color: string;
  label: string;
  dashed?: boolean;
  thick?: boolean;
}

// Legend palettes derive from the canonical token tables (NODE_KIND_COLORS,
// SEVERITY, EDGE_CATEGORY_COLORS, DIMMED). Earlier this file duplicated 30+
// hex literals which silently drifted from tokens.ts; never re-introduce
// raw hex here — add a token first.

const LEGEND_BY_LENS: Record<string, LegendItem[]> = {
  topology: [
    { color: NODE_KIND_COLORS.AgentInstance, label: "Agent" },
    { color: NODE_KIND_COLORS.MCPServer, label: "MCP Server" },
    { color: NODE_KIND_COLORS.MCPTool, label: "Tool" },
    { color: NODE_KIND_COLORS.MCPResource, label: "Resource" },
    { color: EDGE_CATEGORY_COLORS.trust, label: "Trust edge" },
  ],
  "attack-surface": [
    { color: SEVERITY.critical.solid, label: "Critical reach", thick: true },
    { color: SEVERITY.high.solid, label: "High reach", thick: true },
    { color: SEVERITY.medium.solid, label: "Medium" },
    { color: NODE_KIND_COLORS.Identity, label: "Low" },
  ],
  critical: [
    { color: SEVERITY.critical.solid, label: "Critical path", thick: true },
    { color: NODE_KIND_COLORS.A2AAgent, label: "Cross-protocol", dashed: true },
    { color: DIMMED.deep, label: "Dimmed context" },
  ],
  "cross-protocol": [
    { color: NODE_KIND_COLORS.A2AAgent, label: "Cross-protocol edge", dashed: true, thick: true },
    { color: NODE_KIND_COLORS.AgentInstance, label: "A2A Agent" },
    { color: NODE_KIND_COLORS.MCPServer, label: "MCP Server" },
  ],
  credentials: [
    { color: NODE_KIND_COLORS.Credential, label: "Credential" },
    { color: NODE_KIND_COLORS.Identity, label: "Identity" },
    { color: NODE_KIND_COLORS.MCPServer, label: "MCP Server" },
  ],
  poisoning: [
    { color: NODE_KIND_COLORS.InstructionFile, label: "Poisoned tool" },
    { color: SEVERITY.high.solid, label: "Shadowing" },
    { color: SEVERITY.critical.solid, label: "Poisoned instructions" },
  ],
  "blast-radius": [
    { color: NODE_KIND_COLORS.MCPServer, label: "Source node", thick: true },
    { color: NODE_KIND_COLORS.Identity, label: "Reachable" },
    { color: DIMMED.deep, label: "Out of scope" },
  ],
  chokepoints: [
    { color: NODE_KIND_COLORS.AgentInstance, label: "High centrality" },
    { color: NODE_KIND_COLORS.ResourceGroup, label: "Medium centrality" },
    { color: DIMMED.mid, label: "Low centrality" },
  ],
};

export function Legend() {
  const activeLens = useExplorerStore((s) => s.activeLens);
  const lens = getLens(activeLens);
  const items = LEGEND_BY_LENS[activeLens] ?? [];

  if (items.length === 0) return null;

  return (
    <div
      className={cn(
        "pointer-events-auto absolute left-6 bottom-12 z-20 rounded-lg",
        "glass px-3 py-2.5 elev-2",
      )}
    >
      <div className="mb-2 text-[9px] font-semibold uppercase tracking-widest text-muted-foreground">
        {lens.shortLabel} legend
      </div>
      <div className="space-y-1.5">
        {items.map((item, i) => (
          <div key={i} className="flex items-center gap-2">
            <svg width="28" height="10" className="flex-shrink-0">
              {item.dashed ? (
                <line
                  x1="0"
                  y1="5"
                  x2="28"
                  y2="5"
                  stroke={item.color}
                  strokeWidth={item.thick ? 3 : 2}
                  strokeDasharray="4 3"
                />
              ) : (
                <line
                  x1="0"
                  y1="5"
                  x2="28"
                  y2="5"
                  stroke={item.color}
                  strokeWidth={item.thick ? 3 : 2}
                />
              )}
            </svg>
            <div className="text-[10px] text-foreground">{item.label}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
