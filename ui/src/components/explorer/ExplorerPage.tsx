import { ReactFlowProvider } from "@xyflow/react";
import { useExplorerGraph } from "@/hooks/useExplorerGraph";
import { useExplorerStore } from "@/store/explorer";
import { getLens } from "@/lib/explorer/lens-config";
import { buildExplorerGraph } from "@/lib/explorer/graph-builder";
import { useMemo } from "react";

export function ExplorerPage() {
  const { data, isLoading, error } = useExplorerGraph();
  const activeLens = useExplorerStore((s) => s.activeLens);
  const subPresets = useExplorerStore((s) => s.subPresets[activeLens] ?? []);

  const built = useMemo(() => {
    if (!data) return null;
    const lens = getLens(activeLens);
    return buildExplorerGraph(
      { nodes: data.nodes, edges: data.edges },
      {
        lens,
        activeLensId: activeLens,
        subPresets,
        findings: data.findings,
      },
    );
  }, [data, activeLens, subPresets]);

  if (error) {
    return (
      <div className="flex h-full items-center justify-center bg-[#050B18]">
        <p className="text-sm text-red-400">
          Failed to load graph: {error.message}
        </p>
      </div>
    );
  }

  if (isLoading || !built) {
    return (
      <div className="flex h-full items-center justify-center bg-[#050B18]">
        <p className="text-sm text-muted-foreground">Loading Explorer…</p>
      </div>
    );
  }

  return (
    <div className="relative h-full w-full bg-[#050B18]">
      <ReactFlowProvider>
        <div className="flex h-full items-center justify-center">
          <div className="text-center font-mono text-xs text-slate-400">
            <p>Active lens: {activeLens}</p>
            <p>
              Visible: {built.metrics.visibleNodeCount} nodes ·{" "}
              {built.metrics.visibleEdgeCount} edges
            </p>
            <p>
              Severity: {built.metrics.criticalCount} critical ·{" "}
              {built.metrics.highCount} high · {built.metrics.mediumCount} med
            </p>
          </div>
        </div>
      </ReactFlowProvider>
    </div>
  );
}
