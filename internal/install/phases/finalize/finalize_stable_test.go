package finalize

import (
"strings"
"testing"
)

func TestUnpinnedWarning_ZeroCount(t *testing.T) {
msg := UnpinnedWarning(0, nil)
// zero count should still return a string (even if trivial)
_ = msg
}

func TestUnpinnedWarning_OneWithName(t *testing.T) {
msg := UnpinnedWarning(1, []string{"my-dep"})
if !strings.Contains(msg, "my-dep") {
t.Errorf("expected 'my-dep' in message: %q", msg)
}
}

func TestUnpinnedWarning_TwoWithNames(t *testing.T) {
msg := UnpinnedWarning(2, []string{"dep-a", "dep-b"})
if !strings.Contains(msg, "dep-a") {
t.Errorf("expected 'dep-a' in message: %q", msg)
}
if !strings.Contains(msg, "dep-b") {
t.Errorf("expected 'dep-b' in message: %q", msg)
}
}

func TestUnpinnedWarning_SixNames_Truncates(t *testing.T) {
names := []string{"a", "b", "c", "d", "e", "f"}
msg := UnpinnedWarning(6, names)
if !strings.Contains(msg, "and 1 more") {
t.Errorf("expected 'and 1 more', got: %q", msg)
}
}

func TestUnpinnedWarning_TenNames_Truncates(t *testing.T) {
names := make([]string, 10)
for i := range names {
names[i] = "pkg"
}
msg := UnpinnedWarning(10, names)
if !strings.Contains(msg, "more") {
t.Errorf("expected truncation in message: %q", msg)
}
}

func TestInstallStats_zero(t *testing.T) {
s := InstallStats{}
if s.LinksResolved != 0 {
t.Errorf("LinksResolved should be 0, got %d", s.LinksResolved)
}
}

func TestInstallStats_fields(t *testing.T) {
s := InstallStats{
LinksResolved:          1,
CommandsIntegrated:     2,
HooksIntegrated:        3,
InstructionsIntegrated: 4,
}
if s.LinksResolved != 1 {
t.Errorf("LinksResolved = %d, want 1", s.LinksResolved)
}
if s.CommandsIntegrated != 2 {
t.Errorf("CommandsIntegrated = %d, want 2", s.CommandsIntegrated)
}
if s.HooksIntegrated != 3 {
t.Errorf("HooksIntegrated = %d, want 3", s.HooksIntegrated)
}
if s.InstructionsIntegrated != 4 {
t.Errorf("InstructionsIntegrated = %d, want 4", s.InstructionsIntegrated)
}
}

func TestVerboseStatLines_LinksAndHooks(t *testing.T) {
lines := VerboseStatLines(InstallStats{
LinksResolved: 5,
HooksIntegrated: 2,
})
if len(lines) != 2 {
t.Errorf("expected 2 lines, got %d: %v", len(lines), lines)
}
}

func TestVerboseStatLines_OnlyCommands(t *testing.T) {
lines := VerboseStatLines(InstallStats{CommandsIntegrated: 10})
if len(lines) != 1 {
t.Fatalf("expected 1 line, got %d", len(lines))
}
if !strings.Contains(lines[0], "10") {
t.Errorf("expected count 10 in line: %q", lines[0])
}
}

func TestVerboseStatLines_OnlyInstructions(t *testing.T) {
lines := VerboseStatLines(InstallStats{InstructionsIntegrated: 3})
if len(lines) != 1 {
t.Fatalf("expected 1 line, got %d", len(lines))
}
if !strings.Contains(lines[0], "3") {
t.Errorf("expected count 3 in line: %q", lines[0])
}
}

func TestVerboseStatLines_AllNonzero(t *testing.T) {
lines := VerboseStatLines(InstallStats{
LinksResolved:          1,
CommandsIntegrated:     1,
HooksIntegrated:        1,
InstructionsIntegrated: 1,
})
if len(lines) != 4 {
t.Errorf("expected 4 lines, got %d: %v", len(lines), lines)
}
}

func TestInstallResult_zero(t *testing.T) {
r := InstallResult{}
_ = r
}
