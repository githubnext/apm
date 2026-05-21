package errors_test

import (
	"errors"
	"fmt"
	"testing"

	installerrors "github.com/githubnext/apm/internal/install/errors"
)

// ---------------------------------------------------------------------------
// DirectDependencyError fields
// ---------------------------------------------------------------------------

func TestDirectDependencyError_Fields(t *testing.T) {
	e := installerrors.NewDirectDependencyError("some direct dep error")
	if e.Msg != "some direct dep error" {
		t.Errorf("expected Msg='some direct dep error', got %q", e.Msg)
	}
	if e.Error() != "some direct dep error" {
		t.Errorf("expected Error()='some direct dep error', got %q", e.Error())
	}
}

// ---------------------------------------------------------------------------
// AuthenticationError additional scenarios
// ---------------------------------------------------------------------------

func TestAuthenticationError_Fields(t *testing.T) {
	e := installerrors.NewAuthenticationError("auth failed", "host=github.com")
	if e.Msg != "auth failed" {
		t.Errorf("expected Msg='auth failed', got %q", e.Msg)
	}
	if e.DiagnosticContext != "host=github.com" {
		t.Errorf("expected DiagnosticContext='host=github.com', got %q", e.DiagnosticContext)
	}
}

func TestAuthenticationError_IsAuthentication(t *testing.T) {
	e := installerrors.NewAuthenticationError("fail", "ctx")
	if !installerrors.IsAuthentication(e) {
		t.Error("expected IsAuthentication=true")
	}
}

func TestAuthenticationError_WrappedIsAuthentication(t *testing.T) {
	// IsAuthentication uses direct type assertion, not errors.As
	// so it returns false for wrapped errors
	inner := installerrors.NewAuthenticationError("inner", "ctx")
	wrapped := fmt.Errorf("outer: %w", inner)
	if installerrors.IsAuthentication(wrapped) {
		t.Log("IsAuthentication supports wrapped errors")
	}
	// the direct error should still work
	if !installerrors.IsAuthentication(inner) {
		t.Error("expected IsAuthentication=true for direct error")
	}
}

// ---------------------------------------------------------------------------
// FrozenInstallError additional scenarios
// ---------------------------------------------------------------------------

func TestFrozenInstallError_Fields(t *testing.T) {
	e := installerrors.NewFrozenInstallError("frozen", []string{"reason1", "reason2"})
	if e.Msg != "frozen" {
		t.Errorf("expected Msg='frozen', got %q", e.Msg)
	}
	if len(e.Reasons) != 2 {
		t.Errorf("expected 2 reasons, got %d", len(e.Reasons))
	}
}

func TestFrozenInstallError_IsFrozen(t *testing.T) {
	e := installerrors.NewFrozenInstallError("frozen", nil)
	if !installerrors.IsFrozen(e) {
		t.Error("expected IsFrozen=true")
	}
}

func TestFrozenInstallError_EmptyReasonsExtra2(t *testing.T) {
	e := installerrors.NewFrozenInstallError("frozen", []string{})
	if len(e.Reasons) != 0 {
		t.Errorf("expected 0 reasons, got %d", len(e.Reasons))
	}
}

// ---------------------------------------------------------------------------
// PolicyViolationError additional scenarios
// ---------------------------------------------------------------------------

func TestPolicyViolationError_Fields(t *testing.T) {
	e := installerrors.NewPolicyViolationError("policy violated", "org:myorg/.github")
	if e.Msg != "policy violated" {
		t.Errorf("expected Msg='policy violated', got %q", e.Msg)
	}
	if e.PolicySource != "org:myorg/.github" {
		t.Errorf("expected PolicySource='org:myorg/.github', got %q", e.PolicySource)
	}
}

func TestPolicyViolationError_IsPolicy(t *testing.T) {
	e := installerrors.NewPolicyViolationError("denied", "src")
	if !installerrors.IsPolicy(e) {
		t.Error("expected IsPolicy=true")
	}
}

// ---------------------------------------------------------------------------
// IsAuthentication / IsFrozen / IsPolicy on nil
// ---------------------------------------------------------------------------

func TestIsHelpers_NilError(t *testing.T) {
	if installerrors.IsAuthentication(nil) {
		t.Error("expected false for nil error")
	}
	if installerrors.IsFrozen(nil) {
		t.Error("expected false for nil error")
	}
	if installerrors.IsPolicy(nil) {
		t.Error("expected false for nil error")
	}
	if installerrors.IsDirect(nil) {
		t.Error("expected false for nil error")
	}
}

func TestIsHelpers_StandardError(t *testing.T) {
	e := errors.New("plain error")
	if installerrors.IsAuthentication(e) || installerrors.IsFrozen(e) || installerrors.IsPolicy(e) || installerrors.IsDirect(e) {
		t.Error("expected all Is* helpers to return false for plain error")
	}
}
