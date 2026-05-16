package aggregator_test

import (
"os"
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/deps/aggregator"
)

func TestScanWorkflowsForDependencies_Empty(t *testing.T) {
dir := t.TempDir()
servers, err := aggregator.ScanWorkflowsForDependencies(dir)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(servers) != 0 {
t.Errorf("expected empty result, got %v", servers)
}
}

func TestScanWorkflowsForDependencies_WithMCP(t *testing.T) {
dir := t.TempDir()
content := "---\nmcp:\n  - my-mcp-server\n  - other-server\n---\n\n# Prompt"
if err := os.WriteFile(filepath.Join(dir, "test.prompt.md"), []byte(content), 0o644); err != nil {
t.Fatal(err)
}
servers, err := aggregator.ScanWorkflowsForDependencies(dir)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !servers["my-mcp-server"] {
t.Errorf("expected my-mcp-server in result, got %v", servers)
}
if !servers["other-server"] {
t.Errorf("expected other-server in result, got %v", servers)
}
}

func TestScanWorkflowsForDependencies_NoFrontmatter(t *testing.T) {
dir := t.TempDir()
content := "# Just a plain prompt\nNo frontmatter here."
if err := os.WriteFile(filepath.Join(dir, "plain.prompt.md"), []byte(content), 0o644); err != nil {
t.Fatal(err)
}
servers, err := aggregator.ScanWorkflowsForDependencies(dir)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(servers) != 0 {
t.Errorf("expected empty result for no-frontmatter file, got %v", servers)
}
}

func TestScanWorkflowsForDependencies_IgnoresNonPrompt(t *testing.T) {
dir := t.TempDir()
content := "---\nmcp:\n  - ignored-server\n---\n"
if err := os.WriteFile(filepath.Join(dir, "test.md"), []byte(content), 0o644); err != nil {
t.Fatal(err)
}
servers, err := aggregator.ScanWorkflowsForDependencies(dir)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(servers) != 0 {
t.Errorf("expected empty result for .md file (not .prompt.md), got %v", servers)
}
}
