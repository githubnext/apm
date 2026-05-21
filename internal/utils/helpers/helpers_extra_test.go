package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindPluginJSON_Missing(t *testing.T) {
	dir := t.TempDir()
	result := FindPluginJSON(dir)
	if result != "" {
		t.Errorf("expected empty string for missing plugin.json, got %q", result)
	}
}

func TestFindPluginJSON_GithubPlugin(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, ".github", "plugin")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(subdir, "plugin.json")
	if err := os.WriteFile(p, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindPluginJSON(dir)
	if got != p {
		t.Errorf("expected %q, got %q", p, got)
	}
}

func TestDetectPlatform_ReturnsValidString(t *testing.T) {
	p := DetectPlatform()
	valid := map[string]bool{"macos": true, "linux": true, "windows": true, "unknown": true}
	if !valid[p] {
		t.Errorf("DetectPlatform returned unexpected value %q", p)
	}
}

func TestGetAvailablePackageManagers_ReturnsMap(t *testing.T) {
	mgrs := GetAvailablePackageManagers()
	if mgrs == nil {
		t.Error("GetAvailablePackageManagers returned nil")
	}
	// values should equal keys
	for k, v := range mgrs {
		if k != v {
			t.Errorf("map[%q]=%q: expected key==value", k, v)
		}
	}
}

func TestIsToolAvailable_Nonexistent(t *testing.T) {
	if IsToolAvailable("this_tool_should_never_exist_xyzzy_12345") {
		t.Error("expected false for nonexistent tool")
	}
}

func TestFindPluginJSON_TopLevelFirst(t *testing.T) {
	dir := t.TempDir()
	// Create both top-level and .github/plugin/plugin.json
	topPlugin := filepath.Join(dir, "plugin.json")
	if err := os.WriteFile(topPlugin, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	subdir := filepath.Join(dir, ".github", "plugin")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "plugin.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindPluginJSON(dir)
	if got != topPlugin {
		t.Errorf("expected top-level plugin.json to take precedence, got %q", got)
	}
}

func TestFindPluginJSON_NonExistentDir(t *testing.T) {
	result := FindPluginJSON("/nonexistent/path/that/does/not/exist")
	if result != "" {
		t.Errorf("expected empty for nonexistent directory, got %q", result)
	}
}

func TestGetAvailablePackageManagers_Consistency(t *testing.T) {
	mgrs1 := GetAvailablePackageManagers()
	mgrs2 := GetAvailablePackageManagers()
	if len(mgrs1) != len(mgrs2) {
		t.Error("GetAvailablePackageManagers returned different lengths on consecutive calls")
	}
}
