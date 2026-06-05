// cmd_audit.go implements `apm audit` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/audit.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// hiddenUnicodeChars lists Unicode codepoints that can be used to obfuscate code.
var hiddenUnicodeChars = map[rune]string{
	'\u202e': "RIGHT-TO-LEFT OVERRIDE",
	'\u202d': "LEFT-TO-RIGHT OVERRIDE",
	'\u202c': "POP DIRECTIONAL FORMATTING",
	'\u202b': "RIGHT-TO-LEFT EMBEDDING",
	'\u202a': "LEFT-TO-RIGHT EMBEDDING",
	'\u200b': "ZERO WIDTH SPACE",
	'\u200c': "ZERO WIDTH NON-JOINER",
	'\u200d': "ZERO WIDTH JOINER",
	'\ufeff': "BYTE ORDER MARK",
	'\u2060': "WORD JOINER",
	'\u00ad': "SOFT HYPHEN",
	'\u2066': "LEFT-TO-RIGHT ISOLATE",
	'\u2067': "RIGHT-TO-LEFT ISOLATE",
	'\u2068': "FIRST STRONG ISOLATE",
	'\u2069': "POP DIRECTIONAL ISOLATE",
}

// auditFinding records a hidden unicode detection.
type auditFinding struct {
	path string
	char rune
	name string
}

// scanForHiddenUnicode walks dir and returns findings.
func scanForHiddenUnicode(dir string) []auditFinding {
	var findings []auditFinding
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		for _, r := range string(data) {
			if name, ok := hiddenUnicodeChars[r]; ok {
				findings = append(findings, auditFinding{path: path, char: r, name: name})
				break
			}
		}
		return nil
	})
	return findings
}

// runAudit implements `apm audit [OPTIONS] [PACKAGE]`.
func runAudit(args []string) int {
	var (
		flagHelp    bool
		flagCI      bool
		flagVerbose bool
		pkg         string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			flagHelp = true
		case "--ci":
			flagCI = true
		case "-v", "--verbose", "--verbose-output":
			flagVerbose = true
		case "--json", "--summary", "--all":
			// consumed flag
		case "--target", "--runtime", "--exclude", "--only":
			if i+1 < len(args) {
				i++
			}
		default:
			if !startsWith(args[i], "-") && pkg == "" {
				pkg = args[i]
			}
		}
	}

	if flagHelp {
		printCmdHelp("audit")
		return 0
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] No apm.yml found. Run 'apm init' to create one.\n")
		return 1
	}
	proj, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
		return 1
	}

	scanDir := filepath.Join(cwd, "apm_modules")
	if pkg != "" {
		scanDir = filepath.Join(cwd, "apm_modules", pkg)
	}

	if flagVerbose {
		if pkg != "" {
			fmt.Printf("[*] Auditing package '%s' in project '%s'\n", pkg, proj.Name)
		} else {
			fmt.Printf("[*] Auditing project '%s' (%d deps)\n", proj.Name, len(proj.Deps))
		}
	} else {
		fmt.Printf("[*] Auditing project '%s'\n", proj.Name)
	}

	findings := scanForHiddenUnicode(scanDir)
	if len(findings) > 0 {
		for _, f := range findings {
			rel, _ := filepath.Rel(cwd, f.path)
			fmt.Fprintf(os.Stderr, "[x] Hidden Unicode detected: %s (U+%04X %s)\n", rel, f.char, f.name)
		}
		if flagCI {
			return 1
		}
		return 1
	}

	fmt.Println("[+] Audit complete. No hidden Unicode characters found.")
	return 0
}
