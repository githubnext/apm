package auditreport_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/security/auditreport"
)

func sampleFindings() map[string][]auditreport.ScanFinding {
	return map[string][]auditreport.ScanFinding{
		"test.md": {
			{
				Severity:    "critical",
				File:        "test.md",
				Line:        3,
				Column:      5,
				Codepoint:   "U+202E",
				Category:    "bidi-override",
				Description: "Right-to-left override",
			},
		},
	}
}

func TestRelativePathForReport_Relative(t *testing.T) {
	got := auditreport.RelativePathForReport("src/foo.md")
	if got != "src/foo.md" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestRelativePathForReport_BackslashNormalized(t *testing.T) {
	got := auditreport.RelativePathForReport("src\\foo.md")
	if strings.Contains(got, "\\") {
		t.Errorf("backslashes not normalized: %q", got)
	}
}

func TestFindingsToJSON_Structure(t *testing.T) {
	report := auditreport.FindingsToJSON(sampleFindings(), 5, 1)
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if _, ok := report["findings"]; !ok {
		t.Error("expected 'findings' key in JSON report")
	}
	summary, ok := report["summary"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'summary' key in JSON report")
	}
	if _, ok := summary["files_scanned"]; !ok {
		t.Error("expected 'files_scanned' key in summary")
	}
}

func TestFindingsToSARIF_Structure(t *testing.T) {
	report := auditreport.FindingsToSARIF(sampleFindings(), 5)
	if report == nil {
		t.Fatal("expected non-nil SARIF report")
	}
	version, ok := report["version"].(string)
	if !ok || version != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %v", report["version"])
	}
}

func TestFindingsToMarkdown_ContainsFile(t *testing.T) {
	md := auditreport.FindingsToMarkdown(sampleFindings(), 5)
	if !strings.Contains(md, "test.md") {
		t.Errorf("expected 'test.md' in markdown output")
	}
	if !strings.Contains(md, "critical") {
		t.Errorf("expected 'critical' in markdown output")
	}
}

func TestFindingsToMarkdown_NoFindings(t *testing.T) {
	md := auditreport.FindingsToMarkdown(map[string][]auditreport.ScanFinding{}, 10)
	if md == "" {
		t.Fatal("expected non-empty markdown even for no findings")
	}
}

func TestSerializeReport_JSON(t *testing.T) {
	report := map[string]interface{}{"key": "val"}
	out, err := auditreport.SerializeReport(report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "val") {
		t.Errorf("expected serialized output to contain 'val'")
	}
}

func TestDetectFormatFromExtension(t *testing.T) {
	cases := map[string]string{
		"report.json":  "json",
		"report.sarif": "sarif",
		"report.md":    "markdown",
		"report.txt":   "text",
	}
	for path, want := range cases {
		got := auditreport.DetectFormatFromExtension(path)
		if got != want {
			t.Errorf("DetectFormatFromExtension(%q) = %q, want %q", path, got, want)
		}
	}
}
