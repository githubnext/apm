package contentscanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFinding_ZeroValue(t *testing.T) {
	var f ScanFinding
	if f.File != "" || f.Severity != "" || f.Category != "" {
		t.Error("zero value ScanFinding should have empty string fields")
	}
	if f.Line != 0 || f.Column != 0 {
		t.Error("zero value ScanFinding should have zero int fields")
	}
}

func TestScanFinding_AllFields(t *testing.T) {
	f := ScanFinding{
		File:        "test.py",
		Line:        5,
		Column:      12,
		Char:        0x200B,
		Codepoint:   "U+200B",
		Severity:    "warning",
		Category:    "zero-width",
		Description: "Zero-width space",
	}
	if f.File != "test.py" {
		t.Errorf("File: %q", f.File)
	}
	if f.Line != 5 {
		t.Errorf("Line: %d", f.Line)
	}
	if f.Codepoint != "U+200B" {
		t.Errorf("Codepoint: %q", f.Codepoint)
	}
	if f.Severity != "warning" {
		t.Errorf("Severity: %q", f.Severity)
	}
}

func TestScanText_BidiOverrideCritical(t *testing.T) {
	// U+202E is "right-to-left override"
	content := "before\u202Eafter"
	findings := ScanText("test.txt", content)
	if len(findings) == 0 {
		t.Fatal("expected at least one finding for bidi override")
	}
	found := false
	for _, f := range findings {
		if f.Severity == "critical" && f.Category == "bidi-override" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected critical bidi-override finding, got: %v", findings)
	}
}

func TestScanText_NoSuspiciousChars(t *testing.T) {
	findings := ScanText("clean.txt", "Hello, world! This is plain ASCII text.")
	if len(findings) != 0 {
		t.Errorf("expected no findings for clean text, got %d", len(findings))
	}
}

func TestScanText_FilePath(t *testing.T) {
	content := "line1\u200Bline2"
	findings := ScanText("/path/to/file.txt", content)
	if len(findings) == 0 {
		t.Fatal("expected findings")
	}
	if findings[0].File != "/path/to/file.txt" {
		t.Errorf("File: %q", findings[0].File)
	}
}

func TestScanFile_ValidFile_Clean(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "clean.txt")
	if err := os.WriteFile(p, []byte("clean content here"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings, err := ScanFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings for clean file, got %d", len(findings))
	}
}

func TestContentScanner_ScanFiles_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	cs := NewDefaultScanner()
	files := []string{}
	for _, name := range []string{"a.md", "b.md", "c.md"} {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte("clean"), 0o644); err != nil {
			t.Fatal(err)
		}
		files = append(files, p)
	}
	// ScanFiles only returns entries for files with findings; clean files are absent
	results := cs.ScanFiles(files)
	// Verify no false-positive findings for clean files
	for _, p := range files {
		if findings, ok := results[p]; ok && len(findings) > 0 {
			t.Errorf("unexpected findings in clean file %q: %v", p, findings)
		}
	}
}

func TestNewDefaultScanner_NotNil(t *testing.T) {
	cs := NewDefaultScanner()
	if cs == nil {
		t.Error("NewDefaultScanner should return non-nil")
	}
}

func TestNewDefaultScanner_HasExtensions(t *testing.T) {
	cs := NewDefaultScanner()
	if len(cs.Extensions) == 0 {
		t.Error("default scanner should have extensions")
	}
}

func TestScanText_ZeroWidthWarning(t *testing.T) {
	content := "hello\u200Bworld"
	findings := ScanText("f.txt", content)
	if len(findings) == 0 {
		t.Fatal("expected findings for zero-width space")
	}
	if findings[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %q", findings[0].Severity)
	}
}

func TestScanText_ColumnStartsAtOne(t *testing.T) {
	content := "\u200Bstart"
	findings := ScanText("f.txt", content)
	if len(findings) == 0 {
		t.Fatal("expected findings")
	}
	if findings[0].Column < 1 {
		t.Errorf("column should be >= 1, got %d", findings[0].Column)
	}
}

func TestContentScanner_ScanFiles_NonexistentFile(t *testing.T) {
	cs := NewDefaultScanner()
	results := cs.ScanFiles([]string{"/nonexistent/path/file.md"})
	_ = results
}
