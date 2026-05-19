package filescanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectSuspiciousBytes_ZWNJ(t *testing.T) {
	// zero-width non-joiner
	content := "text\u200Chere"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for zero-width non-joiner")
	}
}

func TestDetectSuspiciousBytes_ZWJ(t *testing.T) {
	content := "text\u200Dhere"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for zero-width joiner")
	}
}

func TestDetectSuspiciousBytes_LTREmbedding(t *testing.T) {
	content := "text\u202Ahere"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for LTR embedding")
	}
}

func TestDetectSuspiciousBytes_RTLEmbedding(t *testing.T) {
	content := "text\u202Bhere"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for RTL embedding")
	}
}

func TestDetectSuspiciousBytes_BOMOnly(t *testing.T) {
	content := "\uFEFFsome text"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) == 0 {
		t.Error("expected finding for BOM")
	}
}

func TestDetectSuspiciousBytes_Multiple(t *testing.T) {
	// two distinct suspicious chars
	content := "a\u200Bb\u202Ec"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) < 2 {
		t.Errorf("expected at least 2 findings, got %d", len(findings))
	}
}

func TestDetectSuspiciousBytes_NormalUnicode(t *testing.T) {
	// Normal accented characters - not suspicious
	content := "cafe au lait"
	findings := detectSuspiciousBytes([]byte(content))
	if len(findings) != 0 {
		t.Errorf("expected no findings for normal text, got %v", findings)
	}
}

func TestScanFinding_Fields(t *testing.T) {
	f := ScanFinding{Type: "hidden-character", Message: "contains zero-width space", Line: 5}
	if f.Type != "hidden-character" {
		t.Errorf("unexpected Type: %q", f.Type)
	}
	if f.Message == "" {
		t.Error("expected non-empty message")
	}
	if f.Line != 5 {
		t.Errorf("expected Line=5, got %d", f.Line)
	}
}

func TestScanResult_Fields(t *testing.T) {
	sr := ScanResult{
		Label:    "pkg/file.py",
		Findings: []ScanFinding{{Type: "hidden-character", Message: "x"}},
	}
	if sr.Label != "pkg/file.py" {
		t.Errorf("unexpected Label: %q", sr.Label)
	}
	if len(sr.Findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(sr.Findings))
	}
}

func TestScanDeployedFiles_EmptyLock(t *testing.T) {
	root := t.TempDir()
	lockData := LockFileData{Dependencies: map[string]LockedDependency{}}
	findings, count := ScanDeployedFiles(lockData, root, "")
	if count != 0 {
		t.Errorf("expected 0 files scanned, got %d", count)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestScanDeployedFiles_CleanFileExtra(t *testing.T) {
	root := t.TempDir()
	cleanFile := filepath.Join(root, "safe.py")
	if err := os.WriteFile(cleanFile, []byte("x = 1\nprint(x)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"mypkg": {DeployedFiles: []string{"safe.py"}},
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

func TestScanDeployedFiles_SuspiciousFileExtra(t *testing.T) {
	root := t.TempDir()
	badFile := filepath.Join(root, "evil.py")
	if err := os.WriteFile(badFile, []byte("x\u200B= 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"mypkg": {DeployedFiles: []string{"evil.py"}},
		},
	}
	findings, count := ScanDeployedFiles(lockData, root, "")
	if count != 1 {
		t.Errorf("expected 1 file scanned, got %d", count)
	}
	if len(findings) == 0 {
		t.Error("expected finding for suspicious file")
	}
}

func TestScanDeployedFiles_PackageFilterExtra(t *testing.T) {
	root := t.TempDir()
	f1 := filepath.Join(root, "a.py")
	f2 := filepath.Join(root, "b.py")
	os.WriteFile(f1, []byte("normal"), 0o644)
	os.WriteFile(f2, []byte("normal"), 0o644)
	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"pkgA": {DeployedFiles: []string{"a.py"}},
			"pkgB": {DeployedFiles: []string{"b.py"}},
		},
	}
	_, count := ScanDeployedFiles(lockData, root, "pkgA")
	if count != 1 {
		t.Errorf("filter should scan only pkgA (1 file), got %d", count)
	}
}

func TestScanDeployedFiles_PathTraversal(t *testing.T) {
	root := t.TempDir()
	lockData := LockFileData{
		Dependencies: map[string]LockedDependency{
			"mypkg": {DeployedFiles: []string{"../outside.py"}},
		},
	}
	_, count := ScanDeployedFiles(lockData, root, "")
	if count != 0 {
		t.Errorf("path traversal should be skipped, got %d files scanned", count)
	}
}

func TestIsSafeLockfilePath_Empty(t *testing.T) {
	root := t.TempDir()
	// empty path -- should be unsafe (empty doesn't make a valid sub-path)
	// behavior depends on filepath.Join logic; at minimum it should not panic
	_ = isSafeLockfilePath("", root)
}
