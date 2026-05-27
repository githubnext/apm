// cmd_init.go implements the `apm init` command for the Go rewrite.
// Mirrors src/apm_cli/commands/init.py (non-interactive --yes path).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runInit implements `apm init [OPTIONS] [PROJECT_NAME]`.
// Supports --yes/-y (skip prompts), --verbose/-v, --help/-h.
// Returns an OS exit code.
func runInit(args []string) int {
	var (
		flagYes     bool
		flagVerbose bool
		flagHelp    bool
		projectName string
	)

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--yes", "-y":
			flagYes = true
		case "--verbose", "-v":
			flagVerbose = true
		case "--help", "-h", "-help":
			flagHelp = true
		case "--plugin", "--marketplace":
			// Deprecated flags: warn and continue.
			flag := args[i]
			fmt.Fprintf(os.Stderr, "[!] '%s' is deprecated. See 'apm --help' for alternatives.\n", "apm init "+flag)
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[i])
				fmt.Fprintln(os.Stderr, `Try 'apm init --help' for help.`)
				return 2
			}
			if projectName == "" {
				projectName = args[i]
			}
		}
		i++
	}

	if flagHelp {
		printCmdHelp("init")
		return 0
	}

	// Non-interactive mode required when running without a TTY (CI, tests).
	// With --yes the Python CLI skips all prompts. We always behave that way.
	if !flagYes {
		// If stdout is not a terminal we auto-apply --yes behaviour.
		fi, _ := os.Stdout.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			flagYes = true
		}
	}

	return execInit(projectName, flagYes, flagVerbose)
}

// execInit performs the actual project initialization.
func execInit(projectName string, _ bool, verbose bool) int {
	// Handle explicit current directory.
	if projectName == "." {
		projectName = ""
	}

	// Validate project name.
	if projectName != "" {
		if strings.ContainsAny(projectName, "/\\") || projectName == ".." {
			fmt.Fprintf(os.Stderr, "Error: Invalid project name '%s': must not contain path separators or be '..'.\n", projectName)
			return 1
		}
	}

	// Determine project directory.
	var projectDir string
	var finalName string
	if projectName != "" {
		projectDir = projectName
		if err := os.MkdirAll(projectDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not create directory '%s': %v\n", projectDir, err)
			return 1
		}
		finalName = projectName
		if verbose {
			fmt.Printf("[*] Created project directory: %s\n", projectName)
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not determine working directory: %v\n", err)
			return 1
		}
		projectDir = cwd
		finalName = filepath.Base(cwd)
	}

	apmYMLPath := filepath.Join(projectDir, "apm.yml")

	// Check if apm.yml already exists.
	if _, err := os.Stat(apmYMLPath); err == nil {
		fmt.Fprintf(os.Stderr, "[!] apm.yml already exists in '%s'. Skipping.\n", projectDir)
		return 0
	}

	fmt.Printf("[>] Initializing APM project: %s\n", finalName)

	content := buildApmYML(finalName)
	if err := os.WriteFile(apmYMLPath, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not write apm.yml: %v\n", err)
		return 1
	}

	fmt.Printf("[+] APM project initialized successfully!\n")
	fmt.Printf("    Created: apm.yml\n")
	fmt.Printf("\n")
	fmt.Printf("  Next Steps\n")
	fmt.Printf("  * Install a package:            apm install <owner>/<repo>\n")
	fmt.Printf("  * Run a script:                 apm run <script>\n")
	fmt.Printf("  * Build a plugin? Scaffold one: apm plugin init\n")
	fmt.Printf("  Docs: https://microsoft.github.io/apm  |  Star: https://github.com/microsoft/apm\n")

	return 0
}

// buildApmYML returns the contents of a minimal apm.yml for the given project name.
func buildApmYML(name string) string {
	return fmt.Sprintf(`name: %s
version: 1.0.0
description: APM project for %s
author: Developer
# Which agent platforms to deploy to (uncomment to pin):
# targets:
#   - copilot
#   - claude

dependencies:
  apm: []
  mcp: []
includes: auto
scripts: {}
`, name, name)
}
