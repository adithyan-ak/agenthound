package ingest

// ExposesCredentialEdge builds the standard EXPOSES_CREDENTIAL edge from
// an AIService node to a Credential node. Looters that enumerate exposed
// upstream/provider secrets (litellmloot, openwebuiloot, ...) share this
// constructor so the edge's confidence, risk_weight, and evidence shape
// stay identical across collectors. SourceKind is AIService to satisfy
// the kinds registry's EXPOSES_CREDENTIAL constraint (source must be an
// AIService — see sdk/ingest/kinds.go).
func ExposesCredentialEdge(sourceID, credID, engagementID, source, endpoint string) Edge {
	return Edge{
		Source:     sourceID,
		Target:     credID,
		Kind:       "EXPOSES_CREDENTIAL",
		SourceKind: "AIService",
		TargetKind: "Credential",
		Properties: map[string]any{
			"confidence":  1.0,
			"risk_weight": 0.1,
			"evidence": map[string]any{
				"endpoint":      endpoint,
				"source":        source,
				"engagement_id": engagementID,
			},
		},
	}
}
