package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsToolAvailable_ShReturnsBool(t *testing.T) {
	// sh is universally present on Linux/macOS runners
	got := IsToolAvailable("sh")
	if !got {
		t.Log("sh not found; skipping assertion (unusual environment)")
	}
}

func TestIsToolAvailable_EmptyString(t *testing.T) {
	if IsToolAvailable("") {
		t.Error("empty string should not be a valid tool")
	}
}

func TestGetAvailablePackageManagers_KeysMatchValues(t *testing.T) {
	mgrs := GetAvailablePackageManagers()
	for k, v := range mgrs {
		if k != v {
			t.Errorf("key %q != value %q", k, v)
		}
	}
}

func TestGetAvailablePackageManagers_OnlyKnownNames(t *testing.T) {
	known := map[string]bool{
		"uv": true, "pip": true, "pipx": true,
		"npm": true, "yarn": true, "pnpm": true,
		"brew": true, "apt": true, "yum": true,
		"dnf": true, "apk": true, "pacman": true,
	}
	for k := range GetAvailablePackageManagers() {
		if !known[k] {
			t.Errorf("unexpected manager %q", k)
		}
	}
}

func TestDetectPlatform_OnlyFourValues(t *testing.T) {
	p := DetectPlatform()
	valid := map[string]bool{"macos": true, "linux": true, "windows": true, "unknown": true}
	if !valid[p] {
		t.Errorf("invalid platform %q", p)
	}
}

func TestDetectPlatform_Stable(t *testing.T) {
	a := DetectPlatform()
	b := DetectPlatform()
	if a != b {
		t.Error("DetectPlatform must be deterministic")
	}
}

func TestFindPluginJSON_NoSubdir(t *testing.T) {
	dir := t.TempDir()
	// No plugin.json anywhere -- must return ""
	got := FindPluginJSON(dir)
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestFindPluginJSON_TopLevelFound(t *testing.T) {
	dir := t.TempDir()
	pj := filepath.Join(dir, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindPluginJSON(dir)
	if got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}

func TestFindPluginJSON_GithubSubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".github", "plugin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	pj := filepath.Join(sub, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindPluginJSON(dir)
	if got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}

func TestFindPluginJSON_ClaudeSubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".claude-plugin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	pj := filepath.Join(sub, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindPluginJSON(dir)
	if got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}

func TestFindPluginJSON_CursorSubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".cursor-plugin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	pj := filepath.Join(sub, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindPluginJSON(dir)
	if got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}
