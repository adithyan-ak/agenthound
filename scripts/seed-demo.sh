#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "Seeding AgentHound demo data..."
echo ""

echo "[1/3] Ingesting config collector scan..."
agenthound ingest "$PROJECT_DIR/testdata/demo/config_scan.json"

echo "[2/3] Ingesting MCP collector scan..."
agenthound ingest "$PROJECT_DIR/testdata/demo/mcp_scan.json"

echo "[3/3] Ingesting A2A collector scan..."
agenthound ingest "$PROJECT_DIR/testdata/demo/a2a_scan.json"

echo ""
echo "Demo data loaded successfully."
echo ""
echo "Open http://localhost:8080 to explore the graph."
echo ""
echo "The demo environment contains:"
echo "  - 2 agent instances (Claude Desktop, Cursor)"
echo "  - 6 MCP servers (filesystem, postgres-prod, slack, github, execute-command, s3-backup)"
echo "  - 3 A2A agents (data-pipeline, external-assistant, admin-bot)"
echo "  - 25 tools, 8 resources, 8 skills"
echo "  - 1 poisoned tool (prompt injection in description)"
echo "  - 1 tool shadowing another server's tool"
echo "  - 1 unsigned A2A agent with no authentication"
echo "  - 1 hardcoded high-entropy AWS secret key"
echo "  - 1 suspicious instruction file"
echo "  - Cross-protocol attack paths (A2A -> MCP via shared localhost)"
echo ""
echo "Try these queries:"
echo "  agenthound query --prebuilt agents-shell-access"
echo "  agenthound query --prebuilt cross-protocol-paths"
echo "  agenthound query --prebuilt exfiltration-routes"
echo "  agenthound query --prebuilt poisoned-tools"
echo "  agenthound query --prebuilt no-auth-a2a"
echo "  agenthound query --prebuilt high-entropy-secrets"
