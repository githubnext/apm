package localbundle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMCPServerSpec_ZeroValue(t *testing.T) {
	s := MCPServerSpec{}
	if s.Name != "" || s.Command != "" || s.URL != "" {
		t.Errorf("unexpected non-empty fields: %+v", s)
	}
	if s.Registry {
		t.Error("expected Registry=false by default")
	}
}

func TestMCPServerSpec_AllFields(t *testing.T) {
	s := MCPServerSpec{
		Name:      "my-server",
		Transport: "stdio",
		Command:   "node",
		Args:      []string{"server.js", "--port", "3000"},
		Env:       map[string]string{"NODE_ENV": "production"},
		URL:       "https://mcp.example.com",
		Registry:  true,
	}
	if s.Name != "my-server" {
		t.Errorf("Name = %q", s.Name)
	}
	if s.Transport != "stdio" {
		t.Errorf("Transport = %q", s.Transport)
	}
	if len(s.Args) != 3 {
		t.Errorf("expected 3 args, got %d", len(s.Args))
	}
	if s.Env["NODE_ENV"] != "production" {
		t.Errorf("Env[NODE_ENV] = %q", s.Env["NODE_ENV"])
	}
	if !s.Registry {
		t.Error("expected Registry=true")
	}
}

func TestBundleMCPPresent_NoFile(t *testing.T) {
	dir := t.TempDir()
	if BundleMCPPresent(dir) {
		t.Error("expected BundleMCPPresent=false when no .mcp.json")
	}
}

func TestBundleMCPPresent_WithFile(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"srv": map[string]interface{}{"command": "node", "args": []string{"s.js"}},
		},
	}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	if !BundleMCPPresent(dir) {
		t.Error("expected BundleMCPPresent=true when .mcp.json exists")
	}
}

func TestParseBundleMCPServers_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 0 {
		t.Errorf("expected 0 servers for empty dir, got %d", len(servers))
	}
}

func TestParseBundleMCPServers_MultipleServers_extra2(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"alpha": map[string]interface{}{"command": "alpha-bin", "args": []string{}},
			"beta":  map[string]interface{}{"command": "beta-bin", "args": []string{}},
			"gamma": map[string]interface{}{"url": "https://gamma.example.com/mcp"},
		},
	}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 3 {
		t.Errorf("expected 3 servers, got %d", len(servers))
	}
	// Verify all have names.
	for _, s := range servers {
		if s.Name == "" {
			t.Error("expected non-empty Name for each server")
		}
	}
}

func TestParseBundleMCPServers_WithEnv_extra2(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"env-srv": map[string]interface{}{
				"command": "run",
				"args":    []string{},
				"env":     map[string]interface{}{"KEY": "VAL"},
			},
		},
	}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}
	if servers[0].Env["KEY"] != "VAL" {
		t.Errorf("Env[KEY] = %q", servers[0].Env["KEY"])
	}
}
