import type { Node, Edge } from "@xyflow/react";
import type { APINode, APIEdge } from "@/api/types";
import { getNodeColor, getNodeSize, getNodeLabel } from "./node-styles";
import { getEdgeCategory } from "./edge-styles";

const NODE_TYPE_MAP: Record<string, string> = {
  AgentInstance: "server",
  MCPServer: "server",
  MCPTool: "tool",
  MCPPrompt: "tool",
  MCPResource: "resource",
  A2AAgent: "a2aAgent",
  A2ASkill: "skill",
  Identity: "infra",
  Credential: "infra",
  Host: "infra",
  ConfigFile: "infra",
  InstructionFile: "infra",
  ResourceGroup: "infra",
  TrustZone: "infra",
};

const COLLAPSIBLE_KINDS = new Set(["MCPTool", "MCPPrompt", "A2ASkill"]);
const PARENT_EDGE_KINDS = new Set([
  "PROVIDES_TOOL",
  "PROVIDES_RESOURCE",
  "PROVIDES_PROMPT",
  "ADVERTISES_SKILL",
]);
const MAX_LEAVES_PER_PARENT = 5;

export function buildReactFlowGraph(
  apiNodes: APINode[],
  apiEdges: APIEdge[],
): { nodes: Node[]; edges: Edge[] } {
  const nodeMap = new Map<string, APINode>();
  for (const n of apiNodes) nodeMap.set(n.id, n);

  const childrenByParent = new Map<string, string[]>();
  for (const edge of apiEdges) {
    if (PARENT_EDGE_KINDS.has(edge.kind)) {
      if (!childrenByParent.has(edge.source))
        childrenByParent.set(edge.source, []);
      childrenByParent.get(edge.source)!.push(edge.target);
    }
  }

  const connectedNodes = new Set<string>();
  for (const edge of apiEdges) {
    connectedNodes.add(edge.source);
    connectedNodes.add(edge.target);
  }

  const collapsedIds = new Set<string>();

  for (const [, childIds] of childrenByParent) {
    const collapsible = childIds.filter((id) => {
      const kind = nodeMap.get(id)?.kinds[0] ?? "";
      return COLLAPSIBLE_KINDS.has(kind);
    });

    if (collapsible.length > MAX_LEAVES_PER_PARENT) {
      for (let i = MAX_LEAVES_PER_PARENT; i < collapsible.length; i++) {
        collapsedIds.add(collapsible[i]!);
      }
    }
  }

  const nodes: Node[] = [];

  for (const node of apiNodes) {
    const kind = node.kinds[0] ?? "Unknown";

    if (collapsedIds.has(node.id)) continue;

    const nodeType = NODE_TYPE_MAP[kind] ?? "infra";

    nodes.push({
      id: node.id,
      type: nodeType,
      position: { x: 0, y: 0 },
      data: {
        label: getNodeLabel(node),
        kind,
        color: getNodeColor(node.kinds),
        size: getNodeSize(node),
        riskScore: Number(node.properties.risk_score ?? 0),
        properties: node.properties,
      },
    });
  }

  for (const [parentId, childIds] of childrenByParent) {
    const collapsible = childIds.filter((id) => {
      const kind = nodeMap.get(id)?.kinds[0] ?? "";
      return COLLAPSIBLE_KINDS.has(kind);
    });

    if (collapsible.length > MAX_LEAVES_PER_PARENT) {
      const overflowCount = collapsible.length - MAX_LEAVES_PER_PARENT;
      nodes.push({
        id: `collapse-${parentId}`,
        type: "tool",
        position: { x: 0, y: 0 },
        data: {
          label: `+${overflowCount} more`,
          kind: "MCPTool",
          color: "#F5A623",
          size: 4,
          riskScore: 0,
          properties: {},
          isOverflow: true,
          overflowCount,
        },
      });
    }
  }

  const edges: Edge[] = [];
  const seen = new Set<string>();

  for (const edge of apiEdges) {
    if (!nodeMap.has(edge.source) || !nodeMap.has(edge.target)) continue;
    if (collapsedIds.has(edge.source) || collapsedIds.has(edge.target)) continue;

    const key = `${edge.source}->${edge.target}:${edge.kind}`;
    if (seen.has(key)) continue;
    seen.add(key);

    edges.push({
      id: key,
      source: edge.source,
      target: edge.target,
      type: getEdgeCategory(edge.kind),
      data: { kind: edge.kind, properties: edge.properties },
    });
  }

  for (const [parentId, childIds] of childrenByParent) {
    const collapsible = childIds.filter((id) => {
      const kind = nodeMap.get(id)?.kinds[0] ?? "";
      return COLLAPSIBLE_KINDS.has(kind);
    });
    if (collapsible.length > MAX_LEAVES_PER_PARENT) {
      const collapseNodeId = `collapse-${parentId}`;
      const edgeKey = `${parentId}->${collapseNodeId}:COLLAPSED`;
      if (!seen.has(edgeKey)) {
        seen.add(edgeKey);
        edges.push({
          id: edgeKey,
          source: parentId,
          target: collapseNodeId,
          type: "structure",
          data: { kind: "COLLAPSED", properties: {} },
        });
      }
    }
  }

  return { nodes, edges };
}
