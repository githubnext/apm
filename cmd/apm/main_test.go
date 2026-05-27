package main

import (
	"testing"
)

// TestBuildSmoke verifies that the apm binary scaffolding compiles and links.
// This is the first parity test: the binary exists and builds successfully.
func TestBuildSmoke(t *testing.T) {
	// If this test runs, the package compiled -- that is the assertion.
}

// TestParityHelpIncludesCommands verifies the commands map contains all expected commands.
func TestParityHelpIncludesCommands(t *testing.T) {
	expected := []string{
		"audit", "cache", "compile", "config", "deps", "init", "install",
		"list", "marketplace", "mcp", "outdated", "pack", "plugin", "policy",
		"prune", "run", "runtime", "search", "targets", "uninstall", "unpack",
		"update", "view",
	}
	for _, cmd := range expected {
		if _, ok := commands[cmd]; !ok {
			t.Errorf("commands map missing command: %s", cmd)
		}
	}
}

// TestParityCommandsMapCompleteness verifies that all Python CLI commands are present.
func TestParityCommandsMapCompleteness(t *testing.T) {
	required := []string{
		"audit", "cache", "compile", "config", "deps", "experimental",
		"init", "install", "list", "marketplace", "mcp", "outdated",
		"pack", "plugin", "policy", "preview", "prune", "run", "runtime",
		"search", "self-update", "targets", "uninstall", "unpack", "update", "view",
	}
	for _, cmd := range required {
		if _, ok := commands[cmd]; !ok {
			t.Errorf("commands map missing: %s", cmd)
		}
	}
}

// TestParityVersionFlag verifies --version exits cleanly with version string.
func TestParityVersionFlag(t *testing.T) {
	// run(["--version"]) should return 0.
	code := run([]string{"--version"})
	if code != 0 {
		t.Fatalf("expected exit 0 for --version, got %d", code)
	}
}

// TestParityHelpFlag verifies --help exits cleanly.
func TestParityHelpFlag(t *testing.T) {
	code := run([]string{"--help"})
	if code != 0 {
		t.Fatalf("expected exit 0 for --help, got %d", code)
	}
}

// TestParityNoArgs verifies running with no args shows help (exit 0).
func TestParityNoArgs(t *testing.T) {
	code := run([]string{})
	if code != 0 {
		t.Fatalf("expected exit 0 for no args, got %d", code)
	}
}

// TestParityHelpSubcommand verifies "apm help" exits cleanly.
func TestParityHelpSubcommand(t *testing.T) {
	code := run([]string{"help"})
	if code != 0 {
		t.Fatalf("expected exit 0 for help subcommand, got %d", code)
	}
}

// TestParityUnknownCommandExitsNonZero verifies unknown commands return non-zero.
func TestParityUnknownCommandExitsNonZero(t *testing.T) {
	code := run([]string{"nonexistent-command-xyz"})
	if code == 0 {
		t.Fatal("expected non-zero exit for unknown command")
	}
}

// TestParityInfoAlias verifies "info" is an alias for "view".
func TestParityInfoAlias(t *testing.T) {
	if aliases["info"] != "view" {
		t.Fatalf("expected info -> view alias, got %q", aliases["info"])
	}
}

// TestParitySubcommandHelp verifies each subcommand accepts --help.
func TestParitySubcommandHelp(t *testing.T) {
	for cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			code := run([]string{cmd, "--help"})
			if code != 0 {
				t.Fatalf("apm %s --help returned %d, want 0", cmd, code)
			}
		})
	}
}

// TestParityVersionString verifies the version constant is set (not empty).
func TestParityVersionString(t *testing.T) {
	if version == "" {
		t.Fatal("version string is empty")
	}
}

// TestParityAllCommandsHaveDescriptions verifies each command has a non-empty description.
func TestParityAllCommandsHaveDescriptions(t *testing.T) {
	for cmd, desc := range commands {
		if desc == "" {
			t.Errorf("command %q has empty description", cmd)
		}
	}
}

