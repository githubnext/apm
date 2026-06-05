// cmd_view.go implements `apm view` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/view.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runView implements `apm view [OPTIONS] PACKAGE [FIELD]`.
func runView(args []string) int {
	var flagGlobal, flagHelp bool
	var posArgs []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			flagHelp = true
		case "--global", "-g":
			flagGlobal = true
		default:
			if !startsWith(args[i], "-") {
				posArgs = append(posArgs, args[i])
			}
		}
	}

	if flagHelp {
		printCmdHelp("view")
		return 0
	}

	if len(posArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'PACKAGE'.")
		fmt.Fprintln(os.Stderr, `Try 'apm view --help' for help.`)
		return 2
	}

	pkg := posArgs[0]
	var field string
	if len(posArgs) > 1 {
		field = posArgs[1]
	}

	// Reject path traversal or absolute paths.
	if strings.Contains(filepath.Clean(pkg), "..") || filepath.IsAbs(pkg) {
		fmt.Fprintf(os.Stderr, "[x] Invalid package path: %q (path traversal not allowed)\n", pkg)
		return 2
	}

	// Resolve install directory.
	var installBase string
	if flagGlobal {
		home, _ := os.UserHomeDir()
		installBase = filepath.Join(home, ".apm", "packages")
	} else {
		cwd, _ := os.Getwd()
		installBase = filepath.Join(cwd, ".apm", "packages")
	}

	pkgDir := filepath.Join(installBase, pkg)
	pkgYML := filepath.Join(pkgDir, "apm.yml")

	if _, err := os.Stat(pkgYML); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Package '%s' is not installed.\n", pkg)
		fmt.Fprintf(os.Stderr, "    Looked in: %s\n", pkgDir)
		return 1
	}

	proj, err := parseApmYML(pkgYML)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to read package metadata: %v\n", err)
		return 1
	}

	switch field {
	case "":
		fmt.Printf("Package:     %s\n", pkg)
		fmt.Printf("Name:        %s\n", proj.Name)
		fmt.Printf("Version:     %s\n", proj.Version)
		fmt.Printf("Description: %s\n", proj.Description)
		fmt.Printf("Author:      %s\n", proj.Author)
	case "versions":
		fmt.Printf("[i] Remote version listing requires network access.\n")
		fmt.Printf("    Package: %s\n", pkg)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown field '%s'. Available: versions\n", field)
		return 2
	}
	return 0
}

// startsWith is a simple prefix check to avoid importing strings in every file.
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
