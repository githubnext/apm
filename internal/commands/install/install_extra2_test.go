package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallMode_AllValues(t *testing.T) {
	modes := []InstallMode{InstallModeAll, InstallModePrimitives, InstallModeClients}
	seen := map[InstallMode]bool{}
	for _, m := range modes {
		if m == "" {
			t.Error("InstallMode constant must not be empty")
		}
		if seen[m] {
			t.Errorf("duplicate InstallMode: %q", m)
		}
		seen[m] = true
	}
}

func TestInstallResult_Warnings(t *testing.T) {
	r := InstallResult{
		Warnings: []string{"warn1", "warn2"},
		Errors:   []string{"err1"},
	}
	if len(r.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(r.Warnings))
	}
	if len(r.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(r.Errors))
	}
}

func TestInstallResult_FilesWritten(t *testing.T) {
	r := InstallResult{
		FilesWritten:      []string{"/a/b.txt", "/c/d.txt"},
		PackagesInstalled: 2,
		LockfileUpdated:   true,
	}
	if len(r.FilesWritten) != 2 {
		t.Errorf("expected 2 files, got %d", len(r.FilesWritten))
	}
	if !r.LockfileUpdated {
		t.Error("expected LockfileUpdated=true")
	}
}

func TestLockEntry_Fields(t *testing.T) {
	e := LockEntry{
		Name:   "mypkg",
		Ref:    "v1.2.3",
		Commit: "abc123",
		Source: "github.com/org/repo",
		Hash:   "sha256:deadbeef",
	}
	if e.Name != "mypkg" {
		t.Errorf("Name = %q", e.Name)
	}
	if e.Ref != "v1.2.3" {
		t.Errorf("Ref = %q", e.Ref)
	}
	if e.Hash != "sha256:deadbeef" {
		t.Errorf("Hash = %q", e.Hash)
	}
}

func TestSecurityScanResult_NoFindings(t *testing.T) {
	r := &SecurityScanResult{Package: "safe-pkg"}
	if r.Package != "safe-pkg" {
		t.Errorf("Package = %q", r.Package)
	}
	if r.Blocked {
		t.Error("expected Blocked=false for empty findings")
	}
	if len(r.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(r.Findings))
	}
}

func TestSecurityScanResult_WithFindings(t *testing.T) {
	r := &SecurityScanResult{
		Package:  "risky-pkg",
		Findings: []string{".env file found", "id_rsa found"},
		Blocked:  true,
	}
	if !r.Blocked {
		t.Error("expected Blocked=true")
	}
	if len(r.Findings) != 2 {
		t.Errorf("expected 2 findings, got %d", len(r.Findings))
	}
}

func TestRunPreDeploySecurityScan_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	result, err := RunPreDeploySecurityScan(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Blocked {
		t.Error("expected not blocked for empty directory")
	}
	if len(result.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d: %v", len(result.Findings), result.Findings)
	}
}

func TestRunPreDeploySecurityScan_RiskyFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("SECRET=abc"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := RunPreDeploySecurityScan(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Blocked {
		t.Error("expected blocked when .env is present")
	}
	if len(result.Findings) == 0 {
		t.Error("expected at least 1 finding")
	}
}

func TestFormatInstallSummary_Empty(t *testing.T) {
	r := &InstallResult{}
	s := FormatInstallSummary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestFormatInstallSummary_WithPackages(t *testing.T) {
	r := &InstallResult{PackagesInstalled: 3, PackagesSkipped: 1, LockfileUpdated: true}
	s := FormatInstallSummary(r)
	if !strings.Contains(s, "3") {
		t.Errorf("expected '3' in summary: %q", s)
	}
}

func TestDependencyEntry_AllFields(t *testing.T) {
	d := DependencyEntry{
		Name: "my-pkg",
		Ref:  "main",
		Host: "github.com",
		Org:  "myorg",
		Repo: "myrepo",
	}
	if d.Org != "myorg" {
		t.Errorf("Org = %q", d.Org)
	}
	if d.Repo != "myrepo" {
		t.Errorf("Repo = %q", d.Repo)
	}
}
