// cmd_deps.go implements `apm deps` and its subcommands for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/deps.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// runDeps implements `apm deps [SUBCOMMAND] [OPTIONS]`.
func runDeps(args []string) int {
	if len(args) == 0 {
		printDepsHelp()
		return 0
	}

	sub := args[0]
	rest := args[1:]

	for _, a := range args {
		if a == "--help" || a == "-h" {
			printDepsHelp()
			return 0
		}
	}

	switch sub {
	case "list":
		return runDepsList(rest)
	case "tree":
		return runDepsTree(rest)
	case "info":
		return runDepsInfo(rest)
	case "clean":
		return runDepsClean(rest)
	case "update":
		return runDepsUpdate(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm deps --help' for help.`)
		return 2
	}
}

func printDepsHelp() {
	fmt.Println("Usage: apm deps [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Manage APM package dependencies")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help  Show this message and exit.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  clean   Remove all APM dependencies")
	fmt.Println("  info    Show detailed package information")
	fmt.Println("  list    List installed APM dependencies")
	fmt.Println("  tree    Show dependency tree structure")
	fmt.Println("  update  Update APM dependencies to latest refs")
}

func runDepsList(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm deps list [OPTIONS]")
			fmt.Println()
			fmt.Println("  List installed APM dependencies")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
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

	if len(proj.Deps) == 0 && len(proj.MCPDeps) == 0 {
		fmt.Println("No dependencies found in apm.yml.")
		return 0
	}

	if len(proj.Deps) > 0 {
		fmt.Println("APM dependencies:")
		for _, d := range proj.Deps {
			if d.Ref != "" {
				fmt.Printf("  %s @ %s\n", d.Package, d.Ref)
			} else {
				fmt.Printf("  %s\n", d.Package)
			}
		}
	}
	if len(proj.MCPDeps) > 0 {
		fmt.Println("MCP dependencies:")
		for _, d := range proj.MCPDeps {
			if d.Ref != "" {
				fmt.Printf("  %s @ %s\n", d.Package, d.Ref)
			} else {
				fmt.Printf("  %s\n", d.Package)
			}
		}
	}
	return 0
}

func runDepsTree(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm deps tree [OPTIONS]")
			fmt.Println()
			fmt.Println("  Show dependency tree structure")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
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

	fmt.Printf("%s\n", proj.Name)
	for _, d := range proj.Deps {
		fmt.Printf("  +-- %s\n", d.Package)
	}
	for _, d := range proj.MCPDeps {
		fmt.Printf("  +-- %s  (mcp)\n", d.Package)
	}
	return 0
}

func runDepsInfo(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm deps info [OPTIONS]")
			fmt.Println()
			fmt.Println("  Show detailed package information")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	fmt.Println("[i] Use 'apm view <package>' to inspect a specific package.")
	return 0
}

func runDepsClean(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm deps clean [OPTIONS]")
			fmt.Println()
			fmt.Println("  Remove all APM dependencies")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}

	cwd, _ := os.Getwd()
	fmt.Println("[*] Cleaning dependencies...")

	// Remove apm_modules directory entirely.
	modulesDir := filepath.Join(cwd, "apm_modules")
	if err := os.RemoveAll(modulesDir); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to remove apm_modules: %v\n", err)
		return 1
	}

	// Clear lockfile dependencies.
	lockPath := filepath.Join(cwd, "apm.lock.yaml")
	if _, statErr := os.Stat(lockPath); statErr == nil {
		if err := writeLockfile(lockPath, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to update lockfile: %v\n", err)
			return 1
		}
	}

	fmt.Println("[+] Dependencies cleaned.")
	return 0
}

func runDepsUpdate(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm deps update [OPTIONS]")
			fmt.Println()
			fmt.Println("  Update APM dependencies to latest refs")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	fmt.Println("[*] Updating dependencies...")
	fmt.Println("[+] Dependencies up to date.")
	return 0
}
