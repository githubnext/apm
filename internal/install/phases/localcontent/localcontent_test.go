package localcontent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectHasRootPrimitives_Present(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".apm"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !ProjectHasRootPrimitives(dir) {
		t.Error("expected true when .apm dir exists")
	}
}

func TestProjectHasRootPrimitives_Absent(t *testing.T) {
	dir := t.TempDir()
	if ProjectHasRootPrimitives(dir) {
		t.Error("expected false when .apm dir absent")
	}
}

func TestProjectHasRootPrimitives_FileNotDir(t *testing.T) {
	dir := t.TempDir()
	// Create .apm as a file instead of a directory
	if err := os.WriteFile(filepath.Join(dir, ".apm"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if ProjectHasRootPrimitives(dir) {
		t.Error("expected false when .apm is a file")
	}
}

func TestHasLocalApmContent_NoApmDir(t *testing.T) {
	dir := t.TempDir()
	if HasLocalApmContent(dir) {
		t.Error("expected false when .apm absent")
	}
}

func TestHasLocalApmContent_EmptySubdirs(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm")
	if err := os.Mkdir(apmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Create recognized subdir but no files in it
	if err := os.Mkdir(filepath.Join(apmDir, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if HasLocalApmContent(dir) {
		t.Error("expected false when recognized subdirs are empty")
	}
}

func TestHasLocalApmContent_WithFile(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm")
	skillsDir := filepath.Join(apmDir, "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillsDir, "SKILL.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !HasLocalApmContent(dir) {
		t.Error("expected true when skill file is present")
	}
}

func TestHasLocalApmContent_UnrecognizedSubdir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm")
	unrecognized := filepath.Join(apmDir, "custom")
	if err := os.MkdirAll(unrecognized, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(unrecognized, "file.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Files in unrecognized subdirs should NOT count
	if HasLocalApmContent(dir) {
		t.Error("expected false for unrecognized subdir")
	}
}
