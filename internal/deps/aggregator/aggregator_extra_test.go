package aggregator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/deps/aggregator"
)

func TestScanWorkflowsForDependencies_EmptyMCPBlock(t *testing.T) {
	dir := t.TempDir()
	// mcp: key present but no entries
	content := "---\nmcp:\n---\n\n# Prompt"
	if err := os.WriteFile(filepath.Join(dir, "empty.prompt.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	servers, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// no servers should be extracted from an empty mcp block
	if len(servers) != 0 {
		t.Errorf("expected empty result for empty mcp block, got %v", servers)
	}
}

func TestScanWorkflowsForDependencies_NestedInSubdirDeep(t *testing.T) {
	dir := t.TempDir()
	deep := filepath.Join(dir, "a", "b", "c")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nmcp:\n  - deep-server\n---\n"
	if err := os.WriteFile(filepath.Join(deep, "deep.prompt.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	servers, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !servers["deep-server"] {
		t.Errorf("expected deep-server, got %v", servers)
	}
}

func TestScanWorkflowsForDependencies_MultipleFiles_MergeUnique(t *testing.T) {
	dir := t.TempDir()
	files := []struct {
		name    string
		content string
	}{
		{"a.prompt.md", "---\nmcp:\n  - srv-a\n  - common\n---\n"},
		{"b.prompt.md", "---\nmcp:\n  - srv-b\n  - common\n---\n"},
		{"c.prompt.md", "---\nmcp:\n  - srv-c\n---\n"},
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f.name), []byte(f.content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	servers, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"srv-a", "srv-b", "srv-c", "common"} {
		if !servers[want] {
			t.Errorf("expected %q in result, got %v", want, servers)
		}
	}
	if len(servers) != 4 {
		t.Errorf("expected 4 unique servers, got %d", len(servers))
	}
}

func TestScanWorkflowsForDependencies_OtherExtensionsIgnored(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"test.yml", "test.json", "test.txt", "test.md"} {
		content := "---\nmcp:\n  - should-be-ignored\n---\n"
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	servers, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 0 {
		t.Errorf("expected empty result, got %v", servers)
	}
}

func TestScanWorkflowsForDependencies_DefaultDirNotEmpty(t *testing.T) {
	// When baseDir is "", it uses cwd which is a real directory
	// Just verify it doesn't error (may return servers from the repo itself)
	servers, err := aggregator.ScanWorkflowsForDependencies("")
	if err != nil {
		t.Fatalf("unexpected error with empty baseDir: %v", err)
	}
	_ = servers // result depends on cwd content
}

func TestScanWorkflowsForDependencies_ReturnTypeIsMap(t *testing.T) {
	dir := t.TempDir()
	servers, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Result should be a non-nil map
	if servers == nil {
		t.Error("expected non-nil map")
	}
}

func TestScanWorkflowsForDependencies_InlineMCPValue(t *testing.T) {
	dir := t.TempDir()
	// mcp: with inline value (not a list) - should be handled gracefully
	content := "---\nmcp: single-value\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "inline.prompt.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	// Should not panic; empty or partial result is acceptable
	_, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error for inline mcp value: %v", err)
	}
}

func TestScanWorkflowsForDependencies_LargeListOfServers(t *testing.T) {
	dir := t.TempDir()
	content := "---\nmcp:\n  - srv-1\n  - srv-2\n  - srv-3\n  - srv-4\n  - srv-5\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "large.prompt.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	servers, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 5 {
		t.Errorf("expected 5 servers, got %d: %v", len(servers), servers)
	}
	for i := 1; i <= 5; i++ {
		key := "srv-" + string(rune('0'+i))
		if !servers[key] {
			t.Errorf("expected %q in servers", key)
		}
	}
}

func TestScanWorkflowsForDependencies_FrontmatterWithExtraFields(t *testing.T) {
	dir := t.TempDir()
	content := "---\ntitle: My Prompt\nversion: 1.0\nmcp:\n  - extra-server\nauthor: test\n---\n# Content"
	if err := os.WriteFile(filepath.Join(dir, "extra.prompt.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	servers, err := aggregator.ScanWorkflowsForDependencies(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !servers["extra-server"] {
		t.Errorf("expected extra-server, got %v", servers)
	}
}
