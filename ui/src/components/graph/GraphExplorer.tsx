import { MultiDirectedGraph } from "graphology";
import { SigmaContainer } from "@react-sigma/core";
import "@react-sigma/core/lib/style.css";
import { GraphDataLoader } from "./GraphDataLoader";
import { GraphEvents } from "./GraphEvents";
import { GraphControls } from "./GraphControls";
import { GraphSearch } from "./GraphSearch";
import { GraphFilters } from "./GraphFilters";
import { GraphLegend } from "./GraphLegend";
import { useGraphStore } from "@/store/graph";
import { createNodeReducer, createEdgeReducer } from "@/lib/graph-reducers";

const sigmaSettings = {
  defaultNodeColor: "#999",
  defaultEdgeColor: "#ccc",
  labelRenderedSizeThreshold: 12,
  renderEdgeLabels: false,
  enableEdgeEvents: true,
  labelFont: "Inter, system-ui, sans-serif",
  labelSize: 12,
  labelColor: { color: "#333" },
  edgeLabelFont: "Inter, system-ui, sans-serif",
  edgeLabelSize: 10,
  nodeReducer: undefined as
    | ((node: string, data: Record<string, unknown>) => Record<string, unknown>)
    | undefined,
  edgeReducer: undefined as
    | ((
        edge: string,
        data: Record<string, unknown>,
      ) => Record<string, unknown>)
    | undefined,
};

export function GraphExplorer() {
  const hoveredNode = useGraphStore((s) => s.hoveredNodeId);
  const selectedNode = useGraphStore((s) => s.selectedNodeId);
  const highlightedPath = useGraphStore((s) => s.highlightedPath);
  const filters = useGraphStore((s) => s.activeFilters);

  const settings = {
    ...sigmaSettings,
    nodeReducer: createNodeReducer(filters, hoveredNode, selectedNode, highlightedPath),
    edgeReducer: createEdgeReducer(filters, hoveredNode, highlightedPath),
  };

  return (
    <div className="relative h-full w-full">
      <SigmaContainer
        graph={MultiDirectedGraph}
        settings={settings}
        className="h-full w-full"
      >
        <GraphDataLoader />
        <GraphEvents />
        <GraphSearch />
      </SigmaContainer>
      <GraphControls />
      <GraphFilters />
      <GraphLegend />
    </div>
  );
}
