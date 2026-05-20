package pluginparser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParsePluginManifest_NameField(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{"name": "myplugin"}
	b, _ := json.Marshal(data)
	path := filepath.Join(dir, "plugin.json")
	_ = os.WriteFile(path, b, 0o600)
	m, err := ParsePluginManifest(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "myplugin" {
		t.Errorf("expected myplugin, got %q", m.Name)
	}
}

func TestParsePluginManifest_AgentsField(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{"name": "p", "agents": []string{"agent1", "agent2"}}
	b, _ := json.Marshal(data)
	path := filepath.Join(dir, "plugin.json")
	_ = os.WriteFile(path, b, 0o600)
	m, err := ParsePluginManifest(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(m.Agents))
	}
}

func TestParsePluginManifest_SkillsField(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{"name": "p", "skills": []string{"skillA"}}
	b, _ := json.Marshal(data)
	path := filepath.Join(dir, "plugin.json")
	_ = os.WriteFile(path, b, 0o600)
	m, err := ParsePluginManifest(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Skills) != 1 || m.Skills[0] != "skillA" {
		t.Errorf("unexpected skills: %v", m.Skills)
	}
}

func TestParsePluginManifest_CommandsField(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{"name": "p", "commands": []string{"cmd1", "cmd2", "cmd3"}}
	b, _ := json.Marshal(data)
	path := filepath.Join(dir, "plugin.json")
	_ = os.WriteFile(path, b, 0o600)
	m, err := ParsePluginManifest(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Commands) != 3 {
		t.Errorf("expected 3 commands, got %d", len(m.Commands))
	}
}

func TestPluginManifest_ZeroValue(t *testing.T) {
	var m PluginManifest
	if m.Name != "" {
		t.Error("expected empty Name")
	}
	if len(m.Agents) != 0 {
		t.Error("expected empty Agents")
	}
}

func TestMCPServerConfig_ZeroValue(t *testing.T) {
	var c MCPServerConfig
	if c.Command != "" {
		t.Error("expected empty Command")
	}
	if len(c.Args) != 0 {
		t.Error("expected empty Args")
	}
}

func TestMCPDepEntry_ZeroValue(t *testing.T) {
	var e MCPDepEntry
	if e.Name != "" {
		t.Error("expected empty Name")
	}
	if e.Registry {
		t.Error("expected Registry=false")
	}
}

func TestMCPServerConfig_WithEnv(t *testing.T) {
	c := MCPServerConfig{
		Command: "node",
		Env:     map[string]string{"KEY": "VAL"},
	}
	if c.Env["KEY"] != "VAL" {
		t.Error("expected VAL")
	}
}

func TestMCPDepEntry_WithHeaders(t *testing.T) {
	e := MCPDepEntry{
		Name:    "srv",
		Headers: map[string]string{"Authorization": "Bearer tok"},
	}
	if e.Headers["Authorization"] != "Bearer tok" {
		t.Errorf("unexpected header: %v", e.Headers)
	}
}

func TestParsePluginManifest_NonexistentFile(t *testing.T) {
	_, err := ParsePluginManifest("/nonexistent/path/plugin.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
