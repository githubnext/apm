// cmd/apm is the entry point for the APM CLI (Go rewrite).
// Agent Package Manager (APM) -- Go implementation.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const version = "0.1.0-go"

const helpText = `Agent Package Manager (APM): The package manager for AI-Native Development

Usage:
  apm [command]

Available Commands:
  audit       Audit installed packages for security issues
  cache       Manage the APM package cache
  compile     Compile APM primitives for a project
  config      View or set APM configuration values
  deps        Show or manage package dependencies
  experimental Access experimental features
  init        Initialize a new APM project
  install     Install packages
  list        List installed packages
  marketplace Browse the APM marketplace
  mcp         Manage MCP server integrations
  outdated    Show outdated packages
  pack        Pack a project into a distributable bundle
  plugin      Manage APM plugins
  policy      View or enforce APM policies
  preview     Preview changes before applying them
  prune       Remove unused packages
  run         Run a script or command via APM
  runtime     Manage runtimes
  search      Search the marketplace
  self-update Update APM itself
  targets     List available targets
  uninstall   Remove installed packages
  unpack      Unpack a bundle
  update      Update installed packages
  view        View package information

Flags:
  --help      Show this help and exit
  --version   Show version and exit

Use "apm [command] --help" for more information about a command.`

var commands = map[string]string{
	"audit":       "Audit installed packages for security issues",
	"cache":       "Manage the APM package cache",
	"compile":     "Compile APM primitives for a project",
	"config":      "View or set APM configuration values",
	"deps":        "Show or manage package dependencies",
	"experimental": "Access experimental features",
	"init":        "Initialize a new APM project",
	"install":     "Install packages",
	"list":        "List installed packages",
	"marketplace": "Browse the APM marketplace",
	"mcp":         "Manage MCP server integrations",
	"outdated":    "Show outdated packages",
	"pack":        "Pack a project into a distributable bundle",
	"plugin":      "Manage APM plugins",
	"policy":      "View or enforce APM policies",
	"preview":     "Preview changes before applying them",
	"prune":       "Remove unused packages",
	"run":         "Run a script or command via APM",
	"runtime":     "Manage runtimes",
	"search":      "Search the marketplace",
	"self-update": "Update APM itself",
	"targets":     "List available targets",
	"uninstall":   "Remove installed packages",
	"unpack":      "Unpack a bundle",
	"update":      "Update installed packages",
	"view":        "View package information",
}

// aliases maps legacy or alternate names to canonical commands.
var aliases = map[string]string{
	"info":        "view",
	"self_update": "self-update",
}

func cmdHelp(name string) {
	canonical := name
	if a, ok := aliases[name]; ok {
		canonical = a
	}
	desc, ok := commands[canonical]
	if !ok {
		fmt.Fprintf(os.Stderr, "apm: unknown command %q\n", name)
		fmt.Fprintln(os.Stderr, `Run "apm --help" for usage.`)
		os.Exit(1)
	}
	fmt.Printf("Usage:\n  apm %s [flags]\n\n%s\n\nFlags:\n  --help   Show this help and exit\n", canonical, desc)
}

func run(args []string) int {
	if len(args) == 0 {
		fmt.Println(helpText)
		return 0
	}

	// Top-level flags
	fs := flag.NewFlagSet("apm", flag.ContinueOnError)
	showVersion := fs.Bool("version", false, "Show version and exit")
	showHelp := fs.Bool("help", false, "Show help and exit")

	// Only parse flags that appear before any subcommand.
	// Collect the first non-flag arg as the subcommand.
	var subArgs []string
	i := 0
	for i < len(args) {
		a := args[i]
		if a == "--version" || a == "-version" {
			*showVersion = true
			i++
			continue
		}
		if a == "--help" || a == "-help" || a == "-h" {
			*showHelp = true
			i++
			continue
		}
		// Stop at first non-flag token.
		subArgs = append(subArgs, args[i:]...)
		break
	}

	if *showVersion {
		fmt.Printf("apm version %s (go)\n", version)
		return 0
	}
	if *showHelp && len(subArgs) == 0 {
		fmt.Println(helpText)
		return 0
	}

	if len(subArgs) == 0 {
		fmt.Println(helpText)
		return 0
	}

	cmd := subArgs[0]
	rest := subArgs[1:]

	// "help <command>" dispatches to per-command help.
	if cmd == "help" {
		if len(rest) == 0 {
			fmt.Println(helpText)
			return 0
		}
		cmdHelp(rest[0])
		return 0
	}

	// Resolve aliases.
	if canonical, ok := aliases[cmd]; ok {
		cmd = canonical
	}

	// Unknown command.
	if _, ok := commands[cmd]; !ok {
		fmt.Fprintf(os.Stderr, "apm: unknown command %q\n", cmd)
		fmt.Fprintln(os.Stderr, `Run "apm --help" for usage.`)
		return 1
	}

	// --help on subcommand.
	for _, a := range rest {
		if a == "--help" || a == "-h" || a == "-help" {
			cmdHelp(cmd)
			return 0
		}
	}

	// Subcommand stub: print informative not-yet-implemented message.
	fmt.Fprintf(os.Stderr, "apm %s: not yet fully implemented in the Go rewrite.\n", cmd)
	fmt.Fprintf(os.Stderr, "Use the Python APM CLI for production use: uv run apm %s %s\n",
		cmd, strings.Join(rest, " "))
	return 1
}

func main() {
	os.Exit(run(os.Args[1:]))
}

