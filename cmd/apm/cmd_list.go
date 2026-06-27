// cmd_list.go implements `apm list` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/list.py.
package main

import (
	"fmt"
	"os"
)

// runList implements `apm list [OPTIONS]`.
func runList(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printCmdHelp("list")
			return 0
		}
		if startsWith(a, "-") {
			fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
			fmt.Fprintln(os.Stderr, `Try 'apm list --help' for help.`)
			return 2
		}
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

	if len(proj.Scripts) == 0 {
		fmt.Println("No scripts defined in apm.yml.")
		fmt.Println("Add a 'scripts:' section to define runnable scripts.")
		return 0
	}

	fmt.Printf("Scripts in project '%s':\n", proj.Name)
	for name, desc := range proj.Scripts {
		if desc != "" {
			fmt.Printf("  %-20s %s\n", name, desc)
		} else {
			fmt.Printf("  %s\n", name)
		}
	}
	return 0
}
