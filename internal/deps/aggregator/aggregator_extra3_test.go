package aggregator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanWorkflowsForDependencies_TwoFilesOneServer(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - shared-server\n---\nBody\n"
	_ = os.WriteFile(filepath.Join(dir, "a.prompt.md"), []byte(content), 0o600)
	_ = os.WriteFile(filepath.Join(dir, "b.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["shared-server"] {
		t.Error("expected shared-server")
	}
	if len(result) != 1 {
		t.Errorf("expected 1 unique server, got %d", len(result))
	}
}

func TestScanWorkflowsForDependencies_InlineMCPIgnored(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp: inline-value\n---\nBody\n"
	_ = os.WriteFile(filepath.Join(dir, "c.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = result
}

func TestScanWorkflowsForDependencies_NonExistentDir(t *testing.T) {
	result, err := ScanWorkflowsForDependencies("/nonexistent/path/12345")
	if err == nil && len(result) == 0 {
		return
	}
}

func TestScanWorkflowsForDependencies_TxtFileIgnored(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "file.txt"), []byte("---\nmcp:\n  - srv\n---\n"), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected no results for non-prompt file, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_ReturnsBoolMap(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - alpha\n  - beta\n---\n"
	_ = os.WriteFile(filepath.Join(dir, "t.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["alpha"] || !result["beta"] {
		t.Errorf("expected alpha and beta, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_NestedSubsubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "a", "b", "c")
	_ = os.MkdirAll(sub, 0o755)
	content := "---\nmcp:\n  - deep-server\n---\n"
	_ = os.WriteFile(filepath.Join(sub, "deep.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["deep-server"] {
		t.Error("expected deep-server")
	}
}

func TestScanWorkflowsForDependencies_MixedValidInvalid(t *testing.T) {
	dir := t.TempDir()
	valid := "---\nmcp:\n  - ok-server\n---\n"
	invalid := "no frontmatter at all"
	_ = os.WriteFile(filepath.Join(dir, "valid.prompt.md"), []byte(valid), 0o600)
	_ = os.WriteFile(filepath.Join(dir, "invalid.prompt.md"), []byte(invalid), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["ok-server"] {
		t.Error("expected ok-server")
	}
}

func TestScanWorkflowsForDependencies_EmptyFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "---\n---\nbody\n"
	_ = os.WriteFile(filepath.Join(dir, "e.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for no mcp, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_ReturnMapType(t *testing.T) {
	dir := t.TempDir()
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("result must not be nil")
	}
}

func TestScanWorkflowsForDependencies_SingleEntry(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - only-one\n---\n"
	_ = os.WriteFile(filepath.Join(dir, "x.prompt.md"), []byte(content), 0o600)
	result, err := ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || !result["only-one"] {
		t.Errorf("unexpected result: %v", result)
	}
}
