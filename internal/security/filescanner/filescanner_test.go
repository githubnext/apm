package filescanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSafeLockfilePath(t *testing.T) {
	root := t.TempDir()

	tests := []struct {
		name     string
		relPath  string
		wantSafe bool
	}{
		{"simple file", "subdir/file.txt", true},
		{"traversal", "../outside.txt", false},
		{"double traversal", "a/../../outside.txt", false},
		{"root file", "file.txt", true},
		{"nested", "a/b/c/file.txt", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isSafeLockfilePath(tc.relPath, root)
			if got != tc.wantSafe {
				t.Errorf("isSafeLockfilePath(%q, root) = %v, want %v", tc.relPath, got, tc.wantSafe)
			}
		})
	}
}

func TestDetectSuspiciousBytes_Clean(t *testing.T) {
	findings := detectSuspiciousBytes([]byte("normal ASCII content\nno hidden chars"))
	if len(findings) != 0 {
		t.Errorf("expected no findings for clean content, got %v", findings)
	}
}

func TestDetectSuspiciousBytes_ZeroWidthSpace(t *testing.T) {
	content := "normal\u200Bcontent"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for zero-width space")
	}
	if findings[0].Type != "hidden-character" {
		t.Errorf("expected Type=hidden-character, got %q", findings[0].Type)
	}
}

func TestDetectSuspiciousBytes_RLOverride(t *testing.T) {
	content := "evil\u202Econtent"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for right-to-left override")
	}
}

func TestDetectSuspiciousBytes_BOM(t *testing.T) {
	content := "\uFEFFcontent"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for BOM")
	}
}

func TestScanDeployedFiles_Empty(t *testing.T) {
	lockData := LockFileData{Dependencies: map[string]LockedDependency{}}
	findings, count := ScanDeployedFiles(lockData, t.TempDir(), "")
	if len(findings) != 0 || count != 0 {
		t.Errorf("expected empty results for empty lock data, got findings=%v count=%d", findings, count)
	}
}

func TestScanDeployedFiles_CleanFile(t *testing.T) {
	root := t.TempDir()
	f := filepath.Join(root, "clean.txt")
	if err := os.WriteFile(f, []byte("hello world"), 0600); err != nil {
		t.Fatal(err)
	}

	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"pkg/a": {DeployedFiles: []string{"clean.txt"}},
		},
	}
	findings, count := ScanDeployedFiles(lockData, root, "")
	if count != 1 {
		t.Errorf("expected 1 file scanned, got %d", count)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings for clean file, got %v", findings)
	}
}

func TestScanDeployedFiles_SuspiciousFile(t *testing.T) {
	root := t.TempDir()
	f := filepath.Join(root, "evil.txt")
	if err := os.WriteFile(f, []byte("evil\u202Econtent"), 0600); err != nil {
		t.Fatal(err)
	}

	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"pkg/a": {DeployedFiles: []string{"evil.txt"}},
		},
	}
	findings, count := ScanDeployedFiles(lockData, root, "")
	if count != 1 {
		t.Errorf("expected 1 file scanned, got %d", count)
	}
	if len(findings) == 0 {
		t.Error("expected findings for suspicious file")
	}
}

func TestScanDeployedFiles_PackageFilter(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("ok"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"pkg/a": {DeployedFiles: []string{"a.txt"}},
			"pkg/b": {DeployedFiles: []string{"b.txt"}},
		},
	}
	_, count := ScanDeployedFiles(lockData, root, "pkg/a")
	if count != 1 {
		t.Errorf("expected 1 file when filtering to pkg/a, got %d", count)
	}
}

func TestScanDeployedFiles_UnsafePath(t *testing.T) {
	root := t.TempDir()
	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"pkg/x": {DeployedFiles: []string{"../outside.txt"}},
		},
	}
	_, count := ScanDeployedFiles(lockData, root, "")
	if count != 0 {
		t.Errorf("expected 0 scanned for unsafe path, got %d", count)
	}
}

func TestScanDeployedFiles_Directory(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "subdir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "clean.md"), []byte("# doc"), 0600); err != nil {
		t.Fatal(err)
	}

	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"pkg/a": {DeployedFiles: []string{"subdir/"}},
		},
	}
	_, count := ScanDeployedFiles(lockData, root, "")
	if count != 1 {
		t.Errorf("expected 1 file from dir scan, got %d", count)
	}
}
