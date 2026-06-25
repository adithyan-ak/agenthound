package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestIsTerminal_NonFileWriterIsFalse guards the core safety property: a
// non-*os.File writer (bytes.Buffer here, a pipe in real piped output) is
// never treated as a TTY, so progress carriage returns can't pollute it.
func TestIsTerminal_NonFileWriterIsFalse(t *testing.T) {
	if isTerminal(&bytes.Buffer{}) {
		t.Error("bytes.Buffer must not be reported as a terminal")
	}
}

// TestProgressReporter_DisabledIsNoop verifies that a reporter built over a
// non-TTY writer writes nothing at all — the piped/CI path stays clean.
func TestProgressReporter_DisabledIsNoop(t *testing.T) {
	var buf bytes.Buffer
	r := newProgressReporter(&buf, "[scan]", false)
	r.update(1, 10)
	r.clear()
	if buf.Len() != 0 {
		t.Errorf("disabled reporter wrote %q, want nothing", buf.String())
	}
}

// TestProgressReporter_QuietDisables verifies --quiet disables rendering even
// if the destination were a terminal.
func TestProgressReporter_QuietDisables(t *testing.T) {
	if newProgressReporter(&bytes.Buffer{}, "[scan]", true).enabled {
		t.Error("reporter must be disabled when quiet is set")
	}
}

// TestProgressReporter_RendersWhenEnabled exercises the rendering path with a
// force-enabled reporter (no real TTY needed) and locks the line format.
func TestProgressReporter_RendersWhenEnabled(t *testing.T) {
	var buf bytes.Buffer
	r := &progressReporter{w: &buf, label: "[scan] probing x", enabled: true}

	r.update(3, 10)
	out := buf.String()
	if !strings.HasPrefix(out, "\r") {
		t.Errorf("update output must start with a carriage return, got %q", out)
	}
	if !strings.Contains(out, "[scan] probing x 3/10 (30%)") {
		t.Errorf("update output = %q, want it to contain '3/10 (30%%)'", out)
	}

	buf.Reset()
	r.clear()
	cleared := buf.String()
	if !strings.HasPrefix(cleared, "\r") || !strings.HasSuffix(cleared, "\r") {
		t.Errorf("clear output = %q, want it wrapped in carriage returns", cleared)
	}
	if r.lastLen != 0 {
		t.Errorf("lastLen after clear = %d, want 0", r.lastLen)
	}
}

// TestProgressReporter_PadsShrinkingLine ensures a shorter render overwrites a
// previous longer one so no stale glyphs remain on the line.
func TestProgressReporter_PadsShrinkingLine(t *testing.T) {
	var buf bytes.Buffer
	r := &progressReporter{w: &buf, label: "L", enabled: true}

	r.update(100, 100) // "L 100/100 (100%)" — the longer render
	prevLen := len("L 100/100 (100%)")

	buf.Reset()
	r.update(1, 100) // "L 1/100 (1%)" — shorter, must be padded
	rendered := strings.TrimPrefix(buf.String(), "\r")
	if len(rendered) < prevLen {
		t.Errorf("shrinking update rendered %d cols (%q), want >= %d to overwrite stale chars",
			len(rendered), rendered, prevLen)
	}
}
