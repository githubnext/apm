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

func TestScanWorkflowsForDependencies_SingleMCP(t *testing.T) {
dir := t.TempDir()
content := "---\nmcp:\n  - only-server\n---\n"
if err := os.WriteFile(filepath.Join(dir, "single.prompt.md"), []byte(content), 0o644); err != nil {
t.Fatal(err)
}
servers, err := aggregator.ScanWorkflowsForDependencies(dir)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !servers["only-server"] {
t.Errorf("expected only-server in %v", servers)
}
if len(servers) != 1 {
t.Errorf("expected 1 server, got %d: %v", len(servers), servers)
}
}

func TestScanWorkflowsForDependencies_DeduplicatesAcrossFiles(t *testing.T) {
dir := t.TempDir()
c1 := "---\nmcp:\n  - shared-server\n  - file1-only\n---\n"
c2 := "---\nmcp:\n  - shared-server\n  - file2-only\n---\n"
if err := os.WriteFile(filepath.Join(dir, "a.prompt.md"), []byte(c1), 0o644); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(filepath.Join(dir, "b.prompt.md"), []byte(c2), 0o644); err != nil {
t.Fatal(err)
}
servers, err := aggregator.ScanWorkflowsForDependencies(dir)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !servers["shared-server"] {
t.Errorf("expected shared-server, got %v", servers)
}
if !servers["file1-only"] {
t.Errorf("expected file1-only, got %v", servers)
}
if !servers["file2-only"] {
t.Errorf("expected file2-only, got %v", servers)
}
if len(servers) != 3 {
t.Errorf("expected 3 unique servers, got %d: %v", len(servers), servers)
}
}

func TestScanWorkflowsForDependencies_Recursive(t *testing.T) {
dir := t.TempDir()
sub := filepath.Join(dir, "subdir")
if err := os.MkdirAll(sub, 0o755); err != nil {
t.Fatal(err)
}
content := "---\nmcp:\n  - nested-server\n---\n"
if err := os.WriteFile(filepath.Join(sub, "nested.prompt.md"), []byte(content), 0o644); err != nil {
t.Fatal(err)
}
servers, err := aggregator.ScanWorkflowsForDependencies(dir)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !servers["nested-server"] {
t.Errorf("expected nested-server, got %v", servers)
}
}
