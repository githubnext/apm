package unpacker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBundleLockfile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.lock")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf == nil {
		t.Fatal("expected non-nil lockfile")
	}
}

func TestParseBundleLockfile_WhitespaceOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.lock")
	if err := os.WriteFile(path, []byte("   \n  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf == nil {
		t.Fatal("expected non-nil lockfile")
	}
}

func TestBundleLockfile_Fields(t *testing.T) {
	lf := BundleLockfile{
		Dependencies: []LockEntry{{Name: "dep-a", Version: "1.0.0"}},
		PackMeta:     map[string]interface{}{"name": "test-bundle", "version": "1.2.3"},
	}
	if len(lf.Dependencies) != 1 {
		t.Errorf("Dependencies len: %d", len(lf.Dependencies))
	}
	if lf.Dependencies[0].Name != "dep-a" {
		t.Errorf("Dep[0].Name: %q", lf.Dependencies[0].Name)
	}
	if lf.PackMeta["name"] != "test-bundle" {
		t.Errorf("PackMeta name: %v", lf.PackMeta["name"])
	}
}

func TestLockEntry_AllFields(t *testing.T) {
	e := LockEntry{
		Name:          "my-dep",
		Version:       "v1.2.3",
		DeployedFiles: []string{"/dist/a.txt", "/dist/b.txt"},
	}
	if e.Name != "my-dep" {
		t.Errorf("Name: %q", e.Name)
	}
	if e.Version != "v1.2.3" {
		t.Errorf("Version: %q", e.Version)
	}
	if len(e.DeployedFiles) != 2 {
		t.Errorf("DeployedFiles len: %d", len(e.DeployedFiles))
	}
}

func TestUnpackResult_ZeroValue(t *testing.T) {
	var r UnpackResult
	if r.ExtractedDir != "" {
		t.Errorf("ExtractedDir should be empty, got %q", r.ExtractedDir)
	}
	if r.Files != nil {
		t.Error("Files should be nil")
	}
}

func TestParseBundleLockfile_MixedContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.lock")
	content := `# Bundle Lock File
[pack]
name = mixed-bundle
version = 0.9.0

[deps]
pkg-a = ref=main source=github.com/org/pkg-a
pkg-b = ref=v2.0 source=github.com/org/pkg-b
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = lf
}

func TestUnpackBundle_InvalidTarFile(t *testing.T) {
	dir := t.TempDir()
	fakeBundle := filepath.Join(dir, "fake.tar.gz")
	if err := os.WriteFile(fakeBundle, []byte("not a real tar file"), 0o644); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(dir, "out")
	_, err := UnpackBundle(fakeBundle, outDir, true, false)
	if err == nil {
		t.Error("expected error for invalid tar.gz file")
	}
}

func TestUnpackBundle_DryRun_NoOutputDir(t *testing.T) {
	dir := t.TempDir()
	fakeBundle := filepath.Join(dir, "fake.tar.gz")
	if err := os.WriteFile(fakeBundle, []byte("not a real tar file"), 0o644); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(dir, "nonexistent-out")
	_, err := UnpackBundle(fakeBundle, outDir, true, true)
	if err == nil {
		t.Error("expected error even in dry-run for invalid tar")
	}
}

func TestBundleLockfile_NoDepsNoFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.lock")
	content := "[pack]\nname = empty\nversion = 1.0.0\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 0 {
		t.Errorf("expected 0 deps, got %d", len(lf.Dependencies))
	}
}

func TestParseBundleLockfile_SingleDep(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "apm.lock.yaml")
	content := "pack:\n  name: single-dep\n  version: 0.1.0\ndependencies:\n  - name: alpha\n    version: main\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseBundleLockfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Dependencies) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(lf.Dependencies))
	}
	if lf.Dependencies[0].Name != "alpha" {
		t.Errorf("dep name: %q", lf.Dependencies[0].Name)
	}
}
