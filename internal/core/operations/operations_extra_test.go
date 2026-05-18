package operations_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/operations"
)

func TestConfigureClient_WithConfigUpdates(t *testing.T) {
	res := operations.ConfigureClient(operations.ConfigureClientOptions{
		ClientType:    "vscode",
		ConfigUpdates: map[string]interface{}{"theme": "dark"},
	})
	if !res.Success {
		t.Fatalf("expected success, got error: %s", res.Error)
	}
	if res.Error != "" {
		t.Errorf("expected empty error, got %q", res.Error)
	}
}

func TestConfigureClient_UserScope(t *testing.T) {
	res := operations.ConfigureClient(operations.ConfigureClientOptions{
		ClientType: "claude",
		UserScope:  true,
	})
	if !res.Success {
		t.Fatalf("expected success with user scope, got: %s", res.Error)
	}
}

func TestConfigureClient_ErrorMessage(t *testing.T) {
	res := operations.ConfigureClient(operations.ConfigureClientOptions{})
	if res.Error == "" {
		t.Error("expected non-empty error message when client_type is missing")
	}
}

func TestInstallPackage_WithVersion(t *testing.T) {
	res := operations.InstallPackage(operations.InstallPackageOptions{
		ClientType:  "claude",
		PackageName: "my-tool",
		Version:     "1.2.3",
	})
	if !res.Success {
		t.Fatalf("expected success: %s", res.Error)
	}
	if res.Skipped {
		t.Error("expected Skipped=false")
	}
	if res.Failed {
		t.Error("expected Failed=false")
	}
}

func TestInstallPackage_OnlyClientType(t *testing.T) {
	res := operations.InstallPackage(operations.InstallPackageOptions{
		ClientType: "gemini",
	})
	if res.Success {
		t.Error("expected failure without package_name")
	}
	if !res.Failed {
		t.Error("expected Failed=true")
	}
	if res.Error == "" {
		t.Error("expected non-empty error")
	}
}

func TestInstallPackage_WithSharedEnvVars(t *testing.T) {
	res := operations.InstallPackage(operations.InstallPackageOptions{
		ClientType:    "claude",
		PackageName:   "tool",
		SharedEnvVars: map[string]string{"TOKEN": "abc123"},
	})
	if !res.Success {
		t.Fatalf("unexpected failure: %s", res.Error)
	}
	if !res.Installed {
		t.Error("expected Installed=true")
	}
}

func TestUninstallPackage_OnlyClientType(t *testing.T) {
	res := operations.UninstallPackage(operations.UninstallPackageOptions{
		ClientType: "vscode",
	})
	if res.Success {
		t.Error("expected failure without package_name")
	}
	if res.Error == "" {
		t.Error("expected non-empty error")
	}
}

func TestUninstallPackage_OnlyPackageName(t *testing.T) {
	res := operations.UninstallPackage(operations.UninstallPackageOptions{
		PackageName: "tool",
	})
	if res.Success {
		t.Error("expected failure without client_type")
	}
}
