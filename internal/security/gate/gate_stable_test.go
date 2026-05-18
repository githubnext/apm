package gate

import (
"os"
"path/filepath"
"testing"
)

func TestScanPolicy_OnCritical_field(t *testing.T) {
p := ScanPolicy{OnCritical: OnCriticalWarn}
if p.OnCritical != OnCriticalWarn {
t.Errorf("expected OnCriticalWarn, got %s", p.OnCritical)
}
}

func TestScanPolicy_ForceOverrides_true(t *testing.T) {
p := ScanPolicy{OnCritical: OnCriticalBlock, ForceOverrides: true}
if !p.ForceOverrides {
t.Error("expected ForceOverrides=true")
}
}

func TestOnCriticalBlock_constant(t *testing.T) {
if OnCriticalBlock != "block" {
t.Errorf("unexpected value: %s", OnCriticalBlock)
}
}

func TestOnCriticalWarn_constant(t *testing.T) {
if OnCriticalWarn != "warn" {
t.Errorf("unexpected value: %s", OnCriticalWarn)
}
}

func TestOnCriticalIgnore_constant(t *testing.T) {
if OnCriticalIgnore != "ignore" {
t.Errorf("unexpected value: %s", OnCriticalIgnore)
}
}

func TestBlockPolicy_OnCritical(t *testing.T) {
if BlockPolicy.OnCritical != OnCriticalBlock {
t.Errorf("expected OnCriticalBlock, got %s", BlockPolicy.OnCritical)
}
}

func TestBlockPolicy_ForceOverrides(t *testing.T) {
if !BlockPolicy.ForceOverrides {
t.Error("BlockPolicy.ForceOverrides should be true")
}
}

func TestWarnPolicy_OnCritical(t *testing.T) {
if WarnPolicy.OnCritical != OnCriticalWarn {
t.Errorf("expected OnCriticalWarn, got %s", WarnPolicy.OnCritical)
}
}

func TestWarnPolicy_ForceOverrides(t *testing.T) {
if WarnPolicy.ForceOverrides {
t.Error("WarnPolicy.ForceOverrides should be false")
}
}

func TestReportPolicy_OnCritical(t *testing.T) {
if ReportPolicy.OnCritical != OnCriticalIgnore {
t.Errorf("expected OnCriticalIgnore, got %s", ReportPolicy.OnCritical)
}
}

func TestEffectiveBlock_IgnorePolicy_neverBlocks(t *testing.T) {
p := ScanPolicy{OnCritical: OnCriticalIgnore, ForceOverrides: false}
if p.EffectiveBlock(false) {
t.Error("OnCriticalIgnore should never block")
}
if p.EffectiveBlock(true) {
t.Error("OnCriticalIgnore should never block even with force")
}
}

func TestEffectiveBlock_WarnPolicy_neverBlocks(t *testing.T) {
p := ScanPolicy{OnCritical: OnCriticalWarn, ForceOverrides: false}
if p.EffectiveBlock(false) || p.EffectiveBlock(true) {
t.Error("OnCriticalWarn should never block")
}
}

func TestEffectiveBlock_BlockPolicy_noForce(t *testing.T) {
p := ScanPolicy{OnCritical: OnCriticalBlock, ForceOverrides: true}
if !p.EffectiveBlock(false) {
t.Error("should block without force")
}
}

func TestEffectiveBlock_BlockPolicy_withForce(t *testing.T) {
p := ScanPolicy{OnCritical: OnCriticalBlock, ForceOverrides: true}
if p.EffectiveBlock(true) {
t.Error("should not block when ForceOverrides=true and force=true")
}
}

func TestEffectiveBlock_BlockNoOverride_force(t *testing.T) {
p := ScanPolicy{OnCritical: OnCriticalBlock, ForceOverrides: false}
if !p.EffectiveBlock(true) {
t.Error("should still block when ForceOverrides=false even with force=true")
}
}

func TestScanVerdict_HasFindings_empty(t *testing.T) {
v := ScanVerdict{}
if v.HasFindings() {
t.Error("empty verdict should not have findings")
}
}

func TestScanVerdict_Fields(t *testing.T) {
v := ScanVerdict{
HasCritical:   true,
ShouldBlock:   true,
CriticalCount: 2,
WarningCount:  3,
FilesScanned:  5,
}
if v.CriticalCount != 2 {
t.Errorf("expected CriticalCount=2, got %d", v.CriticalCount)
}
if v.WarningCount != 3 {
t.Errorf("expected WarningCount=3, got %d", v.WarningCount)
}
if v.FilesScanned != 5 {
t.Errorf("expected FilesScanned=5, got %d", v.FilesScanned)
}
}

func TestGate_New_ReportPolicy(t *testing.T) {
g := New(ReportPolicy, false)
if g == nil {
t.Fatal("expected non-nil gate")
}
}

func TestGate_Check_singleCleanFile(t *testing.T) {
dir := t.TempDir()
f := filepath.Join(dir, "clean.md")
_ = os.WriteFile(f, []byte("# Title\nJust some plain text."), 0o644)
g := New(ReportPolicy, false)
v := g.Check([]string{f})
if v.FilesScanned != 1 {
t.Errorf("expected 1 file scanned, got %d", v.FilesScanned)
}
}

func TestGate_Check_twoCleanFiles(t *testing.T) {
dir := t.TempDir()
f1 := filepath.Join(dir, "a.md")
f2 := filepath.Join(dir, "b.md")
_ = os.WriteFile(f1, []byte("Safe content"), 0o644)
_ = os.WriteFile(f2, []byte("Also safe"), 0o644)
g := New(WarnPolicy, false)
v := g.Check([]string{f1, f2})
if v.FilesScanned != 2 {
t.Errorf("expected 2 files scanned, got %d", v.FilesScanned)
}
if v.ShouldBlock {
t.Error("WarnPolicy should not block clean files")
}
}
