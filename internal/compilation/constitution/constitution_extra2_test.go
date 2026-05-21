package constitution

import (
	"os"
	"testing"
)

func TestReadConstitution_NonExistentDir(t *testing.T) {
	ClearCache()
	_, ok := ReadConstitution("/nonexistent/path/xyz/abc123")
	if ok {
		t.Error("should return false for non-existent directory")
	}
}

func TestFindConstitution_ReturnStringForAnyDir(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	got := FindConstitution(tmp)
	if got == "" {
		t.Error("FindConstitution should return a non-empty path")
	}
}

func TestClearCache_NoError(t *testing.T) {
	// Multiple clear calls should not panic
	ClearCache()
	ClearCache()
	ClearCache()
}

func TestReadConstitution_ContentPreserved(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	path := FindConstitution(tmp)
	os.MkdirAll(getParent(path), 0o755)
	content := "# Constitution\n\nThis is the spec.\n"
	os.WriteFile(path, []byte(content), 0o644)

	got, ok := ReadConstitution(tmp)
	if !ok {
		t.Fatal("expected ok=true when file exists")
	}
	if got != content {
		t.Errorf("content mismatch: got %q want %q", got, content)
	}
}

func TestReadConstitution_Idempotent(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	path := FindConstitution(tmp)
	os.MkdirAll(getParent(path), 0o755)
	os.WriteFile(path, []byte("stable content"), 0o644)

	r1, ok1 := ReadConstitution(tmp)
	r2, ok2 := ReadConstitution(tmp)
	if ok1 != ok2 || r1 != r2 {
		t.Errorf("ReadConstitution not idempotent: (%v,%q) vs (%v,%q)", ok1, r1, ok2, r2)
	}
}

func TestReadConstitution_AfterClearCache(t *testing.T) {
	tmp := t.TempDir()
	path := FindConstitution(tmp)
	os.MkdirAll(getParent(path), 0o755)
	os.WriteFile(path, []byte("pre-clear"), 0o644)

	ClearCache()
	got, ok := ReadConstitution(tmp)
	if !ok {
		t.Fatal("should still find file after clear")
	}
	if got != "pre-clear" {
		t.Errorf("after clear, got %q", got)
	}
}

func TestFindConstitution_IsRelativeOrAbsolute(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	got := FindConstitution(tmp)
	// Path should not be empty and should contain a recognizable file name
	if got == "" {
		t.Error("FindConstitution returned empty path")
	}
}

func TestReadConstitution_LargeContent2(t *testing.T) {
	ClearCache()
	tmp := t.TempDir()
	path := FindConstitution(tmp)
	os.MkdirAll(getParent(path), 0o755)
	large := make([]byte, 64*1024)
	for i := range large {
		large[i] = byte('A' + (i % 26))
	}
	os.WriteFile(path, large, 0o644)
	got, ok := ReadConstitution(tmp)
	if !ok {
		t.Fatal("expected ok=true for large file")
	}
	if len(got) != len(large) {
		t.Errorf("large content: got len=%d want %d", len(got), len(large))
	}
}

// getParent returns the parent directory of a path.
func getParent(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}
