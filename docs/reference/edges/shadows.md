# SHADOWS

**Type:** Composite (post-processor generated)  
**Direction:** `MCPTool → MCPTool`  
**Depends on:** Raw edges only  
**Severity:** High  
**OWASP:** MCP02 (Tool Description Manipulation), ASI06 (Tool Shadowing)

## What it means

One MCP tool's description references another tool by name or describes overlapping functionality — potentially tricking an agent into calling the shadowing tool instead of the legitimate one. This is the "tool shadowing" attack: a malicious server registers a tool with a description that mimics a trusted tool, causing the agent's planner to route requests to the attacker's implementation.

## How it's computed

The `shadows` post-processor uses TF-IDF cosine similarity on tool descriptions:

1. For each pair of tools from DIFFERENT servers (same-server tools can't shadow each other — the agent already trusts both)
2. Compute cosine similarity on the description text
3. If similarity > 0.8 AND one tool's description contains the other tool's name as a substring → emit SHADOWS edge

Additionally, explicit cross-references are detected: if tool A's description contains `tool_name: "B"` or references tool B's server by endpoint, that's a direct shadow signal.

## Cypher example

```cypher
MATCH (shadow:MCPTool)-[:SHADOWS]->(legit:MCPTool)
MATCH (shadow)<-[:PROVIDES_TOOL]-(evil:MCPServer)
MATCH (legit)<-[:PROVIDES_TOOL]-(good:MCPServer)
RETURN shadow.name AS shadowing_tool, evil.name AS malicious_server,
       legit.name AS legitimate_tool, good.name AS trusted_server
```

## What an operator does with it

Tool shadowing is a supply-chain attack on agent behavior:
1. Identify which agent trusts the shadowing server (via TRUSTS_SERVER)
2. Check if the shadowing server was added recently (supply-chain compromise)
3. Compare the two tool descriptions side-by-side — is the shadow an exact copy or a subtle modification?
4. Remediate: remove the malicious server from the agent's config, or pin the trusted tool by server+name

## Properties

| Property | Type | Description |
|----------|------|-------------|
| `confidence` | float | Cosine similarity score (0.8–1.0) |
| `risk_weight` | float | 0.4 |
| `evidence` | object | `{similarity_score, cross_reference_detected, description_hash_shadow, description_hash_legit}` |
| `is_composite` | bool | true |
