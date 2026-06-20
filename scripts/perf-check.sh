#!/usr/bin/env bash
# perf-check.sh — POISONS_CONTEXT pair-count invariant.
#
# B8 (POISONS_CONTEXT) widens the narrow SHADOWS guard: an injection-bearing
# tool can poison the shared agent context that drives a high-capability tool.
# To stop that breadth from exploding into a cartesian product, the shadows
# processor caps fan-out at 20 sinks per source tool. This script verifies the
# resulting per-agent ceiling holds against a live graph after a scan.
#
# The math: per-source cap = 20. An agent with N source tools (injection-
# bearing tools it can reach) tops out at N * 20 poisoned pairs. The ceiling
# below is 200, i.e. 10 source tools at the cap. An 11th source tool at the
# cap (220) breaches it — meaning either the per-source cap is too loose or
# the fixture is pathological. Tune MAX_PAIRS deliberately if the fleet shape
# legitimately changes.
#
# This script enforces the operator-facing runtime heuristic (per-agent <= 200
# pairs). The authoritative per-source cap (<= 20 sinks/source) regression gate
# is the Go integration test poisons_context_perf_integration_test.go.
#
# Requires a running graph (Neo4j + Postgres) reachable by agenthound-server.
# When no server binary or database is available the check SKIPS (exit 0) so
# it is safe to invoke in environments without a live stack.

set -euo pipefail

cd "$(dirname "$0")/.."

MAX_PAIRS="${PERF_CHECK_MAX_PAIRS:-200}"

# Resolve the server binary: explicit override, built artifact, or PATH.
BIN="${AGENTHOUND_SERVER_BIN:-}"
if [ -z "$BIN" ]; then
  if [ -x "./bin/agenthound-server" ]; then
    BIN="./bin/agenthound-server"
  elif command -v agenthound-server >/dev/null 2>&1; then
    BIN="agenthound-server"
  fi
fi

if [ -z "$BIN" ]; then
  echo "perf-check: agenthound-server binary not found (set AGENTHOUND_SERVER_BIN or run 'make build-server'); SKIPPING."
  exit 0
fi

QUERY='MATCH (a:AgentInstance)-[:TRUSTS_SERVER]->(:MCPServer)-[:PROVIDES_TOOL]->(t1)-[:POISONS_CONTEXT]->(t2)
RETURN a.name AS agent, count(DISTINCT [t1.objectid, t2.objectid]) AS poisoned_pairs
ORDER BY poisoned_pairs DESC LIMIT 1'

OUT=$(mktemp)
trap 'rm -f "$OUT"' EXIT

if ! "$BIN" query "$QUERY" --format json >"$OUT" 2>/dev/null; then
  echo "perf-check: could not query the graph (no database reachable?); SKIPPING."
  exit 0
fi

# No POISONS_CONTEXT edges -> printRows emits "(no results)"; nothing to gate.
if grep -q "(no results)" "$OUT"; then
  echo "perf-check: no POISONS_CONTEXT edges present; OK."
  exit 0
fi

# Extract the largest poisoned_pairs value from the JSON output.
MAX_OBSERVED=$(grep -oE '"poisoned_pairs"[[:space:]]*:[[:space:]]*[0-9]+' "$OUT" \
  | grep -oE '[0-9]+$' | sort -nr | head -1)
MAX_OBSERVED="${MAX_OBSERVED:-0}"

echo "perf-check: max poisoned_pairs per agent = $MAX_OBSERVED (ceiling $MAX_PAIRS)"

if [ "$MAX_OBSERVED" -gt "$MAX_PAIRS" ]; then
  echo "FAIL: an agent exceeds $MAX_PAIRS POISONS_CONTEXT pairs."
  echo "      Either the per-source fan-out cap (20) is too loose or the fixture is pathological."
  exit 1
fi

echo "perf-check: OK."
