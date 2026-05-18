package securityscan_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/securityscan"
)

func TestPreDeploySecurityScan_SubdirectoryFiles(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "nested.txt"), []byte("clean content"), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "nested-pkg", false)
	if !ok {
		t.Error("expected ok=true for clean nested files")
	}
	if result.HasFindings {
		t.Error("expected no findings for clean nested file")
	}
}

func TestPreDeploySecurityScan_BidiOverride(t *testing.T) {
	dir := t.TempDir()
	// Unicode bidi override (U+202E)
	content := "normal\u202Etext"
	if err := os.WriteFile(filepath.Join(dir, "bidi.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "bidi-pkg", false)
	if ok {
		t.Error("expected ok=false for bidi override character")
	}
	if !result.ShouldBlock {
		t.Error("expected ShouldBlock=true for bidi override")
	}
}

func TestPreDeploySecurityScan_ZeroWidthJoiner(t *testing.T) {
	dir := t.TempDir()
	// Zero-width joiner (U+200D)
	content := "A\u200Dtext"
	if err := os.WriteFile(filepath.Join(dir, "zwj.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "zwj-pkg", false)
	_ = ok
	_ = result
}

func TestPreDeploySecurityScan_MultipleCleanFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.txt", "b.md", "c.go"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("safe content "+name), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "multi-pkg", false)
	if !ok {
		t.Error("expected ok=true for multiple clean files")
	}
	if result.FilesScanned != 3 {
		t.Errorf("expected 3 files scanned, got %d", result.FilesScanned)
	}
}

func TestPreDeploySecurityScan_FindingHasFile(t *testing.T) {
	dir := t.TempDir()
	content := "hidden\u200Bchar"
	if err := os.WriteFile(filepath.Join(dir, "target.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, result := securityscan.PreDeploySecurityScan(dir, "pkg", false)
	if len(result.Findings) == 0 {
		t.Fatal("expected at least one finding")
	}
	if result.Findings[0].FilePath == "" {
		t.Error("finding should have non-empty FilePath field")
	}
}

func TestPreDeploySecurityScan_PackageName(t *testing.T) {
	dir := t.TempDir()
	ok, result := securityscan.PreDeploySecurityScan(dir, "my-special-package", false)
	if !ok {
		t.Error("expected ok=true for empty dir")
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
