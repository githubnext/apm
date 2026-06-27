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

// validConfigKeys is the set of user-settable config keys.
var validConfigKeys = map[string]bool{
	"auto-integrate": true,
	"temp-dir":       true,
}

// runConfig implements `apm config [OPTIONS] [COMMAND] [ARGS...]`.
func runConfig(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
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

	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm config --help' for help.`)
		return 2
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
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Println("Usage: apm config set [OPTIONS] KEY VALUE")
		fmt.Println()
		fmt.Println("  Set a configuration value")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --help  Show this message and exit.")
		return 0
	}
	for _, a := range args {
		if startsWith(a, "-") && a != "--help" && a != "-h" {
			fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
			fmt.Fprintln(os.Stderr, `Try 'apm config set --help' for help.`)
			return 2
		}
	}
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Missing KEY and VALUE arguments.")
		fmt.Fprintln(os.Stderr, `Usage: apm config set KEY VALUE`)
		return 2
	}
	key, value := args[0], args[1]
	if !validConfigKeys[key] {
		fmt.Fprintf(os.Stderr, "[x] Unknown configuration key: '%s'\n", key)
		fmt.Fprintf(os.Stderr, "[>] Valid keys: auto-integrate, temp-dir\n")
		return 1
	}
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
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Println("Usage: apm config get [OPTIONS] [KEY]")
		fmt.Println()
		fmt.Println("  Get a configuration value")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --help  Show this message and exit.")
		return 0
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing KEY argument.")
		return 2
	}
	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm config get --help' for help.`)
		return 2
	}
	key := args[0]
	if !validConfigKeys[key] {
		fmt.Fprintf(os.Stderr, "[x] Unknown configuration key: '%s'\n", key)
		fmt.Fprintf(os.Stderr, "[>] Valid keys: auto-integrate, temp-dir\n")
		return 1
	}
	path := configPath()
	if val, found := readConfigKey(path, key); found {
		fmt.Printf("%s: %s\n", key, val)
		return 0
	}
	// Return default when key is not set in config file.
	switch key {
	case "auto-integrate":
		fmt.Printf("auto-integrate: true\n")
	case "temp-dir":
		fmt.Printf("temp-dir: Not set (using system default)\n")
	}
	return 0
}

func runConfigUnset(args []string) int {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Println("Usage: apm config unset [OPTIONS] KEY")
		fmt.Println()
		fmt.Println("  Unset a configuration value")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --help  Show this message and exit.")
		return 0
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing KEY argument.")
		return 2
	}
	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm config unset --help' for help.`)
		return 2
	}
	key := args[0]
	if !validConfigKeys[key] {
		fmt.Fprintf(os.Stderr, "[x] Unknown configuration key: '%s'\n", key)
		fmt.Fprintf(os.Stderr, "[>] Valid keys: auto-integrate, temp-dir\n")
		return 1
	}
	path := configPath()
	if path == "" {
		fmt.Fprintf(os.Stderr, "[x] Could not determine config path.\n")
		return 1
	}
	if err := removeConfigKey(path, key); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to update config: %v\n", err)
		return 1
	}
	fmt.Printf("[+] Config unset: %s\n", key)
	return 0
}
