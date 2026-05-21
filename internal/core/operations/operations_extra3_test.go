package operations

import (
	"testing"
)

func TestConfigureClient_EmptyClientType(t *testing.T) {
	result := ConfigureClient(ConfigureClientOptions{})
	if result.Success {
		t.Error("expected failure for empty ClientType")
	}
	if result.Error == "" {
		t.Error("expected error message for empty ClientType")
	}
}

func TestConfigureClient_WithClientType(t *testing.T) {
	result := ConfigureClient(ConfigureClientOptions{ClientType: "copilot"})
	if !result.Success {
		t.Errorf("expected success, got error: %q", result.Error)
	}
}

func TestInstallPackage_MissingFields(t *testing.T) {
	result := InstallPackage(InstallPackageOptions{})
	if result.Success {
		t.Error("expected failure for missing required fields")
	}
	if !result.Failed {
		t.Error("expected Failed=true")
	}
}

func TestInstallPackage_WithRequiredFields(t *testing.T) {
	result := InstallPackage(InstallPackageOptions{
		ClientType:  "copilot",
		PackageName: "my-pkg",
	})
	if !result.Success {
		t.Errorf("expected success, got error: %q", result.Error)
	}
	if !result.Installed {
		t.Error("expected Installed=true")
	}
}

func TestUninstallPackageResult_ZeroValue(t *testing.T) {
	var r UninstallPackageResult
	if r.Success {
		t.Error("expected Success=false for zero value")
	}
}

func TestConfigureClientOptions_ZeroValue(t *testing.T) {
	var opts ConfigureClientOptions
	if opts.ClientType != "" {
		t.Error("expected empty ClientType for zero value")
	}
}

func TestInstallPackageOptions_ZeroValue(t *testing.T) {
	var opts InstallPackageOptions
	if opts.UserScope {
		t.Error("expected UserScope=false for zero value")
	}
}

func TestUninstallPackageOptions_ZeroValue(t *testing.T) {
	var opts UninstallPackageOptions
	if opts.UserScope {
		t.Error("expected UserScope=false for zero value")
	}
}
