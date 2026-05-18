package installvalidation_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/installvalidation"
)

func TestAuthenticationError_WithHost(t *testing.T) {
	err := &installvalidation.AuthenticationError{Host: "github.com", Message: "bad token"}
	msg := err.Error()
	if !strings.Contains(msg, "github.com") {
		t.Fatalf("expected host in error, got %q", msg)
	}
	if !strings.Contains(msg, "bad token") {
		t.Fatalf("expected message, got %q", msg)
	}
}

func TestAuthenticationError_NoHost(t *testing.T) {
	err := &installvalidation.AuthenticationError{Message: "no creds"}
	msg := err.Error()
	if !strings.Contains(msg, "no creds") {
		t.Fatalf("expected message, got %q", msg)
	}
}

func TestTLSError_WithHost(t *testing.T) {
	cause := errors.New("x509: cert invalid")
	err := &installvalidation.TLSError{Host: "ghe.example.com", Cause: cause}
	msg := err.Error()
	if !strings.Contains(msg, installvalidation.TLSErrorPrefix) {
		t.Fatalf("expected TLS prefix, got %q", msg)
	}
	if !strings.Contains(msg, "ghe.example.com") {
		t.Fatalf("expected host, got %q", msg)
	}
	if err.Unwrap() != cause {
		t.Fatal("Unwrap should return cause")
	}
}

func TestIsTLSFailure_True(t *testing.T) {
	err := &installvalidation.TLSError{Host: "x", Cause: errors.New("cert")}
	if !installvalidation.IsTLSFailure(err) {
		t.Fatal("IsTLSFailure should be true")
	}
}

func TestIsTLSFailure_False(t *testing.T) {
	if installvalidation.IsTLSFailure(errors.New("other")) {
		t.Fatal("IsTLSFailure should be false for non-TLS error")
	}
	if installvalidation.IsTLSFailure(nil) {
		t.Fatal("IsTLSFailure(nil) should be false")
	}
}

func TestLocalPathFailureReason_Missing(t *testing.T) {
	reason := installvalidation.LocalPathFailureReason("/nonexistent/path/to/pkg")
	if reason == "" {
		t.Fatal("expected a failure reason for missing path")
	}
}

func TestLocalPathNoMarkersHint_EmptyDir(t *testing.T) {
dir := t.TempDir()
hint := installvalidation.LocalPathNoMarkersHint(dir)
if hint != "" {
t.Errorf("expected empty hint for empty dir, got %q", hint)
}
}

func TestLocalPathNoMarkersHint_WithSubpackage(t *testing.T) {
dir := t.TempDir()
sub := dir + "/mypkg"
if err := os.MkdirAll(sub, 0o755); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(sub+"/apm.yml", []byte("name: mypkg\n"), 0o644); err != nil {
t.Fatal(err)
}
hint := installvalidation.LocalPathNoMarkersHint(dir)
if hint == "" {
t.Error("expected a hint for dir with sub-package")
}
}

func TestLocalPathFailureReason_ValidPath(t *testing.T) {
dir := t.TempDir()
if err := os.WriteFile(dir+"/apm.yml", []byte("name: test\n"), 0o644); err != nil {
t.Fatal(err)
}
reason := installvalidation.LocalPathFailureReason(dir)
if reason != "" {
t.Errorf("expected empty reason for valid path, got %q", reason)
}
}

func TestLocalPathFailureReason_NoMarkers(t *testing.T) {
dir := t.TempDir()
reason := installvalidation.LocalPathFailureReason(dir)
if reason == "" {
t.Error("expected failure reason for path with no markers")
}
}

func TestNewPackageProber_Fields(t *testing.T) {
p := installvalidation.NewPackageProber("github.com", "mytoken")
if p == nil {
t.Fatal("NewPackageProber returned nil")
}
if p.Host != "github.com" {
t.Errorf("expected Host=github.com, got %q", p.Host)
}
if p.AuthToken != "mytoken" {
t.Errorf("expected AuthToken=mytoken")
}
if p.Timeout == 0 {
t.Error("expected non-zero timeout")
}
}

func TestProbeResult_Fields(t *testing.T) {
r := installvalidation.ProbeResult{Reachable: true}
if !r.Reachable {
t.Error("Reachable should be true")
}
r2 := installvalidation.ProbeResult{Reachable: false, Reason: "not found", IsAuthError: true}
if r2.Reachable || !r2.IsAuthError {
t.Error("unexpected ProbeResult fields")
}
r3 := installvalidation.ProbeResult{IsTLSError: true, Reason: "tls failed"}
if !r3.IsTLSError {
t.Error("IsTLSError should be true")
}
}

func TestIsADOAuthFailureSignal_Unauthorized(t *testing.T) {
if !installvalidation.IsADOAuthFailureSignal(401, "") {
t.Error("401 should be ADO auth failure")
}
if !installvalidation.IsADOAuthFailureSignal(403, "") {
t.Error("403 should be ADO auth failure")
}
}

func TestIsADOAuthFailureSignal_BodyMatch(t *testing.T) {
if !installvalidation.IsADOAuthFailureSignal(200, "TFS Auth failed") {
t.Error("TFS Auth body should be ADO auth failure")
}
if !installvalidation.IsADOAuthFailureSignal(200, "unauthorized") {
t.Error("unauthorized body should be ADO auth failure")
}
}

func TestIsADOAuthFailureSignal_False(t *testing.T) {
if installvalidation.IsADOAuthFailureSignal(200, "ok response") {
t.Error("200 with ok body should not be ADO auth failure")
}
}

func TestValidatePackageExists_LocalPath(t *testing.T) {
dir := t.TempDir()
if err := os.WriteFile(dir+"/apm.yml", []byte("name: test\n"), 0o644); err != nil {
t.Fatal(err)
}
result := installvalidation.ValidatePackageExists(dir, "github.com", "", false)
if !result.Reachable {
t.Errorf("expected Reachable=true for local path with apm.yml, got: %q", result.Reason)
}
}

func TestValidatePackageExists_InvalidSpec(t *testing.T) {
result := installvalidation.ValidatePackageExists("notapath", "github.com", "", false)
if result.Reachable {
t.Error("expected Reachable=false for invalid spec")
}
}
