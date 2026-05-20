# DEF CON 35 — main-stage CFP abstract

> **Status:** Finalized. v0.5 shipped all required capabilities (scan, discover, loot, poison, revert, extract). Demo validated end-to-end against live Docker lab.
> **Submission target:** DEF CON 35 main stage CFP (typically opens Feb 2027 for August 2027).
> **Track fit:** AppSec / Cloud / AI Village crossover. Main stage preferred — the live demo lands strongest at room scale.

---

## Title

**Hounds in the Stack: Mapping and Hijacking Attack Paths Across MCP, A2A, and AI-Service Infrastructure**

(working alternates: "BloodHound for Agents"; "From Discovery to Detonation: A Live Walk Through an AI-Agent Attack Chain")

---

## Abstract — 250 words

Modern AI deployments are not single hosts. A typical setup runs Ollama on a GPU box, vLLM on another, an LLM gateway like LiteLLM with virtual keys for forty engineers, an Open WebUI fronting it all, MCP servers exposing internal tools to the agents, A2A agents delegating work to peers, and a Jupyter Server somewhere that the data team forgot about. Each one of these components has its own auth posture, its own credential surface, and its own way of pointing at the next.

Existing tooling is single-protocol. There is no tool that maps these together — that says, in a single graph, "this Ollama is exposed by an unauthenticated Open WebUI on port 3000, which speaks to an MCP server that has shell-execution tools, which is referenced by a CLAUDE.md the engineer committed to GitHub, which means an Anthropic API key derived from a virtual key on the LiteLLM gateway can ultimately reach a production database."

We built one. **AgentHound** is BloodHound for AI-agent infrastructure: a Go collector + Neo4j-backed analysis server that discovers MCP, A2A, and AI services on a network, fingerprints them, extracts their credential graph, and uses shortest-path algorithms to surface end-to-end attack paths a single-protocol scanner cannot see.

This talk walks through one of those paths live — discovery, fingerprint, modelfile leak, gateway-key chain, MCP tool poison, revert — in the open-source release operators can run themselves. We close on what the agentic-infrastructure red-team toolbox needs next.

---

## Speaker bio (template — fill at submission)

`{{NAME}}` builds offensive security tooling for AI-agent infrastructure. Author of AgentHound (open-source). Previous work: `{{prior CVEs / talks}}`. `{{social handle}}`.

---

## Detailed outline (45 minutes)

| Time | Section | Content |
|---|---|---|
| 0:00–0:03 | Hook | Live: spin up the 8-host lab. "There are eight services here. Name the attack path." |
| 0:03–0:10 | The agentic-infra threat model | MCP, A2A, AI services. Why they break BloodHound's host/user/group model. The credential-chain that ties them together. |
| 0:10–0:18 | Discovery and fingerprint | Live: `agenthound scan 10.0.0.0/24` and `agenthound discover 10.0.0.0/24`. Watch the graph build. |
| 0:18–0:25 | Loot — modelfile, master key, upstream provider keys | Live: Ollama Looter (system-prompt leak from a fine-tune); LiteLLM Looter (master → upstream chain). |
| 0:25–0:32 | Cross-protocol pathfinding | Cypher walks: agent → MCP tool → host shell; agent → LiteLLM → upstream provider key. The `value_hash` cross-collector merge. |
| 0:32–0:40 | Poison — and revert | Live: poison an MCP tool description. Watch a live agent call the poisoned tool. `agenthound revert <engagement-id>`. Cypher confirms rollback. |
| 0:40–0:43 | What's next | Implanter, Extractor, the 8-VM lab, where contributions land. |
| 0:43–0:45 | Q&A intake / demo loop | Replay any of the live segments on demand. |

---

## Demos

1. **Discover + scan + fingerprint.** End-to-end against a private CIDR. Audience sees the graph.
2. **Modelfile leak.** `agenthound loot --type ollama` against an Ollama with a planted system prompt.
3. **Credential chain.** LiteLLM loot fires the `cross_service_credential_chain` post-processor; Cypher in the Explorer surfaces the path.
4. **Tool poison + revert.** The apex. Live MCP server, live tool description rewrite, live agent invokes the poisoned tool, live revert. Receipts on disk; rollback is byte-identical.

All demos are local — no internet egress required during the talk. Lab is `docker compose up`.

---

## Why this audience, why this venue

DEF CON main stage attendees include the cohort building, defending, and breaking AI infrastructure today. The chain we walk through is reproducible, the tooling is open-source under Apache 2.0, and the threat model maps cleanly to OWASP MCP Top 10 + ASI Top 10 — both of which are about to be referenced in published agent-security policy. This is the operator-facing toolset that puts those threat models in their hands.

---

## Risks and mitigations

- **Network demo failures.** Lab is fully local; no flag day if conference Wi-Fi is hostile.
- **Live agent unpredictability.** The poison demo uses a deterministic agent harness (the receipts capture original content; revert is byte-identical regardless of what the agent did).
- **Audience expectations of CVE drops.** This talk is tooling + methodology, not vuln drop. We frame it that way in the abstract above.

---

## Out-of-scope (do NOT promise in the talk)

- Breaking specific commercial agent platforms by name.
- Demonstrating zero-days against specific MCP server implementations (use generic patterns).
- Anything past v0.4 — Implanter / Extractor land but won't be the headliner.

---

## See also

- AgentHound v0.3.0 release notes (CHANGELOG.md)
- `docs/plans/v0.3-v0.4-implementation.md` — the multi-release plan this talk is built on
- `docs/architecture.md` — full system architecture for the technical deep-dive backup
