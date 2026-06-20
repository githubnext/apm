// cmd_plugin.go implements `apm plugin` group for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/plugin/__init__.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var pluginSubcommands = []struct{ name, desc string }{
	{"init", "Scaffold a plugin (creates plugin.json + apm.yml)"},
}

func printPluginHelp() {
	fmt.Println("Usage: apm plugin [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Scaffold and manage plugins (plugin-author workflows)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help  Show this message and exit.")
	fmt.Println()
	fmt.Println("Commands:")
	for _, sub := range pluginSubcommands {
		fmt.Printf("  %-6s%s\n", sub.name, sub.desc)
	}
}

// runPlugin implements `apm plugin [SUBCOMMAND] [OPTIONS]`.
func runPlugin(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printPluginHelp()
		return 0
	}

	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm plugin --help' for help.`)
		return 2
	}

	sub := args[0]
	rest := args[1:]
	switch sub {
	case "init":
		return runPluginInit(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm plugin --help' for help.`)
		return 2
	}
}

func runPluginInit(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm plugin init [OPTIONS]")
			fmt.Println()
			fmt.Println("  Scaffold a new plugin (plugin.json + apm.yml)")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --yes, -y  Skip confirmation prompt")
			fmt.Println("  --target TEXT  Target harness")
			fmt.Println("  --verbose, -v  Show detailed output")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--yes", "-y", "--verbose", "-v":
			// known no-value flags
		case "--target":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--target=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm plugin init --help' for help.`)
				return 2
			}
		}
	}
	cwd, _ := os.Getwd()
	fmt.Printf("[*] Scaffolding plugin in: %s\n", cwd)

	pluginJSON := `{
  "name": "my-plugin",
  "version": "0.1.0",
  "description": "APM plugin",
  "main": "index.js"
}
`
	pluginPath := filepath.Join(cwd, "plugin.json")
	if err := os.WriteFile(pluginPath, []byte(pluginJSON), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to write plugin.json: %v\n", err)
		return 1
	}

	apmYML := `name: my-plugin
version: 0.1.0
description: APM plugin
type: plugin
targets:
  - copilot
dependencies:
  apm: []
  mcp: []
`
	apmPath := filepath.Join(cwd, "apm.yml")
	if _, statErr := os.Stat(apmPath); os.IsNotExist(statErr) {
		if err := os.WriteFile(apmPath, []byte(apmYML), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to write apm.yml: %v\n", err)
			return 1
		}
	} else {
		if err := appendToApmYML(apmPath, "type: plugin"); err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to update apm.yml: %v\n", err)
			return 1
		}
	}

	fmt.Println("[+] Plugin scaffolded.")
	return 0
}
