package finalize

import (
	"strings"
	"testing"
)

func TestUnpinnedWarning_Zero(t *testing.T) {
	msg := UnpinnedWarning(0, nil)
	if !strings.Contains(msg, "0 depend") {
		t.Errorf("expected count in message, got %q", msg)
	}
}

func TestUnpinnedWarning_SingleName(t *testing.T) {
	msg := UnpinnedWarning(1, []string{"my-pkg"})
	if !strings.Contains(msg, "my-pkg") {
		t.Errorf("expected package name in message: %q", msg)
	}
	if !strings.Contains(msg, "1 dependency") {
		t.Errorf("expected singular 'dependency': %q", msg)
	}
}

func TestUnpinnedWarning_FourNames(t *testing.T) {
	msg := UnpinnedWarning(4, []string{"a", "b", "c", "d"})
	if strings.Contains(msg, "more") {
		t.Errorf("should not show 'more' for 4 names: %q", msg)
	}
	for _, n := range []string{"a", "b", "c", "d"} {
		if !strings.Contains(msg, n) {
			t.Errorf("expected name %q in message: %q", n, msg)
		}
	}
}

func TestUnpinnedWarning_SixNames(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f"}
	msg := UnpinnedWarning(6, names)
	if !strings.Contains(msg, "and 1 more") {
		t.Errorf("expected 'and 1 more': %q", msg)
	}
}

func TestUnpinnedWarning_ContainsDriftHint(t *testing.T) {
	msg := UnpinnedWarning(2, nil)
	if !strings.Contains(msg, "drift") {
		t.Errorf("expected drift hint in message: %q", msg)
	}
}

func TestVerboseStatLines_Prompts(t *testing.T) {
	// TotalPromptsIntegrated is not surfaced by VerboseStatLines currently;
	// zero stats should yield no lines.
	lines := VerboseStatLines(InstallStats{TotalPromptsIntegrated: 3})
	if len(lines) != 0 {
		t.Errorf("TotalPromptsIntegrated is not tracked by VerboseStatLines; expected 0 lines, got %v", lines)
	}
}

func TestVerboseStatLines_AllFields(t *testing.T) {
	stats := InstallStats{
		LinksResolved:          2,
		CommandsIntegrated:     3,
		HooksIntegrated:        1,
		InstructionsIntegrated: 5,
	}
	lines := VerboseStatLines(stats)
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %v", len(lines), lines)
	}
	joined := strings.Join(lines, "\n")
	for _, want := range []string{"2", "3", "1", "5"} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected count %q in output: %s", want, joined)
		}
	}
}

func TestVerboseStatLines_SingleField_Links(t *testing.T) {
	lines := VerboseStatLines(InstallStats{LinksResolved: 1})
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "1") {
		t.Errorf("unexpected content: %q", lines[0])
	}
}

func TestInstallStats_ZeroValue(t *testing.T) {
	var s InstallStats
	if s.LinksResolved != 0 || s.CommandsIntegrated != 0 {
		t.Error("zero value should have all zero fields")
	}
}

func TestInstallResult_ZeroValue(t *testing.T) {
	var r InstallResult
	if r.InstalledCount != 0 || r.PackageTypes != nil {
		t.Error("zero InstallResult should have zero count and nil map")
	}
}
