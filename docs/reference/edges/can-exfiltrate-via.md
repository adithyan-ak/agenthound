# CAN_EXFILTRATE_VIA

**Type:** Composite (post-processor generated)  
**Direction:** `AgentInstance → MCPTool`  
**Depends on:** CAN_REACH  
**Severity:** Critical  
**OWASP:** MCP04 (Tool Misuse), ASI03 (Data Exfiltration)

## What it means

An agent can reach sensitive data (via CAN_REACH to a high-sensitivity MCPResource) AND has access to a tool with an outbound capability (network_outbound, email_send). The combination means the agent can read secrets and send them out.

## How it's computed

The `can_exfiltrate` post-processor runs after CAN_REACH and checks:

1. Does the agent have a CAN_REACH path to a resource with `sensitivity` in {critical, high}?
2. Does the same agent (via TRUSTS_SERVER → PROVIDES_TOOL) have access to a tool whose `capability_surface` includes `network_outbound` OR `email_send`?

If both conditions hold, emit: `(agent)-[:CAN_EXFILTRATE_VIA]->(tool)`

## Cypher example

```cypher
MATCH (a:AgentInstance)-[:CAN_EXFILTRATE_VIA]->(t:MCPTool)
RETURN a.name AS agent, t.name AS exfil_tool,
       t.capability_surface AS capabilities
ORDER BY a.name
```

## What an operator does with it

This is the "game over" edge — the agent can steal data autonomously. Prioritize:
1. Remove the outbound tool from the agent's trusted servers
2. If the tool must stay, restrict its scope (input validation, allowlisted destinations)
3. Monitor: any invocation of the exfil tool by this agent is an incident

## Properties

| Property | Type | Description |
|----------|------|-------------|
| `confidence` | float | 0.0–1.0 |
| `risk_weight` | float | 0.1 (high exploitability) |
| `evidence` | object | `{sensitive_resource, outbound_tool, path_hops}` |
| `is_composite` | bool | true |
