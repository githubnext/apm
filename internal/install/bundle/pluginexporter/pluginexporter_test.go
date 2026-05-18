package pluginexporter

import (
"os"
"path/filepath"
"testing"
)

func TestValidateOutputRel(t *testing.T) {
cases := []struct {
rel  string
want bool
}{
{"agents/my-agent.md", true},
{"skills/foo.md", true},
{"plugin.json", true},
{"/absolute/path", false},
{"../escape", false},
{"a/../../b", false},
{"", true},
}
for _, c := range cases {
got := validateOutputRel(c.rel)
if got != c.want {
t.Errorf("validateOutputRel(%q): want %v, got %v", c.rel, c.want, got)
}
}
}

func TestSanitizeBundleName(t *testing.T) {
cases := []struct {
input string
want  string
}{
{"my-bundle", "my-bundle"},
{"my bundle", "my-bundle"},
{"hello/world", "hello/world"},
{"a!b@c", "a-b-c"},
{"---", "unnamed"},
{"", "unnamed"},
{"valid.name", "valid.name"},
}
for _, c := range cases {
got := sanitizeBundleName(c.input)
if got != c.want {
t.Errorf("sanitizeBundleName(%q): want %q, got %q", c.input, c.want, got)
}
}
}

func TestRenamePrompt(t *testing.T) {
cases := []struct {
input string
want  string
}{
{"foo.prompt.md", "foo.md"},
{"bar.md", "bar.md"},
{"readme.prompt.md", "readme.md"},
{"no-extension", "no-extension"},
}
for _, c := range cases {
got := renamePrompt(c.input)
if got != c.want {
t.Errorf("renamePrompt(%q): want %q, got %q", c.input, c.want, got)
}
}
}

func TestExportPluginBundleDryRun(t *testing.T) {
dir := t.TempDir()
// Create minimal .apm/agents structure
apmDir := filepath.Join(dir, ".apm", "agents")
if err := os.MkdirAll(apmDir, 0o755); err != nil {
t.Fatal(err)
}
if err := os.WriteFile(filepath.Join(apmDir, "my-agent.md"), []byte("# My Agent\n"), 0o644); err != nil {
t.Fatal(err)
}

outDir := t.TempDir()
opts := ExportOptions{
ProjectRoot: dir,
OutputDir:   outDir,
DryRun:      true,
}
result, err := ExportPluginBundle(opts)
if err != nil {
t.Fatalf("ExportPluginBundle dry run: %v", err)
}
if result == nil {
t.Fatal("ExportPluginBundle returned nil result")
}
}
