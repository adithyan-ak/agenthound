package rules

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRuleYAMLRoundtrip(t *testing.T) {
	input := `
id: test-rule
name: Test Rule
description: A test rule
version: 1
enabled: true
scope:
  collector: mcp
  targets:
    - tool.description
severity: high
owasp:
  - MCP04
tags:
  - injection
matcher:
  type: regex
  pattern: '<IMPORTANT>'
  case_insensitive: true
emit:
  finding_type: has_injection_patterns
  property_key: flagged
  property_value: true
  labels:
    - important_tag
tests:
  - input: "<IMPORTANT>override</IMPORTANT>"
    should_match: true
    description: detects important tag
  - input: clean description
    should_match: false
    description: no match on clean text
`
	var r Rule
	if err := yaml.Unmarshal([]byte(input), &r); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if r.ID != "test-rule" {
		t.Errorf("ID = %q, want %q", r.ID, "test-rule")
	}
	if r.Version != 1 {
		t.Errorf("Version = %d, want 1", r.Version)
	}
	if r.Scope.Collector != "mcp" {
		t.Errorf("Scope.Collector = %q, want %q", r.Scope.Collector, "mcp")
	}
	if len(r.Scope.Targets) != 1 || r.Scope.Targets[0] != "tool.description" {
		t.Errorf("Scope.Targets = %v, want [tool.description]", r.Scope.Targets)
	}
	if r.Severity != "high" {
		t.Errorf("Severity = %q, want %q", r.Severity, "high")
	}
	if r.Matcher.Type != "regex" {
		t.Errorf("Matcher.Type = %q, want %q", r.Matcher.Type, "regex")
	}
	if !r.Matcher.CaseInsensitive {
		t.Error("Matcher.CaseInsensitive = false, want true")
	}
	if r.Emit.FindingType != "has_injection_patterns" {
		t.Errorf("Emit.FindingType = %q, want %q", r.Emit.FindingType, "has_injection_patterns")
	}
	if len(r.Tests) != 2 {
		t.Errorf("Tests count = %d, want 2", len(r.Tests))
	}

	out, err := yaml.Marshal(&r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var r2 Rule
	if err := yaml.Unmarshal(out, &r2); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if r2.ID != r.ID {
		t.Errorf("roundtrip ID = %q, want %q", r2.ID, r.ID)
	}
	if r2.Matcher.Pattern != r.Matcher.Pattern {
		t.Errorf("roundtrip Pattern = %q, want %q", r2.Matcher.Pattern, r.Matcher.Pattern)
	}
}

func TestRuleSourceOmittedFromYAML(t *testing.T) {
	r := Rule{Source: "builtin"}
	out, err := yaml.Marshal(&r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var r2 Rule
	if err := yaml.Unmarshal(out, &r2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r2.Source != "" {
		t.Errorf("Source should be omitted from YAML, got %q", r2.Source)
	}
}

func TestScopeYAMLParsing(t *testing.T) {
	input := `
collector: all
targets:
  - tool.description
  - skill.description
  - server.instructions
`
	var s Scope
	if err := yaml.Unmarshal([]byte(input), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Collector != "all" {
		t.Errorf("Collector = %q, want %q", s.Collector, "all")
	}
	if len(s.Targets) != 3 {
		t.Errorf("Targets count = %d, want 3", len(s.Targets))
	}
}

func TestEmitConfigPropertyValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, ec EmitConfig)
	}{
		{
			name:  "bool value",
			input: "finding_type: test\nproperty_value: true",
			check: func(t *testing.T, ec EmitConfig) {
				if ec.PropertyValue != true {
					t.Errorf("PropertyValue = %v (%T), want true", ec.PropertyValue, ec.PropertyValue)
				}
			},
		},
		{
			name:  "string value",
			input: "finding_type: test\nproperty_value: shell_access",
			check: func(t *testing.T, ec EmitConfig) {
				if ec.PropertyValue != "shell_access" {
					t.Errorf("PropertyValue = %v, want shell_access", ec.PropertyValue)
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ec EmitConfig
			if err := yaml.Unmarshal([]byte(tc.input), &ec); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			tc.check(t, ec)
		})
	}
}
