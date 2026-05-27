// Package commands provides CLI command types and helpers for the APM Go rewrite.
package commands

// CommandContext carries shared state for all commands.
type CommandContext struct {
	// ConfigPath overrides the default config file path.
	ConfigPath string
	// Verbose enables verbose output.
	Verbose bool
	// Global applies operations globally rather than per-project.
	Global bool
	// DryRun prints what would be done without making changes.
	DryRun bool
}

// CommandResult is the result of executing a CLI command.
type CommandResult struct {
	// ExitCode is the process exit code (0 = success).
	ExitCode int
	// Output contains the command's stdout.
	Output string
	// Error contains any error message.
	Error string
}

// NewCommandContext returns a CommandContext with default values.
func NewCommandContext() *CommandContext {
	return &CommandContext{}
}

// IsSuccess returns true when ExitCode == 0.
func (r *CommandResult) IsSuccess() bool {
	return r.ExitCode == 0
}
