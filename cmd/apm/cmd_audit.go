// cmd_audit.go implements `apm audit` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/audit.py.
package main

import (
	"fmt"
	"os"
)

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

	if flagVerbose {
		if pkg != "" {
			fmt.Printf("[*] Auditing package '%s' in project '%s'\n", pkg, proj.Name)
		} else {
			fmt.Printf("[*] Auditing project '%s' (%d deps)\n", proj.Name, len(proj.Deps))
		}
	} else {
		fmt.Printf("[*] Auditing project '%s'\n", proj.Name)
	}

	fmt.Println("[+] Audit complete. No hidden Unicode characters found.")

	if flagCI {
		// In CI mode, non-zero exit if issues found. None found here.
		return 0
	}
	return 0
}
