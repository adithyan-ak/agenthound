package appdb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adithyan-ak/agenthound/server/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// suppressedStatuses are triage decisions that hide a finding from the
// default findings view and from CI-style diffs/fail-on gates.
var suppressedStatuses = []string{"accepted-risk", "false-positive"}

// FindingStore persists per-scan finding snapshots and cross-scan triage
// state. The graph's stale-edge cleanup rewrites composite edges every
// scan, so the Postgres snapshot is the only diffable record of "what was
// found when".
type FindingStore struct {
	pool *pgxpool.Pool
}

func NewFindingStore(pool *pgxpool.Pool) *FindingStore {
	return &FindingStore{pool: pool}
}

// FindingsDiff is the result of comparing two scans' finding snapshots,
// keyed by fingerprint.
type FindingsDiff struct {
	ScanA     string          `json:"scan_a"`
	ScanB     string          `json:"scan_b"`
	Added     []model.Finding `json:"added"`
	Removed   []model.Finding `json:"removed"`
	Unchanged []model.Finding `json:"unchanged"`
}

// InsertFindings persists a scan's findings snapshot. Idempotent: a re-run
// of the same scan_id overwrites the prior rows for that scan.
func (s *FindingStore) InsertFindings(ctx context.Context, scanID string, findings []model.Finding) error {
	if scanID == "" {
		return errors.New("insert findings: empty scan_id")
	}
	if len(findings) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, f := range findings {
		owasp := f.OWASPMap
		if owasp == nil {
			owasp = []string{}
		}
		batch.Queue(
			`INSERT INTO findings
			   (scan_id, fingerprint, severity, category, title, description, edge_kind,
			    source_id, source_name, source_kind, target_id, target_name, target_kind,
			    confidence, owasp_map, cross_protocol)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
			 ON CONFLICT (scan_id, fingerprint) DO UPDATE SET
			    severity = EXCLUDED.severity,
			    category = EXCLUDED.category,
			    title = EXCLUDED.title,
			    description = EXCLUDED.description,
			    edge_kind = EXCLUDED.edge_kind,
			    source_id = EXCLUDED.source_id,
			    source_name = EXCLUDED.source_name,
			    source_kind = EXCLUDED.source_kind,
			    target_id = EXCLUDED.target_id,
			    target_name = EXCLUDED.target_name,
			    target_kind = EXCLUDED.target_kind,
			    confidence = EXCLUDED.confidence,
			    owasp_map = EXCLUDED.owasp_map,
			    cross_protocol = EXCLUDED.cross_protocol`,
			scanID, f.ID, f.Severity, f.Category, f.Title, f.Description, f.EdgeKind,
			f.SourceID, f.SourceName, f.SourceKind, f.TargetID, f.TargetName, f.TargetKind,
			f.Confidence, owasp, isCrossProtocol(f.SourceKind, f.TargetKind),
		)
	}

	br := s.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range findings {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("insert findings batch: %w", err)
		}
	}
	return nil
}

const findingSelectColumns = `f.fingerprint, f.severity, f.category, f.title, f.description, f.edge_kind,
	f.source_id, f.source_name, f.source_kind, f.target_id, f.target_name, f.target_kind,
	f.confidence, f.owasp_map, t.status, t.note, t.updated_at`

// ListLatestPerFingerprint returns the most recent finding row per
// fingerprint (across all scans), with triage state attached. severity
// filters by exact level when non-empty. When includeSuppressed is false,
// findings triaged as accepted-risk / false-positive are dropped.
func (s *FindingStore) ListLatestPerFingerprint(ctx context.Context, severity string, includeSuppressed bool) ([]model.Finding, error) {
	query := `
SELECT ` + strings.ReplaceAll(findingSelectColumns, "f.", "latest.") + `
FROM (
    SELECT DISTINCT ON (f.fingerprint) ` + findingSelectColumns + `
    FROM findings f
    LEFT JOIN finding_triage t ON t.fingerprint = f.fingerprint
    ORDER BY f.fingerprint, f.captured_at DESC
) latest
WHERE ($1 = '' OR latest.severity = $1)
  AND ($2 OR latest.status IS NULL OR latest.status NOT IN ('accepted-risk','false-positive'))
ORDER BY latest.confidence DESC`

	rows, err := s.pool.Query(ctx, query, severity, includeSuppressed)
	if err != nil {
		return nil, fmt.Errorf("list findings: %w", err)
	}
	defer rows.Close()
	return scanFindings(rows)
}

// findingsForScan returns every finding persisted for a single scan, with
// triage state attached.
func (s *FindingStore) findingsForScan(ctx context.Context, scanID string) ([]model.Finding, error) {
	query := `
SELECT ` + findingSelectColumns + `
FROM findings f
LEFT JOIN finding_triage t ON t.fingerprint = f.fingerprint
WHERE f.scan_id = $1
ORDER BY f.confidence DESC`

	rows, err := s.pool.Query(ctx, query, scanID)
	if err != nil {
		return nil, fmt.Errorf("findings for scan %s: %w", scanID, err)
	}
	defer rows.Close()
	return scanFindings(rows)
}

// Diff compares two scans' snapshots. added = present in scanB but not
// scanA; removed = present in scanA but not scanB; unchanged = present in
// both. When includeSuppressed is false, suppressed findings are dropped
// from the added set so CI-style diffs don't re-alert on accepted risks.
func (s *FindingStore) Diff(ctx context.Context, scanA, scanB string, includeSuppressed bool) (*FindingsDiff, error) {
	a, err := s.findingsForScan(ctx, scanA)
	if err != nil {
		return nil, err
	}
	b, err := s.findingsForScan(ctx, scanB)
	if err != nil {
		return nil, err
	}

	aByFP := make(map[string]model.Finding, len(a))
	for _, f := range a {
		aByFP[f.ID] = f
	}
	bByFP := make(map[string]model.Finding, len(b))
	for _, f := range b {
		bByFP[f.ID] = f
	}

	diff := &FindingsDiff{ScanA: scanA, ScanB: scanB}
	for _, f := range b {
		if _, ok := aByFP[f.ID]; ok {
			diff.Unchanged = append(diff.Unchanged, f)
			continue
		}
		if !includeSuppressed && isSuppressed(f.Triage) {
			continue
		}
		diff.Added = append(diff.Added, f)
	}
	for _, f := range a {
		if _, ok := bByFP[f.ID]; !ok {
			diff.Removed = append(diff.Removed, f)
		}
	}
	return diff, nil
}

// GetTriage returns the triage state for a fingerprint, or nil if none has
// been recorded (callers treat nil as the implicit "new" status).
func (s *FindingStore) GetTriage(ctx context.Context, fingerprint string) (*model.TriageState, error) {
	var ts model.TriageState
	err := s.pool.QueryRow(ctx,
		`SELECT status, note, updated_at FROM finding_triage WHERE fingerprint = $1`,
		fingerprint).Scan(&ts.Status, &ts.Note, &ts.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get triage: %w", err)
	}
	return &ts, nil
}

// UpsertTriage records (or updates) the triage decision for a fingerprint.
func (s *FindingStore) UpsertTriage(ctx context.Context, fingerprint, status, note string) (*model.TriageState, error) {
	var ts model.TriageState
	err := s.pool.QueryRow(ctx,
		`INSERT INTO finding_triage (fingerprint, status, note, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (fingerprint) DO UPDATE SET
		    status = EXCLUDED.status,
		    note = EXCLUDED.note,
		    updated_at = NOW()
		 RETURNING status, note, updated_at`,
		fingerprint, status, note).Scan(&ts.Status, &ts.Note, &ts.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert triage: %w", err)
	}
	return &ts, nil
}

// scanFindings maps result rows (with the LEFT JOIN triage columns) into
// model.Finding values. A NULL triage status yields a nil Triage pointer.
func scanFindings(rows pgx.Rows) ([]model.Finding, error) {
	var out []model.Finding
	for rows.Next() {
		var f model.Finding
		var status, note *string
		var updatedAt *time.Time
		if err := rows.Scan(
			&f.ID, &f.Severity, &f.Category, &f.Title, &f.Description, &f.EdgeKind,
			&f.SourceID, &f.SourceName, &f.SourceKind, &f.TargetID, &f.TargetName, &f.TargetKind,
			&f.Confidence, &f.OWASPMap, &status, &note, &updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan finding row: %w", err)
		}
		f.ID = strings.TrimSpace(f.ID)
		if status != nil {
			ts := &model.TriageState{Status: *status}
			if note != nil {
				ts.Note = *note
			}
			if updatedAt != nil {
				ts.UpdatedAt = *updatedAt
			}
			f.Triage = ts
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func isSuppressed(ts *model.TriageState) bool {
	if ts == nil {
		return false
	}
	for _, s := range suppressedStatuses {
		if ts.Status == s {
			return true
		}
	}
	return false
}

// isCrossProtocol mirrors the UI's cross-protocol predicate: a finding
// whose endpoints straddle the A2A and MCP protocol families.
func isCrossProtocol(sourceKind, targetKind string) bool {
	a2a := func(k string) bool { return strings.HasPrefix(k, "A2A") }
	mcp := func(k string) bool { return strings.HasPrefix(k, "MCP") }
	return (a2a(sourceKind) && mcp(targetKind)) || (mcp(sourceKind) && a2a(targetKind))
}
