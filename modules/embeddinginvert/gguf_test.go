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

func q8FixturePath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "extract", "synthetic_q8.gguf")
}

func TestParseGGUF_Q8_0Fixture(t *testing.T) {
	gguf, err := ParseGGUF(q8FixturePath())
	if err != nil {
		t.Fatalf("ParseGGUF Q8_0: %v", err)
	}
	if gguf.Version != 3 {
		t.Errorf("Version = %d, want 3", gguf.Version)
	}
	if gguf.TensorType != ggmlTypeQ8_0 {
		t.Errorf("TensorType = %d, want Q8_0 (%d)", gguf.TensorType, ggmlTypeQ8_0)
	}
	if gguf.VocabSize != 10 {
		t.Errorf("VocabSize = %d, want 10", gguf.VocabSize)
	}
	if gguf.EmbedDim != 32 {
		t.Errorf("EmbedDim = %d, want 32", gguf.EmbedDim)
	}
	if len(gguf.Tokens) != 10 {
		t.Fatalf("Tokens len = %d, want 10", len(gguf.Tokens))
	}
	if gguf.Tokens[8] != "[secret_finetune_token]" {
		t.Errorf("Token[8] = %q, want [secret_finetune_token]", gguf.Tokens[8])
	}
	if len(gguf.Embeddings) != 10 {
		t.Fatalf("Embeddings rows = %d, want 10", len(gguf.Embeddings))
	}
	if len(gguf.Embeddings[0]) != 32 {
		t.Errorf("Embeddings cols = %d, want 32", len(gguf.Embeddings[0]))
	}
	// Verify dequantization produced non-zero values in expected range.
	// Normal rows: scale=0.01, values 5-14 → dequant magnitudes ~0.05-0.14
	// Outlier rows: scale=0.1, values 100-127 → dequant magnitudes ~10.0-12.7
	normalMag := l2Norm(gguf.Embeddings[0])
	outlierMag := l2Norm(gguf.Embeddings[8])
	if normalMag == 0 {
		t.Error("normal row has zero magnitude after dequant")
	}
	if outlierMag <= normalMag*2 {
		t.Errorf("outlier magnitude (%.3f) should be much larger than normal (%.3f)", outlierMag, normalMag)
	}
}

func TestParseGGUF_MultipleMetadataKeys(t *testing.T) {
	// The Q8_0 fixture has 3 metadata KVs: general.architecture (string),
	// general.context_length (uint32), tokenizer.ggml.tokens (array).
	// This exercises skipGGUFValue for string + uint32 types before
	// hitting the tokenizer array. If skipGGUFValue is broken, parsing
	// would either panic or produce wrong token data.
	gguf, err := ParseGGUF(q8FixturePath())
	if err != nil {
		t.Fatalf("ParseGGUF with multi-KV: %v", err)
	}
	if len(gguf.Tokens) != 10 {
		t.Errorf("after skipping 2 non-tokenizer KVs, Tokens len = %d, want 10", len(gguf.Tokens))
	}
	if gguf.Tokens[0] != "<pad>" {
		t.Errorf("first token after skip = %q, want <pad>", gguf.Tokens[0])
	}
}

func l2Norm(row []float32) float64 {
	var sum float64
	for _, v := range row {
		sum += float64(v) * float64(v)
	}
	return sum
}
