package claude_test

import (
"os"
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/adapters/client/claude"
)

func TestTargetName(t *testing.T) {
a := claude.New("/tmp", false)
if got := a.TargetName(); got != "claude" {
t.Errorf("TargetName: want claude, got %s", got)
}
}

func TestMCPServersKey(t *testing.T) {
a := claude.New("/tmp", false)
if got := a.MCPServersKey(); got != "mcpServers" {
t.Errorf("MCPServersKey: want mcpServers, got %s", got)
}
}

func TestSupportsUserScope(t *testing.T) {
a := claude.New("/tmp", false)
if !a.SupportsUserScope() {
t.Error("SupportsUserScope should return true")
}
}

func TestGetConfigPathProjectScope(t *testing.T) {
dir := t.TempDir()
a := claude.New(dir, false)
got := a.GetConfigPath()
want := filepath.Join(dir, ".mcp.json")
if got != want {
t.Errorf("GetConfigPath: want %s, got %s", want, got)
}
}

func TestGetConfigPathUserScope(t *testing.T) {
a := claude.New("", true)
got := a.GetConfigPath()
home, _ := os.UserHomeDir()
want := filepath.Join(home, ".claude.json")
if got != want {
t.Errorf("GetConfigPath user scope: want %s, got %s", want, got)
}
}

func TestGetCurrentConfigMissing(t *testing.T) {
a := claude.New(t.TempDir(), false)
cfg := a.GetCurrentConfig()
if cfg == nil {
t.Error("GetCurrentConfig should return empty map, not nil")
}
if len(cfg) != 0 {
t.Errorf("GetCurrentConfig on missing file: want empty map, got %v", cfg)
}
}

func TestUpdateConfigRoundtrip(t *testing.T) {
dir := t.TempDir()
a := claude.New(dir, false)
err := a.UpdateConfig(map[string]interface{}{
"my-server": map[string]interface{}{"command": "node", "args": []string{"server.js"}},
})
if err != nil {
t.Fatalf("UpdateConfig: %v", err)
}
cfg := a.GetCurrentConfig()
servers, ok := cfg["mcpServers"].(map[string]interface{})
if !ok {
t.Fatalf("mcpServers not a map: %T", cfg["mcpServers"])
}
if _, ok := servers["my-server"]; !ok {
t.Error("my-server not found in config")
}
}
