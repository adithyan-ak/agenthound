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

export function buildReactFlowGraph(
  apiNodes: APINode[],
  apiEdges: APIEdge[],
): { nodes: Node[]; edges: Edge[] } {
  const nodeMap = new Map<string, APINode>();
  for (const n of apiNodes) nodeMap.set(n.id, n);

  const nodes: Node[] = [];

  for (const node of apiNodes) {
    const kind = node.kinds[0] ?? "Unknown";
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

  const edges: Edge[] = [];
  const seen = new Set<string>();

  for (const edge of apiEdges) {
    if (!nodeMap.has(edge.source) || !nodeMap.has(edge.target)) continue;

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

  return { nodes, edges };
}
