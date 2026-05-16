package coworkpaths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsCoworkPath(t *testing.T) {
	if !IsCoworkPath("cowork://skills/myskill") {
		t.Error("expected true")
	}
	if IsCoworkPath("/absolute/path") {
		t.Error("expected false")
	}
	if IsCoworkPath("") {
		t.Error("expected false for empty string")
	}
}

func TestToLockfilePath_basic(t *testing.T) {
	root := t.TempDir()
	skillPath := filepath.Join(root, "myplugin")
	lp, err := ToLockfilePath(skillPath, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lp != "cowork://skills/myplugin" {
		t.Errorf("unexpected lockfile path: %q", lp)
	}
}

func TestToLockfilePath_escapeRoot(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "..", "other")
	_, err := ToLockfilePath(outside, root)
	if err == nil {
		t.Fatal("expected error for path escaping root")
	}
}

func TestFromLockfilePath_basic(t *testing.T) {
	root := t.TempDir()
	lp := "cowork://skills/myplugin"
	got, err := FromLockfilePath(lp, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(root, "myplugin")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFromLockfilePath_notCowork(t *testing.T) {
	_, err := FromLockfilePath("/some/path", "/root")
	if err == nil {
		t.Fatal("expected error for non-cowork path")
	}
}

func TestFromLockfilePath_traversal(t *testing.T) {
	root := t.TempDir()
	_, err := FromLockfilePath("cowork://skills/../../../etc/passwd", root)
	if err == nil {
		t.Fatal("expected error for traversal path")
	}
}

func TestRoundTrip(t *testing.T) {
	root := t.TempDir()
	subdir := filepath.Join(root, "plugin", "v2")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	lp, err := ToLockfilePath(subdir, root)
	if err != nil {
		t.Fatalf("ToLockfilePath: %v", err)
	}
	got, err := FromLockfilePath(lp, root)
	if err != nil {
		t.Fatalf("FromLockfilePath: %v", err)
	}
	if got != subdir {
		t.Errorf("round-trip mismatch: want %q, got %q", subdir, got)
	}
}

func TestResolveCoworkSkillsDir_envOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APM_COPILOT_COWORK_SKILLS_DIR", dir)
	got, err := ResolveCoworkSkillsDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty path")
	}
}

func TestResolveCoworkSkillsDir_traversalEnv(t *testing.T) {
	t.Setenv("APM_COPILOT_COWORK_SKILLS_DIR", "/safe/path/../../../etc")
	_, err := ResolveCoworkSkillsDir()
	if err == nil {
		t.Fatal("expected error for traversal in env var")
	}
}
