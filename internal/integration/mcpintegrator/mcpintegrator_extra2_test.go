package mcpintegrator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNormaliseServerName_AtPrefix(t *testing.T) {
	if got := NormaliseServerName("@scope/my-server"); got != "scope/my-server" {
		t.Errorf("got %q", got)
	}
}

func TestNormaliseServerName_AlreadyLower(t *testing.T) {
	if got := NormaliseServerName("my-server"); got != "my-server" {
		t.Errorf("got %q", got)
	}
}

func TestNormaliseServerName_UpperCase(t *testing.T) {
	if got := NormaliseServerName("MyServer"); got != "myserver" {
		t.Errorf("got %q", got)
	}
}

func TestNormaliseServerName_EmptyVariant(t *testing.T) {
	if got := NormaliseServerName(""); got != "" {
		t.Errorf("got %q", got)
	}
}

func TestDetectConflicts_ByPackageMultipleServers(t *testing.T) {
	byPkg := map[string][]MCPServer{
		"pkgA": {{Name: "alpha"}, {Name: "beta"}},
		"pkgB": {{Name: "alpha"}, {Name: "gamma"}},
	}
	conflicts := DetectConflicts(byPkg)
	if len(conflicts) == 0 {
		t.Error("expected at least one conflict")
	}
	for _, c := range conflicts {
		if c.ServerName != "alpha" {
			t.Errorf("unexpected conflict server: %q", c.ServerName)
		}
	}
}

func TestDetectConflicts_EmptyMap(t *testing.T) {
	conflicts := DetectConflicts(map[string][]MCPServer{})
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestDetectConflicts_OnePackage(t *testing.T) {
	byPkg := map[string][]MCPServer{
		"pkgA": {{Name: "alpha"}, {Name: "beta"}},
	}
	conflicts := DetectConflicts(byPkg)
	if len(conflicts) != 0 {
		t.Errorf("single package cannot conflict with itself, got %d", len(conflicts))
	}
}

func TestIntegrateResult_ZeroValue(t *testing.T) {
	var r IntegrateResult
	if r.ServersAdded != nil {
		t.Error("ServersAdded should be nil")
	}
	if r.ServersRemoved != nil {
		t.Error("ServersRemoved should be nil")
	}
}

func TestIntegrateResult_Fields(t *testing.T) {
	r := IntegrateResult{
		ServersAdded:   []string{"a", "b"},
		ServersRemoved: []string{"c"},
		ServersSkipped: []string{"d"},
		ConfigsWritten: []string{"/tmp/mcp.json"},
		Warnings:       []string{"warn1"},
	}
	if len(r.ServersAdded) != 2 {
		t.Errorf("ServersAdded len=%d", len(r.ServersAdded))
	}
	if r.Warnings[0] != "warn1" {
		t.Errorf("Warnings[0]=%q", r.Warnings[0])
	}
}

func TestMCPIntegrator_ClientConfigPath_VSCode(t *testing.T) {
	root := t.TempDir()
	m := &MCPIntegrator{ProjectRoot: root}
	p := m.clientConfigPath("vscode")
	if filepath.Base(p) != "mcp.json" {
		t.Errorf("expected mcp.json, got %q", p)
	}
	if filepath.Base(filepath.Dir(p)) != ".vscode" {
		t.Errorf("expected .vscode dir, got %q", filepath.Dir(p))
	}
}

func TestMCPIntegrator_ClientConfigPath_Cursor(t *testing.T) {
	root := t.TempDir()
	m := &MCPIntegrator{ProjectRoot: root}
	p := m.clientConfigPath("cursor")
	if filepath.Base(filepath.Dir(p)) != ".cursor" {
		t.Errorf("expected .cursor dir, got %q", filepath.Dir(p))
	}
}

func TestMCPIntegrator_ClientConfigPath_Unknown(t *testing.T) {
	root := t.TempDir()
	m := &MCPIntegrator{ProjectRoot: root}
	p := m.clientConfigPath("unknown-client-xyz")
	if p != "" {
		t.Errorf("expected empty path, got %q", p)
	}
}

func TestMCPIntegrator_FindStaleServers_EmptyDir(t *testing.T) {
	root := t.TempDir()
	m := &MCPIntegrator{ProjectRoot: root}
	reports, err := m.FindStaleServers(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("expected 0 reports, got %d", len(reports))
	}
}

func TestMCPIntegrator_FindStaleServers_WithVSCodeConfig(t *testing.T) {
	root := t.TempDir()
	vscodeDir := filepath.Join(root, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"stale-server": map[string]interface{}{"command": "npx"},
		},
	}
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(filepath.Join(vscodeDir, "mcp.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	m := &MCPIntegrator{ProjectRoot: root}
	reports, err := m.FindStaleServers([]MCPServer{{Name: "active-server"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) == 0 {
		t.Error("expected stale report for vscode")
	}
	if reports[0].Client != "vscode" {
		t.Errorf("expected vscode client, got %q", reports[0].Client)
	}
	if len(reports[0].Servers) == 0 || reports[0].Servers[0] != "stale-server" {
		t.Errorf("expected stale-server, got %v", reports[0].Servers)
	}
}

func TestMCPIntegrator_FindStaleServers_NoStale(t *testing.T) {
	root := t.TempDir()
	vscodeDir := filepath.Join(root, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"my-server": map[string]interface{}{"command": "npx"},
		},
	}
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(filepath.Join(vscodeDir, "mcp.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	m := &MCPIntegrator{ProjectRoot: root}
	reports, err := m.FindStaleServers([]MCPServer{{Name: "my-server"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("expected no stale reports, got %d", len(reports))
	}
}

func TestMCPServer_ZeroValue(t *testing.T) {
	var s MCPServer
	if s.Name != "" || s.Command != "" || s.Type != "" {
		t.Error("zero value fields should be empty")
	}
	if s.Args != nil || s.Env != nil {
		t.Error("zero value slices/maps should be nil")
	}
}

func TestIntegrateOptions_ZeroValue(t *testing.T) {
	var opts IntegrateOptions
	if opts.DryRun || opts.Verbose || opts.Force || opts.UserScope {
		t.Error("bool fields should be false")
	}
}
