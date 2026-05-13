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
