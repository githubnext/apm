package securityscan_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/securityscan"
)

func TestPreDeploySecurityScan_NoFindings(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "clean.txt"), []byte("normal ascii text"), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "my-pkg", false)
	if !ok {
		t.Error("expected ok=true for clean directory")
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if result.HasFindings {
		t.Error("expected HasFindings=false for clean directory")
	}
	if result.FilesScanned != 1 {
		t.Errorf("expected FilesScanned=1, got %d", result.FilesScanned)
	}
}

func TestPreDeploySecurityScan_WithHiddenChar(t *testing.T) {
	dir := t.TempDir()
	malicious := "normal text \u200B end"
	if err := os.WriteFile(filepath.Join(dir, "bad.txt"), []byte(malicious), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "evil-pkg", false)
	if ok {
		t.Error("expected ok=false for directory with hidden characters")
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if !result.HasFindings {
		t.Error("expected HasFindings=true")
	}
	if !result.ShouldBlock {
		t.Error("expected ShouldBlock=true")
	}
	if len(result.Findings) == 0 {
		t.Error("expected at least one finding")
	}
}

func TestPreDeploySecurityScan_ForceOverride(t *testing.T) {
	dir := t.TempDir()
	malicious := "normal text \u200B end"
	if err := os.WriteFile(filepath.Join(dir, "bad.txt"), []byte(malicious), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "evil-pkg", true)
	if !ok {
		t.Error("expected ok=true when force=true even with findings")
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if !result.HasFindings {
		t.Error("expected HasFindings=true even in force mode")
	}
}

func TestPreDeploySecurityScan_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	ok, result := securityscan.PreDeploySecurityScan(dir, "empty-pkg", false)
	if !ok {
		t.Error("expected ok=true for empty directory")
	}
	if result.HasFindings {
		t.Error("expected no findings for empty dir")
	}
	if result.FilesScanned != 0 {
		t.Errorf("expected 0 files scanned, got %d", result.FilesScanned)
	}
}

func TestPreDeploySecurityScan_MultipleHiddenChars(t *testing.T) {
	dir := t.TempDir()
	// Mix different hidden characters
	content := "A\u200Btext\u202Emore"
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, result := securityscan.PreDeploySecurityScan(dir, "pkg", false)
	if ok {
		t.Error("expected block on hidden chars")
	}
	// Should have at least 2 findings (one per hidden pattern per file)
	if len(result.Findings) < 2 {
		t.Errorf("expected >=2 findings, got %d", len(result.Findings))
	}
}
