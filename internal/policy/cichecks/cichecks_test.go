package cichecks_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/policy/cichecks"
)

func TestCheckManifestParse(t *testing.T) {
	r := cichecks.CheckManifestParse()
	if !r.Passed {
		t.Error("CheckManifestParse should return passed result")
	}
	if r.HasFailures() {
		t.Error("HasFailures should be false for passing check")
	}
}

func TestCheckManifestParseFailed(t *testing.T) {
	r := cichecks.CheckManifestParseFailed(errors.New("yaml: line 3"))
	if r.Passed {
		t.Error("CheckManifestParseFailed should return failed result")
	}
	if !r.HasFailures() {
		t.Error("HasFailures should be true for failed check")
	}
	if !strings.Contains(r.Message, "yaml: line 3") {
		t.Errorf("Message should contain error text, got %q", r.Message)
	}
}

func TestCheckLockfileExistsNoDeps(t *testing.T) {
	r := cichecks.CheckLockfileExists(t.TempDir(), false)
	if !r.Passed {
		t.Error("no deps => no lockfile required => should pass")
	}
}

func TestCheckLockfileExistsMissingFile(t *testing.T) {
	r := cichecks.CheckLockfileExists(t.TempDir(), true)
	if r.Passed {
		t.Error("missing lockfile with deps should fail")
	}
}

func TestCheckLockfileSyncInSync(t *testing.T) {
	keys := map[string]bool{"a": true, "b": true}
	r := cichecks.CheckLockfileSync(keys, keys)
	if !r.Passed {
		t.Errorf("identical key sets should pass, msg=%q", r.Message)
	}
}

func TestCheckLockfileSyncOutOfSync(t *testing.T) {
	manifest := map[string]bool{"a": true, "b": true}
	lockfile := map[string]bool{"a": true}
	r := cichecks.CheckLockfileSync(manifest, lockfile)
	if r.Passed {
		t.Error("mismatched key sets should fail")
	}
}

func TestCheckRefConsistencyNoDeps(t *testing.T) {
	r := cichecks.CheckRefConsistency(nil)
	if !r.Passed {
		t.Error("no deps should pass ref consistency")
	}
}

func TestCIAuditResultHasFailures(t *testing.T) {
	result := cichecks.CIAuditResult{
		Checks: []cichecks.CheckResult{
			{Name: "a", Passed: true, Message: "ok"},
			{Name: "b", Passed: false, Message: "failed"},
		},
	}
	if !result.HasFailures() {
		t.Error("CIAuditResult.HasFailures should be true")
	}
}

func TestCIAuditResultNoFailures(t *testing.T) {
	result := cichecks.CIAuditResult{
		Checks: []cichecks.CheckResult{
			{Name: "a", Passed: true, Message: "ok"},
		},
	}
	if result.HasFailures() {
		t.Error("CIAuditResult.HasFailures should be false")
	}
}

func TestCIAuditResultRenderSummary(t *testing.T) {
	result := cichecks.CIAuditResult{
		Checks: []cichecks.CheckResult{
			{Name: "lockfile", Passed: true, Message: "in sync"},
			{Name: "refs", Passed: false, Message: "mismatch"},
		},
	}
	summary := result.RenderSummary()
	if !strings.Contains(summary, "[+]") {
		t.Error("summary should contain [+] for passing check")
	}
	if !strings.Contains(summary, "[x]") {
		t.Error("summary should contain [x] for failing check")
	}
	if !strings.Contains(summary, "lockfile") {
		t.Error("summary should contain check name")
	}
}

func TestCheckDriftFindingsNone(t *testing.T) {
	r := cichecks.CheckDriftFindings(nil)
	if !r.Passed {
		t.Error("no drift findings should pass")
	}
}

func TestCheckDriftFindingsWithFindings(t *testing.T) {
	findings := []cichecks.DriftFinding{
		{DepKey: "pkg/foo", FilePath: "some/file.md", Reason: "modified"},
	}
	r := cichecks.CheckDriftFindings(findings)
	if r.Passed {
		t.Error("drift findings should cause check to fail")
	}
}
