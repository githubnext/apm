// Package core provides target resolution, auth, scope, and orchestration
// primitives for the APM CLI.
package core

import "fmt"

// TargetResolutionError is the base error for all target-resolution failures.
type TargetResolutionError struct {
	msg string
}

func (e *TargetResolutionError) Error() string { return e.msg }

// NoHarnessError is returned when no harness signal is detected and no
// explicit target is set.
type NoHarnessError struct{ TargetResolutionError }

// AmbiguousHarnessError is returned when multiple distinct harness signals
// are detected.
type AmbiguousHarnessError struct{ TargetResolutionError }

// UnknownTargetError is returned when a target token is not in the canonical
// set.
type UnknownTargetError struct{ TargetResolutionError }

// ConflictingTargetsError is returned when apm.yml contains both 'target:'
// and 'targets:' (mutex).
type ConflictingTargetsError struct{ TargetResolutionError }

// EmptyTargetsListError is returned when apm.yml 'targets:' is present but
// empty.
type EmptyTargetsListError struct{ TargetResolutionError }

// signal list used in error messages (mirrors _SIGNAL_LIST in errors.py)
const signalList = ".claude/, CLAUDE.md, .cursor/, .cursorrules, " +
	".github/copilot-instructions.md, .codex/, .gemini/, GEMINI.md, " +
	".opencode/, .windsurf/"

// RenderNoHarnessError returns the three-section error string for "no signal
// detected".
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

// RenderAmbiguousError returns the three-section error string for "multiple
// harnesses detected".
func RenderAmbiguousError(detected []string) string {
	if len(detected) == 0 {
		return "[x] Multiple harnesses detected"
	}
	detectedCSV := joinStrings(detected, ", ")
	first := detected[0]
	return fmt.Sprintf("[x] Multiple harnesses detected: %s\n", detectedCSV) +
		"\n" +
		fmt.Sprintf("APM found signals for %s but cannot decide which\n", detectedCSV) +
		"to deploy to. Pin your target explicitly.\n" +
		"\n" +
		"Fix with one of:\n" +
		"\n" +
		fmt.Sprintf("  apm install <pkg> --target %s\n", first) +
		"  apm install <pkg> --dry-run            # preview what each target does\n" +
		"  apm targets                            # see all detected harnesses\n" +
		"\n" +
		"Or declare in apm.yml:\n" +
		"\n" +
		"  targets:\n" +
		fmt.Sprintf("    - %s", first)
}

// RenderUnknownTargetError returns the three-section error string for an
// unknown target token.
func RenderUnknownTargetError(value string, valid []string) string {
	// hide agent-skills from user-facing list
	var visible []string
	for _, t := range valid {
		if t != "agent-skills" {
			visible = append(visible, t)
		}
	}
	sortStrings(visible)
	suggestion := "copilot"
	for _, t := range visible {
		if t == "copilot" {
			suggestion = "copilot"
			break
		}
		suggestion = t
	}
	validCSV := joinStrings(visible, ", ")
	if validCSV == "" {
		validCSV = suggestion
	}
	displayValue := stripBracketNoise(value)
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

// RenderConflictingSchemaError returns the error string for target/targets
// mutex.
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

// NewNoHarnessError constructs a NoHarnessError.
func NewNoHarnessError() *NoHarnessError {
	return &NoHarnessError{TargetResolutionError{RenderNoHarnessError()}}
}

// NewAmbiguousHarnessError constructs an AmbiguousHarnessError.
func NewAmbiguousHarnessError(detected []string) *AmbiguousHarnessError {
	return &AmbiguousHarnessError{TargetResolutionError{RenderAmbiguousError(detected)}}
}

// NewUnknownTargetError constructs an UnknownTargetError.
func NewUnknownTargetError(value string, valid []string) *UnknownTargetError {
	return &UnknownTargetError{TargetResolutionError{RenderUnknownTargetError(value, valid)}}
}

// NewConflictingTargetsError constructs a ConflictingTargetsError.
func NewConflictingTargetsError() *ConflictingTargetsError {
	return &ConflictingTargetsError{TargetResolutionError{RenderConflictingSchemaError()}}
}

// NewEmptyTargetsListError constructs an EmptyTargetsListError.
func NewEmptyTargetsListError() *EmptyTargetsListError {
	msg := "[x] 'targets:' in apm.yml is empty\n" +
		"\n" +
		"The targets list must contain at least one target.\n" +
		"\n" +
		"Fix with one of:\n" +
		"\n" +
		"  apm targets                            # see all supported harnesses\n" +
		"  apm install <pkg> --target claude\n" +
		"  apm init\n" +
		"\n" +
		"Or update apm.yml:\n" +
		"\n" +
		"  targets:\n" +
		"    - claude"
	return &EmptyTargetsListError{TargetResolutionError{msg}}
}
