package installservice

import (
"errors"
"fmt"
"strings"
"testing"
)

func TestInstallNotAvailableError_ErrorString(t *testing.T) {
err := &InstallNotAvailableError{Cause: errors.New("db timeout")}
got := err.Error()
if !strings.Contains(got, "db timeout") {
t.Errorf("expected cause in error string, got %q", got)
}
if !strings.Contains(got, "unavailable") {
t.Errorf("expected 'unavailable' in error string, got %q", got)
}
}

func TestInstallNotAvailableError_NilCause(t *testing.T) {
err := &InstallNotAvailableError{}
got := err.Error()
if got == "" {
t.Error("expected non-empty error string even with nil cause")
}
}

func TestFrozenInstallError_Message(t *testing.T) {
err := &FrozenInstallError{Reason: "outdated lock"}
msg := err.Error()
if !strings.Contains(msg, "outdated lock") {
t.Errorf("expected reason in error: %q", msg)
}
}

func TestFrozenInstallError_EmptyReason(t *testing.T) {
err := &FrozenInstallError{}
msg := err.Error()
_ = msg // should not panic
}

func TestIsFrozenInstallError_Direct(t *testing.T) {
err := &FrozenInstallError{Reason: "direct"}
if !IsFrozenInstallError(err) {
t.Error("direct FrozenInstallError should be detected")
}
}

func TestIsFrozenInstallError_NotFrozen(t *testing.T) {
err := errors.New("ordinary error")
if IsFrozenInstallError(err) {
t.Error("ordinary error should not be a FrozenInstallError")
}
}

func TestIsFrozenInstallError_DoubleWrapped(t *testing.T) {
inner := &FrozenInstallError{Reason: "lock"}
mid := fmt.Errorf("mid: %w", inner)
outer := fmt.Errorf("outer: %w", mid)
if !IsFrozenInstallError(outer) {
t.Error("double-wrapped FrozenInstallError should be detected")
}
}

func TestInstallRequest_AllFields(t *testing.T) {
req := &InstallRequest{
Packages:   []string{"a/b", "c/d"},
Frozen:     true,
UpdateRefs: true,
Scope:      "project",
Target:     "claude",
Verbose:    true,
DryRun:     true,
}
if !req.UpdateRefs {
t.Error("UpdateRefs should be true")
}
if req.Scope != "project" {
t.Errorf("Scope = %q, want project", req.Scope)
}
if !req.DryRun {
t.Error("DryRun should be true")
}
}

func TestInstallResult_AllFields(t *testing.T) {
res := &InstallResult{
Installed: []string{"a"},
Updated:   []string{"b"},
Skipped:   []string{"c"},
Failed:    []string{"d"},
ExitCode:  2,
}
if res.ExitCode != 2 {
t.Errorf("ExitCode = %d, want 2", res.ExitCode)
}
if len(res.Installed) != 1 {
t.Errorf("Installed len = %d, want 1", len(res.Installed))
}
if len(res.Updated) != 1 {
t.Errorf("Updated len = %d, want 1", len(res.Updated))
}
if len(res.Skipped) != 1 {
t.Errorf("Skipped len = %d, want 1", len(res.Skipped))
}
if len(res.Failed) != 1 {
t.Errorf("Failed len = %d, want 1", len(res.Failed))
}
}

func TestInstallService_RunMultiplePackages(t *testing.T) {
svc := New()
req := &InstallRequest{
Packages: []string{"owner/repo1", "owner/repo2", "owner/repo3"},
Verbose:  true,
}
res, err := svc.Run(req)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if res == nil {
t.Fatal("expected non-nil result")
}
}

func TestInstallService_RunWithScope(t *testing.T) {
svc := New()
req := &InstallRequest{
Packages: []string{"owner/pkg"},
Scope:    "user",
}
res, err := svc.Run(req)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if res.ExitCode != 0 {
t.Errorf("ExitCode = %d, want 0", res.ExitCode)
}
}

func TestInstallService_RunUpdateRefs(t *testing.T) {
svc := New()
req := &InstallRequest{
UpdateRefs: true,
Packages:   []string{"owner/pkg"},
}
res, err := svc.Run(req)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if res == nil {
t.Fatal("expected non-nil result")
}
}

func TestInstallService_New_notNil(t *testing.T) {
svc := New()
if svc == nil {
t.Fatal("New() should return non-nil")
}
}
