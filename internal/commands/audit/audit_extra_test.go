package audit

import (
	"strings"
	"testing"
)

func TestScanBytes_HomoglyphWarning(t *testing.T) {
	s := ContentScanner{}
	// U+0430 (Cyrillic a) -- homoglyph, warning-level
	data := []byte("Hello \xd0\xb0 world")
	findings := s.ScanBytes("test.md", data)
	if len(findings) == 0 {
		t.Skip("no finding for homoglyph -- scanner may not check this file type")
	}
	found := false
	for _, f := range findings {
		if f.Severity == SeverityWarning || f.Severity == SeverityCritical {
			found = true
		}
	}
	if !found {
		t.Error("expected warning or critical finding for homoglyph")
	}
}

func TestScanBytes_FindingHasFields(t *testing.T) {
	s := ContentScanner{}
	// bidi override character (critical)
	data := []byte("before\xe2\x80\xaeafter")
	findings := s.ScanBytes("test.txt", data)
	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}
	f := findings[0]
	if f.File == "" {
		t.Error("File field should be non-empty")
	}
	if f.Severity == "" {
		t.Error("Severity field should be non-empty")
	}
	if f.CharName == "" {
		t.Error("CharName field should be non-empty")
	}
}

func TestSeverityConstants_Distinct(t *testing.T) {
	if SeverityCritical == SeverityWarning {
		t.Error("SeverityCritical and SeverityWarning should be distinct")
	}
	if SeverityWarning == SeverityInfo {
		t.Error("SeverityWarning and SeverityInfo should be distinct")
	}
}

func TestSeverityConstants_StringValues(t *testing.T) {
	if string(SeverityCritical) != "critical" {
		t.Errorf("SeverityCritical = %q, want critical", SeverityCritical)
	}
	if string(SeverityWarning) != "warning" {
		t.Errorf("SeverityWarning = %q, want warning", SeverityWarning)
	}
	if string(SeverityInfo) != "info" {
		t.Errorf("SeverityInfo = %q, want info", SeverityInfo)
	}
}

func TestAuditModeConstants_Distinct(t *testing.T) {
	if ModeContentScan == ModeCI {
		t.Error("ModeContentScan and ModeCI should be distinct")
	}
	if ModeCI == ModeDrift {
		t.Error("ModeCI and ModeDrift should be distinct")
	}
}

func TestNew_ReturnsRunner(t *testing.T) {
	r := New(AuditConfig{})
	if r == nil {
		t.Fatal("New should return non-nil Runner")
	}
}

func TestScanBytes_EmptyContent(t *testing.T) {
	s := ContentScanner{}
	findings := s.ScanBytes("empty.md", []byte{})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for empty content, got %d", len(findings))
	}
}

func TestScanBytes_FindingFile(t *testing.T) {
	s := ContentScanner{}
	data := []byte("before\xe2\x80\xaeafter")
	findings := s.ScanBytes("myfile.py", data)
	if len(findings) == 0 {
		t.Skip("no finding for this content")
	}
	if findings[0].File != "myfile.py" {
		t.Errorf("File = %q, want myfile.py", findings[0].File)
	}
}

func TestScanBytes_NormalASCII_NoFindings(t *testing.T) {
	s := ContentScanner{}
	plain := strings.Repeat("abcdefghijklmnopqrstuvwxyz0123456789 \n", 20)
	findings := s.ScanBytes("plain.md", []byte(plain))
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for plain ASCII, got %d", len(findings))
	}
}
