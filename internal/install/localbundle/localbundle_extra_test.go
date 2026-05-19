package localbundle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseBundleMCPServers_URL_SSE_NoCommand(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"remote-srv": map[string]interface{}{
				"url":  "https://mcp.example.com/endpoint",
				"type": "sse",
			},
		},
	}
	b, _ := json.Marshal(data)
	os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}
	s := servers[0]
	if s.URL != "https://mcp.example.com/endpoint" {
		t.Errorf("URL = %q", s.URL)
	}
	if s.Command != "" {
		t.Errorf("Command should be empty for SSE, got %q", s.Command)
	}
}

func TestParseBundleMCPServers_FallbackTransportField(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"srv": map[string]interface{}{
				"command":   "mybin",
				"transport": "stdio",
			},
		},
	}
	b, _ := json.Marshal(data)
	os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 1 {
		t.Fatalf("expected 1, got %d", len(servers))
	}
	// transport field (not type) should be used
	if servers[0].Transport != "stdio" {
		t.Errorf("Transport = %q, want stdio", servers[0].Transport)
	}
}

func TestParseBundleMCPServers_ArgsPreserved(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"srv": map[string]interface{}{
				"command": "node",
				"args":    []interface{}{"--flag", "value", "-x"},
				"type":    "stdio",
			},
		},
	}
	b, _ := json.Marshal(data)
	os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 1 {
		t.Fatalf("expected 1, got %d", len(servers))
	}
	if len(servers[0].Args) != 3 {
		t.Errorf("expected 3 args, got %v", servers[0].Args)
	}
	if servers[0].Args[0] != "--flag" || servers[0].Args[2] != "-x" {
		t.Errorf("args mismatch: %v", servers[0].Args)
	}
}

func TestParseBundleMCPServers_EmptyEnv(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"srv": map[string]interface{}{
				"command": "bin",
			},
		},
	}
	b, _ := json.Marshal(data)
	os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 1 {
		t.Fatalf("expected 1, got %d", len(servers))
	}
	if servers[0].Env != nil && len(servers[0].Env) != 0 {
		t.Errorf("env should be nil/empty, got %v", servers[0].Env)
	}
}

func TestParseBundleMCPServers_InvalidServersMap(t *testing.T) {
	dir := t.TempDir()
	// mcpServers is not a map -- array instead
	data := map[string]interface{}{
		"mcpServers": []interface{}{"srv1"},
	}
	b, _ := json.Marshal(data)
	os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 0 {
		t.Errorf("expected empty when mcpServers is not a map, got %d", len(servers))
	}
}

func TestBundleMCPPresent_CaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	// file named .MCP.JSON (uppercase) should be detected
	os.WriteFile(filepath.Join(dir, ".MCP.JSON"), []byte("{}"), 0644)
	if !BundleMCPPresent(dir) {
		t.Error("BundleMCPPresent should be case-insensitive")
	}
}

func TestParseBundleMCPServers_NonExistentDir(t *testing.T) {
	servers := ParseBundleMCPServers("/nonexistent/path/that/does/not/exist")
	if servers != nil && len(servers) != 0 {
		t.Errorf("expected nil/empty for non-existent dir, got %v", servers)
	}
}

func TestParseBundleMCPServers_RawMapContainsAllKeys(t *testing.T) {
	dir := t.TempDir()
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"srv": map[string]interface{}{
				"command": "cmd",
				"extra":   "value",
			},
		},
	}
	b, _ := json.Marshal(data)
	os.WriteFile(filepath.Join(dir, ".mcp.json"), b, 0644)
	servers := ParseBundleMCPServers(dir)
	if len(servers) != 1 {
		t.Fatalf("expected 1, got %d", len(servers))
	}
	if servers[0].Raw == nil {
		t.Error("Raw should be populated")
	}
	if _, ok := servers[0].Raw["extra"]; !ok {
		t.Error("Raw should contain extra field")
	}
}
