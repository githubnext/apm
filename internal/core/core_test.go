package core_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/core"
)

// ---------------------------------------------------------------------------
// Parity: errors
// ---------------------------------------------------------------------------

func TestParityRenderNoHarnessError(t *testing.T) {
	msg := core.RenderNoHarnessError()
	if !strings.Contains(msg, "[x] No harness detected") {
		t.Errorf("expected headline, got: %s", msg)
	}
	if !strings.Contains(msg, "apm install") {
		t.Error("expected actionable command in error message")
	}
}

func TestParityRenderAmbiguousError(t *testing.T) {
	msg := core.RenderAmbiguousError([]string{".github/", ".claude/"})
	if !strings.Contains(msg, "[x] Multiple harnesses detected") {
		t.Errorf("expected headline, got: %s", msg)
	}
	if !strings.Contains(msg, ".github/") {
		t.Error("expected detected folders in message")
	}
}

func TestParityRenderUnknownTargetError(t *testing.T) {
	msg := core.RenderUnknownTargetError("foo", []string{"claude", "copilot", "cursor"})
	if !strings.Contains(msg, "[x] Unknown target 'foo'") {
		t.Errorf("unexpected message: %s", msg)
	}
	if strings.Contains(msg, "agent-skills") {
		t.Error("agent-skills should be hidden from user-facing error")
	}
}

func TestParityRenderUnknownTargetErrorBracketNoise(t *testing.T) {
	msg := core.RenderUnknownTargetError("['copilot'", []string{"claude", "copilot"})
	if strings.Contains(msg, "['") {
		t.Errorf("bracket noise should be stripped, got: %s", msg)
	}
}

func TestParityRenderConflictingSchemaError(t *testing.T) {
	msg := core.RenderConflictingSchemaError()
	if !strings.Contains(msg, "[x] Cannot use both") {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestParityErrorConstructors(t *testing.T) {
	errs := []error{
		core.NewNoHarnessError(),
		core.NewAmbiguousHarnessError([]string{"a", "b"}),
		core.NewUnknownTargetError("x", []string{"claude"}),
		core.NewConflictingTargetsError(),
		core.NewEmptyTargetsListError(),
	}
	for _, e := range errs {
		if e.Error() == "" {
			t.Error("expected non-empty error message")
		}
	}
}

// ---------------------------------------------------------------------------
// Parity: scope
// ---------------------------------------------------------------------------

func TestParityScopeProject(t *testing.T) {
	s, ok := core.ParseScope("project")
	if !ok || s != core.ScopeProject {
		t.Error("project scope parse failed")
	}
	if s.String() != "project" {
		t.Error("scope.String() wrong")
	}
}

func TestParityScopeUser(t *testing.T) {
	s, ok := core.ParseScope("user")
	if !ok || s != core.ScopeUser {
		t.Error("user scope parse failed")
	}
	if s.String() != "user" {
		t.Error("scope.String() wrong")
	}
}

func TestParityScopeDefault(t *testing.T) {
	s, ok := core.ParseScope("")
	if !ok || s != core.ScopeProject {
		t.Error("empty string should map to project scope")
	}
}

func TestParityGetDeployRoot(t *testing.T) {
	cwd := "/tmp/proj"
	home := "/home/user"
	if core.GetDeployRoot(core.ScopeProject, cwd, home) != cwd {
		t.Error("project deploy root should be cwd")
	}
	if core.GetDeployRoot(core.ScopeUser, cwd, home) != home {
		t.Error("user deploy root should be home")
	}
}

func TestParityGetAPMDir(t *testing.T) {
	cwd := "/tmp/proj"
	home := "/home/user"
	if core.GetAPMDir(core.ScopeProject, cwd, home) != cwd {
		t.Error("project apm dir should be cwd")
	}
	expected := filepath.Join(home, ".apm")
	if core.GetAPMDir(core.ScopeUser, cwd, home) != expected {
		t.Errorf("user apm dir wrong: got %s want %s", core.GetAPMDir(core.ScopeUser, cwd, home), expected)
	}
}

// ---------------------------------------------------------------------------
// Parity: target_detection
// ---------------------------------------------------------------------------

func TestParityDetectTargetExplicit(t *testing.T) {
	cases := []struct{ input, want string }{
		{"copilot", "vscode"},
		{"vscode", "vscode"},
		{"agents", "vscode"},
		{"claude", "claude"},
		{"cursor", "cursor"},
		{"opencode", "opencode"},
		{"codex", "codex"},
		{"gemini", "gemini"},
		{"windsurf", "windsurf"},
		{"all", "all"},
	}
	for _, c := range cases {
		got, reason := core.DetectTarget("/tmp/empty", c.input, "")
		if got != c.want {
			t.Errorf("DetectTarget explicit %q: want %q got %q", c.input, c.want, got)
		}
		if reason != "explicit --target flag" {
			t.Errorf("expected reason 'explicit --target flag', got %q", reason)
		}
	}
}

func TestParityDetectTargetConfig(t *testing.T) {
	got, reason := core.DetectTarget("/tmp/empty", "", "claude")
	if got != "claude" || reason != "apm.yml target" {
		t.Errorf("config target: got %q/%q", got, reason)
	}
}

func TestParityDetectTargetAutoGitHub(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".github"), 0755); err != nil {
		t.Fatal(err)
	}
	got, reason := core.DetectTarget(dir, "", "")
	if got != "vscode" {
		t.Errorf("expected vscode, got %q", got)
	}
	if !strings.Contains(reason, ".github/") {
		t.Errorf("unexpected reason: %q", reason)
	}
}

func TestParityDetectTargetAutoMultiple(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, ".github"), 0755)
	os.Mkdir(filepath.Join(dir, ".claude"), 0755)
	got, _ := core.DetectTarget(dir, "", "")
	if got != "all" {
		t.Errorf("expected all, got %q", got)
	}
}

func TestParityDetectTargetNoFolder(t *testing.T) {
	dir := t.TempDir()
	got, reason := core.DetectTarget(dir, "", "")
	if got != "minimal" {
		t.Errorf("expected minimal, got %q", got)
	}
	if reason != core.ReasonNoTargetFolder {
		t.Errorf("unexpected reason: %q", reason)
	}
}

func TestParityShouldCompile(t *testing.T) {
	agentsTargets := []string{"vscode", "opencode", "codex", "gemini", "windsurf", "all", "minimal"}
	for _, t2 := range agentsTargets {
		if !core.ShouldCompileAgentsMD(t2) {
			t.Errorf("ShouldCompileAgentsMD(%q) should be true", t2)
		}
	}
	if core.ShouldCompileAgentsMD("claude") {
		t.Error("ShouldCompileAgentsMD(claude) should be false")
	}
	if !core.ShouldCompileClaudeMD("claude") || !core.ShouldCompileClaudeMD("all") {
		t.Error("ShouldCompileClaudeMD wrong")
	}
	if core.ShouldCompileClaudeMD("vscode") {
		t.Error("ShouldCompileClaudeMD(vscode) should be false")
	}
	if !core.ShouldCompileGeminiMD("gemini") || !core.ShouldCompileGeminiMD("all") {
		t.Error("ShouldCompileGeminiMD wrong")
	}
	if !core.ShouldCompileCopilotInstructionsMD("vscode") || !core.ShouldCompileCopilotInstructionsMD("all") {
		t.Error("ShouldCompileCopilotInstructionsMD wrong")
	}
	if core.ShouldCompileCopilotInstructionsMD("claude") {
		t.Error("ShouldCompileCopilotInstructionsMD(claude) should be false")
	}
}

func TestParityGetTargetDescription(t *testing.T) {
	if !strings.Contains(core.GetTargetDescription("vscode"), "AGENTS.md") {
		t.Error("vscode description should mention AGENTS.md")
	}
	if !strings.Contains(core.GetTargetDescription("copilot"), "AGENTS.md") {
		t.Error("copilot alias should resolve to vscode description")
	}
	if !strings.Contains(core.GetTargetDescription("claude"), "CLAUDE.md") {
		t.Error("claude description should mention CLAUDE.md")
	}
}

func TestParityNormalizeTargetList(t *testing.T) {
	if core.NormalizeTargetList(nil) != nil {
		t.Error("nil input should return nil")
	}
	got := core.NormalizeTargetList([]string{"copilot"})
	if len(got) != 1 || got[0] != "vscode" {
		t.Errorf("alias resolution failed: %v", got)
	}
	got = core.NormalizeTargetList([]string{"claude", "copilot", "claude"})
	if len(got) != 2 {
		t.Errorf("dedup failed: %v", got)
	}
	got = core.NormalizeTargetList([]string{"all"})
	if len(got) == 0 {
		t.Error("all should expand to all canonical targets")
	}
}

// ---------------------------------------------------------------------------
// Parity: apm_yml
// ---------------------------------------------------------------------------

func TestParityParseTargetsFieldPlural(t *testing.T) {
	data := map[string]interface{}{
		"targets": []interface{}{"claude", "copilot"},
	}
	got, err := core.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "claude" || got[1] != "copilot" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestParityParseTargetsFieldSingular(t *testing.T) {
	data := map[string]interface{}{
		"target": "claude",
	}
	got, err := core.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "claude" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestParityParseTargetsFieldCSV(t *testing.T) {
	data := map[string]interface{}{
		"target": "claude,copilot",
	}
	got, err := core.ParseTargetsField(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("CSV parse failed: %v", got)
	}
}

func TestParityParseTargetsFieldEmpty(t *testing.T) {
	got, err := core.ParseTargetsField(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestParityParseTargetsFieldConflict(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{"claude"}, "target": "copilot"}
	_, err := core.ParseTargetsField(data)
	if err == nil {
		t.Error("expected conflict error")
	}
}

func TestParityParseTargetsFieldEmptyList(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{}}
	_, err := core.ParseTargetsField(data)
	if err == nil {
		t.Error("expected empty list error")
	}
}

func TestParityParseTargetsFieldUnknownTarget(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{"claude", "unknown-target"}}
	_, err := core.ParseTargetsField(data)
	if err == nil {
		t.Error("expected unknown target error")
	}
}
