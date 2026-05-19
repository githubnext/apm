package mcpentry_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpentry"
)

func TestEntryKind_Constants(t *testing.T) {
	// Check that all four EntryKind constants are distinct
	kinds := []mcpentry.EntryKind{
		mcpentry.EntryKindRegistryShorthand,
		mcpentry.EntryKindRegistryDict,
		mcpentry.EntryKindSelfDefinedStdio,
		mcpentry.EntryKindSelfDefinedRemote,
	}
	seen := make(map[mcpentry.EntryKind]bool)
	for _, k := range kinds {
		if seen[k] {
			t.Errorf("duplicate EntryKind value: %d", k)
		}
		seen[k] = true
	}
}

func TestIsSelfDefined_RegistryShorthand(t *testing.T) {
	e := mcpentry.MCPEntry{Kind: mcpentry.EntryKindRegistryShorthand}
	if e.IsSelfDefined() {
		t.Error("RegistryShorthand should not be self-defined")
	}
}

func TestIsSelfDefined_RegistryDict(t *testing.T) {
	e := mcpentry.MCPEntry{Kind: mcpentry.EntryKindRegistryDict}
	if e.IsSelfDefined() {
		t.Error("RegistryDict should not be self-defined")
	}
}

func TestIsSelfDefined_SelfDefinedStdio(t *testing.T) {
	e := mcpentry.MCPEntry{Kind: mcpentry.EntryKindSelfDefinedStdio}
	if !e.IsSelfDefined() {
		t.Error("SelfDefinedStdio should be self-defined")
	}
}

func TestIsSelfDefined_SelfDefinedRemote(t *testing.T) {
	e := mcpentry.MCPEntry{Kind: mcpentry.EntryKindSelfDefinedRemote}
	if !e.IsSelfDefined() {
		t.Error("SelfDefinedRemote should be self-defined")
	}
}

func TestBuildMCPEntry_StdioMultipleArgs(t *testing.T) {
	e, self := mcpentry.BuildMCPEntry("s", "", "", nil, nil, "", []string{"python", "-m", "server", "--port", "9000"}, "")
	if !self {
		t.Fatal("expected self-defined")
	}
	if e.Command != "python" {
		t.Errorf("Command = %q, want python", e.Command)
	}
	if len(e.Args) != 4 {
		t.Errorf("Args len = %d, want 4: %v", len(e.Args), e.Args)
	}
	if e.Args[0] != "-m" {
		t.Errorf("Args[0] = %q, want -m", e.Args[0])
	}
}

func TestBuildMCPEntry_StdioNoArgs(t *testing.T) {
	e, self := mcpentry.BuildMCPEntry("s", "", "", nil, nil, "", []string{"mybin"}, "")
	if !self {
		t.Fatal("expected self-defined")
	}
	if e.Command != "mybin" {
		t.Errorf("Command = %q, want mybin", e.Command)
	}
	if len(e.Args) != 0 {
		t.Errorf("Args should be empty when only binary given, got %v", e.Args)
	}
}

func TestBuildMCPEntry_StdioEnvMultipleKeys(t *testing.T) {
	env := map[string]string{"API_KEY": "secret", "DEBUG": "1", "PORT": "8080"}
	e, _ := mcpentry.BuildMCPEntry("s", "", "", env, nil, "", []string{"server"}, "")
	if len(e.Env) != 3 {
		t.Errorf("expected 3 env keys, got %d", len(e.Env))
	}
	if e.Env["API_KEY"] != "secret" {
		t.Errorf("Env[API_KEY] = %q, want secret", e.Env["API_KEY"])
	}
}

func TestBuildMCPEntry_RemoteWithHeaders(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token", "X-Custom": "val"}
	e, self := mcpentry.BuildMCPEntry("r", "http", "https://api.example.com/mcp", nil, headers, "", nil, "")
	if !self {
		t.Fatal("expected self-defined")
	}
	if e.Kind != mcpentry.EntryKindSelfDefinedRemote {
		t.Errorf("Kind = %v", e.Kind)
	}
	if len(e.Headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(e.Headers))
	}
	if e.Headers["Authorization"] != "Bearer token" {
		t.Errorf("Headers[Authorization] = %q", e.Headers["Authorization"])
	}
}

func TestBuildMCPEntry_RemoteSSETransport(t *testing.T) {
	e, _ := mcpentry.BuildMCPEntry("r", "sse", "https://sse.example.com/mcp", nil, nil, "", nil, "")
	if e.Transport != "sse" {
		t.Errorf("Transport = %q, want sse", e.Transport)
	}
}

func TestBuildMCPEntry_RegistryWithVersion(t *testing.T) {
	e, self := mcpentry.BuildMCPEntry("mypkg", "", "", nil, nil, "~2.0.0", nil, "")
	if self {
		t.Fatal("expected not self-defined")
	}
	if e.Version != "~2.0.0" {
		t.Errorf("Version = %q, want ~2.0.0", e.Version)
	}
}

func TestBuildMCPEntry_RegistryBoolTrue(t *testing.T) {
	e, _ := mcpentry.BuildMCPEntry("mypkg", "", "", nil, nil, "", nil, "")
	if e.Registry != true {
		t.Errorf("Registry for shorthand should be true, got %v", e.Registry)
	}
}

func TestMCPEntry_NameField(t *testing.T) {
	e, _ := mcpentry.BuildMCPEntry("my-server", "", "", nil, nil, "", []string{"cmd"}, "")
	if e.Name != "my-server" {
		t.Errorf("Name = %q, want my-server", e.Name)
	}
}

func TestBuildMCPEntry_StdioRegistryIsFalse(t *testing.T) {
	e, _ := mcpentry.BuildMCPEntry("s", "", "", nil, nil, "", []string{"cmd"}, "")
	if e.Registry != false {
		t.Errorf("Registry for self-defined stdio should be false, got %v", e.Registry)
	}
}
