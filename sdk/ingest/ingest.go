package ingest

type IngestData struct {
	Meta  IngestMeta `json:"meta"`
	Graph GraphData  `json:"graph"`
}

type IngestMeta struct {
	Version          int    `json:"version"`
	Type             string `json:"type"`
	Collector        string `json:"collector"`
	CollectorVersion string `json:"collector_version"`
	Timestamp        string `json:"timestamp"`
	ScanID           string `json:"scan_id"`
	// Extensions is a wire-format-only forward-compat hook for declaring future schema namespaces. Validators ignore unknown values today; future producers may use it to declare extension namespaces for new node/edge kinds. Optional; existing producers may omit it.
	Extensions []string `json:"extensions,omitempty"`
}

type GraphData struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
