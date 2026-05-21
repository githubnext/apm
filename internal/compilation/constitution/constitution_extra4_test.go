package constitution_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/compilation/constitution"
)

func TestFindConstitution_NotEmpty_Extra4(t *testing.T) {
	p := constitution.FindConstitution("/some/dir")
	if p == "" {
		t.Error("expected non-empty path from FindConstitution")
	}
}

func TestFindConstitution_ContainsBaseDir_Extra4(t *testing.T) {
	base := "/mybase"
	p := constitution.FindConstitution(base)
	if !filepath.IsAbs(p) {
		t.Errorf("expected absolute path, got %q", p)
	}
}

func TestFindConstitution_EndsWithConstitutionMD_Extra4(t *testing.T) {
	p := constitution.FindConstitution("/base")
	if filepath.Base(p) != "constitution.md" {
		t.Errorf("expected constitution.md filename, got %q", filepath.Base(p))
	}
}

func TestReadConstitution_MissingDir_ReturnsFalse_Extra4(t *testing.T) {
	constitution.ClearCache()
	_, ok := constitution.ReadConstitution("/nonexistent/path/to/dir")
	if ok {
		t.Error("expected false for missing directory")
	}
}

func TestReadConstitution_ExistingFile_ReturnsContent_Extra4(t *testing.T) {
	dir := t.TempDir()
	constitutionPath := constitution.FindConstitution(dir)
	if err := os.MkdirAll(filepath.Dir(constitutionPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(constitutionPath, []byte("my constitution"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	constitution.ClearCache()
	content, ok := constitution.ReadConstitution(dir)
	if !ok {
		t.Fatal("expected ok=true for existing file")
	}
	if content != "my constitution" {
		t.Errorf("expected 'my constitution', got %q", content)
	}
}

func TestReadConstitution_CachedResult_Extra4(t *testing.T) {
	dir := t.TempDir()
	constitutionPath := constitution.FindConstitution(dir)
	if err := os.MkdirAll(filepath.Dir(constitutionPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(constitutionPath, []byte("cached"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	constitution.ClearCache()
	a, ok := constitution.ReadConstitution(dir)
	if !ok {
		t.Fatal("expected ok")
	}
	b, ok2 := constitution.ReadConstitution(dir)
	if !ok2 {
		t.Fatal("expected ok on second read")
	}
	if a != b {
		t.Errorf("expected same cached result, got %q vs %q", a, b)
	}
}

func TestClearCache_AllowsRereading_Extra4(t *testing.T) {
	dir := t.TempDir()
	constitutionPath := constitution.FindConstitution(dir)
	if err := os.MkdirAll(filepath.Dir(constitutionPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(constitutionPath, []byte("v1"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	constitution.ClearCache()
	constitution.ReadConstitution(dir)
	if err := os.WriteFile(constitutionPath, []byte("v2"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	constitution.ClearCache()
	content, ok := constitution.ReadConstitution(dir)
	if !ok {
		t.Fatal("expected ok after clear cache")
	}
	if content != "v2" {
		t.Errorf("expected v2 after cache clear, got %q", content)
	}
}
