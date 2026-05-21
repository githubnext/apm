package marketplace

import "testing"

func TestIsValidAlias_Dot_Extra4(t *testing.T) {
if !IsValidAlias("a.b") {
t.Error("expected 'a.b' to be valid")
}
}

func TestIsValidAlias_Underscore_Extra4(t *testing.T) {
if !IsValidAlias("my_alias") {
t.Error("expected 'my_alias' to be valid")
}
}

func TestIsValidAlias_Hyphen_Extra4(t *testing.T) {
if !IsValidAlias("my-alias") {
t.Error("expected 'my-alias' to be valid")
}
}

func TestIsValidAlias_Long_Extra4(t *testing.T) {
alias := "abcdefghijklmnopqrstuvwxyz0123456789"
if !IsValidAlias(alias) {
t.Errorf("expected long alphanumeric alias to be valid")
}
}

func TestIsValidAlias_RejectsSpace_Extra4(t *testing.T) {
if IsValidAlias("my alias") {
t.Error("expected 'my alias' to be invalid")
}
}

func TestIsValidAlias_RejectsBracket_Extra4(t *testing.T) {
if IsValidAlias("my[alias]") {
t.Error("expected bracket alias to be invalid")
}
}

func TestMarketplaceConfig_AliasField_Extra4(t *testing.T) {
cfg := MarketplaceConfig{
Alias: "test",
URL:   "https://example.com",
}
if cfg.Alias != "test" {
t.Errorf("expected 'test', got %q", cfg.Alias)
}
}

func TestMarketplaceConfig_URLField_Extra4(t *testing.T) {
cfg := MarketplaceConfig{URL: "https://example.com/pkg"}
if cfg.URL != "https://example.com/pkg" {
t.Errorf("expected URL field, got %q", cfg.URL)
}
}

func TestAddOptions_FieldAccess_Extra4(t *testing.T) {
opts := AddOptions{
ProjectRoot: "/tmp/proj",
Alias:       "mypkg",
URL:         "https://example.com",
}
if opts.ProjectRoot != "/tmp/proj" {
t.Errorf("unexpected ProjectRoot: %q", opts.ProjectRoot)
}
if opts.Alias != "mypkg" {
t.Errorf("unexpected Alias: %q", opts.Alias)
}
}

func TestAddOptions_ZeroValue_Extra4(t *testing.T) {
var opts AddOptions
if opts.ProjectRoot != "" || opts.Alias != "" {
t.Error("zero AddOptions should have empty fields")
}
}
