package apmpackage_test

import (
"os"
"path/filepath"
"testing"

"github.com/githubnext/apm/internal/models/apmpackage"
)

func TestParseContentType_Valid(t *testing.T) {
cases := []struct {
input string
want  apmpackage.PackageContentType
}{
{"instructions", apmpackage.ContentTypeInstructions},
{"skill", apmpackage.ContentTypeSkill},
{"hybrid", apmpackage.ContentTypeHybrid},
{"prompts", apmpackage.ContentTypePrompts},
{"SKILL", apmpackage.ContentTypeSkill},
}
for _, tc := range cases {
got, err := apmpackage.ParseContentType(tc.input)
if err != nil {
t.Errorf("ParseContentType(%q): unexpected error %v", tc.input, err)
}
if got != tc.want {
t.Errorf("ParseContentType(%q): got %v want %v", tc.input, got, tc.want)
}
}
}

func TestParseContentType_Invalid(t *testing.T) {
_, err := apmpackage.ParseContentType("unknown-type")
if err == nil {
t.Error("expected error for unknown content type")
}
}

func TestContentTypeString(t *testing.T) {
cases := []struct {
ct   apmpackage.PackageContentType
want string
}{
{apmpackage.ContentTypeInstructions, "instructions"},
{apmpackage.ContentTypeSkill, "skill"},
{apmpackage.ContentTypeHybrid, "hybrid"},
{apmpackage.ContentTypePrompts, "prompts"},
}
for _, tc := range cases {
if tc.ct.String() != tc.want {
t.Errorf("ContentType.String(): got %q want %q", tc.ct.String(), tc.want)
}
}
}

func TestPackageInfo_HasPrimitives_WithFiles(t *testing.T) {
dir := t.TempDir()
instDir := filepath.Join(dir, "instructions")
if err := os.MkdirAll(instDir, 0o755); err != nil {
t.Fatal(err)
}
apmDir := filepath.Join(dir, ".apm", "instructions")
if err := os.MkdirAll(apmDir, 0o755); err != nil {
t.Fatal(err)
}
f, err := os.Create(filepath.Join(apmDir, "test.md"))
if err != nil {
t.Fatal(err)
}
f.Close()
info := &apmpackage.PackageInfo{InstallPath: dir}
if !info.HasPrimitives() {
t.Error("expected HasPrimitives()=true when .apm/instructions has files")
}
}

func TestPackageInfo_HasPrimitives_Empty(t *testing.T) {
dir := t.TempDir()
info := &apmpackage.PackageInfo{InstallPath: dir}
if info.HasPrimitives() {
t.Error("expected HasPrimitives()=false for empty install dir")
}
}
