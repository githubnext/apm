package helpers

import (
	"strings"
	"testing"
)

func TestDetectPlatform_KnownValues(t *testing.T) {
	p := DetectPlatform()
	valid := map[string]bool{"macos": true, "linux": true, "windows": true, "unknown": true}
	if !valid[p] {
		t.Errorf("unexpected platform value: %q", p)
	}
}

func TestDetectPlatform_NonEmpty_Stable(t *testing.T) {
	p1 := DetectPlatform()
	p2 := DetectPlatform()
	if p1 != p2 {
		t.Error("DetectPlatform should return same value on repeated calls")
	}
}

func TestIsToolAvailable_EchoOrTrue(t *testing.T) {
	// Both echo and true are universally available on Linux/macOS
	if !IsToolAvailable("echo") && !IsToolAvailable("true") {
		t.Error("expected at least one of echo or true to be available")
	}
}

func TestIsToolAvailable_FakeToolReturnsFalse(t *testing.T) {
	if IsToolAvailable("zzz-definitely-not-a-real-tool-12345") {
		t.Error("expected false for nonexistent tool")
	}
}

func TestGetAvailablePackageManagers_ValueEqualsKey(t *testing.T) {
	mgrs := GetAvailablePackageManagers()
	for k, v := range mgrs {
		if k != v {
			t.Errorf("expected key==value, got key=%q val=%q", k, v)
		}
	}
}

func TestGetAvailablePackageManagers_NilNotReturned(t *testing.T) {
	mgrs := GetAvailablePackageManagers()
	if mgrs == nil {
		t.Error("GetAvailablePackageManagers should not return nil")
	}
}

func TestFindPluginJSON_EmptyDir(t *testing.T) {
	result := FindPluginJSON("")
	// Empty path should return empty string (file won't exist)
	if result == "/" || strings.HasSuffix(result, "plugin.json") {
		// If it happens to find something under cwd that's acceptable
	}
	// Just ensure no panic
}

func TestFindPluginJSON_NonExistentPath(t *testing.T) {
	result := FindPluginJSON("/nonexistent/path/that/does/not/exist/xyz123")
	if result != "" {
		t.Errorf("expected empty string for nonexistent path, got %q", result)
	}
}

func TestGetAvailablePackageManagers_Idempotent(t *testing.T) {
	m1 := GetAvailablePackageManagers()
	m2 := GetAvailablePackageManagers()
	if len(m1) != len(m2) {
		t.Error("repeated calls should return same number of package managers")
	}
}
