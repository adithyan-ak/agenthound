import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { describe, it, expect, beforeEach, vi } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter } from "react-router-dom";
import { Pathfinder } from "@/components/pathfinder/Pathfinder";
import type { PathResponse } from "@/api/types";

vi.mock("@/api/analysis", () => ({
  findShortestPath: vi.fn(),
  findAllPaths: vi.fn(),
  findWeightedPath: vi.fn(),
  fetchFindings: vi.fn(),
  fetchPreBuiltQueries: vi.fn(),
  runPreBuiltQuery: vi.fn(),
}));

import { findShortestPath } from "@/api/analysis";

const mockedFindShortest = vi.mocked(findShortestPath);

const mockResult: PathResponse = {
  paths: [
    {
      nodes: [
        { id: "n1", name: "claude-desktop", kinds: ["AgentInstance"] },
        { id: "n2", name: "filesystem-server", kinds: ["MCPServer"] },
        { id: "n3", name: "read_file", kinds: ["MCPTool"] },
      ],
      edges: [
        { kind: "TRUSTS_SERVER", source: "n1", target: "n2" },
        { kind: "PROVIDES_TOOL", source: "n2", target: "n3" },
      ],
      hops: 2,
    },
  ],
};

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return function Wrapper({ children }: { children: React.ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>{children}</MemoryRouter>
      </QueryClientProvider>
    );
  };
}

describe("Pathfinder", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders form elements and instructions", () => {
    render(<Pathfinder />, { wrapper: createWrapper() });

    expect(screen.getByText("Pathfinder")).toBeInTheDocument();
    expect(screen.getByText("Attack Path Analysis")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("e.g. claude-desktop")).toBeInTheDocument();
    expect(screen.getByText("Find Paths")).toBeInTheDocument();
    expect(screen.getByText("Source Kind")).toBeInTheDocument();
    expect(screen.getByText("Algorithm")).toBeInTheDocument();
  });

  it("disables submit when source name is empty", () => {
    render(<Pathfinder />, { wrapper: createWrapper() });

    const submitButton = screen.getByRole("button", { name: /find paths/i });
    expect(submitButton).toBeDisabled();
  });

  it("enables submit when source name is filled", () => {
    render(<Pathfinder />, { wrapper: createWrapper() });

    const input = screen.getByPlaceholderText("e.g. claude-desktop");
    fireEvent.change(input, { target: { value: "my-agent" } });

    const submitButton = screen.getByRole("button", { name: /find paths/i });
    expect(submitButton).not.toBeDisabled();
  });

  it("submits form and displays path results", async () => {
    mockedFindShortest.mockResolvedValue(mockResult);

    render(<Pathfinder />, { wrapper: createWrapper() });

    const input = screen.getByPlaceholderText("e.g. claude-desktop");
    fireEvent.change(input, { target: { value: "claude-desktop" } });

    const form = input.closest("form")!;
    fireEvent.submit(form);

    await waitFor(() => {
      expect(screen.getByText("1 path found")).toBeInTheDocument();
    });

    expect(screen.getByText("claude-desktop")).toBeInTheDocument();
    expect(screen.getByText("filesystem-server")).toBeInTheDocument();
    expect(screen.getByText("read_file")).toBeInTheDocument();
    expect(screen.getByText("TRUSTS_SERVER")).toBeInTheDocument();
    expect(screen.getByText("PROVIDES_TOOL")).toBeInTheDocument();
    expect(screen.getByText("2 hops")).toBeInTheDocument();
  });

  it("shows empty state when no paths found", async () => {
    mockedFindShortest.mockResolvedValue({ paths: [] });

    render(<Pathfinder />, { wrapper: createWrapper() });

    const input = screen.getByPlaceholderText("e.g. claude-desktop");
    fireEvent.change(input, { target: { value: "no-paths-agent" } });

    fireEvent.submit(input.closest("form")!);

    await waitFor(() => {
      expect(screen.getByText("No paths found")).toBeInTheDocument();
    });
  });

  it("renders algorithm radio buttons", () => {
    render(<Pathfinder />, { wrapper: createWrapper() });

    const radios = screen.getAllByRole("radio");
    expect(radios).toHaveLength(3);
    expect(screen.getByLabelText(/shortest/i)).toBeChecked();
  });
});
