// Package finalize emits verbose install stats and returns the final result.
// Mirrors src/apm_cli/install/phases/finalize.py.
package finalize

import "fmt"

// InstallStats holds counters accumulated during the install pipeline.
type InstallStats struct {
	LinksResolved        int
	CommandsIntegrated   int
	HooksIntegrated      int
	InstructionsIntegrated int
	InstalledCount       int
	UnpinnedCount        int
	TotalPromptsIntegrated int
	TotalAgentsIntegrated  int
}

// InstallResult is the value returned from the finalize phase.
type InstallResult struct {
	InstalledCount         int
	TotalPromptsIntegrated int
	TotalAgentsIntegrated  int
	PackageTypes           map[string]int
	Warnings               []string
	Errors                 []string
}

// UnpinnedWarning formats the user-facing warning for unpinned dependencies.
// names is the (possibly empty) list of dep display names. count is total.
func UnpinnedWarning(count int, names []string) string {
	noun := "dependency"
	if count != 1 {
		noun = "dependencies"
	}
	if len(names) == 0 {
		return fmt.Sprintf("%d %s unpinned -- add #tag or #sha to prevent drift", count, noun)
	}
	shown := names
	if len(shown) > 5 {
		shown = shown[:5]
	}
	suffix := ""
	for i, n := range shown {
		if i > 0 {
			suffix += ", "
		}
		suffix += n
	}
	extra := len(names) - len(shown)
	if extra > 0 {
		suffix += fmt.Sprintf(", and %d more", extra)
	}
	return fmt.Sprintf("%d %s unpinned: %s -- add #tag or #sha to prevent drift", count, noun, suffix)
}

// VerboseStatLines returns human-readable lines describing non-zero counters.
func VerboseStatLines(stats InstallStats) []string {
	var lines []string
	if stats.LinksResolved > 0 {
		lines = append(lines, fmt.Sprintf("Resolved %d context file links", stats.LinksResolved))
	}
	if stats.CommandsIntegrated > 0 {
		lines = append(lines, fmt.Sprintf("Integrated %d command(s)", stats.CommandsIntegrated))
	}
	if stats.HooksIntegrated > 0 {
		lines = append(lines, fmt.Sprintf("Integrated %d hook(s)", stats.HooksIntegrated))
	}
	if stats.InstructionsIntegrated > 0 {
		lines = append(lines, fmt.Sprintf("Integrated %d instruction(s)", stats.InstructionsIntegrated))
	}
	return lines
}
