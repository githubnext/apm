// Package gate provides a centralized security scanning gate for APM commands.
// It mirrors src/apm_cli/security/gate.py.
//
// Every command that reads or writes files passes through Gate instead of
// reimplementing scan->classify->decide->report inline. Commands declare
// intent via ScanPolicy; the gate handles the rest.
package gate

import (
	"github.com/githubnext/apm/internal/security/contentscanner"
)

// OnCritical controls how the gate responds to critical findings.
type OnCritical string

const (
	OnCriticalBlock  OnCritical = "block"
	OnCriticalWarn   OnCritical = "warn"
	OnCriticalIgnore OnCritical = "ignore"
)

// ScanPolicy declares how a command handles security findings.
type ScanPolicy struct {
	OnCritical     OnCritical
	ForceOverrides bool // when true, --force downgrades block to warn
}

var (
	// BlockPolicy blocks deployment on critical findings (default).
	BlockPolicy = ScanPolicy{OnCritical: OnCriticalBlock, ForceOverrides: true}
	// WarnPolicy continues with a warning on critical findings.
	WarnPolicy = ScanPolicy{OnCritical: OnCriticalWarn, ForceOverrides: false}
	// ReportPolicy collects findings silently.
	ReportPolicy = ScanPolicy{OnCritical: OnCriticalIgnore, ForceOverrides: false}
)

// EffectiveBlock returns true when this policy would block deployment.
func (p ScanPolicy) EffectiveBlock(force bool) bool {
	return p.OnCritical == OnCriticalBlock && !(p.ForceOverrides && force)
}

// ScanVerdict is the result of a Gate check.
type ScanVerdict struct {
	FindingsByFile map[string][]contentscanner.ScanFinding
	HasCritical    bool
	ShouldBlock    bool
	CriticalCount  int
	WarningCount   int
	FilesScanned   int
}

// HasFindings returns true when any findings were recorded.
func (v ScanVerdict) HasFindings() bool {
	return len(v.FindingsByFile) > 0
}

// Gate wraps a ContentScanner and applies a ScanPolicy.
type Gate struct {
	scanner *contentscanner.ContentScanner
	policy  ScanPolicy
	force   bool
}

// New creates a new Gate with the given policy and force flag.
func New(policy ScanPolicy, force bool) *Gate {
	return &Gate{
		scanner: contentscanner.NewDefaultScanner(),
		policy:  policy,
		force:   force,
	}
}

// Check scans the provided file paths and returns a ScanVerdict.
func (g *Gate) Check(paths []string) ScanVerdict {
	findingsByFile := g.scanner.ScanFiles(paths)

	verdict := ScanVerdict{
		FindingsByFile: findingsByFile,
		FilesScanned:   len(paths),
	}

	for _, findings := range findingsByFile {
		for _, f := range findings {
			switch f.Severity {
			case "critical":
				verdict.HasCritical = true
				verdict.CriticalCount++
			case "warning":
				verdict.WarningCount++
			}
		}
	}

	if verdict.HasCritical {
		verdict.ShouldBlock = g.policy.EffectiveBlock(g.force)
	}
	return verdict
}

// CheckFile is a convenience wrapper for a single file.
func (g *Gate) CheckFile(path string) ScanVerdict {
	return g.Check([]string{path})
}
