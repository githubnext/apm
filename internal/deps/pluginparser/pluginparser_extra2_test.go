package pluginparser

import (
	"strings"
	"testing"
)

func TestYamlString_WithColon(t *testing.T) {
	s := yamlString("key: value")
	if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
		t.Errorf("expected quoted string for colon-containing value, got %q", s)
	}
}

func TestYamlString_WithBraces(t *testing.T) {
	s := yamlString("{key}")
	if !strings.HasPrefix(s, `"`) {
		t.Errorf("expected quoted string for brace-containing value, got %q", s)
	}
}

func TestYamlString_WithNewline(t *testing.T) {
	s := yamlString("line1\nline2")
	if !strings.HasPrefix(s, `"`) {
		t.Errorf("expected quoted string for newline-containing value, got %q", s)
	}
}

func TestYamlString_WithSpace(t *testing.T) {
	s := yamlString("has space")
	if !strings.HasPrefix(s, `"`) {
		t.Errorf("expected quoted string for space-containing value, got %q", s)
	}
}

func TestYamlString_Plain(t *testing.T) {
	s := yamlString("simple-value")
	if s != "simple-value" {
		t.Errorf("expected unquoted plain value, got %q", s)
	}
}

func TestYamlString_WithInnerQuotes(t *testing.T) {
	input := `say "hello"`
	s := yamlString(input)
	if !strings.Contains(s, `\"`) {
		t.Errorf("expected escaped inner quotes, got %q", s)
	}
}

func TestMCPDepEntry_AllFields(t *testing.T) {
	entry := MCPDepEntry{
		Name:      "my-server",
		Transport: "stdio",
		Args:      []string{"--port", "9000"},
		Env:       map[string]string{"KEY": "VAL"},
	}
	if entry.Name != "my-server" {
		t.Errorf("expected Name=my-server, got %q", entry.Name)
	}
	if len(entry.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(entry.Args))
	}
	if entry.Env["KEY"] != "VAL" {
		t.Errorf("expected ENV KEY=VAL")
	}
}

func TestMCPServerConfig_AllFields(t *testing.T) {
	cfg := MCPServerConfig{
		Command: "node",
		Args:    []string{"index.js"},
		Env:     map[string]string{"PORT": "3000"},
		URL:     "https://example.com/mcp",
	}
	if cfg.Command != "node" {
		t.Errorf("expected Command=node, got %q", cfg.Command)
	}
	if cfg.URL != "https://example.com/mcp" {
		t.Errorf("expected URL, got %q", cfg.URL)
	}
}

func TestPluginManifest_MultipleAgents(t *testing.T) {
	pm := PluginManifest{
		Name:   "my-plugin",
		Agents: []string{"agent1", "agent2", "agent3"},
	}
	if len(pm.Agents) != 3 {
		t.Errorf("expected 3 agents, got %d", len(pm.Agents))
	}
}
