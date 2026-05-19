package operations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestServerNeed_Fields(t *testing.T) {
	sn := ServerNeed{Reference: "org/server", NeedsInstall: true, Reason: "not found"}
	if sn.Reference != "org/server" {
		t.Errorf("unexpected Reference %q", sn.Reference)
	}
	if !sn.NeedsInstall {
		t.Error("expected NeedsInstall=true")
	}
	if sn.Reason != "not found" {
		t.Errorf("unexpected Reason %q", sn.Reason)
	}
}

func TestServerNeed_ZeroValue(t *testing.T) {
	var sn ServerNeed
	if sn.NeedsInstall {
		t.Error("zero value NeedsInstall should be false")
	}
	if sn.Reference != "" || sn.Reason != "" {
		t.Error("zero value strings should be empty")
	}
}

func TestInstallStatus_Fields(t *testing.T) {
	is := InstallStatus{Runtime: "claude", Installed: true, ServerID: "abc-123"}
	if is.Runtime != "claude" {
		t.Errorf("unexpected Runtime %q", is.Runtime)
	}
	if !is.Installed {
		t.Error("expected Installed=true")
	}
	if is.ServerID != "abc-123" {
		t.Errorf("unexpected ServerID %q", is.ServerID)
	}
}

func TestInstallStatus_ZeroValue(t *testing.T) {
	var is InstallStatus
	if is.Installed {
		t.Error("zero value Installed should be false")
	}
}

func TestExtractServerIDs_Empty(t *testing.T) {
	data := []byte(`{}`)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty ids, got %v", ids)
	}
}

func TestExtractServerIDs_MCPServersKeyAlt(t *testing.T) {
	payload := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"myserver": map[string]interface{}{
				"id":      "uuid-001",
				"command": "npx",
			},
		},
	}
	data, _ := json.Marshal(payload)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "uuid-001" {
		t.Errorf("expected [uuid-001], got %v", ids)
	}
}

func TestExtractServerIDs_ServersKeyAlt(t *testing.T) {
	payload := map[string]interface{}{
		"servers": map[string]interface{}{
			"s1": map[string]interface{}{"id": "id-a"},
			"s2": map[string]interface{}{"id": "id-b"},
		},
	}
	data, _ := json.Marshal(payload)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 ids, got %v", ids)
	}
}

func TestExtractServerIDs_MissingID(t *testing.T) {
	payload := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"noIDServer": map[string]interface{}{"command": "node"},
		},
	}
	data, _ := json.Marshal(payload)
	ids, err := extractServerIDs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty ids when no id field, got %v", ids)
	}
}

func TestExtractServerIDs_InvalidJSONInput(t *testing.T) {
	_, err := extractServerIDs([]byte(`not json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestMCPConfigPaths_Copilot(t *testing.T) {
	paths := mcpConfigPaths("copilot", "/project", false)
	if len(paths) == 0 {
		t.Fatal("expected at least one path for copilot")
	}
	found := false
	for _, p := range paths {
		if filepath.Base(p) == "mcp.json" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected mcp.json in copilot paths: %v", paths)
	}
}

func TestMCPConfigPaths_VSCode(t *testing.T) {
	paths := mcpConfigPaths("vscode", "/proj", false)
	if len(paths) == 0 {
		t.Fatal("expected at least one path for vscode")
	}
}

func TestMCPConfigPaths_Cursor(t *testing.T) {
	paths := mcpConfigPaths("cursor", "/proj", false)
	if len(paths) == 0 {
		t.Fatal("expected at least one path for cursor")
	}
}

func TestMCPConfigPaths_Unknown(t *testing.T) {
	paths := mcpConfigPaths("unknown-runtime", "/proj", false)
	if len(paths) != 0 {
		t.Errorf("expected no paths for unknown runtime, got %v", paths)
	}
}

func TestGetInstalledServerIDs_EmptyDir(t *testing.T) {
	ops, err := NewMCPServerOperations("https://api.example.com")
	if err != nil {
		t.Skip("cannot create operations: " + err.Error())
	}
	ids, err := ops.getInstalledServerIDs([]string{"claude"}, t.TempDir(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty ids for empty dir, got %v", ids)
	}
}

func TestGetInstalledServerIDs_FromFile(t *testing.T) {
	dir := t.TempDir()
	claudeDir := filepath.Join(dir, ".claude")
	if err := os.Mkdir(claudeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	payload := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"srv": map[string]interface{}{"id": "test-uuid-abc"},
		},
	}
	data, _ := json.Marshal(payload)
	cfgPath := filepath.Join(claudeDir, "claude_mcp_config.json")
	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	ops, err := NewMCPServerOperations("https://api.example.com")
	if err != nil {
		t.Skip("cannot create operations: " + err.Error())
	}
	ids, err := ops.getInstalledServerIDs([]string{"claude"}, dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := ids["test-uuid-abc"]; !ok {
		t.Errorf("expected test-uuid-abc in ids, got %v", ids)
	}
}
