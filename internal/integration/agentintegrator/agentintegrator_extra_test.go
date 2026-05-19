package agentintegrator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/integration/agentintegrator"
	"github.com/githubnext/apm/internal/integration/targets"
)

func TestFindAgentFilesLegacyChatmodes(t *testing.T) {
	dir := t.TempDir()
	chatDir := filepath.Join(dir, ".apm", "chatmodes")
	os.MkdirAll(chatDir, 0755)
	os.WriteFile(filepath.Join(chatDir, "my.chatmode.md"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(chatDir, "notchat.md"), []byte("x"), 0644) // should be ignored
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) != 1 {
		t.Fatalf("expected 1 chatmode file, got %d", len(files))
	}
}

func TestFindAgentFilesApmAgentsMixed(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "agents")
	os.MkdirAll(apmDir, 0755)
	os.WriteFile(filepath.Join(apmDir, "a.agent.md"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(apmDir, "b.md"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(apmDir, "c.txt"), []byte("x"), 0644) // excluded
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(files), files)
	}
}

func TestFindAgentFilesDeduplicated(t *testing.T) {
	dir := t.TempDir()
	// root has agent.md
	os.WriteFile(filepath.Join(dir, "helper.agent.md"), []byte("x"), 0644)
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) != 1 {
		t.Fatalf("dedup: expected 1, got %d", len(files))
	}
}

func TestGetTargetFilenameForTargetChatmode(t *testing.T) {
	source := "/pkg/my.chatmode.md"
	target := targets.KnownTargets["copilot"]
	got := agentintegrator.GetTargetFilenameForTarget(source, target)
	// chatmode stem is "my", extension from copilot is .agent.md
	if got != "my.agent.md" {
		t.Fatalf("expected my.agent.md, got %q", got)
	}
}

func TestGetTargetFilenameForTargetNoMapping(t *testing.T) {
	source := "/pkg/foo.agent.md"
	// Use a target with no agents mapping
	tp := &targets.TargetProfile{
		Primitives: map[string]targets.PrimitiveMapping{},
	}
	got := agentintegrator.GetTargetFilenameForTarget(source, tp)
	// default ext is .agent.md
	if got != "foo.agent.md" {
		t.Fatalf("expected foo.agent.md, got %q", got)
	}
}

func TestPortableRelpath(t *testing.T) {
	got := agentintegrator.PortableRelpath("/a/b/c/file.md", "/a/b")
	if got != "c/file.md" {
		t.Fatalf("expected c/file.md, got %q", got)
	}
}

func TestCopyAgentMissingSource(t *testing.T) {
	dir := t.TempDir()
	_, err := agentintegrator.CopyAgent(filepath.Join(dir, "missing.md"), filepath.Join(dir, "dst.md"))
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestIntegrateAgentsForTargetNoAgentFiles(t *testing.T) {
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "pkg")
	os.MkdirAll(pkgDir, 0755)
	os.MkdirAll(filepath.Join(dir, ".github"), 0755)
	target := targets.KnownTargets["copilot"]
	result := agentintegrator.IntegrateAgentsForTarget(target, pkgDir, dir, false, nil, nil)
	if result.FilesIntegrated != 0 {
		t.Fatalf("expected 0 integrated for empty pkg, got %d", result.FilesIntegrated)
	}
}

func TestIntegrateAgentsForTargetNoTargetDir(t *testing.T) {
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "pkg")
	os.MkdirAll(pkgDir, 0755)
	os.WriteFile(filepath.Join(pkgDir, "agent.agent.md"), []byte("# A"), 0644)
	// Use cursor target which does NOT have AutoCreate=true
	target := targets.KnownTargets["cursor"]
	// Don't create .cursor dir -- target should skip
	result := agentintegrator.IntegrateAgentsForTarget(target, pkgDir, dir, false, nil, nil)
	if result.FilesIntegrated != 0 {
		t.Fatalf("expected 0 for missing target dir, got %d", result.FilesIntegrated)
	}
}

func TestIntegrateAgentsForTargetNoMapping(t *testing.T) {
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "pkg")
	os.MkdirAll(pkgDir, 0755)
	tp := &targets.TargetProfile{
		Primitives: map[string]targets.PrimitiveMapping{},
	}
	result := agentintegrator.IntegrateAgentsForTarget(tp, pkgDir, dir, false, nil, nil)
	if result.FilesIntegrated != 0 {
		t.Fatal("expected 0 without agent mapping")
	}
}

func TestIntegrateAgentsCodexTarget(t *testing.T) {
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "pkg")
	os.MkdirAll(pkgDir, 0755)
	content := "---\nname: MyAgent\ndescription: Does stuff\n---\n# Body\nHello world."
	os.WriteFile(filepath.Join(pkgDir, "myagent.agent.md"), []byte(content), 0644)
	os.MkdirAll(filepath.Join(dir, ".codex"), 0755)
	target := targets.KnownTargets["codex"]
	result := agentintegrator.IntegrateAgentsForTarget(target, pkgDir, dir, false, nil, nil)
	if result.FilesIntegrated != 1 {
		t.Fatalf("expected 1 codex agent integrated, got %d", result.FilesIntegrated)
	}
	// Verify TOML output exists
	tomlPath := filepath.Join(dir, ".codex", "agents", "myagent.toml")
	data, err := os.ReadFile(tomlPath)
	if err != nil {
		t.Fatalf("expected toml output: %v", err)
	}
	if !strings.Contains(string(data), `name = "MyAgent"`) {
		t.Fatalf("toml missing name: %s", string(data))
	}
}

func TestSyncForTargetNoMapping(t *testing.T) {
	dir := t.TempDir()
	tp := &targets.TargetProfile{
		Primitives: map[string]targets.PrimitiveMapping{},
	}
	stats := agentintegrator.SyncForTarget(tp, dir, nil)
	if stats.FilesRemoved != 0 {
		t.Fatalf("expected 0 removed, got %d", stats.FilesRemoved)
	}
}
