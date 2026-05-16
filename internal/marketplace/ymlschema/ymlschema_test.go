package ymlschema

import (
"os"
"path/filepath"
"testing"
)

func TestLoadFromFileMissing(t *testing.T) {
_, err := LoadFromFile("/nonexistent/path/marketplace.yml", false)
if err == nil {
t.Error("expected error for missing file")
}
}

func TestLoadFromFileMinimal(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "marketplace.yml")
content := `name: my-marketplace
description: A test marketplace
version: 1.0.0
owner:
  name: Acme Corp
`
if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
t.Fatal(err)
}
cfg, err := LoadFromFile(path, false)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if cfg.Name != "my-marketplace" {
t.Errorf("unexpected name: %s", cfg.Name)
}
if cfg.Version != "1.0.0" {
t.Errorf("unexpected version: %s", cfg.Version)
}
if cfg.Owner.Name != "Acme Corp" {
t.Errorf("unexpected owner name: %s", cfg.Owner.Name)
}
}

func TestLoadFromFileMissingOwner(t *testing.T) {
dir := t.TempDir()
path := filepath.Join(dir, "marketplace.yml")
content := `name: my-marketplace
description: A test marketplace
version: 1.0.0
`
if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
t.Fatal(err)
}
_, err := LoadFromFile(path, false)
if err == nil {
t.Error("expected error for missing owner")
}
}

func TestValidateSemver(t *testing.T) {
cases := []struct {
version string
valid   bool
}{
{"1.0.0", true},
{"2.3.4", true},
{"1.0.0-alpha.1", true},
{"not-semver", false},
{"1.2", false},
{"", false},
}
for _, tc := range cases {
err := validateSemver(tc.version, "test")
if tc.valid && err != nil {
t.Errorf("validateSemver(%q) unexpected error: %v", tc.version, err)
}
if !tc.valid && err == nil {
t.Errorf("validateSemver(%q) expected error", tc.version)
}
}
}

func TestValidateTagPattern(t *testing.T) {
if err := validateTagPattern("v{version}", "test"); err != nil {
t.Errorf("v{version} should be valid: %v", err)
}
if err := validateTagPattern("{name}-v1", "test"); err != nil {
t.Errorf("{name} pattern should be valid: %v", err)
}
if err := validateTagPattern("no-placeholder", "test"); err == nil {
t.Error("pattern without placeholder should be invalid")
}
}

func TestExtractNestedValue(t *testing.T) {
content := "owner:\n  name: Acme\n  email: test@example.com\n"
val := extractNestedValue(content, "owner", "name")
if val != "Acme" {
t.Errorf("expected 'Acme', got %q", val)
}
email := extractNestedValue(content, "owner", "email")
if email != "test@example.com" {
t.Errorf("expected email, got %q", email)
}
}
