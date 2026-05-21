package aggregator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/deps/aggregator"
)

func TestScanWorkflowsForDependencies_InvalidDirReturnsEmpty(t *testing.T) {
	result, err := aggregator.ScanWorkflowsForDependencies("/nonexistent/path/xyz123")
	if err != nil {
		t.Fatalf("expected no error for missing dir, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestScanWorkflowsForDependencies_NonPromptFilesIgnored(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# readme"), 0o644)
	result, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}

func TestScanWorkflowsForDependencies_SingleServer(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - my-server\n---\nBody.\n"
	os.WriteFile(filepath.Join(dir, "test.prompt.md"), []byte(content), 0o644)
	result, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["my-server"] {
		t.Errorf("expected my-server in results, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_MultipleServers(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - svc-a\n  - svc-b\n---\nBody.\n"
	os.WriteFile(filepath.Join(dir, "test.prompt.md"), []byte(content), 0o644)
	result, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["svc-a"] {
		t.Errorf("expected svc-a, got %v", result)
	}
	if !result["svc-b"] {
		t.Errorf("expected svc-b, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_EmptyPromptFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "empty.prompt.md"), []byte(""), 0o644)
	result, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for empty prompt file")
	}
}

func TestScanWorkflowsForDependencies_NestedDirScanned(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	os.Mkdir(sub, 0o755)
	content := "---\nmcp:\n  - nested-svc\n---\n"
	os.WriteFile(filepath.Join(sub, "nested.prompt.md"), []byte(content), 0o644)
	result, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["nested-svc"] {
		t.Errorf("expected nested-svc in results, got %v", result)
	}
}

func TestScanWorkflowsForDependencies_DedupSameServerTwoFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.prompt.md", "b.prompt.md"} {
		content := "---\nmcp:\n  - dup-svc\n---\n"
		os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
	}
	result, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 unique server, got %d", len(result))
	}
}
