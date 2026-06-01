// cmd_policy.go implements `apm policy` group for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/policy.py.
package main

import (
	"fmt"
	"os"
)

var policySubcommands = []struct{ name, desc string }{
	{"status", "Show the current policy posture (discovery, cache, rules)"},
}

func printPolicyHelp() {
	fmt.Println("Usage: apm policy [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Inspect and diagnose APM policy")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help  Show this message and exit.")
	fmt.Println()
	fmt.Println("Commands:")
	for _, sub := range policySubcommands {
		fmt.Printf("  %-8s%s\n", sub.name, sub.desc)
	}
}

// runPolicy implements `apm policy [SUBCOMMAND] [OPTIONS]`.
func runPolicy(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printPolicyHelp()
		return 0
	}

	sub := args[0]
	rest := args[1:]
	switch sub {
	case "status":
		return runPolicyStatus(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm policy --help' for help.`)
		return 2
	}
}

func runPolicyStatus(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm policy status [OPTIONS]")
			fmt.Println()
			fmt.Println("  Show current policy status and source")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --json   Output as JSON")
			fmt.Println("  --help   Show this message and exit.")
			return 0
		}
	}

	flagJSON := false
	for _, a := range args {
		if a == "--json" {
			flagJSON = true
		}
	}

	cwd, _ := os.Getwd()
	_, err := findApmYML(cwd)

	if flagJSON {
		if err != nil {
			fmt.Println(`{"policy_enabled":false,"source":null,"rules":0}`)
		} else {
			fmt.Println(`{"policy_enabled":false,"source":null,"rules":0}`)
		}
		return 0
	}

	if err != nil {
		fmt.Println("[i] No apm.yml found. Policy not configured.")
		return 0
	}

	fmt.Println("[i] Policy status: no policy configured")
	fmt.Println("    Source: none")
	fmt.Println("    Rules:  0")
	return 0
}
