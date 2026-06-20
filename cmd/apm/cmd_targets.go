// cmd_targets.go implements `apm targets` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/targets.py.
package main

import (
	"fmt"
	"os"
)

// knownTargets is the canonical list of supported target platforms.
var knownTargets = []string{
	"copilot", "claude", "cursor", "opencode", "codex", "gemini", "windsurf",
}

// autoDetectTargets returns a default target list when none is specified.
func autoDetectTargets() []string {
	return []string{"copilot"}
}

// runTargets implements `apm targets [OPTIONS]`.
func runTargets(args []string) int {
	var flagJSON, flagAll, flagHelp bool
	for _, a := range args {
		switch a {
		case "--json":
			flagJSON = true
		case "--all":
			flagAll = true
		case "--help", "-h":
			flagHelp = true
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm targets --help' for help.`)
				return 2
			}
		}
	}
	if flagHelp {
		printCmdHelp("targets")
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

	targets := proj.Targets
	if len(targets) == 0 {
		targets = autoDetectTargets()
	}

	if flagJSON {
		fmt.Print("[")
		for i, t := range targets {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%q", t)
		}
		if flagAll {
			if len(targets) > 0 {
				fmt.Print(", ")
			}
			fmt.Print(`"agent-skills"`)
		}
		fmt.Println("]")
		return 0
	}

	fmt.Printf("Targets for project '%s':\n", proj.Name)
	for _, t := range targets {
		fmt.Printf("  [+] %s\n", t)
	}
	if flagAll {
		fmt.Printf("  [+] agent-skills  (meta)\n")
	}
	return 0
}
