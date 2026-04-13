# CLI Reference

AgentHound uses [Cobra](https://github.com/spf13/cobra) for its CLI. All commands support `--help` for usage details.

## Global flags

These flags apply to all commands. Each can also be set via environment variable.

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `--neo4j-uri` | `AGENTHOUND_NEO4J_URI` | `bolt://localhost:7687` | Neo4j connection URI |
| `--neo4j-user` | `AGENTHOUND_NEO4J_USER` | `neo4j` | Neo4j username |
| `--neo4j-password` | `AGENTHOUND_NEO4J_PASSWORD` | `agenthound` | Neo4j password |
| `--pg-uri` | `AGENTHOUND_PG_URI` | `postgres://agenthound:agenthound@localhost:5432/agenthound?sslmode=disable` | PostgreSQL URI |
| `--log-level` | `AGENTHOUND_LOG_LEVEL` | `info` | Log level: debug, info, warn, error |

Priority: CLI flag > environment variable > default value.

---

## `agenthound serve`

Start the API server and embedded UI.

```bash
agenthound serve [--port 8080]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | API server listen port |

On first start, creates a default `admin` user with the password from `AGENTHOUND_ADMIN_PASSWORD` (default: `agenthound`). Initializes Neo4j schema (constraints + indexes) and PostgreSQL migrations automatically.

The server embeds the React SPA and serves it at the root URL. The API is mounted at `/api/v1/*`.

Graceful shutdown on SIGINT/SIGTERM with a 10-second drain timeout.

---

## `agenthound scan`

Discover and enumerate MCP servers, A2A agents, and client configurations, then analyze the trust graph for attack paths.

```bash
agenthound scan [flags]
```

By default, runs config discovery + MCP enumeration + ingest + post-processing analysis. Use `--config`, `--mcp`, or `--a2a` to run individual collectors.

### Collector selection

| Flag | Description |
|------|-------------|
| `--config` | Run config collector only |
| `--mcp` | Run MCP collector only |
| `--a2a` | Run A2A collector only |

When none of `--config`, `--mcp`, or `--a2a` are specified, runs config + MCP (the default workflow).

### Config collector flags

| Flag | Default | Description |
|------|---------|-------------|
| `--path` | | Path to a single config file (overrides auto-discovery) |
| `--paths` | | Comma-separated paths to multiple config files |
| `--project-dir` | | Project directory for instruction file discovery |
| `--include-credential-values` | `false` | Include raw credential values instead of SHA-256 hashes |

**Supported clients (12):** Claude Desktop, Claude Code, Cursor, VS Code, Windsurf, Continue, Zed, Cline, JetBrains, Kiro, Amazon Q, Augment.

**What it produces:** ConfigFile, AgentInstance, MCPServer, Identity, Credential, Host, InstructionFile nodes plus trust/auth/config edges.

**Security signals:** Unpinned packages, hardcoded secrets (Shannon entropy), instruction file poisoning.

### MCP collector flags

| Flag | Default | Description |
|------|---------|-------------|
| `--url` | | URL of a single HTTP MCP server (overrides auto-discovery) |

AgentHound never calls `tools/call` or `resources/read`. It is read-only and safe to run against production servers.

**Enumerates:** `tools/list`, `resources/list`, `resources/templates/list`, `prompts/list`.

**Transports:** stdio (`mcp.CommandTransport`) and Streamable HTTP (with legacy SSE fallback).

**Security signals per tool:** description hashing (rug pull detection), injection pattern scanning, cross-reference detection, capability surface classification.

### A2A collector flags

| Flag | Default | Description |
|------|---------|-------------|
| `--target` | | URL of a single A2A agent |
| `--targets` | | Comma-separated URLs of multiple agents |
| `--targets-file` | | File with agent URLs (one per line) |
| `--discover-domain` | | Domains to probe for well-known agent cards |
| `--auth-token` | | Bearer token for authenticated agents |

At least one of `--target`, `--targets`, `--targets-file`, or `--discover-domain` is required when using `--a2a`.

`--discover-domain example.com` probes `https://example.com/.well-known/agent-card.json`.

**Version support:** v1.0 (detected by `supportedInterfaces`), v0.3.0 (detected by top-level `url`), legacy fallback to `/.well-known/agent.json`.

**Security signals:** JWS signature verification (RFC 7515), auth posture scoring (none=100 ... mTLS=10), unsigned card flagging.

### Shared flags

| Flag | Default | Description |
|------|---------|-------------|
| `--concurrency` | `5` | Max parallel connections |
| `--timeout` | `120s` | Timeout per server/agent |
| `--insecure` | `false` | Skip TLS verification |

### Output flags

| Flag | Default | Description |
|------|---------|-------------|
| `--output` | | Export merged JSON to file (skips ingest and analysis) |
| `--fail-on` | | Exit 1 if findings at or above severity: `critical`, `high`, `medium`, `low` |

### Examples

```bash
agenthound scan                                        # Full scan (config + MCP)
agenthound scan --config                               # Config files only (offline)
agenthound scan --config --path ~/.cursor/mcp.json     # Single config file
agenthound scan --mcp                                  # MCP servers only
agenthound scan --mcp --url https://mcp.example.com    # Single HTTP MCP server
agenthound scan --a2a --target https://agent.example.com
agenthound scan --a2a --discover-domain example.com
agenthound scan --output scan.json                     # Export without ingesting
agenthound scan --fail-on critical                     # CI/CD gate
```

---

## `agenthound ingest`

Ingest collector JSON output into the graph database.

```bash
agenthound ingest <file.json>
```

Takes exactly one argument: the path to a collector JSON output file.

### Pipeline stages

1. **Validate** -- JSON schema validation (meta, node kinds, edge kinds)
2. **Normalize** -- camelCase to snake_case property keys, timestamps to ISO 8601 UTC, objectid sync
3. **Deduplicate** -- merge by objectid (last-write-wins for properties)
4. **Write** -- batch MERGE into Neo4j (1000 operations per transaction)
5. **Post-process** -- compute composite edges and risk scores

### Output

```
Ingest complete:
  Scan ID:       scan-abc123
  Nodes written: 47
  Edges written: 82
  Duration:      1.23s
```

---

## `agenthound query`

Execute queries against the graph database. Supports four mutually exclusive modes.

### Raw Cypher

```bash
agenthound query "MATCH (n:MCPServer) RETURN n.name, n.transport"
```

### Pre-built query

```bash
agenthound query --prebuilt <query-id>
```

Run `agenthound query --prebuilt ""` to see available query IDs (the command prints the list on unknown ID).

### Findings

```bash
agenthound query --findings [--severity critical|high|medium|low]
agenthound query --findings --fail-on critical         # CI: exit 1 if critical findings
```

Lists all composite edges as security findings with severity classification.

### Shortest path

```bash
agenthound query --shortest-path --from <Kind:name> --to <Kind:name>
```

Node references use `Kind:name` format, e.g., `AgentInstance:claude-desktop` or `MCPResource:postgres://prod`.

### Common flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `table` | Output format: `table` or `json` |
| `--fail-on` | | Exit 1 if findings at or above severity: `critical`, `high`, `medium`, `low` (applies to `--findings` mode) |

### Pre-built query IDs

| ID | Category | Severity |
|----|----------|----------|
| `agents-shell-access` | Critical Paths | critical |
| `shortest-to-database` | Critical Paths | critical |
| `cross-protocol-paths` | Critical Paths | critical |
| `exfiltration-routes` | Critical Paths | critical |
| `credential-chain` | Critical Paths | critical |
| `unpinned-shell` | Combined | critical |
| `poisoned-tools` | Vulnerabilities | high |
| `tool-shadowing` | Vulnerabilities | high |
| `no-auth-servers` | Vulnerabilities | high |
| `no-auth-a2a` | Vulnerabilities | high |
| `rug-pull` | Vulnerabilities | high |
| `instruction-poisoning` | Supply Chain | high |
| `high-entropy-secrets` | Supply Chain | high |
| `unpinned-packages` | Supply Chain | medium |
| `unsigned-cards` | Supply Chain | medium |
| `chokepoint-servers` | Chokepoints | medium |
| `chokepoint-tools` | Chokepoints | medium |

### Valid node kinds for `--from` and `--to`

`MCPServer`, `MCPTool`, `MCPResource`, `MCPPrompt`, `A2AAgent`, `A2ASkill`, `AgentInstance`, `Identity`, `Credential`, `Host`, `ConfigFile`, `InstructionFile`, `ResourceGroup`, `TrustZone`
