package mcpdep_test

import (
	"testing"

	"github.com/githubnext/apm/internal/models/mcpdep"
)

func TestFromString_Valid(t *testing.T) {
	d, err := mcpdep.FromString("github.com/owner/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name != "github.com/owner/repo" {
		t.Errorf("expected name 'github.com/owner/repo', got %q", d.Name)
	}
}

func TestFromString_Empty(t *testing.T) {
	_, err := mcpdep.FromString("")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestIsRegistryResolved_Default(t *testing.T) {
	d, _ := mcpdep.FromString("owner/repo")
	if !d.IsRegistryResolved() {
		t.Error("default dependency should be registry-resolved")
	}
}

func TestIsSelfDefined_RegistryFalse(t *testing.T) {
	d, err := mcpdep.FromDict(map[string]interface{}{
		"name":      "my-server",
		"registry":  false,
		"transport": "stdio",
		"command":   "npx",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.IsSelfDefined() {
		t.Error("expected self-defined with registry: false")
	}
}

func TestFromDict_BasicFields(t *testing.T) {
	d, err := mcpdep.FromDict(map[string]interface{}{
		"name":    "my-mcp",
		"version": "1.0.0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name != "my-mcp" {
		t.Errorf("expected 'my-mcp', got %q", d.Name)
	}
	if d.Version != "1.0.0" {
		t.Errorf("expected '1.0.0', got %q", d.Version)
	}
}

func TestFromDict_MissingName(t *testing.T) {
	_, err := mcpdep.FromDict(map[string]interface{}{})
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestFromDict_LegacyTransportType(t *testing.T) {
	d, err := mcpdep.FromDict(map[string]interface{}{
		"name": "srv",
		"type": "stdio",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Transport != "stdio" {
		t.Errorf("expected transport 'stdio' from legacy 'type', got %q", d.Transport)
	}
}

func TestFromDict_EnvAndHeaders(t *testing.T) {
	d, err := mcpdep.FromDict(map[string]interface{}{
		"name": "srv",
		"env": map[string]interface{}{"KEY": "val"},
		"headers": map[string]interface{}{"X-Token": "tok"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Env["KEY"] != "val" {
		t.Errorf("expected env KEY=val, got %q", d.Env["KEY"])
	}
	if d.Headers["X-Token"] != "tok" {
		t.Errorf("expected header X-Token=tok, got %q", d.Headers["X-Token"])
	}
}

func TestToDict_RoundTrip(t *testing.T) {
	original := map[string]interface{}{
		"name":    "my-mcp",
		"version": "2.0",
	}
	d, _ := mcpdep.FromDict(original)
	out := d.ToDict()
	if out["name"] != "my-mcp" {
		t.Errorf("round-trip name mismatch")
	}
	if out["version"] != "2.0" {
		t.Errorf("round-trip version mismatch")
	}
}

func TestToDict_SelfDefinedRegistry(t *testing.T) {
	d, _ := mcpdep.FromDict(map[string]interface{}{
		"name":      "srv",
		"registry":  false,
		"transport": "stdio",
		"command":   "run",
	})
	out := d.ToDict()
	reg, ok := out["registry"].(bool)
	if !ok || reg != false {
		t.Errorf("expected registry=false in ToDict output, got %v", out["registry"])
	}
}

func TestString_WithTransport(t *testing.T) {
	d := &mcpdep.MCPDependency{Name: "srv", Transport: "stdio"}
	s := d.String()
	if s == "" {
		t.Error("String() should not be empty")
	}
}
