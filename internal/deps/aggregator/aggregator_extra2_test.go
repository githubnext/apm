package aggregator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanWorkflowsForDependencies_SingleServerEntry(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - myserver\n---\nBody\n"
	_ = os.WriteFile(filepath.Join(dir, "tool.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["myserver"] {
		t.Error("expected 'myserver' to be present")
	}
}

func TestScanWorkflowsForDependencies_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "No frontmatter here\n"
	_ = os.WriteFile(filepath.Join(dir, "nofm.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for empty dir, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_MultipleMCPEntries(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - alpha\n  - beta\n  - gamma\n---\n"
	_ = os.WriteFile(filepath.Join(dir, "multi.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, name := range []string{"alpha", "beta", "gamma"} {
		if !result[name] {
			t.Errorf("expected %q in result", name)
		}
	}
}

func TestScanWorkflowsForDependencies_DuplicateServersDeduped(t *testing.T) {
	dir := t.TempDir()
	content1 := "---\nmcp:\n  - dup\n---\n"
	content2 := "---\nmcp:\n  - dup\n---\n"
	_ = os.WriteFile(filepath.Join(dir, "a.prompt.md"), []byte(content1), 0o600)
	_ = os.WriteFile(filepath.Join(dir, "b.prompt.md"), []byte(content2), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["dup"] {
		t.Error("expected 'dup' in merged result")
	}
	if len(result) != 1 {
		t.Errorf("expected exactly 1 unique server, got %d", len(result))
	}
}

func TestScanWorkflowsForDependencies_NotPromptMDIgnored(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - hidden\n---\n"
	_ = os.WriteFile(filepath.Join(dir, "file.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["hidden"] {
		t.Error("non-.prompt.md file should be ignored")
	}
}

func TestScanWorkflowsForDependencies_ReturnMap_NotNil(t *testing.T) {
	dir := t.TempDir()
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("result should be non-nil map")
	}
}

func TestScanWorkflowsForDependencies_SubdirPromptFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub", "nested")
	_ = os.MkdirAll(sub, 0o755)
	content := "---\nmcp:\n  - deep-server\n---\n"
	_ = os.WriteFile(filepath.Join(sub, "deep.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["deep-server"] {
		t.Error("expected 'deep-server' from subdirectory")
	}
}
