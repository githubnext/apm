// Package contentscanner detects hidden Unicode characters in text files.
// It mirrors src/apm_cli/security/content_scanner.py.
//
// Scans for invisible characters (Unicode tags, bidi overrides, variation
// selectors, zero-width characters) that could embed hidden instructions
// in prompt, instruction, and rules files.
package contentscanner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// ScanFinding describes a single suspicious character found during scanning.
type ScanFinding struct {
	File        string
	Line        int
	Column      int
	Char        rune
	Codepoint   string // e.g. "U+200B"
	Severity    string // "critical", "warning", "info"
	Category    string // e.g. "bidi-override", "zero-width"
	Description string
}

type suspiciousRange struct {
	start       rune
	end         rune
	severity    string
	category    string
	description string
}

var suspiciousRanges = []suspiciousRange{
	// Unicode tag characters (invisible ASCII mapping)
	{0xE0001, 0xE007F, "critical", "tag-character", "Unicode tag character (invisible ASCII mapping)"},
	// Bidirectional override characters
	{0x202A, 0x202A, "critical", "bidi-override", "Left-to-right embedding (LRE)"},
	{0x202B, 0x202B, "critical", "bidi-override", "Right-to-left embedding (RLE)"},
	{0x202C, 0x202C, "critical", "bidi-override", "Pop directional formatting (PDF)"},
	{0x202D, 0x202D, "critical", "bidi-override", "Left-to-right override (LRO)"},
	{0x202E, 0x202E, "critical", "bidi-override", "Right-to-left override (RLO)"},
	{0x2066, 0x2066, "critical", "bidi-override", "Left-to-right isolate (LRI)"},
	{0x2067, 0x2067, "critical", "bidi-override", "Right-to-left isolate (RLI)"},
	{0x2068, 0x2068, "critical", "bidi-override", "First strong isolate (FSI)"},
	{0x2069, 0x2069, "critical", "bidi-override", "Pop directional isolate (PDI)"},
	// Variation selectors (Glassworm attack vector)
	{0xE0100, 0xE01EF, "critical", "variation-selector", "Variation selector supplement (hidden payload encoding)"},
	{0xFE00, 0xFE0F, "warning", "variation-selector", "Variation selector (possible payload encoding)"},
	// Zero-width characters
	{0x200B, 0x200B, "warning", "zero-width", "Zero-width space"},
	{0x200C, 0x200C, "warning", "zero-width", "Zero-width non-joiner"},
	{0x200D, 0x200D, "warning", "zero-width", "Zero-width joiner"},
	{0xFEFF, 0xFEFF, "warning", "zero-width", "Zero-width no-break space (BOM)"},
	{0x2060, 0x2060, "warning", "zero-width", "Word joiner"},
	// Invisible separators
	{0x2028, 0x2028, "info", "invisible-separator", "Line separator"},
	{0x2029, 0x2029, "info", "invisible-separator", "Paragraph separator"},
}

func classify(r rune) (severity, category, description string, ok bool) {
	for _, sr := range suspiciousRanges {
		if r >= sr.start && r <= sr.end {
			return sr.severity, sr.category, sr.description, true
		}
	}
	return "", "", "", false
}

// ScanText scans the provided text content for suspicious characters.
// filePath is used only for populating ScanFinding.File.
func ScanText(filePath, content string) []ScanFinding {
	var findings []ScanFinding
	lineNum := 0
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		col := 0
		for i := 0; i < len(line); {
			r, size := utf8.DecodeRuneInString(line[i:])
			col++
			if sev, cat, desc, ok := classify(r); ok {
				findings = append(findings, ScanFinding{
					File:        filePath,
					Line:        lineNum,
					Column:      col,
					Char:        r,
					Codepoint:   fmt.Sprintf("U+%04X", r),
					Severity:    sev,
					Category:    cat,
					Description: desc,
				})
			}
			i += size
		}
	}
	return findings
}

// ScanFile reads and scans a file for suspicious characters.
func ScanFile(filePath string) ([]ScanFinding, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ScanText(filePath, string(data)), nil
}

// ContentScanner scans multiple files and aggregates findings.
type ContentScanner struct {
	Extensions []string // file extensions to scan (e.g. ".md", ".txt")
}

// NewDefaultScanner returns a ContentScanner for typical prompt/instruction files.
func NewDefaultScanner() *ContentScanner {
	return &ContentScanner{
		Extensions: []string{".md", ".txt", ".prompt", ".instructions"},
	}
}

// ScanFiles scans the given list of file paths.
func (cs *ContentScanner) ScanFiles(paths []string) map[string][]ScanFinding {
	results := make(map[string][]ScanFinding)
	for _, p := range paths {
		if !cs.shouldScan(p) {
			continue
		}
		findings, err := ScanFile(p)
		if err != nil {
			continue
		}
		if len(findings) > 0 {
			results[p] = findings
		}
	}
	return results
}

func (cs *ContentScanner) shouldScan(p string) bool {
	if len(cs.Extensions) == 0 {
		return true
	}
	lower := strings.ToLower(p)
	for _, ext := range cs.Extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
