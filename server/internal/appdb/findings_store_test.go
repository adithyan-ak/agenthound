package appdb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/adithyan-ak/agenthound/server/model"
)

// TestIntegrationFindingStore exercises the FindingStore against a real
// Postgres (skipped without AGENTHOUND_PG_URI; CI's test-integration job
// wires it). It is the regression guard for the ListLatestPerFingerprint
// outer-SELECT alias bug, where the triage columns referenced the inner
// subquery alias `t` from the outer query and produced
// "missing FROM-clause entry for table t" on every findings read.
func TestIntegrationFindingStore(t *testing.T) {
	skipIfNoPG(t)
	ctx := context.Background()

	pool, err := NewPool(os.Getenv("AGENTHOUND_PG_URI"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := RunMigrations(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	scanStore := NewScanStore(pool)
	fs := NewFindingStore(pool)

	fpA := "aaaaaaaaaaaaaaaa"
	fpB := "bbbbbbbbbbbbbbbb"
	fpD := "dddddddddddddddd"

	// Make the test hermetic regardless of leftover state: clear any prior
	// fs-test scans (cascade-deletes their findings) and the fixed triage
	// fingerprints up front. finding_triage has no FK, so it survives scan
	// deletion and must be cleared explicitly.
	cleanup := func() {
		_, _ = pool.Exec(ctx, "DELETE FROM scans WHERE id LIKE 'fs-test-%'")
		_, _ = pool.Exec(ctx, "DELETE FROM finding_triage WHERE fingerprint = ANY($1)", []string{fpA, fpB, fpD})
	}
	cleanup()
	// Use defer (not t.Cleanup) so this runs BEFORE the deferred pool.Close()
	// above (defers are LIFO) — the t.Cleanup variant would run after the
	// pool is already closed and silently no-op.
	defer cleanup()

	scanID := "fs-test-" + time.Now().Format("20060102150405.000000")
	scanID2 := scanID + "-2"
	mustScan := func(id string) {
		if err := scanStore.CreateScan(ctx, &model.Scan{
			ID: id, Collector: "mcp", Status: model.ScanStatusRunning, StartedAt: time.Now().UTC(),
		}); err != nil {
			t.Fatalf("create scan %s: %v", id, err)
		}
	}
	mustScan(scanID)
	mustScan(scanID2)

	findings := []model.Finding{
		{ID: fpA, Severity: "critical", Category: "Data Exfiltration", Title: "exfil", EdgeKind: "CAN_EXFILTRATE_VIA",
			SourceID: "s1", SourceName: "agent", SourceKind: "AgentInstance", TargetID: "t1", TargetName: "tool", TargetKind: "MCPTool",
			Confidence: 0.9, OWASPMap: []string{"MCP04", "ASI08"}},
		{ID: fpB, Severity: "high", Category: "Tool Shadowing", Title: "shadow", EdgeKind: "SHADOWS",
			SourceKind: "MCPTool", TargetKind: "MCPTool", Confidence: 0.6},
	}
	if err := fs.InsertFindings(ctx, scanID, findings); err != nil {
		t.Fatalf("insert findings: %v", err)
	}

	// The load-bearing assertion: this read regressed to a 500 before the
	// alias fix.
	all, err := fs.ListLatestPerFingerprint(ctx, "", false)
	if err != nil {
		t.Fatalf("ListLatestPerFingerprint: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(all))
	}

	crit, err := fs.ListLatestPerFingerprint(ctx, "critical", false)
	if err != nil {
		t.Fatalf("list critical: %v", err)
	}
	if len(crit) != 1 || crit[0].Severity != "critical" {
		t.Fatalf("severity filter: got %+v", crit)
	}

	// Suppression: a false-positive is hidden by default, shown with the flag.
	if _, err := fs.UpsertTriage(ctx, fpA, "false-positive", "benign"); err != nil {
		t.Fatalf("upsert triage: %v", err)
	}
	visible, err := fs.ListLatestPerFingerprint(ctx, "", false)
	if err != nil {
		t.Fatalf("list after suppress: %v", err)
	}
	if len(visible) != 1 {
		t.Fatalf("suppressed finding should be hidden by default; got %d visible", len(visible))
	}
	withSupp, err := fs.ListLatestPerFingerprint(ctx, "", true)
	if err != nil {
		t.Fatalf("list include-suppressed: %v", err)
	}
	if len(withSupp) != 2 {
		t.Fatalf("include_suppressed should show all; got %d", len(withSupp))
	}
	var sawInlineTriage bool
	for _, f := range withSupp {
		if f.ID == fpA {
			if f.Triage == nil || f.Triage.Status != "false-positive" || f.Triage.Note != "benign" {
				t.Fatalf("expected inline triage {false-positive, benign}, got %+v", f.Triage)
			}
			sawInlineTriage = true
		}
	}
	if !sawInlineTriage {
		t.Fatal("suppressed finding not returned with include_suppressed")
	}

	// GetTriage: present and absent.
	ts, err := fs.GetTriage(ctx, fpA)
	if err != nil || ts == nil || ts.Status != "false-positive" {
		t.Fatalf("GetTriage(fpA): ts=%+v err=%v", ts, err)
	}
	if none, err := fs.GetTriage(ctx, "cccccccccccccccc"); err != nil || none != nil {
		t.Fatalf("GetTriage(unknown): want (nil,nil), got (%+v,%v)", none, err)
	}

	// Diff: scan2 keeps fpA, drops fpB, adds fpD.
	findings2 := []model.Finding{
		findings[0],
		{ID: fpD, Severity: "high", Category: "Cross-Tool Taint", Title: "taint", EdgeKind: "TAINTS",
			SourceKind: "MCPTool", TargetKind: "MCPTool", Confidence: 0.7},
	}
	if err := fs.InsertFindings(ctx, scanID2, findings2); err != nil {
		t.Fatalf("insert findings2: %v", err)
	}
	diff, err := fs.Diff(ctx, scanID, scanID2, false)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if len(diff.Added) != 1 || diff.Added[0].ID != fpD {
		t.Fatalf("diff.Added = %+v, want [%s]", diff.Added, fpD)
	}
	if len(diff.Removed) != 1 || diff.Removed[0].ID != fpB {
		t.Fatalf("diff.Removed = %+v, want [%s]", diff.Removed, fpB)
	}
	if len(diff.Unchanged) != 1 || diff.Unchanged[0].ID != fpA {
		t.Fatalf("diff.Unchanged = %+v, want [%s]", diff.Unchanged, fpA)
	}
}
