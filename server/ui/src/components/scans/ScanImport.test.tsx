import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { describe, it, expect, beforeEach, vi } from "vitest";
import { ScanImport } from "./ScanImport";

vi.mock("@/api/scans", () => ({
  uploadScan: vi.fn(),
}));

import { uploadScan } from "@/api/scans";

const mockedUploadScan = vi.mocked(uploadScan);

function makeJSONFile(name: string, content: string): File {
  return new File([content], name, { type: "application/json" });
}

const validScanJSON = JSON.stringify({
  meta: {
    version: 1,
    type: "agenthound-ingest",
    collector: "config",
    collector_version: "0.1.0",
    timestamp: "2026-04-23T12:00:00Z",
    scan_id: "test-scan-1",
  },
  graph: { nodes: [], edges: [] },
});

describe("ScanImport", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("uploads a dropped JSON file and calls onSuccess", async () => {
    mockedUploadScan.mockResolvedValue({
      scan_id: "test-scan-1",
      nodes_written: 5,
      edges_written: 3,
    });
    const onSuccess = vi.fn();

    render(<ScanImport open={true} onClose={() => {}} onSuccess={onSuccess} />);

    const dropzone = screen.getByTestId("dropzone");
    const file = makeJSONFile("scan.json", validScanJSON);

    fireEvent.drop(dropzone, {
      dataTransfer: { files: [file] },
    });

    await waitFor(() => {
      expect(mockedUploadScan).toHaveBeenCalledWith(file);
    });
    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled();
    });
    expect(screen.getByText(/imported scan\.json/i)).toBeInTheDocument();
    expect(screen.getByText(/5 nodes, 3 edges written/i)).toBeInTheDocument();
  });

  it("shows an error and does not upload when the file is not valid JSON", async () => {
    const onSuccess = vi.fn();
    render(<ScanImport open={true} onClose={() => {}} onSuccess={onSuccess} />);

    const dropzone = screen.getByTestId("dropzone");
    const badFile = makeJSONFile("scan.json", "not json {{");

    fireEvent.drop(dropzone, {
      dataTransfer: { files: [badFile] },
    });

    await waitFor(() => {
      expect(screen.getByText(/import failed/i)).toBeInTheDocument();
    });
    expect(screen.getByText(/not valid json/i)).toBeInTheDocument();
    expect(mockedUploadScan).not.toHaveBeenCalled();
    expect(onSuccess).not.toHaveBeenCalled();
  });

  it("shows an error when the server rejects the upload", async () => {
    mockedUploadScan.mockRejectedValue(
      new Error("server error (500): check server logs"),
    );
    const onSuccess = vi.fn();

    render(<ScanImport open={true} onClose={() => {}} onSuccess={onSuccess} />);

    const dropzone = screen.getByTestId("dropzone");
    const file = makeJSONFile("scan.json", validScanJSON);

    fireEvent.drop(dropzone, {
      dataTransfer: { files: [file] },
    });

    await waitFor(() => {
      expect(screen.getByText(/import failed/i)).toBeInTheDocument();
    });
    expect(
      screen.getByText(/server error \(500\): check server logs/i),
    ).toBeInTheDocument();
    expect(onSuccess).not.toHaveBeenCalled();
  });
});
