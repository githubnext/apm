package policychecks_test

import (
	"os"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/policy/policychecks"
)

func TestCIAuditResult_Empty(t *testing.T) {
	r := policychecks.CIAuditResult{}
	if r.HasFailures() {
		t.Error("empty audit result should not have failures")
	}
	s := r.RenderSummary()
	if s != "" {
		t.Errorf("expected empty summary, got %q", s)
	}
}

func TestCheckResult_Name(t *testing.T) {
	r := policychecks.CheckResult{Name: "my-check", Passed: true, Message: "ok"}
	if r.Name != "my-check" {
		t.Errorf("Name = %q", r.Name)
	}
}

func TestRenderSummary_MultipleDetails(t *testing.T) {
	r := policychecks.CIAuditResult{
		Checks: []policychecks.CheckResult{
			{Name: "chk", Passed: false, Message: "fail", Details: []string{"d1", "d2", "d3"}},
		},
	}
	s := r.RenderSummary()
	if !strings.Contains(s, "d1") || !strings.Contains(s, "d2") || !strings.Contains(s, "d3") {
		t.Errorf("expected all details in summary, got %q", s)
	}
}

func TestCheckDependencyAllowlist_MultipleAllowed(t *testing.T) {
	deps := []policychecks.DependencyRef{
		{CanonicalString: "github.com/owner/a"},
		{CanonicalString: "github.com/owner/b"},
	}
	policy := policychecks.DependencyPolicy{Allow: []string{"github.com/owner/*"}}
	r := policychecks.CheckDependencyAllowlist(deps, policy)
	if !r.Passed {
		t.Errorf("all in allowlist should pass: %s", r.Message)
	}
}

func TestCheckDependencyAllowlist_MultiplePatterns(t *testing.T) {
	deps := []policychecks.DependencyRef{
		{CanonicalString: "github.com/owner/a"},
		{CanonicalString: "gitlab.com/other/b"},
	}
	policy := policychecks.DependencyPolicy{Allow: []string{"github.com/owner/*", "gitlab.com/other/*"}}
	r := policychecks.CheckDependencyAllowlist(deps, policy)
	if !r.Passed {
		t.Errorf("all in multi-pattern allowlist should pass: %s", r.Message)
	}
}

func TestCheckDependencyDenylist_NotInDeny(t *testing.T) {
	deps := []policychecks.DependencyRef{{CanonicalString: "github.com/safe/pkg"}}
	policy := policychecks.DependencyPolicy{Deny: []string{"github.com/bad/*"}}
	r := policychecks.CheckDependencyDenylist(deps, policy)
	if !r.Passed {
		t.Errorf("safe dep should pass deny check: %s", r.Message)
	}
}

func TestCheckDependencyDenylist_EmptyDeps(t *testing.T) {
	policy := policychecks.DependencyPolicy{Deny: []string{"github.com/bad/*"}}
	r := policychecks.CheckDependencyDenylist(nil, policy)
	if !r.Passed {
		t.Errorf("empty deps should pass deny check: %s", r.Message)
	}
}

func TestCheckRequiredPackages_NoneRequired(t *testing.T) {
	r := policychecks.CheckRequiredPackages(nil, policychecks.DependencyPolicy{})
	if !r.Passed {
		t.Error("no required packages should pass")
	}
}

func TestCheckRequiredPackages_MultiplePresent(t *testing.T) {
	deps := []policychecks.DependencyRef{
		{CanonicalString: "github.com/owner/a"},
		{CanonicalString: "github.com/owner/b"},
	}
	policy := policychecks.DependencyPolicy{Require: []string{"github.com/owner/a", "github.com/owner/b"}}
	r := policychecks.CheckRequiredPackages(deps, policy)
	if !r.Passed {
		t.Errorf("all required present should pass: %s", r.Message)
	}
}

func TestCheckCompilationTarget_EmptyActual(t *testing.T) {
	r := policychecks.CheckCompilationTarget("", "vscode")
	if r.Passed {
		t.Error("empty actual should fail when required is set")
	}
}

func TestCheckExtensionsPresent_EmptyRequired(t *testing.T) {
	r := policychecks.CheckExtensionsPresent(map[string]bool{}, nil)
	if !r.Passed {
		t.Error("no required extensions should pass")
	}
}

func TestCheckExtensionsPresent_MultiplePresent(t *testing.T) {
	present := map[string]bool{"x-team": true, "x-env": true}
	r := policychecks.CheckExtensionsPresent(present, []string{"x-team", "x-env"})
	if !r.Passed {
		t.Errorf("all extensions present should pass: %s", r.Message)
	}
}

func TestCheckExtensionsPresent_PartiallyMissing(t *testing.T) {
	present := map[string]bool{"x-team": true}
	r := policychecks.CheckExtensionsPresent(present, []string{"x-team", "x-required"})
	if r.Passed {
		t.Error("missing extension should fail")
	}
	if len(r.Details) != 1 || r.Details[0] != "x-required" {
		t.Errorf("expected 'x-required' in details, got %v", r.Details)
	}
}

func TestLoadRawApmYML_Missing(t *testing.T) {
	dir := t.TempDir()
	result := policychecks.LoadRawApmYML(dir)
	if result != nil {
		t.Errorf("expected nil for missing apm.yml, got %v", result)
	}
}

func TestLoadRawApmYML_Present(t *testing.T) {
	dir := t.TempDir()
	content := "target: vscode\nextensions:\n  x-team: true\n"
	if err := os.WriteFile(dir+"/apm.yml", []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	result := policychecks.LoadRawApmYML(dir)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["target"] != "vscode" {
		t.Errorf("expected target=vscode, got %v", result["target"])
	}
}

func TestLoadRawApmYML_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	content := "# this is a comment\ntarget: cursor\n"
	if err := os.WriteFile(dir+"/apm.yml", []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	result := policychecks.LoadRawApmYML(dir)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if _, ok := result["target"]; !ok {
		t.Error("target should be parsed")
	}
}
