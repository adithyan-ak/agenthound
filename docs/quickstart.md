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

## 3. Collect data

### Discover MCP client configurations

Scans your system for all known MCP client config files (Claude Desktop, Cursor, VS Code, Windsurf, Continue, Zed, Cline, JetBrains, Kiro, Amazon Q, Augment):

```bash
agenthound collect config --discover --output config-scan.json
```

### Enumerate MCP servers

Connects to each discovered MCP server and enumerates tools, resources, and prompts:

```bash
agenthound collect mcp --discover --output mcp-scan.json
```

### (Optional) Scan A2A agents

If you have A2A agents running:

```bash
agenthound collect a2a --target https://agent.example.com --output a2a-scan.json
```

## 4. Ingest into the graph

```bash
agenthound ingest config-scan.json
agenthound ingest mcp-scan.json
```

Each ingest builds graph nodes and edges in Neo4j, then runs post-processing to compute composite attack paths (CAN_REACH, CAN_EXFILTRATE_VIA, SHADOWS, etc.) and risk scores.

## 5. Open the UI

Navigate to [http://localhost:8080](http://localhost:8080) and log in:

- **Username:** `admin`
- **Password:** `agenthound`

> Change the default password in production by setting `AGENTHOUND_ADMIN_PASSWORD`.

## 6. Explore

- **Dashboard** -- overview of node/edge counts, risk distribution, top findings
- **Graph Explorer** -- interactive graph visualization, click nodes to inspect
- **Pathfinder** -- find shortest/weighted attack paths between any two nodes
- **Query Library** -- 17 pre-built security queries mapped to OWASP
- **Scan Manager** -- view scan history, trigger new scans

## 7. Run queries from the CLI

```bash
# List all findings
agenthound query --findings

# Critical findings only
agenthound query --findings --severity critical

# Pre-built query
agenthound query --prebuilt agents-shell-access

# Shortest path between two nodes
agenthound query --shortest-path \
  --from AgentInstance:claude-desktop \
  --to MCPResource:postgres://prod

# Raw Cypher
agenthound query "MATCH (a:AgentInstance)-[:TRUSTS_SERVER]->(s) RETURN a.name, s.name"
```

## Quick collect-and-ingest

Collectors support `--ingest` to skip the JSON file and write directly to the graph:

```bash
agenthound collect config --discover --ingest
agenthound collect mcp --discover --ingest
```

This requires a running Neo4j and PostgreSQL (via Docker or configured via environment variables).

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AGENTHOUND_NEO4J_URI` | `bolt://localhost:7687` | Neo4j connection URI |
| `AGENTHOUND_NEO4J_USER` | `neo4j` | Neo4j username |
| `AGENTHOUND_NEO4J_PASSWORD` | `agenthound` | Neo4j password |
| `AGENTHOUND_PG_URI` | `postgres://agenthound:agenthound@localhost:5432/agenthound?sslmode=disable` | PostgreSQL URI |
| `AGENTHOUND_API_PORT` | `8080` | API server port |
| `AGENTHOUND_LOG_LEVEL` | `info` | Log level: debug, info, warn, error |
| `AGENTHOUND_JWT_SECRET` | (random) | JWT signing secret (set for stable tokens) |
| `AGENTHOUND_ADMIN_PASSWORD` | `agenthound` | Initial admin password |

## Next steps

- [CLI Reference](cli-reference.md) -- full command documentation
- [API Reference](api-reference.md) -- REST API endpoints
- [Graph Model](graph-model.md) -- node types, edge types, risk scoring
- [Detection Rules](detection-rules.md) -- what AgentHound detects and why
