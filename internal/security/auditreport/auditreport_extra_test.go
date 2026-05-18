package auditreport_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/security/auditreport"
)

func TestRelativePathForReport_relative(t *testing.T) {
	got := auditreport.RelativePathForReport("src/foo/bar.md")
	if got != "src/foo/bar.md" {
		t.Errorf("expected src/foo/bar.md, got %s", got)
	}
}

func TestRelativePathForReport_backslash(t *testing.T) {
	got := auditreport.RelativePathForReport(`src\foo\bar.md`)
	if strings.Contains(got, `\`) {
		t.Errorf("expected forward slashes, got %s", got)
	}
}

func TestFindingsToJSON_empty(t *testing.T) {
	result := auditreport.FindingsToJSON(nil, 10, 0)
	if result["exit_code"] != 0 {
		t.Errorf("expected exit_code 0")
	}
	summary := result["summary"].(map[string]interface{})
	if summary["files_scanned"].(int) != 10 {
		t.Errorf("expected 10 files_scanned")
	}
	if summary["critical"].(int) != 0 {
		t.Errorf("expected 0 critical")
	}
}

func TestFindingsToJSON_mixed(t *testing.T) {
	findings := map[string][]auditreport.ScanFinding{
		"foo.md": {
			{Severity: "critical", File: "foo.md", Line: 1, Column: 5, Codepoint: "U+200B", Category: "zero-width", Description: "ZWSP"},
			{Severity: "warning", File: "foo.md", Line: 2, Column: 1, Codepoint: "U+200C", Category: "zero-width", Description: "ZWNJ"},
		},
		"bar.md": {
			{Severity: "info", File: "bar.md", Line: 10, Column: 2, Codepoint: "U+00AD", Category: "soft-hyphen", Description: "SHY"},
		},
	}
	result := auditreport.FindingsToJSON(findings, 5, 1)
	summary := result["summary"].(map[string]interface{})
	if summary["critical"].(int) != 1 {
		t.Errorf("expected 1 critical")
	}
	if summary["warning"].(int) != 1 {
		t.Errorf("expected 1 warning")
	}
	if summary["info"].(int) != 1 {
		t.Errorf("expected 1 info")
	}
	if summary["files_affected"].(int) != 2 {
		t.Errorf("expected 2 files_affected")
	}
}

func TestFindingsToSARIF_empty(t *testing.T) {
	result := auditreport.FindingsToSARIF(nil, 0)
	if result["version"] != "2.1.0" {
		t.Errorf("expected SARIF 2.1.0, got %v", result["version"])
	}
}

func TestFindingsToSARIF_withFindings(t *testing.T) {
	findings := map[string][]auditreport.ScanFinding{
		"a.md": {
			{Severity: "critical", File: "a.md", Line: 1, Column: 1, Codepoint: "U+200B", Category: "zero-width", Description: "ZWSP"},
		},
	}
	result := auditreport.FindingsToSARIF(findings, 1)
	runs := result["runs"].([]interface{})
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	run := runs[0].(map[string]interface{})
	results := run["results"].([]interface{})
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestFindingsToMarkdown_clean(t *testing.T) {
	md := auditreport.FindingsToMarkdown(nil, 10)
	if !strings.Contains(md, "Clean") {
		t.Errorf("expected clean message, got: %s", md)
	}
}

func TestFindingsToMarkdown_withFindings(t *testing.T) {
	findings := map[string][]auditreport.ScanFinding{
		"doc.md": {
			{Severity: "critical", File: "doc.md", Line: 3, Column: 7, Codepoint: "U+200B", Category: "zero-width", Description: "zero-width space"},
		},
	}
	md := auditreport.FindingsToMarkdown(findings, 5)
	if !strings.Contains(md, "doc.md") {
		t.Errorf("expected doc.md in markdown output")
	}
	if !strings.Contains(md, "CRITICAL") {
		t.Errorf("expected CRITICAL severity in output")
	}
	if !strings.Contains(md, "U+200B") {
		t.Errorf("expected codepoint in output")
	}
}

func TestDetectFormatFromExtension_variants(t *testing.T) {
	cases := []struct {
		path   string
		expect string
	}{
		{"report.SARIF.JSON", "sarif"},
		{"output.sarif", "sarif"},
		{"data.json", "json"},
		{"notes.md", "markdown"},
		{"output.txt", "text"},
		{"plain", "text"},
	}
	for _, c := range cases {
		got := auditreport.DetectFormatFromExtension(c.path)
		if got != c.expect {
			t.Errorf("DetectFormatFromExtension(%q) = %q, want %q", c.path, got, c.expect)
		}
	}
}

func TestSerializeReport(t *testing.T) {
	report := map[string]interface{}{"version": "1", "exit_code": 0}
	s, err := auditreport.SerializeReport(report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(s, `"version"`) {
		t.Errorf("expected version in serialized output")
	}
}

func TestFindingsToMarkdown_warningCount(t *testing.T) {
	findings := map[string][]auditreport.ScanFinding{
		"a.md": {
			{Severity: "warning", File: "a.md", Line: 1, Column: 1, Codepoint: "U+00AD", Category: "soft-hyphen", Description: "shy"},
			{Severity: "warning", File: "a.md", Line: 2, Column: 1, Codepoint: "U+00AD", Category: "soft-hyphen", Description: "shy2"},
		},
	}
	md := auditreport.FindingsToMarkdown(findings, 3)
	if !strings.Contains(md, "2 warnings") {
		t.Errorf("expected '2 warnings', got: %s", md)
	}
}

func TestFindingsToMarkdown_singleWarning(t *testing.T) {
	findings := map[string][]auditreport.ScanFinding{
		"b.md": {
			{Severity: "warning", File: "b.md", Line: 1, Column: 1, Codepoint: "U+00AD", Category: "soft-hyphen", Description: "shy"},
		},
	}
	md := auditreport.FindingsToMarkdown(findings, 1)
	if !strings.Contains(md, "1 warning") {
		t.Errorf("expected '1 warning', got: %s", md)
	}
	if strings.Contains(md, "1 warnings") {
		t.Errorf("should not have '1 warnings' (plural), got: %s", md)
	}
}
