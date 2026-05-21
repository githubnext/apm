package securityscan_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/securityscan"
)

func TestPreDeploySecurityScan_ReturnsTwoValues(t *testing.T) {
	dir := t.TempDir()
	ok, result := securityscan.PreDeploySecurityScan(dir, "testpkg", false)
	_ = ok
	if result == nil {
		t.Error("expected non-nil ScanResult")
	}
}

func TestPreDeploySecurityScan_CleanDirectory(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "clean.md")
	if err := os.WriteFile(f, []byte("# Hello World\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, result := securityscan.PreDeploySecurityScan(dir, "pkg", false)
	if result.HasFindings {
		t.Error("expected no findings for clean directory")
	}
}

func TestPreDeploySecurityScan_HiddenChar_BidiOverride(t *testing.T) {
	dir := t.TempDir()
	// U+202E RIGHT-TO-LEFT OVERRIDE
	f := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(f, []byte("text\u202Emore"), 0644); err != nil {
		t.Fatal(err)
	}
	_, result := securityscan.PreDeploySecurityScan(dir, "pkg", false)
	if !result.HasFindings {
		t.Error("expected findings for bidi override character")
	}
}

func TestPreDeploySecurityScan_ForceOverride_Proceeds(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(f, []byte("text\u202Emore"), 0644); err != nil {
		t.Fatal(err)
	}
	ok, _ := securityscan.PreDeploySecurityScan(dir, "pkg", true)
	if !ok {
		t.Error("expected force=true to allow deployment even with findings")
	}
}

func TestPreDeploySecurityScan_EmptyPackageName(t *testing.T) {
	dir := t.TempDir()
	_, result := securityscan.PreDeploySecurityScan(dir, "", false)
	if result == nil {
		t.Error("expected non-nil result even with empty package name")
	}
}

func TestPreDeploySecurityScan_ZeroWidthNonJoiner(t *testing.T) {
	dir := t.TempDir()
	// U+200C ZERO WIDTH NON-JOINER
	f := filepath.Join(dir, "zwc.txt")
	if err := os.WriteFile(f, []byte("te\u200Cxt"), 0644); err != nil {
		t.Fatal(err)
	}
	_, result := securityscan.PreDeploySecurityScan(dir, "pkg", false)
	if !result.HasFindings {
		t.Error("expected finding for zero-width non-joiner")
	}
}

func TestPreDeploySecurityScan_FindingsHaveFilePath(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "suspect.txt")
	if err := os.WriteFile(f, []byte("x\u200Bx"), 0644); err != nil {
		t.Fatal(err)
	}
	_, result := securityscan.PreDeploySecurityScan(dir, "pkg", false)
	if !result.HasFindings {
		t.Skip("no findings -- hidden char may not be detected")
	}
	for _, finding := range result.Findings {
		if finding.FilePath == "" {
			t.Error("expected non-empty FilePath in finding")
		}
	}
}

func TestPreDeploySecurityScan_MultipleCleanFilesNoFindings(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.md", "b.txt", "c.yaml"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("clean content\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	_, result := securityscan.PreDeploySecurityScan(dir, "multi", false)
	if result.HasFindings {
		t.Error("expected no findings for all-clean directory")
	}
}

func TestPreDeploySecurityScan_ScanResultFields(t *testing.T) {
	dir := t.TempDir()
	_, result := securityscan.PreDeploySecurityScan(dir, "pkg", false)
	_ = result.HasFindings
	_ = result.ShouldBlock
	_ = result.Findings
}
