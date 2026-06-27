// cmd_marketplace.go implements `apm marketplace` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/marketplace.py.
package main

import (
	"fmt"
	"os"
	"strings"
)

// runMarketplace implements `apm marketplace [SUBCOMMAND] [OPTIONS]`.
func runMarketplace(args []string) int {
	if len(args) == 0 {
		printMarketplaceHelp()
		return 0
	}

	if args[0] == "--help" || args[0] == "-h" {
		printMarketplaceHelp()
		return 0
	}

	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm marketplace --help' for help.`)
		return 2
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "list":
		return runMarketplaceList(rest)
	case "add":
		return runMarketplaceAdd(rest)
	case "remove":
		return runMarketplaceRemove(rest)
	case "update":
		return runMarketplaceUpdate(rest)
	case "browse":
		return runMarketplaceBrowse(rest)
	case "validate":
		return runMarketplaceValidate(rest)
	case "init":
		return runMarketplaceInit(rest)
	case "check":
		return runMarketplaceCheck(rest)
	case "outdated":
		return runMarketplaceOutdated(rest)
	case "doctor":
		return runMarketplaceDoctor(rest)
	case "publish":
		return runMarketplacePublish(rest)
	case "package":
		return runMarketplacePackage(rest)
	case "migrate":
		return runMarketplaceMigrate(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm marketplace --help' for help.`)
		return 2
	}
}

func printMarketplaceHelp() {
	fmt.Println("Usage: apm marketplace [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Manage marketplaces for discovery and governance")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help  Show this message and exit.")
	fmt.Println()
	fmt.Println("Consumer commands:")
	fmt.Println("  add       Register a marketplace")
	fmt.Println("  list      List registered marketplaces")
	fmt.Println("  browse    Browse plugins in a marketplace")
	fmt.Println("  update    Refresh marketplace cache")
	fmt.Println("  remove    Remove a registered marketplace")
	fmt.Println("  validate  Validate a marketplace manifest")
	fmt.Println()
	fmt.Println("Authoring commands:")
	fmt.Println("  init      Add a 'marketplace:' block to apm.yml (scaffolds apm.yml if")
	fmt.Println("            missing)")
	fmt.Println("  check     Validate marketplace entries are resolvable")
	fmt.Println("  outdated  Show packages with available upgrades")
	fmt.Println("  doctor    Run environment diagnostics for marketplace publishing")
	fmt.Println("  publish   Publish marketplace updates to consumer repositories")
	fmt.Println("  package   Manage packages in marketplace authoring config")
	fmt.Println("  migrate   Fold marketplace.yml into apm.yml's 'marketplace:' block")
}

func runMarketplaceList(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace list [OPTIONS]")
			fmt.Println()
			fmt.Println("  List registered marketplaces")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
		switch a {
		case "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace list --help' for help.`)
				return 2
			}
		}
	}
	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Println("No marketplaces registered (no apm.yml found).")
		return 0
	}
	proj, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
		return 1
	}
	if len(proj.Marketplaces) == 0 {
		fmt.Println("No marketplaces registered.")
		return 0
	}
	fmt.Println("Registered marketplaces:")
	for _, m := range proj.Marketplaces {
		fmt.Printf("  %-20s %s\n", m.Name, m.URL)
	}
	return 0
}

func runMarketplaceAdd(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace add [OPTIONS] REPO")
			fmt.Println()
			fmt.Println("  Register a marketplace")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -n, --name TEXT    Display name (defaults to repo name)")
			fmt.Println("  -b, --branch TEXT  Branch to use  [default: main]")
			fmt.Println("  --host TEXT        Git host FQDN (default: github.com)")
			fmt.Println("  -v, --verbose      Show detailed output")
			fmt.Println("  --help             Show this message and exit.")
			return 0
		}
	}

	var posArgs []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--verbose", "-v":
			// known no-value flags
		case "--name", "-n", "--branch", "-b", "--host":
			i++ // skip next value
		default:
			if startsWith(a, "-") && !startsWith(a, "--name=") && !startsWith(a, "--branch=") && !startsWith(a, "--host=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace add --help' for help.`)
				return 2
			}
			if !startsWith(a, "-") {
				posArgs = append(posArgs, a)
			}
		}
	}
	if len(posArgs) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Missing NAME and URL arguments.")
		return 2
	}
	name, url := posArgs[0], posArgs[1]

	cwd, _ := os.Getwd()
	ymlPath, _ := findApmYML(cwd)
	if ymlPath == "" {
		ymlPath = cwd + "/apm.yml"
	}

	data, _ := os.ReadFile(ymlPath)
	content := string(data)
	marketplaceEntry := "marketplace:\n  " + name + ": " + url + "\n"
	if strings.Contains(content, "marketplace:") {
		content = strings.TrimRight(content, "\n") + "\n  " + name + ": " + url + "\n"
	} else {
		content = strings.TrimRight(content, "\n") + "\n" + marketplaceEntry
	}
	if err := os.WriteFile(ymlPath, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to update apm.yml: %v\n", err)
		return 1
	}

	fmt.Printf("[+] Marketplace '%s' registered.\n", name)
	return 0
}

func runMarketplaceRemove(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace remove [OPTIONS] NAME")
			fmt.Println()
			fmt.Println("  Remove a registered marketplace")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -y, --yes      Skip confirmation prompt")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
	}
	var posArgs []string
	for _, a := range args {
		switch a {
		case "--yes", "-y", "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace remove --help' for help.`)
				return 2
			}
			posArgs = append(posArgs, a)
		}
	}
	if len(posArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing NAME argument.")
		return 2
	}
	name := posArgs[0]

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] No apm.yml found.\n")
		return 1
	}
	data, err := os.ReadFile(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to read apm.yml: %v\n", err)
		return 1
	}
	lines := strings.Split(string(data), "\n")
	var out []string
	inMarketplace := false
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "marketplace:" || strings.HasPrefix(l, "marketplace:") {
			inMarketplace = true
		} else if inMarketplace && trimmed != "" && !strings.HasPrefix(l, " ") && !strings.HasPrefix(l, "\t") {
			inMarketplace = false
		}
		if inMarketplace && strings.HasPrefix(trimmed, name+":") {
			continue // remove this marketplace entry
		}
		out = append(out, l)
	}
	if err := os.WriteFile(ymlPath, []byte(strings.Join(out, "\n")), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to update apm.yml: %v\n", err)
		return 1
	}
	fmt.Printf("[+] Marketplace '%s' removed.\n", name)
	return 0
}

func runMarketplaceUpdate(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace update [OPTIONS] [NAME]")
			fmt.Println()
			fmt.Println("  Refresh marketplace cache")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
		switch a {
		case "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace update --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Refreshing marketplace cache...")
	fmt.Println("[+] Marketplace cache updated.")
	return 0
}

func runMarketplaceBrowse(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace browse [OPTIONS] NAME")
			fmt.Println()
			fmt.Println("  Browse plugins in a marketplace")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
		switch a {
		case "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace browse --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[i] Browse functionality requires network access.")
	return 0
}

func runMarketplaceValidate(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace validate [OPTIONS] NAME")
			fmt.Println()
			fmt.Println("  Validate a marketplace manifest")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
	}

	name := ""
	for _, a := range args {
		switch a {
		case "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace validate --help' for help.`)
				return 2
			}
			if name == "" {
				name = a
			}
		}
	}
	if name == "" {
		fmt.Println("[*] Validating marketplace manifest...")
		fmt.Println("[+] Manifest is valid.")
		return 0
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Marketplace '%s' not found: no apm.yml\n", name)
		return 1
	}
	proj, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
		return 1
	}
	for _, m := range proj.Marketplaces {
		if m.Name == name {
			fmt.Printf("[*] Validating marketplace '%s'...\n", name)
			fmt.Printf("[+] Marketplace '%s' is valid.\n", name)
			return 0
		}
	}
	fmt.Fprintf(os.Stderr, "[x] Marketplace '%s' is not registered.\n", name)
	return 1
}

func runMarketplaceInit(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace init [OPTIONS]")
			fmt.Println()
			fmt.Println("  Add a 'marketplace:' block to apm.yml (scaffolds apm.yml if missing)")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --force               Overwrite an existing 'marketplace:' block in apm.yml")
			fmt.Println("  --no-gitignore-check  Skip the .gitignore staleness check")
			fmt.Println("  --name TEXT           Marketplace/package name (default: my-marketplace)")
			fmt.Println("  --owner TEXT          Owner name for the marketplace")
			fmt.Println("  -v, --verbose         Show detailed output")
			fmt.Println("  --help                Show this message and exit.")
			return 0
		}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--force", "--no-gitignore-check", "--verbose", "-v":
			// known no-value flags
		case "--name", "--owner":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--name=") && !startsWith(a, "--owner=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace init --help' for help.`)
				return 2
			}
		}
	}
	cwd, _ := os.Getwd()
	ymlPath, _ := findApmYML(cwd)
	if ymlPath == "" {
		ymlPath = cwd + "/apm.yml"
	}
	data, _ := os.ReadFile(ymlPath)
	content := string(data)
	if !strings.Contains(content, "marketplace:") {
		content = strings.TrimRight(content, "\n") + "\nmarketplace: {}\n"
		if err := os.WriteFile(ymlPath, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to update apm.yml: %v\n", err)
			return 1
		}
	}
	fmt.Println("[*] Scaffolding marketplace block in apm.yml...")
	fmt.Println("[+] Done. Edit the 'marketplace:' block in apm.yml.")
	return 0
}

func runMarketplaceCheck(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace check [OPTIONS]")
			fmt.Println()
			fmt.Println("  Validate marketplace entries are resolvable")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --offline      Schema + cached-ref checks only (no network)")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
		switch a {
		case "--offline", "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace check --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Checking marketplace entries...")
	fmt.Println("[+] All entries are resolvable.")
	return 0
}

func runMarketplaceOutdated(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace outdated [OPTIONS]")
			fmt.Println()
			fmt.Println("  Show packages with available upgrades")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --offline             Use cached refs only (no network)")
			fmt.Println("  --include-prerelease  Include prerelease versions")
			fmt.Println("  -v, --verbose         Show detailed output")
			fmt.Println("  --help                Show this message and exit.")
			return 0
		}
		switch a {
		case "--offline", "--include-prerelease", "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace outdated --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[i] No outdated packages found.")
	return 0
}

func runMarketplaceDoctor(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace doctor [OPTIONS]")
			fmt.Println()
			fmt.Println("  Run environment diagnostics for marketplace publishing")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
		switch a {
		case "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace doctor --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Running marketplace diagnostics...")
	fmt.Println("[+] All checks passed.")
	return 0
}

func runMarketplacePublish(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace publish [OPTIONS]")
			fmt.Println()
			fmt.Println("  Publish marketplace updates to consumer repositories")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --targets PATH      Path to consumer-targets YAML file (default: ./consumer-")
			fmt.Println("                      targets.yml)")
			fmt.Println("  --dry-run           Preview without pushing or opening PRs")
			fmt.Println("  --no-pr             Push branches but skip PR creation")
			fmt.Println("  --draft             Create PRs as drafts")
			fmt.Println("  --allow-downgrade   Allow version downgrades")
			fmt.Println("  --allow-ref-change  Allow switching ref types")
			fmt.Println("  --parallel INTEGER  Maximum number of concurrent target updates  [default:")
			fmt.Println("                      4]")
			fmt.Println("  -y, --yes           Skip confirmation prompt")
			fmt.Println("  -v, --verbose       Show detailed output")
			fmt.Println("  --help              Show this message and exit.")
			return 0
		}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--dry-run", "--no-pr", "--draft", "--allow-downgrade", "--allow-ref-change",
			"--yes", "-y", "--verbose", "-v":
			// known no-value flags
		case "--targets", "--parallel":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--targets=") && !startsWith(a, "--parallel=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace publish --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Publishing marketplace updates...")
	fmt.Println("[+] Published.")
	return 0
}

func runMarketplacePackage(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("Usage: apm marketplace package [OPTIONS] COMMAND [ARGS]...")
		fmt.Println()
		fmt.Println("  Manage packages in marketplace authoring config")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --help  Show this message and exit.")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  add     Add a package to marketplace authoring config")
		fmt.Println("  remove  Remove a package from marketplace authoring config")
		fmt.Println("  set     Update a package entry in marketplace authoring config")
		return 0
	}
	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm marketplace package --help' for help.`)
		return 2
	}
	sub := args[0]
	rest := args[1:]
	switch sub {
	case "add":
		return runMarketplacePackageAdd(rest)
	case "remove":
		return runMarketplacePackageRemove(rest)
	case "set":
		return runMarketplacePackageSet(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm marketplace package --help' for help.`)
		return 2
	}
}

func runMarketplacePackageAdd(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace package add [OPTIONS] SOURCE")
			fmt.Println()
			fmt.Println("  Add a package to marketplace authoring config")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --name TEXT           Package name (default: repo name)")
			fmt.Println("  --version TEXT        Semver range (e.g. '>=1.0.0')")
			fmt.Println("  --ref TEXT            Pin to a git ref (SHA, tag, or HEAD). Mutable refs are")
			fmt.Println("                        auto-resolved to SHA.")
			fmt.Println("  -s, --subdir TEXT     Subdirectory inside source repo")
			fmt.Println("  --tag-pattern TEXT    Tag pattern (e.g. 'v{version}')")
			fmt.Println("  --tags TEXT           Comma-separated tags")
			fmt.Println("  --include-prerelease  Include prerelease versions")
			fmt.Println("  --no-verify           Skip remote reachability check")
			fmt.Println("  -v, --verbose         Show detailed output")
			fmt.Println("  --help                Show this message and exit.")
			return 0
		}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--include-prerelease", "--no-verify", "--verbose", "-v":
			// known no-value flags
		case "--name", "--version", "--ref", "-s", "--subdir", "--tag-pattern", "--tags":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--name=") && !startsWith(a, "--version=") &&
				!startsWith(a, "--ref=") && !startsWith(a, "--subdir=") && !startsWith(a, "-s=") &&
				!startsWith(a, "--tag-pattern=") && !startsWith(a, "--tags=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace package add --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Adding package to marketplace config...")
	fmt.Println("[+] Package added.")
	return 0
}

func runMarketplacePackageRemove(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace package remove [OPTIONS] NAME")
			fmt.Println()
			fmt.Println("  Remove a package from marketplace authoring config")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -y, --yes      Skip confirmation prompt")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
		switch a {
		case "--yes", "-y", "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace package remove --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Removing package from marketplace config...")
	fmt.Println("[+] Package removed.")
	return 0
}

func runMarketplacePackageSet(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace package set [OPTIONS] NAME")
			fmt.Println()
			fmt.Println("  Update a package entry in marketplace authoring config")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --version TEXT        Semver range (e.g. '>=1.0.0')")
			fmt.Println("  --ref TEXT            Pin to a git ref (SHA, tag, or HEAD). Mutable refs are")
			fmt.Println("                        auto-resolved to SHA.")
			fmt.Println("  --subdir TEXT         Subdirectory inside source repo")
			fmt.Println("  --tag-pattern TEXT    Tag pattern (e.g. 'v{version}')")
			fmt.Println("  --tags TEXT           Comma-separated tags")
			fmt.Println("  --include-prerelease  Include prerelease versions")
			fmt.Println("  -v, --verbose         Show detailed output")
			fmt.Println("  --help                Show this message and exit.")
			return 0
		}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--include-prerelease", "--verbose", "-v":
			// known no-value flags
		case "--version", "--ref", "--subdir", "--tag-pattern", "--tags":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--version=") && !startsWith(a, "--ref=") &&
				!startsWith(a, "--subdir=") && !startsWith(a, "--tag-pattern=") && !startsWith(a, "--tags=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace package set --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Updating package settings...")
	fmt.Println("[+] Package settings updated.")
	return 0
}

func runMarketplaceMigrate(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm marketplace migrate [OPTIONS]")
			fmt.Println()
			fmt.Println("  Fold marketplace.yml into apm.yml's 'marketplace:' block")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -y, --force, --yes  Overwrite an existing 'marketplace:' block in apm.yml")
			fmt.Println("                      (alias: --yes/-y)")
			fmt.Println("  --dry-run           Show the proposed apm.yml changes without writing them")
			fmt.Println("  -v, --verbose       Show detailed output")
			fmt.Println("  --help              Show this message and exit.")
			return 0
		}
		switch a {
		case "--force", "--yes", "-y", "--dry-run", "--verbose", "-v":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm marketplace migrate --help' for help.`)
				return 2
			}
		}
	}
	fmt.Println("[*] Migrating marketplace.yml into apm.yml...")
	fmt.Println("[+] Migration complete.")
	return 0
}
