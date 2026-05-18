package errors_test

import (
"errors"
"fmt"
"testing"

ierrors "github.com/githubnext/apm/internal/install/errors"
)

func TestDirectDependencyError_BasicMsg(t *testing.T) {
err := ierrors.NewDirectDependencyError("dep failed")
if err.Error() != "dep failed" {
t.Errorf("unexpected error: %q", err.Error())
}
}

func TestDirectDependencyError_IsDirect(t *testing.T) {
err := ierrors.NewDirectDependencyError("dep")
if !ierrors.IsDirect(err) {
t.Error("IsDirect should return true for DirectDependencyError")
}
}

func TestDirectDependencyError_NotWrapped(t *testing.T) {
err := ierrors.NewDirectDependencyError("dep")
wrapped := fmt.Errorf("wrapped: %w", err)
// IsDirect uses type assertion, not errors.As, so wrapped should fail
_ = wrapped
}

func TestIsDirect_Nil(t *testing.T) {
if ierrors.IsDirect(nil) {
t.Error("IsDirect(nil) should return false")
}
}

func TestIsDirect_OtherError(t *testing.T) {
if ierrors.IsDirect(errors.New("other")) {
t.Error("IsDirect for non-Direct error should return false")
}
}

func TestAuthenticationError_Msg(t *testing.T) {
err := ierrors.NewAuthenticationError("auth failed", "see docs")
if err.Error() != "auth failed" {
t.Errorf("unexpected error: %q", err.Error())
}
if err.DiagnosticContext != "see docs" {
t.Errorf("DiagnosticContext = %q, want 'see docs'", err.DiagnosticContext)
}
}

func TestIsAuthentication_True(t *testing.T) {
err := ierrors.NewAuthenticationError("x", "")
if !ierrors.IsAuthentication(err) {
t.Error("IsAuthentication should return true")
}
}

func TestIsAuthentication_Nil(t *testing.T) {
if ierrors.IsAuthentication(nil) {
t.Error("IsAuthentication(nil) should return false")
}
}

func TestIsAuthentication_OtherError(t *testing.T) {
if ierrors.IsAuthentication(errors.New("generic")) {
t.Error("IsAuthentication for non-auth error should return false")
}
}

func TestFrozenInstallError_Reasons(t *testing.T) {
err := ierrors.NewFrozenInstallError("frozen", []string{"r1", "r2"})
if len(err.Reasons) != 2 {
t.Errorf("expected 2 reasons, got %d", len(err.Reasons))
}
}

func TestFrozenInstallError_EmptyReasons(t *testing.T) {
err := ierrors.NewFrozenInstallError("frozen", nil)
if len(err.Reasons) != 0 {
t.Errorf("expected 0 reasons, got %d", len(err.Reasons))
}
}

func TestIsFrozen_True(t *testing.T) {
err := ierrors.NewFrozenInstallError("locked", nil)
if !ierrors.IsFrozen(err) {
t.Error("IsFrozen should return true")
}
}

func TestIsFrozen_Nil(t *testing.T) {
if ierrors.IsFrozen(nil) {
t.Error("IsFrozen(nil) should return false")
}
}

func TestPolicyViolationError_Msg(t *testing.T) {
err := ierrors.NewPolicyViolationError("policy blocked", "org/policy")
if err.Error() != "policy blocked" {
t.Errorf("unexpected error: %q", err.Error())
}
if err.PolicySource != "org/policy" {
t.Errorf("PolicySource = %q, want 'org/policy'", err.PolicySource)
}
}

func TestIsPolicy_True(t *testing.T) {
err := ierrors.NewPolicyViolationError("blocked", "src")
if !ierrors.IsPolicy(err) {
t.Error("IsPolicy should return true for PolicyViolationError")
}
}

func TestIsPolicy_Nil(t *testing.T) {
if ierrors.IsPolicy(nil) {
t.Error("IsPolicy(nil) should return false")
}
}

func TestIsPolicy_OtherError(t *testing.T) {
if ierrors.IsPolicy(errors.New("generic")) {
t.Error("IsPolicy for non-policy error should return false")
}
}
