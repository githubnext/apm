package audit

import (
	"strings"
	"testing"
)

func TestScanFinding_ZeroValue_Extra3(t *testing.T) {
	var f ScanFinding
	if f.File != "" || f.Line != 0 || f.Column != 0 {
		t.Error("zero ScanFinding should have empty fields")
	}
}

func TestScanFinding_AssignFields_Extra3(t *testing.T) {
	f := ScanFinding{
		File:     "test.md",
		Line:     5,
		Column:   10,
		CharCode: 0x200B,
		CharName: "ZERO WIDTH SPACE",
		Severity: SeverityCritical,
	}
	if f.File != "test.md" {
		t.Errorf("expected File=test.md, got %q", f.File)
	}
	if f.Line != 5 {
		t.Errorf("expected Line=5, got %d", f.Line)
	}
	if f.Severity != SeverityCritical {
		t.Errorf("expected critical severity")
	}
}

func TestAuditConfig_ZeroValue_Extra3(t *testing.T) {
	var cfg AuditConfig
	if cfg.ProjectRoot != "" || cfg.Verbose || cfg.OutputFormat != "" {
		t.Error("zero AuditConfig should have zero fields")
	}
}

func TestScanOptions_ZeroValue_Extra3(t *testing.T) {
	var opts ScanOptions
	if opts.Strip || opts.Preview || opts.MaxFindings != 0 {
		t.Error("zero ScanOptions should have zero fields")
	}
}

func TestAuditMode_ContentScanValue_Extra3(t *testing.T) {
	if string(ModeContentScan) == "" {
		t.Error("ModeContentScan must not be empty")
	}
	if string(ModeCI) == "" {
		t.Error("ModeCI must not be empty")
	}
	if string(ModeDrift) == "" {
		t.Error("ModeDrift must not be empty")
	}
}

func TestAuditMode_AllDistinct_Extra3(t *testing.T) {
	modes := []AuditMode{ModeContentScan, ModeCI, ModeDrift}
	seen := map[AuditMode]bool{}
	for _, m := range modes {
		if seen[m] {
			t.Errorf("duplicate mode: %q", m)
		}
		seen[m] = true
	}
}

func TestSeverity_InformationalValue_Extra3(t *testing.T) {
	if string(SeverityInfo) == "" {
		t.Error("SeverityInfo must not be empty")
	}
	if SeverityInfo == SeverityCritical {
		t.Error("info and critical must be distinct")
	}
	if SeverityInfo == SeverityWarning {
		t.Error("info and warning must be distinct")
	}
}

func TestScanFinding_ContextField_Extra3(t *testing.T) {
	f := ScanFinding{Context: "some surrounding text"}
	if !strings.Contains(f.Context, "surrounding") {
		t.Error("Context field should store surrounding text")
	}
}
