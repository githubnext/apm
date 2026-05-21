package audit

import "testing"

func TestScanFinding_SeverityField_Extra4(t *testing.T) {
f := ScanFinding{Severity: SeverityCritical}
if f.Severity != SeverityCritical {
t.Errorf("unexpected severity: %s", f.Severity)
}
}

func TestScanFinding_LineField_Extra4(t *testing.T) {
f := ScanFinding{Line: 42}
if f.Line != 42 {
t.Errorf("unexpected line: %d", f.Line)
}
}

func TestScanFinding_ColumnField_Extra4(t *testing.T) {
f := ScanFinding{Column: 5}
if f.Column != 5 {
t.Errorf("unexpected column: %d", f.Column)
}
}

func TestAuditConfig_ProjectRootField_Extra4(t *testing.T) {
cfg := AuditConfig{ProjectRoot: "/my/project"}
if cfg.ProjectRoot != "/my/project" {
t.Errorf("unexpected project root: %s", cfg.ProjectRoot)
}
}

func TestAuditConfig_VerboseField_Extra4(t *testing.T) {
cfg := AuditConfig{Verbose: true}
if !cfg.Verbose {
t.Error("expected Verbose true")
}
}

func TestAuditConfig_OutputFormatField_Extra4(t *testing.T) {
cfg := AuditConfig{OutputFormat: "json"}
if cfg.OutputFormat != "json" {
t.Errorf("unexpected output format: %s", cfg.OutputFormat)
}
}

func TestContentScanner_ScanBytes_ZeroWidth_Extra4(t *testing.T) {
s := ContentScanner{}
content := []byte("normal text​with zero-width space")
findings := s.ScanBytes("test.txt", content)
if len(findings) == 0 {
t.Error("expected finding for zero-width space")
}
}

func TestContentScanner_ScanBytes_Clean_Extra4(t *testing.T) {
s := ContentScanner{}
content := []byte("just regular ASCII text here")
findings := s.ScanBytes("test.txt", content)
if len(findings) != 0 {
t.Errorf("expected no findings, got %d", len(findings))
}
}

func TestSeverity_DistinctValues_Extra4(t *testing.T) {
if SeverityWarning == SeverityCritical {
t.Error("expected Warning and Critical to be distinct")
}
if SeverityWarning == SeverityInfo {
t.Error("expected Warning and Info to be distinct")
}
}

func TestAuditMode_ContentScan_Extra4(t *testing.T) {
m := ModeContentScan
if m == ModeCI {
t.Error("expected content scan and CI to be distinct")
}
}

func TestScanFinding_CharNameField_Extra4(t *testing.T) {
f := ScanFinding{CharName: "ZERO WIDTH SPACE"}
if f.CharName != "ZERO WIDTH SPACE" {
t.Errorf("unexpected CharName: %s", f.CharName)
}
}
