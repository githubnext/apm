package mcpentry_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpentry"
)

func TestIsSelfDefined(t *testing.T) {
	tests := []struct {
		kind mcpentry.EntryKind
		want bool
	}{
		{mcpentry.EntryKindRegistryShorthand, false},
		{mcpentry.EntryKindRegistryDict, false},
		{mcpentry.EntryKindSelfDefinedStdio, true},
		{mcpentry.EntryKindSelfDefinedRemote, true},
	}
	for _, tc := range tests {
		e := mcpentry.MCPEntry{Kind: tc.kind}
		if got := e.IsSelfDefined(); got != tc.want {
			t.Errorf("kind=%d IsSelfDefined()=%v, want %v", tc.kind, got, tc.want)
		}
	}
}

func TestBuildMCPEntry_SelfDefinedStdio(t *testing.T) {
	e, self := mcpentry.BuildMCPEntry("myserver", "", "", nil, nil, "", []string{"npx", "-y", "server"}, "")
	if !self {
		t.Fatal("expected isSelfDefined=true")
	}
	if e.Kind != mcpentry.EntryKindSelfDefinedStdio {
		t.Errorf("kind=%v, want SelfDefinedStdio", e.Kind)
	}
	if e.Transport != "stdio" {
		t.Errorf("transport=%q, want stdio", e.Transport)
	}
	if e.Command != "npx" {
		t.Errorf("command=%q, want npx", e.Command)
	}
	if len(e.Args) != 2 || e.Args[0] != "-y" || e.Args[1] != "server" {
		t.Errorf("args=%v, want [-y server]", e.Args)
	}
}

func TestBuildMCPEntry_SelfDefinedStdio_WithEnv(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	e, self := mcpentry.BuildMCPEntry("srv", "", "", env, nil, "", []string{"cmd"}, "")
	if !self || e.Env["FOO"] != "bar" {
		t.Errorf("unexpected: self=%v env=%v", self, e.Env)
	}
}

func TestBuildMCPEntry_SelfDefinedRemote(t *testing.T) {
	e, self := mcpentry.BuildMCPEntry("remote", "sse", "https://example.com/mcp", nil, nil, "", nil, "")
	if !self {
		t.Fatal("expected isSelfDefined=true")
	}
	if e.Kind != mcpentry.EntryKindSelfDefinedRemote {
		t.Errorf("kind=%v, want SelfDefinedRemote", e.Kind)
	}
	if e.URL != "https://example.com/mcp" {
		t.Errorf("url=%q", e.URL)
	}
	if e.Transport != "sse" {
		t.Errorf("transport=%q, want sse", e.Transport)
	}
}

func TestBuildMCPEntry_SelfDefinedRemote_DefaultTransport(t *testing.T) {
	e, _ := mcpentry.BuildMCPEntry("r", "", "https://x.com", nil, nil, "", nil, "")
	if e.Transport != "http" {
		t.Errorf("transport=%q, want http", e.Transport)
	}
}

func TestBuildMCPEntry_RegistryDict(t *testing.T) {
	e, self := mcpentry.BuildMCPEntry("pkg", "", "", nil, nil, "^1.0.0", nil, "")
	if self {
		t.Fatal("expected isSelfDefined=false")
	}
	if e.Kind != mcpentry.EntryKindRegistryDict {
		t.Errorf("kind=%v, want RegistryDict", e.Kind)
	}
	if e.Version != "^1.0.0" {
		t.Errorf("version=%q", e.Version)
	}
	if e.Registry != true {
		t.Errorf("registry=%v, want true", e.Registry)
	}
}

func TestBuildMCPEntry_RegistryDict_WithRegistryURL(t *testing.T) {
	e, _ := mcpentry.BuildMCPEntry("pkg", "", "", nil, nil, "", nil, "https://my.registry.com")
	if e.Kind != mcpentry.EntryKindRegistryDict {
		t.Errorf("kind=%v, want RegistryDict", e.Kind)
	}
	if e.Registry != "https://my.registry.com" {
		t.Errorf("registry=%v", e.Registry)
	}
}

func TestBuildMCPEntry_RegistryShorthand(t *testing.T) {
	e, self := mcpentry.BuildMCPEntry("pkg", "", "", nil, nil, "", nil, "")
	if self {
		t.Fatal("expected isSelfDefined=false")
	}
	if e.Kind != mcpentry.EntryKindRegistryShorthand {
		t.Errorf("kind=%v, want RegistryShorthand", e.Kind)
	}
	if e.Registry != true {
		t.Errorf("registry=%v, want true", e.Registry)
	}
}
