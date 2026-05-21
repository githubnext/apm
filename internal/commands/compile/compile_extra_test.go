package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeHash_Deterministic(t *testing.T) {
	h1 := computeHash("hello world")
	h2 := computeHash("hello world")
	if h1 != h2 {
		t.Errorf("hash not deterministic: %q vs %q", h1, h2)
	}
}

func TestComputeHash_Different(t *testing.T) {
	h1 := computeHash("foo")
	h2 := computeHash("bar")
	if h1 == h2 {
		t.Error("different content should produce different hashes")
	}
}

func TestComputeHash_Length(t *testing.T) {
	h := computeHash("test content")
	// 8 bytes = 16 hex chars
	if len(h) != 16 {
		t.Errorf("expected 16 hex chars, got %d: %q", len(h), h)
	}
}

func TestExtractTitle_FromHeading(t *testing.T) {
	content := "# My Title\n\nSome content."
	title := extractTitle(content, "test.md")
	if title != "My Title" {
		t.Errorf("expected 'My Title', got %q", title)
	}
}

func TestExtractTitle_Fallback(t *testing.T) {
	content := "no heading here\njust text"
	title := extractTitle(content, "myfile.instructions.md")
	if title != "myfile.instructions" {
		t.Errorf("expected filename fallback 'myfile.instructions', got %q", title)
	}
}

func TestExtractTitle_EmptyContent(t *testing.T) {
	title := extractTitle("", "somefile.md")
	if title != "somefile" {
		t.Errorf("expected 'somefile' from empty content, got %q", title)
	}
}

func TestExtractTitle_SecondaryHeading(t *testing.T) {
	// ## heading should NOT be used as title
	content := "## secondary\n# primary"
	title := extractTitle(content, "f.md")
	if title != "primary" {
		t.Errorf("expected 'primary', got %q", title)
	}
}

func TestBuildConstitution_Empty(t *testing.T) {
	result := buildConstitution(nil)
	if !strings.Contains(result, "No primitives found") {
		t.Errorf("expected 'No primitives found' in empty constitution: %q", result)
	}
}

func TestBuildConstitution_WithSection(t *testing.T) {
	sections := []PrimitiveSection{
		{Title: "MyInst", Kind: "instruction", Content: "# MyInst\nDo this.", Path: "a.md"},
	}
	result := buildConstitution(sections)
	if !strings.Contains(result, "MyInst") {
		t.Errorf("expected title in constitution: %q", result)
	}
	if !strings.Contains(result, "Do this.") {
		t.Errorf("expected content in constitution: %q", result)
	}
}

func TestBuildConstitution_MultipleTypes(t *testing.T) {
	sections := []PrimitiveSection{
		{Title: "CtxA", Kind: "context", Content: "context content", Path: "ctx.md"},
		{Title: "InstB", Kind: "instruction", Content: "instruction content", Path: "inst.md"},
	}
	result := buildConstitution(sections)
	if !strings.Contains(result, "instruction content") || !strings.Contains(result, "context content") {
		t.Error("expected both sections in constitution")
	}
}

func TestFileMatchesContent_Match(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !fileMatchesContent(path, "hello") {
		t.Error("expected match for same content")
	}
}

func TestFileMatchesContent_Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if fileMatchesContent(path, "world") {
		t.Error("expected mismatch for different content")
	}
}

func TestFileMatchesContent_Missing(t *testing.T) {
	if fileMatchesContent("/nonexistent/path/abc.txt", "content") {
		t.Error("expected false for missing file")
	}
}

func TestWriteAtomic_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	if err := writeAtomic(path, []byte("written")); err != nil {
		t.Fatalf("writeAtomic failed: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "written" {
		t.Errorf("expected 'written', got %q", string(got))
	}
}

func TestWriteAtomic_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	os.WriteFile(path, []byte("old"), 0o644)
	if err := writeAtomic(path, []byte("new")); err != nil {
		t.Fatalf("writeAtomic overwrite failed: %v", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "new" {
		t.Errorf("expected 'new', got %q", string(got))
	}
}

func TestCompileOptions_ZeroValue(t *testing.T) {
	var opts CompileOptions
	if opts.DryRun || opts.Watch || opts.Force || opts.Strict || opts.Verbose {
		t.Error("zero value should have all bool options false")
	}
}

func TestCompileStats_ZeroValue(t *testing.T) {
	var stats CompileStats
	if stats.Instructions != 0 || stats.Contexts != 0 || stats.Chatmodes != 0 || stats.Primitives != 0 {
		t.Error("zero value should have all counts 0")
	}
	if len(stats.Warnings) != 0 {
		t.Error("zero value should have no warnings")
	}
}

func TestCompileResult_ZeroValue(t *testing.T) {
	var r CompileResult
	if r.DryRun {
		t.Error("DryRun should default false")
	}
}
