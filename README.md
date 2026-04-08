# AgentHound

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)](https://go.dev)
[![Neo4j](https://img.shields.io/badge/Neo4j-4.4+-008CC1.svg)](https://neo4j.com)

**BloodHound for AI Agent Infrastructure.**

AgentHound enumerates MCP servers and A2A agents, builds a directed trust graph in Neo4j, and uses shortest-path algorithms to discover attack paths across protocol boundaries. It is the first tool to perform cross-protocol graph-based attack path analysis for AI agent infrastructure.

## Key features

- **Three collectors** -- Config (12 MCP client formats), MCP (tool/resource/prompt enumeration via official Go SDK), A2A (Agent Card fetching with JWS verification)
- **Directed trust graph** -- 14 node types, 23 edge types (13 direct + 10 computed), deterministic SHA-256 node IDs for cross-collector merge
- **Attack path discovery** -- shortest path, all paths, Dijkstra weighted path across agent-server-tool-resource chains
- **Cross-protocol analysis** -- finds A2A-to-MCP attack paths via host correlation that no single-protocol tool can detect
- **17 pre-built security queries** -- mapped to OWASP MCP Top 10 and OWASP Agentic Top 10
- **Risk scoring** -- weighted scores (0-100) for agents, servers, and tools based on auth posture, capabilities, exposure, and poisoning
- **Detections** -- tool poisoning, tool shadowing, rug pulls, unauthenticated servers, credential exposure, instruction poisoning, supply chain risks, data exfiltration routes
- **Interactive UI** -- Sigma.js graph explorer (100K+ nodes), pathfinder, entity inspector, dashboard, query library

## Quickstart

```bash
git clone https://github.com/adithyan-ak/agenthound.git
cd agenthound

# Start Neo4j + PostgreSQL + AgentHound
docker compose -f docker/docker-compose.yml up -d

# Discover and scan MCP infrastructure
agenthound collect config --discover --ingest
agenthound collect mcp --discover --ingest

# Open the UI
open http://localhost:8080    # login: admin / agenthound

# Or query from the CLI
agenthound query --findings --severity critical
```

See the [full quickstart guide](docs/quickstart.md) for prerequisites and detailed steps.

## Architecture

```
+------------------+     +------------------+     +------------------+
|  Config Collector|     |   MCP Collector  |     |   A2A Collector  |
|  (12 parsers)    |     |  (Go SDK v1.5.0) |     |  (HTTP + JWS)   |
+--------+---------+     +--------+---------+     +--------+---------+
         |                        |                        |
         v                        v                        v
    +----+------------------------+------------------------+----+
    |                    Ingest Pipeline                         |
    |  validate -> normalize -> deduplicate -> write -> post-   |
    |                                                  process  |
    +----------------------------+------------------------------+
                                 |
                    +------------+------------+
                    |                         |
               +----+----+            +------+------+
               | Neo4j   |            | PostgreSQL  |
               | (graph) |            | (app data)  |
               +---------+            +-------------+
                    |
         +----------+----------+
         |                     |
    +----+----+         +------+------+
    | REST API|         |   React UI  |
    | (chi)   |         | (Sigma.js)  |
    +---------+         +-------------+
```

**Three collectors** produce standardized JSON. The **ingest pipeline** validates, normalizes, and writes to Neo4j, then runs **9 post-processors** to compute composite attack paths and risk scores. The **REST API** serves the embedded **React SPA** and exposes all graph, analysis, and management endpoints.

## Documentation

| Document | Description |
|----------|-------------|
| [Quickstart](docs/quickstart.md) | 5-minute setup guide |
| [CLI Reference](docs/cli-reference.md) | All commands and flags |
| [API Reference](docs/api-reference.md) | REST API endpoints |
| [Graph Model](docs/graph-model.md) | Node types, edge types, risk scoring |
| [Detection Rules](docs/detection-rules.md) | What AgentHound detects, mapped to OWASP |
| [Contributing](CONTRIBUTING.md) | How to add collectors, detections, queries |
| [Changelog](CHANGELOG.md) | Release notes |
| [Security](SECURITY.md) | Vulnerability reporting |

## Tech stack

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.25+, Cobra CLI, chi/v5 router |
| Graph DB | Neo4j 4.4+ Community (Cypher + APOC) |
| App DB | PostgreSQL 16 (pgx/v5) |
| MCP SDK | `github.com/modelcontextprotocol/go-sdk` v1.5.0 |
| Frontend | React 18, TypeScript, Vite 6, Sigma.js 3, shadcn/ui |
| Auth | bcrypt, JWT (HMAC-SHA256), API tokens, RBAC |
| Deployment | Docker Compose (Neo4j + PostgreSQL + AgentHound) |

## License

[Apache License 2.0](LICENSE)
