package mcpdep

import (
	"testing"
)

func TestIsValidName_ValidNames(t *testing.T) {
	valid := []string{
		"my-server",
		"io.github.acme/cool-server",
		"server123",
		"@scope/name",
		"_underscore",
		"a",
		"A",
	}
	for _, name := range valid {
		if !isValidName(name) {
			t.Errorf("isValidName(%q) should be true", name)
		}
	}
}

func TestIsValidName_InvalidNames(t *testing.T) {
	invalid := []string{
		"",
		"has space",
		"has!exclamation",
		"has#hash",
	}
	for _, name := range invalid {
		if isValidName(name) {
			t.Errorf("isValidName(%q) should be false", name)
		}
	}
}

func TestIsValidName_TooLong(t *testing.T) {
	long := ""
	for i := 0; i < 129; i++ {
		long += "a"
	}
	if isValidName(long) {
		t.Error("name > 128 chars should be invalid")
	}
	// exactly 128 is valid
	exact := ""
	for i := 0; i < 128; i++ {
		exact += "a"
	}
	if !isValidName(exact) {
		t.Error("128-char name should be valid")
	}
}

func TestValidate_EmptyName(t *testing.T) {
	d := &MCPDependency{Name: ""}
	if err := d.Validate(false); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestValidate_InvalidName(t *testing.T) {
	d := &MCPDependency{Name: "has space"}
	if err := d.Validate(false); err == nil {
		t.Error("expected error for invalid name")
	}
}

func TestValidate_DotDotSegment(t *testing.T) {
	d := &MCPDependency{Name: "foo/../bar"}
	if err := d.Validate(false); err == nil {
		t.Error("expected error for .. in name")
	}
}

func TestValidate_ValidName(t *testing.T) {
	d := &MCPDependency{Name: "my-server"}
	if err := d.Validate(false); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestString_BasicFormat(t *testing.T) {
	d := &MCPDependency{Name: "my-server", Transport: "stdio"}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty String()")
	}
}

func TestToDict_BasicFields(t *testing.T) {
	d := &MCPDependency{
		Name:      "my-server",
		Transport: "sse",
		URL:       "https://example.com",
		Command:   "npx",
	}
	dict := d.ToDict()
	if dict["name"] != "my-server" {
		t.Errorf("expected name my-server, got %v", dict["name"])
	}
}

func TestToDict_WithEnv(t *testing.T) {
	d := &MCPDependency{
		Name: "env-server",
		Env:  map[string]string{"KEY": "value"},
	}
	dict := d.ToDict()
	if dict["name"] != "env-server" {
		t.Errorf("unexpected name: %v", dict["name"])
	}
	if dict["env"] == nil {
		t.Error("expected env in dict")
	}
}

func TestFromDict_NilArgs(t *testing.T) {
	m := map[string]interface{}{
		"name": "simple",
	}
	d, err := FromDict(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name != "simple" {
		t.Errorf("expected name simple, got %q", d.Name)
	}
}

func TestFromDict_WithTransport(t *testing.T) {
	m := map[string]interface{}{
		"name":      "sse-server",
		"transport": "sse",
		"url":       "https://sse.example.com",
	}
	d, err := FromDict(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Transport != "sse" {
		t.Errorf("expected transport sse, got %q", d.Transport)
	}
	if d.URL != "https://sse.example.com" {
		t.Errorf("expected url, got %q", d.URL)
	}
}

func TestFromString_NameOnly(t *testing.T) {
	d, err := FromString("my-pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name != "my-pkg" {
		t.Errorf("expected name my-pkg, got %q", d.Name)
	}
}

func TestFromString_WithOrg(t *testing.T) {
	d, err := FromString("acme/my-server")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name == "" {
		t.Error("expected non-empty name")
	}
}

func TestMCPDependency_Fields(t *testing.T) {
	d := MCPDependency{
		Name:      "test",
		Transport: "stdio",
		Command:   "npx",
		URL:       "http://localhost:3000",
	}
	if d.Name != "test" || d.Transport != "stdio" || d.Command != "npx" {
		t.Errorf("unexpected fields: %+v", d)
	}
}
