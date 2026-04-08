import { NODE_COLORS } from "./node-styles";

interface ActiveFilters {
  nodeKinds: Set<string>;
  edgeKinds: Set<string>;
  minRiskScore: number;
}

interface HighlightedPath {
  nodeIds: string[];
  edgeKeys: string[];
}

export function createNodeReducer(
  filters: ActiveFilters,
  hoveredNode: string | null,
  selectedNode: string | null,
  highlightedPath: HighlightedPath | null,
) {
  return (node: string, data: Record<string, unknown>) => {
    const res = { ...data };
    const kind = data._kind as string;

    if (!filters.nodeKinds.has(kind)) {
      res.hidden = true;
      return res;
    }

    const riskScore = Number(data._riskScore ?? 0);
    if (riskScore < filters.minRiskScore) {
      res.hidden = true;
      return res;
    }

    if (highlightedPath) {
      const onPath = highlightedPath.nodeIds.includes(node);
      res.color = onPath ? (NODE_COLORS[kind] ?? "#999") : "#333";
      res.size = onPath
        ? (data.size as number) * 1.5
        : (data.size as number) * 0.4;
      res.zIndex = onPath ? 1 : 0;
      return res;
    }

    if (hoveredNode) {
      if (node === hoveredNode || node === selectedNode) {
        res.zIndex = 1;
      } else {
        res.color = "#333";
        res.size = (data.size as number) * 0.6;
        res.label = "";
      }
    }

    if (node === selectedNode) {
      res.highlighted = true;
      res.zIndex = 2;
    }

    return res;
  };
}

export function createEdgeReducer(
  filters: ActiveFilters,
  hoveredNode: string | null,
  highlightedPath: HighlightedPath | null,
) {
  return (edge: string, data: Record<string, unknown>) => {
    const res = { ...data };
    const kind = data._kind as string;

    if (!filters.edgeKinds.has(kind)) {
      res.hidden = true;
      return res;
    }

    if (highlightedPath) {
      const onPath = highlightedPath.edgeKeys.includes(edge);
      res.color = onPath ? "#FF0000" : "#222";
      res.size = onPath ? 3 : 0.3;
      res.zIndex = onPath ? 1 : 0;
      return res;
    }

    if (hoveredNode) {
      res.color = "#222";
      res.size = 0.5;
    }

    return res;
  };
}
