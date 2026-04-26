package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestRule(t *testing.T, dir, filename, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestNewEngine_LoadsBuiltins(t *testing.T) {
	eng, err := NewEngine(LoadOptions{})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	if eng.RuleCount() == 0 {
		t.Error("expected builtin rules to be loaded")
	}
}

func TestNewEngine_CustomRules(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "test.yaml", `
id: custom-engine-test
name: Engine Test Rule
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: regex
  pattern: '<IMPORTANT>'
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{CustomDir: dir, EnableOnly: []string{"custom-engine-test"}})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	if eng.RuleCount() != 1 {
		t.Errorf("RuleCount = %d, want 1", eng.RuleCount())
	}
}

func TestEngine_Evaluate_ScopeFiltering(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "mcp-only.yaml", `
id: mcp-only-rule-test
name: MCP Only
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [dangerous]
  case_insensitive: true
emit:
  finding_type: test
`)
	writeTestRule(t, dir, "all-scope.yaml", `
id: all-scope-rule-test
name: All Scope
version: 1
severity: medium
scope:
  collector: all
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [warning]
  case_insensitive: true
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{
		CustomDir:  dir,
		EnableOnly: []string{"mcp-only-rule-test", "all-scope-rule-test"},
	})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	if eng.RuleCount() != 2 {
		t.Fatalf("RuleCount = %d, want 2", eng.RuleCount())
	}

	matches := eng.Evaluate("mcp", "tool.description", "this is dangerous and a warning")
	ruleIDs := make(map[string]bool)
	for _, m := range matches {
		ruleIDs[m.RuleID] = true
	}
	if !ruleIDs["mcp-only-rule-test"] {
		t.Error("mcp-only-rule-test should match for collector=mcp")
	}
	if !ruleIDs["all-scope-rule-test"] {
		t.Error("all-scope-rule-test should match for collector=mcp (scope=all)")
	}

	matches = eng.Evaluate("a2a", "tool.description", "this is dangerous and a warning")
	ruleIDs = make(map[string]bool)
	for _, m := range matches {
		ruleIDs[m.RuleID] = true
	}
	if ruleIDs["mcp-only-rule-test"] {
		t.Error("mcp-only-rule-test should NOT match for collector=a2a")
	}
	if !ruleIDs["all-scope-rule-test"] {
		t.Error("all-scope-rule-test should match for collector=a2a (scope=all)")
	}
}

func TestEngine_Evaluate_NoMatchOnWrongTarget(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "tool-only.yaml", `
id: tool-only-target-test
name: Tool Only
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [zzzuniquezzz]
  case_insensitive: true
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{CustomDir: dir, EnableOnly: []string{"tool-only-target-test"}})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}

	matches := eng.Evaluate("mcp", "tool.name", "zzzuniquezzz")
	if len(matches) != 0 {
		t.Errorf("expected 0 matches for wrong target, got %d", len(matches))
	}

	matches = eng.Evaluate("mcp", "tool.description", "zzzuniquezzz")
	if len(matches) != 1 {
		t.Errorf("expected 1 match for correct target, got %d", len(matches))
	}
}

func TestEngine_CustomOverridesBuiltin(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "override.yaml", `
id: injection-important-tag
name: Overridden Important Tag
version: 1
severity: info
scope:
  collector: all
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [override_keyword_unique]
  case_insensitive: true
emit:
  finding_type: overridden
`)

	eng, err := NewEngine(LoadOptions{CustomDir: dir, EnableOnly: []string{"injection-important-tag"}})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}

	rules := eng.Rules()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Emit.FindingType != "overridden" {
		t.Errorf("custom rule should override builtin, got finding_type=%q", rules[0].Emit.FindingType)
	}
	if rules[0].Severity != "info" {
		t.Errorf("custom rule severity should override, got %q", rules[0].Severity)
	}
}

func TestEngine_DisableRules(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "rule1.yaml", `
id: disable-test-one
name: Rule One
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [zzzdisabletestzzz]
  case_insensitive: true
emit:
  finding_type: test
`)
	writeTestRule(t, dir, "rule2.yaml", `
id: disable-test-two
name: Rule Two
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [zzzdisabletestzzz]
  case_insensitive: true
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{
		CustomDir:  dir,
		EnableOnly: []string{"disable-test-one", "disable-test-two"},
		DisableIDs: []string{"disable-test-one"},
	})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	if eng.RuleCount() != 1 {
		t.Errorf("RuleCount = %d, want 1 (one disabled)", eng.RuleCount())
	}

	matches := eng.Evaluate("mcp", "tool.description", "zzzdisabletestzzz")
	for _, m := range matches {
		if m.RuleID == "disable-test-one" {
			t.Error("disabled rule should not produce matches")
		}
	}
}

func TestEngine_EnableOnly(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "rule1.yaml", `
id: enable-only-one
name: Rule One
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [test]
  case_insensitive: true
emit:
  finding_type: test
`)
	writeTestRule(t, dir, "rule2.yaml", `
id: enable-only-two
name: Rule Two
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [test]
  case_insensitive: true
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{
		CustomDir:  dir,
		EnableOnly: []string{"enable-only-one"},
	})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	if eng.RuleCount() != 1 {
		t.Errorf("RuleCount = %d, want 1", eng.RuleCount())
	}
}

func TestEngine_EvaluateAll(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "rule.yaml", `
id: evaluate-all-test
name: Evaluate All
version: 1
severity: medium
scope:
  collector: mcp
  targets: [tool.description, tool.name]
matcher:
  type: keyword
  keywords: [zzzevalallzzz]
  case_insensitive: true
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{CustomDir: dir, EnableOnly: []string{"evaluate-all-test"}})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}

	matches := eng.EvaluateAll("mcp", map[string]string{
		"tool.description": "zzzevalallzzz command",
		"tool.name":        "zzzevalallzzz_tool",
	})
	if len(matches) != 2 {
		t.Errorf("expected 2 matches (one per field), got %d", len(matches))
	}
}

func TestEngine_InputTruncation(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "rule.yaml", `
id: truncation-test-rule
name: Truncation Test
version: 1
severity: low
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: regex
  pattern: 'ZZZUNIQUEMARKERZZZ'
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{CustomDir: dir, EnableOnly: []string{"truncation-test-rule"}})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}

	input := strings.Repeat("x", maxInputBytes+100) + "ZZZUNIQUEMARKERZZZ"
	matches := eng.Evaluate("mcp", "tool.description", input)
	if len(matches) != 0 {
		t.Error("marker beyond 1MiB should be truncated and not matched")
	}
}

func TestEngine_MatchTextTruncation(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "rule.yaml", `
id: match-text-trunc-test
name: Match Truncation
version: 1
severity: low
scope:
  collector: config
  targets: [credential.value]
matcher:
  type: regex
  pattern: 'Z+'
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{CustomDir: dir, EnableOnly: []string{"match-text-trunc-test"}})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}

	input := strings.Repeat("Z", 200)
	matches := eng.Evaluate("config", "credential.value", input)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if len(matches[0].Text) > 100 {
		t.Errorf("match text length = %d, should be truncated to 100", len(matches[0].Text))
	}
}

func TestEngine_InvalidRuleSkipped(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "invalid.yaml", `
id: x
name: ""
severity: bogus
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: regex
  pattern: test
emit:
  finding_type: test
`)
	writeTestRule(t, dir, "valid.yaml", `
id: valid-after-invalid
name: Valid Rule
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: regex
  pattern: 'test'
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{
		CustomDir:  dir,
		EnableOnly: []string{"x", "valid-after-invalid"},
	})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	if eng.RuleCount() != 1 {
		t.Errorf("RuleCount = %d, want 1 (invalid rule skipped)", eng.RuleCount())
	}
}

func TestEngine_DisabledRuleSkipped(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "disabled.yaml", `
id: explicitly-disabled
name: Disabled Rule
version: 1
enabled: false
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [test]
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{
		CustomDir:  dir,
		EnableOnly: []string{"explicitly-disabled"},
	})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	if eng.RuleCount() != 0 {
		t.Errorf("RuleCount = %d, want 0 (disabled rule)", eng.RuleCount())
	}
}

func TestEngine_Rules_ReturnsAllActive(t *testing.T) {
	dir := t.TempDir()
	writeTestRule(t, dir, "r1.yaml", `
id: rules-list-test-aaa
name: Rule A
version: 1
severity: high
scope:
  collector: mcp
  targets: [tool.description]
matcher:
  type: keyword
  keywords: [test]
emit:
  finding_type: test
`)
	writeTestRule(t, dir, "r2.yaml", `
id: rules-list-test-bbb
name: Rule B
version: 1
severity: low
scope:
  collector: a2a
  targets: [skill.description]
matcher:
  type: keyword
  keywords: [test]
emit:
  finding_type: test
`)

	eng, err := NewEngine(LoadOptions{
		CustomDir:  dir,
		EnableOnly: []string{"rules-list-test-aaa", "rules-list-test-bbb"},
	})
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	rules := eng.Rules()
	if len(rules) != 2 {
		t.Errorf("Rules() returned %d, want 2", len(rules))
	}
}
