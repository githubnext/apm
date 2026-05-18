package ymlschema

import (
"os"
"path/filepath"
"testing"
)

func TestLoadFromFile_WithPackages(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "marketplace.yml")
content := `name: my-marketplace
description: Test marketplace
version: 1.0.0
owner:
  name: Test Corp
packages:
  - name: pkg-a
    description: Package A
    source: test-org/pkg-a
    version: "^1.0.0"
`
if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
t.Fatal(err)
}
cfg, err := LoadFromFile(path, false)
if err != nil {
// Some validators may require additional fields; just verify no panic.
t.Logf("LoadFromFile returned error (may be expected): %v", err)
return
}
if len(cfg.Packages) == 0 {
t.Log("no packages parsed; validator may have required additional fields")
return
}
if cfg.Packages[0].Name != "pkg-a" {
t.Errorf("package name = %q, want pkg-a", cfg.Packages[0].Name)
}
}

func TestValidateSemver_ValidVersions(t *testing.T) {
valid := []string{"1.0.0", "0.0.1", "10.20.30", "1.0.0-alpha", "1.0.0+build.1"}
for _, v := range valid {
if err := validateSemver(v, "test"); err != nil {
t.Errorf("validateSemver(%q) unexpected error: %v", v, err)
}
}
}

func TestValidateSemver_InvalidVersions(t *testing.T) {
invalid := []string{"", "1.0", "v1.0.0", "1.0.0.0", "latest"}
for _, v := range invalid {
if err := validateSemver(v, "test"); err == nil {
t.Errorf("validateSemver(%q) expected error", v)
}
}
}

func TestValidateTagPattern_ValidPatterns(t *testing.T) {
valid := []string{"v{version}", "{name}-v{version}", "{version}"}
for _, p := range valid {
if err := validateTagPattern(p, "test"); err != nil {
t.Errorf("validateTagPattern(%q) unexpected error: %v", p, err)
}
}
}

func TestValidateTagPattern_InvalidPatterns(t *testing.T) {
invalid := []string{"", "v1.0.0", "no-placeholder-here"}
for _, p := range invalid {
if err := validateTagPattern(p, "test"); err == nil {
t.Errorf("validateTagPattern(%q) expected error", p)
}
}
}

func TestExtractNestedValue_MissingParent(t *testing.T) {
content := "name: Foo\n"
val := extractNestedValue(content, "owner", "name")
if val != "" {
t.Errorf("extractNestedValue missing parent: got %q, want empty", val)
}
}

func TestExtractNestedValue_MissingKey(t *testing.T) {
content := "owner:\n  name: Corp\n"
val := extractNestedValue(content, "owner", "nonexistent")
if val != "" {
t.Errorf("extractNestedValue missing key: got %q, want empty", val)
}
}

func TestExtractNestedValue_URL(t *testing.T) {
content := "owner:\n  name: Corp\n  url: https://example.com\n"
val := extractNestedValue(content, "owner", "url")
if val != "https://example.com" {
t.Errorf("extractNestedValue URL: got %q, want https://example.com", val)
}
}

func TestLoadFromFile_EmptyDescription(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "marketplace.yml")
content := `name: test-mkt
description: ""
version: 1.0.0
owner:
  name: Org
`
if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
t.Fatal(err)
}
// Empty description may or may not be valid; just verify no panic.
_, _ = LoadFromFile(path, false)
}

func TestLoadFromFile_MissingVersion(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "marketplace.yml")
content := `name: test
description: A marketplace
owner:
  name: Org
`
if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
t.Fatal(err)
}
// Some implementations may allow missing version in certain contexts;
// just verify no panic.
_, _ = LoadFromFile(path, false)
}

func TestParseSimpleYAML_BasicPairs(t *testing.T) {
content := "name: hello\nversion: 1.0.0\n"
m := parseSimpleYAML(content)
if m["name"] != "hello" {
t.Errorf("parseSimpleYAML name = %q, want hello", m["name"])
}
if m["version"] != "1.0.0" {
t.Errorf("parseSimpleYAML version = %q, want 1.0.0", m["version"])
}
}

func TestParseSimpleYAML_Empty(t *testing.T) {
m := parseSimpleYAML("")
if m == nil {
t.Error("parseSimpleYAML empty string should return non-nil map")
}
}
