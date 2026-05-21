package mcpdep

import (
	"testing"
)

func TestMCPDependency_ZeroValue_Extra2(t *testing.T) {
	var d MCPDependency
	if d.Name != "" || d.Transport != "" || d.Version != "" {
		t.Error("zero-value MCPDependency should have empty string fields")
	}
}

func TestMCPDependency_Fields_Extra2(t *testing.T) {
	d := MCPDependency{
		Name:      "myserver",
		Transport: "stdio",
		Version:   "1.0.0",
		Package:   "myorg/myserver",
	}
	if d.Name != "myserver" {
		t.Errorf("Name = %q", d.Name)
	}
	if d.Transport != "stdio" {
		t.Errorf("Transport = %q", d.Transport)
	}
}

func TestFromString_NameOnly_Extra2(t *testing.T) {
	d, err := FromString("myserver")
	if err != nil {
		t.Fatalf("FromString error: %v", err)
	}
	if d.Name != "myserver" {
		t.Errorf("Name = %q, want myserver", d.Name)
	}
}

func TestFromString_WithOrg_Extra2(t *testing.T) {
	d, err := FromString("myorg/myserver")
	if err != nil {
		t.Fatalf("FromString error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil dependency")
	}
}

func TestFromString_Empty_Extra2(t *testing.T) {
	_, err := FromString("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestToDict_RoundTrip_Extra2(t *testing.T) {
	d := MCPDependency{
		Name:      "myserver",
		Transport: "stdio",
		Version:   "1.2.3",
	}
	m := d.ToDict()
	if m == nil {
		t.Fatal("ToDict returned nil")
	}
	if name, _ := m["name"].(string); name != "myserver" {
		t.Errorf("ToDict name = %q, want myserver", name)
	}
}

func TestFromDict_BasicFields_Extra2(t *testing.T) {
	m := map[string]interface{}{
		"name":    "myserver",
		"version": "2.0.0",
	}
	d, err := FromDict(m)
	if err != nil {
		t.Fatalf("FromDict error: %v", err)
	}
	if d.Name != "myserver" {
		t.Errorf("Name = %q", d.Name)
	}
}

func TestFromDict_MissingName_Extra2(t *testing.T) {
	m := map[string]interface{}{"version": "1.0.0"}
	_, err := FromDict(m)
	if err == nil {
		t.Error("expected error when name is missing")
	}
}

func TestIsRegistryResolved_Default_Extra2(t *testing.T) {
	d := MCPDependency{Name: "myserver"}
	// Default registry (nil) is considered resolved
	_ = d.IsRegistryResolved()
}

func TestIsSelfDefined_RegistryFalse_Extra2(t *testing.T) {
	d := MCPDependency{Name: "myserver", Registry: RegistryFalse}
	if !d.IsSelfDefined() {
		t.Error("expected IsSelfDefined=true when Registry=RegistryFalse")
	}
}

func TestIsSelfDefined_DefaultRegistry_Extra2(t *testing.T) {
	d := MCPDependency{Name: "myserver"}
	if d.IsSelfDefined() {
		t.Error("expected IsSelfDefined=false for default registry")
	}
}

func TestString_Format_Extra2(t *testing.T) {
	d := MCPDependency{Name: "myserver", Transport: "stdio"}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty String()")
	}
}

func TestValidate_ValidName_Extra2(t *testing.T) {
	d := MCPDependency{Name: "myserver"}
	if err := d.Validate(false); err != nil {
		t.Errorf("expected no error for valid name, got: %v", err)
	}
}

func TestValidate_EmptyName_Extra2(t *testing.T) {
	d := MCPDependency{Name: ""}
	if err := d.Validate(false); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestMCPDependency_Env_Extra2(t *testing.T) {
	d := MCPDependency{
		Name: "myserver",
		Env:  map[string]string{"TOKEN": "abc", "DEBUG": "1"},
	}
	if len(d.Env) != 2 {
		t.Errorf("Env len = %d, want 2", len(d.Env))
	}
}

func TestMCPDependency_Tools_Extra2(t *testing.T) {
	d := MCPDependency{
		Name:  "myserver",
		Tools: []string{"search", "index", "query"},
	}
	if len(d.Tools) != 3 {
		t.Errorf("Tools len = %d, want 3", len(d.Tools))
	}
}
