package gemini_test

import (
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/adapters/client/gemini"
)

func TestTargetName(t *testing.T) {
a := gemini.New("/tmp", false)
if got := a.TargetName(); got != "gemini" {
t.Errorf("TargetName: want gemini, got %s", got)
}
}

func TestMCPServersKey(t *testing.T) {
a := gemini.New("/tmp", false)
if got := a.MCPServersKey(); got != "mcpServers" {
t.Errorf("MCPServersKey: want mcpServers, got %s", got)
}
}

func TestSupportsUserScope(t *testing.T) {
a := gemini.New("/tmp", false)
if !a.SupportsUserScope() {
t.Error("SupportsUserScope should return true")
}
}

func TestGetConfigPath(t *testing.T) {
dir := t.TempDir()
a := gemini.New(dir, false)
got := a.GetConfigPath()
want := filepath.Join(dir, ".gemini", "settings.json")
if got != want {
t.Errorf("GetConfigPath: want %s, got %s", want, got)
}
}

func TestGetCurrentConfigMissing(t *testing.T) {
a := gemini.New(t.TempDir(), false)
cfg := a.GetCurrentConfig()
if cfg == nil {
t.Error("GetCurrentConfig should return empty map, not nil")
}
if len(cfg) != 0 {
t.Errorf("GetCurrentConfig on missing file: want empty, got %v", cfg)
}
}
