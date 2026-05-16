package installservice

import (
	"errors"
	"testing"
)

func TestInstallNotAvailableError_WithCause(t *testing.T) {
	err := &InstallNotAvailableError{Cause: errors.New("db down")}
	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
	if msg != "APM install subsystem unavailable: db down" {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestInstallNotAvailableError_NoCause(t *testing.T) {
	err := &InstallNotAvailableError{}
	if err.Error() != "APM install subsystem unavailable" {
		t.Errorf("unexpected: %s", err.Error())
	}
}

func TestFrozenInstallError_WithReason(t *testing.T) {
	err := &FrozenInstallError{Reason: "lockfile missing"}
	if err.Error() != "frozen install failed: lockfile missing" {
		t.Errorf("unexpected: %s", err.Error())
	}
}

func TestFrozenInstallError_NoReason(t *testing.T) {
	err := &FrozenInstallError{}
	if err.Error() != "frozen install failed: lockfile missing or out of sync" {
		t.Errorf("unexpected: %s", err.Error())
	}
}

func TestIsFrozenInstallError(t *testing.T) {
	frozen := &FrozenInstallError{Reason: "test"}
	if !IsFrozenInstallError(frozen) {
		t.Error("expected IsFrozenInstallError to return true")
	}
	other := errors.New("other error")
	if IsFrozenInstallError(other) {
		t.Error("expected IsFrozenInstallError to return false for non-FrozenInstallError")
	}
}

func TestInstallService_RunNilRequest(t *testing.T) {
	svc := New()
	_, err := svc.Run(nil)
	if err == nil {
		t.Error("expected error for nil request")
	}
	var notAvail *InstallNotAvailableError
	if !errors.As(err, &notAvail) {
		t.Errorf("expected InstallNotAvailableError, got %T: %v", err, err)
	}
}

func TestInstallService_RunValidRequest(t *testing.T) {
	svc := New()
	req := &InstallRequest{Packages: []string{"foo/bar"}}
	result, err := svc.Run(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.ExitCode != 0 {
		t.Errorf("expected ExitCode 0, got %d", result.ExitCode)
	}
}

func TestInstallService_Interface(t *testing.T) {
	var _ InstallServicer = New()
}
