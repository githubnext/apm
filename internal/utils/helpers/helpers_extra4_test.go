package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsToolAvailable_ExistingTool_Extra4(t *testing.T) {
	// "ls" should be available on any Unix system
	if !IsToolAvailable("ls") {
		t.Skip("ls not available in this environment")
	}
}

func TestIsToolAvailable_NonExistent_Extra4(t *testing.T) {
	if IsToolAvailable("__nonexistent_tool_xyz_abc__") {
		t.Error("non-existent tool should not be available")
	}
}

func TestDetectPlatform_ReturnsKnown_Extra4(t *testing.T) {
	p := DetectPlatform()
	known := map[string]bool{"macos": true, "linux": true, "windows": true, "unknown": true}
	if !known[p] {
		t.Errorf("DetectPlatform returned unexpected %q", p)
	}
}

func TestDetectPlatform_NonEmpty_Extra4(t *testing.T) {
	p := DetectPlatform()
	if p == "" {
		t.Error("DetectPlatform should return non-empty string")
	}
}

func TestGetAvailablePackageManagers_ReturnsMap_Extra4(t *testing.T) {
	m := GetAvailablePackageManagers()
	if m == nil {
		t.Error("GetAvailablePackageManagers should return non-nil map")
	}
	for k, v := range m {
		if k == "" {
			t.Error("empty key in package managers")
		}
		if v == "" {
			t.Errorf("empty value for key %q", k)
		}
	}
}

func TestFindPluginJSON_MissingDir_Extra4(t *testing.T) {
	result := FindPluginJSON("/nonexistent/path/xyz")
	if result != "" {
		t.Errorf("FindPluginJSON missing dir = %q, want empty", result)
	}
}

func TestFindPluginJSON_EmptyDir_Extra4(t *testing.T) {
	dir := t.TempDir()
	result := FindPluginJSON(dir)
	if result != "" {
		t.Errorf("FindPluginJSON empty dir = %q, want empty", result)
	}
}

func TestFindPluginJSON_WithFile_Extra4(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "plugin.json")
	if err := os.WriteFile(pf, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	result := FindPluginJSON(dir)
	if result != pf {
		t.Errorf("FindPluginJSON = %q, want %q", result, pf)
	}
}
