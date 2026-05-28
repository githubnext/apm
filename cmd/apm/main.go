// cmd/apm is the entry point for the APM CLI (Go rewrite).
// Agent Package Manager (APM) -- Go implementation.
package main

import (
	"fmt"
	"os"
	"strings"
)

// version mirrors the Python CLI version for parity. The build tag may
// override this at link time.
const version = "0.14.1"

// commandOrder defines the display order for the top-level help (matches Python CLI).
var commandOrder = []string{
	"audit", "cache", "compile", "config", "deps", "experimental",
	"init", "install", "list", "marketplace", "mcp", "outdated", "pack",
	"plugin", "policy", "preview", "prune", "run", "runtime", "search",
	"self-update", "targets", "uninstall", "unpack", "update", "view",
}

// commands maps each command name to its one-line description (matches Python CLI).
var commands = map[string]string{
	"audit":        "Scan installed packages for hidden Unicode characters",
	"cache":        "Manage the local package cache",
	"compile":      "Compile APM context into distributed AGENTS.md files",
	"config":       "Configure APM CLI",
	"deps":         "Manage APM package dependencies",
	"experimental": "Manage experimental feature flags",
	"init":         "Initialize a new APM project",
	"install":      "Install APM and MCP dependencies (supports APM packages,...",
	"list":         "List available scripts in the current project",
	"marketplace":  "Manage marketplaces for discovery and governance",
	"mcp":          "Discover, inspect, and install MCP servers",
	"outdated":     "Show outdated locked dependencies",
	"pack":         "Pack distributable artifacts from your APM project.",
	"plugin":       "Scaffold and manage plugins (plugin-author workflows)",
	"policy":       "Inspect and diagnose APM policy",
	"preview":      "Preview a script's compiled prompt files",
	"prune":        "Remove APM packages not listed in apm.yml",
	"run":          "Run a script with parameters (experimental)",
	"runtime":      "Manage AI runtimes (experimental)",
	"search":       "Search plugins in a marketplace (QUERY@MARKETPLACE)",
	"self-update":  "Update the APM CLI binary itself to the latest version",
	"targets":      "Show resolved targets for the current project.",
	"uninstall":    "Remove APM packages, their integrated files, and apm.yml...",
	"unpack":       "[Deprecated] Extract an APM bundle into the current project.",
	"update":       "Refresh APM dependencies to the latest matching refs",
	"view":         "View package metadata or list remote versions.",
}

// aliases maps legacy or alternate names to canonical commands.
var aliases = map[string]string{
	"info":        "view",
	"self_update": "self-update",
}

func printHelp() {
	fmt.Println("Usage: apm [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Agent Package Manager (APM): The package manager for AI-Native Development")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --version  Show version and exit.")
	fmt.Println("  --help     Show this message and exit.")
	fmt.Println()
	fmt.Println("Commands:")
	for _, name := range commandOrder {
		desc := commands[name]
		fmt.Printf("  %-14s%s\n", name, desc)
	}
}

// isGroupCmd returns true for commands that have subcommands and manage their own --help.
func isGroupCmd(name string) bool {
	switch name {
	case "cache", "deps", "marketplace", "mcp", "policy", "runtime", "plugin", "experimental", "config":
		return true
	}
	return false
}

func printCmdHelp(name string) {
	canonical := name
	if a, ok := aliases[name]; ok {
		canonical = a
	}
	desc, ok := commands[canonical]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", name)
		fmt.Fprintln(os.Stderr, `Try 'apm --help' for help.`)
		os.Exit(2)
	}
	// Full description for selected commands (matches Python Click output).
	fullDesc := commandFullDesc[canonical]
	if fullDesc == "" {
		fullDesc = desc
	}
	fmt.Printf("Usage: apm %s [OPTIONS]", canonical)
	// Commands with positional args -- match Python CLI usage strings exactly.
	switch canonical {
	case "install":
		fmt.Printf(" [PACKAGES]...")
	case "uninstall":
		fmt.Printf(" PACKAGES...")
	case "view":
		fmt.Printf(" PACKAGE [FIELD]")
	case "search":
		fmt.Printf(" QUERY@MARKETPLACE")
	case "run":
		fmt.Printf(" [SCRIPT_NAME]")
	case "preview":
		fmt.Printf(" [SCRIPT_NAME]")
	case "audit":
		fmt.Printf(" [PACKAGE]")
	case "unpack":
		fmt.Printf(" BUNDLE_PATH")
	case "init":
		fmt.Printf(" [PROJECT_NAME]")
	case "targets":
		fmt.Printf(" COMMAND [ARGS]...")
	}
	fmt.Println()
	fmt.Println()
	// Print full description: indent each line with two spaces.
	for _, line := range strings.Split(fullDesc, "\n") {
		if line == "" {
			fmt.Println()
		} else {
			fmt.Printf("  %s\n", line)
		}
	}
	fmt.Println()
	fmt.Println("Options:")
	if opts, ok := commandOptions[canonical]; ok {
		for _, opt := range opts {
			fmt.Println(opt)
		}
	} else {
		// Default: only --help option.
		fmt.Println("  --help  Show this message and exit.")
	}
}

func run(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 0
	}

	var subArgs []string
	showVersion := false
	showHelp := false

	for i := 0; i < len(args); {
		a := args[i]
		switch {
		case a == "--version" || a == "-version":
			showVersion = true
			i++
		case a == "--help" || a == "-help" || a == "-h":
			showHelp = true
			i++
		default:
			subArgs = append(subArgs, args[i:]...)
			i = len(args)
		}
	}

	if showVersion {
		fmt.Printf("Agent Package Manager (APM) CLI version %s (go)\n", version)
		return 0
	}
	if showHelp && len(subArgs) == 0 {
		printHelp()
		return 0
	}
	if len(subArgs) == 0 {
		printHelp()
		return 0
	}

	cmd := subArgs[0]
	rest := subArgs[1:]

	if cmd == "help" {
		if len(rest) == 0 {
			printHelp()
			return 0
		}
		printCmdHelp(rest[0])
		return 0
	}

	if canonical, ok := aliases[cmd]; ok {
		cmd = canonical
	}

	if _, ok := commands[cmd]; !ok {
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", cmd)
		fmt.Fprintln(os.Stderr, `Try 'apm --help' for help.`)
		return 2
	}

	for _, a := range rest {
		if a == "--help" || a == "-h" || a == "-help" {
			// Group commands handle their own --help (they list subcommands).
			if isGroupCmd(cmd) {
				break
			}
			printCmdHelp(cmd)
			return 0
		}
	}

	// Dispatch to implemented commands.
	switch cmd {
	case "init":
		return runInit(rest)
	case "targets":
		return runTargets(rest)
	case "list":
		return runList(rest)
	case "deps":
		return runDeps(rest)
	case "cache":
		return runCache(rest)
	case "config":
		return runConfig(rest)
	case "view":
		return runView(rest)
	case "marketplace":
		return runMarketplace(rest)
	case "compile":
		return runCompile(rest)
	case "pack":
		return runPack(rest)
	case "unpack":
		return runUnpack(rest)
	case "install":
		return runInstall(rest)
	case "uninstall":
		return runUninstall(rest)
	case "update":
		return runUpdate(rest)
	case "prune":
		return runPrune(rest)
	case "audit":
		return runAudit(rest)
	case "policy":
		return runPolicy(rest)
	case "runtime":
		return runRuntime(rest)
	case "mcp":
		return runMCP(rest)
	case "plugin":
		return runPlugin(rest)
	case "search":
		return runSearch(rest)
	case "run":
		return runRun(rest)
	case "outdated":
		return runOutdated(rest)
	case "self-update":
		return runSelfUpdate(rest)
	case "experimental":
		return runExperimental(rest)
	case "preview":
		return runPreview(rest)
	}

	// Commands not yet fully wired to Go business logic.
	fmt.Fprintf(os.Stderr, "apm: %s command is not yet implemented in the Go rewrite.\n", cmd)
	fmt.Fprintf(os.Stderr, "Run 'apm --help' for usage.\n")
	_ = strings.Join(rest, " ")
	return 1
}

func main() {
	os.Exit(run(os.Args[1:]))
}

