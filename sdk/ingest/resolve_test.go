package ingest

import "testing"

func TestResolveEdgeEndpoints_KnownKind(t *testing.T) {
	src, tgt := ResolveEdgeEndpoints("PROVIDES_TOOL", "", "")
	if src != "MCPServer" {
		t.Errorf("source: got %q, want %q", src, "MCPServer")
	}
	if tgt != "MCPTool" {
		t.Errorf("target: got %q, want %q", tgt, "MCPTool")
	}
}

func TestResolveEdgeEndpoints_UnknownKind(t *testing.T) {
	src, tgt := ResolveEdgeEndpoints("MADE_UP_KIND", "", "")
	if src != "" {
		t.Errorf("source: got %q, want empty", src)
	}
	if tgt != "" {
		t.Errorf("target: got %q, want empty", tgt)
	}
}

func TestResolveEdgeEndpoints_ExplicitOverride(t *testing.T) {
	src, tgt := ResolveEdgeEndpoints("PROVIDES_TOOL", "CustomSource", "CustomTarget")
	if src != "CustomSource" {
		t.Errorf("source: got %q, want %q", src, "CustomSource")
	}
	if tgt != "CustomTarget" {
		t.Errorf("target: got %q, want %q", tgt, "CustomTarget")
	}
}

func TestResolveEdgeEndpoints_PartialOverrideSource(t *testing.T) {
	src, tgt := ResolveEdgeEndpoints("PROVIDES_TOOL", "CustomSource", "")
	if src != "CustomSource" {
		t.Errorf("source: got %q, want %q", src, "CustomSource")
	}
	if tgt != "MCPTool" {
		t.Errorf("target: got %q, want %q (registry fallback)", tgt, "MCPTool")
	}
}

func TestResolveEdgeEndpoints_PartialOverrideTarget(t *testing.T) {
	src, tgt := ResolveEdgeEndpoints("PROVIDES_TOOL", "", "CustomTarget")
	if src != "MCPServer" {
		t.Errorf("source: got %q, want %q (registry fallback)", src, "MCPServer")
	}
	if tgt != "CustomTarget" {
		t.Errorf("target: got %q, want %q", tgt, "CustomTarget")
	}
}

func TestResolveEdgeEndpoints_AllRegisteredKindsResolvable(t *testing.T) {
	for kind, ep := range EdgeKindEndpoints {
		src, tgt := ResolveEdgeEndpoints(kind, "", "")
		if src == "" {
			t.Errorf("%s: source resolved to empty (registry has %v)", kind, ep.SourceKinds)
		}
		if tgt == "" {
			t.Errorf("%s: target resolved to empty (registry has %v)", kind, ep.TargetKinds)
		}
	}
}
