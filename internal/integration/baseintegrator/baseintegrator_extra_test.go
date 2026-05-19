package baseintegrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/integration/baseintegrator"
)

func TestCheckCollision_FileNotExist(t *testing.T) {
	managed := map[string]struct{}{"other.md": {}}
	// Non-existent file should never collide
	if baseintegrator.CheckCollision("/nonexistent/path/file.md", "file.md", managed, false, nil) {
		t.Fatal("non-existent file should not collide")
	}
}

func TestCheckCollision_NilManagedNeverCollides(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.md")
	os.WriteFile(f, []byte("x"), 0644)
	if baseintegrator.CheckCollision(f, "file.md", nil, false, nil) {
		t.Fatal("nil managed should never collide")
	}
}

func TestCheckCollision_ManagedFileNoCollision(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.md")
	os.WriteFile(f, []byte("x"), 0644)
	managed := map[string]struct{}{"file.md": {}}
	if baseintegrator.CheckCollision(f, "file.md", managed, false, nil) {
		t.Fatal("file in managed set should not collide")
	}
}

func TestCheckCollision_UserAuthoredCollides(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "user.md")
	os.WriteFile(f, []byte("x"), 0644)
	managed := map[string]struct{}{"other.md": {}}
	if !baseintegrator.CheckCollision(f, "user.md", managed, false, nil) {
		t.Fatal("user-authored file not in managed should collide")
	}
}

func TestCheckCollision_ForceOverrides(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "user.md")
	os.WriteFile(f, []byte("x"), 0644)
	managed := map[string]struct{}{"other.md": {}}
	if baseintegrator.CheckCollision(f, "user.md", managed, true, nil) {
		t.Fatal("force=true should suppress collision")
	}
}

func TestNormalizeManagedFiles_BackslashToForward(t *testing.T) {
	in := map[string]struct{}{`a\b\c.md`: {}}
	out := baseintegrator.NormalizeManagedFiles(in)
	if _, ok := out["a/b/c.md"]; !ok {
		t.Fatal("backslash should be normalized to forward slash")
	}
}

func TestNormalizeManagedFiles_AlreadyForward(t *testing.T) {
	in := map[string]struct{}{"a/b/c.md": {}}
	out := baseintegrator.NormalizeManagedFiles(in)
	if _, ok := out["a/b/c.md"]; !ok {
		t.Fatal("forward slash path should be preserved")
	}
}

func TestNormalizeManagedFiles_Empty(t *testing.T) {
	out := baseintegrator.NormalizeManagedFiles(map[string]struct{}{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

func TestPartitionBucketKey_KnownAlias(t *testing.T) {
	cases := []struct{ prim, target, want string }{
		{"prompts", "copilot", "prompts"},
		{"agents", "copilot", "agents_github"},
		{"commands", "claude", "commands"},
		{"instructions", "copilot", "instructions"},
		{"instructions", "cursor", "rules_cursor"},
		{"instructions", "claude", "rules_claude"},
	}
	for _, c := range cases {
		got := baseintegrator.PartitionBucketKey(c.prim, c.target)
		if got != c.want {
			t.Errorf("PartitionBucketKey(%q, %q) = %q, want %q", c.prim, c.target, got, c.want)
		}
	}
}

func TestPartitionBucketKey_Unknown(t *testing.T) {
	got := baseintegrator.PartitionBucketKey("unknown", "target")
	if got != "unknown_target" {
		t.Errorf("expected unknown_target, got %q", got)
	}
}

func TestValidateDeployPath_DotDotRejected(t *testing.T) {
	if baseintegrator.ValidateDeployPath("../etc/passwd", "/project", nil, nil) {
		t.Error("path with .. should be rejected")
	}
}

func TestValidateDeployPath_NoAllowedPrefixMatch(t *testing.T) {
	if baseintegrator.ValidateDeployPath("hidden/secret", "/project", []string{".github/"}, nil) {
		t.Error("path not matching allowed prefixes should be rejected")
	}
}

func TestBucketAliases_NotEmpty(t *testing.T) {
	if len(baseintegrator.BucketAliases) == 0 {
		t.Error("BucketAliases should not be empty")
	}
}
