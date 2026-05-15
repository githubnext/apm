package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanBytes_Clean(t *testing.T) {
	s := ContentScanner{}
	findings := s.ScanBytes("test.md", []byte("Hello, world! Normal text here."))
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestScanBytes_BiDiOverride(t *testing.T) {
	s := ContentScanner{}
	// 0x202E is RIGHT-TO-LEFT OVERRIDE -- critical severity
	data := []byte("before\xe2\x80\xaeafter")
	findings := s.ScanBytes("test.md", data)
	if len(findings) == 0 {
		t.Fatal("expected findings for bidi override character")
	}
	if findings[0].Severity != SeverityCritical {
		t.Errorf("expected critical, got %s", findings[0].Severity)
	}
}

func TestScanBytes_ZeroWidth(t *testing.T) {
	s := ContentScanner{}
	// 0x200B is ZERO WIDTH SPACE -- warning
	data := []byte("hello\xe2\x80\x8bworld")
	findings := s.ScanBytes("test.md", data)
	if len(findings) == 0 {
		t.Fatal("expected findings for zero-width space")
	}
	if findings[0].Severity != SeverityWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestScanBytes_LineColumn(t *testing.T) {
	s := ContentScanner{}
	data := []byte("line1\nline2\xe2\x80\x8bafter")
	findings := s.ScanBytes("f.md", data)
	if len(findings) == 0 {
		t.Fatal("no findings")
	}
	if findings[0].Line != 2 {
		t.Errorf("expected line 2, got %d", findings[0].Line)
	}
}

func TestScanFile_NotFound(t *testing.T) {
	s := ContentScanner{}
	_, err := s.ScanFile("/nonexistent/path.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestScanFile_Clean(t *testing.T) {
	f, err := os.CreateTemp("", "audit_test_*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("Clean content with no hidden chars.")
	f.Close()

	s := ContentScanner{}
	findings, err := s.ScanFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestRunnerRun_NoFiles(t *testing.T) {
	dir := t.TempDir()
	r := New(AuditConfig{ProjectRoot: dir})
	result, err := r.Run(ScanOptions{AuditConfig: AuditConfig{ProjectRoot: dir}})
	if err != nil {
		t.Fatal(err)
	}
	if result.FilesScanned != 0 {
		t.Errorf("expected 0 scanned, got %d", result.FilesScanned)
	}
}

func TestRunnerRun_WithCleanFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "readme.md")
	os.WriteFile(f, []byte("# Hello World\n\nClean content."), 0644)

	r := New(AuditConfig{ProjectRoot: dir})
	result, err := r.Run(ScanOptions{
		AuditConfig: AuditConfig{ProjectRoot: dir},
		Files:       []string{f},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.FilesScanned != 1 {
		t.Errorf("expected 1 scanned, got %d", result.FilesScanned)
	}
	if result.HasCritical || result.HasWarnings {
		t.Error("expected clean result")
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
}

func TestRunnerRun_WithCriticalFinding(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad.md")
	// embed bidi override
	os.WriteFile(f, []byte("text\xe2\x80\xaeevil"), 0644)

	r := New(AuditConfig{ProjectRoot: dir})
	result, err := r.Run(ScanOptions{
		AuditConfig: AuditConfig{ProjectRoot: dir},
		Files:       []string{f},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasCritical {
		t.Error("expected critical finding")
	}
	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
}

func TestRenderSummary_Clean(t *testing.T) {
	result := &ScanResult{FilesScanned: 5}
	s := RenderSummary(result)
	if !strings.Contains(s, "Clean") {
		t.Errorf("expected 'Clean' in summary, got: %s", s)
	}
}

func TestRenderSummary_Critical(t *testing.T) {
	result := &ScanResult{
		HasCritical:    true,
		FindingsByFile: map[string][]ScanFinding{"f.md": {{}}},
	}
	s := RenderSummary(result)
	if !strings.Contains(s, "Critical") {
		t.Errorf("expected 'Critical' in summary, got: %s", s)
	}
}

func TestRenderSummary_Warning(t *testing.T) {
	result := &ScanResult{
		HasWarnings:    true,
		FindingsByFile: map[string][]ScanFinding{"f.md": {{}}},
	}
	s := RenderSummary(result)
	if !strings.Contains(s, "Warning") {
		t.Errorf("expected 'Warning' in summary, got: %s", s)
	}
}

func TestRenderFindingsTable_Empty(t *testing.T) {
	result := &ScanResult{FindingsByFile: map[string][]ScanFinding{}}
	out := RenderFindingsTable(result)
	if !strings.Contains(out, "No hidden") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestRenderFindingsJSON(t *testing.T) {
	result := &ScanResult{
		FindingsByFile: map[string][]ScanFinding{
			"f.md": {{File: "f.md", Line: 1, CharCode: 0x202E, Severity: SeverityCritical}},
		},
	}
	out, err := RenderFindingsJSON(result)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "f.md") {
		t.Error("expected file name in JSON output")
	}
}

func TestAuditOutcomeCause_NoGitRemote(t *testing.T) {
	s := AuditOutcomeCause("no_git_remote", "", "")
	if !strings.Contains(s, "org from git remote") {
		t.Errorf("unexpected: %s", s)
	}
}

func TestAuditOutcomeCause_Absent(t *testing.T) {
	s := AuditOutcomeCause("absent", "https://example.com/policy", "")
	if !strings.Contains(s, "No org policy") {
		t.Errorf("unexpected: %s", s)
	}
}

func TestAuditOutcomeCause_Empty(t *testing.T) {
	s := AuditOutcomeCause("empty", "https://example.com/policy", "")
	if !strings.Contains(s, "empty") {
		t.Errorf("unexpected: %s", s)
	}
}

func TestStripFindings_DryRun(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad.md")
	content := "text\xe2\x80\x8bhidden"
	os.WriteFile(f, []byte(content), 0644)

	findings := map[string][]ScanFinding{
		f: {{File: f, Severity: SeverityWarning}},
	}
	results, err := StripFindings(findings, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 strip result, got %d", len(results))
	}
	// dry run: file should be unchanged
	data, _ := os.ReadFile(f)
	if string(data) != content {
		t.Error("dry run should not modify file")
	}
}

func TestStripFindings_Live(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "clean.md")
	os.WriteFile(f, []byte("hello\xe2\x80\x8bworld"), 0644)

	findings := map[string][]ScanFinding{
		f: {{File: f, Severity: SeverityWarning}},
	}
	results, err := StripFindings(findings, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected strip results")
	}
	data, _ := os.ReadFile(f)
	if strings.Contains(string(data), "\xe2\x80\x8b") {
		t.Error("hidden char should have been removed")
	}
}

func TestScanLockfilePackages(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "apm.lock.yaml")
	os.WriteFile(f, []byte("packages:\n  - name: foo\n    version: 1.0\n"), 0644)

	s := ContentScanner{}
	result, err := ScanLockfilePackages(f, s)
	if err != nil {
		t.Fatal(err)
	}
	if result.FilesScanned > 1 {
		t.Error("unexpected file count")
	}
}
