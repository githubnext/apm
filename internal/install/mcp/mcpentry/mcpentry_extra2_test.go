package mcpentry

import "testing"

func TestMCPEntry_ZeroValue(t *testing.T) {
	e := MCPEntry{}
	if e.Name != "" || e.Transport != "" || e.URL != "" || e.Command != "" {
		t.Error("zero value MCPEntry should have empty string fields")
	}
	if e.IsSelfDefined() {
		t.Error("zero value MCPEntry should not be self-defined")
	}
}

func TestBuildMCPEntry_BareRegistry(t *testing.T) {
	e, selfDefined := BuildMCPEntry("myplugin", "", "", nil, nil, "", nil, "")
	if selfDefined {
		t.Error("bare registry entry should not be self-defined")
	}
	if e.Kind != EntryKindRegistryShorthand {
		t.Errorf("expected shorthand kind, got %v", e.Kind)
	}
	if e.Name != "myplugin" {
		t.Errorf("expected myplugin, got %q", e.Name)
	}
}

func TestBuildMCPEntry_RemoteNoHeaders(t *testing.T) {
	e, selfDefined := BuildMCPEntry("srv", "", "https://example.com", nil, nil, "", nil, "")
	if !selfDefined {
		t.Error("remote entry should be self-defined")
	}
	if e.URL != "https://example.com" {
		t.Errorf("expected url, got %q", e.URL)
	}
	if e.Headers != nil {
		t.Error("no headers should be nil")
	}
}

func TestBuildMCPEntry_RemoteDefaultTransport(t *testing.T) {
	e, _ := BuildMCPEntry("srv", "", "https://example.com", nil, nil, "", nil, "")
	if e.Transport != "http" {
		t.Errorf("expected default http transport, got %q", e.Transport)
	}
}

func TestBuildMCPEntry_RegistryWithCustomURL(t *testing.T) {
	e, selfDefined := BuildMCPEntry("pkg", "", "", nil, nil, "", nil, "https://myregistry.example.com")
	if selfDefined {
		t.Error("registry entry should not be self-defined")
	}
	if e.Kind != EntryKindRegistryDict {
		t.Errorf("expected registry dict kind, got %v", e.Kind)
	}
	if e.Registry != "https://myregistry.example.com" {
		t.Errorf("unexpected registry: %v", e.Registry)
	}
}

func TestBuildMCPEntry_StdioSingleArg(t *testing.T) {
	e, selfDefined := BuildMCPEntry("tool", "", "", nil, nil, "", []string{"mytool"}, "")
	if !selfDefined {
		t.Error("stdio entry should be self-defined")
	}
	if e.Command != "mytool" {
		t.Errorf("expected mytool, got %q", e.Command)
	}
	if len(e.Args) != 0 {
		t.Errorf("expected no args, got %v", e.Args)
	}
}

func TestIsSelfDefined_AllKinds(t *testing.T) {
	cases := []struct {
		kind     EntryKind
		expected bool
	}{
		{EntryKindRegistryShorthand, false},
		{EntryKindRegistryDict, false},
		{EntryKindSelfDefinedStdio, true},
		{EntryKindSelfDefinedRemote, true},
	}
	for _, c := range cases {
		e := MCPEntry{Kind: c.kind}
		if e.IsSelfDefined() != c.expected {
			t.Errorf("kind %v: expected isSelfDefined=%v", c.kind, c.expected)
		}
	}
}
