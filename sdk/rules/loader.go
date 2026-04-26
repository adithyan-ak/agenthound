package rules

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed builtin
var builtinFS embed.FS

func loadBuiltinRules() ([]Rule, error) {
	return loadRulesFromFS(builtinFS, "builtin", "builtin")
}

func loadRulesFromFS(fsys fs.FS, root string, source string) ([]Rule, error) {
	var rules []Rule
	entries, err := fs.ReadDir(fsys, root)
	if err != nil {
		return nil, fmt.Errorf("reading %s directory: %w", root, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(fsys, filepath.Join(root, entry.Name()))
		if err != nil {
			slog.Warn("failed to read rule file", "file", entry.Name(), "error", err)
			continue
		}
		r, err := parseRuleFile(data, source)
		if err != nil {
			slog.Warn("failed to parse rule file", "file", entry.Name(), "error", err)
			continue
		}
		rules = append(rules, *r)
	}
	return rules, nil
}

func loadCustomRules(dir string) ([]Rule, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat custom rules dir: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("custom rules path is not a directory: %s", dir)
	}

	var rules []Rule
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading custom rules dir: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			slog.Warn("failed to read custom rule", "path", path, "error", err)
			continue
		}
		r, err := parseRuleFile(data, path)
		if err != nil {
			slog.Warn("failed to parse custom rule", "path", path, "error", err)
			continue
		}
		rules = append(rules, *r)
	}
	return rules, nil
}

func parseRuleFile(data []byte, source string) (*Rule, error) {
	var r Rule
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("yaml unmarshal: %w", err)
	}
	r.Source = source
	if r.Version == 0 {
		r.Version = 1
	}
	if !r.yamlHasEnabled(data) {
		r.Enabled = true
	}
	return &r, nil
}

func (r *Rule) yamlHasEnabled(data []byte) bool {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return false
	}
	_, ok := raw["enabled"]
	return ok
}
