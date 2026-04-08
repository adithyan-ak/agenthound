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

## `agenthound collect config`

Parse MCP client configuration files to discover agent-server trust relationships, credentials, and instruction files.

```bash
agenthound collect config [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--discover` | `false` | Auto-discover all known MCP client config files |
| `--path` | | Path to a single config file |
| `--paths` | | Comma-separated paths to multiple config files |
| `--output` | stdout | Write JSON output to file |
| `--ingest` | `false` | Ingest directly into graph database (requires running Neo4j + PG) |
| `--include-credential-values` | `false` | Include raw credential values instead of SHA-256 hashes |
| `--project-dir` | | Project directory for instruction file discovery |

At least one of `--discover`, `--path`, or `--paths` is required.

### Supported clients

Claude Desktop, Claude Code, Cursor, VS Code, Windsurf, Continue, Zed, Cline, JetBrains, Kiro, Amazon Q, Augment.

### What it produces

- **ConfigFile** nodes for each config file found
- **AgentInstance** nodes for each MCP client
- **MCPServer** nodes for each configured server
- **Identity** and **Credential** nodes for auth configuration
- **Host** nodes extracted from server endpoints
- **InstructionFile** nodes for discovered instruction/rules files
- Trust, auth, and config edges between all nodes

### Security signals detected

- Unpinned packages (`npx -y @pkg` without version pin)
- Hardcoded secrets (Shannon entropy > 4.5 for base64, > 3.0 for hex)
- Instruction file poisoning (imperative overrides, exfiltration patterns, hidden Unicode)

---

## `agenthound collect mcp`

Connect to MCP servers and enumerate their capabilities.

```bash
agenthound collect mcp [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--discover` | `false` | Enumerate all servers from discovered config files |
| `--config` | | Path to MCP client config file |
| `--url` | | URL of a single HTTP MCP server |
| `--output` | stdout | Write JSON output to file |
| `--ingest` | `false` | Ingest directly into graph database |
| `--concurrency` | `5` | Max parallel server connections |
| `--timeout` | `120s` | Timeout per server |
| `--insecure` | `false` | Skip TLS verification for HTTP servers |

At least one of `--discover`, `--config`, or `--url` is required.

### What it enumerates

- `tools/list` -- tool names, descriptions, input schemas, annotations
- `resources/list` and `resources/templates/list` -- resource URIs, types, sizes
- `prompts/list` -- prompt names, descriptions, arguments

AgentHound never calls `tools/call` or `resources/read`. It is read-only and safe to run against production servers.

### Transport support

- **stdio** -- launches server process via `mcp.CommandTransport`
- **Streamable HTTP** -- connects via `mcp.StreamableClientTransport`, falls back to legacy SSE on 400/404/405

### Security signals per tool

- **description_hash** -- SHA-256 of canonical description for rug pull detection
- **injection_patterns** -- prompt injection markers in descriptions
- **cross_references** -- references to other tools/servers in descriptions
- **capability_surface** -- classified into: shell_access, file_read, file_write, network_outbound, database_access, email_send, code_execution, credential_access

---

## `agenthound collect a2a`

Fetch A2A Agent Cards and collect agent/skill data.

```bash
agenthound collect a2a [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--target` | | URL of a single A2A agent |
| `--targets` | | Comma-separated URLs of multiple agents |
| `--targets-file` | | File with agent URLs (one per line) |
| `--discover-domain` | | Domains to probe for well-known agent cards |
| `--output` | stdout | Write JSON output to file |
| `--ingest` | `false` | Ingest directly into graph database |
| `--auth-token` | | Bearer token for authenticated agents |
| `--insecure` | `false` | Skip TLS verification |
| `--concurrency` | `5` | Max parallel agent fetches |
| `--timeout` | `15s` | Timeout per agent |

At least one of `--target`, `--targets`, `--targets-file`, or `--discover-domain` is required.

### Domain discovery

`--discover-domain example.com` probes `https://example.com/.well-known/agent-card.json`.

### Version support

- **v1.0** -- detected by presence of `supportedInterfaces` field
- **v0.3.0** -- detected by top-level `url` field
- **Legacy** -- falls back to `/.well-known/agent.json` if the v0.3.0+ path returns 404

### Security signals

- JWS signature verification (RFC 7515) when `signatures` field is present
- Auth posture scoring: none=100, apiKey=70, bearer=50, oauth=25, oidc=20, mTLS=10
- Unsigned cards are flagged

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
