# Changelog

## v0.1.0 (2026-04-08)

Initial release. AgentHound is a BloodHound-style security tool for AI agent infrastructure.

### Collectors

- **Config Collector** -- parses 12 MCP client configuration formats (Claude Desktop, Claude Code, Cursor, VS Code, Windsurf, Continue, Zed, Cline, JetBrains, Kiro, Amazon Q, Augment) to discover agent-server trust relationships, credentials, instruction files, and host information
- **MCP Collector** -- connects to MCP servers via stdio and Streamable HTTP transports using the official Go SDK (v1.5.0), enumerates tools/resources/prompts, classifies capability surfaces, detects injection patterns and cross-references
- **A2A Collector** -- fetches A2A Agent Cards (v0.3.0 and v1.0) via HTTP, verifies JWS signatures (RFC 7515), scores auth posture, supports domain discovery (`--discover-domain`)

### Graph engine

- Neo4j 4.4+ with auto-detected schema syntax (4.4 `ON...ASSERT` / 5.x `FOR...REQUIRE`)
- 14 node types, 13 direct edge types, 10 composite edge types
- Deterministic content-based SHA-256 node IDs for cross-collector merge
- Batch MERGE writes with UNWIND (1000 ops/txn)
- APOC Dijkstra with non-APOC fallback for weighted pathfinding

### Ingest pipeline

- JSON schema validation, camelCase-to-snake_case normalization, deduplication by objectid
- Rug pull detection via description_hash tracking across scans
- Stale composite edge cleanup scoped to source collector

### Post-processors (9)

- HAS_ACCESS_TO -- capability surface to URI scheme matching
- CAN_EXECUTE -- shell/code execution tool to host mapping
- SHADOWS -- cross-server tool shadowing detection
- POISONED_DESCRIPTION -- injection pattern detection in tool descriptions
- CAN_REACH -- transitive agent-to-resource access paths
- CAN_EXFILTRATE_VIA -- sensitive data access + outbound channel combination
- CAN_IMPERSONATE -- TF-IDF cosine similarity on A2A skill descriptions
- POISONED_INSTRUCTIONS -- suspicious pattern detection in instruction files
- Cross-protocol CAN_REACH -- A2A-to-MCP boundary traversal via host correlation
- Risk scoring -- weighted scores for agents, servers, and tools (0-100)

### 17 pre-built queries

Mapped to OWASP MCP Top 10 and OWASP Agentic Top 10:
- Critical Paths: agents-shell-access, shortest-to-database, cross-protocol-paths, exfiltration-routes, credential-chain
- Vulnerabilities: poisoned-tools, tool-shadowing, no-auth-servers, no-auth-a2a, rug-pull
- Supply Chain: unpinned-packages, instruction-poisoning, unsigned-cards, high-entropy-secrets
- Chokepoints: chokepoint-servers, chokepoint-tools
- Combined: unpinned-shell

### REST API

- Full CRUD for graph nodes/edges
- Pathfinding: shortest-path, all-paths, weighted-path (Dijkstra)
- Findings endpoint with severity filtering
- Pre-built query execution
- Scan history and management
- Rate limiting (100/min general, 20/min ingest, 10/min raw Cypher)

### Authentication and authorization

- bcrypt password hashing, JWT tokens (24h, HMAC-SHA256)
- API tokens with `ah_` prefix for programmatic access
- RBAC: admin, analyst, viewer roles
- First-run admin user creation
- Audit logging for security-relevant actions

### Frontend

- React 18 + TypeScript + Vite 6 SPA embedded in Go binary via `go:embed`
- Graph Explorer with Sigma.js 3 (WebGL, 100K+ nodes) and ForceAtlas2 layout
- Interactive Pathfinder with shortest/weighted path visualization
- Entity Inspector with properties, connections, risk breakdown
- Dashboard with stats, risk distribution, auth coverage, top findings
- Query Library with all 17 pre-built queries
- Scan Manager for history and triggering scans
- Login flow with JWT session management

### Infrastructure

- Docker Compose: Neo4j 4.4 + PostgreSQL 16 + AgentHound
- Multi-stage Dockerfile (golang:1.25-alpine build, alpine:3.19 runtime)
- Non-root container user (UID 1001)
- Makefile: build, test, lint, docker, ui-build, seed
- CI: golangci-lint, tests with Neo4j+PG services, build verification
