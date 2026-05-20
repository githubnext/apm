package operations

import (
	"testing"
)

func TestConfigureClientResult_ZeroValue_Extra2(t *testing.T) {
	var r ConfigureClientResult
	if r.Success || r.Error != "" {
		t.Error("zero-value ConfigureClientResult should have false/empty fields")
	}
}

func TestConfigureClientResult_Fields_Extra2(t *testing.T) {
	r := ConfigureClientResult{Success: true, Error: ""}
	if !r.Success {
		t.Error("Success should be true")
	}
}

func TestInstallPackageResult_ZeroValue_Extra2(t *testing.T) {
	var r InstallPackageResult
	if r.Success || r.Installed || r.Skipped || r.Failed || r.Error != "" {
		t.Error("zero-value InstallPackageResult should have false/empty fields")
	}
}

func TestInstallPackageResult_Fields_Extra2(t *testing.T) {
	r := InstallPackageResult{
		Success:   true,
		Installed: true,
		Skipped:   false,
		Failed:    false,
		Error:     "",
	}
	if !r.Success || !r.Installed {
		t.Error("Success and Installed should be true")
	}
}

func TestUninstallPackageResult_ZeroValue_Extra2(t *testing.T) {
	var r UninstallPackageResult
	if r.Success || r.Error != "" {
		t.Error("zero-value UninstallPackageResult should have false/empty fields")
	}
}

func TestConfigureClientOptions_ZeroValue_Extra2(t *testing.T) {
	var opts ConfigureClientOptions
	if opts.ClientType != "" || opts.ProjectRoot != "" || opts.UserScope {
		t.Error("zero-value ConfigureClientOptions should have empty/false fields")
	}
}

func TestConfigureClientOptions_Fields_Extra2(t *testing.T) {
	opts := ConfigureClientOptions{
		ClientType:    "copilot",
		ConfigUpdates: map[string]interface{}{"key": "val"},
		ProjectRoot:   "/proj",
		UserScope:     true,
	}
	if opts.ClientType != "copilot" {
		t.Errorf("ClientType = %q", opts.ClientType)
	}
	if opts.ConfigUpdates["key"] != "val" {
		t.Errorf("ConfigUpdates[key] = %v", opts.ConfigUpdates["key"])
	}
}

func TestInstallPackageOptions_ZeroValue_Extra2(t *testing.T) {
	var opts InstallPackageOptions
	if opts.ClientType != "" || opts.PackageName != "" || opts.UserScope {
		t.Error("zero-value InstallPackageOptions should have empty/false fields")
	}
}

func TestInstallPackageOptions_Fields_Extra2(t *testing.T) {
	opts := InstallPackageOptions{
		ClientType:  "claude",
		PackageName: "mypkg",
		Version:     "1.0.0",
		ProjectRoot: "/proj",
		UserScope:   true,
	}
	if opts.ClientType != "claude" {
		t.Errorf("ClientType = %q", opts.ClientType)
	}
	if opts.PackageName != "mypkg" {
		t.Errorf("PackageName = %q", opts.PackageName)
	}
}

func TestUninstallPackageOptions_Fields_Extra2(t *testing.T) {
	opts := UninstallPackageOptions{
		ClientType:  "vscode",
		PackageName: "pkg",
		ProjectRoot: "/proj",
		UserScope:   false,
	}
	if opts.ClientType != "vscode" {
		t.Errorf("ClientType = %q", opts.ClientType)
	}
}

func TestConfigureClient_EmptyClientType_Extra2(t *testing.T) {
	result := ConfigureClient(ConfigureClientOptions{})
	if result.Success {
		t.Error("expected failure for empty client type")
	}
	if result.Error == "" {
		t.Error("expected non-empty error for empty client type")
	}
}

func TestInstallPackage_EmptyClientType_Extra2(t *testing.T) {
	result := InstallPackage(InstallPackageOptions{PackageName: "pkg"})
	if result.Success {
		t.Error("expected failure for empty client type")
	}
}

func TestUninstallPackage_EmptyPackageName_Extra2(t *testing.T) {
	result := UninstallPackage(UninstallPackageOptions{ClientType: "copilot"})
	if result.Success {
		t.Error("expected failure for empty package name")
	}
}

func TestInstallPackage_EmptyPackageName_Extra2(t *testing.T) {
	result := InstallPackage(InstallPackageOptions{ClientType: "copilot"})
	if result.Success {
		t.Error("expected failure for empty package name")
	}
}

func TestInstallPackageResult_ErrorField_Extra2(t *testing.T) {
	r := InstallPackageResult{Success: false, Error: "something went wrong", Failed: true}
	if r.Error != "something went wrong" {
		t.Errorf("Error = %q", r.Error)
	}
	if !r.Failed {
		t.Error("Failed should be true")
	}
}
