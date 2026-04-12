import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  ReactFlow,
  ReactFlowProvider,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  useReactFlow,
  type Node,
  type Edge,
  type NodeMouseHandler,
  type EdgeMouseHandler,
  BackgroundVariant,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { useGraphData } from "@/hooks/useGraph";
import { computeLayout } from "@/lib/layout";
import { useGraphStore } from "@/store/graph";
import { useUIStore } from "@/store/ui";
import { GraphFilters } from "./GraphFilters";
import { GraphSearch } from "./GraphSearch";
import { GraphPathfinder } from "./GraphPathfinder";
import {
  NodeContextMenu,
  type ContextMenuState,
} from "./NodeContextMenu";

import { ServerNode } from "./nodes/ServerNode";
import { ToolNode } from "./nodes/ToolNode";
import { ResourceNode } from "./nodes/ResourceNode";
import { A2AAgentNode } from "./nodes/A2AAgentNode";
import { SkillNode } from "./nodes/SkillNode";
import { InfraNode } from "./nodes/InfraNode";
import { AttackEdge } from "./edges/AttackEdge";
import { TrustEdge } from "./edges/TrustEdge";
import { StructureEdge } from "./edges/StructureEdge";

const nodeTypes = {
  server: ServerNode,
  tool: ToolNode,
  resource: ResourceNode,
  a2aAgent: A2AAgentNode,
  skill: SkillNode,
  infra: InfraNode,
};

const edgeTypes = {
  attack: AttackEdge,
  trust: TrustEdge,
  structure: StructureEdge,
};

function minimapColor(node: Node): string {
  return ((node.data as Record<string, unknown>)?.color as string) ?? "#666";
}

function GraphCanvas() {
  const { data, isLoading, error } = useGraphData();
  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);
  const [layoutReady, setLayoutReady] = useState(false);
  const [contextMenu, setContextMenu] = useState<ContextMenuState | null>(null);
  const reactFlow = useReactFlow();
  const reactFlowRef = useRef(reactFlow);
  reactFlowRef.current = reactFlow;
  const hasInitialLayoutRef = useRef(false);

  const selectNode = useGraphStore((s) => s.selectNode);
  const selectEdge = useGraphStore((s) => s.selectEdge);
  const clearSelection = useGraphStore((s) => s.clearSelection);
  const highlightPath = useGraphStore((s) => s.highlightPath);
  const filters = useGraphStore((s) => s.activeFilters);
  const highlightedPath = useGraphStore((s) => s.highlightedPath);
  const openSidebar = useUIStore((s) => s.openSidebar);

  // Build layout when data changes.
  // IMPORTANT: do NOT flip layoutReady back to false on refetches —
  // that causes the canvas to unmount and flash "Loading graph...".
  // Also do NOT put reactFlow in deps (use ref instead) — the ReactFlowInstance
  // reference can be unstable across renders which would re-trigger the effect.
  useEffect(() => {
    if (!data) return;
    let cancelled = false;
    computeLayout(data.nodes, data.edges).then((positioned) => {
      if (cancelled) return;
      setNodes(positioned);
      setEdges(data.edges);
      if (!hasInitialLayoutRef.current) {
        hasInitialLayoutRef.current = true;
        setLayoutReady(true);
        setTimeout(
          () => reactFlowRef.current.fitView({ padding: 0.15, duration: 300 }),
          100,
        );
      }
    });
    return () => {
      cancelled = true;
    };
  }, [data, setNodes, setEdges]);

  // Fit view to highlighted path when the path changes.
  // Read current nodes via reactFlow.getNodes() so we don't depend on the
  // local `nodes` state (which updates on every refetch).
  useEffect(() => {
    if (!highlightedPath || !layoutReady) return;
    const pathNodeIds = new Set(highlightedPath.nodeIds);
    const timer = setTimeout(() => {
      const currentNodes = reactFlowRef.current.getNodes();
      const pathNodes = currentNodes.filter((n) => pathNodeIds.has(n.id));
      if (pathNodes.length === 0) return;
      reactFlowRef.current.fitView({
        nodes: pathNodes.map((n) => ({ id: n.id })),
        padding: 0.35,
        duration: 600,
        maxZoom: 1.2,
      });
    }, 120);
    return () => clearTimeout(timer);
  }, [highlightedPath, layoutReady]);

  // Keyboard: Esc clears selection + highlight + context menu
  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      const target = e.target as HTMLElement;
      if (target?.tagName === "INPUT" || target?.tagName === "TEXTAREA") return;
      if (e.key === "Escape") {
        setContextMenu(null);
        clearSelection();
      }
    }
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [clearSelection]);

  const filteredNodes = useMemo(() => {
    const pathNodeIds = highlightedPath
      ? new Set(highlightedPath.nodeIds)
      : null;
    return nodes.map((n) => {
      const d = n.data as Record<string, unknown>;
      const kind = d?.kind as string;
      const riskScore = (d?.riskScore as number) ?? 0;
      const hidden =
        !filters.nodeKinds.has(kind) || riskScore < filters.minRiskScore;
      let style = n.style;
      if (pathNodeIds && !pathNodeIds.has(n.id)) {
        style = { ...style, opacity: 0.15 };
      } else if (pathNodeIds) {
        style = { ...style, opacity: 1 };
      }
      return { ...n, hidden, style };
    });
  }, [nodes, filters, highlightedPath]);

  const filteredEdges = useMemo(() => {
    const pathEdgeIds = highlightedPath
      ? new Set(highlightedPath.edgeKeys)
      : null;
    return edges.map((e) => {
      const kind = (e.data as Record<string, unknown>)?.kind as string;
      const hidden = !filters.edgeKinds.has(kind);
      let animated = false;
      let style = e.style;
      if (pathEdgeIds) {
        if (pathEdgeIds.has(e.id)) {
          animated = true;
          style = { ...style, opacity: 1 };
        } else {
          style = { ...style, opacity: 0.08 };
        }
      }
      return { ...e, hidden, animated, style };
    });
  }, [edges, filters, highlightedPath]);

  const onNodeClick: NodeMouseHandler = useCallback(
    (_, node) => {
      selectNode(node.id);
      openSidebar();
    },
    [selectNode, openSidebar],
  );

  const onEdgeClick: EdgeMouseHandler = useCallback(
    (_, edge) => {
      selectEdge(edge.id);
      openSidebar();
    },
    [selectEdge, openSidebar],
  );

  const onPaneClick = useCallback(() => {
    clearSelection();
    setContextMenu(null);
  }, [clearSelection]);

  const onNodeContextMenu = useCallback(
    (e: React.MouseEvent, node: Node) => {
      e.preventDefault();
      const d = node.data as Record<string, unknown>;
      setContextMenu({
        nodeId: node.id,
        nodeLabel: String(d?.label ?? node.id.slice(0, 12)),
        nodeKind: String(d?.kind ?? "Unknown"),
        x: e.clientX,
        y: e.clientY,
      });
    },
    [],
  );

  const onSearchSelect = useCallback(
    (nodeId: string) => {
      selectNode(nodeId);
      openSidebar();
      const node = nodes.find((n) => n.id === nodeId);
      if (node) {
        reactFlow.fitView({
          nodes: [{ id: nodeId }],
          padding: 0.8,
          duration: 500,
          maxZoom: 1.5,
        });
      }
    },
    [nodes, reactFlow, selectNode, openSidebar],
  );

  const onContextFocus = useCallback(
    (nodeId: string) => {
      // BFS 2-hop neighborhood highlight
      const neighborIds = new Set<string>([nodeId]);
      const edgeKeys = new Set<string>();
      const frontier = [nodeId];
      for (let hop = 0; hop < 2; hop++) {
        const next: string[] = [];
        for (const current of frontier) {
          for (const e of edges) {
            if (e.source === current && !neighborIds.has(e.target)) {
              neighborIds.add(e.target);
              edgeKeys.add(e.id);
              next.push(e.target);
            } else if (e.target === current && !neighborIds.has(e.source)) {
              neighborIds.add(e.source);
              edgeKeys.add(e.id);
              next.push(e.source);
            } else if (
              (e.source === current || e.target === current) &&
              neighborIds.has(e.source) &&
              neighborIds.has(e.target)
            ) {
              edgeKeys.add(e.id);
            }
          }
        }
        frontier.splice(0, frontier.length, ...next);
      }
      highlightPath({
        nodeIds: Array.from(neighborIds),
        edgeKeys: Array.from(edgeKeys),
        title: "2-hop neighborhood",
      });
      selectNode(nodeId);
      openSidebar();
    },
    [edges, highlightPath, selectNode, openSidebar],
  );

  const onContextBlastRadius = useCallback(
    (nodeId: string) => {
      // BFS outbound only
      const visitedNodes = new Set<string>([nodeId]);
      const visitedEdges = new Set<string>();
      const queue = [nodeId];
      let maxHops = 6;
      while (queue.length > 0 && maxHops-- > 0) {
        const nextQueue: string[] = [];
        for (const current of queue) {
          for (const e of edges) {
            if (e.source === current && !visitedNodes.has(e.target)) {
              visitedNodes.add(e.target);
              visitedEdges.add(e.id);
              nextQueue.push(e.target);
            }
          }
        }
        queue.splice(0, queue.length, ...nextQueue);
      }
      highlightPath({
        nodeIds: Array.from(visitedNodes),
        edgeKeys: Array.from(visitedEdges),
        title: "Outbound blast radius",
      });
      selectNode(nodeId);
      openSidebar();
    },
    [edges, highlightPath, selectNode, openSidebar],
  );

  const onContextInbound = useCallback(
    (nodeId: string) => {
      const visitedNodes = new Set<string>([nodeId]);
      const visitedEdges = new Set<string>();
      const queue = [nodeId];
      let maxHops = 6;
      while (queue.length > 0 && maxHops-- > 0) {
        const nextQueue: string[] = [];
        for (const current of queue) {
          for (const e of edges) {
            if (e.target === current && !visitedNodes.has(e.source)) {
              visitedNodes.add(e.source);
              visitedEdges.add(e.id);
              nextQueue.push(e.source);
            }
          }
        }
        queue.splice(0, queue.length, ...nextQueue);
      }
      highlightPath({
        nodeIds: Array.from(visitedNodes),
        edgeKeys: Array.from(visitedEdges),
        title: "Inbound reach",
      });
      selectNode(nodeId);
      openSidebar();
    },
    [edges, highlightPath, selectNode, openSidebar],
  );

  if (error) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-destructive">
          Error loading graph: {error.message}
        </p>
      </div>
    );
  }

  if (!layoutReady) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-sm text-muted-foreground">
          {isLoading ? "Loading graph..." : "Computing layout..."}
        </p>
      </div>
    );
  }

  return (
    <>
      <ReactFlow
        nodes={filteredNodes}
        edges={filteredEdges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={onNodeClick}
        onEdgeClick={onEdgeClick}
        onPaneClick={onPaneClick}
        onNodeContextMenu={onNodeContextMenu}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        fitView
        minZoom={0.05}
        maxZoom={2}
        proOptions={{ hideAttribution: true }}
        defaultEdgeOptions={{ type: "structure" }}
        onlyRenderVisibleElements
      >
        <Background
          variant={BackgroundVariant.Dots}
          gap={24}
          size={1}
          color="#1a2030"
        />
        <Controls
          position="bottom-right"
          showInteractive={false}
          className="!bg-card !border-border !shadow-md"
        />
        <MiniMap
          nodeColor={minimapColor}
          maskColor="rgba(0,0,0,0.7)"
          className="!bg-card !border-border"
          style={{ height: 120, width: 200 }}
          position="bottom-left"
        />
      </ReactFlow>

      <GraphSearch nodes={nodes} onSelect={onSearchSelect} />

      <NodeContextMenu
        state={contextMenu}
        onClose={() => setContextMenu(null)}
        onFocus={onContextFocus}
        onShowBlastRadius={onContextBlastRadius}
        onShowInbound={onContextInbound}
      />
    </>
  );
}

export function GraphExplorer() {
  return (
    <div className="relative h-full w-full">
      <ReactFlowProvider>
        <GraphCanvas />
        <GraphFilters />
        <GraphPathfinder />
      </ReactFlowProvider>
    </div>
  );
}
