package opencode

import (
	"testing"
)

func TestServerEntry_ZeroValue(t *testing.T) {
	var se ServerEntry
	if se.Type != "" || se.URL != "" || se.Enabled {
		t.Error("zero ServerEntry should have empty Type, URL and Enabled=false")
	}
	if se.Command != nil || se.Headers != nil || se.Environment != nil {
		t.Error("zero ServerEntry should have nil slices/maps")
	}
}

func TestCopilotEntry_ZeroValue(t *testing.T) {
	var ce CopilotEntry
	if ce.Command != "" || ce.URL != "" {
		t.Error("zero CopilotEntry fields should be empty")
	}
}

func TestToOpenCodeFormat_EnabledFlagRespected(t *testing.T) {
	entry := CopilotEntry{Command: "npx", Args: []string{"-y", "pkg"}}
	enabled := ToOpenCodeFormat(entry, true)
	disabled := ToOpenCodeFormat(entry, false)
	if !enabled.Enabled {
		t.Error("expected Enabled=true")
	}
	if disabled.Enabled {
		t.Error("expected Enabled=false")
	}
}

func TestToOpenCodeFormat_TypeLocalForCommand(t *testing.T) {
	entry := CopilotEntry{Command: "node", Args: []string{"server.js"}}
	result := ToOpenCodeFormat(entry, true)
	if result.Type != "local" {
		t.Errorf("expected Type=local, got %q", result.Type)
	}
}

func TestToOpenCodeFormat_TypeRemoteForURL(t *testing.T) {
	entry := CopilotEntry{URL: "https://mcp.example.com/server"}
	result := ToOpenCodeFormat(entry, true)
	if result.Type != "remote" {
		t.Errorf("expected Type=remote, got %q", result.Type)
	}
	if result.URL != "https://mcp.example.com/server" {
		t.Errorf("expected URL preserved, got %q", result.URL)
	}
}

func TestToOpenCodeFormat_CommandPrependedToArgs(t *testing.T) {
	entry := CopilotEntry{Command: "npx", Args: []string{"-y", "mypkg"}}
	result := ToOpenCodeFormat(entry, true)
	if len(result.Command) != 3 {
		t.Fatalf("expected 3 command elements, got %d", len(result.Command))
	}
	if result.Command[0] != "npx" {
		t.Errorf("expected Command[0]=npx, got %q", result.Command[0])
	}
	if result.Command[1] != "-y" {
		t.Errorf("expected Command[1]=-y, got %q", result.Command[1])
	}
}

func TestToOpenCodeFormat_EnvMapped(t *testing.T) {
	entry := CopilotEntry{
		Command: "uvx",
		Env:     map[string]string{"TOKEN": "secret"},
	}
	result := ToOpenCodeFormat(entry, true)
	if result.Environment["TOKEN"] != "secret" {
		t.Errorf("expected TOKEN=secret in Environment, got %q", result.Environment["TOKEN"])
	}
}

func TestToOpenCodeFormat_NoEnvWhenEmpty(t *testing.T) {
	entry := CopilotEntry{Command: "uvx"}
	result := ToOpenCodeFormat(entry, true)
	if result.Environment != nil {
		t.Error("expected nil Environment when no env vars")
	}
}

func TestToOpenCodeFormat_URLHeadersPreserved(t *testing.T) {
	entry := CopilotEntry{
		URL:     "https://api.example.com/mcp",
		Headers: map[string]string{"Authorization": "Bearer tok"},
	}
	result := ToOpenCodeFormat(entry, true)
	if result.Headers["Authorization"] != "Bearer tok" {
		t.Errorf("expected Authorization header preserved, got %q", result.Headers["Authorization"])
	}
}

func TestToOpenCodeFormat_EmptyCommandEntry(t *testing.T) {
	// Neither Command nor URL set -- should not panic
	entry := CopilotEntry{}
	result := ToOpenCodeFormat(entry, false)
	if result.Type != "local" {
		t.Errorf("expected default Type=local, got %q", result.Type)
	}
}
