import { ArrowUpRight, ArrowDownLeft } from "lucide-react";
import type { APIEdge } from "@entities/graph/dto";
import { useGraphStore } from "../model/graph-store";
import { useUIStore } from "@shared/model/ui-store";
import { Badge } from "@shared/ui/primitives/badge";
import { Button } from "@shared/ui/primitives/button";
import { Separator } from "@shared/ui/primitives/separator";

interface NodeConnectionsProps {
  edges: APIEdge[];
  nodeId: string;
}

export function NodeConnections({ edges, nodeId }: NodeConnectionsProps) {
  const selectNode = useGraphStore((s) => s.selectNode);
  const openSidebar = useUIStore((s) => s.openSidebar);

  if (edges.length === 0) {
    return (
      <div className="py-4 text-sm text-muted-foreground text-center">
        No connections
      </div>
    );
  }

  const grouped = new Map<string, APIEdge[]>();
  for (const edge of edges) {
    const list = grouped.get(edge.kind) ?? [];
    list.push(edge);
    grouped.set(edge.kind, list);
  }

  function handleClick(edge: APIEdge) {
    const otherId = edge.source === nodeId ? edge.target : edge.source;
    selectNode(otherId);
    openSidebar();
  }

  const groupEntries = Array.from(grouped.entries());

  return (
    <div className="space-y-3">
      {groupEntries.map(([kind, kindEdges], groupIdx) => (
        <div key={kind}>
          <div className="flex items-center justify-between mb-1">
            <Badge variant="outline" className="text-[10px]">
              {kind}
            </Badge>
            <span className="text-[10px] text-muted-foreground">{kindEdges.length}</span>
          </div>
          <div className="space-y-0.5">
            {kindEdges.map((edge, i) => {
              const isOutgoing = edge.source === nodeId;
              const otherId = isOutgoing ? edge.target : edge.source;
              const otherName =
                (isOutgoing
                  ? edge.properties?.target_name
                  : edge.properties?.source_name) ?? otherId.slice(0, 12);

              return (
                <Button
                  key={`${edge.kind}-${i}`}
                  variant="ghost"
                  onClick={() => handleClick(edge)}
                  className="flex w-full items-center gap-2 h-auto px-2 py-1 justify-start text-xs"
                >
                  {isOutgoing ? (
                    <ArrowUpRight className="h-3 w-3 text-muted-foreground flex-shrink-0" />
                  ) : (
                    <ArrowDownLeft className="h-3 w-3 text-muted-foreground flex-shrink-0" />
                  )}
                  <span className="text-foreground truncate">
                    {String(otherName)}
                  </span>
                  <span className="ml-auto text-[10px] text-muted-foreground">
                    {isOutgoing ? "out" : "in"}
                  </span>
                </Button>
              );
            })}
          </div>
          {groupIdx < groupEntries.length - 1 && <Separator className="mt-3" />}
        </div>
      ))}
    </div>
  );
}
