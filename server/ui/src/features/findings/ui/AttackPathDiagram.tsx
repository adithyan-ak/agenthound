import { Waypoints } from "lucide-react";
import { WidgetCard } from "@shared/ui/widgets";
import type { AttackPath, AttackPathNode } from "@entities/finding/model";
import { PathHexNode } from "./PathHexNode";
import { PathEdgeArrow } from "./PathEdgeArrow";

interface AttackPathDiagramProps {
  path: AttackPath | null;
  severity: string;
  sourceId: string;
  sourceName: string;
  sourceKind: string;
  targetId: string;
  targetName: string;
  targetKind: string;
}

export function AttackPathDiagram({
  path,
  severity,
  sourceId,
  sourceName,
  sourceKind,
  targetId,
  targetName,
  targetKind,
}: AttackPathDiagramProps) {
  const hasPath = !!path && path.nodes.length > 0;

  const fallbackNodes: AttackPathNode[] = [
    { id: sourceId, kinds: [sourceKind], properties: { name: sourceName } },
    { id: targetId, kinds: [targetKind], properties: { name: targetName } },
  ];
  const orderedNodes =
    path && path.nodes.length > 0 ? orderNodesFromEdges(path) : fallbackNodes;

  return (
    <WidgetCard
      title="Attack Path"
      icon={Waypoints}
      action={
        <span className="font-mono text-[10px] uppercase tracking-[0.12em] text-muted-foreground">
          {String(Math.max(orderedNodes.length - 1, 0)).padStart(2, "0")} hops
        </span>
      }
    >
      <div className="hud-grid overflow-x-auto rounded-[3px] border border-border/60 bg-black/20 p-4">
        <div className="flex min-w-max items-center justify-center gap-0">
          {orderedNodes.map((node, i) => {
            const isFirst = i === 0;
            const isLast = i === orderedNodes.length - 1;
            const edgeKind = hasPath ? path?.edges[i]?.kind : isLast ? undefined : "\u2014";
            return (
              <div key={node.id} className="flex items-center">
                <PathHexNode node={node} isFirst={isFirst} isLast={isLast} severity={severity} />
                {!isLast && edgeKind && <PathEdgeArrow kind={edgeKind} />}
              </div>
            );
          })}
        </div>
        {!hasPath && (
          <p className="mt-3 text-center font-mono text-[10px] uppercase tracking-[0.12em] text-muted-foreground">
            Intermediate path details unavailable
          </p>
        )}
      </div>
    </WidgetCard>
  );
}

function orderNodesFromEdges(path: AttackPath): AttackPathNode[] {
  if (path.edges.length === 0) return path.nodes;

  const targetSet = new Set(path.edges.map((e) => e.target));

  let startId = path.edges[0]!.source;
  for (const edge of path.edges) {
    if (!targetSet.has(edge.source)) {
      startId = edge.source;
      break;
    }
  }

  const edgeMap = new Map<string, (typeof path.edges)[0]>();
  for (const e of path.edges) {
    edgeMap.set(e.source, e);
  }

  const order: string[] = [];
  order.push(startId);
  let current = startId;
  const visited = new Set<string>();
  visited.add(current);

  while (edgeMap.has(current)) {
    const edge = edgeMap.get(current)!;
    if (visited.has(edge.target)) break;
    order.push(edge.target);
    visited.add(edge.target);
    current = edge.target;
  }

  const nodeMap = new Map(path.nodes.map((n) => [n.id, n]));
  return order
    .map((id) => nodeMap.get(id))
    .filter((n): n is NonNullable<typeof n> => n != null);
}
