package localbundle

import (
"encoding/json"
"os"
"path/filepath"
"testing"
)

func TestParseBundleMCPServers(t *testing.T) {
dir := t.TempDir()
data := map[string]interface{}{
"mcpServers": map[string]interface{}{
"my-server": map[string]interface{}{
"command": "npx",
"args":    []interface{}{"-y", "my-pkg"},
"type":    "stdio",
},
},
}
b, _ := json.Marshal(data)
if err := os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644); err != nil {
t.Fatal(err)
}
servers := ParseBundleMCPServers(dir)
if len(servers) != 1 {
t.Fatalf("expected 1 server, got %d", len(servers))
}
s := servers[0]
if s.Name != "my-server" {
t.Errorf("expected my-server, got %s", s.Name)
}
if s.Command != "npx" {
t.Errorf("expected npx, got %s", s.Command)
}
if s.Transport != "stdio" {
t.Errorf("expected stdio, got %s", s.Transport)
}
}

func TestParseBundleMCPServersMissing(t *testing.T) {
dir := t.TempDir()
servers := ParseBundleMCPServers(dir)
if len(servers) != 0 {
t.Errorf("expected no servers, got %d", len(servers))
}
}

func TestBundleMCPPresent(t *testing.T) {
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("{}"), 0644)
if !BundleMCPPresent(dir) {
t.Error("expected true")
}
}

func TestBundleMCPPresentFalse(t *testing.T) {
dir := t.TempDir()
if BundleMCPPresent(dir) {
t.Error("expected false")
}
}
