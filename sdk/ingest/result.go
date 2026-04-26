package ingest

import "time"

type IngestResult struct {
	ScanID              string               `json:"scan_id"`
	NodesWritten        int                  `json:"nodes_written"`
	EdgesWritten        int                  `json:"edges_written"`
	Warnings            []string             `json:"warnings,omitempty"`
	PostProcessingStats []PostProcessingStat `json:"post_processing_stats,omitempty"`
	Duration            time.Duration        `json:"duration"`
}

type PostProcessingStat struct {
	ProcessorName string        `json:"processor_name"`
	EdgesCreated  int           `json:"edges_created"`
	NodesUpdated  int           `json:"nodes_updated"`
	Duration      time.Duration `json:"duration"`
	Error         string        `json:"error,omitempty"`
}
