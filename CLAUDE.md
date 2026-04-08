# AgentHound — Project Context for Implementation

AgentHound is a BloodHound-style open-source security tool that enumerates MCP servers and A2A agents, builds a directed trust graph in Neo4j, and uses shortest-path algorithms to discover attack paths across protocol boundaries. **No existing tool does cross-protocol graph-based attack path analysis for AI agent infrastructure.**

## Pre-Commit Checks (MANDATORY)

Before every commit, run these checks locally and fix all issues:

```bash
gofmt -l .                  # Must produce no output (no formatting issues)
go build ./...              # Must pass with zero errors
go vet ./...                # Must pass with zero warnings
```

If any check fails, fix the issues before committing. Do NOT commit code that fails these checks. The CI runs `golangci-lint` which enforces `errcheck` (all error return values must be handled — use `_, _ =` for intentionally discarded ones like `fmt.Fprintf` to stderr) and `gofmt` (no manual alignment padding).

## Tech Stack

| Component | Choice | Key Details |
|-----------|--------|-------------|
| Backend | **Go 1.25+** | Required by MCP Go SDK v1.5.0 (`iter.Seq2`). Single binary. |
| CLI | **cobra** | `agenthound serve`, `collect mcp/a2a/config`, `ingest`, `query` |
| HTTP router | **chi/v5** | REST API at `/api/v1/*` |
| Graph DB | **Neo4j 4.4+ Community** | Cypher pathfinding, APOC for Dijkstra. Dual syntax: 4.4 uses `ON...ASSERT`, 5.x uses `FOR...REQUIRE`. |
| App DB | **PostgreSQL 16** | Users, scans, audit_log, api_tokens tables. Driver: `pgx/v5`. |
| MCP SDK | `github.com/modelcontextprotocol/go-sdk` **v1.5.0** | `mcp.NewClient()`, `mcp.CommandTransport`, auto-paginating iterators (`session.Tools(ctx, nil)`) |
| Neo4j driver | `neo4j-go-driver/v5` v5.28+ | v5 for 4.4 compat. Batch writes with `UNWIND` + `MERGE`. |
| JSON validation | `santhosh-tekuri/jsonschema/v6` | Ingest schema validation |
| Frontend | **React 18 + TypeScript + Vite 6** | SPA embedded in Go binary via `go:embed` |
| Graph viz | **sigma 3.0.2 + graphology 0.26.0 + @react-sigma/core 5.0.6** | WebGL, 100K+ nodes. Pass `MultiDirectedGraph` constructor to `SigmaContainer`, NOT an instance. |
| Layout | **graphology-layout-forceatlas2 0.10.1** | Force-directed. Run 2s then freeze. |
| UI | **shadcn/ui + Tailwind CSS 3.x + Zustand 5.x + TanStack Query 5.x + Recharts 2.x** | |
| Deployment | **Docker Compose** | 3 containers: `graph-db` (neo4j:4.4-community), `app-db` (postgres:16-alpine), `agenthound` |
| License | Apache 2.0 | |

## Project Structure

```
agenthound/
├── cmd/agenthound/main.go              # Cobra entry point
├── internal/
│   ├── api/                             # REST API server (chi router)
│   │   ├── server.go                    # Routes, go:embed for UI
│   │   ├── middleware/                   # logging, cors
│   │   ├── handlers/                    # health, graph, ingest, analysis, scans
│   │   └── ui/dist/                     # Copied from ui/dist by Makefile (go:embed target)
│   ├── cli/                             # Cobra subcommands: root, serve, ingest, collect_*, query
│   ├── graph/                           # Neo4j layer
│   │   ├── driver.go                    # Connection management
│   │   ├── schema.go                    # Constraints + indexes, version-detected syntax
│   │   ├── writer.go                    # Batch MERGE writer (UNWIND). APOC preferred, per-kind fallback.
│   │   └── reader.go                    # Query execution, GetNode, GetStats
│   ├── ingest/                          # Pipeline: validate → normalize → deduplicate → write
│   │   ├── pipeline.go                  # Orchestrator, triggers post-processing after write
│   │   ├── validator.go                 # Schema validation (meta, node kinds, edge kinds)
│   │   └── normalizer.go               # camelCase→snake_case, timestamps, objectid sync
│   ├── collector/
│   │   ├── common/                      # hasher.go, patterns.go, capability.go, entropy.go
│   │   ├── config/                      # Config Collector + parsers/ (12 client formats)
│   │   ├── mcp/                         # MCP Collector (wraps Go SDK)
│   │   └── a2a/                         # A2A Collector (HTTP GET + JSON parse)
│   ├── analysis/                        # Post-processing
│   │   ├── postprocessor.go             # Runner with dependency ordering
│   │   ├── processors/                  # has_access_to, can_execute, shadows, poisoned, can_reach,
│   │   │                                # can_exfiltrate, can_impersonate, cross_protocol, risk_score
│   │   ├── riskscore/                   # agent.go, server.go, tool.go, weights.go
│   │   ├── similarity/tfidf.go          # For CAN_IMPERSONATE cosine similarity
│   │   └── sensitivity/classifier.go    # Resource URI → sensitivity level
│   ├── model/                           # Node, Edge, IngestData, IngestMeta, GraphData, Scan
│   ├── appdb/                           # PostgreSQL: driver, migrations, scans CRUD
│   ├── auth/                            # Phase 5: bcrypt passwords, JWT (24h, HMAC-SHA256), API tokens (ah_ prefix), RBAC
│   ├── audit/                           # Phase 5: audit_log writes
│   └── config/config.go                 # Env-based config (AGENTHOUND_NEO4J_URI, etc.)
├── pkg/
│   ├── collector/collector.go           # Collector interface
│   └── analysis/postprocessor.go        # PostProcessor interface
├── ui/                                  # React SPA (Vite)
│   └── src/
│       ├── api/                         # client.ts (ky), graph.ts, analysis.ts, scans.ts
│       ├── store/                       # Zustand: graph.ts, search.ts, ui.ts
│       ├── components/
│       │   ├── graph/                   # GraphExplorer, Controls, Search, Filters, Legend
│       │   ├── pathfinder/              # Pathfinder, PathSelector, PathResults, PathHighlight
│       │   ├── dashboard/               # Dashboard, StatCards, RiskChart, AuthCoverage, TopFindings
│       │   ├── inspector/               # EntityInspector, NodeProperties, Connections, RiskBreakdown
│       │   ├── scans/                   # ScanManager, NewScan, ScanHistory
│       │   └── queries/                 # QueryLibrary (17 pre-built), QueryResult
│       ├── hooks/                       # useGraph, usePathfinding, useNodeSearch, useGraphEvents
│       └── lib/                         # graph-builder.ts, node-styles.ts, edge-styles.ts, layout.ts
├── docker/
│   ├── Dockerfile                       # Multi-stage: golang:1.25-alpine → alpine:3.19
│   └── docker-compose.yml               # neo4j:4.4 + postgres:16 + agenthound
├── testdata/                            # valid_mcp_scan.json, valid_a2a_scan.json, valid_config_scan.json, merged_scan.json, invalid_scan.json
├── scripts/                             # seed-test-data.sh, seed-demo.sh
├── Makefile                             # build, test, lint, docker, up, down, ui-build
├── .github/workflows/ci.yml            # lint, test (neo4j+pg services), build, docker
└── .golangci.yml
```

## Graph Data Model

**Core principle:** Edges = exploitable relationships. Direction = access flow. `Agent → Server → Tool → Resource`.

### Node Types (12 collector-produced + 2 synthetic)

| Label | Source | Key Properties |
|-------|--------|----------------|
| `MCPServer` | Config + MCP | `name`, `endpoint`, `transport` (stdio/http), `auth_method`, `protocol_version`, `instructions`, `capabilities`, `is_pinned`, `has_tasks_capability` |
| `MCPTool` | MCP | `name`, `description`, `input_schema`, `output_schema`, `annotations`, `description_hash` (SHA-256), `capability_surface[]`, `has_injection_patterns`, `has_cross_references` |
| `MCPResource` | MCP | `uri`, `name`, `mime_type`, `size`, `uri_scheme`, `sensitivity` (auto-classified) |
| `MCPPrompt` | MCP | `name`, `description`, `arguments` |
| `A2AAgent` | A2A | `name`, `description`, `url`, `provider`, `version`, `protocol_versions`, `capabilities`, `security_schemes`, `auth_method`, `is_signed`, `signature_valid`, `card_hash` |
| `A2ASkill` | A2A | `id`, `name`, `description`, `input_modes`, `output_modes`, `description_hash`, `has_injection_patterns` |
| `AgentInstance` | Config | `name`, `framework`, `config_path` |
| `Identity` | Config + MCP | `type` (none/apiKey/oauth/bearer/mtls), `scope`, `is_static` |
| `Credential` | Config | `type` (envVar/hardcoded/vaultRef/inputPrompt), `name`, `source`, `is_exposed`, `high_entropy` |
| `Host` | Config + A2A | `hostname`, `ip`, `is_local`, `is_private`, `is_public` |
| `ConfigFile` | Config | `path`, `client`, `server_count` |
| `InstructionFile` | Config | `path`, `type` (agents.md/claude.md/cursorrules/copilot-instructions/memory.md), `hash`, `is_suspicious` |
| `ResourceGroup` | Post-processor (synthetic) | `type`, `sensitivity` |
| `TrustZone` | Post-processor (synthetic) | `name`, `level`, `node_count` |

### Node ID Strategy (deterministic, content-based SHA-256)

```
MCPServer:       SHA-256("MCPServer:" + transport + ":" + endpoint + ":" + sorted_args)
MCPTool:         SHA-256("MCPTool:" + server_id + ":" + tool_name)
MCPResource:     SHA-256("MCPResource:" + server_id + ":" + resource_uri)
A2AAgent:        SHA-256("A2AAgent:" + agent_card_url)
A2ASkill:        SHA-256("A2ASkill:" + agent_id + ":" + skill_id)
AgentInstance:   SHA-256("AgentInstance:" + config_file_id + ":" + client_name)
ConfigFile:      SHA-256("ConfigFile:" + absolute_path)
Host:            SHA-256("Host:" + hostname_or_ip)
```

**CRITICAL:** MCPServer ID MUST match between Config Collector and MCP Collector — this is the merge point that connects trust (who trusts what) to capabilities (what it exposes).

### Edge Types

**Directly collected (from collectors):**

| Edge | Source → Target | Collector |
|------|----------------|-----------|
| `TRUSTS_SERVER` | AgentInstance → MCPServer | Config |
| `PROVIDES_TOOL` | MCPServer → MCPTool | MCP |
| `PROVIDES_RESOURCE` | MCPServer → MCPResource | MCP |
| `PROVIDES_PROMPT` | MCPServer → MCPPrompt | MCP |
| `ADVERTISES_SKILL` | A2AAgent → A2ASkill | A2A |
| `DELEGATES_TO` | A2AAgent → A2AAgent | A2A |
| `AUTHENTICATES_WITH` | MCPServer/A2AAgent → Identity | Config/A2A |
| `USES_CREDENTIAL` | Identity → Credential | Config |
| `RUNS_ON` | MCPServer/A2AAgent → Host | Config/A2A |
| `CONFIGURED_IN` | MCPServer → ConfigFile | Config |
| `HAS_ENV_VAR` | MCPServer → Credential | Config |
| `LOADS_INSTRUCTIONS` | AgentInstance → InstructionFile | Config |
| `SAME_AUTH_DOMAIN` | A2AAgent → A2AAgent | A2A |

**Post-processed composite edges (computed from graph state):**

| Edge | Source → Target | Depends On | Key Detail |
|------|----------------|------------|------------|
| `HAS_ACCESS_TO` | MCPTool → MCPResource | Raw edges | Capability surface + URI scheme match |
| `CAN_EXECUTE` | MCPTool → Host | Raw edges | shell_access/code_execution capability |
| `SHADOWS` | MCPTool → MCPTool | Raw edges | Cross-server description reference |
| `POISONED_DESCRIPTION` | MCPTool → MCPTool (self) | Raw edges | Injection patterns detected |
| `CAN_REACH` | AgentInstance → MCPResource | HAS_ACCESS_TO | **THE critical edge** — transitive access. Also: credential chain variant (6 hops). |
| `CAN_EXFILTRATE_VIA` | AgentInstance → MCPTool | CAN_REACH | Agent reaches sensitive data AND outbound channel |
| `CAN_IMPERSONATE` | A2AAgent → A2AAgent | Raw edges | TF-IDF cosine similarity > 0.8 on skill descriptions |
| `POISONED_INSTRUCTIONS` | InstructionFile → InstructionFile (self) | Raw edges | Suspicious patterns in instruction files (imperative overrides, exfiltration commands, hidden Unicode) |
| Cross-protocol `CAN_REACH` | A2AAgent → MCPResource | HAS_ACCESS_TO + DELEGATES_TO | A2A→MCP boundary via host correlation. **What no other tool does.** |

**All edges carry:** `scan_id`, `last_seen`, `confidence` (0.0-1.0), `risk_weight`, `is_composite`, `evidence`. Composite edges also carry: `source_collector` (`'mcp'` or `'a2a'`) — used by stale edge cleanup to scope partial scan deletions.

**Post-processor execution order (dependencies):**
1. HAS_ACCESS_TO → 2. CAN_EXECUTE → 3. SHADOWS → 4. POISONED_DESCRIPTION → 5. CAN_REACH (depends on 1) → 6. CAN_EXFILTRATE_VIA (depends on 5) → 7. CAN_IMPERSONATE → 8. Cross-protocol CAN_REACH (depends on 1) → 9. RiskScore (depends on 1-8)

## Three Collectors

### Config Collector (`internal/collector/config/`)
- **Parses** 12+ MCP client config formats (Claude Desktop, Claude Code, Cursor, VS Code, Windsurf, Continue, Zed, Cline, JetBrains, Kiro, Amazon Q, Augment)
- **Key format differences:** VS Code uses `servers` key (not `mcpServers`). Windsurf uses `serverUrl` (not `url`). Zed uses `context_servers`. Cline has `autoApprove` array. Continue uses YAML.
- **Produces:** ConfigFile, AgentInstance, MCPServer, Identity, Credential, Host, InstructionFile nodes + all trust/auth edges
- **Detects:** Unpinned packages (`npx -y @pkg` without `@version`), high-entropy secrets (Shannon entropy >4.5 base64, >3.0 hex), credential patterns, instruction file poisoning
- **Parser architecture:** `ConfigParser` interface per client — `ClientName()`, `ConfigPaths()`, `Parse(path, data)`

### MCP Collector (`internal/collector/mcp/`)
- **Wraps** official Go MCP SDK. Connection: `mcp.NewClient()` → `client.Connect(ctx, transport, nil)` → auto-paginating `session.Tools(ctx, nil)`
- **Enumerates:** tools/list, resources/list, resources/templates/list, prompts/list. NEVER calls tools/call or resources/read.
- **Transports:** stdio (`mcp.CommandTransport{Command: cmd}`, env via `cmd.Env`) and Streamable HTTP (`mcp.StreamableClientTransport{Endpoint: url}`). Falls back to legacy SSE on 400/404/405.
- **Security signals per tool:** description_hash (SHA-256 canonical JSON), injection patterns, cross-references, capability_surface classification (8 categories: shell_access, file_read, file_write, network_outbound, database_access, email_send, code_execution, credential_access), annotations (readOnlyHint, destructiveHint, idempotentHint, openWorldHint — all untrusted hints)
- **Parallel:** goroutines with configurable concurrency, 120s total timeout per server, 30s init timeout, 100-page pagination safety valve

### A2A Collector (`internal/collector/a2a/`)
- **Pure HTTP client** — GET `/.well-known/agent-card.json` (v0.3.0+), fallback `/.well-known/agent.json` (legacy)
- **Version detection:** `supportedInterfaces` present → v1.0, top-level `url` → v0.3.0
- **Handles both v0.3.0 and v1.0 Agent Card formats.** v1.0 moves `url` into `supportedInterfaces[].url`, removes top-level `protocolVersion`.
- **JWS signature verification** (RFC 7515) when `signatures` field present. Unsigned = flagged.
- **Auth posture scoring:** none=100, apiKey=70, bearer=50, oauth=25, oidc=20, mTLS=10
- **Produces:** A2AAgent, A2ASkill, Host nodes + ADVERTISES_SKILL, DELEGATES_TO, SAME_AUTH_DOMAIN, RUNS_ON edges

## Ingest Format

All collectors output the same JSON schema (aligned with BloodHound OpenGraph):
```json
{
  "meta": { "version": 1, "type": "agenthound-ingest", "collector": "mcp|a2a|config", "collector_version": "0.1.0", "timestamp": "ISO8601", "scan_id": "scan-xxx" },
  "graph": {
    "nodes": [{ "id": "sha256:...", "kinds": ["MCPServer"], "properties": {...} }],
    "edges": [{ "source": "sha256:...", "target": "sha256:...", "kind": "PROVIDES_TOOL", "properties": {...} }]
  }
}
```

**Merge strategy:** MERGE by `objectid`. Same MCPServer node from Config + MCP collectors merges properties (last-write-wins). `ON MATCH SET n.previous_description_hash = n.description_hash` preserves old hash for rug pull detection.

**Normalizer:** camelCase → snake_case property keys. Timestamps → ISO 8601 UTC. Ensure `objectid` matches node `id`.

## Risk Scoring

**Edge risk weights (lower = easier to exploit):**
- TRUSTS_SERVER: none=0.1, static_key=0.3, oauth=0.7, mtls=0.9
- PROVIDES_TOOL: 0.1 (always available)
- HAS_ACCESS_TO: 0.2
- CAN_EXECUTE: 0.1
- DELEGATES_TO: none=0.1, authed=0.5
- SHADOWS: 0.4
- CAN_IMPERSONATE: 0.6

**Node risk scores (0-100):**
- Agent: 0.30×credential + 0.25×blast_radius + 0.20×auth_posture + 0.15×tool_surface + 0.10×poisoning
- Server: 0.35×auth_strength + 0.25×tool_risk + 0.20×exposure + 0.20×credential_handling
- Tool: 0.30×capability_class + 0.25×poisoning + 0.25×access_sensitivity + 0.20×input_validation

**Resource sensitivity auto-classification:** postgres/mysql/mongodb+prod → critical, file:///etc/ → critical, *.env/*.key/*.pem → critical, redis+prod → critical, DB non-prod → high, file:/// general → medium

## API Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/health` | GET | Neo4j + PG connectivity |
| `/api/v1/graph/stats` | GET | Node/edge counts by kind |
| `/api/v1/graph/nodes` | GET | List nodes (filter: kind, limit) |
| `/api/v1/graph/nodes/{id}` | GET | Node + connected edges |
| `/api/v1/graph/edges` | GET | List edges (filter: kind, source, target) |
| `/api/v1/ingest` | POST | Upload collector JSON → pipeline → post-process |
| `/api/v1/query` | POST | Raw Cypher (admin only) |
| `/api/v1/analysis/shortest-path` | POST | `{source, target, max_hops, algorithm}` |
| `/api/v1/analysis/all-paths` | POST | Bounded path enumeration |
| `/api/v1/analysis/weighted-path` | POST | Dijkstra via APOC |
| `/api/v1/analysis/findings` | GET | All composite edges as findings with severity |
| `/api/v1/analysis/prebuilt/{id}` | GET | 17 pre-built queries |
| `/api/v1/scans` | GET/POST | List history / trigger scan |
| `/api/v1/scans/{id}` | GET | Scan status |
| `/api/v1/auth/login` | POST | Username/password → JWT |
| `/api/v1/auth/tokens` | POST | Create API token (ah_ prefix) |
| `/api/v1/audit` | GET | Audit log (admin only) |

## CLI Commands

```bash
agenthound serve                              # Start API server on :8080
agenthound collect config --discover          # Discover all MCP client configs
agenthound collect mcp --discover             # Enumerate all discovered MCP servers
agenthound collect mcp --config <path>        # Enumerate servers from specific config
agenthound collect mcp --url <url>            # Scan single HTTP MCP server
agenthound collect a2a --target <url>         # Fetch single Agent Card
agenthound collect a2a --targets <url1,url2>  # Fetch multiple Agent Cards
agenthound ingest <file.json>                 # Ingest collector output → Neo4j + post-process
agenthound query "<cypher>"                   # Execute raw Cypher
agenthound query --prebuilt <query-id>        # Run pre-built query
agenthound query --findings --severity critical
```

## Implementation Phases

### Phase 1: Foundation (Weeks 1-2)
Go scaffolding, Docker Compose, Neo4j schema (14 constraints + indexes), ingest pipeline, basic API, CI/CD.
**Exit:** `docker compose up` → healthy. `agenthound ingest testdata/valid_mcp_scan.json` → nodes in Neo4j.

### Phase 2: Collectors (Weeks 3-5)
Config Collector (week 3, no network, 12 parsers) → MCP Collector (week 4, Go SDK) → A2A Collector (week 5, HTTP+JSON).
**Exit:** `collect config --discover` produces valid JSON. `collect mcp --config` enumerates real servers. MCPServer IDs match across collectors.

### Phase 3: Post-Processing & Attack Paths (Weeks 6-7)
9 post-processors in dependency order, risk scoring, pathfinding API, 17 pre-built queries, CLI query mode.
**Exit:** After scanning test env, AgentHound finds known attack paths. CAN_REACH edges exist. Risk scores ordered correctly.

### Phase 4: Frontend (Weeks 8-10)
React+Vite SPA: Dashboard, Graph Explorer (Sigma.js), Pathfinder, Entity Inspector, Scan Manager, Query Library. Embedded via `go:embed`.
**Exit:** User opens localhost:8080, sees dashboard, clicks nodes, runs pathfinding, sees highlighted paths.

### Phase 5: Hardening & Release (Weeks 11-12)
Auth (bcrypt+JWT+API tokens+RBAC), audit logging, error handling, docs, testing (>80% coverage), performance validation, security review, release artifacts (GHCR images + cross-compiled binaries via GoReleaser).

## Key Implementation Constraints

1. **Neo4j version compat:** Schema init must detect version via `CALL dbms.components()` and use 4.4 (`ON...ASSERT`) or 5.x (`FOR...REQUIRE`) syntax.
2. **APOC fallback:** All APOC-dependent code needs non-APOC fallbacks. APOC only required for Dijkstra. Node writes: group by kind, run separate MERGE per kind. Edge writes: `edgeKindCypher` map with per-kind Cypher strings.
3. **Property keys:** Neo4j is case-sensitive. All properties stored as snake_case. Normalizer converts camelCase from collector JSON.
4. **Batch writes:** 1000 operations per Neo4j transaction. Use `UNWIND $nodes AS node` pattern.
5. **go:embed constraint:** Go forbids `..` in embed paths. Makefile copies `ui/dist` → `internal/api/ui/dist` before `go build`.
6. **Sigma.js v3:** Pass `MultiDirectedGraph` constructor (not instance) to `SigmaContainer`. Use `useLoadGraph()` hook. `staleTime: 30_000` on React Query to prevent unnecessary re-renders.
7. **MCP SDK:** `mcp.CommandTransport` env vars set on `exec.Cmd.Env`, not on transport. `client.Connect()` handles full init handshake. Auto-paginating iterators handle `NextCursor`.
8. **A2A version detection:** `supportedInterfaces` → v1.0, top-level `url` → v0.3.0. Must handle both.
9. **Credential safety:** Config Collector hashes credential values by default (SHA-256). `--include-credential-values` for audit mode.
10. **Stale edge cleanup:** Only delete composite edges whose source collector ran in current scan. Prevents ping-pong on partial scans.

## OWASP Coverage

AgentHound maps all findings to OWASP MCP Top 10 (MCP01-MCP10) and OWASP Agentic Top 10 (ASI01-ASI10). Full/partial coverage documented in `architecture/07-threat-model.md`.

## Pre-Built Queries (17)

| ID | Category |
|----|----------|
| `agents-shell-access` | Critical Paths |
| `shortest-to-database` | Critical Paths |
| `cross-protocol-paths` | Critical Paths |
| `exfiltration-routes` | Critical Paths |
| `credential-chain` | Critical Paths |
| `poisoned-tools` | Vulnerabilities |
| `tool-shadowing` | Vulnerabilities |
| `no-auth-servers` | Vulnerabilities |
| `no-auth-a2a` | Vulnerabilities |
| `rug-pull` | Vulnerabilities |
| `unpinned-packages` | Supply Chain |
| `instruction-poisoning` | Supply Chain |
| `unsigned-cards` | Supply Chain |
| `high-entropy-secrets` | Supply Chain |
| `chokepoint-servers` | Chokepoints |
| `chokepoint-tools` | Chokepoints |
| `unpinned-shell` | Combined |

## Node Visual Encoding (Frontend)

| Kind | Color | Size Basis |
|------|-------|-----------|
| AgentInstance | `#4A90D9` blue | risk_score |
| MCPServer | `#50C878` green | tool count |
| MCPTool | `#F5A623` orange | capability risk |
| MCPResource | `#D0021B` red | sensitivity |
| A2AAgent | `#7B68EE` purple | skill count |
| A2ASkill | `#9B59B6` light purple | fixed |
| Identity | `#8E8E93` gray | fixed |
| Credential | `#FF6B6B` warning red | exposure risk |
| ConfigFile | `#95A5A6` silver | server count |
| Host | `#2C3E50` dark | fixed |

## Config Defaults

```
AGENTHOUND_NEO4J_URI=bolt://localhost:7687
AGENTHOUND_NEO4J_USER=neo4j
AGENTHOUND_NEO4J_PASSWORD=agenthound
AGENTHOUND_PG_URI=postgres://agenthound:agenthound@localhost:5432/agenthound?sslmode=disable
AGENTHOUND_API_PORT=8080
AGENTHOUND_LOG_LEVEL=info
```

## Reference Docs (in repo)

| Path | Content |
|------|---------|
| `architecture/01-vision.md` through `architecture/09-roadmap.md` | PRD sections (gitignored, local only) |
| `collectors/01-mcp-collector.md` | MCP collector spec + Go SDK usage |
| `collectors/02-a2a-collector.md` | A2A collector spec + card schema |
| `collectors/03-config-collector.md` | Config collector spec + 12 client formats |
| `collectors/04-graph-pipeline.md` | Full ingest → post-process → query pipeline |
| `collectors/05-visual-architecture.md` | Mermaid diagrams of the graph |
| `collectors/06-ui-scenarios.md` | UI graph rendering scenarios |
| `plan/phase01.md` through `plan/phase05.md` | Detailed implementation plans per phase |
| `research/mcp-spec-2025-11-25.md` | MCP protocol reference (JSON-RPC, transports, OAuth, Tasks) |
| `research/sigma-react-graphology-stack.md` | Frontend stack research |
