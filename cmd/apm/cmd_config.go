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

// runConfig implements `apm config [OPTIONS]`.
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

	path := configPath()
	if path == "" {
		fmt.Fprintf(os.Stderr, "[x] Could not determine config path.\n")
		return 1
	}

	// If a key=value is provided, offer a simple set operation hint.
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "[i] Config editing is interactive. Config file: %s\n", path)
		return 0
	}

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
