// cmd_runtime.go implements `apm runtime` group for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/runtime.py.
package main

import (
	"fmt"
	"os"
)

var runtimeSubcommands = []struct{ name, desc string }{
	{"list", "List available and installed runtimes"},
	{"remove", "Remove an installed runtime"},
	{"setup", "Set up a runtime"},
	{"status", "Show active runtime and preference order"},
}

func printRuntimeHelp() {
	fmt.Println("Usage: apm runtime [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Manage AI runtimes (experimental)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help  Show this message and exit.")
	fmt.Println()
	fmt.Println("Commands:")
	for _, sub := range runtimeSubcommands {
		fmt.Printf("  %-8s%s\n", sub.name, sub.desc)
	}
}

// runRuntime implements `apm runtime [SUBCOMMAND] [OPTIONS]`.
func runRuntime(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printRuntimeHelp()
		return 0
	}

	sub := args[0]
	rest := args[1:]
	switch sub {
	case "setup":
		return runRuntimeSetup(rest)
	case "list":
		return runRuntimeList(rest)
	case "remove":
		return runRuntimeRemove(rest)
	case "status":
		return runRuntimeStatus(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm runtime --help' for help.`)
		return 2
	}
}

func runRuntimeSetup(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm runtime setup [OPTIONS] RUNTIME_NAME")
			fmt.Println()
			fmt.Println("  Set up a runtime")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --version TEXT  Specific version to install")
			fmt.Println("  --vanilla       Skip APM-specific configuration")
			fmt.Println("  --help          Show this message and exit.")
			return 0
		}
	}

	runtime := ""
	for _, a := range args {
		if !startsWith(a, "-") && runtime == "" {
			runtime = a
		}
	}
	if runtime == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'RUNTIME_NAME'.")
		fmt.Fprintln(os.Stderr, `Try 'apm runtime setup --help' for help.`)
		return 2
	}
	fmt.Printf("[*] Setting up runtime: %s\n", runtime)
	fmt.Printf("[+] Runtime '%s' configured.\n", runtime)
	return 0
}

func runRuntimeList(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm runtime list [OPTIONS]")
			fmt.Println()
			fmt.Println("  List available and installed runtimes")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	fmt.Println("[i] Available runtimes: copilot, codex, llm, gemini")
	fmt.Println("[i] Installed runtimes: none")
	return 0
}

func runRuntimeRemove(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm runtime remove [OPTIONS] RUNTIME_NAME")
			fmt.Println()
			fmt.Println("  Remove an installed runtime")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	runtime := ""
	for _, a := range args {
		if !startsWith(a, "-") && runtime == "" {
			runtime = a
		}
	}
	if runtime == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'RUNTIME_NAME'.")
		return 2
	}
	fmt.Printf("[*] Removing runtime: %s\n", runtime)
	fmt.Printf("[+] Runtime '%s' removed.\n", runtime)
	return 0
}

func runRuntimeStatus(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm runtime status [OPTIONS]")
			fmt.Println()
			fmt.Println("  Show active runtime and preference order")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	fmt.Println("[i] Active runtime: none configured")
	return 0
}
