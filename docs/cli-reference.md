# CLI Reference

AgentHound ships as **two binaries**: `agenthound` (collector) and `agenthound-server` (single-user analysis server). Both use [Cobra](https://github.com/spf13/cobra); all commands support `--help`.

## Collector global flags (`agenthound`)

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `--output` | `AGENTHOUND_OUTPUT` | (auto-named in CWD) | Write scan JSON to this path. Use `-` for stdout (pipe target). When unset, defaults to `./scan-<scan_id>.json` in the current working directory. |
| `--concurrency` | `AGENTHOUND_CONCURRENCY` | `0` (auto) | Max parallel collector workers. |
| `--log-level` | `AGENTHOUND_LOG_LEVEL` | `info` | Log level: debug, info, warn, error. |
| `--quiet` | `AGENTHOUND_QUIET=1` | `false` | Suppress non-error log output. |
| `--log-json` | `AGENTHOUND_LOG_JSON=1` | `false` | Emit logs as JSON instead of text. |

The collector is offline-by-default. It does not phone home to a server. To ingest the resulting JSON on the operator's box, use `agenthound-server ingest <file>` or `agenthound-server ingest -` (stdin), or drag-drop the file in the UI's `Scan Manager → Import scan` dialog.

## Server global flags (`agenthound-server`)

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `--bind` | `AGENTHOUND_BIND` | `127.0.0.1:8080` | Bind address `host:port`. Set to `0.0.0.0:8080` only inside a trusted network. |
| `--neo4j-uri` | `AGENTHOUND_NEO4J_URI` | `bolt://localhost:7687` | Neo4j connection URI. |
| `--neo4j-user` | `AGENTHOUND_NEO4J_USER` | `neo4j` | Neo4j username. |
| `--neo4j-password` | `AGENTHOUND_NEO4J_PASSWORD` | `agenthound` | Neo4j password. |
| `--pg-uri` | `AGENTHOUND_PG_URI` | `postgres://agenthound:agenthound@localhost:5432/agenthound?sslmode=disable` | PostgreSQL URI. |
| `--cors-origins` | `AGENTHOUND_CORS_ORIGINS` | `http://localhost:8080` | Comma-separated CORS origins. |
| `--log-level` | `AGENTHOUND_LOG_LEVEL` | `info` | Log level. |

Priority: CLI flag > environment variable > default value.

---

## `agenthound-server serve`

Start the API server and embedded UI.

```bash
agenthound-server serve [--bind 127.0.0.1:8080]
```

The server has **no application-layer authentication**. Default bind is loopback only; expose remotely only over your own VPN/SSH tunnel/firewall. See [`security.md`](security.md).

Initializes Neo4j schema (constraints + indexes) and PostgreSQL migrations automatically. Embeds the React SPA at the root URL; API at `/api/v1/*`. Graceful shutdown on SIGINT/SIGTERM with a 10-second drain timeout.

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
| `--output` | (auto-named in CWD) | Write merged JSON to this path. Use `-` for stdout (pipe target). When unset, defaults to `./scan-<scan_id>.json` in CWD. |

### Examples

```bash
agenthound scan                                        # Full scan (config + MCP); writes ./scan-<scan_id>.json in CWD
agenthound scan --config                               # Config files only (offline)
agenthound scan --config --path ~/.cursor/mcp.json     # Single config file
agenthound scan --mcp                                  # MCP servers only
agenthound scan --mcp --url https://mcp.example.com    # Single HTTP MCP server
agenthound scan --a2a --target https://agent.example.com
agenthound scan --a2a --discover-domain example.com
agenthound scan --output scan.json                     # Explicit file path
agenthound scan --output - | agenthound-server ingest - # Stream JSON over stdin
agenthound scan --output - | ssh op-box 'agenthound-server ingest -' # Pipe over SSH
```

For CI/CD gating on findings, run the analysis on the operator's server:

```bash
agenthound-server query --findings --severity critical --fail-on critical
```

---

## `agenthound-server ingest`

Ingest collector JSON output into the graph database.

```bash
agenthound-server ingest <file.json>
agenthound-server ingest -
```

Takes exactly one argument: the path to a collector JSON output file, or `-` to read from stdin. The stdin form is the standard pipe target for `agenthound scan --output -`:

```bash
agenthound scan --output - | agenthound-server ingest -
ssh target 'agenthound scan --output -' | agenthound-server ingest -
```

The same pipeline (validate → normalize → deduplicate → write → post-process) backs the `POST /api/v1/ingest` HTTP endpoint and the UI's `Scan Manager → Import scan` drag-drop dialog. All three entry points produce the same graph state.

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

## `agenthound-server query`

Execute queries against the graph database (runs on the operator's box, not the collector). Supports four mutually exclusive modes.

### Raw Cypher

```bash
agenthound-server query "MATCH (n:MCPServer) RETURN n.name, n.transport"
```

### Pre-built query

```bash
agenthound-server query --prebuilt <query-id>
```

Run `agenthound-server query --prebuilt ""` to see available query IDs (the command prints the list on unknown ID).

### Findings

```bash
agenthound-server query --findings [--severity critical|high|medium|low]
agenthound-server query --findings --fail-on critical         # CI: exit 1 if critical findings
```

Lists all composite edges as security findings with severity classification.

### Shortest path

```bash
agenthound-server query --shortest-path --from <Kind:name> --to <Kind:name>
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
