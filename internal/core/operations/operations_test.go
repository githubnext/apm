package operations_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/operations"
)

func TestConfigureClient_MissingClientType(t *testing.T) {
	res := operations.ConfigureClient(operations.ConfigureClientOptions{})
	if res.Success {
		t.Fatal("expected failure when client_type is empty")
	}
	if res.Error == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestConfigureClient_WithClientType(t *testing.T) {
	res := operations.ConfigureClient(operations.ConfigureClientOptions{
		ClientType:  "claude",
		ProjectRoot: "/tmp/proj",
	})
	if !res.Success {
		t.Fatalf("expected success, got error: %s", res.Error)
	}
}

func TestInstallPackage_MissingFields(t *testing.T) {
	res := operations.InstallPackage(operations.InstallPackageOptions{})
	if res.Success {
		t.Fatal("expected failure when required fields are empty")
	}
	if !res.Failed {
		t.Fatal("expected Failed=true")
	}
}

func TestInstallPackage_MissingPackageName(t *testing.T) {
	res := operations.InstallPackage(operations.InstallPackageOptions{ClientType: "vscode"})
	if res.Success {
		t.Fatal("expected failure with empty package_name")
	}
}

func TestInstallPackage_ValidFields(t *testing.T) {
	res := operations.InstallPackage(operations.InstallPackageOptions{
		ClientType:  "claude",
		PackageName: "my-tool",
	})
	if !res.Success {
		t.Fatalf("expected success, got: %s", res.Error)
	}
	if !res.Installed {
		t.Fatal("expected Installed=true")
	}
	if res.Skipped {
		t.Fatal("expected Skipped=false")
	}
	if res.Failed {
		t.Fatal("expected Failed=false")
	}
}

func TestUninstallPackage_MissingFields(t *testing.T) {
	res := operations.UninstallPackage(operations.UninstallPackageOptions{})
	if res.Success {
		t.Fatal("expected failure when required fields are empty")
	}
}

func TestUninstallPackage_ValidFields(t *testing.T) {
	res := operations.UninstallPackage(operations.UninstallPackageOptions{
		ClientType:  "claude",
		PackageName: "my-tool",
		ProjectRoot: "/tmp/proj",
	})
	if !res.Success {
		t.Fatalf("expected success, got: %s", res.Error)
	}
}
