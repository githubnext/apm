package commands_test

import (
	"testing"

	"github.com/githubnext/apm/internal/commands"
)

// TestParityCommandContextFields verifies field parity with Python CommandContext.
func TestParityCommandContextFields(t *testing.T) {
	ctx := commands.NewCommandContext()
	if ctx.Verbose {
		t.Fatal("expected Verbose=false")
	}
	if ctx.Global {
		t.Fatal("expected Global=false")
	}
	if ctx.DryRun {
		t.Fatal("expected DryRun=false")
	}
}

// TestParityCommandContextConfigPath verifies ConfigPath field.
func TestParityCommandContextConfigPath(t *testing.T) {
	ctx := &commands.CommandContext{ConfigPath: "/tmp/apm.yml"}
	if ctx.ConfigPath != "/tmp/apm.yml" {
		t.Fatalf("unexpected ConfigPath: %s", ctx.ConfigPath)
	}
}

// TestParityCommandResultSuccess verifies IsSuccess.
func TestParityCommandResultSuccess(t *testing.T) {
	r := &commands.CommandResult{ExitCode: 0}
	if !r.IsSuccess() {
		t.Fatal("expected IsSuccess for exit code 0")
	}
}

// TestParityCommandResultFailure verifies failure detection.
func TestParityCommandResultFailure(t *testing.T) {
	r := &commands.CommandResult{ExitCode: 1, Error: "something failed"}
	if r.IsSuccess() {
		t.Fatal("expected !IsSuccess for exit code 1")
	}
}

// TestParityCommandResultOutput verifies Output field.
func TestParityCommandResultOutput(t *testing.T) {
	r := &commands.CommandResult{ExitCode: 0, Output: "[+] Done"}
	if r.Output != "[+] Done" {
		t.Fatalf("unexpected output: %s", r.Output)
	}
}

// TestParityNewCommandContextDefaults verifies zero values.
func TestParityNewCommandContextDefaults(t *testing.T) {
	ctx := commands.NewCommandContext()
	if ctx == nil {
		t.Fatal("NewCommandContext returned nil")
	}
	if ctx.ConfigPath != "" {
		t.Fatalf("expected empty ConfigPath, got %s", ctx.ConfigPath)
	}
}

// TestParityCommandContextVerboseFlag verifies setting verbose.
func TestParityCommandContextVerboseFlag(t *testing.T) {
	ctx := commands.NewCommandContext()
	ctx.Verbose = true
	if !ctx.Verbose {
		t.Fatal("expected Verbose=true after set")
	}
}

// TestParityCommandContextDryRunFlag verifies DryRun field.
func TestParityCommandContextDryRunFlag(t *testing.T) {
	ctx := &commands.CommandContext{DryRun: true}
	if !ctx.DryRun {
		t.Fatal("expected DryRun=true")
	}
}

// TestParityCommandContextGlobalFlag verifies Global field.
func TestParityCommandContextGlobalFlag(t *testing.T) {
	ctx := &commands.CommandContext{Global: true}
	if !ctx.Global {
		t.Fatal("expected Global=true")
	}
}

// TestParityCommandResultErrorField verifies Error field.
func TestParityCommandResultErrorField(t *testing.T) {
	r := &commands.CommandResult{ExitCode: 2, Error: "permission denied"}
	if r.Error != "permission denied" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
}

// TestParityCommandResultExitCode verifies exit code field.
func TestParityCommandResultExitCode(t *testing.T) {
	r := &commands.CommandResult{ExitCode: 42}
	if r.ExitCode != 42 {
		t.Fatalf("unexpected exit code: %d", r.ExitCode)
	}
}
