package mcp

import (
	"encoding/json"
	"net/url"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/adithyan-ak/agenthound/internal/collector/common"
)

type ToolSignals struct {
	DescriptionHash    string
	CapabilitySurface  []string
	HasInjection       bool
	HasCrossReferences bool
	Annotations        map[string]any
}

type ResourceSignals struct {
	URIScheme   string
	Sensitivity string
}

func computeToolSignals(tool *mcpsdk.Tool, allToolNames map[string]bool) ToolSignals {
	schemaMap := inputSchemaAsMap(tool.InputSchema)

	sig := ToolSignals{
		DescriptionHash:   common.DescriptionHash(tool.Name, tool.Description, schemaMap),
		CapabilitySurface: common.ClassifyCapabilities(tool.Name, tool.Description, schemaMap),
		HasInjection:      common.HasInjectionPatterns(tool.Description),
	}

	if tool.Description != "" {
		descLower := strings.ToLower(tool.Description)
		for name := range allToolNames {
			if name != tool.Name && strings.Contains(descLower, strings.ToLower(name)) {
				sig.HasCrossReferences = true
				break
			}
		}
	}

	sig.Annotations = flattenAnnotations(tool.Annotations)

	return sig
}

func flattenAnnotations(ann *mcpsdk.ToolAnnotations) map[string]any {
	if ann == nil {
		return nil
	}
	m := make(map[string]any)
	m["read_only_hint"] = ann.ReadOnlyHint
	m["idempotent_hint"] = ann.IdempotentHint
	if ann.DestructiveHint != nil {
		m["destructive_hint"] = *ann.DestructiveHint
	}
	if ann.OpenWorldHint != nil {
		m["open_world_hint"] = *ann.OpenWorldHint
	}
	if ann.Title != "" {
		m["title"] = ann.Title
	}
	return m
}

func computeResourceSignals(uri string) ResourceSignals {
	sig := ResourceSignals{
		Sensitivity: string(common.ClassifyResourceSensitivity(uri)),
	}

	if u, err := url.Parse(uri); err == nil && u.Scheme != "" {
		sig.URIScheme = u.Scheme
	}

	return sig
}

func inputSchemaAsMap(schema any) map[string]any {
	if schema == nil {
		return nil
	}

	if m, ok := schema.(map[string]any); ok {
		return m
	}

	if s, ok := schema.(string); ok {
		var m map[string]any
		if err := json.Unmarshal([]byte(s), &m); err == nil {
			return m
		}
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	return m
}
