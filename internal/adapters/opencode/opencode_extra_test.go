package opencode_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/adapters/opencode"
)

func TestToOpenCodeFormat_NoCommand(t *testing.T) {
	entry := opencode.CopilotEntry{}
	got := opencode.ToOpenCodeFormat(entry, true)
	if got.Type != "local" {
		t.Errorf("expected local type, got %q", got.Type)
	}
	if len(got.Command) != 0 {
		t.Errorf("expected no command, got %v", got.Command)
	}
}

func TestToOpenCodeFormat_ArgsPreserved(t *testing.T) {
	entry := opencode.CopilotEntry{
		Command: "node",
		Args:    []string{"server.js", "--port", "3000"},
	}
	got := opencode.ToOpenCodeFormat(entry, true)
	if len(got.Command) != 4 {
		t.Fatalf("expected 4 command parts, got %d: %v", len(got.Command), got.Command)
	}
	if got.Command[0] != "node" {
		t.Errorf("Command[0] = %q, want node", got.Command[0])
	}
	if got.Command[3] != "3000" {
		t.Errorf("Command[3] = %q, want 3000", got.Command[3])
	}
}

func TestToOpenCodeFormat_URLWithNoHeaders(t *testing.T) {
	entry := opencode.CopilotEntry{URL: "https://example.com/mcp"}
	got := opencode.ToOpenCodeFormat(entry, true)
	if got.Type != "remote" {
		t.Errorf("expected remote type, got %q", got.Type)
	}
	if got.URL != "https://example.com/mcp" {
		t.Errorf("URL mismatch: %q", got.URL)
	}
	if len(got.Headers) != 0 {
		t.Errorf("expected no headers, got %v", got.Headers)
	}
}

func TestToOpenCodeFormat_EmptyEnv(t *testing.T) {
	entry := opencode.CopilotEntry{Command: "cmd"}
	got := opencode.ToOpenCodeFormat(entry, true)
	if len(got.Environment) != 0 {
		t.Errorf("expected no environment, got %v", got.Environment)
	}
}

func TestServerEntry_Fields(t *testing.T) {
	e := opencode.ServerEntry{
		Type:    "local",
		Command: []string{"npx", "-y", "pkg"},
		Enabled: true,
		Environment: map[string]string{"API_KEY": "secret"},
	}
	if e.Type != "local" {
		t.Errorf("Type mismatch: %q", e.Type)
	}
	if len(e.Command) != 3 {
		t.Errorf("Command length: %d", len(e.Command))
	}
	if !e.Enabled {
		t.Error("Enabled should be true")
	}
	if e.Environment["API_KEY"] != "secret" {
		t.Errorf("Environment mismatch")
	}
}

func TestCopilotEntry_Fields(t *testing.T) {
	e := opencode.CopilotEntry{
		Command: "npx",
		Args:    []string{"-y", "pkg"},
		Env:     map[string]string{"KEY": "val"},
		URL:     "",
	}
	if e.Command != "npx" {
		t.Errorf("Command mismatch: %q", e.Command)
	}
	if len(e.Args) != 2 {
		t.Errorf("Args length: %d", len(e.Args))
	}
}

func TestGetCurrentConfig_WithFile(t *testing.T) {
	dir := t.TempDir()
	openDir := filepath.Join(dir, ".opencode")
	if err := os.Mkdir(openDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgFile := filepath.Join(dir, "opencode.json")
	cfg := map[string]interface{}{
		"mcp": map[string]interface{}{
			"my-server": map[string]interface{}{
				"type":    "local",
				"command": []string{"node", "s.js"},
				"enabled": true,
			},
		},
	}
	b, _ := json.Marshal(cfg)
	if err := os.WriteFile(cfgFile, b, 0o644); err != nil {
		t.Fatal(err)
	}
	adapter := opencode.New(dir)
	got := adapter.GetCurrentConfig()
	if got == nil {
		t.Fatal("expected non-nil config")
	}
	if _, ok := got["mcp"]; !ok {
		t.Error("expected 'mcp' key in config")
	}
}

func TestIsOptedIn_WithFile(t *testing.T) {
	dir := t.TempDir()
	openDir := filepath.Join(dir, ".opencode")
	if err := os.Mkdir(openDir, 0o755); err != nil {
		t.Fatal(err)
	}
	adapter := opencode.New(dir)
	if !adapter.IsOptedIn() {
		t.Error("IsOptedIn should return true when .opencode/ exists")
	}
}
