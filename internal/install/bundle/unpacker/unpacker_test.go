package unpacker

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "lockfile-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseBundleLockfile_Empty(t *testing.T) {
	path := writeTempFile(t, "")
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 0 {
		t.Errorf("expected 0 dependencies, got %d", len(lf.Dependencies))
	}
}

func TestParseBundleLockfile_Dependencies(t *testing.T) {
	content := `dependencies:
  - name: owner/repo
    version: "1.0.0"
    deployed_files:
      - .claude/skills/foo.md
      - .github/skills/bar.md
  - name: other/pkg
    version: "2.0.0"
`
	path := writeTempFile(t, content)
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(lf.Dependencies))
	}
	dep := lf.Dependencies[0]
	if dep.Name != "owner/repo" {
		t.Errorf("expected name 'owner/repo', got %q", dep.Name)
	}
	if dep.Version != `"1.0.0"` {
		t.Errorf("expected version '1.0.0', got %q", dep.Version)
	}
	if len(dep.DeployedFiles) != 2 {
		t.Errorf("expected 2 deployed files, got %d", len(dep.DeployedFiles))
	}
}

func TestParseBundleLockfile_PackMeta(t *testing.T) {
	content := `pack:
  format: plugin-v1
  target: claude
  packed_at: 2025-01-01T00:00:00Z
dependencies:
`
	path := writeTempFile(t, content)
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf.PackMeta["target"] != "claude" {
		t.Errorf("expected target 'claude', got %v", lf.PackMeta["target"])
	}
	if lf.PackMeta["format"] != "plugin-v1" {
		t.Errorf("expected format 'plugin-v1', got %v", lf.PackMeta["format"])
	}
}

func TestParseBundleLockfile_MissingFile(t *testing.T) {
	_, err := ParseBundleLockfile(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}
