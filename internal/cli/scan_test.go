package cli

import (
	"strings"
	"testing"

	"github.com/adithyan-ak/agenthound/internal/analysis"
	"github.com/adithyan-ak/agenthound/internal/model"
	"github.com/spf13/cobra"
)

func TestCountFindingsBySeverity(t *testing.T) {
	findings := []analysis.Finding{
		{ID: "1", Severity: "critical"},
		{ID: "2", Severity: "high"},
		{ID: "3", Severity: "high"},
		{ID: "4", Severity: "medium"},
		{ID: "5", Severity: "low"},
		{ID: "6", Severity: "low"},
		{ID: "7", Severity: "low"},
	}
	counts := countFindingsBySeverity(findings)

	want := map[string]int{"critical": 1, "high": 2, "medium": 1, "low": 3}
	for sev, expected := range want {
		if counts[sev] != expected {
			t.Errorf("severity %q: got %d, want %d", sev, counts[sev], expected)
		}
	}
}

func TestCountFindingsBySeverity_Empty(t *testing.T) {
	counts := countFindingsBySeverity(nil)
	for _, sev := range []string{"critical", "high", "medium", "low"} {
		if counts[sev] != 0 {
			t.Errorf("severity %q: got %d, want 0", sev, counts[sev])
		}
	}
}

func TestCountAtOrAbove(t *testing.T) {
	findings := []analysis.Finding{
		{ID: "1", Severity: "critical"},
		{ID: "2", Severity: "high"},
		{ID: "3", Severity: "high"},
		{ID: "4", Severity: "medium"},
		{ID: "5", Severity: "low"},
	}
	got := countAtOrAbove(findings, severityRank["high"])
	if got != 3 {
		t.Errorf("countAtOrAbove(threshold=high): got %d, want 3", got)
	}
}

func TestCountAtOrAbove_NoneAbove(t *testing.T) {
	findings := []analysis.Finding{
		{ID: "1", Severity: "medium"},
		{ID: "2", Severity: "medium"},
		{ID: "3", Severity: "low"},
	}
	got := countAtOrAbove(findings, severityRank["critical"])
	if got != 0 {
		t.Errorf("countAtOrAbove(threshold=critical): got %d, want 0", got)
	}
}

func TestCollectorDetail_Config(t *testing.T) {
	data := &model.IngestData{
		Graph: model.GraphData{
			Nodes: []model.Node{
				{ID: "1", Kinds: []string{"AgentInstance"}},
				{ID: "2", Kinds: []string{"AgentInstance"}},
				{ID: "3", Kinds: []string{"MCPServer"}},
				{ID: "4", Kinds: []string{"Credential"}},
				{ID: "5", Kinds: []string{"Credential"}},
				{ID: "6", Kinds: []string{"Credential"}},
			},
		},
	}
	got := collectorDetail(data, "config")
	want := "2 agents, 1 servers, 3 credentials"
	if got != want {
		t.Errorf("collectorDetail(config): got %q, want %q", got, want)
	}
}

func TestCollectorDetail_MCP(t *testing.T) {
	data := &model.IngestData{
		Graph: model.GraphData{
			Nodes: []model.Node{
				{ID: "1", Kinds: []string{"MCPTool"}},
				{ID: "2", Kinds: []string{"MCPTool"}},
				{ID: "3", Kinds: []string{"MCPResource"}},
				{ID: "4", Kinds: []string{"MCPPrompt"}},
			},
		},
	}
	got := collectorDetail(data, "mcp")
	want := "2 tools, 1 resources, 1 prompts"
	if got != want {
		t.Errorf("collectorDetail(mcp): got %q, want %q", got, want)
	}
}

func TestCollectorDetail_A2A(t *testing.T) {
	data := &model.IngestData{
		Graph: model.GraphData{
			Nodes: []model.Node{
				{ID: "1", Kinds: []string{"A2AAgent"}},
				{ID: "2", Kinds: []string{"A2ASkill"}},
				{ID: "3", Kinds: []string{"A2ASkill"}},
				{ID: "4", Kinds: []string{"A2ASkill"}},
			},
		},
	}
	got := collectorDetail(data, "a2a")
	want := "1 agents, 3 skills"
	if got != want {
		t.Errorf("collectorDetail(a2a): got %q, want %q", got, want)
	}
}

func TestCollectorDetail_Default(t *testing.T) {
	data := &model.IngestData{
		Graph: model.GraphData{
			Nodes: []model.Node{
				{ID: "1", Kinds: []string{"MCPServer"}},
				{ID: "2", Kinds: []string{"MCPTool"}},
			},
		},
	}
	got := collectorDetail(data, "unknown")
	want := "2 nodes"
	if got != want {
		t.Errorf("collectorDetail(unknown): got %q, want %q", got, want)
	}
}

func TestCollectorDetail_Empty(t *testing.T) {
	data := &model.IngestData{}
	got := collectorDetail(data, "config")
	want := "0 agents, 0 servers, 0 credentials"
	if got != want {
		t.Errorf("collectorDetail(config empty): got %q, want %q", got, want)
	}
}

func TestOrderResults(t *testing.T) {
	results := []scanResult{
		{name: "A2A", status: "ok"},
		{name: "MCP", status: "ok"},
		{name: "Config", status: "ok"},
	}
	ordered := orderResults(results)

	wantOrder := []string{"Config", "MCP", "A2A"}
	for i, want := range wantOrder {
		if ordered[i].name != want {
			t.Errorf("position %d: got %q, want %q", i, ordered[i].name, want)
		}
	}

	if results[0].name != "A2A" {
		t.Error("orderResults mutated the original slice")
	}
}

func TestOrderResults_SingleElement(t *testing.T) {
	results := []scanResult{
		{name: "MCP", status: "ok"},
	}
	ordered := orderResults(results)
	if len(ordered) != 1 || ordered[0].name != "MCP" {
		t.Errorf("single element: got %v", ordered)
	}
}

func TestOrderResults_AlreadyOrdered(t *testing.T) {
	results := []scanResult{
		{name: "Config", status: "ok"},
		{name: "MCP", status: "ok"},
		{name: "A2A", status: "ok"},
	}
	ordered := orderResults(results)
	wantOrder := []string{"Config", "MCP", "A2A"}
	for i, want := range wantOrder {
		if ordered[i].name != want {
			t.Errorf("position %d: got %q, want %q", i, ordered[i].name, want)
		}
	}
}

func TestCountAtOrAbove_AllLevels(t *testing.T) {
	findings := []analysis.Finding{
		{ID: "1", Severity: "critical"},
		{ID: "2", Severity: "high"},
		{ID: "3", Severity: "medium"},
		{ID: "4", Severity: "low"},
	}

	tests := []struct {
		threshold string
		want      int
	}{
		{"critical", 1},
		{"high", 2},
		{"medium", 3},
		{"low", 4},
	}

	for _, tt := range tests {
		t.Run(tt.threshold, func(t *testing.T) {
			got := countAtOrAbove(findings, severityRank[tt.threshold])
			if got != tt.want {
				t.Errorf("countAtOrAbove(%s) = %d, want %d", tt.threshold, got, tt.want)
			}
		})
	}
}

func TestRunScan_ConfigWithURL(t *testing.T) {
	cmd := &cobra.Command{RunE: runScan}
	cmd.Flags().Bool("config", false, "")
	cmd.Flags().Bool("mcp", false, "")
	cmd.Flags().Bool("a2a", false, "")
	cmd.Flags().String("path", "", "")
	cmd.Flags().StringSlice("paths", nil, "")
	cmd.Flags().String("project-dir", "", "")
	cmd.Flags().Bool("include-credential-values", false, "")
	cmd.Flags().String("url", "", "")
	cmd.Flags().String("target", "", "")
	cmd.Flags().StringSlice("targets", nil, "")
	cmd.Flags().StringSlice("discover-domain", nil, "")
	cmd.Flags().String("targets-file", "", "")
	cmd.Flags().String("auth-token", "", "")
	cmd.Flags().Int("concurrency", 5, "")
	cmd.Flags().Duration("timeout", 0, "")
	cmd.Flags().Bool("insecure", false, "")
	cmd.Flags().String("output", "", "")
	cmd.Flags().String("fail-on", "", "")
	_ = cmd.Flags().Set("config", "true")
	_ = cmd.Flags().Set("url", "http://example.com")
	err := runScan(cmd, nil)
	if err == nil {
		t.Fatal("expected error for --url with --config")
	}
	if !strings.Contains(err.Error(), "--url requires --mcp") {
		t.Errorf("error = %q, want '--url requires --mcp'", err.Error())
	}
}

func TestRunScan_A2ANoTarget(t *testing.T) {
	cmd := &cobra.Command{RunE: runScan}
	cmd.Flags().Bool("config", false, "")
	cmd.Flags().Bool("mcp", false, "")
	cmd.Flags().Bool("a2a", false, "")
	cmd.Flags().String("path", "", "")
	cmd.Flags().StringSlice("paths", nil, "")
	cmd.Flags().String("project-dir", "", "")
	cmd.Flags().Bool("include-credential-values", false, "")
	cmd.Flags().String("url", "", "")
	cmd.Flags().String("target", "", "")
	cmd.Flags().StringSlice("targets", nil, "")
	cmd.Flags().StringSlice("discover-domain", nil, "")
	cmd.Flags().String("targets-file", "", "")
	cmd.Flags().String("auth-token", "", "")
	cmd.Flags().Int("concurrency", 5, "")
	cmd.Flags().Duration("timeout", 0, "")
	cmd.Flags().Bool("insecure", false, "")
	cmd.Flags().String("output", "", "")
	cmd.Flags().String("fail-on", "", "")
	_ = cmd.Flags().Set("a2a", "true")
	err := runScan(cmd, nil)
	if err == nil {
		t.Fatal("expected error for A2A without target")
	}
	if !strings.Contains(err.Error(), "A2A requires") {
		t.Errorf("error = %q, want 'A2A requires'", err.Error())
	}
}

func TestRunScan_InvalidFailOn(t *testing.T) {
	cmd := &cobra.Command{RunE: runScan}
	cmd.Flags().Bool("config", false, "")
	cmd.Flags().Bool("mcp", false, "")
	cmd.Flags().Bool("a2a", false, "")
	cmd.Flags().String("path", "", "")
	cmd.Flags().StringSlice("paths", nil, "")
	cmd.Flags().String("project-dir", "", "")
	cmd.Flags().Bool("include-credential-values", false, "")
	cmd.Flags().String("url", "", "")
	cmd.Flags().String("target", "", "")
	cmd.Flags().StringSlice("targets", nil, "")
	cmd.Flags().StringSlice("discover-domain", nil, "")
	cmd.Flags().String("targets-file", "", "")
	cmd.Flags().String("auth-token", "", "")
	cmd.Flags().Int("concurrency", 5, "")
	cmd.Flags().Duration("timeout", 0, "")
	cmd.Flags().Bool("insecure", false, "")
	cmd.Flags().String("output", "", "")
	cmd.Flags().String("fail-on", "", "")
	_ = cmd.Flags().Set("fail-on", "banana")
	err := runScan(cmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid --fail-on")
	}
	if !strings.Contains(err.Error(), "invalid --fail-on") {
		t.Errorf("error = %q, want 'invalid --fail-on'", err.Error())
	}
}

func TestSeverityRank(t *testing.T) {
	if severityRank["critical"] != 0 {
		t.Errorf("critical rank = %d, want 0", severityRank["critical"])
	}
	if severityRank["high"] != 1 {
		t.Errorf("high rank = %d, want 1", severityRank["high"])
	}
	if severityRank["medium"] != 2 {
		t.Errorf("medium rank = %d, want 2", severityRank["medium"])
	}
	if severityRank["low"] != 3 {
		t.Errorf("low rank = %d, want 3", severityRank["low"])
	}
	if len(severityRank) != 4 {
		t.Errorf("severity rank has %d entries, want 4", len(severityRank))
	}
}
