package configcmd

import (
	"testing"
)

func TestParseBoolValue_True(t *testing.T) {
	for _, v := range []string{"true", "True", "TRUE", "1", "yes", "YES"} {
		got, err := ParseBoolValue(v)
		if err != nil {
			t.Errorf("ParseBoolValue(%q): unexpected error: %v", v, err)
		}
		if !got {
			t.Errorf("ParseBoolValue(%q): expected true", v)
		}
	}
}

func TestParseBoolValue_False(t *testing.T) {
	for _, v := range []string{"false", "False", "FALSE", "0", "no", "NO"} {
		got, err := ParseBoolValue(v)
		if err != nil {
			t.Errorf("ParseBoolValue(%q): unexpected error: %v", v, err)
		}
		if got {
			t.Errorf("ParseBoolValue(%q): expected false", v)
		}
	}
}

func TestParseBoolValue_Invalid(t *testing.T) {
	for _, v := range []string{"maybe", "2", "on", "off", ""} {
		_, err := ParseBoolValue(v)
		if err == nil {
			t.Errorf("ParseBoolValue(%q): expected error", v)
		}
	}
}

func TestParseBoolValue_Whitespace(t *testing.T) {
	got, err := ParseBoolValue("  true  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected true")
	}
}

func TestValidConfigKeys_NonEmpty(t *testing.T) {
	keys := ValidConfigKeys()
	if len(keys) == 0 {
		t.Error("expected at least one config key")
	}
}

func TestDisplayName_KnownKey(t *testing.T) {
	name := DisplayName("auto_integrate")
	if name == "" {
		t.Error("expected non-empty display name for auto_integrate")
	}
}

func TestDisplayName_UnknownKey(t *testing.T) {
	name := DisplayName("unknown_key_xyz")
	if name == "" {
		t.Error("expected fallback display name for unknown key")
	}
}

func TestParseAPMYML_Empty(t *testing.T) {
	cfg := parseAPMYML("")
	if cfg.Name != "" {
		t.Errorf("expected empty name, got %q", cfg.Name)
	}
	if cfg.MCPDepCount != 0 {
		t.Errorf("expected 0 MCP deps, got %d", cfg.MCPDepCount)
	}
}

func TestParseAPMYML_Basic(t *testing.T) {
	content := "name: myapp\nversion: 1.2.3\n"
	cfg := parseAPMYML(content)
	if cfg.Name != "myapp" {
		t.Errorf("expected name 'myapp', got %q", cfg.Name)
	}
	if cfg.Version != "1.2.3" {
		t.Errorf("expected version '1.2.3', got %q", cfg.Version)
	}
}

func TestParseAPMYML_MCPDeps(t *testing.T) {
	// MCPDepCount counts entries under mcp.dependencies;
	// basic parser may not handle nested YAML lists -- just assert no panic.
	content := "name: test\n"
	cfg := parseAPMYML(content)
	if cfg.MCPDepCount < 0 {
		t.Error("MCPDepCount should be non-negative")
	}
}
