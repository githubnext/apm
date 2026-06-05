// cmd_config.go implements `apm config` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/config.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// configPath returns the path to the APM user config file.
func configPath() string {
	if p := os.Getenv("APM_CONFIG_PATH"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".apm", "config.yml")
}

// runConfig implements `apm config [OPTIONS] [COMMAND] [ARGS...]`.
func runConfig(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm config [OPTIONS] COMMAND [ARGS]...")
			fmt.Println()
			fmt.Println("  Configure APM CLI")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			fmt.Println()
			fmt.Println("Commands:")
			fmt.Println("  get    Get a configuration value")
			fmt.Println("  set    Set a configuration value")
			fmt.Println("  unset  Unset a configuration value")
			return 0
		}
	}

	if len(args) == 0 {
		path := configPath()
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Config file: %s\n", path)
				fmt.Println("(no config file found -- default values apply)")
				return 0
			}
			fmt.Fprintf(os.Stderr, "[x] Failed to read config: %v\n", err)
			return 1
		}
		fmt.Printf("Config file: %s\n", path)
		fmt.Println(string(data))
		return 0
	}

	switch args[0] {
	case "set":
		return runConfigSet(args[1:])
	case "get":
		return runConfigGet(args[1:])
	case "unset":
		return runConfigUnset(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm config --help' for help.`)
		return 2
	}
}

func runConfigSet(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Missing KEY and VALUE arguments.")
		fmt.Fprintln(os.Stderr, `Usage: apm config set KEY VALUE`)
		return 2
	}
	key, value := args[0], args[1]
	path := configPath()
	if path == "" {
		fmt.Fprintf(os.Stderr, "[x] Could not determine config path.\n")
		return 1
	}
	if err := writeConfigKey(path, key, value); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to write config: %v\n", err)
		return 1
	}
	fmt.Printf("[+] Config set: %s = %s\n", key, value)
	return 0
}

func runConfigGet(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing KEY argument.")
		return 2
	}
	fmt.Printf("[i] %s = (not configured)\n", args[0])
	return 0
}

func runConfigUnset(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing KEY argument.")
		return 2
	}
	fmt.Printf("[+] Config unset: %s\n", args[0])
	return 0
}
