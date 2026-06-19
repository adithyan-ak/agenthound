-- 003_findings_triage.sql
--
-- Tier-0 productization: persist a per-scan snapshot of findings so the
-- graph (which the post-processor's stale-edge cleanup rewrites every scan)
-- stays diffable, and track cross-scan triage state keyed by the stable
-- 16-char finding fingerprint.
--
-- Two tables, deliberately separate:
--   * findings        — immutable per-scan record (FK to scans, ON DELETE
--                       CASCADE so deleting a scan reaps its snapshot).
--   * finding_triage  — cross-scan analyst state. NO FK to findings: triage
--                       decisions (accepted-risk / false-positive) must
--                       survive scan deletion and re-detection.

CREATE TABLE IF NOT EXISTS findings (
    scan_id        TEXT NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    fingerprint    CHAR(16) NOT NULL,
    severity       TEXT NOT NULL DEFAULT '',
    category       TEXT NOT NULL DEFAULT '',
    title          TEXT NOT NULL DEFAULT '',
    description    TEXT NOT NULL DEFAULT '',
    edge_kind      TEXT NOT NULL DEFAULT '',
    source_id      TEXT NOT NULL DEFAULT '',
    source_name    TEXT NOT NULL DEFAULT '',
    source_kind    TEXT NOT NULL DEFAULT '',
    target_id      TEXT NOT NULL DEFAULT '',
    target_name    TEXT NOT NULL DEFAULT '',
    target_kind    TEXT NOT NULL DEFAULT '',
    confidence     DOUBLE PRECISION NOT NULL DEFAULT 0,
    owasp_map      JSONB NOT NULL DEFAULT '[]',
    cross_protocol BOOLEAN NOT NULL DEFAULT FALSE,
    captured_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (scan_id, fingerprint)
);

CREATE INDEX IF NOT EXISTS idx_findings_fingerprint ON findings(fingerprint);
CREATE INDEX IF NOT EXISTS idx_findings_scan_id ON findings(scan_id);
CREATE INDEX IF NOT EXISTS idx_findings_severity ON findings(severity);
CREATE INDEX IF NOT EXISTS idx_findings_edge_kind ON findings(edge_kind);

CREATE TABLE IF NOT EXISTS finding_triage (
    fingerprint CHAR(16) PRIMARY KEY,
    status      TEXT NOT NULL DEFAULT 'new'
                CHECK (status IN ('new','triaging','confirmed','accepted-risk','false-positive')),
    note        TEXT NOT NULL DEFAULT '',
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
