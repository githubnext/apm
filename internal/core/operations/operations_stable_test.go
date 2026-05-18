package operations_test

import (
"testing"

"github.com/githubnext/apm/internal/core/operations"
)

func TestConfigureClient_AllClients(t *testing.T) {
clients := []string{"claude", "vscode", "gemini", "cursor", "windsurf", "copilot"}
for _, c := range clients {
res := operations.ConfigureClient(operations.ConfigureClientOptions{ClientType: c})
if !res.Success {
t.Errorf("ConfigureClient(%q) failed: %s", c, res.Error)
}
}
}

func TestConfigureClient_ProjectRoot_set(t *testing.T) {
res := operations.ConfigureClient(operations.ConfigureClientOptions{
ClientType:  "claude",
ProjectRoot: "/tmp/my-project",
})
if !res.Success {
t.Fatalf("expected success, got: %s", res.Error)
}
}

func TestConfigureClient_EmptyConfigUpdates(t *testing.T) {
res := operations.ConfigureClient(operations.ConfigureClientOptions{
ClientType:    "vscode",
ConfigUpdates: map[string]interface{}{},
})
if !res.Success {
t.Fatalf("expected success with empty config updates: %s", res.Error)
}
}

func TestConfigureClient_NilConfigUpdates(t *testing.T) {
res := operations.ConfigureClient(operations.ConfigureClientOptions{
ClientType:    "claude",
ConfigUpdates: nil,
})
if !res.Success {
t.Fatalf("expected success with nil config updates: %s", res.Error)
}
}

func TestConfigureClient_UserScope_and_ProjectRoot(t *testing.T) {
res := operations.ConfigureClient(operations.ConfigureClientOptions{
ClientType:  "cursor",
UserScope:   true,
ProjectRoot: "/tmp/proj",
})
if !res.Success {
t.Fatalf("expected success: %s", res.Error)
}
}

func TestInstallPackage_ProjectRoot(t *testing.T) {
res := operations.InstallPackage(operations.InstallPackageOptions{
ClientType:  "claude",
PackageName: "tool",
ProjectRoot: "/tmp/proj",
})
if !res.Success {
t.Fatalf("expected success: %s", res.Error)
}
}

func TestInstallPackage_UserScope(t *testing.T) {
res := operations.InstallPackage(operations.InstallPackageOptions{
ClientType:  "gemini",
PackageName: "tool",
UserScope:   true,
})
if !res.Success {
t.Fatalf("expected success: %s", res.Error)
}
}

func TestInstallPackage_AllClients(t *testing.T) {
clients := []string{"claude", "vscode", "gemini", "cursor", "windsurf"}
for _, c := range clients {
res := operations.InstallPackage(operations.InstallPackageOptions{
ClientType:  c,
PackageName: "mypkg",
})
if !res.Success {
t.Errorf("InstallPackage(%q) failed: %s", c, res.Error)
}
if !res.Installed {
t.Errorf("InstallPackage(%q) Installed should be true", c)
}
}
}

func TestInstallPackage_EmptyVersion(t *testing.T) {
res := operations.InstallPackage(operations.InstallPackageOptions{
ClientType:  "claude",
PackageName: "tool",
Version:     "",
})
if !res.Success {
t.Fatalf("expected success with empty version: %s", res.Error)
}
}

func TestInstallPackage_NilSharedEnvVars(t *testing.T) {
res := operations.InstallPackage(operations.InstallPackageOptions{
ClientType:    "claude",
PackageName:   "tool",
SharedEnvVars: nil,
})
if !res.Success {
t.Fatalf("expected success with nil SharedEnvVars: %s", res.Error)
}
}

func TestInstallPackage_EmptySharedEnvVars(t *testing.T) {
res := operations.InstallPackage(operations.InstallPackageOptions{
ClientType:    "claude",
PackageName:   "tool",
SharedEnvVars: map[string]string{},
})
if !res.Success {
t.Fatalf("expected success with empty SharedEnvVars: %s", res.Error)
}
}

func TestUninstallPackage_AllClients(t *testing.T) {
clients := []string{"claude", "vscode", "gemini", "cursor", "windsurf"}
for _, c := range clients {
res := operations.UninstallPackage(operations.UninstallPackageOptions{
ClientType:  c,
PackageName: "mypkg",
})
if !res.Success {
t.Errorf("UninstallPackage(%q) failed: %s", c, res.Error)
}
}
}

func TestUninstallPackage_ProjectRoot(t *testing.T) {
res := operations.UninstallPackage(operations.UninstallPackageOptions{
ClientType:  "claude",
PackageName: "tool",
ProjectRoot: "/tmp/proj",
})
if !res.Success {
t.Fatalf("expected success: %s", res.Error)
}
}

func TestUninstallPackage_UserScope(t *testing.T) {
res := operations.UninstallPackage(operations.UninstallPackageOptions{
ClientType:  "claude",
PackageName: "tool",
UserScope:   true,
})
if !res.Success {
t.Fatalf("expected success: %s", res.Error)
}
}

func TestUninstallPackage_ErrorOnEmpty(t *testing.T) {
res := operations.UninstallPackage(operations.UninstallPackageOptions{})
if res.Success {
t.Error("expected failure for empty options")
}
if res.Error == "" {
t.Error("expected non-empty error message")
}
}

func TestInstallPackageResult_Fields(t *testing.T) {
res := operations.InstallPackage(operations.InstallPackageOptions{
ClientType:  "claude",
PackageName: "pkg",
})
if res.Skipped {
t.Error("expected Skipped=false for fresh install")
}
if res.Failed {
t.Error("expected Failed=false for fresh install")
}
}
