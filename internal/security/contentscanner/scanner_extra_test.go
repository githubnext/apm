package contentscanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/githubnext/apm/internal/security/contentscanner"
)

func TestScanText_TagCharacter(t *testing.T) {
	// U+E0001 is a unicode tag character (critical)
	findings := contentscanner.ScanText("f.md", "hello\U000E0001world")
	if len(findings) == 0 {
		t.Fatal("expected finding for tag character")
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected critical, got %q", findings[0].Severity)
	}
	if findings[0].Category != "tag-character" {
		t.Errorf("expected tag-character, got %q", findings[0].Category)
	}
}

func TestScanText_VariationSelectorWarning(t *testing.T) {
	// U+FE00 variation selector (warning)
	findings := contentscanner.ScanText("f.md", "x\uFE00y")
	if len(findings) == 0 {
		t.Fatal("expected finding for variation selector")
	}
	if findings[0].Severity != "warning" {
		t.Errorf("expected warning, got %q", findings[0].Severity)
	}
}

func TestScanText_InvisibleSeparatorInfo(t *testing.T) {
	// U+2028 line separator (info)
	findings := contentscanner.ScanText("f.md", "a\u2028b")
	if len(findings) == 0 {
		t.Fatal("expected finding for line separator")
	}
	if findings[0].Severity != "info" {
		t.Errorf("expected info, got %q", findings[0].Severity)
	}
}

func TestScanText_CodepointFormat(t *testing.T) {
	findings := contentscanner.ScanText("f.md", "\u200B")
	if len(findings) == 0 {
		t.Fatal("expected finding")
	}
	if findings[0].Codepoint != "U+200B" {
		t.Errorf("expected U+200B, got %q", findings[0].Codepoint)
	}
}

func TestScanText_MultipleFindings(t *testing.T) {
	// two zero-width spaces on same line
	findings := contentscanner.ScanText("f.md", "a\u200Bb\u200Cc")
	if len(findings) != 2 {
		t.Errorf("expected 2 findings, got %d", len(findings))
	}
}

func TestScanText_ColumnTracking(t *testing.T) {
	// U+200B at column 4 (0-indexed pos 3)
	findings := contentscanner.ScanText("f.md", "abc\u200Bdef")
	if len(findings) == 0 {
		t.Fatal("expected finding")
	}
	if findings[0].Column != 4 {
		t.Errorf("expected column 4, got %d", findings[0].Column)
	}
}

func TestScanText_EmptyInput(t *testing.T) {
	findings := contentscanner.ScanText("f.md", "")
	if len(findings) != 0 {
		t.Errorf("expected no findings for empty input, got %d", len(findings))
	}
}

func TestScanText_MultilineTracking(t *testing.T) {
	content := "line1\nline2\nline3\u202Eend"
	findings := contentscanner.ScanText("f.md", content)
	if len(findings) == 0 {
		t.Fatal("expected finding")
	}
	if findings[0].Line != 3 {
		t.Errorf("expected line 3, got %d", findings[0].Line)
	}
}

func TestScanFile_WithHiddenChar(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "test.md")
	if err := os.WriteFile(fp, []byte("hello\u200Bworld"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings, err := contentscanner.ScanFile(fp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) == 0 {
		t.Fatal("expected finding in file")
	}
}

func TestContentScanner_ScanFiles_MatchesExtension(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "test.md")
	if err := os.WriteFile(fp, []byte("x\u200By"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := contentscanner.NewDefaultScanner()
	results := s.ScanFiles([]string{fp})
	if len(results) == 0 {
		t.Fatal("expected .md file to be scanned")
	}
}

func TestContentScanner_EmptyExtensions_ScansAll(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "test.go")
	if err := os.WriteFile(fp, []byte("x\u200By"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := &contentscanner.ContentScanner{Extensions: nil}
	results := s.ScanFiles([]string{fp})
	if len(results) == 0 {
		t.Fatal("empty extensions should scan all files")
	}
}

func TestContentScanner_ScanFiles_EmptyList(t *testing.T) {
	s := contentscanner.NewDefaultScanner()
	results := s.ScanFiles(nil)
	if len(results) != 0 {
		t.Errorf("expected empty results for nil path list")
	}
}

func TestScanText_BOMCharacter(t *testing.T) {
	// U+FEFF zero-width no-break space / BOM
	findings := contentscanner.ScanText("f.md", "\uFEFFcontent")
	if len(findings) == 0 {
		t.Fatal("expected finding for BOM character")
	}
	if findings[0].Category != "zero-width" {
		t.Errorf("expected zero-width category, got %q", findings[0].Category)
	}
}
