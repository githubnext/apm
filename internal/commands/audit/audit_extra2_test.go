package audit

import (
	"strings"
	"testing"
)

func TestSeverityConstants_Extra2(t *testing.T) {
	if SeverityCritical == SeverityWarning {
		t.Error("critical and warning should be distinct")
	}
	if string(SeverityCritical) == "" {
		t.Error("critical severity should have non-empty string value")
	}
	if string(SeverityWarning) == "" {
		t.Error("warning severity should have non-empty string value")
	}
}

func TestAuditModeConstants_Extra2(t *testing.T) {
	if ModeContentScan == ModeCI {
		t.Error("scan and CI mode should be distinct")
	}
}

func TestScanFinding_ZeroValue_Extra2(t *testing.T) {
	var f ScanFinding
	if f.File != "" || f.Line != 0 || f.Column != 0 {
		t.Error("zero-value ScanFinding should have empty fields")
	}
}

func TestScanFinding_Fields_Extra2(t *testing.T) {
	f := ScanFinding{
		File:     "test.py",
		Line:     10,
		Column:   5,
		CharCode: 0x202E,
		CharName: "RIGHT-TO-LEFT OVERRIDE",
		Severity: SeverityCritical,
	}
	if f.File != "test.py" {
		t.Errorf("File = %q, want test.py", f.File)
	}
	if f.Line != 10 {
		t.Errorf("Line = %d, want 10", f.Line)
	}
	if f.Severity != SeverityCritical {
		t.Errorf("Severity = %q, want critical", f.Severity)
	}
}

func TestAuditConfig_ZeroValue_Extra2(t *testing.T) {
	var cfg AuditConfig
	if cfg.ProjectRoot != "" || cfg.Verbose {
		t.Error("zero-value AuditConfig should have empty fields")
	}
}

func TestAuditConfig_Fields_Extra2(t *testing.T) {
	cfg := AuditConfig{
		ProjectRoot:  "/my/project",
		Verbose:      true,
		OutputFormat: "json",
		OutputPath:   "/tmp/out.json",
	}
	if cfg.ProjectRoot != "/my/project" {
		t.Errorf("ProjectRoot = %q", cfg.ProjectRoot)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("OutputFormat = %q", cfg.OutputFormat)
	}
}

func TestScanOptions_InheritsAuditConfig_Extra2(t *testing.T) {
	opts := ScanOptions{
		AuditConfig: AuditConfig{ProjectRoot: "/proj", OutputFormat: "text"},
		Strip:       true,
	}
	if opts.ProjectRoot != "/proj" {
		t.Errorf("embedded ProjectRoot = %q", opts.ProjectRoot)
	}
	if !opts.Strip {
		t.Error("Strip should be true")
	}
}

func TestCIOptions_InheritsAuditConfig_Extra2(t *testing.T) {
	opts := CIOptions{
		AuditConfig: AuditConfig{ProjectRoot: "/proj"},
		FailFast:    true,
	}
	if opts.ProjectRoot != "/proj" {
		t.Errorf("embedded ProjectRoot = %q", opts.ProjectRoot)
	}
	if !opts.FailFast {
		t.Error("FailFast should be true")
	}
}

func TestAuditOutcomeCause_EmptyOutcome_Extra2(t *testing.T) {
	result := AuditOutcomeCause("", "", "")
	_ = result // should not panic
}

func TestAuditOutcomeCause_AbsentOutcome_Extra2(t *testing.T) {
	result := AuditOutcomeCause("absent", "lockfile", "")
	if result == "" {
		t.Error("expected non-empty cause for absent outcome")
	}
}

func TestNew_ReturnsRunner_Extra2(t *testing.T) {
	cfg := AuditConfig{ProjectRoot: "/tmp"}
	r := New(cfg)
	if r == nil {
		t.Error("New returned nil")
	}
}

func TestScanBytes_EmptyInput_Extra2(t *testing.T) {
	s := ContentScanner{}
	findings := s.ScanBytes("empty.py", []byte{})
	if findings == nil {
		findings = []ScanFinding{}
	}
	if len(findings) != 0 {
		t.Errorf("empty content should yield no findings, got %d", len(findings))
	}
}

func TestRenderSummary_NoFindings_Extra2(t *testing.T) {
	result := &ScanResult{
		FindingsByFile: map[string][]ScanFinding{},
		FilesScanned:   5,
	}
	summary := RenderSummary(result)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestRenderFindingsTable_Extra2(t *testing.T) {
	result := &ScanResult{
		FindingsByFile: map[string][]ScanFinding{
			"file.py": {
				{File: "file.py", Line: 1, Column: 2, CharCode: 0x200B, CharName: "ZERO WIDTH SPACE", Severity: SeverityWarning},
			},
		},
		FilesScanned: 1,
		HasWarnings:  true,
	}
	table := RenderFindingsTable(result)
	if !strings.Contains(table, "file.py") {
		t.Errorf("expected file.py in table output, got %q", table)
	}
}

func TestCIFinding_Fields_Extra2(t *testing.T) {
	f := CIFinding{
		Outcome: "missing",
		Source:  "npm",
		ErrText: "not found",
		Level:   "block",
	}
	if f.Level != "block" {
		t.Errorf("Level = %q, want block", f.Level)
	}
}

func TestCIAuditResult_ZeroValue_Extra2(t *testing.T) {
	var r CIAuditResult
	if len(r.Findings) != 0 || r.ExitCode != 0 {
		t.Error("zero-value CIAuditResult should have empty fields")
	}
}

func TestLockfilePackage_Fields_Extra2(t *testing.T) {
	p := LockfilePackage{
		Name:    "pkg",
		Version: "1.2.3",
		Path:    "/lockfile",
	}
	if p.Name != "pkg" || p.Version != "1.2.3" {
		t.Errorf("LockfilePackage fields not set: %+v", p)
	}
}
