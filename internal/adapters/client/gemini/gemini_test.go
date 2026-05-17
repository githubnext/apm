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

func TestGetConfigPathEmptyRoot(t *testing.T) {
a := gemini.New("", false)
got := a.GetConfigPath()
if got == "" {
t.Error("GetConfigPath with empty root should not return empty string")
}
if filepath.Base(got) != "settings.json" {
t.Errorf("GetConfigPath should end with settings.json, got %q", got)
}
}

func TestUpdateConfigNoGeminiDir(t *testing.T) {
dir := t.TempDir()
a := gemini.New(dir, false)
err := a.UpdateConfig(map[string]interface{}{"mcpServers": map[string]interface{}{}})
if err != nil {
t.Errorf("UpdateConfig with no .gemini dir should be a no-op, got: %v", err)
}
}

func TestTargetNameIsStable(t *testing.T) {
a1 := gemini.New("/tmp/a", false)
a2 := gemini.New("/tmp/b", true)
if a1.TargetName() != a2.TargetName() {
t.Error("TargetName should not depend on constructor args")
}
}

func TestMCPServersKeyIsStable(t *testing.T) {
a := gemini.New("/tmp", true)
if a.MCPServersKey() != "mcpServers" {
t.Errorf("MCPServersKey: want mcpServers, got %s", a.MCPServersKey())
}
}

func TestGetConfigPathContainsGemini(t *testing.T) {
dir := t.TempDir()
a := gemini.New(dir, false)
got := a.GetConfigPath()
if !filepath.IsAbs(got) {
t.Errorf("GetConfigPath should be absolute, got %q", got)
}
if filepath.Dir(filepath.Dir(got)) != dir {
t.Errorf("expected path under %s/.gemini/, got %q", dir, got)
}
}
