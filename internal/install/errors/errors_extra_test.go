package errors_test

import (
	"errors"
	"testing"

	ierrors "github.com/githubnext/apm/internal/install/errors"
)

func TestDirectDependencyError_EmptyMsg(t *testing.T) {
	err := ierrors.NewDirectDependencyError("")
	if err.Msg != "" {
		t.Errorf("expected empty msg, got %q", err.Msg)
	}
	if err.Error() != "" {
		t.Errorf("Error() should return empty, got %q", err.Error())
	}
}

func TestAuthenticationError_EmptyContext(t *testing.T) {
	err := ierrors.NewAuthenticationError("auth failed", "")
	if err.DiagnosticContext != "" {
		t.Errorf("expected empty DiagnosticContext, got %q", err.DiagnosticContext)
	}
}

func TestAuthenticationError_AsTarget(t *testing.T) {
	err := ierrors.NewAuthenticationError("x", "ctx")
	var target *ierrors.AuthenticationError
	if !errors.As(err, &target) {
		t.Fatal("errors.As failed for AuthenticationError")
	}
	if target.Msg != "x" {
		t.Errorf("Msg = %q, want x", target.Msg)
	}
}

func TestFrozenInstallError_AsTarget(t *testing.T) {
	err := ierrors.NewFrozenInstallError("frozen", []string{"r1"})
	var target *ierrors.FrozenInstallError
	if !errors.As(err, &target) {
		t.Fatal("errors.As failed for FrozenInstallError")
	}
	if len(target.Reasons) != 1 || target.Reasons[0] != "r1" {
		t.Errorf("Reasons = %v, want [r1]", target.Reasons)
	}
}

func TestPolicyViolationError_AsTarget(t *testing.T) {
	err := ierrors.NewPolicyViolationError("blocked", "src")
	var target *ierrors.PolicyViolationError
	if !errors.As(err, &target) {
		t.Fatal("errors.As failed for PolicyViolationError")
	}
	if target.PolicySource != "src" {
		t.Errorf("PolicySource = %q, want src", target.PolicySource)
	}
}

func TestIsHelpers_CrossTypes(t *testing.T) {
	direct := ierrors.NewDirectDependencyError("d")
	auth := ierrors.NewAuthenticationError("a", "")
	frozen := ierrors.NewFrozenInstallError("f", nil)
	policy := ierrors.NewPolicyViolationError("p", "src")

	// Negative checks
	if ierrors.IsAuthentication(direct) {
		t.Error("direct is not auth")
	}
	if ierrors.IsFrozen(auth) {
		t.Error("auth is not frozen")
	}
	if ierrors.IsDirect(policy) {
		t.Error("policy is not direct")
	}
	if ierrors.IsPolicy(frozen) {
		t.Error("frozen is not policy")
	}
}

func TestFrozenInstallError_SingleReason(t *testing.T) {
	err := ierrors.NewFrozenInstallError("frozen", []string{"only reason"})
	if len(err.Reasons) != 1 {
		t.Fatalf("expected 1 reason, got %d", len(err.Reasons))
	}
	if err.Reasons[0] != "only reason" {
		t.Errorf("reason = %q, want 'only reason'", err.Reasons[0])
	}
}

func TestIsPolicy_DirectError(t *testing.T) {
	err := ierrors.NewDirectDependencyError("x")
	if ierrors.IsPolicy(err) {
		t.Error("direct error should not be policy error")
	}
}
