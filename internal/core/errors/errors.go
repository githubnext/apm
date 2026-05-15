// Package errors provides the error hierarchy and renderers for target resolution.
//
// Mirrors src/apm_cli/core/errors.py.
package errors

import (
	"fmt"
	"sort"
	"strings"
)

// TargetResolutionError is the base error for all target-resolution user errors.
type TargetResolutionError struct {
	Message string
}

func (e *TargetResolutionError) Error() string {
	return e.Message
}

// NoHarnessError is returned when no harness signal is detected and no explicit target is set.
type NoHarnessError struct {
	TargetResolutionError
}

// AmbiguousHarnessError is returned when multiple distinct harness signals are detected.
type AmbiguousHarnessError struct {
	TargetResolutionError
}

// UnknownTargetError is returned when a target token is not in the canonical set.
type UnknownTargetError struct {
	TargetResolutionError
}

// ConflictingTargetsError is returned when apm.yml has both 'target:' and 'targets:'.
type ConflictingTargetsError struct {
	TargetResolutionError
}

// EmptyTargetsListError is returned when apm.yml 'targets:' is present but empty.
type EmptyTargetsListError struct {
	TargetResolutionError
}

const signalList = ".claude/, CLAUDE.md, .cursor/, .cursorrules, " +
	".github/copilot-instructions.md, .codex/, .gemini/, GEMINI.md, " +
	".opencode/, .windsurf/"

// RenderNoHarnessError returns the 3-section error for 'no signal detected'.
func RenderNoHarnessError() string {
	return "[x] No harness detected\n" +
		"\n" +
		"APM scanned for harness markers (" + signalList + ")" +
		" but found none in this project.\n" +
		"\n" +
		"Previously APM defaulted to copilot; this is now explicit.\n" +
		"\n" +
		"Fix with one of:\n" +
		"\n" +
		"  apm targets                            # see all supported harnesses\n" +
		"  apm install <pkg> --target claude      # deploy to a specific harness\n" +
		"  apm install <pkg> --target copilot     # or any supported target\n" +
		"\n" +
		"Or declare in apm.yml:\n" +
		"\n" +
		"  targets:\n" +
		"    - claude"
}

// RenderAmbiguousError returns the 3-section error for 'multiple harnesses detected'.
func RenderAmbiguousError(detected []string) string {
	detectedCSV := strings.Join(detected, ", ")
	suggestion := "claude"
	if len(detected) > 0 {
		suggestion = detected[0]
	}
	return fmt.Sprintf("[x] Multiple harnesses detected: %s\n", detectedCSV) +
		"\n" +
		fmt.Sprintf("APM found signals for %s but cannot decide which\n", detectedCSV) +
		"to deploy to. Pin your target explicitly.\n" +
		"\n" +
		"Fix with one of:\n" +
		"\n" +
		fmt.Sprintf("  apm install <pkg> --target %s\n", suggestion) +
		"  apm install <pkg> --dry-run            # preview what each target does\n" +
		"  apm targets                            # see all detected harnesses\n" +
		"\n" +
		"Or declare in apm.yml:\n" +
		"\n" +
		"  targets:\n" +
		fmt.Sprintf("    - %s", suggestion)
}

// RenderUnknownTargetError returns the 3-section error for an unknown target token.
func RenderUnknownTargetError(value string, valid []string) string {
	visible := make([]string, 0, len(valid))
	for _, t := range valid {
		if t != "agent-skills" {
			visible = append(visible, t)
		}
	}
	sort.Strings(visible)

	suggestion := "copilot"
	for _, t := range visible {
		if t == "copilot" {
			suggestion = "copilot"
			break
		}
	}
	if suggestion != "copilot" && len(visible) > 0 {
		suggestion = visible[0]
	}

	validCSV := strings.Join(visible, ", ")
	if validCSV == "" {
		validCSV = suggestion
	}

	// Strip bracket/quote noise
	displayValue := strings.Trim(value, "[]'\" ")
	if displayValue == "" {
		displayValue = value
	}
	if displayValue == "" {
		displayValue = "<empty>"
	}

	return fmt.Sprintf("[x] Unknown target '%s'\n", displayValue) +
		"\n" +
		fmt.Sprintf("Valid targets: %s\n", validCSV) +
		"\n" +
		"Fix with one of:\n" +
		"\n" +
		"  apm targets                            # see all supported harnesses\n" +
		fmt.Sprintf("  apm install <pkg> --target %s\n", suggestion) +
		"  apm install <pkg> --dry-run\n" +
		"\n" +
		"Or declare in apm.yml:\n" +
		"\n" +
		"  targets:\n" +
		fmt.Sprintf("    - %s", suggestion)
}

// RenderConflictingSchemaError returns the 3-section error for target/targets mutex.
func RenderConflictingSchemaError() string {
	return "[x] Cannot use both 'target:' and 'targets:' in apm.yml\n" +
		"\n" +
		"Use the canonical plural form:\n" +
		"\n" +
		"Fix with one of:\n" +
		"\n" +
		"  apm targets                            # see all supported harnesses\n" +
		"  apm install <pkg> --target claude\n" +
		"  apm init                               # regenerate apm.yml\n" +
		"\n" +
		"Or update apm.yml to use the canonical form:\n" +
		"\n" +
		"  targets:\n" +
		"    - claude\n" +
		"    - copilot"
}
