package agentintegrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/integration/agentintegrator"
	"github.com/githubnext/apm/internal/integration/targets"
)

func TestFindAgentFilesEmpty(t *testing.T) {
	dir := t.TempDir()
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}
}

func TestFindAgentFilesRoot(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "my.agent.md"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "chat.chatmode.md"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "other.txt"), []byte("x"), 0644)
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) != 2 {
		t.Fatalf("expected 2 agent files, got %d", len(files))
	}
}

func TestFindAgentFilesApmAgentsDir(t *testing.T) {
	dir := t.TempDir()
	apmDir := filepath.Join(dir, ".apm", "agents")
	os.MkdirAll(apmDir, 0755)
	os.WriteFile(filepath.Join(apmDir, "helper.agent.md"), []byte("x"), 0644)
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
}

func TestGetTargetFilenameForTargetCopilot(t *testing.T) {
	source := "/pkg/.apm/agents/reviewer.agent.md"
	target := targets.KnownTargets["copilot"]
	got := agentintegrator.GetTargetFilenameForTarget(source, target)
	if got != "reviewer.agent.md" {
		t.Fatalf("expected reviewer.agent.md, got %q", got)
	}
}

func TestGetTargetFilenameForTargetClaude(t *testing.T) {
	source := "/pkg/.apm/agents/reviewer.agent.md"
	target := targets.KnownTargets["claude"]
	got := agentintegrator.GetTargetFilenameForTarget(source, target)
	// claude uses .md extension
	if got != "reviewer.md" {
		t.Fatalf("expected reviewer.md, got %q", got)
	}
}

func TestCopyAgent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.agent.md")
	dst := filepath.Join(dir, "dst.agent.md")
	os.WriteFile(src, []byte("# Agent\nHello"), 0644)
	n, err := agentintegrator.CopyAgent(src, dst)
	if err != nil {
		t.Fatalf("copy error: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 links, got %d", n)
	}
	data, _ := os.ReadFile(dst)
	if string(data) != "# Agent\nHello" {
		t.Fatal("content mismatch")
	}
}

func TestIntegrateAgentsForTarget(t *testing.T) {
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "pkg")
	os.MkdirAll(pkgDir, 0755)
	os.WriteFile(filepath.Join(pkgDir, "helper.agent.md"), []byte("# Helper"), 0644)

	// Create .github dir so copilot target is active
	os.MkdirAll(filepath.Join(dir, ".github"), 0755)

	target := targets.KnownTargets["copilot"]
	result := agentintegrator.IntegrateAgentsForTarget(target, pkgDir, dir, false, nil, nil)
	if result.FilesIntegrated != 1 {
		t.Fatalf("expected 1 integrated, got %d", result.FilesIntegrated)
	}
	expected := filepath.Join(dir, ".github", "agents", "helper.agent.md")
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Fatalf("expected output file at %s", expected)
	}
}

func TestSyncForTarget(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, ".github", "agents")
	os.MkdirAll(agentsDir, 0755)
	f := filepath.Join(agentsDir, "foo-apm.agent.md")
	os.WriteFile(f, []byte("x"), 0644)

	target := targets.KnownTargets["copilot"]
	stats := agentintegrator.SyncForTarget(target, dir, nil)
	if stats.FilesRemoved != 1 {
		t.Fatalf("expected 1 removed, got %d", stats.FilesRemoved)
	}
}
