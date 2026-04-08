# API Reference

All endpoints are served at `/api/v1/*` on the configured port (default: 8080).

## Authentication

All endpoints except `/health` and `/auth/login` require authentication via Bearer token in the `Authorization` header:

```
Authorization: Bearer <jwt-or-api-token>
```

Two token types are supported:

- **JWT** -- obtained from `/auth/login`, expires in 24 hours
- **API token** -- created via `/auth/tokens`, prefixed with `ah_`, no expiry

### Roles

| Role | Access level |
|------|-------------|
| `admin` | Full access: raw Cypher queries, user management, audit log |
| `analyst` | Ingest data, run pathfinding, manage API tokens, trigger scans |
| `viewer` | Read-only: view graph, findings, scans, pre-built queries |

### Error format

All errors return a structured JSON response:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "description of the problem"
  }
}
```

Error codes: `VALIDATION_ERROR` (400), `UNAUTHORIZED` (401), `FORBIDDEN` (403), `NOT_FOUND` (404), `RATE_LIMITED` (429), `SERVICE_UNAVAILABLE` (503), `INTERNAL_ERROR` (500).

### Rate limiting

- General: 100 requests/minute per IP
- Ingest: 20 requests/minute per IP
- Raw Cypher: 10 requests/minute per IP

Returns `429 Too Many Requests` when exceeded.

---

## Health

### `GET /api/v1/health`

**Auth:** None

Returns connectivity status for Neo4j and PostgreSQL.

```json
{
  "status": "healthy",
  "neo4j": "connected",
  "postgres": "connected"
}
```

---

## Auth

### `POST /api/v1/auth/login`

**Auth:** None

```json
// Request
{
  "username": "admin",
  "password": "agenthound"
}

// Response (200)
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2026-04-09T12:00:00Z",
  "user": {
    "id": "uuid",
    "username": "admin",
    "role": "admin"
  }
}
```

### `POST /api/v1/auth/tokens`

**Auth:** analyst+

Create an API token for programmatic access.

```json
// Request
{
  "name": "ci-pipeline"
}

// Response (201)
{
  "id": "uuid",
  "name": "ci-pipeline",
  "token": "ah_abc123...",
  "created_at": "2026-04-08T12:00:00Z"
}
```

The full token value is only returned at creation time.

### `GET /api/v1/auth/tokens`

**Auth:** analyst+

List API tokens for the current user.

### `DELETE /api/v1/auth/tokens/{id}`

**Auth:** analyst+

Revoke an API token.

### `POST /api/v1/auth/users`

**Auth:** admin

Create a new user.

```json
{
  "username": "analyst1",
  "password": "secure-password",
  "role": "analyst"
}
```

### `GET /api/v1/auth/users`

**Auth:** admin

List all users.

---

## Graph

### `GET /api/v1/graph/stats`

**Auth:** viewer+

Returns node and edge counts by kind.

```json
{
  "nodes": {
    "MCPServer": 12,
    "MCPTool": 47,
    "MCPResource": 23,
    "AgentInstance": 3
  },
  "edges": {
    "TRUSTS_SERVER": 15,
    "PROVIDES_TOOL": 47,
    "CAN_REACH": 8
  }
}
```

### `GET /api/v1/graph/nodes`

**Auth:** viewer+

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `kind` | string | (all) | Filter by node label |
| `limit` | int | 100 | Max results (1-10000) |

### `GET /api/v1/graph/nodes/{id}`

**Auth:** viewer+

Returns a single node with all connected edges.

```json
{
  "node": {
    "id": "sha256:...",
    "kinds": ["MCPServer"],
    "properties": { "name": "...", "transport": "stdio" }
  },
  "edges": [
    { "source": "...", "target": "...", "kind": "PROVIDES_TOOL", "properties": {} }
  ]
}
```

### `GET /api/v1/graph/edges`

**Auth:** viewer+

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `kind` | string | (all) | Filter by edge kind |
| `source` | string | | Filter by source node ID |
| `target` | string | | Filter by target node ID |
| `limit` | int | 100 | Max results (1-10000) |

---

## Ingest

### `POST /api/v1/ingest`

**Auth:** analyst+
**Rate limit:** 20/min
**Max body:** 100 MB

Upload collector JSON output. Runs the full pipeline: validate, normalize, deduplicate, write, post-process.

```json
// Request body: collector JSON output (see graph-model.md for schema)

// Response (200)
{
  "scan_id": "scan-abc123",
  "nodes_written": 47,
  "edges_written": 82,
  "duration": "1.23s",
  "warnings": []
}
```

---

## Analysis

### `POST /api/v1/analysis/shortest-path`

**Auth:** analyst+

Find the shortest path between two nodes.

```json
// Request
{
  "source": "my-agent",
  "source_kind": "AgentInstance",
  "target": "postgres://prod",
  "target_kind": "MCPResource",
  "max_hops": 10
}

// Response
{
  "paths": [
    {
      "nodes": [
        { "id": "sha256:...", "name": "my-agent", "kinds": ["AgentInstance"] }
      ],
      "edges": [
        { "kind": "TRUSTS_SERVER", "source": "sha256:...", "target": "sha256:..." }
      ],
      "hops": 3
    }
  ]
}
```

### `POST /api/v1/analysis/all-paths`

**Auth:** analyst+

Enumerate all paths between two nodes (bounded).

Same request format as shortest-path, plus:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `limit` | int | 10 | Max paths returned (1-100) |

### `POST /api/v1/analysis/weighted-path`

**Auth:** analyst+

Find the lowest-risk-weight path using Dijkstra (APOC) or shortestPath+reduce fallback.

Same request format as shortest-path. Response includes `"algorithm": "dijkstra"` or `"algorithm": "shortestPath+reduce"`.

### `GET /api/v1/analysis/findings`

**Auth:** viewer+

List all composite edges as security findings.

| Param | Type | Description |
|-------|------|-------------|
| `severity` | string | Filter: `critical`, `high`, `medium`, `low` |

### `GET /api/v1/analysis/prebuilt`

**Auth:** viewer+

List all 17 pre-built queries with metadata.

### `GET /api/v1/analysis/prebuilt/{id}`

**Auth:** viewer+

Execute a pre-built query and return results.

```json
{
  "query": {
    "id": "agents-shell-access",
    "name": "Agents with Shell Access",
    "severity": "critical",
    "category": "Critical Paths",
    "owasp_map": ["MCP01", "ASI06"]
  },
  "rows": [...]
}
```

---

## Query

### `POST /api/v1/query`

**Auth:** admin
**Rate limit:** 10/min

Execute raw Cypher against Neo4j.

```json
// Request
{
  "query": "MATCH (n:MCPServer) RETURN n.name LIMIT 10"
}
```

---

## Scans

### `GET /api/v1/scans`

**Auth:** viewer+

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `limit` | int | 50 | Max results |
| `offset` | int | 0 | Pagination offset |

### `GET /api/v1/scans/{id}`

**Auth:** viewer+

Get scan details by ID.

---

## Audit

### `GET /api/v1/audit`

**Auth:** admin

View audit log entries (login events, ingest operations, query executions, token management).
