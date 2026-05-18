package configcmd

import (
"testing"
)

func TestParseBoolValue_CaseInsensitive(t *testing.T) {
trueVals := []string{"TRUE", "True", "TrUe", "YES", "Yes", "1"}
for _, v := range trueVals {
got, err := ParseBoolValue(v)
if err != nil {
t.Errorf("ParseBoolValue(%q) unexpected error: %v", v, err)
}
if !got {
t.Errorf("ParseBoolValue(%q) = false, want true", v)
}
}
}

func TestParseBoolValue_FalseCaseInsensitive(t *testing.T) {
falseVals := []string{"FALSE", "False", "FaLsE", "NO", "No", "0"}
for _, v := range falseVals {
got, err := ParseBoolValue(v)
if err != nil {
t.Errorf("ParseBoolValue(%q) unexpected error: %v", v, err)
}
if got {
t.Errorf("ParseBoolValue(%q) = true, want false", v)
}
}
}

func TestParseBoolValue_InvalidValues(t *testing.T) {
invalid := []string{"on", "off", "enabled", "disabled", "t", "f", "y", "n", "2", "-1", "  "}
for _, v := range invalid {
_, err := ParseBoolValue(v)
if err == nil {
t.Errorf("ParseBoolValue(%q) expected error, got nil", v)
}
}
}

func TestValidConfigKeys_ContainsKnownKeys(t *testing.T) {
keys := ValidConfigKeys()
knownKeys := []string{"auto-integrate", "temp-dir"}
keySet := make(map[string]bool, len(keys))
for _, k := range keys {
keySet[k] = true
}
for _, k := range knownKeys {
if !keySet[k] {
t.Errorf("ValidConfigKeys missing expected key %q", k)
}
}
}

func TestDisplayName_AutoIntegrate(t *testing.T) {
name := DisplayName("auto_integrate")
if name != "auto-integrate" {
t.Errorf("DisplayName(auto_integrate) = %q, want auto-integrate", name)
}
}

func TestDisplayName_TempDir(t *testing.T) {
name := DisplayName("temp_dir")
if name != "temp-dir" {
t.Errorf("DisplayName(temp_dir) = %q, want temp-dir", name)
}
}

func TestDisplayName_UnknownFallback(t *testing.T) {
name := DisplayName("unknown_key")
// Should return a non-empty fallback (the raw key or similar).
if name == "" {
t.Error("DisplayName for unknown key should return non-empty fallback")
}
}

func TestParseAPMYML_WithVersion(t *testing.T) {
content := "name: myapp\nversion: 2.0.0\n"
cfg := parseAPMYML(content)
if cfg.Version != "2.0.0" {
t.Errorf("Version = %q, want 2.0.0", cfg.Version)
}
}

func TestParseAPMYML_WithName(t *testing.T) {
content := "name: testapp\n"
cfg := parseAPMYML(content)
if cfg.Name != "testapp" {
t.Errorf("Name = %q, want testapp", cfg.Name)
}
}

func TestParseAPMYML_NoNameVersionEmpty(t *testing.T) {
cfg := parseAPMYML("description: just a description\n")
if cfg.Name != "" {
t.Errorf("Name should be empty when absent, got %q", cfg.Name)
}
if cfg.Version != "" {
t.Errorf("Version should be empty when absent, got %q", cfg.Version)
}
}

func TestParseAPMYML_MultipleFields(t *testing.T) {
content := "name: full-app\nversion: 3.1.4\nentrypoint: main.go\n"
cfg := parseAPMYML(content)
if cfg.Name != "full-app" {
t.Errorf("Name = %q, want full-app", cfg.Name)
}
if cfg.Version != "3.1.4" {
t.Errorf("Version = %q, want 3.1.4", cfg.Version)
}
}
