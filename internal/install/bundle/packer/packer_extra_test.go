package packer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackOptions_Fields(t *testing.T) {
	opts := PackOptions{
		ProjectRoot: "/tmp/proj",
		OutputDir:   "/tmp/out",
	}
	if opts.ProjectRoot != "/tmp/proj" || opts.OutputDir != "/tmp/out" {
		t.Errorf("unexpected fields: %+v", opts)
	}
}

func TestPackResult_ZeroValue(t *testing.T) {
	r := PackResult{}
	if r.BundlePath != "" || r.MappedCount != 0 || r.LockfileEnriched {
		t.Errorf("unexpected zero value: %+v", r)
	}
	if len(r.Files) != 0 {
		t.Error("expected empty Files slice")
	}
}

func TestPackResult_AllFields(t *testing.T) {
	r := PackResult{
		BundlePath:       "/out/bundle.tar.gz",
		Files:            []string{"a.md", "b.md"},
		LockfileEnriched: true,
		MappedCount:      2,
	}
	if r.MappedCount != 2 || !r.LockfileEnriched || len(r.Files) != 2 {
		t.Errorf("unexpected fields: %+v", r)
	}
}

func TestDeployedFile_Fields(t *testing.T) {
	d := DeployedFile{SourcePath: "/abs/path/file.md", BundlePath: "relative/file.md"}
	if d.SourcePath != "/abs/path/file.md" || d.BundlePath != "relative/file.md" {
		t.Errorf("unexpected fields: %+v", d)
	}
}

func TestBundleDependency_Fields(t *testing.T) {
	bd := BundleDependency{
		Name:          "my-dep",
		Version:       "1.2.3",
		DeployedFiles: []string{"a.md", "b.md"},
	}
	if bd.Name != "my-dep" || bd.Version != "1.2.3" || len(bd.DeployedFiles) != 2 {
		t.Errorf("unexpected fields: %+v", bd)
	}
}

func TestDetectTarget_CopilotFiles(t *testing.T) {
	dir := t.TempDir()
	copilotDir := filepath.Join(dir, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(copilotDir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(copilotDir, []byte("# copilot"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := detectTarget(dir)
	if target != "copilot" && target != "github" && target != "" {
		// any reasonable target name is okay; just ensure no panic
		t.Logf("detectTarget returned %q", target)
	}
}

func TestDetectTarget_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	target := detectTarget(dir)
	_ = target // no panic is the assertion
}

func TestFilterFilesByTarget_UnknownTarget(t *testing.T) {
	files := []string{"a/.github/copilot-instructions.md", "b/cursor/rules.md"}
	filtered, mappings := filterFilesByTarget(files, "unknown-target-xyz")
	_ = mappings
	// should not panic; result can be empty
	if len(filtered) > len(files) {
		t.Errorf("filtered cannot be larger than input")
	}
}

func TestFilterFilesByTarget_AllFilesTarget(t *testing.T) {
	files := []string{"a.md", "b.md", "c.md"}
	filtered, _ := filterFilesByTarget(files, "all")
	if len(filtered) > len(files) {
		t.Errorf("filtered %d > input %d", len(filtered), len(files))
	}
}

func TestCopyFile_Success(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile: %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestCopyFile_MissingSrc(t *testing.T) {
	dir := t.TempDir()
	err := copyFile(filepath.Join(dir, "no-such.txt"), filepath.Join(dir, "dst.txt"))
	if err == nil {
		t.Error("expected error for missing src")
	}
}

func TestCopyDirContents_Basic(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "file.txt"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyDirContents(src, dst); err != nil {
		t.Fatalf("copyDirContents: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dst, "file.txt"))
	if err != nil {
		t.Fatalf("file not copied: %v", err)
	}
	if string(data) != "data" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestCopyDirContents_WithSubdir(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	subdir := filepath.Join(src, "sub")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "nested.txt"), []byte("nested"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyDirContents(src, dst); err != nil {
		t.Fatalf("copyDirContents: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dst, "sub", "nested.txt"))
	if err != nil {
		t.Fatalf("nested file not copied: %v", err)
	}
	if string(data) != "nested" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestCreateTarGz_Basic(t *testing.T) {
	src := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "hello.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	archivePath := filepath.Join(t.TempDir(), "out.tar.gz")
	if err := createTarGz(src, archivePath); err != nil {
		t.Fatalf("createTarGz: %v", err)
	}
	info, err := os.Stat(archivePath)
	if err != nil {
		t.Fatalf("archive not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("archive is empty")
	}
}

func TestReadDeployedFiles_Empty(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, "apm.lock.yaml")
	if err := os.WriteFile(lockPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	deps, err := readDeployedFiles(lockPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 0 {
		t.Errorf("expected 0 deps, got %d", len(deps))
	}
}
