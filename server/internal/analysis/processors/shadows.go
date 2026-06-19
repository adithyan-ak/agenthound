package processors

import (
	"context"
	"time"

	"github.com/adithyan-ak/agenthound/server/internal/graph"
)

type Shadows struct{}

func (p *Shadows) Name() string           { return "shadows" }
func (p *Shadows) Dependencies() []string { return nil }

func (p *Shadows) Process(ctx context.Context, db graph.GraphDB, scanID string) (graph.ProcessingStats, error) {
	start := time.Now()

	// A SHADOWS edge means t1 specifically impersonates/overrides t2 by
	// naming it in t1's description. The match must be target-specific:
	// t1's description has to contain t2's name. We deliberately do NOT
	// branch on t1.has_cross_references — that flag is target-blind (it
	// is true if t1 references *any* sibling tool, see modules/mcp/
	// signals.go), so OR-ing it in made one flagged tool shadow every
	// tool on every other server, a cartesian blow-up of false positives.
	// has_cross_references still feeds tool risk scoring as a node
	// property (server/internal/analysis/riskscore/tool.go); it just no
	// longer manufactures SHADOWS edges.
	shadowsCypher := `
MATCH (s1:MCPServer)-[:PROVIDES_TOOL]->(t1:MCPTool),
      (s2:MCPServer)-[:PROVIDES_TOOL]->(t2:MCPTool)
WHERE s1 <> s2
  AND t1 <> t2
  AND t1.description IS NOT NULL
  AND t2.name IS NOT NULL
  AND toLower(t1.description) CONTAINS toLower(t2.name)
MERGE (t1)-[e:SHADOWS]->(t2)
ON CREATE SET e.confidence = CASE WHEN t1.has_injection_patterns = true THEN 0.9 ELSE 0.6 END,
              e.is_composite = true,
              e.source_collector = 'mcp',
              e.scan_id = $scan_id,
              e.risk_weight = 0.4,
              e.last_seen = datetime()
ON MATCH SET  e.scan_id = $scan_id,
              e.last_seen = datetime()
RETURN count(*) AS written`

	shadowsN, err := db.ExecuteWrite(ctx, shadowsCypher, map[string]any{"scan_id": scanID})
	if err != nil {
		return graph.ProcessingStats{
			ProcessorName: p.Name(),
			Duration:      time.Since(start),
		}, err
	}

	// POISONS_CONTEXT is the deliberate widening of the narrow SHADOWS
	// guard above: an injection-bearing tool can poison the shared agent
	// context that drives a high-capability tool, even without naming it.
	// To keep that breadth from exploding into a cartesian product, the
	// fan-out is capped at 20 sinks per source tool. With the perf-check
	// ceiling of 200 poisoned pairs per agent, an agent maxes out at 10
	// source tools * 20 sinks = 200 (see scripts/perf-check.sh).
	poisonsCypher := `
MATCH (src:MCPTool)
WHERE src.has_injection_patterns = true
MATCH (snk:MCPTool)
WHERE src <> snk
  AND any(cap IN snk.capability_surface WHERE cap IN ['shell_access', 'code_execution', 'credential_access', 'email_send'])
WITH src, collect(DISTINCT snk) AS sinks
WHERE size(sinks) <= 20
UNWIND sinks AS snk
MERGE (src)-[e:POISONS_CONTEXT]->(snk)
SET e.scan_id = $scan_id, e.last_seen = datetime(), e.is_composite = true,
    e.source_collector = 'mcp', e.confidence = 0.6, e.risk_weight = 0.4
RETURN count(*) AS written`

	poisonsN, err := db.ExecuteWrite(ctx, poisonsCypher, map[string]any{"scan_id": scanID})
	if err != nil {
		return graph.ProcessingStats{
			ProcessorName: p.Name(),
			EdgesCreated:  shadowsN,
			Duration:      time.Since(start),
		}, err
	}

	return graph.ProcessingStats{
		ProcessorName: p.Name(),
		EdgesCreated:  shadowsN + poisonsN,
		Duration:      time.Since(start),
	}, nil
}
