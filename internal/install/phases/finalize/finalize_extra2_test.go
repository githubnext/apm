package finalize

import (
	"strings"
	"testing"
)

func TestUnpinnedWarning_TwoNames(t *testing.T) {
	msg := UnpinnedWarning(2, []string{"a/b", "c/d"})
	if !strings.Contains(msg, "a/b") {
		t.Errorf("expected a/b in message, got %q", msg)
	}
	if !strings.Contains(msg, "dependencies") {
		t.Errorf("expected plural 'dependencies', got %q", msg)
	}
}

func TestUnpinnedWarning_OneName(t *testing.T) {
	msg := UnpinnedWarning(1, []string{"x/y"})
	if !strings.Contains(msg, "dependency") {
		t.Errorf("expected singular 'dependency', got %q", msg)
	}
	if !strings.Contains(msg, "x/y") {
		t.Errorf("expected x/y in message, got %q", msg)
	}
}

func TestUnpinnedWarning_ZeroCountNoNames(t *testing.T) {
	msg := UnpinnedWarning(0, nil)
	if msg == "" {
		t.Error("expected non-empty message even for 0")
	}
}

func TestInstallStats_Fields(t *testing.T) {
	s := InstallStats{
		LinksResolved:      3,
		CommandsIntegrated: 2,
		InstalledCount:     5,
	}
	if s.LinksResolved != 3 {
		t.Errorf("LinksResolved = %d, want 3", s.LinksResolved)
	}
	if s.InstalledCount != 5 {
		t.Errorf("InstalledCount = %d, want 5", s.InstalledCount)
	}
}

func TestInstallResult_Fields(t *testing.T) {
	r := InstallResult{
		InstalledCount:         4,
		TotalPromptsIntegrated: 2,
		TotalAgentsIntegrated:  1,
		Warnings:               []string{"w1"},
		Errors:                 []string{},
	}
	if r.InstalledCount != 4 {
		t.Errorf("InstalledCount = %d, want 4", r.InstalledCount)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestVerboseStatLines_HooksIntegrated(t *testing.T) {
	lines := VerboseStatLines(InstallStats{HooksIntegrated: 2})
	if len(lines) == 0 {
		t.Error("expected lines for HooksIntegrated")
	}
}

func TestVerboseStatLines_InstructionsNonZero(t *testing.T) {
	lines := VerboseStatLines(InstallStats{InstructionsIntegrated: 4})
	if len(lines) == 0 {
		t.Error("expected lines for InstructionsIntegrated")
	}
}

func TestUnpinnedWarning_DriftHintPresent(t *testing.T) {
	msg := UnpinnedWarning(3, []string{"a", "b", "c"})
	if !strings.Contains(msg, "drift") {
		t.Errorf("expected drift hint in message, got %q", msg)
	}
}
