import { api } from "@shared/api/client";
import type { Scan } from "./model";

export interface IngestResult {
  scan_id: string;
  nodes_written: number;
  edges_written: number;
  warnings?: string[];
  post_processing_stats?: Array<{
    processor_name: string;
    edges_created?: number;
  }>;
  duration?: number;
}

export async function fetchScans(limit = 50, offset = 0): Promise<Scan[]> {
  return api
    .get("scans", {
      searchParams: { limit: String(limit), offset: String(offset) },
    })
    .json<Scan[]>();
}

export async function deleteScan(id: string): Promise<void> {
  await api.delete(`scans/${encodeURIComponent(id)}`);
}

// uploadScan POSTs a collector JSON file to /api/v1/ingest. The file is read
// as text and posted as the raw request body with Content-Type:
// application/json (matches the existing handler contract).
export async function uploadScan(file: File): Promise<IngestResult> {
  const text = await readFileAsText(file);
  return api
    .post("ingest", {
      body: text,
      headers: { "Content-Type": "application/json" },
      timeout: 120_000,
    })
    .json<IngestResult>();
}

function readFileAsText(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result ?? ""));
    reader.onerror = () => reject(reader.error ?? new Error("read failed"));
    reader.readAsText(file);
  });
}
