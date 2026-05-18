package mcpintegrator

import (
	"testing"
)

func TestMCPServer_Fields(t *testing.T) {
	s := MCPServer{
		Name:        "my-server",
		Command:     "npx",
		Args:        []string{"-y", "my-mcp"},
		Env:         map[string]string{"TOKEN": "abc"},
		Type:        "stdio",
		URL:         "",
		Description: "My server",
		Scope:       "project",
	}
	if s.Name != "my-server" {
		t.Errorf("Name: %q", s.Name)
	}
	if len(s.Args) != 2 {
		t.Errorf("Args length: %d", len(s.Args))
	}
	if s.Env["TOKEN"] != "abc" {
		t.Errorf("Env TOKEN: %q", s.Env["TOKEN"])
	}
}

func TestMCPLockEntry_Fields(t *testing.T) {
	e := MCPLockEntry{
		Name:        "server-x",
		ResolvedRef: "refs/heads/main",
		Commit:      "abc1234",
		Source:      "github",
	}
	if e.Name != "server-x" {
		t.Errorf("Name: %q", e.Name)
	}
	if e.Commit != "abc1234" {
		t.Errorf("Commit: %q", e.Commit)
	}
}

func TestIntegrateOptions_Fields(t *testing.T) {
	opts := IntegrateOptions{
		ProjectRoot: "/my/project",
		DryRun:      true,
		Verbose:     false,
		Force:       true,
		UserScope:   false,
		Targets:     []string{"copilot", "cursor"},
	}
	if opts.ProjectRoot != "/my/project" {
		t.Errorf("ProjectRoot: %q", opts.ProjectRoot)
	}
	if !opts.DryRun {
		t.Error("DryRun should be true")
	}
	if len(opts.Targets) != 2 {
		t.Errorf("Targets length: %d", len(opts.Targets))
	}
}

func TestNormaliseServerName_AtPrefixLong(t *testing.T) {
	if got := NormaliseServerName("@Org/Server"); got != "org/server" {
		t.Errorf("got %q", got)
	}
}

func TestNormaliseServerName_Underscore(t *testing.T) {
	if got := NormaliseServerName("my_server"); got != "my_server" {
		t.Errorf("got %q", got)
	}
}

func TestDetectConflicts_MultipleConflicts(t *testing.T) {
	byPackage := map[string][]MCPServer{
		"pkgA": {{Name: "s1"}, {Name: "s2"}},
		"pkgB": {{Name: "s1"}, {Name: "s3"}},
		"pkgC": {{Name: "s2"}},
	}
	results := DetectConflicts(byPackage)
	if len(results) < 2 {
		t.Errorf("expected >=2 conflicts, got %d", len(results))
	}
}

func TestDetectConflicts_SinglePackage(t *testing.T) {
	byPackage := map[string][]MCPServer{
		"pkgA": {{Name: "server1"}, {Name: "server2"}},
	}
	results := DetectConflicts(byPackage)
	if len(results) != 0 {
		t.Errorf("single package should not produce conflicts, got %d", len(results))
	}
}

func TestNew_VerboseMode(t *testing.T) {
	mi, err := New("/tmp", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mi == nil {
		t.Fatal("expected non-nil MCPIntegrator")
	}
}

func TestIsVSCodeAvailable_EmptyPath(t *testing.T) {
	result := IsVSCodeAvailable("")
	_ = result
}

func TestIsCursorAvailable_EmptyPath(t *testing.T) {
	result := IsCursorAvailable("")
	_ = result
}

func TestStaleReport_Fields(t *testing.T) {
	report := StaleReport{
		Client:  "copilot",
		Servers: []string{"old-server", "stale-server"},
	}
	if report.Client != "copilot" {
		t.Errorf("Client: %q", report.Client)
	}
	if len(report.Servers) != 2 {
		t.Errorf("Servers length: %d", len(report.Servers))
	}
}

func TestConflictResult_Fields(t *testing.T) {
	cr := ConflictResult{
		ServerName: "conflicting",
		PackageA:   "pkgA",
		PackageB:   "pkgB",
	}
	if cr.ServerName != "conflicting" {
		t.Errorf("ServerName: %q", cr.ServerName)
	}
	if cr.PackageA != "pkgA" {
		t.Errorf("PackageA: %q", cr.PackageA)
	}
	if cr.PackageB != "pkgB" {
		t.Errorf("PackageB: %q", cr.PackageB)
	}
}
