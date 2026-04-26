# AgentHound fingerprint templates

This directory will hold YAML fingerprint templates following the project's
matcher engine in `sdk/rules`. The first concrete template lands with the first
non-MCP module PR.

Until then the directory is intentionally empty. Templates are NOT shipped or
embedded; their location and packaging are part of `docs/future-modules.md`.

## Adding a fingerprint template

1. Define an `agenthound.yaml` matching the schema in `sdk/rules/`.
2. Add it under `templates/fingerprint/<service>/<service>.yaml`.
3. Add a Go test in `modules/<service>/` that loads and exercises it.

For the planned full DSL (status / header / json_path / extractors), see
`docs/future-modules.md`.
