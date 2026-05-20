package ymlschema

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageEntry_ZeroValue(t *testing.T) {
	var e PackageEntry
	if e.Name != "" || e.Source != "" || e.Version != "" {
		t.Error("zero value strings should be empty")
	}
	if e.IncludePrerelease {
		t.Error("IncludePrerelease should be false")
	}
	if e.Tags != nil {
		t.Error("Tags should be nil")
	}
}

func TestMarketplaceOwner_ZeroValue(t *testing.T) {
	var o MarketplaceOwner
	if o.Name != "" || o.Email != "" || o.URL != "" {
		t.Error("zero value strings should be empty")
	}
}

func TestMarketplaceBuild_ZeroValue(t *testing.T) {
	var b MarketplaceBuild
	if b.TagPattern != "" {
		t.Error("TagPattern should be empty")
	}
}

func TestMarketplaceConfig_ZeroValue(t *testing.T) {
	var c MarketplaceConfig
	if c.Name != "" || c.Description != "" {
		t.Error("zero value strings should be empty")
	}
	if c.Packages != nil {
		t.Error("Packages should be nil")
	}
}

func TestLoadFromFile_MinimalOwnerOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marketplace.yml")
	content := `name: test-marketplace
description: Test
version: 1.0.0
owner:
  name: Test Corp
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFromFile(path, false)
	if err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	if cfg.Name != "test-marketplace" {
		t.Errorf("Name: %q", cfg.Name)
	}
	if cfg.Owner.Name != "Test Corp" {
		t.Errorf("Owner.Name: %q", cfg.Owner.Name)
	}
}

func TestLoadFromFile_VersionParsed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marketplace.yml")
	content := `name: versioned
description: Versioned
version: 2.5.0
owner:
  name: Org
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFromFile(path, false)
	if err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	if cfg.Version != "2.5.0" {
		t.Errorf("Version: %q", cfg.Version)
	}
}

func TestLoadFromFile_WithMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marketplace.yml")
	content := `name: meta-marketplace
description: Metadata test
version: 1.0.0
owner:
  name: Test Corp
packages:
  - name: p1
    description: Package one
    source: org/p1
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFromFile(path, false)
	if err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	if cfg.Name != "meta-marketplace" {
		t.Errorf("Name: %q", cfg.Name)
	}
	if cfg.Version != "1.0.0" {
		t.Errorf("Version: %q", cfg.Version)
	}
}

func TestMarketplaceYmlError_Message(t *testing.T) {
	e := &MarketplaceYmlError{Msg: "test error message"}
	if e.Error() != "test error message" {
		t.Errorf("Error(): %q", e.Error())
	}
}

func TestMarketplaceYmlError_EmptyMsg(t *testing.T) {
	e := &MarketplaceYmlError{}
	if e.Error() != "" {
		t.Errorf("Error(): %q", e.Error())
	}
}

func TestParseSimpleYAML_MultilineValues(t *testing.T) {
	content := "key1: value1\nkey2: value2\nkey3: value3\n"
	result := parseSimpleYAML(content)
	if result["key1"] != "value1" {
		t.Errorf("key1: %q", result["key1"])
	}
	if result["key2"] != "value2" {
		t.Errorf("key2: %q", result["key2"])
	}
	if result["key3"] != "value3" {
		t.Errorf("key3: %q", result["key3"])
	}
}

func TestParseSimpleYAML_IndentedLine(t *testing.T) {
	content := "name: test\n  email: test@example.com\n"
	result := parseSimpleYAML(content)
	if _, ok := result["name"]; !ok {
		t.Error("name key should be present")
	}
}

func TestExtractNestedValue_DeepNesting(t *testing.T) {
	content := "owner:\n  name: Deep Corp\n  email: deep@example.com\n"
	name := extractNestedValue(content, "owner", "name")
	if name != "Deep Corp" {
		t.Errorf("expected 'Deep Corp', got %q", name)
	}
	email := extractNestedValue(content, "owner", "email")
	if email != "deep@example.com" {
		t.Errorf("expected 'deep@example.com', got %q", email)
	}
}

func TestLoadFromFile_MissingOwner_Returns_Error(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marketplace.yml")
	content := `name: test
description: No owner
version: 1.0.0
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadFromFile(path, false)
	if err == nil {
		t.Error("expected error for missing owner field")
	}
}

func TestLoadFromFile_LegacyMode_OutputDiffers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marketplace.yml")
	content := `name: legacy
description: Legacy mode
version: 1.0.0
owner:
  name: Corp
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfgLegacy, err := LoadFromFile(path, true)
	if err != nil {
		t.Fatalf("LoadFromFile legacy error: %v", err)
	}
	cfgNormal, err := LoadFromFile(path, false)
	if err != nil {
		t.Fatalf("LoadFromFile normal error: %v", err)
	}
	if cfgLegacy.Output == cfgNormal.Output {
		t.Error("legacy and normal output paths should differ")
	}
	if !cfgLegacy.IsLegacy {
		t.Error("IsLegacy should be true")
	}
	if cfgNormal.IsLegacy {
		t.Error("IsLegacy should be false for normal mode")
	}
}

func TestValidateTagPattern_WithVersion(t *testing.T) {
	err := validateTagPattern("v{version}", "test")
	if err != nil {
		t.Errorf("valid pattern should not error: %v", err)
	}
}

func TestValidateTagPattern_WithName(t *testing.T) {
	err := validateTagPattern("{name}-v1.0", "test")
	if err != nil {
		t.Errorf("valid pattern with {name} should not error: %v", err)
	}
}

func TestValidateSemver_PreRelease(t *testing.T) {
	err := validateSemver("1.0.0-alpha.1", "test")
	if err != nil {
		t.Errorf("pre-release semver should be valid: %v", err)
	}
}

func TestValidateSemver_BuildMetadata(t *testing.T) {
	err := validateSemver("1.0.0+build.1", "test")
	if err != nil {
		t.Errorf("build metadata semver should be valid: %v", err)
	}
}

func TestValidateSemver_TwoComponents(t *testing.T) {
	err := validateSemver("1.0", "test")
	if err == nil {
		t.Error("two-component version should be invalid")
	}
}

func TestLoadFromFile_OwnerFieldsParsed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marketplace.yml")
	content := `name: owner-test
description: Test
version: 1.0.0
owner:
  name: Acme Corp
  email: hello@acme.com
  url: https://acme.com
packages:
  - name: p1
    description: package one
    source: acme/p1
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFromFile(path, false)
	if err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	if cfg.Owner.Name != "Acme Corp" {
		t.Errorf("Owner.Name: %q", cfg.Owner.Name)
	}
	if !strings.Contains(cfg.Owner.Email, "acme.com") {
		t.Errorf("Owner.Email: %q", cfg.Owner.Email)
	}
}
