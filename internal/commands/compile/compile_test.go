package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompile_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	result, err := Compile(CompileOptions{ProjectRoot: dir})
	// May return error if no .apm dir found, or succeed with zero sections
	_ = err
	_ = result
}

func TestCompile_WithApmDir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "instructions")
	os.MkdirAll(apmDir, 0755)
	os.WriteFile(filepath.Join(apmDir, "test.instructions.md"), []byte("# Test\nHello world."), 0644)

	result, err := Compile(CompileOptions{ProjectRoot: dir})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCompile_DryRun(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "instructions")
	os.MkdirAll(apmDir, 0755)
	os.WriteFile(filepath.Join(apmDir, "a.instructions.md"), []byte("# A\nContent."), 0644)

	result, err := Compile(CompileOptions{ProjectRoot: dir, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if !result.DryRun {
		t.Error("expected DryRun flag set")
	}
}

func TestCompile_OutputPath(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "instructions")
	os.MkdirAll(apmDir, 0755)
	os.WriteFile(filepath.Join(apmDir, "a.instructions.md"), []byte("# A\nContent."), 0644)

	outPath := filepath.Join(dir, "AGENTS.md")
	result, err := Compile(CompileOptions{ProjectRoot: dir, Output: outPath})
	if err != nil {
		t.Fatal(err)
	}
	if result.OutputPath != outPath {
		t.Errorf("expected output path %s, got %s", outPath, result.OutputPath)
	}
	// File should be written
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}

func TestCompile_ForceRewrite(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "instructions")
	os.MkdirAll(apmDir, 0755)
	os.WriteFile(filepath.Join(apmDir, "a.instructions.md"), []byte("# A\nContent."), 0644)

	outPath := filepath.Join(dir, "AGENTS.md")
	// Write once
	Compile(CompileOptions{ProjectRoot: dir, Output: outPath})
	// Write again with Force
	result, err := Compile(CompileOptions{ProjectRoot: dir, Output: outPath, Force: true})
	if err != nil {
		t.Fatal(err)
	}
	_ = result
}

func TestCompile_MultipleInstructions(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "instructions")
	os.MkdirAll(apmDir, 0755)
	for i, name := range []string{"a", "b", "c"} {
		_ = i
		os.WriteFile(filepath.Join(apmDir, name+".instructions.md"), []byte("# "+name+"\nContent."), 0644)
	}

	result, err := Compile(CompileOptions{ProjectRoot: dir})
	if err != nil {
		t.Fatal(err)
	}
	if result.Stats.Instructions < 3 {
		t.Errorf("expected >=3 instructions, got %d", result.Stats.Instructions)
	}
}

func TestCompileStats_Accumulate(t *testing.T) {
	s := CompileStats{}
	s.Instructions++
	s.Contexts++
	s.Primitives = s.Instructions + s.Contexts
	if s.Primitives != 2 {
		t.Errorf("expected 2, got %d", s.Primitives)
	}
}

func TestExtractTitle(t *testing.T) {
	cases := []struct {
		content  string
		filename string
		want     string
	}{
		{"# My Title\nContent", "test.md", "My Title"},
		{"No heading here", "myfile.instructions.md", "myfile"},
		{"## Second level\nContent", "file.md", "file"},
	}
	for _, c := range cases {
		got := extractTitle(c.content, c.filename)
		if !strings.Contains(got, c.want) {
			t.Errorf("extractTitle(%q, %q) = %q, want to contain %q", c.content, c.filename, got, c.want)
		}
	}
}

func TestComputeHash(t *testing.T) {
	h1 := computeHash("hello")
	h2 := computeHash("hello")
	h3 := computeHash("world")
	if h1 != h2 {
		t.Error("same input should produce same hash")
	}
	if h1 == h3 {
		t.Error("different inputs should produce different hashes")
	}
	if len(h1) < 8 {
		t.Errorf("expected non-trivial hash length, got %d chars", len(h1))
	}
}

func TestWriteAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.md")
	err := writeAtomic(path, []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "hello world" {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestFileMatchesContent(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "f.md")
	os.WriteFile(f, []byte("match me"), 0644)
	if !fileMatchesContent(f, "match me") {
		t.Error("expected match")
	}
	if fileMatchesContent(f, "different") {
		t.Error("expected no match")
	}
	if fileMatchesContent("/nonexistent", "x") {
		t.Error("nonexistent file should not match")
	}
}
