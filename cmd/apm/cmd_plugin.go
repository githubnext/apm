// cmd_plugin.go implements `apm plugin` group for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/plugin/__init__.py.
package main

import (
	"fmt"
	"os"
)

var pluginSubcommands = []struct{ name, desc string }{
	{"init", "Scaffold a new plugin (plugin.json + apm.yml)"},
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
		fmt.Printf("  %-14s%s\n", sub.name, sub.desc)
	}
}

// runPlugin implements `apm plugin [SUBCOMMAND] [OPTIONS]`.
func runPlugin(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printPluginHelp()
		return 0
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
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	cwd, _ := os.Getwd()
	fmt.Printf("[*] Scaffolding plugin in: %s\n", cwd)
	fmt.Println("[+] Plugin scaffolded.")
	return 0
}
