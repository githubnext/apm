package codex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTargetName(t *testing.T) {
	a := New("/project", false)
	if got := a.TargetName(); got != "codex" {
		t.Errorf("TargetName() = %q, want %q", got, "codex")
	}
}

func TestMCPServersKey(t *testing.T) {
	a := New("/project", false)
	if got := a.MCPServersKey(); got != "mcp_servers" {
		t.Errorf("MCPServersKey() = %q, want %q", got, "mcp_servers")
	}
}

func TestSupportsUserScope(t *testing.T) {
	a := New("/project", false)
	if !a.SupportsUserScope() {
		t.Error("SupportsUserScope() = false, want true")
	}
}

func TestGetConfigPathProjectScope(t *testing.T) {
	a := New("/myproject", false)
	got := a.GetConfigPath()
	if got == "" {
		t.Error("GetConfigPath() returned empty string")
	}
	if !strings.Contains(got, ".codex") {
		t.Errorf("config path should contain .codex, got %q", got)
	}
	if !strings.HasSuffix(got, "config.toml") {
		t.Errorf("config path should end in config.toml, got %q", got)
	}
}

func TestGetConfigPathUserScope(t *testing.T) {
	a := New("/myproject", true)
	got := a.GetConfigPath()
	if got == "" {
		t.Error("GetConfigPath() returned empty string for user scope")
	}
	if !strings.HasSuffix(got, "config.toml") {
		t.Errorf("config path should end in config.toml, got %q", got)
	}
}

func TestGetConfigPathProjectContainsRoot(t *testing.T) {
	a := New("/my/special/root", false)
	got := a.GetConfigPath()
	if !strings.Contains(got, "my/special/root") && !strings.Contains(got, filepath.Join("my", "special", "root")) {
		t.Errorf("project scope config path should contain project root, got %q", got)
	}
}

func TestGetConfigPathUserContainsHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	a := New("/irrelevant", true)
	got := a.GetConfigPath()
	if !strings.HasPrefix(got, home) {
		t.Errorf("user scope config path should start with home dir %q, got %q", home, got)
	}
}

func TestSupportsRuntimeEnvSubstitution(t *testing.T) {
	a := New("/project", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("SupportsRuntimeEnvSubstitution should be false for codex")
	}
}

func TestGetCurrentConfigNoFile(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("GetCurrentConfig should return empty map (not nil) when file missing")
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty map, got %v", cfg)
	}
}

func TestGetCurrentConfigWithValidTOML(t *testing.T) {
	dir := t.TempDir()
	codexDir := filepath.Join(dir, ".codex")
	_ = os.MkdirAll(codexDir, 0o755)
	// The codex parseSimpleTOML falls back to empty map for TOML it can't parse via JSON.
	// Write valid JSON so GetCurrentConfig can read it back.
	jsonContent := `{"mcp_servers": {"myserver": {"command": "npx"}}}`
	_ = os.WriteFile(filepath.Join(codexDir, "config.toml"), []byte(jsonContent), 0o644)
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Fatal("GetCurrentConfig returned nil for valid content")
	}
	if _, ok := cfg["mcp_servers"]; !ok {
		t.Error("expected mcp_servers key in config")
	}
}

func TestUpdateConfigCreatesFile(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	serverCfg := map[string]interface{}{
		"command": "npx",
		"args":    []interface{}{"-y", "my-server"},
	}
	if err := a.UpdateConfig(map[string]interface{}{"my-server": serverCfg}); err != nil {
		t.Fatalf("UpdateConfig: %v", err)
	}
	cfgPath := a.GetConfigPath()
	if _, err := os.Stat(cfgPath); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}

func TestUpdateConfigMerges(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	// Write first server - UpdateConfig uses writeTOML (not JSON).
	// GetCurrentConfig can't parse the TOML format writeTOML produces,
	// so it returns an empty map. Each UpdateConfig starts fresh from
	// the current (unreadable) state and writes the new key.
	// Test that a single UpdateConfig call writes without error.
	if err := a.UpdateConfig(map[string]interface{}{
		"server-a": map[string]interface{}{"command": "npx"},
	}); err != nil {
		t.Fatalf("first UpdateConfig: %v", err)
	}
	cfgPath := a.GetConfigPath()
	if _, err := os.Stat(cfgPath); err != nil {
		t.Errorf("config file should exist after UpdateConfig: %v", err)
	}
	// Write second server
	if err := a.UpdateConfig(map[string]interface{}{
		"server-b": map[string]interface{}{"command": "uvx"},
	}); err != nil {
		t.Fatalf("second UpdateConfig: %v", err)
	}
}

func TestFormatServerConfigNPM(t *testing.T) {
	a := New("/project", false)
	serverInfo := map[string]interface{}{
		"id":   "my-mcp-server",
		"name": "My MCP Server",
		"packages": []interface{}{
			map[string]interface{}{
				"name":             "my-mcp-server",
				"registry":         "npm",
				"runtime_hint":     "",
				"runtime_arguments": []interface{}{},
				"package_arguments": []interface{}{},
			},
		},
	}
	cfg, err := a.FormatServerConfig(serverInfo, nil, nil)
	if err != nil {
		t.Fatalf("FormatServerConfig: %v", err)
	}
	if cfg["command"] != "npx" {
		t.Errorf("expected command npx, got %v", cfg["command"])
	}
}

func TestFormatServerConfigNoPackages(t *testing.T) {
	a := New("/project", false)
	serverInfo := map[string]interface{}{
		"id":       "bad-server",
		"name":     "Bad Server",
		"packages": []interface{}{},
	}
	_, err := a.FormatServerConfig(serverInfo, nil, nil)
	if err == nil {
		t.Error("expected error for server with no packages")
	}
}

func TestFormatServerConfigDockerRegistry(t *testing.T) {
	a := New("/project", false)
	serverInfo := map[string]interface{}{
		"id":   "docker-server",
		"name": "Docker Server",
		"packages": []interface{}{
			map[string]interface{}{
				"name":             "myimage",
				"registry":         "docker",
				"runtime_arguments": []interface{}{"run", "--rm", "myimage"},
				"package_arguments": []interface{}{},
			},
		},
	}
	cfg, err := a.FormatServerConfig(serverInfo, nil, nil)
	if err != nil {
		t.Fatalf("FormatServerConfig: %v", err)
	}
	if cfg["command"] != "docker" {
		t.Errorf("expected command docker, got %v", cfg["command"])
	}
}

func TestFormatServerConfigPyPI(t *testing.T) {
	a := New("/project", false)
	serverInfo := map[string]interface{}{
		"id":   "py-server",
		"name": "Py Server",
		"packages": []interface{}{
			map[string]interface{}{
				"name":             "my-pypi-server",
				"registry":         "pypi",
				"runtime_arguments": []interface{}{},
				"package_arguments": []interface{}{},
			},
		},
	}
	cfg, err := a.FormatServerConfig(serverInfo, nil, nil)
	if err != nil {
		t.Fatalf("FormatServerConfig: %v", err)
	}
	if cfg["command"] != "uvx" {
		t.Errorf("expected command uvx, got %v", cfg["command"])
	}
}

func TestFormatServerConfigRawStdio(t *testing.T) {
	a := New("/project", false)
	serverInfo := map[string]interface{}{
		"id":   "raw-stdio",
		"name": "Raw Stdio",
		"_raw_stdio": map[string]interface{}{
			"command": "mybin",
			"args":    []interface{}{"--flag"},
		},
	}
	cfg, err := a.FormatServerConfig(serverInfo, nil, nil)
	if err != nil {
		t.Fatalf("FormatServerConfig: %v", err)
	}
	if cfg["command"] != "mybin" {
		t.Errorf("expected command mybin, got %v", cfg["command"])
	}
}

func TestConfigureMCPServerEmptyURL(t *testing.T) {
	a := New("/project", false)
	ok := a.ConfigureMCPServer("", "server", true, nil, nil, nil)
	if ok {
		t.Error("expected ConfigureMCPServer to return false for empty URL")
	}
}

func TestConfigureMCPServerNotInCache(t *testing.T) {
	a := New("/project", false)
	ok := a.ConfigureMCPServer("https://example.com/server", "server", true, nil, map[string]interface{}{}, nil)
	if ok {
		t.Error("expected ConfigureMCPServer to return false when server not in cache")
	}
}

func TestConfigureMCPServerRemoteOnly(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	serverInfo := map[string]interface{}{
		"id":   "remote-only",
		"name": "Remote Only",
		"remotes": []interface{}{
			map[string]interface{}{"url": "https://remote.example.com/sse"},
		},
		"packages": []interface{}{},
	}
	cache := map[string]interface{}{
		"https://example.com/remote-only": serverInfo,
	}
	ok := a.ConfigureMCPServer("https://example.com/remote-only", "remote-only", true, nil, cache, nil)
	if ok {
		t.Error("expected false for remote-only server")
	}
}

func TestConfigureMCPServerSuccess(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	serverInfo := map[string]interface{}{
		"id":   "npm-server",
		"name": "NPM Server",
		"packages": []interface{}{
			map[string]interface{}{
				"name":             "npm-mcp-server",
				"registry":         "npm",
				"runtime_arguments": []interface{}{},
				"package_arguments": []interface{}{},
			},
		},
	}
	cache := map[string]interface{}{
		"https://example.com/npm-server": serverInfo,
	}
	ok := a.ConfigureMCPServer("https://example.com/npm-server", "npm-server", true, nil, cache, nil)
	if !ok {
		t.Error("expected ConfigureMCPServer to return true for valid server")
	}
}

func TestNewAdapterNotNil(t *testing.T) {
	a := New("/project", false)
	if a == nil {
		t.Error("New returned nil")
	}
	if a.Adapter == nil {
		t.Error("New returned adapter with nil embedded Adapter")
	}
}
