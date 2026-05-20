package auditreport

import (
	"strings"
	"testing"
)

func TestScanFinding_ZeroValue(t *testing.T) {
	var f ScanFinding
	if f.Severity != "" || f.File != "" || f.Line != 0 {
		t.Error("zero value should have empty fields")
	}
}

func TestScanFinding_Fields(t *testing.T) {
	f := ScanFinding{
		Severity:    "critical",
		File:        "src/foo.py",
		Line:        42,
		Column:      7,
		Codepoint:   "U+200B",
		Category:    "zero-width",
		Description: "zero-width space",
	}
	if f.Severity != "critical" {
		t.Errorf("Severity = %q, want critical", f.Severity)
	}
	if f.Line != 42 {
		t.Errorf("Line = %d, want 42", f.Line)
	}
}

func TestRelativePathForReport_ForwardSlash(t *testing.T) {
	p := RelativePathForReport("src/foo.py")
	if strings.Contains(p, "\\") {
		t.Errorf("expected forward slashes, got %q", p)
	}
}

func TestDetectFormatFromExtension_JSON(t *testing.T) {
	if DetectFormatFromExtension("report.json") != "json" {
		t.Error("expected json for .json")
	}
}

func TestDetectFormatFromExtension_SARIF(t *testing.T) {
	if DetectFormatFromExtension("report.sarif") != "sarif" {
		t.Error("expected sarif for .sarif")
	}
}

func TestDetectFormatFromExtension_MD(t *testing.T) {
	if DetectFormatFromExtension("report.md") != "markdown" {
		t.Error("expected markdown for .md")
	}
}

func TestFindingsToMarkdown_EmptyClean(t *testing.T) {
	md := FindingsToMarkdown(nil, 5)
	if md == "" {
		t.Error("expected non-empty markdown even for clean scan")
	}
}

func TestFindingsToJSON_ExitCode(t *testing.T) {
	result := FindingsToJSON(nil, 3, 0)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestFindingsToSARIF_NonNil(t *testing.T) {
	result := FindingsToSARIF(nil, 0)
	if result == nil {
		t.Fatal("expected non-nil SARIF result")
	}
}
