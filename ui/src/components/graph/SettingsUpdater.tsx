import { useEffect, useMemo } from "react";
import { useSetSettings } from "@react-sigma/core";
import { useGraphStore } from "@/store/graph";
import { createNodeReducer, createEdgeReducer } from "@/lib/graph-reducers";

export function SettingsUpdater() {
  const setSettings = useSetSettings();
  const hoveredNode = useGraphStore((s) => s.hoveredNodeId);
  const selectedNode = useGraphStore((s) => s.selectedNodeId);
  const highlightedPath = useGraphStore((s) => s.highlightedPath);
  const filters = useGraphStore((s) => s.activeFilters);

  const nodeReducer = useMemo(
    () => createNodeReducer(filters, hoveredNode, selectedNode, highlightedPath),
    [filters, hoveredNode, selectedNode, highlightedPath],
  );

  const edgeReducer = useMemo(
    () => createEdgeReducer(filters, hoveredNode, highlightedPath),
    [filters, hoveredNode, highlightedPath],
  );

  useEffect(() => {
    setSettings({ nodeReducer, edgeReducer });
  }, [setSettings, nodeReducer, edgeReducer]);

  return null;
}
