// Package auditreport provides serialization helpers for apm audit results.
// Supports JSON, SARIF 2.1.0, and Markdown output formats.
// Migrated from src/apm_cli/security/audit_report.py
package auditreport

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ScanFinding represents a single security finding from a content scan.
type ScanFinding struct {
	// Severity is "critical", "warning", or "info".
	Severity string
	// File is the path to the file containing the finding.
	File string
	// Line is the 1-based line number.
	Line int
	// Column is the 1-based column number.
	Column int
	// Codepoint is the Unicode codepoint string (e.g. "U+200B").
	Codepoint string
	// Category classifies the finding type (e.g. "zero-width").
	Category string
	// Description is a human-readable explanation.
	Description string
}

const (
	sarifVersion = "2.1.0"
	sarifSchema  = "https://docs.oasis-open.org/sarif/sarif/v2.1.0/cos02/schemas/sarif-schema-2.1.0.json"
	toolName     = "apm-audit"
	toolInfoURI  = "https://apm.github.io/apm/enterprise/security/"
)

// severityMap maps APM severity strings to SARIF level strings.
var severityMap = map[string]string{
	"critical": "error",
	"warning":  "warning",
	"info":     "note",
}

// RelativePathForReport normalizes a file path to a relative forward-slash path.
func RelativePathForReport(filePath string) string {
	p := filepath.Clean(filePath)
	if filepath.IsAbs(p) {
		cwd, err := os.Getwd()
		if err == nil {
			rel, err2 := filepath.Rel(cwd, p)
			if err2 == nil {
				return filepath.ToSlash(rel)
			}
		}
		return filepath.Base(p)
	}
	return strings.ReplaceAll(filePath, "\\", "/")
}

// ruleID builds a SARIF rule ID from a finding category.
func ruleID(category string) string {
	return "apm/hidden-unicode/" + category
}

// allFindings flattens a map of findings by file into a single slice.
func allFindings(findingsByFile map[string][]ScanFinding) []ScanFinding {
	var out []ScanFinding
	for _, ff := range findingsByFile {
		out = append(out, ff...)
	}
	return out
}

// FindingsToJSON converts scan findings to APM's JSON report format.
func FindingsToJSON(findingsByFile map[string][]ScanFinding, filesScanned int, exitCode int) map[string]interface{} {
	all := allFindings(findingsByFile)

	critical, warning, info := 0, 0, 0
	for _, f := range all {
		switch f.Severity {
		case "critical":
			critical++
		case "warning":
			warning++
		case "info":
			info++
		}
	}

	items := make([]map[string]interface{}, 0, len(all))
	for _, f := range all {
		items = append(items, map[string]interface{}{
			"severity":    f.Severity,
			"file":        RelativePathForReport(f.File),
			"line":        f.Line,
			"column":      f.Column,
			"codepoint":   f.Codepoint,
			"category":    f.Category,
			"description": f.Description,
		})
	}

	return map[string]interface{}{
		"version":   "1",
		"exit_code": exitCode,
		"summary": map[string]interface{}{
			"files_scanned":  filesScanned,
			"files_affected": len(findingsByFile),
			"critical":       critical,
			"warning":        warning,
			"info":           info,
		},
		"findings": items,
	}
}

// FindingsToSARIF converts scan findings to SARIF 2.1.0 format.
func FindingsToSARIF(findingsByFile map[string][]ScanFinding, filesScanned int) map[string]interface{} {
	all := allFindings(findingsByFile)

	seenRules := map[string]map[string]interface{}{}
	for _, f := range all {
		rid := ruleID(f.Category)
		if _, exists := seenRules[rid]; !exists {
			seenRules[rid] = map[string]interface{}{
				"id": rid,
				"shortDescription": map[string]interface{}{
					"text": strings.Title(strings.ReplaceAll(f.Category, "-", " ")),
				},
				"defaultConfiguration": map[string]interface{}{
					"level": func() string {
						if v, ok := severityMap[f.Severity]; ok {
							return v
						}
						return "note"
					}(),
				},
				"helpUri": toolInfoURI,
			}
		}
	}

	rulesList := make([]interface{}, 0, len(seenRules))
	for _, r := range seenRules {
		rulesList = append(rulesList, r)
	}

	results := make([]interface{}, 0, len(all))
	for _, f := range all {
		level := "note"
		if v, ok := severityMap[f.Severity]; ok {
			level = v
		}
		results = append(results, map[string]interface{}{
			"ruleId": ruleID(f.Category),
			"level":  level,
			"message": map[string]interface{}{
				"text": fmt.Sprintf("%s (%s)", f.Description, f.Codepoint),
			},
			"locations": []interface{}{
				map[string]interface{}{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]interface{}{
							"uri": RelativePathForReport(f.File),
						},
						"region": map[string]interface{}{
							"startLine":   f.Line,
							"startColumn": f.Column,
						},
					},
				},
			},
			"properties": map[string]interface{}{
				"codepoint": f.Codepoint,
				"category":  f.Category,
			},
		})
	}

	return map[string]interface{}{
		"$schema": sarifSchema,
		"version": sarifVersion,
		"runs": []interface{}{
			map[string]interface{}{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":           toolName,
						"informationUri": toolInfoURI,
						"rules":          rulesList,
					},
				},
				"results": results,
				"invocations": []interface{}{
					map[string]interface{}{
						"executionSuccessful": true,
						"properties": map[string]interface{}{
							"filesScanned": filesScanned,
						},
					},
				},
			},
		},
	}
}

// WriteReport writes a report dict as JSON to the given path.
func WriteReport(report map[string]interface{}, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, append(data, '\n'), 0o644)
}

// SerializeReport serializes a report dict to a JSON string.
func SerializeReport(report map[string]interface{}) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FindingsToMarkdown converts scan findings to GitHub-Flavored Markdown.
func FindingsToMarkdown(findingsByFile map[string][]ScanFinding, filesScanned int) string {
	all := allFindings(findingsByFile)

	if len(all) == 0 {
		return fmt.Sprintf("## APM Audit Report\n\n**Clean** -- no security findings across %d files.\n", filesScanned)
	}

	critical, warning, info := 0, 0, 0
	for _, f := range all {
		switch f.Severity {
		case "critical":
			critical++
		case "warning":
			warning++
		case "info":
			info++
		}
	}
	affected := len(findingsByFile)
	total := len(all)

	parts := []string{}
	if critical > 0 {
		parts = append(parts, fmt.Sprintf("%d critical", critical))
	}
	if warning > 0 {
		s := "s"
		if warning == 1 {
			s = ""
		}
		parts = append(parts, fmt.Sprintf("%d warning%s", warning, s))
	}
	if info > 0 {
		parts = append(parts, fmt.Sprintf("%d info", info))
	}

	countLabel := fmt.Sprintf("**%d finding", total)
	if total != 1 {
		countLabel += "s"
	}
	countLabel += "**"

	affectedStr := "files"
	if affected == 1 {
		affectedStr = "file"
	}

	summary := fmt.Sprintf("%s across %d %s (%s) | %d files scanned",
		countLabel, affected, affectedStr, strings.Join(parts, ", "), filesScanned)

	severityOrder := map[string]int{"critical": 0, "warning": 1, "info": 2}
	sort.SliceStable(all, func(i, j int) bool {
		si := severityOrder[all[i].Severity]
		sj := severityOrder[all[j].Severity]
		if si != sj {
			return si < sj
		}
		if all[i].File != all[j].File {
			return all[i].File < all[j].File
		}
		return all[i].Line < all[j].Line
	})

	var sb strings.Builder
	sb.WriteString("## APM Audit Report\n\n")
	sb.WriteString(summary + "\n\n")
	sb.WriteString("| Severity | File | Location | Codepoint | Description |\n")
	sb.WriteString("|----------|------|----------|-----------|-------------|\n")
	for _, f := range all {
		sev := strings.ToUpper(f.Severity)
		desc := strings.ReplaceAll(f.Description, "|", "\\|")
		sb.WriteString(fmt.Sprintf("| %s | `%s` | %d:%d | `%s` | %s |\n",
			sev, RelativePathForReport(f.File), f.Line, f.Column, f.Codepoint, desc))
	}
	sb.WriteString("\nRun `apm audit --strip` to remove flagged characters.\n")

	return sb.String()
}

// DetectFormatFromExtension auto-detects output format from file extension.
func DetectFormatFromExtension(path string) string {
	name := strings.ToLower(filepath.Base(path))
	if strings.HasSuffix(name, ".sarif.json") || strings.HasSuffix(name, ".sarif") {
		return "sarif"
	}
	if strings.HasSuffix(name, ".json") {
		return "json"
	}
	if strings.HasSuffix(name, ".md") {
		return "markdown"
	}
	return "text"
}
