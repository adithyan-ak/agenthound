import { Component, type ReactNode } from "react";
import { MultiDirectedGraph } from "graphology";
import { SigmaContainer } from "@react-sigma/core";
import "@react-sigma/core/lib/style.css";
import { GraphDataLoader } from "./GraphDataLoader";
import { GraphEvents } from "./GraphEvents";
import { GraphControls } from "./GraphControls";
import { GraphSearch } from "./GraphSearch";
import { GraphFilters } from "./GraphFilters";
import { GraphLegend } from "./GraphLegend";
import { SettingsUpdater } from "./SettingsUpdater";

const sigmaSettings = {
  defaultNodeColor: "#999",
  defaultEdgeColor: "#ccc",
  labelRenderedSizeThreshold: 12,
  renderEdgeLabels: false,
  enableEdgeEvents: true,
  labelFont: "Inter, system-ui, sans-serif",
  labelSize: 12,
  labelColor: { color: "#333" } as const,
  edgeLabelFont: "Inter, system-ui, sans-serif",
  edgeLabelSize: 10,
};

interface ErrorBoundaryState {
  error: Error | null;
}

class GraphErrorBoundary extends Component<
  { children: ReactNode },
  ErrorBoundaryState
> {
  state: ErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error) {
    return { error };
  }

  render() {
    if (this.state.error) {
      return (
        <div className="absolute inset-0 flex items-center justify-center bg-background">
          <div className="text-center space-y-2 p-4">
            <p className="text-destructive font-medium">Graph rendering error</p>
            <pre className="text-xs text-muted-foreground max-w-lg whitespace-pre-wrap">
              {this.state.error.message}
            </pre>
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}

export function GraphExplorer() {
  return (
    <div className="relative h-full w-full">
      <GraphErrorBoundary>
        <SigmaContainer
          graph={MultiDirectedGraph}
          settings={sigmaSettings}
          className="h-full w-full"
        >
          <SettingsUpdater />
          <GraphDataLoader />
          <GraphEvents />
          <GraphSearch />
          <GraphControls />
        </SigmaContainer>
      </GraphErrorBoundary>
      <GraphFilters />
      <GraphLegend />
    </div>
  );
}
