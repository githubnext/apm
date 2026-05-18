package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/security/contentscanner"
)

func TestReportPolicy_NeverBlocks(t *testing.T) {
	for _, force := range []bool{true, false} {
		if ReportPolicy.EffectiveBlock(force) {
			t.Errorf("ReportPolicy.EffectiveBlock(force=%v) should never block", force)
		}
	}
}

func TestWarnPolicy_NeverBlocks(t *testing.T) {
	for _, force := range []bool{true, false} {
		if WarnPolicy.EffectiveBlock(force) {
			t.Errorf("WarnPolicy.EffectiveBlock(force=%v) should never block", force)
		}
	}
}

func TestBlockPolicy_BlocksWithoutForce(t *testing.T) {
	if !BlockPolicy.EffectiveBlock(false) {
		t.Error("BlockPolicy should block when force=false")
	}
	if BlockPolicy.EffectiveBlock(true) {
		t.Error("BlockPolicy should not block when force=true (ForceOverrides=true)")
	}
}

func TestScanVerdict_HasFindings_WithEntries(t *testing.T) {
	v := ScanVerdict{}
	v.FindingsByFile = map[string][]contentscanner.ScanFinding{}
	v.FindingsByFile["file.md"] = nil
	if !v.HasFindings() {
		t.Error("expected HasFindings=true when map has an entry")
	}
}

func TestGate_CheckEmptyPaths(t *testing.T) {
	g := New(BlockPolicy, false)
	v := g.Check([]string{})
	if v.FilesScanned != 0 {
		t.Errorf("expected 0 files scanned for empty input, got %d", v.FilesScanned)
	}
}

func TestGate_CheckNilPaths(t *testing.T) {
	g := New(WarnPolicy, false)
	v := g.Check(nil)
	if v.FilesScanned != 0 {
		t.Errorf("expected 0 files scanned for nil input, got %d", v.FilesScanned)
	}
}

func TestGate_CheckFileClearsFindings(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "content.txt")
	if err := os.WriteFile(f, []byte("safe text no issues"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := New(BlockPolicy, false)
	v := g.CheckFile(f)
	if v.ShouldBlock {
		t.Error("clean file should not trigger block")
	}
}

func TestGate_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	_ = os.WriteFile(f1, []byte("content a"), 0o644)
	_ = os.WriteFile(f2, []byte("content b"), 0o644)
	g := New(ReportPolicy, false)
	v := g.Check([]string{f1, f2})
	if v.FilesScanned != 2 {
		t.Errorf("expected 2 files scanned, got %d", v.FilesScanned)
	}
}
