package policychecks_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/policy/policychecks"
)

func TestCheckResult_HasFailures(t *testing.T) {
	if (policychecks.CheckResult{Passed: true}).HasFailures() {
		t.Error("Passed=true should not have failures")
	}
	if !(policychecks.CheckResult{Passed: false}).HasFailures() {
		t.Error("Passed=false should have failures")
	}
}

func TestCIAuditResult_HasFailures_AllPassed(t *testing.T) {
	r := policychecks.CIAuditResult{
		Checks: []policychecks.CheckResult{
			{Passed: true},
			{Passed: true},
		},
	}
	if r.HasFailures() {
		t.Error("all passed should not have failures")
	}
}

func TestCIAuditResult_HasFailures_OneFailed(t *testing.T) {
	r := policychecks.CIAuditResult{
		Checks: []policychecks.CheckResult{
			{Passed: true},
			{Passed: false},
		},
	}
	if !r.HasFailures() {
		t.Error("one failed check should report HasFailures=true")
	}
}

func TestRenderSummary(t *testing.T) {
	r := policychecks.CIAuditResult{
		Checks: []policychecks.CheckResult{
			{Name: "a", Passed: true, Message: "ok"},
			{Name: "b", Passed: false, Message: "bad", Details: []string{"detail1"}},
		},
	}
	s := r.RenderSummary()
	if !strings.Contains(s, "[+] a: ok") {
		t.Errorf("expected [+] a: ok in %q", s)
	}
	if !strings.Contains(s, "[x] b: bad") {
		t.Errorf("expected [x] b: bad in %q", s)
	}
	if !strings.Contains(s, "detail1") {
		t.Errorf("expected detail1 in %q", s)
	}
}

func TestCheckDependencyAllowlist_NoPolicy(t *testing.T) {
	r := policychecks.CheckDependencyAllowlist(nil, policychecks.DependencyPolicy{})
	if !r.Passed {
		t.Error("empty allow list should pass")
	}
}

func TestCheckDependencyAllowlist_Pass(t *testing.T) {
	deps := []policychecks.DependencyRef{{CanonicalString: "github.com/owner/repo"}}
	policy := policychecks.DependencyPolicy{Allow: []string{"github.com/owner/*"}}
	r := policychecks.CheckDependencyAllowlist(deps, policy)
	if !r.Passed {
		t.Errorf("expected pass, got: %v", r.Message)
	}
}

func TestCheckDependencyAllowlist_Fail(t *testing.T) {
	deps := []policychecks.DependencyRef{{CanonicalString: "github.com/other/repo"}}
	policy := policychecks.DependencyPolicy{Allow: []string{"github.com/owner/*"}}
	r := policychecks.CheckDependencyAllowlist(deps, policy)
	if r.Passed {
		t.Error("expected failure for dep not in allow list")
	}
}

func TestCheckDependencyAllowlist_LocalSkipped(t *testing.T) {
	deps := []policychecks.DependencyRef{{CanonicalString: "./local", IsLocal: true}}
	policy := policychecks.DependencyPolicy{Allow: []string{"github.com/owner/*"}}
	r := policychecks.CheckDependencyAllowlist(deps, policy)
	if !r.Passed {
		t.Error("local deps should be skipped")
	}
}

func TestCheckDependencyDenylist_NoPolicy(t *testing.T) {
	r := policychecks.CheckDependencyDenylist(nil, policychecks.DependencyPolicy{})
	if !r.Passed {
		t.Error("empty deny list should pass")
	}
}

func TestCheckDependencyDenylist_Denied(t *testing.T) {
	deps := []policychecks.DependencyRef{{CanonicalString: "github.com/bad/pkg"}}
	policy := policychecks.DependencyPolicy{Deny: []string{"github.com/bad/*"}}
	r := policychecks.CheckDependencyDenylist(deps, policy)
	if r.Passed {
		t.Error("expected failure for denied dep")
	}
}

func TestCheckRequiredPackages_Pass(t *testing.T) {
	deps := []policychecks.DependencyRef{{CanonicalString: "github.com/owner/required"}}
	policy := policychecks.DependencyPolicy{Require: []string{"github.com/owner/required"}}
	r := policychecks.CheckRequiredPackages(deps, policy)
	if !r.Passed {
		t.Errorf("expected pass, got: %v", r.Message)
	}
}

func TestCheckRequiredPackages_Missing(t *testing.T) {
	deps := []policychecks.DependencyRef{}
	policy := policychecks.DependencyPolicy{Require: []string{"github.com/owner/required"}}
	r := policychecks.CheckRequiredPackages(deps, policy)
	if r.Passed {
		t.Error("expected failure for missing required package")
	}
}

func TestCheckCompilationTarget_Match(t *testing.T) {
	r := policychecks.CheckCompilationTarget("vscode", "vscode")
	if !r.Passed {
		t.Error("matching target should pass")
	}
}

func TestCheckCompilationTarget_Mismatch(t *testing.T) {
	r := policychecks.CheckCompilationTarget("cursor", "vscode")
	if r.Passed {
		t.Error("mismatched target should fail")
	}
}

func TestCheckCompilationTarget_NoRequirement(t *testing.T) {
	r := policychecks.CheckCompilationTarget("anything", "")
	if !r.Passed {
		t.Error("no required target should pass")
	}
}

func TestCheckExtensionsPresent_Pass(t *testing.T) {
	present := map[string]bool{"x-team": true}
	r := policychecks.CheckExtensionsPresent(present, []string{"x-team"})
	if !r.Passed {
		t.Error("present extension should pass")
	}
}

func TestCheckExtensionsPresent_Missing(t *testing.T) {
	present := map[string]bool{}
	r := policychecks.CheckExtensionsPresent(present, []string{"x-required"})
	if r.Passed {
		t.Error("missing extension should fail")
	}
}
