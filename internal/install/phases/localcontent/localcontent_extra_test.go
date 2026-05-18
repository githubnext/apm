package localcontent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectHasRootPrimitives_NonExistentDir(t *testing.T) {
	if ProjectHasRootPrimitives("/nonexistent/path/abc123") {
		t.Error("expected false for non-existent root")
	}
}

func TestHasLocalApmContent_MultipleSubdirs(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm")
	// Create multiple recognized subdirs, only one with a file
	for _, sub := range []string{"skills", "instructions", "chatmodes"} {
		os.MkdirAll(filepath.Join(apmDir, sub), 0o755)
	}
	os.WriteFile(filepath.Join(apmDir, "instructions", "lint.instructions.md"), []byte("x"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true when instructions has a file")
	}
}

func TestHasLocalApmContent_AgentsSubdir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "agents")
	os.MkdirAll(apmDir, 0o755)
	os.WriteFile(filepath.Join(apmDir, "my.agent.md"), []byte("agent"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true for agents subdir with file")
	}
}

func TestHasLocalApmContent_PromptsSubdir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "prompts")
	os.MkdirAll(apmDir, 0o755)
	os.WriteFile(filepath.Join(apmDir, "foo.prompt.md"), []byte("x"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true for prompts subdir with file")
	}
}

func TestHasLocalApmContent_HooksSubdir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "hooks")
	os.MkdirAll(apmDir, 0o755)
	os.WriteFile(filepath.Join(apmDir, "pre-install.sh"), []byte("#!/bin/bash"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true for hooks subdir with file")
	}
}

func TestHasLocalApmContent_CommandsSubdir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "commands")
	os.MkdirAll(apmDir, 0o755)
	os.WriteFile(filepath.Join(apmDir, "custom.md"), []byte("x"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true for commands subdir with file")
	}
}

func TestHasLocalApmContent_ApmFileNotDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".apm"), []byte("x"), 0o644)
	if HasLocalApmContent(dir) {
		t.Error("expected false when .apm is a file, not a directory")
	}
}

func TestHasLocalApmContent_DeepNestedFile(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, ".apm", "skills", "sub", "deep")
	os.MkdirAll(nested, 0o755)
	os.WriteFile(filepath.Join(nested, "SKILL.md"), []byte("x"), 0o644)
	if !HasLocalApmContent(dir) {
		t.Error("expected true for deeply nested file in recognized subdir")
	}
}
