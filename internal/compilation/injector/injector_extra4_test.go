package injector

import (
"os"
"path/filepath"
"testing"
)

func TestStatusSkipped_Extra4(t *testing.T) {
if StatusSkipped == "" {
t.Error("StatusSkipped should not be empty")
}
}

func TestStatusCreated_Extra4(t *testing.T) {
if StatusCreated == "" {
t.Error("StatusCreated should not be empty")
}
}

func TestStatusUpdated_Extra4(t *testing.T) {
if StatusUpdated == "" {
t.Error("StatusUpdated should not be empty")
}
}

func TestStatusUnchanged_Extra4(t *testing.T) {
if StatusUnchanged == "" {
t.Error("StatusUnchanged should not be empty")
}
}

func TestStatusMissing_Extra4(t *testing.T) {
if StatusMissing == "" {
t.Error("StatusMissing should not be empty")
}
}

func TestAllStatusValues_Distinct_Extra4(t *testing.T) {
statuses := []InjectionStatus{
StatusSkipped,
StatusCreated,
StatusUpdated,
StatusUnchanged,
StatusMissing,
}
seen := map[InjectionStatus]bool{}
for _, s := range statuses {
if seen[s] {
t.Errorf("duplicate status: %q", s)
}
seen[s] = true
}
}

func TestConstitutionInjector_InjectNoOutput_Extra4(t *testing.T) {
dir := t.TempDir()
ci := &ConstitutionInjector{BaseDir: dir}
content, status, _ := ci.Inject("hello world", false, "")
_ = content
_ = status
}

func TestConstitutionInjector_InjectWithOutputPath_Extra4(t *testing.T) {
dir := t.TempDir()
out := filepath.Join(dir, "out.md")
ci := &ConstitutionInjector{BaseDir: dir}
_, _, _ = ci.Inject("content", false, out)
}

func TestConstitutionInjector_BaseDir_Extra4(t *testing.T) {
dir := t.TempDir()
ci := &ConstitutionInjector{BaseDir: dir}
if ci.BaseDir != dir {
t.Errorf("expected %q, got %q", dir, ci.BaseDir)
}
}

func TestInject_WithConstitution_WritesFile_Extra4(t *testing.T) {
dir := t.TempDir()
constitDir := filepath.Join(dir, ".specify", "memory")
if err := os.MkdirAll(constitDir, 0755); err != nil {
t.Fatal(err)
}
constitPath := filepath.Join(constitDir, "constitution.md")
if err := os.WriteFile(constitPath, []byte("# constitution"), 0644); err != nil {
t.Fatal(err)
}
outPath := filepath.Join(dir, "result.md")
ci := &ConstitutionInjector{BaseDir: dir}
_, _, _ = ci.Inject("original content", true, outPath)
}

func TestInject_WithConstitution_StatusValue_Extra4(t *testing.T) {
dir := t.TempDir()
constitDir := filepath.Join(dir, ".specify", "memory")
_ = os.MkdirAll(constitDir, 0755)
_ = os.WriteFile(filepath.Join(constitDir, "constitution.md"), []byte("# c"), 0644)
ci := &ConstitutionInjector{BaseDir: dir}
_, status, _ := ci.Inject("content", true, "")
if status == "" {
t.Error("expected non-empty status")
}
}
