package finalize

import (
	"strings"
	"testing"
)

func TestUnpinnedWarning_SingleNoNames(t *testing.T) {
	msg := UnpinnedWarning(1, nil)
	if !strings.Contains(msg, "1 dependency unpinned") {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestUnpinnedWarning_PluralNoNames(t *testing.T) {
	msg := UnpinnedWarning(3, nil)
	if !strings.Contains(msg, "3 dependencies unpinned") {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestUnpinnedWarning_WithNames(t *testing.T) {
	msg := UnpinnedWarning(2, []string{"foo", "bar"})
	if !strings.Contains(msg, "foo") || !strings.Contains(msg, "bar") {
		t.Errorf("expected names in message: %q", msg)
	}
	if !strings.Contains(msg, "2 dependencies") {
		t.Errorf("expected count in message: %q", msg)
	}
}

func TestUnpinnedWarning_TruncatesAt5(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	msg := UnpinnedWarning(7, names)
	if !strings.Contains(msg, "and 2 more") {
		t.Errorf("expected 'and 2 more' in message: %q", msg)
	}
}

func TestUnpinnedWarning_ExactlyFive(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e"}
	msg := UnpinnedWarning(5, names)
	if strings.Contains(msg, "more") {
		t.Errorf("unexpected 'more' with exactly 5 names: %q", msg)
	}
}

func TestVerboseStatLines_AllZero(t *testing.T) {
	lines := VerboseStatLines(InstallStats{})
	if len(lines) != 0 {
		t.Errorf("expected no lines for zero stats, got %v", lines)
	}
}

func TestVerboseStatLines_LinksResolved(t *testing.T) {
	lines := VerboseStatLines(InstallStats{LinksResolved: 3})
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "3") {
		t.Errorf("expected count in line: %q", lines[0])
	}
}

func TestVerboseStatLines_Commands(t *testing.T) {
	lines := VerboseStatLines(InstallStats{CommandsIntegrated: 5})
	if len(lines) != 1 || !strings.Contains(lines[0], "5 command") {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestVerboseStatLines_Hooks(t *testing.T) {
	lines := VerboseStatLines(InstallStats{HooksIntegrated: 2})
	if len(lines) != 1 || !strings.Contains(lines[0], "2 hook") {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestVerboseStatLines_Instructions(t *testing.T) {
	lines := VerboseStatLines(InstallStats{InstructionsIntegrated: 7})
	if len(lines) != 1 || !strings.Contains(lines[0], "7 instruction") {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestVerboseStatLines_Multiple(t *testing.T) {
	lines := VerboseStatLines(InstallStats{
		LinksResolved:          1,
		CommandsIntegrated:     2,
		HooksIntegrated:        3,
		InstructionsIntegrated: 4,
	})
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d: %v", len(lines), lines)
	}
}
