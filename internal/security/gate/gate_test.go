package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/security/contentscanner"
)

func TestEffectiveBlock(t *testing.T) {
	tests := []struct {
		policy ScanPolicy
		force  bool
		want   bool
	}{
		{BlockPolicy, false, true},
		{BlockPolicy, true, false},
		{WarnPolicy, false, false},
		{WarnPolicy, true, false},
		{ReportPolicy, false, false},
		{ReportPolicy, true, false},
	}
	for _, tt := range tests {
		got := tt.policy.EffectiveBlock(tt.force)
		if got != tt.want {
			t.Errorf("EffectiveBlock(%v, force=%v) = %v, want %v", tt.policy, tt.force, got, tt.want)
		}
	}
}

func TestScanVerdict_HasFindings(t *testing.T) {
	empty := ScanVerdict{}
	if empty.HasFindings() {
		t.Error("expected no findings for empty verdict")
	}
	nonEmpty := ScanVerdict{
		FindingsByFile: make(map[string][]contentscanner.ScanFinding),
	}
	// Empty map means no findings.
	if nonEmpty.HasFindings() {
		t.Error("expected no findings for verdict with empty map")
	}
}

func TestGate_CheckCleanFile(t *testing.T) {
	dir := t.TempDir()
	clean := filepath.Join(dir, "clean.md")
	if err := os.WriteFile(clean, []byte("# Hello world\nThis is safe content.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := New(BlockPolicy, false)
	v := g.Check([]string{clean})
	if v.HasCritical {
		t.Error("expected no critical findings in clean file")
	}
	if v.ShouldBlock {
		t.Error("expected no block for clean file")
	}
	if v.FilesScanned != 1 {
		t.Errorf("expected 1 file scanned, got %d", v.FilesScanned)
	}
}

func TestGate_CheckFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "safe.txt")
	if err := os.WriteFile(f, []byte("just text"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := New(WarnPolicy, false)
	v := g.CheckFile(f)
	if v.FilesScanned != 1 {
		t.Errorf("expected 1 file scanned, got %d", v.FilesScanned)
	}
}

func TestGate_CheckMissingFile(t *testing.T) {
	g := New(BlockPolicy, false)
	v := g.Check([]string{"/nonexistent/path/file.md"})
	if v.FilesScanned != 1 {
		t.Errorf("expected 1 file (even missing), got %d", v.FilesScanned)
	}
}
