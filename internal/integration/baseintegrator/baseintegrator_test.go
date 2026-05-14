package baseintegrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/integration/baseintegrator"
)

func TestCheckCollisionNilManaged(t *testing.T) {
	if baseintegrator.CheckCollision("/any/path", "any/path", nil, false, nil) {
		t.Fatal("nil managed should never collide")
	}
}

func TestCheckCollisionManagedContains(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.md")
	os.WriteFile(f, []byte("x"), 0644)
	managed := map[string]struct{}{"file.md": {}}
	if baseintegrator.CheckCollision(f, "file.md", managed, false, nil) {
		t.Fatal("file in managed set should not collide")
	}
}

func TestCheckCollisionUserAuthored(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.md")
	os.WriteFile(f, []byte("x"), 0644)
	managed := map[string]struct{}{"other.md": {}}
	if !baseintegrator.CheckCollision(f, "file.md", managed, false, nil) {
		t.Fatal("user-authored file should collide")
	}
}

func TestCheckCollisionForceOverrides(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.md")
	os.WriteFile(f, []byte("x"), 0644)
	managed := map[string]struct{}{"other.md": {}}
	if baseintegrator.CheckCollision(f, "file.md", managed, true, nil) {
		t.Fatal("force should override collision")
	}
}

func TestNormalizeManagedFilesBackslash(t *testing.T) {
	in := map[string]struct{}{`a\b\c.md`: {}}
	out := baseintegrator.NormalizeManagedFiles(in)
	if _, ok := out["a/b/c.md"]; !ok {
		t.Fatal("backslash should be normalized to forward slash")
	}
}

func TestPartitionBucketKeyAlias(t *testing.T) {
	got := baseintegrator.PartitionBucketKey("prompts", "copilot")
	if got != "prompts" {
		t.Fatalf("expected 'prompts', got %q", got)
	}
}

func TestPartitionBucketKeyPassthrough(t *testing.T) {
	got := baseintegrator.PartitionBucketKey("agents", "cursor")
	if got != "agents_cursor" {
		t.Fatalf("expected 'agents_cursor', got %q", got)
	}
}

func TestValidateDeployPathTraversal(t *testing.T) {
	if baseintegrator.ValidateDeployPath("../etc/passwd", "/project", []string{".github/"}, nil) {
		t.Fatal("path traversal should be rejected")
	}
}

func TestValidateDeployPathDisallowedPrefix(t *testing.T) {
	if baseintegrator.ValidateDeployPath(".hidden/secret", "/project", []string{".github/"}, nil) {
		t.Fatal("disallowed prefix should be rejected")
	}
}

func TestCleanupEmptyParents(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "a", "b")
	os.MkdirAll(sub, 0755)
	f := filepath.Join(sub, "file.md")
	os.WriteFile(f, []byte("x"), 0644)
	os.Remove(f)
	baseintegrator.CleanupEmptyParents([]string{f}, dir)
	if _, err := os.Stat(sub); !os.IsNotExist(err) {
		t.Fatal("empty sub directory should have been removed")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("stop-at directory should NOT be removed")
	}
}

func TestSyncRemoveFilesLegacyGlob(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "foo-apm.agent.md")
	os.WriteFile(f, []byte("x"), 0644)
	stats := baseintegrator.SyncRemoveFiles(dir, nil, ".github/agents/", dir, "*-apm.agent.md", nil, nil)
	if stats.FilesRemoved != 1 {
		t.Fatalf("expected 1 removed, got %d", stats.FilesRemoved)
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Fatal("file should have been removed")
	}
}

func TestFindFilesByGlob(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.prompt.md"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "b.prompt.md"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "other.txt"), []byte("x"), 0644)
	results := baseintegrator.FindFilesByGlob(dir, "*.prompt.md", nil)
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}
