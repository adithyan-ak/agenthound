Agenthound: Mapping Multi-Hop Credential Chains Across MCP, A2A, and LLM Gateways
Description
An MCP env var leaks the LiteLLM master key. The master key fronts your Anthropic production key. A teammate's Cursor config references a CLAUDE.md committed to GitHub. No single config file declares this chain but in a graph, it's a five-hop path you can walk in 200 milliseconds. Modern AI deployments aren't a single host. Ollama on a GPU box, vLLM on another, a LiteLLM gateway holding aggregated provider keys for forty engineers, Open WebUI fronting it, MCP servers exposing internal tools to agents, A2A agents delegating to peers, a Jupyter the data team forgot about. Every component has its own auth posture and credential surface. None of them know about the next. Single-protocol AI-security tooling can't surface paths that cross MCP, A2A, and the AI services behind them.

AgentHound (github.com/adithyan-ak/agenthound) is BloodHound for AI agent infrastructure: a Go collector + Neo4j-backed analysis server that enumerates MCP servers, A2A agents, and AI services on a network, fingerprints them, extracts their credential graph, and runs shortest-path queries to surface end-to-end attack paths a single-protocol scanner cannot see. The hero detection is credential-chain CAN_REACH: multi-hop paths where Server A reads a credential, Server B uses that credential, and an agent reaches B's resources without trusting B directly. value_hash, a SHA-256 of the credential value populated by every Looter, is the cross-collector merge primitive that collapses the same secret across services into one node letting Cypher walk the chain.

This 30-minute workshop frames the thesis end-to-end:
- The threat model: why MCP/A2A/AI-service sprawl breaks the host/user/group model BloodHound was built on, mapped to OWASP MCP Top 10 (MCP03 - Tool Poisoning) and the OWASP Agentic AI Top 10 (ASI04 - Agentic Supply Chain Vulnerabilities).
- The graph schema: 23 collector node kinds, 30 edge kinds, 15 post-processors in dependency-validated order.
- The credential-chain detection demo: `make demo` starts a Docker lab, runs the collector from inside the lab network, ingests a config scan plus network scan, protocol discovery, LiteLLM loot, and Ollama loot through the HTTP API, then validates that the `litellm-credential-leak` path is present. The graph surfaces "agent → MCP env-var → LiteLLM master → upstream Anthropic key" as a single finding backed by `value_hash`.
- A 2-minute Extract cameo, framed as a pre-staged residue check: `agenthound extract --type embedding-invert` against a GGUF fine-tune. Pure-Go GGUF parser; statistical outliers on the embedding matrix surface vocabulary tokens as `:ExtractedTrainingSignal` nodes - a training-signal residue check, not a model-stealing primitive.
Open source under Apache 2.0. v0.7.0 is tagged.

What type of session is your submission	
Workshop and Tactic

Does your workshop come with a tactic	
Yes

How much time would you like for you session	
30 minutes for the workshop + 2 hours for the paired tactic (3 waves of 10–15 attendees, ~35 min activity per wave). Workshop kept short by design, the tactic is where the hands-on lives.

Submitted
31 May 2026 3:07 pm
