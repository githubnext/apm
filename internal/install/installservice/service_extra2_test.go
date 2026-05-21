package installservice

import (
	"errors"
	"testing"
)

func TestInstallRequest_ZeroValue(t *testing.T) {
	var req InstallRequest
	if req.Frozen {
		t.Error("Frozen should default to false")
	}
	if req.DryRun {
		t.Error("DryRun should default to false")
	}
	if req.Verbose {
		t.Error("Verbose should default to false")
	}
	if len(req.Packages) != 0 {
		t.Error("Packages should default to nil")
	}
}

func TestInstallResult_ZeroValue(t *testing.T) {
	var res InstallResult
	if res.ExitCode != 0 {
		t.Errorf("ExitCode should default to 0, got %d", res.ExitCode)
	}
	if len(res.Installed) != 0 || len(res.Failed) != 0 {
		t.Error("slices should default to nil")
	}
}

func TestInstallNotAvailableError_Interface(t *testing.T) {
	var err error = &InstallNotAvailableError{Cause: errors.New("x")}
	if err == nil {
		t.Fatal("should implement error")
	}
}

func TestFrozenInstallError_Interface(t *testing.T) {
	var err error = &FrozenInstallError{Reason: "r"}
	if err == nil {
		t.Fatal("should implement error")
	}
}

func TestInstallNotAvailableError_CauseNil(t *testing.T) {
	e := &InstallNotAvailableError{}
	if e.Error() != "APM install subsystem unavailable" {
		t.Errorf("unexpected: %s", e.Error())
	}
}

func TestFrozenInstallError_ReasonSet(t *testing.T) {
	e := &FrozenInstallError{Reason: "stale lock"}
	msg := e.Error()
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestIsFrozenInstallError_Nil(t *testing.T) {
	if IsFrozenInstallError(nil) {
		t.Error("nil should not be FrozenInstallError")
	}
}

func TestIsFrozenInstallError_OtherError(t *testing.T) {
	if IsFrozenInstallError(errors.New("other")) {
		t.Error("plain error should not match")
	}
}

func TestInstallResult_Populated(t *testing.T) {
	res := &InstallResult{
		Installed: []string{"a/b", "c/d"},
		Updated:   []string{"e/f"},
		Failed:    []string{},
		ExitCode:  0,
	}
	if len(res.Installed) != 2 {
		t.Errorf("expected 2 installed, got %d", len(res.Installed))
	}
	if res.ExitCode != 0 {
		t.Error("expected success exit code")
	}
}

func TestInstallService_NewIsNotNil(t *testing.T) {
	svc := New()
	if svc == nil {
		t.Fatal("New() should not return nil")
	}
}
