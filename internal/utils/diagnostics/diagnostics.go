// Package diagnostics provides a collect-then-render diagnostic reporting system.
//
// Integrators push diagnostics during install (or any command), and the
// collector renders a clean, grouped summary at the end. Thread-safe.
package diagnostics

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Category constants for diagnostic grouping.
const (
	CategoryCollision = "collision"
	CategoryOverwrite = "overwrite"
	CategoryWarning   = "warning"
	CategoryError     = "error"
	CategorySecurity  = "security"
	CategoryPolicy    = "policy"
	CategoryAuth      = "auth"
	CategoryDrift     = "drift"
	CategoryInfo      = "info"
)

// Drift severity constants.
const (
	DriftModified      = "modified"
	DriftUnintegrated  = "unintegrated"
	DriftOrphaned      = "orphaned"
)

var categoryOrder = []string{
	CategorySecurity,
	CategoryPolicy,
	CategoryAuth,
	CategoryDrift,
	CategoryCollision,
	CategoryOverwrite,
	CategoryWarning,
	CategoryError,
	CategoryInfo,
}

// Diagnostic is a single diagnostic message produced during an operation.
type Diagnostic struct {
	Message  string
	Category string
	Package  string
	Detail   string
	Severity string // "critical", "warning", "info" -- used by security category
}

// DiagnosticCollector collects diagnostics during a multi-package operation
// and renders a grouped summary at the end. Thread-safe.
type DiagnosticCollector struct {
	verbose     bool
	diagnostics []Diagnostic
	mu          sync.Mutex
	Out         io.Writer
}

// New creates a new DiagnosticCollector.
func New(verbose bool) *DiagnosticCollector {
	return &DiagnosticCollector{verbose: verbose, Out: os.Stdout}
}

// Skip records a collision skip (file exists, not managed by APM).
func (d *DiagnosticCollector) Skip(path, pkg string) {
	d.add(Diagnostic{Message: path, Category: CategoryCollision, Package: pkg})
}

// Overwrite records a sub-skill or file overwrite.
func (d *DiagnosticCollector) Overwrite(path, pkg, detail string) {
	d.add(Diagnostic{Message: path, Category: CategoryOverwrite, Package: pkg, Detail: detail})
}

// Warn records a general warning.
func (d *DiagnosticCollector) Warn(message, pkg, detail string) {
	d.add(Diagnostic{Message: message, Category: CategoryWarning, Package: pkg, Detail: detail})
}

// Error records an error (download failure, integration failure, etc.).
func (d *DiagnosticCollector) Error(message, pkg, detail string) {
	d.add(Diagnostic{Message: message, Category: CategoryError, Package: pkg, Detail: detail})
}

// Security records a security finding (hidden characters, etc.).
func (d *DiagnosticCollector) Security(message, pkg, detail, severity string) {
	if severity == "" {
		severity = "warning"
	}
	d.add(Diagnostic{Message: message, Category: CategorySecurity, Package: pkg, Detail: detail, Severity: severity})
}

// Info records an informational hint (non-blocking, actionable guidance).
func (d *DiagnosticCollector) Info(message, pkg, detail string) {
	d.add(Diagnostic{Message: message, Category: CategoryInfo, Package: pkg, Detail: detail})
}

// Policy records a policy enforcement finding.
func (d *DiagnosticCollector) Policy(message, pkg, detail string) {
	d.add(Diagnostic{Message: message, Category: CategoryPolicy, Package: pkg, Detail: detail})
}

// Auth records an authentication issue.
func (d *DiagnosticCollector) Auth(message, pkg, detail string) {
	d.add(Diagnostic{Message: message, Category: CategoryAuth, Package: pkg, Detail: detail})
}

// Drift records a drift finding.
func (d *DiagnosticCollector) Drift(message, pkg, detail string) {
	d.add(Diagnostic{Message: message, Category: CategoryDrift, Package: pkg, Detail: detail})
}

// HasDiagnostics returns true if any diagnostics have been recorded.
func (d *DiagnosticCollector) HasDiagnostics() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.diagnostics) > 0
}

// HasErrors returns true if any error diagnostics have been recorded.
func (d *DiagnosticCollector) HasErrors() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, diag := range d.diagnostics {
		if diag.Category == CategoryError || diag.Category == CategorySecurity {
			return true
		}
	}
	return false
}

// All returns a copy of all collected diagnostics.
func (d *DiagnosticCollector) All() []Diagnostic {
	d.mu.Lock()
	defer d.mu.Unlock()
	result := make([]Diagnostic, len(d.diagnostics))
	copy(result, d.diagnostics)
	return result
}

// RenderSummary prints a grouped summary of all diagnostics.
func (d *DiagnosticCollector) RenderSummary() {
	d.mu.Lock()
	diags := make([]Diagnostic, len(d.diagnostics))
	copy(diags, d.diagnostics)
	d.mu.Unlock()

	if len(diags) == 0 {
		return
	}

	// Group by category
	grouped := make(map[string][]Diagnostic)
	for _, diag := range diags {
		grouped[diag.Category] = append(grouped[diag.Category], diag)
	}

	headers := map[string]string{
		CategorySecurity:  "[!] Security findings",
		CategoryPolicy:    "[!] Policy enforcement",
		CategoryAuth:      "[!] Authentication issues",
		CategoryDrift:     "[!] Drift detected",
		CategoryCollision: "[!] File collisions (skipped)",
		CategoryOverwrite: "[~] Overwrites",
		CategoryWarning:   "[!] Warnings",
		CategoryError:     "[x] Errors",
		CategoryInfo:      "[i] Notes",
	}

	out := d.Out
	for _, cat := range categoryOrder {
		items, ok := grouped[cat]
		if !ok || len(items) == 0 {
			continue
		}
		header, ok := headers[cat]
		if !ok {
			header = "[i] " + cat
		}
		fmt.Fprintf(out, "\n%s (%d):\n", header, len(items))
		for _, item := range items {
			line := "  - " + item.Message
			if item.Package != "" {
				line += " [" + item.Package + "]"
			}
			if item.Detail != "" && d.verbose {
				line += "\n    " + item.Detail
			}
			fmt.Fprintln(out, line)
		}
	}
}

func (d *DiagnosticCollector) add(diag Diagnostic) {
	d.mu.Lock()
	d.diagnostics = append(d.diagnostics, diag)
	d.mu.Unlock()
}
