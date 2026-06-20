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

	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm policy --help' for help.`)
		return 2
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
			fmt.Println("  --policy-source TEXT  Policy source")
			fmt.Println("  --no-cache  Force fresh policy fetch")
			fmt.Println("  -o, --output PATH  Write output to file")
			fmt.Println("  --check  Exit non-zero when policy is not satisfied")
			fmt.Println("  --help   Show this message and exit.")
			return 0
		}
	}

	flagJSON := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--json", "--no-cache", "--check":
			if a == "--json" {
				flagJSON = true
			}
		case "--policy-source", "-o", "--output":
			i++ // skip value
		default:
			if startsWith(a, "-") && !startsWith(a, "--policy-source=") && !startsWith(a, "--output=") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm policy status --help' for help.`)
				return 2
			}
		}
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		if flagJSON {
			fmt.Println(`{"policy_enabled":false,"source":null,"rules":0}`)
		} else {
			fmt.Println("[i] No apm.yml found. Policy not configured.")
		}
		return 0
	}

	proj, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
		return 1
	}

	if len(proj.PolicyDeny) == 0 {
		if flagJSON {
			fmt.Println(`{"policy_enabled":false,"source":null,"rules":0}`)
		} else {
			fmt.Println("[i] Policy status: no policy configured")
			fmt.Println("    Source: none")
			fmt.Println("    Rules:  0")
		}
		return 0
	}

	// Check deps against deny list.
	denySet := make(map[string]bool, len(proj.PolicyDeny))
	for _, d := range proj.PolicyDeny {
		denySet[d] = true
	}

	var violations []string
	for _, dep := range proj.Deps {
		if denySet[dep.Package] {
			violations = append(violations, dep.Package)
		}
	}
	for _, dep := range proj.MCPDeps {
		if denySet[dep.Package] {
			violations = append(violations, dep.Package)
		}
	}

	if len(violations) > 0 {
		if flagJSON {
			fmt.Printf(`{"policy_enabled":true,"violations":%d}`, len(violations))
			fmt.Println()
		} else {
			fmt.Printf("[x] Policy violation: %d denied package(s) found\n", len(violations))
			for _, v := range violations {
				fmt.Printf("    [x] Denied: %s\n", v)
			}
		}
		return 1
	}

	if flagJSON {
		fmt.Printf(`{"policy_enabled":true,"rules":%d,"violations":0}`, len(proj.PolicyDeny))
		fmt.Println()
	} else {
		fmt.Printf("[+] Policy OK: %d rules, no violations\n", len(proj.PolicyDeny))
	}
	return 0
}
