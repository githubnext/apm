package targetdetection

import (
"os"
"path/filepath"
"testing"
)

func TestNormalizeTarget_All(t *testing.T) {
cases := []struct{ in, want string }{
{"copilot", "vscode"},
{"vscode", "vscode"},
{"agents", "vscode"},
{"claude", "claude"},
{"cursor", "cursor"},
{"codex", "codex"},
{"gemini", "gemini"},
{"opencode", "opencode"},
{"windsurf", "windsurf"},
{"all", "all"},
{"minimal", "minimal"},
{"unknown", "unknown"},
}
for _, c := range cases {
if got := NormalizeTarget(c.in); got != c.want {
t.Errorf("NormalizeTarget(%q) = %q, want %q", c.in, got, c.want)
}
}
}

func TestValidTargets_Contents(t *testing.T) {
for _, name := range []string{"vscode", "claude", "cursor", "codex", "gemini", "opencode", "windsurf", "all", "minimal"} {
if !ValidTargets[name] {
t.Errorf("expected %q in ValidTargets", name)
}
}
}

func TestCanonicalTargetsOrdered_Length(t *testing.T) {
if len(CanonicalTargetsOrdered) == 0 {
t.Fatal("CanonicalTargetsOrdered must not be empty")
}
}

func TestCanonicalDeployDirs_Coverage(t *testing.T) {
for _, name := range CanonicalTargetsOrdered {
if _, ok := CanonicalDeployDirs[name]; !ok {
t.Errorf("CanonicalDeployDirs missing entry for %q", name)
}
}
}

func TestDetectSignals_EmptyDir(t *testing.T) {
dir := t.TempDir()
sigs := DetectSignals(dir)
if len(sigs) != 0 {
t.Errorf("expected no signals in empty dir, got %v", sigs)
}
}

func TestDetectSignals_CopilotFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".github"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".github", "copilot-instructions.md"), []byte("# CI"), 0644); err != nil {
		t.Fatal(err)
	}
	sigs := DetectSignals(dir)
	found := false
	for _, s := range sigs {
		if s.Target == "copilot" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected copilot signal for copilot-instructions.md, got %v", sigs)
	}
}

func TestDetectSignals_ClaudeFile(t *testing.T) {
dir := t.TempDir()
if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Claude"), 0644); err != nil {
t.Fatal(err)
}
sigs := DetectSignals(dir)
found := false
for _, s := range sigs {
if s.Target == "claude" {
found = true
}
}
if !found {
t.Errorf("expected claude signal for CLAUDE.md, got %v", sigs)
}
}

func TestDetectSignals_CursorDir(t *testing.T) {
dir := t.TempDir()
if err := os.MkdirAll(filepath.Join(dir, ".cursor"), 0755); err != nil {
t.Fatal(err)
}
sigs := DetectSignals(dir)
found := false
for _, s := range sigs {
if s.Target == "cursor" {
found = true
}
}
if !found {
t.Errorf("expected cursor signal for .cursor dir, got %v", sigs)
}
}

func TestResolveTargets_FlagOverride(t *testing.T) {
dir := t.TempDir()
res, err := ResolveTargets(dir, []string{"claude"}, []string{"cursor"})
if err != nil {
t.Fatal(err)
}
if len(res.Targets) != 1 || res.Targets[0] != "claude" {
t.Errorf("flag should override yaml targets; got %v", res.Targets)
}
}

func TestResolveTargets_YAMLTargets(t *testing.T) {
dir := t.TempDir()
res, err := ResolveTargets(dir, nil, []string{"cursor", "claude"})
if err != nil {
t.Fatal(err)
}
if len(res.Targets) != 2 {
t.Errorf("expected 2 targets, got %v", res.Targets)
}
if res.Source != "apm.yml" {
t.Errorf("expected source apm.yml, got %q", res.Source)
}
}

func TestResolveTargets_UnknownFlag(t *testing.T) {
dir := t.TempDir()
_, err := ResolveTargets(dir, []string{"unknown-tool"}, nil)
if err == nil {
t.Fatal("expected error for unknown target")
}
}

func TestResolveTargets_NoHarness(t *testing.T) {
dir := t.TempDir()
_, err := ResolveTargets(dir, nil, nil)
if err == nil {
t.Fatal("expected error when no harness found")
}
}

func TestResolveTargets_AutoDetect(t *testing.T) {
dir := t.TempDir()
if err := os.MkdirAll(filepath.Join(dir, ".cursor"), 0755); err != nil {
t.Fatal(err)
}
res, err := ResolveTargets(dir, nil, nil)
if err != nil {
t.Fatal(err)
}
if len(res.Targets) == 0 {
t.Fatal("expected at least one target")
}
}

func TestResolveTargets_DedupFlag(t *testing.T) {
dir := t.TempDir()
res, err := ResolveTargets(dir, []string{"claude", "claude", "cursor"}, nil)
if err != nil {
t.Fatal(err)
}
if len(res.Targets) != 2 {
t.Errorf("expected deduped targets, got %v", res.Targets)
}
}

func TestExpandAllTargets_NoHarness(t *testing.T) {
dir := t.TempDir()
_, err := ExpandAllTargets(dir, nil)
if err == nil {
t.Fatal("expected error for empty dir with no yaml targets")
}
}

func TestExpandAllTargets_WithYAML(t *testing.T) {
dir := t.TempDir()
targets, err := ExpandAllTargets(dir, []string{"claude", "cursor"})
if err != nil {
t.Fatal(err)
}
if len(targets) != 2 {
t.Errorf("expected 2 targets, got %v", targets)
}
}

func TestExpandAllTargets_Dedup(t *testing.T) {
dir := t.TempDir()
if err := os.MkdirAll(filepath.Join(dir, ".cursor"), 0755); err != nil {
t.Fatal(err)
}
targets, err := ExpandAllTargets(dir, []string{"cursor"})
if err != nil {
t.Fatal(err)
}
// cursor should appear only once even if from both signal and yaml
count := 0
for _, t2 := range targets {
if t2 == "cursor" {
count++
}
}
if count != 1 {
t.Errorf("expected cursor once, got %d times in %v", count, targets)
}
}

func TestFormatProvenance_Single(t *testing.T) {
r := ResolvedTargets{Targets: []string{"claude"}, Source: "apm.yml"}
got := FormatProvenance(r)
if got != "Targets: claude  (source: apm.yml)" {
t.Errorf("unexpected provenance: %q", got)
}
}

func TestFormatProvenance_Empty(t *testing.T) {
r := ResolvedTargets{Targets: []string{}, Source: "manual"}
got := FormatProvenance(r)
if got == "" {
t.Error("expected non-empty provenance string")
}
}

func TestDetectTarget_ConfigTarget(t *testing.T) {
target, reason := DetectTarget("/nonexistent", "", "cursor")
if target != "cursor" {
t.Errorf("expected cursor, got %q", target)
}
if reason != "apm.yml target" {
t.Errorf("unexpected reason: %q", reason)
}
}

func TestDetectTarget_NoFolders(t *testing.T) {
dir := t.TempDir()
target, reason := DetectTarget(dir, "", "")
if target != "minimal" {
t.Errorf("expected minimal, got %q", target)
}
if reason == "" {
t.Error("expected non-empty reason")
}
}

func TestDetectTarget_GithubFolder(t *testing.T) {
dir := t.TempDir()
if err := os.MkdirAll(filepath.Join(dir, ".github"), 0755); err != nil {
t.Fatal(err)
}
target, _ := DetectTarget(dir, "", "")
if target != "vscode" {
t.Errorf("expected vscode for .github folder, got %q", target)
}
}

func TestDetectTarget_MultipleDetected(t *testing.T) {
dir := t.TempDir()
for _, sub := range []string{".github", ".claude"} {
if err := os.MkdirAll(filepath.Join(dir, sub), 0755); err != nil {
t.Fatal(err)
}
}
target, reason := DetectTarget(dir, "", "")
if target != "all" {
t.Errorf("expected 'all' for multiple folders, got %q", target)
}
if reason == "" {
t.Error("expected non-empty reason")
}
}

func TestSignal_Fields(t *testing.T) {
s := Signal{Target: "claude", Source: "CLAUDE.md"}
if s.Target != "claude" || s.Source != "CLAUDE.md" {
t.Errorf("unexpected signal fields: %+v", s)
}
}

func TestResolvedTargets_AutoCreate(t *testing.T) {
dir := t.TempDir()
res, err := ResolveTargets(dir, []string{"claude"}, nil)
if err != nil {
t.Fatal(err)
}
if !res.AutoCreate {
t.Error("expected AutoCreate to be true when flag provided")
}
}
