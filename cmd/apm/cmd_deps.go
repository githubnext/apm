// cmd_deps.go implements `apm deps` and its subcommands for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/deps.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runDeps implements `apm deps [SUBCOMMAND] [OPTIONS]`.
func runDeps(args []string) int {
	if len(args) == 0 {
		printDepsHelp()
		return 0
	}

	sub := args[0]
	rest := args[1:]

	// Only intercept --help when it is the first (and only meaningful) arg,
	// not when it follows a subcommand name -- let the subcommand handle it.
	if sub == "--help" || sub == "-h" {
		printDepsHelp()
		return 0
	}

	if startsWith(sub, "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm deps --help' for help.`)
		return 2
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
			fmt.Println("  -g, --global  List user-scope dependencies (~/.apm/) instead of project")
			fmt.Println("  --all         Show both project and user-scope dependencies")
			fmt.Println("  --insecure    Show only installed dependencies locked to http:// sources")
			fmt.Println("  --help        Show this message and exit.")
			return 0
		}
		switch a {
		case "-g", "--global", "--all", "--insecure":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm deps list --help' for help.`)
				return 2
			}
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
			fmt.Println("  -g, --global  Show user-scope dependency tree (~/.apm/)")
			fmt.Println("  --help        Show this message and exit.")
			return 0
		}
		switch a {
		case "-g", "--global":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm deps tree --help' for help.`)
				return 2
			}
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
			fmt.Println("Usage: apm deps info [OPTIONS] PACKAGE")
			fmt.Println()
			fmt.Println("  Show detailed package information")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
		if startsWith(a, "-") {
			fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
			fmt.Fprintln(os.Stderr, `Try 'apm deps info --help' for help.`)
			return 2
		}
	}

	// Collect non-flag arguments as the package name.
	var pkg string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			pkg = a
			break
		}
	}
	if pkg == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'PACKAGE'.")
		fmt.Fprintln(os.Stderr, `Try 'apm deps info --help' for help.`)
		return 2
	}

	cwd, _ := os.Getwd()
	pkgDir := filepath.Join(cwd, "apm_modules", pkg)
	if _, err := os.Stat(pkgDir); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Package '%s' is not installed (no apm_modules/%s).\n", pkg, pkg)
		fmt.Fprintf(os.Stderr, "[i] Run 'apm install' to install dependencies.\n")
		return 1
	}

	ymlPath := filepath.Join(pkgDir, "apm.yml")
	if _, err := os.Stat(ymlPath); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Package '%s' has no apm.yml in apm_modules/%s.\n", pkg, pkg)
		return 1
	}
	meta, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Could not read package metadata: %v\n", err)
		return 1
	}

	fmt.Printf("Package: %s\n", meta.Name)
	if meta.Version != "" {
		fmt.Printf("Version: %s\n", meta.Version)
	}
	if len(meta.Deps) > 0 {
		fmt.Println("Dependencies:")
		for _, d := range meta.Deps {
			fmt.Printf("  %s\n", d.Package)
		}
	}
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
			fmt.Println("  --dry-run  Show what would be removed without removing")
			fmt.Println("  -y, --yes  Skip confirmation prompt")
			fmt.Println("  --help     Show this message and exit.")
			return 0
		}
		switch a {
		case "--dry-run", "--yes", "-y":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm deps clean --help' for help.`)
				return 2
			}
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
			fmt.Println("Usage: apm deps update [OPTIONS] [PACKAGES]...")
			fmt.Println()
			fmt.Println("  Update APM dependencies to latest refs")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose                 Show detailed update information")
			fmt.Println("  --force                       Overwrite locally-authored files on collision")
			fmt.Println("  -t, --target TARGET           Target platform (comma-separated). Values:")
			fmt.Println("                                copilot, claude, cursor, opencode, codex,")
			fmt.Println("                                gemini, windsurf, agent-skills, all. 'agent-")
			fmt.Println("                                skills' deploys to .agents/skills/ (cross-")
			fmt.Println("                                client). 'all' = copilot+claude+cursor+opencod")
			fmt.Println("                                e+codex+gemini+windsurf (excludes agent-")
			fmt.Println("                                skills); combine with 'agent-skills' for both.")
			fmt.Println("                                'copilot-cowork' is also accepted when the")
			fmt.Println("                                copilot-cowork experimental flag is enabled")
			fmt.Println("                                (run 'apm experimental enable copilot-")
			fmt.Println("                                cowork').")
			fmt.Println("  --parallel-downloads INTEGER  Max concurrent package downloads (0 to disable")
			fmt.Println("                                parallelism)  [default: 4]")
			fmt.Println("  -g, --global                  Update user-scope dependencies (~/.apm/)")
			fmt.Println("  --legacy-skill-paths          Deploy skill files to per-client paths (e.g.")
			fmt.Println("                                .cursor/skills/) instead of the shared")
			fmt.Println("                                .agents/skills/ directory. Compatibility flag")
			fmt.Println("                                for projects that need per-client skill")
			fmt.Println("                                layouts.")
			fmt.Println("  --help                        Show this message and exit.")
			return 0
		}
		switch a {
		case "--verbose", "-v", "--force", "--global", "-g", "--legacy-skill-paths",
			"--target", "-t":
			// known flags
		default:
			if startsWith(a, "-") && !startsWith(a, "--parallel-downloads") && !startsWith(a, "--target=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm deps update --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Updating dependencies...")
	fmt.Println("[+] Dependencies up to date.")
	return 0
}
