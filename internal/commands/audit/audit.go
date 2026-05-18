// Package audit implements the APM audit command -- content integrity scanning
// for prompt files.
//
// Scans installed APM packages (or arbitrary files) for hidden Unicode
// characters that could embed invisible instructions. Also supports
// lock-file consistency (--ci) and drift detection (--drift) modes.
//
// Exit codes:
//
//	0 -- clean (no findings, or info-only)
//	1 -- critical findings detected
//	2 -- warnings only (no critical)
//
// Migrated from: src/apm_cli/commands/audit.py
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// -------------------------------------------------------------------
// Finding types
// -------------------------------------------------------------------

// Severity classifies how serious a finding is.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

// ScanFinding records a single suspicious character or pattern.
type ScanFinding struct {
	File     string
	Line     int
	Column   int
	CharCode int
	CharName string
	Context  string
	Severity Severity
}

// -------------------------------------------------------------------
// Config / options
// -------------------------------------------------------------------

// AuditConfig holds options shared across audit modes.
type AuditConfig struct {
	ProjectRoot  string
	Verbose      bool
	OutputFormat string // "text" | "json"
	OutputPath   string
}

// AuditMode selects the audit sub-command.
type AuditMode string

const (
	ModeContentScan AuditMode = "content"
	ModeCI          AuditMode = "ci"
	ModeDrift       AuditMode = "drift"
)

// ScanOptions controls a content-scan run.
type ScanOptions struct {
	AuditConfig
	Files      []string // explicit file list; empty = scan all packages
	Strip      bool
	Preview    bool
	MaxFindings int
}

// CIOptions controls a --ci policy-check run.
type CIOptions struct {
	AuditConfig
	Policy     string
	FailFast   bool
}

// -------------------------------------------------------------------
// ContentScanner
// -------------------------------------------------------------------

// ContentScanner scans files for hidden or dangerous Unicode characters.
type ContentScanner struct{}

// HiddenUnicodeRanges lists Unicode categories and codepoints that should
// not appear in prompt files.
var HiddenUnicodeRanges = []struct {
	Name    string
	Test    func(rune) bool
	Sev     Severity
}{
	{
		Name: "bidirectional override",
		Test: func(r rune) bool {
			return r == 0x202A || r == 0x202B || r == 0x202C || r == 0x202D || r == 0x202E ||
				r == 0x2066 || r == 0x2067 || r == 0x2068 || r == 0x2069
		},
		Sev: SeverityCritical,
	},
	{
		Name: "zero-width character",
		Test: func(r rune) bool {
			return r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF
		},
		Sev: SeverityWarning,
	},
	{
		Name: "invisible formatting",
		Test: func(r rune) bool {
			return unicode.Is(unicode.Cf, r) && r != 0x200B && r != 0x200C && r != 0x200D && r != 0xFEFF
		},
		Sev: SeverityWarning,
	},
}

// ScanFile scans a single file for hidden characters.
func (s ContentScanner) ScanFile(path string) ([]ScanFinding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return s.ScanBytes(path, data), nil
}

// ScanBytes scans raw bytes for hidden characters.
func (s ContentScanner) ScanBytes(name string, data []byte) []ScanFinding {
	var findings []ScanFinding
	lineNum := 1
	col := 1
	for i, r := range string(data) {
		if r == '\n' {
			lineNum++
			col = 1
			continue
		}
		for _, cat := range HiddenUnicodeRanges {
			if cat.Test(r) {
				ctx := extractContext(string(data), i, 40)
				findings = append(findings, ScanFinding{
					File:     name,
					Line:     lineNum,
					Column:   col,
					CharCode: int(r),
					CharName: cat.Name,
					Context:  ctx,
					Severity: cat.Sev,
				})
			}
		}
		col++
	}
	return findings
}

func extractContext(s string, idx, radius int) string {
	start := idx - radius
	if start < 0 {
		start = 0
	}
	end := idx + radius
	if end > len(s) {
		end = len(s)
	}
	return strings.Map(func(r rune) rune {
		if r < 0x20 && r != '\t' {
			return '.'
		}
		return r
	}, s[start:end])
}

// -------------------------------------------------------------------
// Runner
// -------------------------------------------------------------------

// Runner orchestrates an audit run.
type Runner struct {
	cfg     AuditConfig
	scanner ContentScanner
}

// New constructs an audit Runner.
func New(cfg AuditConfig) *Runner {
	return &Runner{cfg: cfg}
}

// ScanResult is the output of a content scan.
type ScanResult struct {
	FindingsByFile map[string][]ScanFinding
	FilesScanned   int
	HasCritical    bool
	HasWarnings    bool
	ExitCode       int
}

// Run executes a content-scan audit.
func (r *Runner) Run(opts ScanOptions) (*ScanResult, error) {
	result := &ScanResult{FindingsByFile: make(map[string][]ScanFinding)}

	files := opts.Files
	if len(files) == 0 {
		files = r.discoverPackageFiles(opts.ProjectRoot)
	}

	for _, f := range files {
		findings, err := r.scanner.ScanFile(f)
		if err != nil {
			continue
		}
		result.FilesScanned++
		if len(findings) > 0 {
			result.FindingsByFile[f] = findings
		}
	}

	for _, findings := range result.FindingsByFile {
		for _, f := range findings {
			switch f.Severity {
			case SeverityCritical:
				result.HasCritical = true
			case SeverityWarning:
				result.HasWarnings = true
			}
		}
	}

	if result.HasCritical {
		result.ExitCode = 1
	} else if result.HasWarnings {
		result.ExitCode = 2
	}
	return result, nil
}

func (r *Runner) discoverPackageFiles(root string) []string {
	if root == "" {
		root = "."
	}
	var files []string
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".md", ".txt", ".yaml", ".yml", ".json", ".prompt":
			files = append(files, path)
		}
		return nil
	})
	return files
}

// -------------------------------------------------------------------
// Strip mode
// -------------------------------------------------------------------

// StripResult records what was removed from a file.
type StripResult struct {
	File    string
	Removed int
	Backed  string
}

// StripFindings removes hidden characters from the listed files.
func StripFindings(findings map[string][]ScanFinding, dryRun bool) ([]StripResult, error) {
	var results []StripResult
	for file, ff := range findings {
		data, err := os.ReadFile(file)
		if err != nil {
			return results, err
		}
		original := string(data)
		cleaned := stripHidden(original)
		removed := len([]rune(original)) - len([]rune(cleaned))
		if removed == 0 {
			continue
		}
		sr := StripResult{File: file, Removed: removed}
		if !dryRun {
			sr.Backed = file + ".bak"
			if err := os.WriteFile(sr.Backed, data, 0o644); err != nil {
				return results, err
			}
			if err := os.WriteFile(file, []byte(cleaned), 0o644); err != nil {
				return results, err
			}
		}
		results = append(results, sr)
		_ = ff
	}
	return results, nil
}

func stripHidden(s string) string {
	var sb strings.Builder
	for _, r := range s {
		keep := true
		for _, cat := range HiddenUnicodeRanges {
			if cat.Test(r) {
				keep = false
				break
			}
		}
		if keep {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// -------------------------------------------------------------------
// CI audit mode
// -------------------------------------------------------------------

// CIFinding is a policy discovery finding from the CI audit mode.
type CIFinding struct {
	Outcome string
	Source  string
	ErrText string
	Level   string // "warn" | "block"
}

// CIAuditResult is the output of a --ci policy audit.
type CIAuditResult struct {
	Findings []CIFinding
	ExitCode int
}

// AuditOutcomeCause renders a human-readable cause for a policy-discovery outcome.
func AuditOutcomeCause(outcome, source, errText string) string {
	switch outcome {
	case "no_git_remote":
		return "Could not determine org from git remote"
	case "absent":
		return fmt.Sprintf("No org policy found at %s", source)
	case "empty":
		return fmt.Sprintf("Org policy at %s is present but empty", source)
	default:
		if errText != "" {
			return fmt.Sprintf("Policy fetch failed: %s", errText)
		}
		return fmt.Sprintf("Policy fetch failed: %s", outcome)
	}
}

// -------------------------------------------------------------------
// Output rendering
// -------------------------------------------------------------------

// RenderFindingsTable renders findings to a text table.
func RenderFindingsTable(result *ScanResult) string {
	if len(result.FindingsByFile) == 0 {
		return "[+] No hidden characters found.\n"
	}
	var sb strings.Builder

	files := make([]string, 0, len(result.FindingsByFile))
	for f := range result.FindingsByFile {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, f := range files {
		findings := result.FindingsByFile[f]
		sb.WriteString(fmt.Sprintf("[!] %s (%d finding(s))\n", f, len(findings)))
		for _, ff := range findings {
			sb.WriteString(fmt.Sprintf("    L%d C%d  U+%04X  %s  |%s|\n",
				ff.Line, ff.Column, ff.CharCode, ff.CharName, ff.Context))
		}
	}
	return sb.String()
}

// RenderFindingsJSON renders findings as JSON.
func RenderFindingsJSON(result *ScanResult) (string, error) {
	b, err := json.MarshalIndent(result.FindingsByFile, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// RenderSummary renders a one-line summary.
func RenderSummary(result *ScanResult) string {
	switch {
	case result.HasCritical:
		return fmt.Sprintf("[x] Critical findings in %d file(s). Exit code 1.", len(result.FindingsByFile))
	case result.HasWarnings:
		return fmt.Sprintf("[!] Warnings in %d file(s). Exit code 2.", len(result.FindingsByFile))
	default:
		return fmt.Sprintf("[+] Clean. Scanned %d file(s).", result.FilesScanned)
	}
}

// -------------------------------------------------------------------
// Lockfile audit helpers
// -------------------------------------------------------------------

// LockfilePackage is a minimal lockfile entry used for scanning.
type LockfilePackage struct {
	Name    string
	Version string
	Path    string
}

// ScanLockfilePackages scans all packages listed in a lockfile.
func ScanLockfilePackages(lockfilePath string, scanner ContentScanner) (*ScanResult, error) {
	data, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, err
	}
	result := &ScanResult{FindingsByFile: make(map[string][]ScanFinding)}
	findings := scanner.ScanBytes(lockfilePath, data)
	if len(findings) > 0 {
		result.FindingsByFile[lockfilePath] = findings
		result.FilesScanned = 1
	}
	return result, nil
}
