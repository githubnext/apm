package cichecks

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckManifestParse_Passed(t *testing.T) {
	r := CheckManifestParse()
	if !r.Passed {
		t.Error("expected Passed=true for successful manifest parse")
	}
}

func TestCheckManifestParseFailed_NotPassed(t *testing.T) {
	r := CheckManifestParseFailed(errors.New("parse error"))
	if r.Passed {
		t.Error("expected Passed=false for failed manifest parse")
	}
}

func TestCheckManifestParseFailed_MessageContainsError(t *testing.T) {
	r := CheckManifestParseFailed(errors.New("syntax error at line 5"))
	if !strings.Contains(r.Message, "syntax error at line 5") {
		t.Errorf("expected error in message, got %q", r.Message)
	}
}

func TestCheckLockfileExists_NoFile(t *testing.T) {
	dir := t.TempDir()
	r := CheckLockfileExists(dir, true)
	if r.Passed {
		t.Error("expected failure when lockfile missing")
	}
}

func TestCheckLockfileExists_NoDeps(t *testing.T) {
	dir := t.TempDir()
	r := CheckLockfileExists(dir, false)
	if !r.Passed {
		t.Error("expected pass when no deps required")
	}
}

func TestCheckLockfileExists_FilePresent(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "apm.lock.yaml"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	r := CheckLockfileExists(dir, true)
	if !r.Passed {
		t.Errorf("expected pass when lockfile present, got: %q", r.Message)
	}
}

func TestCheckDriftFindings_Empty(t *testing.T) {
	r := CheckDriftFindings(nil)
	if !r.Passed {
		t.Error("expected pass with no drift findings")
	}
}

func TestCheckDriftFindings_WithFinding(t *testing.T) {
	findings := []DriftFinding{{DepKey: "owner/repo", FilePath: "file.py", Reason: "modified"}}
	r := CheckDriftFindings(findings)
	if r.Passed {
		t.Error("expected failure with drift findings")
	}
}

func TestCheckResult_Name(t *testing.T) {
	r := CheckResult{Name: "lockfile_sync", Passed: true}
	if r.Name != "lockfile_sync" {
		t.Errorf("unexpected name: %q", r.Name)
	}
}

func TestCheckResult_Details(t *testing.T) {
	r := CheckResult{Details: []string{"detail1", "detail2"}}
	if len(r.Details) != 2 {
		t.Errorf("expected 2 details, got %d", len(r.Details))
	}
}

func TestLockedDepInfo_WithDeployedFiles(t *testing.T) {
	dep := LockedDepInfo{
		Key:           "owner/repo",
		DeployedFiles: []string{"a.py", "b.py"},
		ContentHash:   "abc123",
	}
	if len(dep.DeployedFiles) != 2 {
		t.Errorf("expected 2 deployed files")
	}
}

func TestCIAuditResult_RenderSummary_AllPassed(t *testing.T) {
	r := CIAuditResult{Checks: []CheckResult{
		{Name: "a", Passed: true},
		{Name: "b", Passed: true},
	}}
	summary := r.RenderSummary()
	_ = summary // verify no panic
}
