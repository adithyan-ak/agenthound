---
paths:
  - "docs/**"
  - "mkdocs.yml"
  - ".github/workflows/docs.yml"
---
# Docs Site Rules (MkDocs + Material → GitHub Pages → docs.agenthound.io)

- Public docs are built by **MkDocs + Material** from `docs/`. `mkdocs.yml` lives at repo root; toolchain pinned in `docs/requirements.txt` (`mkdocs-material==9.*`).
- The site builds with **`mkdocs build --strict`** in CI (`.github/workflows/docs.yml`). Strict mode FAILS on any broken internal link, missing anchor, orphan page, or absolute link. A broken link is a red build — fix the link, don't relax the validation.
- **Adding a new doc page:** create the `.md` under `docs/`, then add it to the `nav:` in `mkdocs.yml`. A page not in `nav` triggers an `omitted_files` warning → strict failure. Verify locally with `mkdocs build --strict` before pushing.
- **Internal links:** use RELATIVE `.md` paths (e.g. `../operator/security.md`). MkDocs rewrites them to URLs. Never absolute paths. Never link into `cfp/` or `plans/`.
- `docs/README.md` is the site homepage (MkDocs treats it as index). Keep it as the nav hub; do NOT add a `docs/index.md` (it would shadow the README).
- **Excluded from the public site:** `cfp/` and `plans/` via `exclude_docs` in `mkdocs.yml` (internal strategy). Don't nav-link them. (Squash/relocate later per CLAUDE.md.)
- **Pages source** is "GitHub Actions" (Settings → Pages), NOT a `gh-pages` branch. Do NOT commit a `CNAME` file — the custom domain `docs.agenthound.io` is configured in repo Settings only; a committed CNAME is ignored by the Actions flow and just rots.
- **DNS:** `docs.agenthound.io` is a subdomain → single `CNAME` record `docs` → `adithyan-ak.github.io`. (Subdomains use CNAME, not the apex A-records used for a bare domain.)
- `site/` is the build artifact — gitignored, never committed.
- GitHub Actions are SHA-pinned (dependabot-managed) to match the other workflows; keep new doc-workflow actions pinned with a `# vX.Y.Z` comment.
- Local preview: `pip install -r docs/requirements.txt && mkdocs serve`.
