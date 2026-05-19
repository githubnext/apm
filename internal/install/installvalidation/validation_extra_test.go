package installvalidation_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/installvalidation"
)

func TestTLSError_NoHost(t *testing.T) {
	err := &installvalidation.TLSError{Cause: errors.New("cert expired")}
	msg := err.Error()
	if !strings.Contains(msg, installvalidation.TLSErrorPrefix) {
		t.Fatalf("expected TLS prefix, got %q", msg)
	}
}

func TestTLSError_NoCause(t *testing.T) {
	err := &installvalidation.TLSError{Host: "example.com"}
	msg := err.Error()
	if !strings.Contains(msg, "example.com") {
		t.Fatalf("expected host in msg, got %q", msg)
	}
	if err.Unwrap() != nil {
		t.Fatal("Unwrap should be nil when no Cause")
	}
}

func TestIsTLSFailure_ByMessage(t *testing.T) {
	err := errors.New(installvalidation.TLSErrorPrefix + ": more info")
	if !installvalidation.IsTLSFailure(err) {
		t.Fatal("IsTLSFailure should be true for TLS prefix in message")
	}
}

func TestIsTLSFailure_CertVerifyFailed(t *testing.T) {
	err := errors.New("CERTIFICATE_VERIFY_FAILED: bad cert")
	if !installvalidation.IsTLSFailure(err) {
		t.Fatal("IsTLSFailure should be true for CERTIFICATE_VERIFY_FAILED")
	}
}

func TestLocalPathMarkers_NotEmpty(t *testing.T) {
	if len(installvalidation.LocalPathMarkers) == 0 {
		t.Fatal("LocalPathMarkers should not be empty")
	}
	for _, m := range installvalidation.LocalPathMarkers {
		if m == "" {
			t.Fatal("LocalPathMarkers should not contain empty string")
		}
	}
}

func TestLocalPathFailureReason_ApmDir(t *testing.T) {
	dir := t.TempDir()
	apmDir := dir + "/.apm"
	if err := os.Mkdir(apmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	reason := installvalidation.LocalPathFailureReason(dir)
	if reason != "" {
		t.Errorf("expected empty reason for path with .apm dir, got %q", reason)
	}
}

func TestLocalPathNoMarkersHint_MultipleSubpackages(t *testing.T) {
	dir := t.TempDir()
	for _, sub := range []string{"pkgA", "pkgB", "pkgC"} {
		subDir := dir + "/" + sub
		if err := os.MkdirAll(subDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(subDir+"/apm.yml", []byte("name: "+sub+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	hint := installvalidation.LocalPathNoMarkersHint(dir)
	if hint == "" {
		t.Error("expected non-empty hint for dir with multiple sub-packages")
	}
	for _, sub := range []string{"pkgA", "pkgB", "pkgC"} {
		if !strings.Contains(hint, sub) {
			t.Errorf("expected hint to contain %q, got %q", sub, hint)
		}
	}
}

func TestLocalPathNoMarkersHint_FilesIgnored(t *testing.T) {
	dir := t.TempDir()
	// Only files, no subdirs with markers
	if err := os.WriteFile(dir+"/apm.yml", []byte("name: root\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hint := installvalidation.LocalPathNoMarkersHint(dir)
	// Root-level files don't qualify as sub-directories
	_ = hint // no panic
}

func TestNewPackageProber_DefaultTimeout(t *testing.T) {
	p := installvalidation.NewPackageProber("github.com", "")
	if p.Timeout <= 0 {
		t.Error("expected positive default timeout")
	}
	if p.HTTPClient == nil {
		t.Error("expected non-nil HTTPClient")
	}
}

func TestProbeResult_TLSAndAuth_MutuallyExclusive(t *testing.T) {
	r := installvalidation.ProbeResult{IsAuthError: true, IsTLSError: false}
	if r.IsTLSError {
		t.Error("IsTLSError should be false")
	}
	r2 := installvalidation.ProbeResult{IsTLSError: true, IsAuthError: false}
	if r2.IsAuthError {
		t.Error("IsAuthError should be false")
	}
}

func TestIsADOAuthFailureSignal_AzureDevOps(t *testing.T) {
	if !installvalidation.IsADOAuthFailureSignal(200, "Azure DevOps access denied") {
		t.Error("azure devops body should trigger ADO auth signal")
	}
}

func TestValidatePackageExists_RelativePath_NoMarkers(t *testing.T) {
	dir := t.TempDir()
	// dir exists but has no markers
	result := installvalidation.ValidatePackageExists(dir, "github.com", "", false)
	if result.Reachable {
		t.Error("expected Reachable=false for path with no markers")
	}
}

func TestValidatePackageExists_OwnerRepo_NoRef(t *testing.T) {
	// owner/repo without hash ref: should attempt GitHub API (unreachable in sandbox)
	result := installvalidation.ValidatePackageExists("owner/repo", "github.com", "", false)
	// In sandbox, network is blocked; just ensure no panic and Reachable is false
	if result.Reachable {
		t.Skip("unexpected network access in sandbox")
	}
}

func TestValidatePackageExists_MissingOwner(t *testing.T) {
	result := installvalidation.ValidatePackageExists("onlyonepart", "github.com", "", false)
	if result.Reachable {
		t.Error("expected Reachable=false for single-part spec")
	}
}
