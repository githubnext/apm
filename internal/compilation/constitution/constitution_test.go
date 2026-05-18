package constitution

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestFindConstitution(t *testing.T) {
	path := FindConstitution("/repo")
	want := filepath.Join("/repo", compilationconst.ConstitutionRelativePath)
	if path != want {
		t.Errorf("got %q, want %q", path, want)
	}
}

func TestReadConstitutionMissing(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	_, ok := ReadConstitution(tmp)
	if ok {
		t.Error("expected false for missing constitution")
	}
}

func TestReadConstitutionPresent(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	constitutionPath := FindConstitution(tmp)
	if err := os.MkdirAll(filepath.Dir(constitutionPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Constitution\n\nProject rules here."
	if err := os.WriteFile(constitutionPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, ok := ReadConstitution(tmp)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got != content {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestReadConstitutionCached(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	constitutionPath := FindConstitution(tmp)
	if err := os.MkdirAll(filepath.Dir(constitutionPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "cached content"
	if err := os.WriteFile(constitutionPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// First read
	got1, ok1 := ReadConstitution(tmp)
	if !ok1 || got1 != content {
		t.Fatalf("first read failed: ok=%v got=%q", ok1, got1)
	}

	// Modify file on disk -- cache should still return old value
	if err := os.WriteFile(constitutionPath, []byte("modified"), 0o644); err != nil {
		t.Fatal(err)
	}
	got2, ok2 := ReadConstitution(tmp)
	if !ok2 || got2 != content {
		t.Errorf("cache miss: ok=%v got=%q, want %q", ok2, got2, content)
	}
}

func TestClearCache(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	constitutionPath := FindConstitution(tmp)
	if err := os.MkdirAll(filepath.Dir(constitutionPath), 0o755); err != nil {
		t.Fatal(err)
	}

	// First read: missing
	_, ok := ReadConstitution(tmp)
	if ok {
		t.Error("should be missing")
	}

	// Write file and clear cache
	if err := os.WriteFile(constitutionPath, []byte("new content"), 0o644); err != nil {
		t.Fatal(err)
	}
	ClearCache()

	// Second read: should pick up the new file
	got, ok := ReadConstitution(tmp)
	if !ok || got != "new content" {
		t.Errorf("after ClearCache: ok=%v got=%q", ok, got)
	}
}
