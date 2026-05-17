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

func TestParseBundleMCPServers_MultipleServers(t *testing.T) {
dir := t.TempDir()
data := map[string]interface{}{
"mcpServers": map[string]interface{}{
"server-a": map[string]interface{}{"command": "cmd-a", "type": "stdio"},
"server-b": map[string]interface{}{"url": "http://localhost:8080", "type": "sse"},
},
}
b, _ := json.Marshal(data)
os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
servers := ParseBundleMCPServers(dir)
if len(servers) != 2 {
t.Fatalf("expected 2 servers, got %d", len(servers))
}
}

func TestParseBundleMCPServers_WithEnv(t *testing.T) {
dir := t.TempDir()
data := map[string]interface{}{
"mcpServers": map[string]interface{}{
"srv": map[string]interface{}{
"command": "node",
"args":    []interface{}{"index.js"},
"env":     map[string]interface{}{"KEY": "val"},
},
},
}
b, _ := json.Marshal(data)
os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
servers := ParseBundleMCPServers(dir)
if len(servers) != 1 {
t.Fatalf("expected 1, got %d", len(servers))
}
if servers[0].Env["KEY"] != "val" {
t.Errorf("expected env KEY=val, got %v", servers[0].Env)
}
}

func TestParseBundleMCPServers_MalformedJSON(t *testing.T) {
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte("not json"), 0644)
servers := ParseBundleMCPServers(dir)
if len(servers) != 0 {
t.Error("expected empty on malformed JSON")
}
}

func TestParseBundleMCPServers_NoMCPServersKey(t *testing.T) {
dir := t.TempDir()
os.WriteFile(filepath.Join(dir, ".mcp.json"), []byte(`{"other":"value"}`), 0644)
servers := ParseBundleMCPServers(dir)
if len(servers) != 0 {
t.Error("expected empty when mcpServers key missing")
}
}

func TestParseBundleMCPServers_SSETransport(t *testing.T) {
dir := t.TempDir()
data := map[string]interface{}{
"mcpServers": map[string]interface{}{
"remote": map[string]interface{}{
"url":       "https://example.com/mcp",
"transport": "sse",
},
},
}
b, _ := json.Marshal(data)
os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
servers := ParseBundleMCPServers(dir)
if len(servers) != 1 {
t.Fatalf("expected 1, got %d", len(servers))
}
if servers[0].URL != "https://example.com/mcp" {
t.Errorf("URL mismatch: %s", servers[0].URL)
}
if servers[0].Transport != "sse" {
t.Errorf("expected sse transport, got %s", servers[0].Transport)
}
}
