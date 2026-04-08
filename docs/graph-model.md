# Graph Data Model

AgentHound builds a directed trust graph where nodes represent entities in AI agent infrastructure and edges represent exploitable relationships. Direction follows the flow of access and control.

The core traversal pattern: `Agent -> Server -> Tool -> Resource`.

## Node types

### Collector-produced (12 types)

| Label | Source | Key properties |
|-------|--------|----------------|
| `MCPServer` | Config + MCP | `name`, `endpoint`, `transport` (stdio/http), `auth_method`, `protocol_version`, `instructions`, `capabilities`, `is_pinned`, `has_tasks_capability` |
| `MCPTool` | MCP | `name`, `description`, `input_schema`, `output_schema`, `annotations`, `description_hash`, `capability_surface[]`, `has_injection_patterns`, `has_cross_references` |
| `MCPResource` | MCP | `uri`, `name`, `mime_type`, `size`, `uri_scheme`, `sensitivity` |
| `MCPPrompt` | MCP | `name`, `description`, `arguments` |
| `A2AAgent` | A2A | `name`, `description`, `url`, `provider`, `version`, `protocol_versions`, `capabilities`, `security_schemes`, `auth_method`, `is_signed`, `signature_valid`, `card_hash` |
| `A2ASkill` | A2A | `id`, `name`, `description`, `input_modes`, `output_modes`, `description_hash`, `has_injection_patterns` |
| `AgentInstance` | Config | `name`, `framework`, `config_path` |
| `Identity` | Config + MCP | `type` (none/apiKey/oauth/bearer/mtls), `scope`, `is_static` |
| `Credential` | Config | `type` (envVar/hardcoded/vaultRef/inputPrompt), `name`, `source`, `is_exposed`, `high_entropy` |
| `Host` | Config + A2A | `hostname`, `ip`, `is_local`, `is_private`, `is_public` |
| `ConfigFile` | Config | `path`, `client`, `server_count` |
| `InstructionFile` | Config | `path`, `type`, `hash`, `is_suspicious` |

Instruction file types: `agents.md`, `claude.md`, `cursorrules`, `copilot-instructions`, `memory.md`.

### Synthetic (2 types, created by post-processors)

| Label | Source | Key properties |
|-------|--------|----------------|
| `ResourceGroup` | Post-processor | `type`, `sensitivity` |
| `TrustZone` | Post-processor | `name`, `level`, `node_count` |

## Node ID strategy

All node IDs are deterministic, content-based SHA-256 hashes. This ensures identical entities from different collectors merge correctly.

| Node type | ID computation |
|-----------|---------------|
| `MCPServer` | `SHA-256("MCPServer:" + transport + ":" + endpoint + ":" + sorted_args)` |
| `MCPTool` | `SHA-256("MCPTool:" + server_id + ":" + tool_name)` |
| `MCPResource` | `SHA-256("MCPResource:" + server_id + ":" + resource_uri)` |
| `A2AAgent` | `SHA-256("A2AAgent:" + agent_card_url)` |
| `A2ASkill` | `SHA-256("A2ASkill:" + agent_id + ":" + skill_id)` |
| `AgentInstance` | `SHA-256("AgentInstance:" + config_file_id + ":" + client_name)` |
| `ConfigFile` | `SHA-256("ConfigFile:" + absolute_path)` |
| `Host` | `SHA-256("Host:" + hostname_or_ip)` |

The MCPServer ID must match between the Config Collector and MCP Collector. This is the merge point that connects trust relationships (who trusts what) to capabilities (what a server exposes).

## Edge types

### Direct edges (from collectors)

| Edge | Direction | Collector | Meaning |
|------|-----------|-----------|---------|
| `TRUSTS_SERVER` | AgentInstance -> MCPServer | Config | Agent trusts this server to provide tools |
| `PROVIDES_TOOL` | MCPServer -> MCPTool | MCP | Server exposes this tool |
| `PROVIDES_RESOURCE` | MCPServer -> MCPResource | MCP | Server exposes this resource |
| `PROVIDES_PROMPT` | MCPServer -> MCPPrompt | MCP | Server exposes this prompt |
| `ADVERTISES_SKILL` | A2AAgent -> A2ASkill | A2A | Agent advertises this skill |
| `DELEGATES_TO` | A2AAgent -> A2AAgent | A2A | Agent delegates tasks to another agent |
| `AUTHENTICATES_WITH` | MCPServer/A2AAgent -> Identity | Config/A2A | Entity uses this auth identity |
| `USES_CREDENTIAL` | Identity -> Credential | Config | Identity backed by this credential |
| `RUNS_ON` | MCPServer/A2AAgent -> Host | Config/A2A | Entity runs on this host |
| `CONFIGURED_IN` | MCPServer -> ConfigFile | Config | Server defined in this config file |
| `HAS_ENV_VAR` | MCPServer -> Credential | Config | Server has access to this env var |
| `LOADS_INSTRUCTIONS` | AgentInstance -> InstructionFile | Config | Agent loads this instruction file |
| `SAME_AUTH_DOMAIN` | A2AAgent -> A2AAgent | A2A | Agents share an auth domain |

### Composite edges (computed by post-processors)

These edges are computed from graph state after ingestion. They represent derived attack paths and security findings.

| Edge | Direction | Depends on | Meaning |
|------|-----------|------------|---------|
| `HAS_ACCESS_TO` | MCPTool -> MCPResource | Raw edges | Tool's capability surface matches resource URI scheme |
| `CAN_EXECUTE` | MCPTool -> Host | Raw edges | Tool has shell_access or code_execution capability |
| `SHADOWS` | MCPTool -> MCPTool | Raw edges | Tool on another server references this tool's name/description |
| `POISONED_DESCRIPTION` | MCPTool -> MCPTool (self) | Raw edges | Tool description contains injection patterns |
| `CAN_REACH` | AgentInstance -> MCPResource | HAS_ACCESS_TO | Transitive access: agent can reach a resource through trust chain. Also: credential chain variant (up to 6 hops) |
| `CAN_EXFILTRATE_VIA` | AgentInstance -> MCPTool | CAN_REACH | Agent can reach sensitive data AND has an outbound exfiltration channel |
| `CAN_IMPERSONATE` | A2AAgent -> A2AAgent | Raw edges | TF-IDF cosine similarity > 0.8 on skill descriptions |
| `POISONED_INSTRUCTIONS` | InstructionFile -> InstructionFile (self) | Raw edges | Suspicious patterns: imperative overrides, exfiltration commands, hidden Unicode |
| Cross-protocol `CAN_REACH` | A2AAgent -> MCPResource | HAS_ACCESS_TO + DELEGATES_TO | A2A agent reaches MCP resources via host correlation across protocol boundary |

### Post-processor execution order

Processors run in dependency order:

1. HAS_ACCESS_TO
2. CAN_EXECUTE
3. SHADOWS
4. POISONED_DESCRIPTION
5. CAN_REACH (depends on 1)
6. CAN_EXFILTRATE_VIA (depends on 5)
7. CAN_IMPERSONATE
8. POISONED_INSTRUCTIONS
9. Cross-protocol CAN_REACH (depends on 1)
10. RiskScore (depends on 1-9)

### Edge properties

All edges carry:

| Property | Type | Description |
|----------|------|-------------|
| `scan_id` | string | Scan that created/updated this edge |
| `last_seen` | ISO 8601 | Timestamp of last observation |
| `confidence` | float | 0.0 to 1.0 confidence score |
| `risk_weight` | float | Lower = easier to exploit (used by Dijkstra) |
| `is_composite` | bool | True for post-processed edges |
| `evidence` | string | Human-readable explanation |

Composite edges also carry `source_collector` (`mcp` or `a2a`) for scoped stale edge cleanup.

## Risk scoring

### Edge risk weights

Lower weight = easier to exploit = higher risk.

| Edge | Weight | Notes |
|------|--------|-------|
| `TRUSTS_SERVER` (no auth) | 0.1 | Trivial to exploit |
| `TRUSTS_SERVER` (static key) | 0.3 | Key may be in config |
| `TRUSTS_SERVER` (oauth) | 0.7 | Requires token theft |
| `TRUSTS_SERVER` (mtls) | 0.9 | Requires cert theft |
| `PROVIDES_TOOL` | 0.1 | Always available |
| `HAS_ACCESS_TO` | 0.2 | |
| `CAN_EXECUTE` | 0.1 | |
| `DELEGATES_TO` (no auth) | 0.1 | |
| `DELEGATES_TO` (authed) | 0.5 | |
| `SHADOWS` | 0.4 | |
| `CAN_IMPERSONATE` | 0.6 | |

### Node risk scores (0-100)

**Agent risk:**

| Factor | Weight | Source |
|--------|--------|--------|
| Credential exposure | 0.30 | Hardcoded secrets, env vars |
| Blast radius | 0.25 | Number of reachable resources |
| Auth posture | 0.20 | Weakest auth in trust chain |
| Tool surface | 0.15 | Dangerous capability classes |
| Poisoning | 0.10 | Poisoned tools/instructions |

**Server risk:**

| Factor | Weight | Source |
|--------|--------|--------|
| Auth strength | 0.35 | Auth method configured |
| Tool risk | 0.25 | Highest capability class |
| Exposure | 0.20 | Public/private/local host |
| Credential handling | 0.20 | Env var exposure, static keys |

**Tool risk:**

| Factor | Weight | Source |
|--------|--------|--------|
| Capability class | 0.30 | shell_access, code_execution, etc. |
| Poisoning | 0.25 | Injection patterns, cross-references |
| Access sensitivity | 0.25 | Sensitivity of reachable resources |
| Input validation | 0.20 | Schema presence and restrictiveness |

### Resource sensitivity classification

| Pattern | Sensitivity |
|---------|------------|
| postgres/mysql/mongodb + prod | critical |
| `file:///etc/` | critical |
| `*.env`, `*.key`, `*.pem` | critical |
| redis + prod | critical |
| Database (non-prod) | high |
| `file:///` (general) | medium |

## Merge strategy

Nodes are merged by `objectid` using Cypher `MERGE`. When the same MCPServer appears from both Config and MCP collectors:

- Properties use last-write-wins semantics
- `ON MATCH SET n.previous_description_hash = n.description_hash` preserves the old hash for rug pull detection
- Edges accumulate (both collectors contribute different edge types)

## Stale edge cleanup

On partial scans (e.g., only MCP collector ran), only composite edges whose `source_collector` matches the current scan's collector are deleted and recomputed. This prevents ping-pong deletion when collectors run independently.
