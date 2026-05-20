package policychecks

import (
	"strings"
	"testing"
)

func TestCheckResult_HasFailures_TrueWhenNotPassed(t *testing.T) {
	r := CheckResult{Passed: false}
	if !r.HasFailures() {
		t.Error("expected HasFailures=true when Passed=false")
	}
}

func TestCheckResult_HasFailures_FalseWhenPassed(t *testing.T) {
	r := CheckResult{Passed: true}
	if r.HasFailures() {
		t.Error("expected HasFailures=false when Passed=true")
	}
}

func TestCIAuditResult_RenderSummary_WithFailures(t *testing.T) {
	r := CIAuditResult{Checks: []CheckResult{
		{Name: "allowlist", Passed: false, Message: "blocked"},
	}}
	summary := r.RenderSummary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestCheckDependencyAllowlist_EmptyDeps(t *testing.T) {
	r := CheckDependencyAllowlist(nil, DependencyPolicy{Allow: []string{"owner/*"}})
	if !r.Passed {
		t.Error("expected pass for empty dep list")
	}
}

func TestCheckDependencyDenylist_EmptyPolicy(t *testing.T) {
	deps := []DependencyRef{{CanonicalString: "owner/repo"}}
	r := CheckDependencyDenylist(deps, DependencyPolicy{})
	if !r.Passed {
		t.Error("expected pass when no deny list")
	}
}

func TestCheckRequiredPackages_NotPresent(t *testing.T) {
	r := CheckRequiredPackages(nil, DependencyPolicy{Require: []string{"owner/required"}})
	if r.Passed {
		t.Error("expected failure when required package not present")
	}
}

func TestCheckRequiredPackages_Present(t *testing.T) {
	deps := []DependencyRef{{CanonicalString: "owner/required"}}
	r := CheckRequiredPackages(deps, DependencyPolicy{Require: []string{"owner/required"}})
	if !r.Passed {
		t.Errorf("expected pass, got: %q", r.Message)
	}
}

func TestCheckCompilationTarget_Match(t *testing.T) {
	r := CheckCompilationTarget("vscode", "vscode")
	if !r.Passed {
		t.Errorf("expected pass for matching targets, got: %q", r.Message)
	}
}

func TestCheckCompilationTarget_NoRequired(t *testing.T) {
	r := CheckCompilationTarget("vscode", "")
	if !r.Passed {
		t.Error("expected pass when no required target")
	}
}

func TestCheckCompilationTarget_Mismatch(t *testing.T) {
	r := CheckCompilationTarget("cursor", "vscode")
	if r.Passed {
		t.Error("expected failure for mismatched targets")
	}
	if !strings.Contains(r.Message, "cursor") && !strings.Contains(r.Message, "vscode") {
		t.Errorf("expected target names in message, got: %q", r.Message)
	}
}

func TestDependencyRef_IsLocal(t *testing.T) {
	ref := DependencyRef{CanonicalString: "/local/path", IsLocal: true}
	if !ref.IsLocal {
		t.Error("expected IsLocal=true")
	}
}

func TestDependencyPolicy_ZeroValue(t *testing.T) {
	var dp DependencyPolicy
	if len(dp.Allow) != 0 || len(dp.Deny) != 0 || len(dp.Require) != 0 {
		t.Error("expected zero value")
	}
}

func TestCheckExtensionsPresent_OneRequired(t *testing.T) {
	present := map[string]bool{"copilot": true}
	r := CheckExtensionsPresent(present, []string{"copilot"})
	if !r.Passed {
		t.Errorf("expected pass, got: %q", r.Message)
	}
}

func TestCheckExtensionsPresent_Missing(t *testing.T) {
	r := CheckExtensionsPresent(map[string]bool{}, []string{"copilot"})
	if r.Passed {
		t.Error("expected failure when required extension missing")
	}
}
