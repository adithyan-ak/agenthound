# Security Policy

## Reporting a vulnerability

If you discover a security vulnerability in AgentHound, please report it responsibly through [GitHub Security Advisories](https://github.com/adithyan-ak/agenthound/security/advisories/new).

**Do not** open a public GitHub issue for security vulnerabilities.

### What to include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response timeline

- **Acknowledgment:** within 48 hours
- **Initial assessment:** within 7 days
- **Fix timeline:** depends on severity, typically within 30 days for critical issues

### Scope

The following are in scope:

- Authentication and authorization bypasses
- Injection vulnerabilities (Cypher injection, command injection)
- Data exposure (credentials, tokens, PII leaking through API responses or logs)
- Supply chain issues in AgentHound's own dependencies
- Container escape or privilege escalation in the Docker deployment

The following are out of scope:

- Vulnerabilities in MCP servers or A2A agents that AgentHound scans (report those to their respective maintainers)
- Denial of service through intentionally large graph data (the tool is designed for trusted operator use)
- Issues requiring physical access to the host

## Security design

AgentHound handles sensitive data (credentials, infrastructure topology, attack paths). Key security measures:

- **Credential hashing:** Config Collector hashes credential values by default (SHA-256). Raw values require explicit `--include-credential-values` flag.
- **Authentication:** bcrypt password hashing, JWT with HMAC-SHA256, API tokens with secure random generation
- **Authorization:** Role-based access control (admin/analyst/viewer) on all API endpoints
- **Audit logging:** All security-relevant actions are logged with actor, action, resource, and timestamp
- **Rate limiting:** API endpoints are rate-limited to prevent abuse
- **Input validation:** All API inputs are validated. Cypher queries are admin-only.
- **Container security:** Non-root user, minimal base image (alpine:3.19), no unnecessary packages

## Supported versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | Yes       |
