package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompileOptions_ZeroValue_Extra3(t *testing.T) {
	var opts CompileOptions
	if opts.DryRun || opts.Watch || opts.Force || opts.Strict || opts.Verbose {
		t.Error("zero CompileOptions should have all bool fields false")
	}
}

func TestCompileStats_ZeroValue_Extra3(t *testing.T) {
	var s CompileStats
	if s.Instructions != 0 || s.Contexts != 0 || s.Chatmodes != 0 || s.Primitives != 0 {
		t.Error("zero CompileStats should have all counts zero")
	}
	if len(s.Warnings) != 0 {
		t.Error("zero CompileStats should have no warnings")
	}
}

func TestCompileStats_AssignFields_Extra3(t *testing.T) {
	s := CompileStats{
		Instructions: 3,
		Contexts:     2,
		Chatmodes:    1,
		Primitives:   6,
		Warnings:     []string{"warn1"},
	}
	if s.Primitives != 6 {
		t.Errorf("expected Primitives=6, got %d", s.Primitives)
	}
	if len(s.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(s.Warnings))
	}
}

func TestCompileResult_ZeroValue_Extra3(t *testing.T) {
	var r CompileResult
	if r.OutputPath != "" || r.Status != "" || r.DryRun {
		t.Error("zero CompileResult should have empty/false fields")
	}
}

func TestComputeHash_NonEmpty_Extra3(t *testing.T) {
	h := computeHash("some content")
	if h == "" {
		t.Error("computeHash should return non-empty string")
	}
}

func TestComputeHash_DifferentInputs_DifferentHashes_Extra3(t *testing.T) {
	a := computeHash("abc")
	b := computeHash("xyz")
	if a == b {
		t.Error("different inputs should produce different hashes")
	}
}

func TestCompile_WithApmDir_NoError_Extra3(t *testing.T) {
	tmp := t.TempDir()
	apmDir := filepath.Join(tmp, ".apm", "instructions")
	if err := os.MkdirAll(apmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	inst := filepath.Join(apmDir, "test.instructions.md")
	if err := os.WriteFile(inst, []byte("# Test instruction\nDo something."), 0o644); err != nil {
		t.Fatal(err)
	}
	opts := CompileOptions{ProjectRoot: tmp, DryRun: true}
	result, err := Compile(opts)
	if err != nil {
		t.Fatalf("Compile with instructions returned error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCompile_DryRun_ResultFlagSet_Extra3(t *testing.T) {
	tmp := t.TempDir()
	// Create a minimal .apm directory so Compile does not error.
	apmDir := filepath.Join(tmp, ".apm")
	if err := os.MkdirAll(apmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	opts := CompileOptions{ProjectRoot: tmp, DryRun: true}
	result, err := Compile(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DryRun != true {
		t.Error("result.DryRun should be true when DryRun option is true")
	}
	_ = strings.Contains("unused import prevention", "x")
}
