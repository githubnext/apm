package contentscanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/security/contentscanner"
)

func TestScanText_NoFindings(t *testing.T) {
	findings := contentscanner.ScanText("test.md", "Hello, world!\nNo hidden chars here.")
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d: %v", len(findings), findings)
	}
}

func TestScanText_ZeroWidthSpace(t *testing.T) {
	// U+200B zero-width space
	findings := contentscanner.ScanText("test.md", "Hello\u200Bworld")
	if len(findings) == 0 {
		t.Fatal("expected finding for zero-width space")
	}
	f := findings[0]
	if f.Category != "zero-width" {
		t.Errorf("expected category 'zero-width', got %q", f.Category)
	}
	if f.Severity != "warning" {
		t.Errorf("expected severity 'warning', got %q", f.Severity)
	}
}

func TestScanText_BidiOverride(t *testing.T) {
	// U+202E right-to-left override
	findings := contentscanner.ScanText("file.md", "text\u202Emore")
	if len(findings) == 0 {
		t.Fatal("expected finding for bidi override")
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected critical severity for bidi override, got %q", findings[0].Severity)
	}
}

func TestScanText_LineNumberTracking(t *testing.T) {
	findings := contentscanner.ScanText("f.md", "line1\nline2\u200Bsuffix\nline3")
	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}
	if findings[0].Line != 2 {
		t.Errorf("expected line 2, got %d", findings[0].Line)
	}
}

func TestScanText_FilePathPopulated(t *testing.T) {
	findings := contentscanner.ScanText("myfile.md", "x\u200By")
	if len(findings) == 0 {
		t.Fatal("expected finding")
	}
	if findings[0].File != "myfile.md" {
		t.Errorf("expected File='myfile.md', got %q", findings[0].File)
	}
}

func TestScanFile_ReturnsErrorForMissing(t *testing.T) {
	_, err := contentscanner.ScanFile("/nonexistent/path/file.md")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestScanFile_ScansRealFile(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "test.md")
	if err := os.WriteFile(fp, []byte("clean content"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings, err := contentscanner.ScanFile(fp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for clean file, got %d", len(findings))
	}
}

func TestNewDefaultScanner_Extensions(t *testing.T) {
	s := contentscanner.NewDefaultScanner()
	if len(s.Extensions) == 0 {
		t.Fatal("expected default extensions")
	}
}

func TestContentScanner_ScanFiles_SkipsUnknownExtension(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "file.go")
	// embed a zero-width space in a .go file
	if err := os.WriteFile(fp, []byte("x\u200By"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := contentscanner.NewDefaultScanner()
	results := s.ScanFiles([]string{fp})
	if len(results) != 0 {
		t.Fatalf("expected .go file to be skipped, got %v", results)
	}
}
