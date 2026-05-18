package gemini_test

import (
"encoding/json"
"os"
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/adapters/client/gemini"
)

func TestUpdateConfig_WithGeminiDir(t *testing.T) {
dir := t.TempDir()
geminiDir := filepath.Join(dir, ".gemini")
if err := os.MkdirAll(geminiDir, 0o755); err != nil {
t.Fatal(err)
}
a := gemini.New(dir, false)
updates := map[string]interface{}{
"mcpServers": map[string]interface{}{
"my-server": map[string]interface{}{"command": "go", "args": []string{"run", "."}},
},
}
if err := a.UpdateConfig(updates); err != nil {
t.Fatalf("UpdateConfig unexpected error: %v", err)
}
data, err := os.ReadFile(filepath.Join(geminiDir, "settings.json"))
if err != nil {
t.Fatalf("settings.json not created: %v", err)
}
var cfg map[string]interface{}
if err := json.Unmarshal(data, &cfg); err != nil {
t.Fatalf("invalid JSON: %v", err)
}
if _, ok := cfg["mcpServers"]; !ok {
t.Error("settings.json should contain mcpServers key")
}
}

func TestGetCurrentConfig_ValidJSON(t *testing.T) {
dir := t.TempDir()
geminiDir := filepath.Join(dir, ".gemini")
if err := os.MkdirAll(geminiDir, 0o755); err != nil {
t.Fatal(err)
}
content := `{"mcpServers":{"s1":{"command":"node"}}}`
if err := os.WriteFile(filepath.Join(geminiDir, "settings.json"), []byte(content), 0o644); err != nil {
t.Fatal(err)
}
a := gemini.New(dir, false)
cfg := a.GetCurrentConfig()
if _, ok := cfg["mcpServers"]; !ok {
t.Error("GetCurrentConfig should return mcpServers")
}
}

func TestGetCurrentConfig_InvalidJSON(t *testing.T) {
dir := t.TempDir()
geminiDir := filepath.Join(dir, ".gemini")
if err := os.MkdirAll(geminiDir, 0o755); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(filepath.Join(geminiDir, "settings.json"), []byte("not json"), 0o644); err != nil {
t.Fatal(err)
}
a := gemini.New(dir, false)
cfg := a.GetCurrentConfig()
// Should return empty map, not panic.
if cfg == nil {
t.Error("GetCurrentConfig should return empty map on invalid JSON, not nil")
}
}

func TestGetConfigPath_UserScope(t *testing.T) {
dir := t.TempDir()
a := gemini.New(dir, true)
got := a.GetConfigPath()
// Even in user scope the path ends with settings.json.
if filepath.Base(got) != "settings.json" {
t.Errorf("GetConfigPath (user scope) should end with settings.json, got %q", got)
}
}

func TestNew_ReturnNonNil(t *testing.T) {
a := gemini.New("/tmp", false)
if a == nil {
t.Error("New should return non-nil adapter")
}
}

func TestTargetName_IsGemini(t *testing.T) {
for _, root := range []string{"/tmp", "", t.TempDir()} {
a := gemini.New(root, false)
if got := a.TargetName(); got != "gemini" {
t.Errorf("TargetName(%q): got %q, want gemini", root, got)
}
}
}

func TestMCPServersKey_IsConstant(t *testing.T) {
a := gemini.New(t.TempDir(), false)
k1 := a.MCPServersKey()
k2 := a.MCPServersKey()
if k1 != k2 || k1 != "mcpServers" {
t.Errorf("MCPServersKey not stable: %q / %q", k1, k2)
}
}

func TestUpdateConfig_EmptyRoot(t *testing.T) {
a := gemini.New("", false)
err := a.UpdateConfig(map[string]interface{}{})
// Should not panic; may return nil (no-op) since .gemini/ won't exist in cwd.
_ = err
}
