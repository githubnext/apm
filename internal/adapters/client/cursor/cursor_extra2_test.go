package cursor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeMCPEntryForCursor_PreservesCommand(t *testing.T) {
	entry := map[string]interface{}{
		"command": "npx",
		"args":    []interface{}{"-y", "my-server"},
	}
	out := normalizeMCPEntryForCursor(entry)
	if out["command"] != "npx" {
		t.Errorf("command should be preserved: %v", out["command"])
	}
}

func TestNormalizeMCPEntryForCursor_PreservesEnv(t *testing.T) {
	entry := map[string]interface{}{
		"command": "node",
		"env": map[string]interface{}{
			"TOKEN": "abc",
		},
	}
	out := normalizeMCPEntryForCursor(entry)
	if out["env"] == nil {
		t.Error("env should be preserved")
	}
}

func TestNormalizeMCPEntryForCursor_EmptyEntry(t *testing.T) {
	out := normalizeMCPEntryForCursor(map[string]interface{}{})
	if out == nil {
		t.Error("output should not be nil for empty entry")
	}
}

func TestServerKeyFor_WithName(t *testing.T) {
	k := serverKeyFor("github.com/org/repo", "custom-name")
	if k != "custom-name" {
		t.Errorf("expected custom-name, got %q", k)
	}
}

func TestServerKeyFor_EmptyName_ExtractsFromURL(t *testing.T) {
	k := serverKeyFor("github.com/org/my-repo", "")
	if k != "my-repo" {
		t.Errorf("expected my-repo, got %q", k)
	}
}

func TestServerKeyFor_NoSlashInURL(t *testing.T) {
	k := serverKeyFor("my-server", "")
	if k != "my-server" {
		t.Errorf("expected my-server, got %q", k)
	}
}

func TestGetConfigPath_ContainsCursor(t *testing.T) {
	a := New("/my/project", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, ".cursor") {
		t.Errorf("expected .cursor in path: %q", p)
	}
	if !strings.HasSuffix(p, "mcp.json") {
		t.Errorf("expected mcp.json suffix: %q", p)
	}
}

func TestUpdateConfig_CreatesDirectoryAndFile(t *testing.T) {
	dir := t.TempDir()
	// cursor UpdateConfig only writes if .cursor/ dir exists
	cursorDir := filepath.Join(dir, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{
		"my-server": map[string]interface{}{"command": "npx"},
	})
	if err != nil {
		t.Fatalf("UpdateConfig error: %v", err)
	}
	cfgPath := a.GetConfigPath()
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("config file should exist: %v", err)
	}
}

func TestUpdateConfig_JSONIsValid(t *testing.T) {
	dir := t.TempDir()
	cursorDir := filepath.Join(dir, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	a := New(dir, false)
	if err := a.UpdateConfig(map[string]interface{}{
		"srv": map[string]interface{}{"command": "node"},
	}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(a.GetConfigPath())
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("config is not valid JSON: %v", err)
	}
}

func TestGetCurrentConfig_InvalidJSONExtra(t *testing.T) {
	dir := t.TempDir()
	cursorDir := filepath.Join(dir, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cursorDir, "mcp.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if len(cfg) != 0 {
		t.Errorf("expected empty map for invalid JSON: %v", cfg)
	}
}

func TestNew_ProjectRootSet(t *testing.T) {
	a := New("/tmp/proj", false)
	if a == nil {
		t.Fatal("New should return non-nil")
	}
}

func TestGetCurrentConfig_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	cursorDir := filepath.Join(dir, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"my-server": map[string]interface{}{"command": "node"},
		},
	}
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(filepath.Join(cursorDir, "mcp.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	a := New(dir, false)
	got := a.GetCurrentConfig()
	if got["mcpServers"] == nil {
		t.Error("mcpServers key should be present")
	}
}
