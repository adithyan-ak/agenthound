import { create } from "zustand";

interface HighlightedPath {
  nodeIds: string[];
  edgeKeys: string[];
  title?: string;
}

interface GraphState {
  selectedNodeId: string | null;
  selectedEdgeId: string | null;
  highlightedPath: HighlightedPath | null;
}

interface GraphActions {
  selectNode: (id: string | null) => void;
  selectEdge: (id: string | null) => void;
  highlightPath: (path: HighlightedPath) => void;
  clearHighlight: () => void;
  clearSelection: () => void;
}

export const useGraphStore = create<GraphState & GraphActions>()((set) => ({
  selectedNodeId: null,
  selectedEdgeId: null,
  highlightedPath: null,

  selectNode: (id) => set({ selectedNodeId: id, selectedEdgeId: null }),

  selectEdge: (id) => set({ selectedEdgeId: id, selectedNodeId: null }),

  highlightPath: (path) => set({ highlightedPath: path }),

  clearHighlight: () => set({ highlightedPath: null }),

  clearSelection: () =>
    set({
      selectedNodeId: null,
      selectedEdgeId: null,
      highlightedPath: null,
    }),
}));
