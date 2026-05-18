package pluginparser

import (
"encoding/json"
"os"
"path/filepath"
"strings"
"testing"
)

func TestParsePluginManifestMissing(t *testing.T) {
_, err := ParsePluginManifest("/nonexistent/plugin.json")
if err == nil {
t.Error("expected error for missing file")
}
}

func TestParsePluginManifestMinimal(t *testing.T) {
dir := t.TempDir()
pluginJSON := filepath.Join(dir, "plugin.json")
data := map[string]interface{}{"name": "my-plugin"}
b, _ := json.Marshal(data)
if err := os.WriteFile(pluginJSON, b, 0o644); err != nil {
t.Fatal(err)
}
manifest, err := ParsePluginManifest(pluginJSON)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if manifest.Name != "my-plugin" {
t.Errorf("unexpected name: %s", manifest.Name)
}
}

func TestParsePluginManifestInvalidJSON(t *testing.T) {
dir := t.TempDir()
pluginJSON := filepath.Join(dir, "plugin.json")
if err := os.WriteFile(pluginJSON, []byte("{invalid}"), 0o644); err != nil {
t.Fatal(err)
}
_, err := ParsePluginManifest(pluginJSON)
if err == nil {
t.Error("expected error for invalid JSON")
}
}

func TestYamlString(t *testing.T) {
cases := []struct {
input    string
contains string
}{
{"simple", "simple"},
{"with space", "with space"},
{"with: colon", ":"},
}
for _, tc := range cases {
got := yamlString(tc.input)
if !strings.Contains(got, tc.contains) {
t.Errorf("yamlString(%q) = %q, expected to contain %q", tc.input, got, tc.contains)
}
}
}

func TestSynthesizeApmYMLFromPluginMinimal(t *testing.T) {
dir := t.TempDir()
manifest := &PluginManifest{Name: "test-plugin"}
apmYMLPath, err := SynthesizeApmYMLFromPlugin(dir, manifest)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if apmYMLPath == "" {
t.Error("expected non-empty apm.yml path")
}
content, err := os.ReadFile(apmYMLPath)
if err != nil {
t.Fatalf("could not read generated apm.yml: %v", err)
}
if !strings.Contains(string(content), "test-plugin") {
t.Errorf("generated apm.yml doesn't contain plugin name: %s", string(content))
}
}

func TestSynthesizeApmYMLFromPluginDefaultsName(t *testing.T) {
dir := t.TempDir()
manifest := &PluginManifest{} // no name — should default to dir basename
_, err := SynthesizeApmYMLFromPlugin(dir, manifest)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if manifest.Name == "" {
t.Error("manifest name should have been set to dir basename")
}
}
