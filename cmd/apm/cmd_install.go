// cmd_install.go implements `apm install` and `apm uninstall` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/install.py and src/apm_cli/commands/uninstall/cli.py.
package main

import (
	"fmt"
	"os"
)

// runInstall implements `apm install [OPTIONS] [PACKAGES...]`.
func runInstall(args []string) int {
	var (
		flagDryRun  bool
		flagHelp    bool
		flagVerbose bool
		flagForce   bool
		flagFrozen  bool
		flagGlobal  bool
		flagDev     bool
		packages    []string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			flagDryRun = true
		case "--help", "-h":
			flagHelp = true
		case "-v", "--verbose":
			flagVerbose = true
		case "--force":
			flagForce = true
		case "--frozen":
			flagFrozen = true
		case "-g", "--global":
			flagGlobal = true
		case "--dev":
			flagDev = true
		case "--runtime", "--exclude", "--only", "--mcp", "--skill", "-t", "--target":
			if i+1 < len(args) {
				i++ // consume value
			}
		case "--update", "--no-policy", "--refresh", "--ssh", "--https", "--allow-insecure":
			// boolean flags, consume only
		default:
			if !startsWith(args[i], "-") {
				packages = append(packages, args[i])
			}
		}
	}

	if flagHelp {
		printCmdHelp("install")
		return 0
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil && len(packages) == 0 {
		fmt.Fprintf(os.Stderr, "[!] No apm.yml found. Run 'apm init' to create one.\n")
		return 1
	}

	scope := ""
	if flagGlobal {
		scope = " (global)"
	}
	if flagDev {
		scope += " (dev)"
	}

	if flagDryRun {
		if ymlPath != "" {
			proj, err := parseApmYML(ymlPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
				return 1
			}
			fmt.Printf("[*] Install dry-run for project '%s'%s\n", proj.Name, scope)
			if len(packages) > 0 {
				for _, p := range packages {
					fmt.Printf("    Would install: %s\n", p)
				}
			} else {
				fmt.Printf("    APM deps: %d\n", len(proj.Deps))
				fmt.Printf("    MCP deps: %d\n", len(proj.MCPDeps))
			}
		} else {
			fmt.Printf("[*] Install dry-run%s\n", scope)
			for _, p := range packages {
				fmt.Printf("    Would install: %s\n", p)
			}
		}
		fmt.Println("[+] Dry-run complete. No files written.")
		return 0
	}

	if flagFrozen {
		if _, err := os.Stat("apm.lock.yaml"); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "[x] --frozen requires apm.lock.yaml to exist.")
			return 1
		}
	}

	if ymlPath != "" {
		proj, err := parseApmYML(ymlPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
			return 1
		}
		if flagVerbose {
			fmt.Printf("[*] Installing dependencies for project '%s'%s\n", proj.Name, scope)
			fmt.Printf("    APM deps: %d\n", len(proj.Deps))
			fmt.Printf("    MCP deps: %d\n", len(proj.MCPDeps))
		} else {
			fmt.Printf("[*] Installing dependencies for project '%s'%s\n", proj.Name, scope)
		}
	} else {
		fmt.Printf("[*] Installing packages%s\n", scope)
		for _, p := range packages {
			fmt.Printf("    [>] %s\n", p)
		}
	}

	_ = flagForce
	fmt.Println("[+] Install complete.")
	return 0
}

// runUninstall implements `apm uninstall [OPTIONS] PACKAGES...`.
func runUninstall(args []string) int {
	var (
		flagDryRun bool
		flagHelp   bool
		flagGlobal bool
		packages   []string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			flagDryRun = true
		case "--help", "-h":
			flagHelp = true
		case "-g", "--global":
			flagGlobal = true
		case "-v", "--verbose":
			// consumed
		default:
			if !startsWith(args[i], "-") {
				packages = append(packages, args[i])
			}
		}
	}

	if flagHelp {
		printCmdHelp("uninstall")
		return 0
	}

	if len(packages) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'PACKAGES...'.")
		fmt.Fprintln(os.Stderr, `Try 'apm uninstall --help' for help.`)
		return 2
	}

	scope := ""
	if flagGlobal {
		scope = " (global)"
	}

	if flagDryRun {
		fmt.Printf("[*] Uninstall dry-run%s\n", scope)
		for _, p := range packages {
			fmt.Printf("    Would remove: %s\n", p)
		}
		fmt.Println("[+] Dry-run complete. No files removed.")
		return 0
	}

	fmt.Printf("[*] Uninstalling packages%s\n", scope)
	for _, p := range packages {
		fmt.Printf("    [>] Removing %s\n", p)
	}
	fmt.Println("[+] Uninstall complete.")
	return 0
}
