// Package security provides types for the security scanning gate.
// Mirrors Python apm_cli.security.gate.
package security

// ScanFinding represents a single security finding from content scanning.
type ScanFinding struct {
	Rule     string // rule identifier
	Severity string // "critical", "warning", "info"
	Message  string // human-readable description
	Line     int    // line number (0 = unknown)
	Column   int    // column (0 = unknown)
}

// ScanPolicy declares how a command handles security findings.
type ScanPolicy struct {
	OnCritical     string // "block", "warn", "ignore"
	ForceOverrides bool   // when true, --force downgrades block to warn
}

// BlockPolicy blocks on critical findings; force can override.
var BlockPolicy = ScanPolicy{OnCritical: "block", ForceOverrides: true}

// WarnPolicy warns on critical findings; force has no effect.
var WarnPolicy = ScanPolicy{OnCritical: "warn", ForceOverrides: false}

// ReportPolicy collects findings silently.
var ReportPolicy = ScanPolicy{OnCritical: "ignore", ForceOverrides: false}

// EffectiveBlock returns true when this policy would block deployment.
func (p ScanPolicy) EffectiveBlock(force bool) bool {
	return p.OnCritical == "block" && !(p.ForceOverrides && force)
}

// ScanVerdict is the result of a SecurityGate check.
type ScanVerdict struct {
	FindingsByFile map[string][]ScanFinding
	HasCritical    bool
	ShouldBlock    bool
	CriticalCount  int
	WarningCount   int
	FilesScanned   int
}

// HasFindings returns true if any file has findings.
func (v *ScanVerdict) HasFindings() bool {
	return len(v.FindingsByFile) > 0
}

// AllFindings returns a flat list of all findings across files.
func (v *ScanVerdict) AllFindings() []ScanFinding {
	var out []ScanFinding
	for _, findings := range v.FindingsByFile {
		out = append(out, findings...)
	}
	return out
}

// ContentPattern is a compiled pattern used by the content scanner.
type ContentPattern struct {
	ID       string
	Pattern  string
	Severity string
	Message  string
}

// DefaultContentPatterns returns the built-in security patterns.
func DefaultContentPatterns() []ContentPattern {
	return []ContentPattern{
		{ID: "hardcoded-secret", Pattern: `(?i)(password|passwd|secret|token|api_key)\s*=\s*["'][^"']{8,}["']`, Severity: "critical", Message: "Possible hardcoded secret"},
		{ID: "private-key", Pattern: `-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY`, Severity: "critical", Message: "Private key material"},
		{ID: "aws-key", Pattern: `AKIA[0-9A-Z]{16}`, Severity: "critical", Message: "Possible AWS access key"},
		{ID: "github-token", Pattern: `gh[pousr]_[A-Za-z0-9_]{36,}`, Severity: "critical", Message: "Possible GitHub token"},
	}
}

// AuditSeverity is the severity level for audit report entries.
type AuditSeverity string

const (
	SeverityCritical AuditSeverity = "critical"
	SeverityHigh     AuditSeverity = "high"
	SeverityMedium   AuditSeverity = "medium"
	SeverityLow      AuditSeverity = "low"
	SeverityInfo     AuditSeverity = "info"
)

// AuditReportEntry represents one finding in an audit report.
type AuditReportEntry struct {
	CheckName string
	Severity  AuditSeverity
	Message   string
	Details   []string
	File      string
}

// AuditReport is the top-level output of a security audit.
type AuditReport struct {
	Entries []AuditReportEntry
}

// Critical returns all critical-severity entries.
func (r *AuditReport) Critical() []AuditReportEntry {
	var out []AuditReportEntry
	for _, e := range r.Entries {
		if e.Severity == SeverityCritical {
			out = append(out, e)
		}
	}
	return out
}

// HasBlockers returns true if there are any critical or high findings.
func (r *AuditReport) HasBlockers() bool {
	for _, e := range r.Entries {
		if e.Severity == SeverityCritical || e.Severity == SeverityHigh {
			return true
		}
	}
	return false
}
