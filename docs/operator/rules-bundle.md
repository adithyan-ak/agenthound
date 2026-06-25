# `--rules-bundle` — out-of-band fingerprint rule updates

The fingerprint rules engine ships rules embedded in the AgentHound binary (`sdk/rules/builtin/fingerprints/*.yaml`). AgentHound provides a `--rules-bundle <path>` override so operators can pick up rule fixes without rebuilding the collector.

> **The operator is responsible for verifying the cosign signature on the bundle BEFORE pointing AgentHound at it.** Verification is optional by design — the loader does not call cosign automatically and does not refuse unsigned bundles. Always run `cosign verify-blob` against the tarball yourself (see below) before pointing AgentHound at it.

---

## Quick start

```bash
# Download a published bundle.
gh release download rules-v2026.06.01 \
    --repo adithyan-ak/agenthound \
    --pattern 'agenthound-rules-*.tar.gz*'

# Verify the cosign signature BEFORE running anything.
cosign verify-blob \
    --bundle agenthound-rules-rules-v2026.06.01.tar.gz.sigstore.json \
    --certificate-identity-regexp 'https://github.com/.*' \
    --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
    agenthound-rules-rules-v2026.06.01.tar.gz

# Use the bundle in any command that runs fingerprinters.
agenthound --rules-bundle ./agenthound-rules-rules-v2026.06.01.tar.gz scan 10.0.0.0/24

# Or point at a directory of YAML files (during development / lab work).
agenthound --rules-bundle ./my-custom-rules/ scan 10.0.0.0/24

# Env var alternative.
AGENTHOUND_RULES_BUNDLE=./bundle.tar.gz agenthound scan 10.0.0.0/24
```

---

## Override semantics

The bundle merges into the embedded rule set with **same-id rules from the bundle winning**:

| Embedded rule | Bundle rule | Effective rule |
|---|---|---|
| `id: ollama` | (absent) | embedded |
| `id: ollama` | `id: ollama` | bundle (override wins) |
| (absent) | `id: my-custom` | bundle (additive) |

This is what you want for hot-fixing a broken regex in a shipped rule, or for adding a new fingerprinter rule the binary's rule set doesn't yet include.

---

## Bundle format

A bundle is one of:

- A directory containing `*.yaml` files (one rule per file).
- A `.tar.gz` archive containing `fingerprints/*.yaml` entries.

The loader reads `*.yaml` and skips other file types. Each YAML file follows the same shape as the embedded rules at `sdk/rules/builtin/fingerprints/`:

```yaml
id: ollama-hotfix-2026-06
name: Ollama (CVE-2026-XXXXX hotfix)
description: refines the Ollama version regex to catch the new patch series
version: 2
service_kind: ollama
probes:
  - method: GET
    path: /api/version
    matchers:
      - type: http_status
        status_code: 200
      - type: json_path
        path: "$.version"
        regex: '^\d+\.\d+\.\d+(-rc\d+)?$'
    captures:
      version: "$.version"
emit:
  node_kinds:
    - OllamaInstance
    - AIService
  properties:
    service_kind: ollama
    auth_method: none
    is_anonymous_loot: "true"
    version: "{capture:version}"
```

Rule IDs MUST be globally unique within a bundle (and across the bundle + embedded merge — the bundle's same-id rule wins, but two same-id rules within one bundle is a load-time conflict).

---

## Release cadence

Bundles are published by the `rules-bundle.yml` GitHub Actions workflow. Triggers:

- **`workflow_dispatch`** — manual release, used for ad-hoc rule fixes.
- **`on: push: tags: ['rules-v*']`** — pushing a `rules-vYYYY.MM.DD` tag automatically cuts a release.

There is **no `on: schedule` trigger**. Bundles are content-driven — a no-changes month produces no bundle. An empty release would confuse cosign verification on the consumer side.

The release artifacts:

- `agenthound-rules-<tag>.tar.gz` — the bundle.
- `agenthound-rules-<tag>.tar.gz.sha256` — checksum.
- `agenthound-rules-<tag>.tar.gz.sigstore.json` — cosign keyless bundle (signature + certificate, cosign v3 format).

---

## Troubleshooting

**Bundle doesn't load at all.** Check the path. `agenthound --rules-bundle <path>` surfaces the error from `LoadFingerprintBundle` if the path doesn't exist or doesn't unpack. Check the format with `tar -tzf <bundle>.tar.gz` showing one or more `*.yaml` entries.

**Bundle loads but my override doesn't take effect.** Same-id-wins requires the bundle's rule ID to match the embedded rule ID exactly. Check the embedded set with `agenthound rules list`.

**Bundle loads but a rule is silently dropped.** The loader skips files that fail YAML parsing. Run `agenthound rules validate <yaml-path>` against each file in your bundle to catch parse errors. Per-file size cap is 1 MiB.

**Cosign verification fails.** The `--certificate-identity-regexp` in the example matches GitHub Actions OIDC. If you forked the repo and re-released, your cert identity will differ — adjust the regex. Time skew can cause verification failures if your machine clock is off; sync NTP.

---

## See also

- `sdk/rules/bundle.go` — implementation of `LoadFingerprintBundle` and `MergeFingerprintRules`.
- `.github/workflows/rules-bundle.yml` — the release pipeline.
- [`agenthound scan` operator guide](scanner.md) — the network scanner that consumes fingerprint rules.
