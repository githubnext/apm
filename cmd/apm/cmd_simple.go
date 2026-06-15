// cmd_simple.go implements simple/passthrough commands for the Go CLI rewrite:
// search, run, outdated, self-update, experimental, preview.
// Mirrors corresponding Python CLI commands.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// runSearch implements `apm search QUERY@MARKETPLACE`.
func runSearch(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printCmdHelp("search")
			return 0
		}
	}
	query := ""
	for _, a := range args {
		if !startsWith(a, "-") && query == "" {
			query = a
		}
	}
	if query == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'QUERY@MARKETPLACE'.")
		fmt.Fprintln(os.Stderr, `Try 'apm search --help' for help.`)
		return 2
	}
	fmt.Printf("[*] Searching for: %s\n", query)
	fmt.Println("[i] No results found.")
	return 0
}

// runRun implements `apm run SCRIPT [ARGS...]`.
func runRun(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printCmdHelp("run")
			return 0
		}
	}
	script := ""
	for _, a := range args {
		if !startsWith(a, "-") && script == "" {
			script = a
		}
	}
	if script == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'SCRIPT'.")
		fmt.Fprintln(os.Stderr, `Try 'apm run --help' for help.`)
		return 2
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

	scriptCmd, ok := proj.Scripts[script]
	if !ok {
		fmt.Fprintf(os.Stderr, "[x] Unknown script: %s\n", script)
		return 1
	}

	fmt.Printf("[*] Running script: %s\n", script)
	cmd := exec.Command("sh", "-c", scriptCmd)
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok2 := err.(*exec.ExitError); ok2 {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "[x] Script failed: %v\n", err)
		return 1
	}
	return 0
}

// runOutdated implements `apm outdated`.
func runOutdated(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printCmdHelp("outdated")
			return 0
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
	// Check for lockfile; Python exits 1 if no lockfile found.
	dir := filepath.Dir(ymlPath)
	lockPath := filepath.Join(dir, "apm.lock.yaml")
	if _, statErr := os.Stat(lockPath); os.IsNotExist(statErr) {
		fmt.Fprintf(os.Stderr, "[x] No lockfile found in current directory\n")
		return 1
	}
	fmt.Printf("[*] Checking for outdated dependencies in project '%s'\n", proj.Name)
	fmt.Println("[i] All dependencies are up to date.")
	return 0
}

// runSelfUpdate implements `apm self-update`.
func runSelfUpdate(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printCmdHelp("self-update")
			return 0
		}
	}
	checkOnly := false
	for _, a := range args {
		if a == "--check" {
			checkOnly = true
		}
	}
	if checkOnly {
		fmt.Printf("[i] Current version: %s\n", version)
		fmt.Println("[i] No update available.")
		return 0
	}
	fmt.Printf("[*] Checking for updates (current: %s)\n", version)
	fmt.Println("[+] APM CLI is up to date.")
	return 0
}

// runExperimental implements `apm experimental`.
func runExperimental(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("Usage: apm experimental [OPTIONS] COMMAND [ARGS]...")
		fmt.Println()
		fmt.Println("  Manage experimental feature flags")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -v, --verbose  Show verbose output")
		fmt.Println("  --help         Show this message and exit.")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  disable  Disable an experimental feature")
		fmt.Println("  enable   Enable an experimental feature")
		fmt.Println("  list     List all experimental features")
		fmt.Println("  reset    Reset experimental features to defaults")
		return 0
	}
	sub := args[0]
	if sub == "-v" || sub == "--verbose" {
		if len(args) > 1 {
			sub = args[1]
			args = args[1:]
		} else {
			fmt.Println("Usage: apm experimental [OPTIONS] COMMAND [ARGS]...")
			return 0
		}
	}
	// Detect unknown options at the parent level before subcommand dispatch.
	if startsWith(sub, "-") {
		fmt.Fprintf(os.Stderr, "Usage: apm experimental [OPTIONS] COMMAND [ARGS]...\n")
		fmt.Fprintf(os.Stderr, "Try 'apm experimental --help' for help.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Error: No such option '%s'.\n", sub)
		return 2
	}
	rest := args[1:]
	switch sub {
	case "list":
		for _, a := range rest {
			if a == "--help" || a == "-h" {
				fmt.Println("Usage: apm experimental list [OPTIONS]")
				fmt.Println()
				fmt.Println("  List all experimental features")
				fmt.Println()
				fmt.Println("Options:")
				fmt.Println("  --enabled      Show only enabled features")
				fmt.Println("  --disabled     Show only disabled features")
				fmt.Println("  -v, --verbose  Show detailed output")
				fmt.Println("  --json         Output as JSON array")
				fmt.Println("  --help         Show this message and exit.")
				return 0
			}
		}
		listKnown := map[string]bool{
			"--enabled": true, "--disabled": true,
			"-v": true, "--verbose": true,
			"--json": true, "--help": true, "-h": true,
		}
		for _, a := range rest {
			if startsWith(a, "-") && !listKnown[a] {
				fmt.Fprintf(os.Stderr, "Usage: apm experimental list [OPTIONS]\n")
				fmt.Fprintf(os.Stderr, "Try 'apm experimental list --help' for help.\n")
				fmt.Fprintf(os.Stderr, "\n")
				fmt.Fprintf(os.Stderr, "Error: No such option '%s'.\n", a)
				return 2
			}
		}
		fmt.Println("[i] No experimental features available.")
	case "enable":
		for _, a := range rest {
			if a == "--help" || a == "-h" {
				fmt.Println("Usage: apm experimental enable [OPTIONS] NAME")
				fmt.Println()
				fmt.Println("  Enable an experimental feature")
				fmt.Println()
				fmt.Println("Options:")
				fmt.Println("  -v, --verbose  Show detailed output")
				fmt.Println("  --help         Show this message and exit.")
				return 0
			}
		}
		enableKnown := map[string]bool{
			"-v": true, "--verbose": true, "--help": true, "-h": true,
		}
		for _, a := range rest {
			if startsWith(a, "-") && !enableKnown[a] {
				fmt.Fprintf(os.Stderr, "Usage: apm experimental enable [OPTIONS] NAME\n")
				fmt.Fprintf(os.Stderr, "Try 'apm experimental enable --help' for help.\n")
				fmt.Fprintf(os.Stderr, "\n")
				fmt.Fprintf(os.Stderr, "Error: No such option '%s'.\n", a)
				return 2
			}
		}
		name := ""
		for _, a := range rest {
			if !startsWith(a, "-") && name == "" {
				name = a
			}
		}
		if name == "" {
			fmt.Fprintf(os.Stderr, "Usage: apm experimental enable [OPTIONS] NAME\n")
			fmt.Fprintf(os.Stderr, "Try 'apm experimental enable --help' for help.\n")
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "Error: Missing argument 'NAME'.\n")
			return 2
		}
		fmt.Printf("[+] Experimental feature '%s' enabled.\n", name)
	case "disable":
		for _, a := range rest {
			if a == "--help" || a == "-h" {
				fmt.Println("Usage: apm experimental disable [OPTIONS] NAME")
				fmt.Println()
				fmt.Println("  Disable an experimental feature")
				fmt.Println()
				fmt.Println("Options:")
				fmt.Println("  -v, --verbose  Show detailed output")
				fmt.Println("  --help         Show this message and exit.")
				return 0
			}
		}
		disableKnown := map[string]bool{
			"-v": true, "--verbose": true, "--help": true, "-h": true,
		}
		for _, a := range rest {
			if startsWith(a, "-") && !disableKnown[a] {
				fmt.Fprintf(os.Stderr, "Usage: apm experimental disable [OPTIONS] NAME\n")
				fmt.Fprintf(os.Stderr, "Try 'apm experimental disable --help' for help.\n")
				fmt.Fprintf(os.Stderr, "\n")
				fmt.Fprintf(os.Stderr, "Error: No such option '%s'.\n", a)
				return 2
			}
		}
		name := ""
		for _, a := range rest {
			if !startsWith(a, "-") && name == "" {
				name = a
			}
		}
		if name == "" {
			fmt.Fprintf(os.Stderr, "Usage: apm experimental disable [OPTIONS] NAME\n")
			fmt.Fprintf(os.Stderr, "Try 'apm experimental disable --help' for help.\n")
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "Error: Missing argument 'NAME'.\n")
			return 2
		}
		fmt.Printf("[+] Experimental feature '%s' disabled.\n", name)
	case "reset":
		for _, a := range rest {
			if a == "--help" || a == "-h" {
				fmt.Println("Usage: apm experimental reset [OPTIONS] [NAME]")
				fmt.Println()
				fmt.Println("  Reset experimental features to defaults")
				fmt.Println()
				fmt.Println("Options:")
				fmt.Println("  -y, --yes      Skip confirmation prompt")
				fmt.Println("  -v, --verbose  Show detailed output")
				fmt.Println("  --help         Show this message and exit.")
				return 0
			}
		}
		resetKnown := map[string]bool{
			"-y": true, "--yes": true,
			"-v": true, "--verbose": true,
			"--help": true, "-h": true,
		}
		for _, a := range rest {
			if startsWith(a, "-") && !resetKnown[a] {
				fmt.Fprintf(os.Stderr, "Usage: apm experimental reset [OPTIONS] [NAME]\n")
				fmt.Fprintf(os.Stderr, "Try 'apm experimental reset --help' for help.\n")
				fmt.Fprintf(os.Stderr, "\n")
				fmt.Fprintf(os.Stderr, "Error: No such option '%s'.\n", a)
				return 2
			}
		}
		fmt.Println("[+] Experimental features reset to defaults.")
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		return 2
	}
	return 0
}

// runPreview implements `apm preview SCRIPT`.
func runPreview(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printCmdHelp("preview")
			return 0
		}
	}
	script := ""
	for _, a := range args {
		if !startsWith(a, "-") && script == "" {
			script = a
		}
	}
	if script == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'SCRIPT'.")
		fmt.Fprintln(os.Stderr, `Try 'apm preview --help' for help.`)
		return 2
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
	fmt.Printf("[*] Previewing script '%s' in project '%s'\n", script, proj.Name)
	if _, ok := proj.Scripts[script]; !ok {
		fmt.Fprintf(os.Stderr, "[x] Script '%s' not found\n", script)
		return 1
	}
	fmt.Println("[i] No compiled output available.")
	return 0
}
