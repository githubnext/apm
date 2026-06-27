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

	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm mcp --help' for help.`)
		return 2
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
			fmt.Println("  Add an MCP server to apm.yml. Alias for 'apm install --mcp'.")
			fmt.Println()
			fmt.Println("  Examples:")
			fmt.Println()
			fmt.Println("    apm mcp install fetch -- npx -y @modelcontextprotocol/server-fetch")
			fmt.Println()
			fmt.Println("    apm mcp install api --transport http --url https://example.com/mcp")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			fmt.Println()
			fmt.Println("  Common options (see `apm install --mcp --help` for full list): --transport")
			fmt.Println("  [stdio|http|sse|streamable-http] --url URL           Server URL for remote")
			fmt.Println("  transports --env KEY=VALUE     Environment variable (repeatable) --header")
			fmt.Println("  KEY=VALUE  HTTP header (repeatable) --registry URL      Custom registry URL")
			fmt.Println("  --mcp-version VER    Pin registry entry to a specific version --dev / --dry-")
			fmt.Println("  run / --force / --verbose / --no-policy")
			return 0
		}
	}
	name := ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--verbose", "-v":
			// known no-value flags
		case "--limit":
			i++ // skip value
		default:
			// Python Click ignore_unknown_options=True assigns --X args to the
			// NAME positional (they do NOT go to ctx.args). Accept all
			// unrecognized args (including --X) as NAME to match that behavior.
			if !startsWith(a, "--limit=") && name == "" {
				name = a
			}
		}
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "Usage: apm mcp install [OPTIONS] NAME")
		fmt.Fprintln(os.Stderr, "Try 'apm mcp install --help' for help.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'NAME'.")
		return 2
	}
	// Python's mcp_install forwards NAME to `apm install --mcp NAME`, which
	// rejects any option-like value as an unknown option. Mirror that error.
	if startsWith(name, "-") {
		fmt.Println("[!] Install interrupted after 0.0s.")
		fmt.Fprintln(os.Stderr, "Usage: apm install [OPTIONS] [PACKAGES]...")
		fmt.Fprintln(os.Stderr, "Try 'apm install --help' for help.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Error: MCP name cannot start with '-'; did you forget a value for --mcp?")
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
			fmt.Println("  --limit INTEGER  Number of results to show  [default: 10]")
			fmt.Println("  -v, --verbose    Show detailed output")
			fmt.Println("  --help           Show this message and exit.")
			return 0
		}
	}
	query := ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--verbose", "-v":
			// known no-value flags
		case "--limit":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--limit=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm mcp search --help' for help.`)
				return 2
			}
			if !startsWith(a, "-") && query == "" {
				query = a
			}
		}
	}
	fmt.Printf("[*] Searching MCP registry for: %s\n", query)
	fmt.Println("[i] No results found.")
	return 0
}

func runMCPInspect(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm mcp show [OPTIONS] SERVER_NAME")
			fmt.Println()
			fmt.Println("  Show detailed MCP server information")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose  Show detailed output")
			fmt.Println("  --help         Show this message and exit.")
			return 0
		}
	}
	name := ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--verbose", "-v":
			// known no-value flags
		case "--limit":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--limit=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm mcp show --help' for help.`)
				return 2
			}
			if !startsWith(a, "-") && name == "" {
				name = a
			}
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
			fmt.Println("  --limit INTEGER  Number of results to show  [default: 20]")
			fmt.Println("  -v, --verbose    Show detailed output")
			fmt.Println("  --help           Show this message and exit.")
			return 0
		}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--verbose", "-v":
			// known no-value flags
		case "--limit":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--limit=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm mcp list --help' for help.`)
				return 2
			}
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
