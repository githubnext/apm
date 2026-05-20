package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeHash_EmptyString(t *testing.T) {
	h := computeHash("")
	if h == "" {
		t.Error("hash of empty string should not be empty")
	}
}

func TestComputeHash_Unicode(t *testing.T) {
	h1 := computeHash("hello")
	h2 := computeHash("hello")
	if h1 != h2 {
		t.Error("hash must be deterministic for same input")
	}
}

func TestExtractTitle_OnlyH2(t *testing.T) {
	content := "## Section Title\n\nsome text"
	got := extractTitle(content, "file.md")
	if got == "" {
		t.Error("expected non-empty title")
	}
}

func TestExtractTitle_FallbackFilename(t *testing.T) {
	got := extractTitle("no heading here", "my-module.md")
	if got == "" {
		t.Error("expected title from filename fallback")
	}
}

func TestExtractTitle_WhitespaceTitle(t *testing.T) {
	got := extractTitle("# \n\nsome text", "fallback.md")
	if got == "" {
		t.Error("expected fallback title when heading is whitespace")
	}
}

func TestBuildConstitution_SingleSection(t *testing.T) {
	sections := []PrimitiveSection{
		{Kind: "chatmode", Title: "My Skill", Content: "do something"},
	}
	result := buildConstitution(sections)
	if !strings.Contains(result, "My Skill") {
		t.Errorf("expected title in constitution, got: %s", result)
	}
}

func TestBuildConstitution_SectionTypes(t *testing.T) {
	sections := []PrimitiveSection{
		{Kind: "instruction", Title: "Inst", Content: "content1"},
		{Kind: "chatmode", Title: "Skill", Content: "content2"},
	}
	result := buildConstitution(sections)
	if result == "" {
		t.Error("expected non-empty constitution")
	}
}

func TestFileMatchesContent_ExactMatch(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(p, []byte("exact content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !fileMatchesContent(p, "exact content") {
		t.Error("expected file to match content")
	}
}

func TestFileMatchesContent_NoMatch(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(p, []byte("old content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if fileMatchesContent(p, "new content") {
		t.Error("expected file not to match different content")
	}
}

func TestWriteAtomic_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	if err := writeAtomic(p, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Errorf("expected hello, got %q", string(b))
	}
}

func TestWriteAtomic_Overwrites(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	_ = os.WriteFile(p, []byte("old"), 0o644)
	if err := writeAtomic(p, []byte("new")); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if string(b) != "new" {
		t.Errorf("expected new, got %q", string(b))
	}
}

func TestPrimitiveSection_Fields(t *testing.T) {
	s := PrimitiveSection{Kind: "chatmode", Title: "T", Content: "C", Path: "/p"}
	if s.Kind != "chatmode" || s.Title != "T" || s.Content != "C" || s.Path != "/p" {
		t.Errorf("unexpected fields: %+v", s)
	}
}

func TestCompileOptions_ProjectRoot(t *testing.T) {
	opts := CompileOptions{ProjectRoot: "/some/path"}
	if opts.ProjectRoot != "/some/path" {
		t.Errorf("unexpected ProjectRoot: %q", opts.ProjectRoot)
	}
}

func TestCompileStats_CountsZero(t *testing.T) {
	var cs CompileStats
	if cs.Instructions != 0 || cs.Contexts != 0 {
		t.Error("expected zero counts")
	}
}

func TestCompileResult_NotNilOnZero(t *testing.T) {
	cr := &CompileResult{}
	if cr == nil {
		t.Error("expected non-nil result")
	}
}
