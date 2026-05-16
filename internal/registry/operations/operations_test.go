package operations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMCPConfigPaths_Claude_UserScope(t *testing.T) {
	paths := mcpConfigPaths("claude", "/project", true)
	if len(paths) == 0 {
		t.Skip("no home dir available")
	}
	for _, p := range paths {
		if filepath.Base(filepath.Dir(p)) != ".claude" &&
			filepath.Base(filepath.Dir(filepath.Dir(p))) != "Claude" {
			// Accept either .claude or Library/Application Support/Claude
		}
	}
	// Should contain config file names
	for _, p := range paths {
		base := filepath.Base(p)
		if base != "claude_desktop_config.json" && base != "claude_mcp_config.json" {
			t.Errorf("unexpected config file: %s", base)
		}
	}
}

func TestMCPConfigPaths_Claude_ProjectScope(t *testing.T) {
	paths := mcpConfigPaths("claude", "/my/project", false)
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	want := filepath.Join("/my/project", ".claude", "claude_mcp_config.json")
	if paths[0] != want {
		t.Errorf("path = %q; want %q", paths[0], want)
	}
}

func TestMCPConfigPaths_Copilot_ProjectScope(t *testing.T) {
	paths := mcpConfigPaths("copilot", "/proj", false)
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	if filepath.Base(paths[0]) != "mcp.json" {
		t.Errorf("expected mcp.json, got %s", filepath.Base(paths[0]))
	}
}

func TestMCPConfigPaths_VSCode_ProjectScope(t *testing.T) {
	paths := mcpConfigPaths("vscode", "/proj", false)
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	want := filepath.Join("/proj", ".vscode", "mcp.json")
	if paths[0] != want {
		t.Errorf("path = %q; want %q", paths[0], want)
	}
}

func TestMCPConfigPaths_Cursor_ProjectScope(t *testing.T) {
	paths := mcpConfigPaths("cursor", "/proj", false)
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	want := filepath.Join("/proj", ".cursor", "mcp.json")
	if paths[0] != want {
		t.Errorf("path = %q; want %q", paths[0], want)
	}
}

func TestMCPConfigPaths_Unknown_Runtime(t *testing.T) {
	paths := mcpConfigPaths("unknown-runtime", "/proj", false)
	if len(paths) != 0 {
		t.Errorf("expected 0 paths for unknown runtime, got %d", len(paths))
	}
}

func TestMCPConfigPaths_EmptyProjectRoot(t *testing.T) {
	paths := mcpConfigPaths("copilot", "", false)
	if len(paths) != 0 {
		t.Errorf("expected 0 paths for empty project root, got %d", len(paths))
	}
}

func TestExtractServerIDs_MCPServersKey(t *testing.T) {
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"myserver": map[string]interface{}{
				"id":      "srv-uuid-001",
				"command": "npx",
			},
		},
	}
	raw, _ := json.Marshal(data)
	ids, err := extractServerIDs(raw)
	if err != nil {
		t.Fatalf("extractServerIDs: %v", err)
	}
	if len(ids) != 1 || ids[0] != "srv-uuid-001" {
		t.Errorf("ids = %v; want [srv-uuid-001]", ids)
	}
}

func TestExtractServerIDs_ServersKey(t *testing.T) {
	data := map[string]interface{}{
		"servers": map[string]interface{}{
			"server1": map[string]interface{}{"id": "id-aaa"},
			"server2": map[string]interface{}{"id": "id-bbb"},
		},
	}
	raw, _ := json.Marshal(data)
	ids, err := extractServerIDs(raw)
	if err != nil {
		t.Fatalf("extractServerIDs: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 ids, got %d: %v", len(ids), ids)
	}
}

func TestExtractServerIDs_NoIDField(t *testing.T) {
	data := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"myserver": map[string]interface{}{"command": "npx"},
		},
	}
	raw, _ := json.Marshal(data)
	ids, err := extractServerIDs(raw)
	if err != nil {
		t.Fatalf("extractServerIDs: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 ids, got %v", ids)
	}
}

func TestExtractServerIDs_InvalidJSON(t *testing.T) {
	_, err := extractServerIDs([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestExtractServerIDs_EmptyObject(t *testing.T) {
	ids, err := extractServerIDs([]byte("{}"))
	if err != nil {
		t.Fatalf("extractServerIDs: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 ids, got %v", ids)
	}
}

func TestGetInstalledServerIDs_ReadsConfigFile(t *testing.T) {
	dir := t.TempDir()
	vscodeDir := filepath.Join(dir, ".vscode")
	os.MkdirAll(vscodeDir, 0o755)

	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"myserver": map[string]interface{}{"id": "test-server-id"},
		},
	}
	raw, _ := json.Marshal(cfg)
	os.WriteFile(filepath.Join(vscodeDir, "mcp.json"), raw, 0o644)

	ops := &MCPServerOperations{}
	ids, err := ops.getInstalledServerIDs([]string{"copilot"}, dir, false)
	if err != nil {
		t.Fatalf("getInstalledServerIDs: %v", err)
	}
	if _, ok := ids["test-server-id"]; !ok {
		t.Errorf("expected test-server-id in ids, got %v", ids)
	}
}

func TestGetInstalledServerIDs_MissingFile(t *testing.T) {
	ops := &MCPServerOperations{}
	ids, err := ops.getInstalledServerIDs([]string{"copilot"}, "/nonexistent/proj", false)
	if err != nil {
		t.Fatalf("getInstalledServerIDs: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty ids for missing file, got %v", ids)
	}
}
