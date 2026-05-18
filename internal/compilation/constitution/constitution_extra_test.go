package constitution

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConstitution_MultipleRoots(t *testing.T) {
	ClearCache()
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	path1 := FindConstitution(tmp1)
	if err := os.MkdirAll(filepath.Dir(path1), 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(path1, []byte("root1 content"), 0o644)

	// Only tmp1 has constitution; tmp2 does not
	got, ok := ReadConstitution(tmp1)
	if !ok || got != "root1 content" {
		t.Errorf("root1: ok=%v got=%q", ok, got)
	}

	_, ok2 := ReadConstitution(tmp2)
	if ok2 {
		t.Error("root2 should not find constitution")
	}
}

func TestReadConstitution_EmptyContent(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	path := FindConstitution(tmp)
	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, []byte(""), 0o644)

	got, ok := ReadConstitution(tmp)
	if !ok {
		t.Error("expected ok=true for empty file")
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestReadConstitution_LargeContent(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	path := FindConstitution(tmp)
	os.MkdirAll(filepath.Dir(path), 0o755)
	content := string(make([]byte, 10000))
	os.WriteFile(path, []byte(content), 0o644)

	got, ok := ReadConstitution(tmp)
	if !ok {
		t.Error("expected ok=true for large file")
	}
	if len(got) != len(content) {
		t.Errorf("content length mismatch: got %d, want %d", len(got), len(content))
	}
}

func TestClearCacheMultipleTimes(t *testing.T) {
	// Multiple ClearCache calls should not panic
	ClearCache()
	ClearCache()
	ClearCache()
}

func TestFindConstitutionRelative(t *testing.T) {
	// FindConstitution should return a path that ends with the expected relative path component
	got := FindConstitution("/some/repo")
	if got == "/some/repo" {
		t.Error("FindConstitution should return a path under baseDir")
	}
	if len(got) <= len("/some/repo") {
		t.Errorf("FindConstitution returned too short path: %q", got)
	}
}

func TestReadConstitutionCacheIsolation(t *testing.T) {
	ClearCache()
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	p1 := FindConstitution(tmp1)
	os.MkdirAll(filepath.Dir(p1), 0o755)
	os.WriteFile(p1, []byte("content-A"), 0o644)

	// Read tmp1 into cache
	got1, ok1 := ReadConstitution(tmp1)
	if !ok1 || got1 != "content-A" {
		t.Fatalf("tmp1 read failed: ok=%v got=%q", ok1, got1)
	}

	// tmp2 still missing -- cache should not confuse the two
	_, ok2 := ReadConstitution(tmp2)
	if ok2 {
		t.Error("tmp2 should not find constitution (cache isolation)")
	}
}
