# Quickstart

Get AgentHound running and scanning your AI agent infrastructure in 5 minutes.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go 1.25+](https://go.dev/dl/) (for building from source)
- At least one MCP client configured (Claude Desktop, Cursor, VS Code, etc.)

## 1. Start the infrastructure

```bash
git clone https://github.com/adithyan-ak/agenthound.git
cd agenthound
docker compose -f docker/docker-compose.yml up -d
```

This starts Neo4j (graph database), PostgreSQL (app database), and the AgentHound server on port 8080.

Wait for all services to be healthy:

```bash
docker compose -f docker/docker-compose.yml ps
```

## 2. Build from source (optional)

If you want to run collectors locally instead of through Docker:

```bash
make build
```

The binary is at `bin/agenthound`.

## 3. Scan your infrastructure

Run a full scan that discovers configs and enumerates MCP servers. The collector writes JSON to a local file (auto-named `./scan-<scan_id>.json` in the current directory):

```bash
agenthound scan
```

This auto-discovers all MCP client config files (Claude Desktop, Cursor, VS Code, Windsurf, Continue, Zed, Cline, JetBrains, Kiro, Amazon Q, Augment) and connects to each configured MCP server.

### (Optional) Scan A2A agents

If you have A2A agents running:

```bash
agenthound scan --a2a --target https://agent.example.com
```

## 4. Ingest the scan into the graph

The collector is offline-by-default. Move the scan JSON to the operator's box and ingest it. There are three equivalent paths:

```bash
# (a) File + CLI ingest
agenthound-server ingest scan-*.json

# (b) Stream over stdin (no file at rest)
agenthound scan --output - | agenthound-server ingest -

# (c) UI drag-drop import
# Open http://localhost:8080 → Scan Manager → "Import scan" → drag scan.json into the dropzone
```

The same ingest pipeline (validate → normalize → deduplicate → write → post-process) runs for all three. After ingest, the post-processors compute composite attack paths and risk scores.

## 5. Open the UI

Navigate to [http://localhost:8080](http://localhost:8080).

The server has **no application-layer authentication**. It binds to `127.0.0.1:8080` by default; access from another machine should go over a network you trust (VPN, SSH tunnel, Tailscale). See [`security.md`](security.md) for the full posture.

## 6. Explore

- **Dashboard** -- overview of node/edge counts, risk distribution, top findings
- **Graph Explorer** -- interactive graph visualization, click nodes to inspect
- **Findings** -- per-finding detail page with the embedded attack-path diagram
- **Query Library** -- 17 pre-built security queries mapped to OWASP
- **Scan Manager** -- view scan history, trigger new scans

## 7. Run queries from the CLI (operator's box)

```bash
# List all findings
agenthound-server query --findings

# Critical findings only
agenthound-server query --findings --severity critical

# Pre-built query
agenthound-server query --prebuilt agents-shell-access

# Shortest path between two nodes
agenthound-server query --shortest-path \
  --from AgentInstance:claude-desktop \
  --to MCPResource:postgres://prod

# Raw Cypher
agenthound-server query "MATCH (a:AgentInstance)-[:TRUSTS_SERVER]->(s) RETURN a.name, s.name"
```

## Quick scan tips

Run individual collectors when you only need partial data:

```bash
agenthound scan --config                           # Config files only (offline, no network)
agenthound scan --mcp --url https://mcp.example.com  # Single MCP server
```

For CI/CD pipelines, gate on findings on the operator's server (not on the collector):

```bash
agenthound-server query --findings --severity critical --fail-on critical
```

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AGENTHOUND_NEO4J_URI` | `bolt://localhost:7687` | Neo4j connection URI |
| `AGENTHOUND_NEO4J_USER` | `neo4j` | Neo4j username |
| `AGENTHOUND_NEO4J_PASSWORD` | `agenthound` | Neo4j password |
| `AGENTHOUND_PG_URI` | `postgres://agenthound:agenthound@localhost:5432/agenthound?sslmode=disable` | PostgreSQL URI |
| `AGENTHOUND_BIND` | `127.0.0.1:8080` | Server bind address (`host:port`). Set to `0.0.0.0:8080` only inside a trusted network. |
| `AGENTHOUND_LOG_LEVEL` | `info` | Log level: debug, info, warn, error |
| `AGENTHOUND_CORS_ORIGINS` | `http://localhost:8080` | Comma-separated CORS origins for the UI |

## Next steps

- [CLI Reference](cli-reference.md) -- full command documentation
- [API Reference](api-reference.md) -- REST API endpoints
- [Graph Model](graph-model.md) -- node types, edge types, risk scoring
- [Detection Rules](detection-rules.md) -- what AgentHound detects and why
