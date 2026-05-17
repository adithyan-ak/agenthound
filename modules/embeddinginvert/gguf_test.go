package embeddinginvert

import (
	"path/filepath"
	"runtime"
	"testing"
)

func fixturePath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "extract", "synthetic.gguf")
}

func TestParseGGUF_SyntheticFixture(t *testing.T) {
	gguf, err := ParseGGUF(fixturePath())
	if err != nil {
		t.Fatalf("ParseGGUF: %v", err)
	}
	if gguf.Version != 3 {
		t.Errorf("Version = %d, want 3", gguf.Version)
	}
	if gguf.VocabSize != 10 {
		t.Errorf("VocabSize = %d, want 10", gguf.VocabSize)
	}
	if gguf.EmbedDim != 8 {
		t.Errorf("EmbedDim = %d, want 8", gguf.EmbedDim)
	}
	if len(gguf.Tokens) != 10 {
		t.Errorf("Tokens len = %d, want 10", len(gguf.Tokens))
	}
	if gguf.Tokens[8] != "[fine_tune_secret]" {
		t.Errorf("Token[8] = %q, want [fine_tune_secret]", gguf.Tokens[8])
	}
	if gguf.Tokens[9] != "[internal_tool_xyz]" {
		t.Errorf("Token[9] = %q, want [internal_tool_xyz]", gguf.Tokens[9])
	}
	if len(gguf.Embeddings) != 10 {
		t.Fatalf("Embeddings rows = %d, want 10", len(gguf.Embeddings))
	}
	if len(gguf.Embeddings[0]) != 8 {
		t.Errorf("Embeddings cols = %d, want 8", len(gguf.Embeddings[0]))
	}
}

func TestParseGGUF_InvalidMagic(t *testing.T) {
	_, err := ParseGGUF("/dev/null")
	if err == nil {
		t.Error("expected error on /dev/null")
	}
}

func TestParseGGUF_NotFound(t *testing.T) {
	_, err := ParseGGUF("/nonexistent/path.gguf")
	if err == nil {
		t.Error("expected error on missing file")
	}
}
