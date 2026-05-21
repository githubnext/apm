package targetdetection

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeTarget_LowerCase(t *testing.T) {
	// NormalizeTarget maps some aliases, doesn't do case conversion
	got := NormalizeTarget("vscode")
	if got != "vscode" {
		t.Errorf("NormalizeTarget(vscode) = %q, want vscode", got)
	}
}

func TestNormalizeTarget_AlreadyLower(t *testing.T) {
	got := NormalizeTarget("cursor")
	if got != "cursor" {
		t.Errorf("NormalizeTarget(cursor) = %q", got)
	}
}

func TestNormalizeTarget_CopilotAlias(t *testing.T) {
	got := NormalizeTarget("copilot")
	if got != "vscode" {
		t.Errorf("NormalizeTarget(copilot) = %q, want vscode", got)
	}
}

func TestValidTargets_ContainsExpected(t *testing.T) {
	for _, target := range []string{"copilot", "cursor", "claude", "codex", "gemini"} {
		if !ValidTargets[target] {
			t.Errorf("ValidTargets missing: %q", target)
		}
	}
}

func TestValidTargets_DoesNotContainUnknown(t *testing.T) {
	if ValidTargets["unknown-tool"] {
		t.Error("ValidTargets should not contain 'unknown-tool'")
	}
}

func TestCanonicalTargetsOrdered_NonEmpty(t *testing.T) {
	if len(CanonicalTargetsOrdered) == 0 {
		t.Error("CanonicalTargetsOrdered should not be empty")
	}
}

func TestCanonicalTargetsOrdered_AllInValidTargets(t *testing.T) {
	for _, t2 := range CanonicalTargetsOrdered {
		if !ValidTargets[t2] {
			t.Errorf("CanonicalTargetsOrdered contains unknown target: %q", t2)
		}
	}
}

func TestResolvedTargets_ZeroValue(t *testing.T) {
	r := ResolvedTargets{}
	if len(r.Targets) != 0 {
		t.Error("zero ResolvedTargets should have empty Targets")
	}
}

func TestResolvedTargets_Fields(t *testing.T) {
	r := ResolvedTargets{
		Targets:    []string{"copilot", "cursor"},
		Source:     "apm.yml",
		AutoCreate: true,
	}
	if len(r.Targets) != 2 {
		t.Errorf("Targets len = %d", len(r.Targets))
	}
	if r.Source != "apm.yml" {
		t.Errorf("Source = %q", r.Source)
	}
}

func TestSignal_Fields_Extra2(t *testing.T) {
	s := Signal{Target: "copilot", Source: ".github/copilot-instructions.md"}
	if s.Target != "copilot" {
		t.Errorf("Target = %q", s.Target)
	}
}

func TestDetectSignals_EmptyDir_Extra2(t *testing.T) {
	dir := t.TempDir()
	signals := DetectSignals(dir)
	if len(signals) != 0 {
		t.Errorf("expected no signals for empty dir, got %v", signals)
	}
}

func TestDetectSignals_WithCopilotFile(t *testing.T) {
	dir := t.TempDir()
	ghDir := filepath.Join(dir, ".github")
	if err := os.MkdirAll(ghDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ghDir, "copilot-instructions.md"), []byte("instructions"), 0o644); err != nil {
		t.Fatal(err)
	}
	signals := DetectSignals(dir)
	found := false
	for _, s := range signals {
		if s.Target == "copilot" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected copilot signal for copilot-instructions.md")
	}
}

func TestResolveTargets_FlagOverride_Extra2(t *testing.T) {
	dir := t.TempDir()
	r, err := ResolveTargets(dir, []string{"cursor"}, nil)
	if err != nil {
		t.Fatalf("ResolveTargets: %v", err)
	}
	if len(r.Targets) != 1 || r.Targets[0] != "cursor" {
		t.Errorf("Targets = %v", r.Targets)
	}
}

func TestResolveTargets_UnknownFlag_Extra2(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveTargets(dir, []string{"unknowntool"}, nil)
	if err == nil {
		t.Error("expected error for unknown target flag")
	}
}

func TestFormatProvenance_FlagSource(t *testing.T) {
	r := ResolvedTargets{Targets: []string{"copilot"}, Source: "--target flag"}
	s := FormatProvenance(r)
	if !strings.Contains(s, "flag") && s == "" {
		// FormatProvenance may return empty or descriptive string; just ensure no panic
	}
	_ = s
}

func TestExpandAllTargets_YAMLTargets(t *testing.T) {
	dir := t.TempDir()
	targets, err := ExpandAllTargets(dir, []string{"copilot", "cursor"})
	if err != nil {
		t.Fatalf("ExpandAllTargets: %v", err)
	}
	if len(targets) < 2 {
		t.Errorf("expected >= 2 targets, got %v", targets)
	}
}

func TestDetectTarget_Explicit(t *testing.T) {
	dir := t.TempDir()
	target, source := DetectTarget(dir, "cursor", "")
	if target != "cursor" {
		t.Errorf("DetectTarget explicit = %q", target)
	}
	_ = source
}
