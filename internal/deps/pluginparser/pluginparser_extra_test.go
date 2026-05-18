package pluginparser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParsePluginManifest_WithAgentsAndSkills(t *testing.T) {
	dir := t.TempDir()
	pluginJSON := filepath.Join(dir, "plugin.json")
	data := map[string]interface{}{
		"name":   "full-plugin",
		"agents": []string{"agent1.md", "agent2.md"},
		"skills": []string{"skill1.md"},
	}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(pluginJSON, b, 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := ParsePluginManifest(pluginJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(m.Agents))
	}
	if len(m.Skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(m.Skills))
	}
}

func TestParsePluginManifest_WithCommands(t *testing.T) {
	dir := t.TempDir()
	pluginJSON := filepath.Join(dir, "plugin.json")
	data := map[string]interface{}{
		"name":     "cmd-plugin",
		"commands": []string{"cmd1.md", "cmd2.md", "cmd3.md"},
	}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(pluginJSON, b, 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := ParsePluginManifest(pluginJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Commands) != 3 {
		t.Errorf("expected 3 commands, got %d", len(m.Commands))
	}
}

func TestParsePluginManifest_EmptyName(t *testing.T) {
	dir := t.TempDir()
	pluginJSON := filepath.Join(dir, "plugin.json")
	data := map[string]interface{}{"agents": []string{"a.md"}}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(pluginJSON, b, 0o644); err != nil {
		t.Fatal(err)
	}
	// Should not error even if name is empty; logs a warning
	m, err := ParsePluginManifest(pluginJSON)
	if err != nil {
		t.Fatalf("unexpected error for empty name: %v", err)
	}
	if m.Name != "" {
		t.Errorf("expected empty name, got %q", m.Name)
	}
}

func TestYamlString_NoQuoting(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"with space", `"with space"`},
		{"with:colon", `"with:colon"`},
		{"with#hash", `"with#hash"`},
		{"", ``},
	}
	for _, c := range cases {
		got := yamlString(c.input)
		if got != c.want {
			t.Errorf("yamlString(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestMCPServerConfig_Fields(t *testing.T) {
	cfg := MCPServerConfig{
		Command: "npx",
		Args:    []string{"-y", "some-pkg"},
		URL:     "http://localhost",
		Type:    "local",
	}
	if cfg.Command != "npx" {
		t.Errorf("Command mismatch: %q", cfg.Command)
	}
	if len(cfg.Args) != 2 {
		t.Errorf("Args length mismatch: %d", len(cfg.Args))
	}
}

func TestMCPDepEntry_Fields(t *testing.T) {
	dep := MCPDepEntry{
		Name:      "my-server",
		Transport: "stdio",
		Command:   "node",
		Args:      []string{"server.js"},
	}
	if dep.Name != "my-server" {
		t.Errorf("Name mismatch: %q", dep.Name)
	}
	if dep.Transport != "stdio" {
		t.Errorf("Transport mismatch: %q", dep.Transport)
	}
}

func TestPluginManifest_StructFields(t *testing.T) {
	m := PluginManifest{
		Name:   "test",
		Agents: []string{"a.md"},
	}
	if m.Name != "test" {
		t.Errorf("Name mismatch")
	}
	if len(m.Agents) != 1 {
		t.Errorf("Agents length: %d", len(m.Agents))
	}
}
