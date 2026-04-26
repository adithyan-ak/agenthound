# Future modules — deferred surface and planning notes

This file documents work that is **deliberately not shipped** in the current
release but is informed by design choices already made. Treat it as a
roadmap for whoever picks the framework up next, not a commitment to ship.

## How to add a new module today

The collector is built around a self-registration pattern. To add a new
module:

1. Create `modules/<name>/`.
2. Implement an action interface from `sdk/action/`. The current
   shipped action is `Enumerator`. Future actions (`Fingerprinter`,
   `Looter`, `Extractor`, `Poisoner`, `Implanter`) are declared but not
   yet driven by any caller.
3. Add `modules/<name>/register.go` with an `init()` that calls
   `sdk/module.Register(...)` to expose the module to the registry.
4. Blank-import `_ "github.com/adithyan-ak/agenthound/modules/<name>"`
   in `collector/cmd/agenthound/main.go`.
5. Wire the new module into the `scan` CLI verb (or add a new verb).

That's the contract. See `modules/README.md` for the cleanest existing
example.

## Future shape commitments

These are areas where the current API is intentionally narrow because the
right shape isn't known yet. Expect them to evolve before 1.0.

### `sdk/action.Target`

`Target` is currently a flat struct. By 1.0 it will likely grow typed
sub-structs for different target classes (host, MCP server, A2A agent,
config file). Module authors should not hard-code field positions.

### `sdk/action.Action.IsDestructive() bool`

A boolean is sufficient for shipping today's collectors (none are
destructive). Once `Poisoner` and `Implanter` modules are real, a
`SideEffects` enum is more honest:

```go
type SideEffects int
const (
    SideEffectsNone SideEffects = iota
    SideEffectsReadOnly
    SideEffectsModifiesTarget
    SideEffectsPersistsImplant
)
```

The migration is straightforward: `IsDestructive() bool` becomes a method
that returns `e.SideEffects() >= SideEffectsModifiesTarget`.

### `sdk/ingest.IngestMeta.Extensions`

`Extensions` is a `map[string]any` today. It is wire-format-only — the
validator does not type-check it. The first new node kind that needs
structured metadata will introduce per-extension schemas in
`sdk/ingest/extensions/`. Until then, prefer typed top-level fields over
stuffing into `Extensions`.

## Future Nuclei-pattern community ecosystem

The architect review (C1–C5) flagged a future direction worth recording
here so the design doesn't drift away from it.

### Templates likely live in a separate repo

Nuclei's lesson: bundling templates and engine in the same repo causes
release-cadence friction. AgentHound templates (fingerprint patterns,
poison patterns, instruction-poisoning indicators) should live at
`adithyan-ak/agenthound-templates` once the matcher is stable. The
engine (this repo) loads templates by reference, not by inclusion.

### `sdk/rules.MatcherSpec` extensions

Today's matcher supports keyword, prefix, regex, entropy, and composite
matchers. To support fingerprinting, add:

- `http_status` — match HTTP response status code
- `http_header` — match a named header against a regex
- `json_path` — extract a JSONPath value and match it
- `extractor` block — emit captured groups for downstream rules
- structured emit — turn a match into a node/edge with typed fields

These extensions should land behind a `MatcherSpec.Version: 2` boundary
so existing rules don't silently re-interpret.

### `agenthound template` CLI

Three subcommands for the future template ecosystem:

- `agenthound template validate <path>` — schema-check a template file.
- `agenthound template test <path> --against <fixture>` — run the
  template's matchers against a fixture and emit the resulting nodes/edges.
- `agenthound template sign <path> --key <ecdsa-key>` — sign a template
  for distribution.

### Template signing

**This is non-negotiable for any module that performs poison or implant
actions.** Nuclei CVE-2024-43405 (signature bypass via YAML newline
injection) is the cautionary tale. AgentHound's destructive surface
(when implemented) is broader than Nuclei's, so the signing bar is
higher:

- ECDSA P-256 signatures, never just hashes.
- Sign over the canonical-encoded YAML, not the raw bytes.
- Verify on load, not on use, so a tampered template fails fast.
- Maintain a community keyring at `adithyan-ak/agenthound-templates/.keyring`.

### `-update-templates` flag

Once templates are externalized, the collector should be able to fetch
the latest template bundle on demand. Default off (red-team OPSEC); opt
in via `-update-templates` or `agenthound template update`.

## Operator-PII redaction

`--redact-paths <regex>` is a known-future addition to the collector.
The use case: an operator collecting on a third-party customer's host
where filesystem paths embed PII. Apply at the normalizer step, before
the JSON ever lands on disk or wire.

## Revert paths for destructive modules

The `Reverter` interface lives in `sdk/action/` from day one. Any module
that satisfies `Poisoner` or `Implanter` MUST also satisfy `Reverter`.
This is enforced at module registration time.

The `Reverter` contract is "best-effort": it MUST be idempotent and
SHOULD log every action it takes, but it cannot guarantee revert success
on a host the operator has lost access to. That's a design constraint,
not a defect.

## What this file is NOT

This is not a backlog. Items here are recorded so future contributors
know which design choices are load-bearing and which are placeholders.
Nothing here is a commitment to ship.
