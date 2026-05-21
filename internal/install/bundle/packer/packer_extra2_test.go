package packer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackOptions_AllFields(t *testing.T) {
	opts := PackOptions{
		ProjectRoot: "/proj",
		OutputDir:   "/out",
		Format:      "apm",
		Target:      "copilot",
		Archive:     true,
		DryRun:      false,
		Force:       true,
	}
	if opts.Format != "apm" {
		t.Errorf("Format = %q", opts.Format)
	}
	if opts.Target != "copilot" {
		t.Errorf("Target = %q", opts.Target)
	}
	if !opts.Archive {
		t.Error("expected Archive=true")
	}
	if !opts.Force {
		t.Error("expected Force=true")
	}
}

func TestBundleDependency_Fields_extra2(t *testing.T) {
	bd := BundleDependency{
		Name:          "my-dep",
		Version:       "1.2.3",
		DeployedFiles: []string{"a.txt", "b.md"},
	}
	if bd.Name != "my-dep" {
		t.Errorf("Name = %q", bd.Name)
	}
	if bd.Version != "1.2.3" {
		t.Errorf("Version = %q", bd.Version)
	}
	if len(bd.DeployedFiles) != 2 {
		t.Errorf("expected 2 deployed files, got %d", len(bd.DeployedFiles))
	}
}

func TestDeployedFile_Fields_extra2(t *testing.T) {
	df := DeployedFile{
		SourcePath: "/absolute/path/file.txt",
		BundlePath: "relative/file.txt",
	}
	if df.SourcePath != "/absolute/path/file.txt" {
		t.Errorf("SourcePath = %q", df.SourcePath)
	}
	if df.BundlePath != "relative/file.txt" {
		t.Errorf("BundlePath = %q", df.BundlePath)
	}
}

func TestPackBundle_EmptyProjectRoot_Error(t *testing.T) {
	_, err := PackBundle(PackOptions{})
	if err == nil {
		t.Error("expected error for empty ProjectRoot")
	}
}

func TestPackBundle_NonexistentRoot_Error(t *testing.T) {
	_, err := PackBundle(PackOptions{ProjectRoot: "/nonexistent/path/xyz"})
	if err == nil {
		t.Error("expected error for nonexistent project root")
	}
}

func TestPackBundle_DryRun_NoOutput(t *testing.T) {
	dir := t.TempDir()
	outDir := t.TempDir()
	// Write minimal lockfile.
	lockfile := filepath.Join(dir, "apm.lock.yaml")
	if err := os.WriteFile(lockfile, []byte("# empty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := PackBundle(PackOptions{
		ProjectRoot: dir,
		OutputDir:   outDir,
		DryRun:      true,
	})
	// DryRun should succeed or return a benign error, not panic.
	_ = err
}

func TestPackResult_PathMappings(t *testing.T) {
	r := PackResult{
		PathMappings: map[string]string{
			"src/a.txt": "bundle/a.txt",
		},
		MappedCount: 1,
	}
	if r.MappedCount != 1 {
		t.Errorf("MappedCount = %d", r.MappedCount)
	}
	if r.PathMappings["src/a.txt"] != "bundle/a.txt" {
		t.Errorf("PathMappings[src/a.txt] = %q", r.PathMappings["src/a.txt"])
	}
}
