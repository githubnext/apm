package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeMCPEntryForClaudeCode_RemoteType(t *testing.T) {
	entry := map[string]interface{}{
		"type":  "http",
		"url":   "https://example.com/mcp",
		"tools": []string{"tool1"},
		"id":    "some-id",
	}
	out := normalizeMCPEntryForClaudeCode(entry)
	if _, ok := out["tools"]; ok {
		t.Error("tools should be removed for remote type")
	}
	if _, ok := out["id"]; ok {
		t.Error("id should be removed for remote type")
	}
	if out["url"] != "https://example.com/mcp" {
		t.Errorf("url should be preserved: %v", out["url"])
	}
}

func TestNormalizeMCPEntryForClaudeCode_StdioType(t *testing.T) {
	entry := map[string]interface{}{
		"command": "npx",
		"args":    []string{"-y", "my-server"},
		"tools":   []string{"tool1"},
		"id":      "",
	}
	out := normalizeMCPEntryForClaudeCode(entry)
	if _, ok := out["tools"]; ok {
		t.Error("tools should be removed")
	}
	if out["type"] != "stdio" {
		t.Errorf("type should be stdio for command entries: %v", out["type"])
	}
}

func TestNormalizeMCPEntryForClaudeCode_NoCommand(t *testing.T) {
	entry := map[string]interface{}{
		"args": []string{"--flag"},
	}
	out := normalizeMCPEntryForClaudeCode(entry)
	if _, ok := out["type"]; ok {
		t.Error("type should not be set when no command present")
	}
}

func TestNormalizeMCPEntryForClaudeCode_PreservesNonTools(t *testing.T) {
	entry := map[string]interface{}{
		"command": "node",
		"env":     map[string]string{"KEY": "val"},
	}
	out := normalizeMCPEntryForClaudeCode(entry)
	if out["env"] == nil {
		t.Error("env should be preserved")
	}
	if out["command"] != "node" {
		t.Errorf("command should be preserved: %v", out["command"])
	}
}

func TestServerKeyFor_NameTakesPriority(t *testing.T) {
	k := serverKeyFor("github.com/org/repo", "my-name")
	if k != "my-name" {
		t.Errorf("expected my-name, got %q", k)
	}
}

func TestServerKeyFor_EmptyName(t *testing.T) {
	k := serverKeyFor("github.com/org/repo", "")
	if k != "repo" {
		t.Errorf("expected repo, got %q", k)
	}
}

func TestServerKeyFor_NoSlash(t *testing.T) {
	k := serverKeyFor("my-server", "")
	if k != "my-server" {
		t.Errorf("expected my-server, got %q", k)
	}
}

func TestAdapterGetConfigPath_ContainsClaude(t *testing.T) {
	a := New("/project", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, "mcp") {
		t.Errorf("expected mcp in path: %q", p)
	}
	if !strings.Contains(p, "project") {
		t.Errorf("expected project root in path: %q", p)
	}
}

func TestAdapterUserScope_DifferentPath(t *testing.T) {
	ap := New("/project", false)
	au := New("/project", true)
	pp := ap.GetConfigPath()
	up := au.GetConfigPath()
	if pp == up {
		t.Error("project scope and user scope should have different paths")
	}
}

func TestUpdateConfig_WritesJSON(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{
		"test-server": map[string]interface{}{"command": "node"},
	})
	if err != nil {
		t.Fatalf("UpdateConfig error: %v", err)
	}
	data, err := os.ReadFile(a.GetConfigPath())
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
}

func TestGetCurrentConfig_WithValidFile(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	cfgPath := a.GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"srv": map[string]interface{}{"command": "npx"},
		},
	}
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	got := a.GetCurrentConfig()
	if got["mcpServers"] == nil {
		t.Error("mcpServers key should be present")
	}
}

func TestNormalizeMCPEntryForClaudeCode_RemoteTypePreservesType(t *testing.T) {
	entry := map[string]interface{}{
		"type": "remote",
		"url":  "https://mcp.example.com",
	}
	out := normalizeMCPEntryForClaudeCode(entry)
	if out["type"] != "remote" {
		t.Errorf("remote type should be preserved: %v", out["type"])
	}
}

func TestNormalizeMCPEntryForClaudeCode_NonEmptyID_Preserved(t *testing.T) {
	entry := map[string]interface{}{
		"command": "node",
		"id":      "my-id",
	}
	out := normalizeMCPEntryForClaudeCode(entry)
	if out["id"] != "my-id" {
		t.Errorf("non-empty id should be preserved: %v", out["id"])
	}
}
