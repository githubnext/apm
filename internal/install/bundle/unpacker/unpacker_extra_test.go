package unpacker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBundleLockfile_CommentLines(t *testing.T) {
	content := "# This is a comment\ndependencies:\n  # another comment\n  - name: my/pkg\n    version: \"1.0\"\n"
	path := writeTempFile(t, content)
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 1 {
		t.Errorf("expected 1 dep, got %d", len(lf.Dependencies))
	}
	if lf.Dependencies[0].Name != "my/pkg" {
		t.Errorf("unexpected dep name: %q", lf.Dependencies[0].Name)
	}
}

func TestParseBundleLockfile_DeployedFiles(t *testing.T) {
	content := "dependencies:\n  - name: a/b\n    version: v1\n    deployed_files:\n      - .claude/skills/x.md\n      - .github/skills/y.md\n      - dir/file.txt\n"
	path := writeTempFile(t, content)
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(lf.Dependencies))
	}
	if len(lf.Dependencies[0].DeployedFiles) != 3 {
		t.Errorf("expected 3 deployed files, got %d", len(lf.Dependencies[0].DeployedFiles))
	}
}

func TestParseBundleLockfile_MultiDepsNoFiles(t *testing.T) {
	content := "dependencies:\n  - name: a/b\n    version: v1\n  - name: c/d\n    version: v2\n  - name: e/f\n    version: v3\n"
	path := writeTempFile(t, content)
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 3 {
		t.Errorf("expected 3 deps, got %d: %v", len(lf.Dependencies), lf.Dependencies)
	}
}

func TestParseBundleLockfile_StructFields(t *testing.T) {
	result := &UnpackResult{
		ExtractedDir:     "/tmp/out",
		Verified:         true,
		SkippedCount:     2,
		SecurityWarnings: 1,
		SecurityCritical: 0,
	}
	if result.ExtractedDir != "/tmp/out" {
		t.Errorf("ExtractedDir mismatch")
	}
	if !result.Verified {
		t.Error("Verified should be true")
	}
	if result.SkippedCount != 2 {
		t.Errorf("SkippedCount mismatch: %d", result.SkippedCount)
	}
}

func TestLockEntry_Fields(t *testing.T) {
	e := LockEntry{
		Name:          "owner/repo",
		Version:       "v2.0.0",
		DeployedFiles: []string{"a.md", "b.md"},
	}
	if e.Name != "owner/repo" {
		t.Errorf("Name mismatch: %q", e.Name)
	}
	if e.Version != "v2.0.0" {
		t.Errorf("Version mismatch: %q", e.Version)
	}
	if len(e.DeployedFiles) != 2 {
		t.Errorf("DeployedFiles length mismatch: %d", len(e.DeployedFiles))
	}
}

func TestBundleLockfile_EmptyStruct(t *testing.T) {
	lf := &BundleLockfile{
		Dependencies: []LockEntry{},
		PackMeta:     map[string]interface{}{},
		RawData:      map[string]interface{}{},
	}
	if len(lf.Dependencies) != 0 {
		t.Errorf("expected empty deps")
	}
}

func TestUnpackBundle_NonexistentPath(t *testing.T) {
	outDir := t.TempDir()
	_, err := UnpackBundle("/nonexistent/bundle.tar.gz", outDir, false, false)
	if err == nil {
		t.Error("expected error for nonexistent bundle path")
	}
}

func TestUnpackBundle_DryRunNotExist(t *testing.T) {
	outDir := t.TempDir()
	_, err := UnpackBundle("/nonexistent/bundle.tar.gz", outDir, false, true)
	if err == nil {
		t.Error("expected error for nonexistent path even in dry-run")
	}
}

func TestParseBundleLockfile_OnlyPackSection(t *testing.T) {
	content := "pack:\n  version: 1.0\n  name: myapp\n"
	f, err := os.CreateTemp(t.TempDir(), "lf-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	lf, err := ParseBundleLockfile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 0 {
		t.Errorf("expected 0 deps, got %d", len(lf.Dependencies))
	}
}

func TestUnpackResult_FilesList(t *testing.T) {
	result := &UnpackResult{
		Files: []string{"a.md", "b.md", "c.md"},
	}
	if len(result.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(result.Files))
	}
}

func TestParseBundleLockfile_WritesAndReads(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, "apm.lock.yaml")
	content := "dependencies:\n  - name: test/pkg\n    version: abc123\n    deployed_files:\n      - skills/test.md\n"
	if err := os.WriteFile(lockPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseBundleLockfile(lockPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(lf.Dependencies))
	}
	if lf.Dependencies[0].DeployedFiles[0] != "skills/test.md" {
		t.Errorf("unexpected deployed file: %q", lf.Dependencies[0].DeployedFiles[0])
	}
}
