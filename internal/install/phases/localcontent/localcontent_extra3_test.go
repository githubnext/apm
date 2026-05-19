package localcontent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/install/phases/localcontent"
)

func TestHasLocalApmContent_InstructionsSubdirWithFile(t *testing.T) {
	root := t.TempDir()
	instDir := filepath.Join(root, ".apm", "instructions")
	_ = os.MkdirAll(instDir, 0o755)
	_ = os.WriteFile(filepath.Join(instDir, "rules.md"), []byte("rules"), 0o644)
	if !localcontent.HasLocalApmContent(root) {
		t.Error("expected HasLocalApmContent=true for instructions subdir")
	}
}

func TestHasLocalApmContent_ChatmodesSubdirWithFile(t *testing.T) {
	root := t.TempDir()
	chatDir := filepath.Join(root, ".apm", "chatmodes")
	_ = os.MkdirAll(chatDir, 0o755)
	_ = os.WriteFile(filepath.Join(chatDir, "mode.yml"), []byte("mode"), 0o644)
	if !localcontent.HasLocalApmContent(root) {
		t.Error("expected HasLocalApmContent=true for chatmodes subdir")
	}
}

func TestHasLocalApmContent_AllSubdirsEmpty(t *testing.T) {
	root := t.TempDir()
	for _, sub := range []string{"skills", "instructions", "chatmodes", "agents", "prompts", "hooks", "commands"} {
		_ = os.MkdirAll(filepath.Join(root, ".apm", sub), 0o755)
	}
	if localcontent.HasLocalApmContent(root) {
		t.Error("expected HasLocalApmContent=false when all subdirs are empty")
	}
}

func TestHasLocalApmContent_OnlyDirInSubdir(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, ".apm", "skills", "subsubdir")
	_ = os.MkdirAll(nested, 0o755)
	// No files, only a directory -- should return false
	if localcontent.HasLocalApmContent(root) {
		t.Error("expected HasLocalApmContent=false when subdir contains only dirs, no files")
	}
}

func TestProjectHasRootPrimitives_EmptyApmDir(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".apm"), 0o755)
	if !localcontent.ProjectHasRootPrimitives(root) {
		t.Error("expected true: .apm dir exists even if empty")
	}
}

func TestHasLocalApmContent_MultipleFilesInDifferentSubdirs(t *testing.T) {
	root := t.TempDir()
	for _, sub := range []string{"skills", "agents"} {
		dir := filepath.Join(root, ".apm", sub)
		_ = os.MkdirAll(dir, 0o755)
		_ = os.WriteFile(filepath.Join(dir, "item.md"), []byte("content"), 0o644)
	}
	if !localcontent.HasLocalApmContent(root) {
		t.Error("expected HasLocalApmContent=true")
	}
}

func TestHasLocalApmContent_PromptsSubdirNoFile(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".apm", "prompts"), 0o755)
	if localcontent.HasLocalApmContent(root) {
		t.Error("expected false when prompts subdir is empty")
	}
}

func TestProjectHasRootPrimitives_NonExistentRoot(t *testing.T) {
	if localcontent.ProjectHasRootPrimitives("/nonexistent/xyz/abc/123") {
		t.Error("expected false for non-existent root")
	}
}
