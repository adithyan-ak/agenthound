package processors

import (
	"context"
	"time"

	"github.com/adithyan-ak/agenthound/server/internal/graph"
)

// CrossServiceCredentialChain wires the v0.2 credential-chain demo:
// when the Config Collector emits a Credential C1 (an env var or
// hardcoded value referenced by an MCP server) AND the LiteLLM Looter
// emits a Credential C1master (the master key the operator supplied),
// AND C1.value_hash == C1master.value_hash, those two nodes describe
// the same secret. The LiteLLM gateway then exposes upstream provider
// keys (C2) via EXPOSES_CREDENTIAL — and a CAN_REACH edge from the
// agent through that chain to the upstream credential is the finding.
//
// Path:
//
//	(:AgentInstance)-[:TRUSTS_SERVER]->(:MCPServer)
//	    -[:HAS_ENV_VAR]->(:Credential C1)
//	    where C1.value_hash matches a Credential C1master from a LiteLLM Looter run
//	(:LiteLLMGateway:AIService gw)-[:EXPOSES_CREDENTIAL]->(:Credential C1master)
//	(gw)-[:EXPOSES_CREDENTIAL]->(:Credential C2)
//	    where C2.type IN ["apiKey", "virtual_key"]
//
// We emit (:AgentInstance)-[:CAN_REACH]->(C2) with metadata that
// records the merge endpoint (hash + LiteLLM gateway name).
//
// Dependencies: ["has_access_to", "can_reach"] — has_access_to so the
// graph has resource accessibility wired, can_reach so this processor
// runs AFTER the existing transitive can_reach work and we don't
// double-emit edges that the Phase 4 chain already covers.
type CrossServiceCredentialChain struct{}

func (p *CrossServiceCredentialChain) Name() string { return "cross_service_credential_chain" }

func (p *CrossServiceCredentialChain) Dependencies() []string {
	return []string{"has_access_to", "can_reach"}
}

func (p *CrossServiceCredentialChain) Process(ctx context.Context, db graph.GraphDB, scanID string) (graph.ProcessingStats, error) {
	start := time.Now()

	// The join: c1.value_hash = c1master.value_hash. c1 comes from the
	// Config Collector (MCPServer-[:HAS_ENV_VAR]->c1). c1master comes
	// from the LiteLLM Looter ((gw:LiteLLMGateway)-[:EXPOSES_CREDENTIAL]->c1master).
	// gw also -[:EXPOSES_CREDENTIAL]->c2, the upstream provider Credential.
	//
	// We require c1 != c1master (otherwise this would only fire on
	// hand-loaded test fixtures where both nodes happen to share an
	// objectid; in real graphs they always have different objectids
	// because the Config Collector and Looter compute IDs differently).
	cypher := `
MATCH (a:AgentInstance)-[:TRUSTS_SERVER]->(s:MCPServer)
      -[:HAS_ENV_VAR]->(c1:Credential)
WHERE c1.value_hash IS NOT NULL AND c1.value_hash <> ''
MATCH (gw:LiteLLMGateway)-[:EXPOSES_CREDENTIAL]->(c1master:Credential)
WHERE c1master.value_hash = c1.value_hash AND c1master.objectid <> c1.objectid
MATCH (gw)-[:EXPOSES_CREDENTIAL]->(c2:Credential)
WHERE c2.type IN ['apiKey', 'virtual_key'] AND c2.objectid <> c1master.objectid
MERGE (a)-[e:CAN_REACH]->(c2)
SET e.scan_id = $scan_id, e.last_seen = datetime(), e.is_composite = true,
    e.source_collector = 'cross_service_credential_chain',
    e.via_server = s.name,
    e.via_credential = c1.name,
    e.via_gateway = gw.name,
    e.merge_value_hash = c1.value_hash,
    e.upstream_provider = COALESCE(c2.provider, 'unknown'),
    e.hops = 5,
    e.confidence = 0.95,
    e.risk_weight = 0.1
RETURN count(e) AS written`

	written, err := db.ExecuteWrite(ctx, cypher, map[string]any{"scan_id": scanID})
	if err != nil {
		return graph.ProcessingStats{
			ProcessorName: p.Name(),
			Duration:      time.Since(start),
		}, err
	}
	return graph.ProcessingStats{
		ProcessorName: p.Name(),
		EdgesCreated:  written,
		Duration:      time.Since(start),
	}, nil
}
