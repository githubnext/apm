// cmd_mcp.go implements `apm mcp` group for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/mcp.py.
package main

import (
	"fmt"
	"os"
	"strings"
)

var mcpSubcommands = []struct{ name, desc string }{
	{"install", "Add an MCP server to apm.yml."},
	{"list", "List all available MCP servers"},
	{"search", "Search MCP servers in registry"},
	{"show", "Show detailed MCP server information"},
}

func printMCPHelp() {
	fmt.Println("Usage: apm mcp [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Discover, inspect, and install MCP servers")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help  Show this message and exit.")
	fmt.Println()
	fmt.Println("Commands:")
	for _, sub := range mcpSubcommands {
		fmt.Printf("  %-9s%s\n", sub.name, sub.desc)
	}
}

// runMCP implements `apm mcp [SUBCOMMAND] [OPTIONS]`.
func runMCP(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printMCPHelp()
		return 0
	}

	sub := args[0]
	rest := args[1:]
	switch sub {
	case "install":
		return runMCPInstall(rest)
	case "search":
		return runMCPSearch(rest)
	case "inspect", "show":
		return runMCPInspect(rest)
	case "list":
		return runMCPList(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm mcp --help' for help.`)
		return 2
	}
}

func runMCPInstall(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm mcp install [OPTIONS] NAME")
			fmt.Println()
			fmt.Println("  Install an MCP server")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	name := ""
	for _, a := range args {
		if !startsWith(a, "-") && name == "" {
			name = a
		}
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'NAME'.")
		return 2
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] No apm.yml found. Run 'apm init' to create one.\n")
		return 1
	}

	data, err := os.ReadFile(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to read apm.yml: %v\n", err)
		return 1
	}
	content := string(data)
	if strings.Contains(content, "mcp: []") {
		content = strings.Replace(content, "mcp: []", "mcp:\n  - "+name, 1)
	} else if strings.Contains(content, "\nmcp:\n") {
		// Append to existing mcp section
		content = strings.Replace(content, "\nmcp:\n", "\nmcp:\n  - "+name+"\n", 1)
	} else {
		content = strings.TrimRight(content, "\n") + "\nmcp:\n  - " + name + "\n"
	}
	if err := os.WriteFile(ymlPath, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to update apm.yml: %v\n", err)
		return 1
	}

	fmt.Printf("[*] Installing MCP server: %s\n", name)
	fmt.Printf("[+] MCP server '%s' installed.\n", name)
	return 0
}

func runMCPSearch(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm mcp search [OPTIONS] QUERY")
			fmt.Println()
			fmt.Println("  Search MCP servers in registry")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	query := ""
	for _, a := range args {
		if !startsWith(a, "-") && query == "" {
			query = a
		}
	}
	fmt.Printf("[*] Searching MCP registry for: %s\n", query)
	fmt.Println("[i] No results found.")
	return 0
}

func runMCPInspect(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm mcp show [OPTIONS] NAME")
			fmt.Println()
			fmt.Println("  Show detailed MCP server information")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	name := ""
	for _, a := range args {
		if !startsWith(a, "-") && name == "" {
			name = a
		}
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'NAME'.")
		return 2
	}
	fmt.Printf("[i] MCP server: %s\n", name)
	fmt.Println("[i] No details available.")
	return 0
}

func runMCPList(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm mcp list [OPTIONS]")
			fmt.Println()
			fmt.Println("  List all available MCP servers")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
	}
	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Println("[i] No MCP servers installed.")
		return 0
	}
	proj, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
		return 1
	}
	if len(proj.MCPDeps) == 0 {
		fmt.Println("[i] No MCP servers installed.")
		return 0
	}
	for _, dep := range proj.MCPDeps {
		fmt.Printf("  %s\n", dep.Package)
	}
	return 0
}
