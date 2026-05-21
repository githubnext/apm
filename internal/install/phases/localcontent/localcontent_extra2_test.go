package localcontent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectHasRootPrimitives_WithApmDir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm")
	if err := os.MkdirAll(apmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if !ProjectHasRootPrimitives(dir) {
		t.Error("expected true when .apm dir exists")
	}
}

func TestProjectHasRootPrimitives_ApmIsFile(t *testing.T) {
	dir := t.TempDir()
	// .apm is a file, not a directory
	apmFile := filepath.Join(dir, ".apm")
	if err := os.WriteFile(apmFile, []byte("not a dir"), 0o644); err != nil {
		t.Fatal(err)
	}
	if ProjectHasRootPrimitives(dir) {
		t.Error("expected false when .apm is a file, not a directory")
	}
}

func TestHasLocalApmContent_EmptyApmDir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm")
	if err := os.MkdirAll(apmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if HasLocalApmContent(dir) {
		t.Error("expected false when .apm is empty")
	}
}

func TestHasLocalApmContent_OnlyUnrecognizedSubdir(t *testing.T) {
	dir := t.TempDir()
	// Create .apm/custom-unknown which is not in primitiveDirs
	sub := filepath.Join(dir, ".apm", "custom-unknown")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(sub, "file.md"), []byte("x"), 0o644)
	if HasLocalApmContent(dir) {
		t.Error("expected false when only unrecognized subdirs have files")
	}
}

func TestHasLocalApmContent_SkillsSubdirWithFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".apm", "skills")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(sub, "my-skill.md"), []byte("skill"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true when skills dir has a file")
	}
}

func TestHasLocalApmContent_HooksSubdirWithFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".apm", "hooks")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(sub, "hook.sh"), []byte("#!/bin/sh"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true when hooks dir has a file")
	}
}

func TestHasLocalApmContent_CommandsSubdirWithFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".apm", "commands")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(sub, "cmd.md"), []byte("cmd"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true when commands dir has a file")
	}
}

func TestHasLocalApmContent_ApmDirMissing(t *testing.T) {
	dir := t.TempDir()
	// no .apm dir at all
	if HasLocalApmContent(dir) {
		t.Error("expected false when no .apm dir exists")
	}
}

func TestHasLocalApmContent_NestedFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, ".apm", "instructions", "subdir")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(sub, "deep.instructions.md"), []byte("x"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true when nested file exists under instructions")
	}
}
