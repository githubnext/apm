package errors_test

import (
	"errors"
	"testing"

	ierrors "github.com/githubnext/apm/internal/install/errors"
)

func TestDirectDependencyError(t *testing.T) {
	err := ierrors.NewDirectDependencyError("dep failed")
	if err.Msg != "dep failed" {
		t.Fatalf("expected 'dep failed', got %q", err.Msg)
	}
	if err.Error() != "dep failed" {
		t.Fatalf("Error() mismatch")
	}
	var target *ierrors.DirectDependencyError
	if !errors.As(err, &target) {
		t.Fatal("errors.As failed for DirectDependencyError")
	}
}

func TestAuthenticationError(t *testing.T) {
	err := ierrors.NewAuthenticationError("auth failed", "ctx detail")
	if err.Msg != "auth failed" {
		t.Fatalf("expected 'auth failed', got %q", err.Msg)
	}
	if err.DiagnosticContext != "ctx detail" {
		t.Fatalf("DiagnosticContext mismatch")
	}
	if err.Error() != "auth failed" {
		t.Fatalf("Error() mismatch")
	}
}

func TestFrozenInstallError(t *testing.T) {
	reasons := []string{"lockfile missing", "out of sync"}
	err := ierrors.NewFrozenInstallError("frozen", reasons)
	if err.Msg != "frozen" {
		t.Fatalf("expected 'frozen', got %q", err.Msg)
	}
	if len(err.Reasons) != 2 {
		t.Fatalf("expected 2 reasons, got %d", len(err.Reasons))
	}
	// Mutation of original slice should not affect stored copy.
	reasons[0] = "mutated"
	if err.Reasons[0] == "mutated" {
		t.Fatal("Reasons slice was not copied")
	}
}

func TestFrozenInstallErrorNilReasons(t *testing.T) {
	err := ierrors.NewFrozenInstallError("frozen", nil)
	if len(err.Reasons) != 0 {
		t.Fatalf("expected empty reasons, got %v", err.Reasons)
	}
}

func TestPolicyViolationError(t *testing.T) {
	err := ierrors.NewPolicyViolationError("policy blocked", "org.yml")
	if err.Msg != "policy blocked" {
		t.Fatalf("expected 'policy blocked', got %q", err.Msg)
	}
	if err.PolicySource != "org.yml" {
		t.Fatalf("PolicySource mismatch")
	}
}

func TestIsHelpers(t *testing.T) {
	direct := ierrors.NewDirectDependencyError("x")
	auth := ierrors.NewAuthenticationError("x", "ctx")
	frozen := ierrors.NewFrozenInstallError("x", nil)
	policy := ierrors.NewPolicyViolationError("x", "src")

	if !ierrors.IsDirect(direct) {
		t.Fatal("IsDirect false")
	}
	if ierrors.IsDirect(auth) {
		t.Fatal("IsDirect true for auth")
	}
	if !ierrors.IsAuthentication(auth) {
		t.Fatal("IsAuthentication false")
	}
	if !ierrors.IsFrozen(frozen) {
		t.Fatal("IsFrozen false")
	}
	if !ierrors.IsPolicy(policy) {
		t.Fatal("IsPolicy false")
	}
	if ierrors.IsPolicy(nil) {
		t.Fatal("IsPolicy(nil) should be false")
	}
}
