package helpers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/utils/helpers"
)

func TestIsToolAvailable(t *testing.T) {
	// "sh" or "cat" should exist on any POSIX CI runner.
	if !helpers.IsToolAvailable("sh") {
		t.Error("expected 'sh' to be found on PATH")
	}
	if helpers.IsToolAvailable("definitely-not-a-real-binary-xyz") {
		t.Error("expected nonexistent tool to return false")
	}
}

func TestDetectPlatform(t *testing.T) {
	p := helpers.DetectPlatform()
	valid := map[string]bool{"macos": true, "linux": true, "windows": true, "unknown": true}
	if !valid[p] {
		t.Errorf("unexpected platform %q", p)
	}
}

func TestGetAvailablePackageManagers(t *testing.T) {
	// Just check it returns a map (may be empty in a minimal container).
	m := helpers.GetAvailablePackageManagers()
	if m == nil {
		t.Error("expected non-nil map")
	}
}

func TestFindPluginJSON(t *testing.T) {
	dir := t.TempDir()

	// No plugin.json yet.
	if got := helpers.FindPluginJSON(dir); got != "" {
		t.Errorf("expected empty, got %q", got)
	}

	// Create the top-level plugin.json.
	pj := filepath.Join(dir, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := helpers.FindPluginJSON(dir); got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}

func TestFindPluginJSONSubdirs(t *testing.T) {
	dir := t.TempDir()

	// Create under .github/plugin/
	sub := filepath.Join(dir, ".github", "plugin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	pj := filepath.Join(sub, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := helpers.FindPluginJSON(dir); got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}

func TestFindPluginJSON_ClaudePlugin(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".claude-plugin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	pj := filepath.Join(sub, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Top-level not present, .claude-plugin should be found.
	if got := helpers.FindPluginJSON(dir); got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}

func TestFindPluginJSON_CursorPlugin(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".cursor-plugin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	pj := filepath.Join(sub, "plugin.json")
	if err := os.WriteFile(pj, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := helpers.FindPluginJSON(dir); got != pj {
		t.Errorf("expected %q, got %q", pj, got)
	}
}

func TestFindPluginJSON_TopLevelTakesPrecedence(t *testing.T) {
	dir := t.TempDir()
	// Create both top-level and sub-directory plugin.json
	topPJ := filepath.Join(dir, "plugin.json")
	if err := os.WriteFile(topPJ, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, ".claude-plugin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	subPJ := filepath.Join(sub, "plugin.json")
	if err := os.WriteFile(subPJ, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Top-level should win.
	if got := helpers.FindPluginJSON(dir); got != topPJ {
		t.Errorf("expected top-level %q, got %q", topPJ, got)
	}
}

func TestIsToolAvailable_Cat(t *testing.T) {
	// cat should be available on any POSIX system
	if !helpers.IsToolAvailable("cat") {
		t.Skip("cat not available on this platform")
	}
}

func TestDetectPlatform_NotEmpty(t *testing.T) {
	p := helpers.DetectPlatform()
	if p == "" {
		t.Error("DetectPlatform() should never return empty string")
	}
}
