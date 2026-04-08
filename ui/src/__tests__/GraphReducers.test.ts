import { describe, it, expect } from "vitest";
import { createNodeReducer, createEdgeReducer } from "@/lib/graph-reducers";

const allKinds = new Set(["MCPServer", "MCPTool", "AgentInstance"]);
const allEdgeKinds = new Set(["TRUSTS_SERVER", "PROVIDES_TOOL", "CAN_REACH"]);

const baseFilters = {
  nodeKinds: allKinds,
  edgeKinds: allEdgeKinds,
  minRiskScore: 0,
};

const nodeData = (kind: string, size = 10, riskScore = 50) => ({
  _kind: kind,
  _riskScore: riskScore,
  size,
  color: "#50C878",
  label: "test-node",
});

const edgeData = (kind: string) => ({
  _kind: kind,
  color: "#ccc",
  size: 1,
});

describe("createNodeReducer", () => {
  it("passes through nodes matching filters", () => {
    const reducer = createNodeReducer(baseFilters, null, null, null);
    const result = reducer("n1", nodeData("MCPServer"));
    expect(result.hidden).toBeUndefined();
    expect(result.color).toBe("#50C878");
  });

  it("hides nodes whose kind is filtered out", () => {
    const filters = { ...baseFilters, nodeKinds: new Set(["MCPTool"]) };
    const reducer = createNodeReducer(filters, null, null, null);
    const result = reducer("n1", nodeData("MCPServer"));
    expect(result.hidden).toBe(true);
  });

  it("hides nodes below minRiskScore", () => {
    const filters = { ...baseFilters, minRiskScore: 60 };
    const reducer = createNodeReducer(filters, null, null, null);
    const result = reducer("n1", nodeData("MCPServer", 10, 40));
    expect(result.hidden).toBe(true);
  });

  it("does not hide nodes at or above minRiskScore", () => {
    const filters = { ...baseFilters, minRiskScore: 50 };
    const reducer = createNodeReducer(filters, null, null, null);
    const result = reducer("n1", nodeData("MCPServer", 10, 50));
    expect(result.hidden).toBeUndefined();
  });

  it("highlights path nodes and dims non-path nodes", () => {
    const path = { nodeIds: ["n1"], edgeKeys: [] };
    const reducer = createNodeReducer(baseFilters, null, null, path);

    const onPath = reducer("n1", nodeData("MCPServer", 10));
    expect(onPath.size).toBe(15);
    expect(onPath.zIndex).toBe(1);
    expect(onPath.color).toBe("#50C878");

    const offPath = reducer("n2", nodeData("MCPTool", 10));
    expect(offPath.size).toBe(4);
    expect(offPath.zIndex).toBe(0);
    expect(offPath.color).toBe("#333");
  });

  it("dims non-hovered nodes when a node is hovered", () => {
    const reducer = createNodeReducer(baseFilters, "n1", null, null);

    const hovered = reducer("n1", nodeData("MCPServer", 10));
    expect(hovered.zIndex).toBe(1);
    expect(hovered.color).toBe("#50C878");

    const other = reducer("n2", nodeData("MCPTool", 10));
    expect(other.color).toBe("#333");
    expect(other.size).toBe(6);
    expect(other.label).toBe("");
  });

  it("marks selected node as highlighted", () => {
    const reducer = createNodeReducer(baseFilters, null, "n1", null);
    const result = reducer("n1", nodeData("MCPServer"));
    expect(result.highlighted).toBe(true);
    expect(result.zIndex).toBe(2);
  });

  it("preserves hovered + selected node visibility together", () => {
    const reducer = createNodeReducer(baseFilters, "n2", "n1", null);
    const selected = reducer("n1", nodeData("MCPServer", 10));
    expect(selected.zIndex).toBe(2);
    expect(selected.highlighted).toBe(true);
  });
});

describe("createEdgeReducer", () => {
  it("passes through edges matching filters", () => {
    const reducer = createEdgeReducer(baseFilters, null, null);
    const result = reducer("e1", edgeData("TRUSTS_SERVER"));
    expect(result.hidden).toBeUndefined();
  });

  it("hides edges whose kind is filtered out", () => {
    const filters = { ...baseFilters, edgeKinds: new Set(["PROVIDES_TOOL"]) };
    const reducer = createEdgeReducer(filters, null, null);
    const result = reducer("e1", edgeData("TRUSTS_SERVER"));
    expect(result.hidden).toBe(true);
  });

  it("highlights path edges red and dims non-path edges", () => {
    const path = { nodeIds: [], edgeKeys: ["e1"] };
    const reducer = createEdgeReducer(baseFilters, null, path);

    const onPath = reducer("e1", edgeData("TRUSTS_SERVER"));
    expect(onPath.color).toBe("#FF0000");
    expect(onPath.size).toBe(3);

    const offPath = reducer("e2", edgeData("PROVIDES_TOOL"));
    expect(offPath.color).toBe("#222");
    expect(offPath.size).toBe(0.3);
  });

  it("dims edges when a node is hovered", () => {
    const reducer = createEdgeReducer(baseFilters, "n1", null);
    const result = reducer("e1", edgeData("TRUSTS_SERVER"));
    expect(result.color).toBe("#222");
    expect(result.size).toBe(0.5);
  });
});
