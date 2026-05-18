package installservice

import (
	"errors"
	"fmt"
	"testing"
)

func TestInstallRequest_Fields(t *testing.T) {
	req := &InstallRequest{
		Packages:   []string{"owner/repo", "other/pkg"},
		Frozen:     true,
		UpdateRefs: false,
		Scope:      "user",
		Target:     "claude",
		Verbose:    true,
		DryRun:     false,
	}
	if len(req.Packages) != 2 {
		t.Errorf("expected 2 packages, got %d", len(req.Packages))
	}
	if !req.Frozen {
		t.Error("Frozen should be true")
	}
	if req.Scope != "user" {
		t.Errorf("Scope = %q, want user", req.Scope)
	}
}

func TestInstallResult_Fields(t *testing.T) {
	res := &InstallResult{
		Installed: []string{"a/b"},
		Updated:   []string{"c/d"},
		Skipped:   []string{"e/f"},
		Failed:    []string{"g/h"},
		ExitCode:  1,
	}
	if len(res.Installed) != 1 {
		t.Errorf("expected 1 installed, got %d", len(res.Installed))
	}
	if res.ExitCode != 1 {
		t.Errorf("ExitCode = %d, want 1", res.ExitCode)
	}
}

func TestInstallNotAvailableError_Wraps(t *testing.T) {
	inner := errors.New("inner cause")
	outer := &InstallNotAvailableError{Cause: inner}
	if !errors.Is(outer.Cause, inner) {
		t.Error("expected Cause to be the inner error")
	}
}

func TestFrozenInstallError_TypeAssertion(t *testing.T) {
	err := error(&FrozenInstallError{Reason: "missing"})
	var fe *FrozenInstallError
	if !errors.As(err, &fe) {
		t.Fatal("expected FrozenInstallError via errors.As")
	}
	if fe.Reason != "missing" {
		t.Errorf("Reason = %q, want missing", fe.Reason)
	}
}

func TestIsFrozenInstallError_Wrapped(t *testing.T) {
	inner := &FrozenInstallError{Reason: "wrapped"}
	wrapped := fmt.Errorf("context: %w", inner)
	if !IsFrozenInstallError(wrapped) {
		t.Error("IsFrozenInstallError should detect wrapped FrozenInstallError")
	}
}

func TestInstallService_RunEmptyPackages(t *testing.T) {
	svc := New()
	req := &InstallRequest{Packages: []string{}}
	res, err := svc.Run(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", res.ExitCode)
	}
}

func TestInstallService_RunFrozen(t *testing.T) {
	svc := New()
	req := &InstallRequest{Frozen: true, Packages: []string{"x/y"}}
	res, err := svc.Run(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestInstallService_RunDryRun(t *testing.T) {
	svc := New()
	req := &InstallRequest{DryRun: true}
	res, err := svc.Run(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result for dry-run")
	}
}
