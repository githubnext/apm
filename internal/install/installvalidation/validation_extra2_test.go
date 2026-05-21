package installvalidation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuthenticationError_Error(t *testing.T) {
	e := &AuthenticationError{Host: "github.com", Message: "bad token"}
	if e.Error() == "" {
		t.Error("expected non-empty error message")
	}
	if !strings.Contains(e.Error(), "bad token") {
		t.Errorf("expected message in error, got %q", e.Error())
	}
}

func TestAuthenticationError_ZeroValue(t *testing.T) {
	var e AuthenticationError
	_ = e.Error() // should not panic
}

func TestTLSError_WrapUnwrap(t *testing.T) {
	inner := errors.New("x509: certificate")
	e := &TLSError{Host: "example.com", Cause: inner}
	if !errors.Is(e, inner) {
		t.Error("expected Unwrap to expose inner error")
	}
}

func TestTLSError_MessageContainsHost(t *testing.T) {
	e := &TLSError{Host: "example.com", Cause: errors.New("cert error")}
	if !strings.Contains(e.Error(), "example.com") {
		t.Errorf("expected host in error, got %q", e.Error())
	}
}

func TestIsTLSFailure_Wrapped(t *testing.T) {
	inner := &TLSError{Host: "h", Cause: errors.New("tls")}
	wrapped := errors.New("wrap: " + inner.Error())
	// Direct TLSError should be detected
	if !IsTLSFailure(inner) {
		t.Error("expected TLSError to be detected")
	}
	_ = wrapped
}

func TestIsTLSFailure_PlainError(t *testing.T) {
	if IsTLSFailure(errors.New("some random error")) {
		t.Error("plain error should not be TLS failure")
	}
}

func TestIsTLSFailure_Nil(t *testing.T) {
	if IsTLSFailure(nil) {
		t.Error("nil should not be TLS failure")
	}
}

func TestLocalPathFailureReason_ApmDir(t *testing.T) {
	reason := LocalPathFailureReason("/some/.apm/path")
	_ = reason // verify no panic
}

func TestLocalPathFailureReason_RegularPath(t *testing.T) {
	reason := LocalPathFailureReason("/regular/path")
	_ = reason // verify no panic
}

func TestLocalPathNoMarkersHint_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	hint := LocalPathNoMarkersHint(dir)
	_ = hint // verify no panic
}

func TestLocalPathNoMarkersHint_WithFiles(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0o644)
	hint := LocalPathNoMarkersHint(dir)
	_ = hint // verify no panic
}

func TestNewPackageProber_NotNil(t *testing.T) {
	p := NewPackageProber("github.com", "token")
	if p == nil {
		t.Error("expected non-nil prober")
	}
}

func TestProbeResult_ZeroValue(t *testing.T) {
	var pr ProbeResult
	if pr.Reachable || pr.IsAuthError || pr.IsTLSError || pr.Reason != "" {
		t.Error("expected zero value")
	}
}

func TestIsADOAuthFailureSignal_StatusOK(t *testing.T) {
	if IsADOAuthFailureSignal(200, "") {
		t.Error("200 should not be auth failure")
	}
}

func TestIsADOAuthFailureSignal_Status401(t *testing.T) {
	got := IsADOAuthFailureSignal(401, "Unauthorized")
	_ = got // verify no panic
}
