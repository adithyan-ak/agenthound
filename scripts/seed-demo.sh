#!/usr/bin/env bash
# AgentHound v0.2 demo seed.
#
# Brings up the docker-compose lab, runs the collector against the lab
# CIDR + the operator's MCP config, runs the LiteLLM looter, and pipes
# all three ingest envelopes into the analysis server. After this
# completes, the Findings panel in the UI surfaces the
# cross_service_credential_chain finding tying the agent's env-var
# credential to the LiteLLM-extracted upstream provider keys.
#
# Prerequisites:
#   - Docker Compose installed
#   - agenthound-server running on http://localhost:8080
#     (start with `docker compose -f docker/docker-compose.yml up -d`)
#   - bin/agenthound and bin/agenthound-server built (`make build`)
#
# Idempotent: re-running this script tears down the lab cleanly and
# starts fresh. Use --keep to skip teardown.
#
# v0.1 LEGACY: previous versions of this script ingested pre-baked
# testdata/demo/{config,mcp,a2a}_scan.json fixtures. v0.2 generates
# real ingest envelopes from a live lab so the credential-chain demo
# uses authentic data — see docs/plans/v0.2-implementation.md
# Phase 7 acceptance criteria.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LAB_COMPOSE="$REPO_ROOT/docker/demo/docker-compose.yml"
LAB_CIDR="172.20.0.0/24"
LITELLM_HOST="172.20.0.20:4000"
MASTER_KEY="sk-DEMO-CHAIN-KEY-NOT-REAL"
OPERATOR_CONFIG="$REPO_ROOT/docker/demo/operator-config/claude_desktop_config.json"
ENGAGEMENT_ID="V02-DEMO-LOCAL"

CMD_AGENTHOUND="$REPO_ROOT/bin/agenthound"
CMD_SERVER="$REPO_ROOT/bin/agenthound-server"

if [[ ! -x "$CMD_AGENTHOUND" || ! -x "$CMD_SERVER" ]]; then
    echo "[seed-demo] error: bin/agenthound and bin/agenthound-server must exist."
    echo "[seed-demo]        run: cd $REPO_ROOT && make build"
    exit 1
fi

KEEP=0
if [[ "${1:-}" == "--keep" ]]; then
    KEEP=1
    shift
fi

mkdir -p "$REPO_ROOT/testdata/demo"

echo "[seed-demo] tearing down any prior lab state..."
docker compose -f "$LAB_COMPOSE" down -v --remove-orphans >/dev/null 2>&1 || true

echo "[seed-demo] starting demo lab (Ollama + LiteLLM stub + operator)..."
docker compose -f "$LAB_COMPOSE" up -d --build

echo "[seed-demo] waiting for litellm stub to come up..."
for i in $(seq 1 30); do
    if docker compose -f "$LAB_COMPOSE" exec -T litellm-stub curl -fsS http://127.0.0.1:4000/health/liveliness >/dev/null 2>&1; then
        echo "[seed-demo] litellm stub healthy."
        break
    fi
    sleep 1
done

# 1. Network scan against the lab CIDR. Phase 2/3 fingerprinters fire
#    on the open ports for both Ollama (11434) and LiteLLM (4000) and
#    emit multi-label :OllamaInstance:AIService and
#    :LiteLLMGateway:AIService nodes.
SCAN_OUT="$REPO_ROOT/testdata/demo/scan-output.json"
echo "[seed-demo] running network scan against $LAB_CIDR..."
"$CMD_AGENTHOUND" scan "$LAB_CIDR" --output "$SCAN_OUT" --network-scan-concurrency 10 \
    --timeout 5s

# 2. Config collector against the operator's MCP config. Emits the
#    :Credential node with value_hash matching the master Credential
#    the looter is about to emit — the load-bearing cross-collector
#    merge.
CONFIG_OUT="$REPO_ROOT/testdata/demo/config-output.json"
echo "[seed-demo] running config collector against operator config..."
"$CMD_AGENTHOUND" scan --config --path "$OPERATOR_CONFIG" --output "$CONFIG_OUT"

# 3. LiteLLM Looter against the stub gateway. Use --master-key sugar
#    rather than --credential KEY=VALUE because that's the documented
#    operator workflow.
LOOT_OUT="$REPO_ROOT/testdata/demo/loot-output.json"
echo "[seed-demo] running litellm looter against $LITELLM_HOST..."
# AUTHORIZED prompt sentinel — first invocation prompts; subsequent
# invocations skip. Pipe AUTHORIZED unconditionally so seed-demo.sh
# is idempotent on a fresh machine.
echo "AUTHORIZED" | "$CMD_AGENTHOUND" loot "$LITELLM_HOST" \
    --type litellm \
    --master-key "$MASTER_KEY" \
    --engagement-id "$ENGAGEMENT_ID" \
    --output "$LOOT_OUT"

# 4. Ingest all three envelopes into the running analysis server.
echo "[seed-demo] ingesting scan envelope..."
"$CMD_SERVER" ingest "$SCAN_OUT"
echo "[seed-demo] ingesting config envelope..."
"$CMD_SERVER" ingest "$CONFIG_OUT"
echo "[seed-demo] ingesting loot envelope..."
"$CMD_SERVER" ingest "$LOOT_OUT"

echo
echo "[seed-demo] done. Open http://localhost:8080 → Findings →"
echo "[seed-demo]   look for 'litellm-credential-leak' in Critical Paths."
echo
echo "[seed-demo] credential-chain check (raw API):"
echo "[seed-demo]   curl -s localhost:8080/api/v1/analysis/findings | \\"
echo "[seed-demo]     jq '.[] | select(.processor==\"cross_service_credential_chain\")'"
echo

if [[ "$KEEP" -eq 0 ]]; then
    echo "[seed-demo] (lab containers still running; 'docker compose -f $LAB_COMPOSE down' to stop)"
fi
