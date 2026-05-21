package gate

import (
	"testing"
)

func TestOnCritical_Constants(t *testing.T) {
	if OnCriticalBlock != "block" {
		t.Errorf("block constant wrong: %s", OnCriticalBlock)
	}
	if OnCriticalWarn != "warn" {
		t.Errorf("warn constant wrong: %s", OnCriticalWarn)
	}
	if OnCriticalIgnore != "ignore" {
		t.Errorf("ignore constant wrong: %s", OnCriticalIgnore)
	}
}

func TestScanPolicy_ZeroValue(t *testing.T) {
	var p ScanPolicy
	if p.ForceOverrides {
		t.Error("ForceOverrides should default to false")
	}
}

func TestWarnPolicy_Fields(t *testing.T) {
	if WarnPolicy.OnCritical != OnCriticalWarn {
		t.Errorf("WarnPolicy.OnCritical = %s, want warn", WarnPolicy.OnCritical)
	}
	if WarnPolicy.ForceOverrides {
		t.Error("WarnPolicy should not override force")
	}
}

func TestReportPolicy_Fields(t *testing.T) {
	if ReportPolicy.OnCritical != OnCriticalIgnore {
		t.Errorf("ReportPolicy.OnCritical = %s, want ignore", ReportPolicy.OnCritical)
	}
}

func TestBlockPolicy_NeverBlocksWhenForceAndOverridable(t *testing.T) {
	if BlockPolicy.EffectiveBlock(true) {
		t.Error("block policy with ForceOverrides=true should not block when force=true")
	}
}

func TestBlockPolicy_BlocksWhenNoForce(t *testing.T) {
	if !BlockPolicy.EffectiveBlock(false) {
		t.Error("block policy should block when force=false")
	}
}

func TestWarnPolicy_DoesNotBlock(t *testing.T) {
	if WarnPolicy.EffectiveBlock(false) {
		t.Error("warn policy should never block")
	}
	if WarnPolicy.EffectiveBlock(true) {
		t.Error("warn policy should never block")
	}
}

func TestReportPolicy_DoesNotBlock(t *testing.T) {
	if ReportPolicy.EffectiveBlock(false) {
		t.Error("report policy should never block")
	}
}

func TestScanVerdict_ZeroValue(t *testing.T) {
	var v ScanVerdict
	if v.HasCritical {
		t.Error("default HasCritical should be false")
	}
	if v.HasFindings() {
		t.Error("empty verdict should have no findings")
	}
}

func TestGate_NewReturnsNonNil(t *testing.T) {
	g := New(BlockPolicy, false)
	if g == nil {
		t.Fatal("New should return non-nil Gate")
	}
}

func TestGate_CheckNilSlice(t *testing.T) {
	g := New(BlockPolicy, false)
	v := g.Check(nil)
	if v.FilesScanned != 0 {
		t.Errorf("expected 0 files scanned, got %d", v.FilesScanned)
	}
	if v.HasFindings() {
		t.Error("no files should produce no findings")
	}
}

func TestScanVerdict_FilesScanned(t *testing.T) {
	v := ScanVerdict{FilesScanned: 5, CriticalCount: 0}
	if v.FilesScanned != 5 {
		t.Errorf("FilesScanned = %d, want 5", v.FilesScanned)
	}
}
