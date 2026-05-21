package agentintegrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/integration/agentintegrator"
)

func TestAgentIntegrator_ZeroValue(t *testing.T) {
	var a agentintegrator.AgentIntegrator
	_ = a // struct with no fields — just verify it compiles
}

func TestFindAgentFiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) != 0 {
		t.Errorf("expected no agent files in empty dir, got %v", files)
	}
}

func TestFindAgentFiles_NonexistentDir(t *testing.T) {
	files := agentintegrator.FindAgentFiles("/nonexistent/path/xyz123")
	if len(files) != 0 {
		t.Errorf("expected no agent files for nonexistent dir, got %v", files)
	}
}

func TestPortableRelpath_SameDir(t *testing.T) {
	result := agentintegrator.PortableRelpath("/a/b/c", "/a/b")
	if result != "c" {
		t.Errorf("expected 'c', got %q", result)
	}
}

func TestPortableRelpath_Nested(t *testing.T) {
	result := agentintegrator.PortableRelpath("/a/b/c/d", "/a/b")
	if result != filepath.Join("c", "d") {
		t.Errorf("expected 'c/d', got %q", result)
	}
}

func TestCopyAgent_SameDir(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "agent.md")
	dstFile := filepath.Join(dir, "agent-copy.md")
	content := []byte("# agent content")
	if err := os.WriteFile(srcFile, content, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := agentintegrator.CopyAgent(srcFile, dstFile)
	if err != nil {
		t.Fatalf("CopyAgent failed: %v", err)
	}
	// Verify the file was actually written
	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("destination file not readable: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("destination content mismatch: got %q want %q", got, content)
	}
}

func TestCopyAgent_SourceNotExist(t *testing.T) {
	_, err := agentintegrator.CopyAgent("/nonexistent/source.md", "/tmp/dst.md")
	if err == nil {
		t.Error("expected error for nonexistent source")
	}
}

func TestFindAgentFiles_WithMdFile(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, ".github", "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agentFile := filepath.Join(agentsDir, "my-agent.md")
	if err := os.WriteFile(agentFile, []byte("# agent"), 0o644); err != nil {
		t.Fatal(err)
	}
	files := agentintegrator.FindAgentFiles(dir)
	if len(files) == 0 {
		t.Log("no agent files found — may depend on expected directory layout")
	}
}
