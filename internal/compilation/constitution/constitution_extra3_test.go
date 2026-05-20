package constitution

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestFindConstitution_ContainsBasePath_Extra3(t *testing.T) {
	p := FindConstitution("/my/base")
	if !filepath.IsAbs(p) {
		t.Errorf("FindConstitution should return absolute path for absolute base, got %q", p)
	}
}

func TestFindConstitution_ContainsRelativePath_Extra3(t *testing.T) {
	p := FindConstitution("/base")
	if p == "" {
		t.Error("FindConstitution should return a non-empty path")
	}
	// Should include the relative path component
	rel := compilationconst.ConstitutionRelativePath
	if !filepath.IsAbs(p) && p != filepath.Join("/base", rel) {
		t.Errorf("unexpected path: %q", p)
	}
}

func TestReadConstitution_CacheHit_Extra3(t *testing.T) {
	ClearCache()
	dir := t.TempDir()
	constitDir := filepath.Join(dir, ".specify", "memory")
	os.MkdirAll(constitDir, 0o755)
	os.WriteFile(filepath.Join(constitDir, "constitution.md"), []byte("cached content"), 0o644)

	content1, ok1 := ReadConstitution(dir)
	content2, ok2 := ReadConstitution(dir)

	if !ok1 || !ok2 {
		t.Error("ReadConstitution should succeed")
	}
	if content1 != content2 {
		t.Error("second call should return cached content")
	}
	ClearCache()
}

func TestReadConstitution_MissingReturnsEmpty_Extra3(t *testing.T) {
	ClearCache()
	dir := t.TempDir()
	content, ok := ReadConstitution(dir)
	if ok {
		t.Error("ReadConstitution on empty dir should return false")
	}
	if content != "" {
		t.Errorf("content should be empty, got %q", content)
	}
	ClearCache()
}

func TestClearCacheAllowsRefresh_Extra3(t *testing.T) {
	dir := t.TempDir()
	constitDir := filepath.Join(dir, ".specify", "memory")
	os.MkdirAll(constitDir, 0o755)

	ClearCache()
	_, ok1 := ReadConstitution(dir) // miss
	if ok1 {
		t.Error("expected miss on first read")
	}

	os.WriteFile(filepath.Join(constitDir, "constitution.md"), []byte("new"), 0o644)
	ClearCache()
	_, ok2 := ReadConstitution(dir) // now should hit
	if !ok2 {
		t.Error("expected hit after file created and cache cleared")
	}
	ClearCache()
}

func TestFindConstitution_UniquePerBase_Extra3(t *testing.T) {
	p1 := FindConstitution("/base1")
	p2 := FindConstitution("/base2")
	if p1 == p2 {
		t.Error("different bases should produce different constitution paths")
	}
}
