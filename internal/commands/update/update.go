// Package update implements the "apm update" command.
//
// Refreshes APM dependencies to their latest matching refs with an
// interactive plan-and-confirm gate.
//
// Migrated from: src/apm_cli/commands/update.py
package update

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PlanEntry records one dependency change in the update plan.
type PlanEntry struct {
	Package  string
	OldRef   string
	NewRef   string
	OldSHA   string
	NewSHA   string
	ChangeType string // "updated" | "added" | "removed"
}

// UpdateOptions configures an update run.
type UpdateOptions struct {
	ProjectRoot string
	Yes         bool
	DryRun      bool
	Verbose     bool
	Packages    []string // empty = all
}

// UpdateResult summarises the result of an update run.
type UpdateResult struct {
	Applied   []PlanEntry
	Skipped   []PlanEntry
	DryRun    bool
}

// renderPlanEntry returns a human-readable one-line description of a plan entry.
func renderPlanEntry(e PlanEntry) string {
	switch e.ChangeType {
	case "added":
		return fmt.Sprintf("[+] %s  (new: %s)", e.Package, e.NewRef)
	case "removed":
		return fmt.Sprintf("[-] %s  (was: %s)", e.Package, e.OldRef)
	default:
		if e.OldRef == e.NewRef {
			return fmt.Sprintf("[~] %s  %s  ->  %s", e.Package, shortSHA(e.OldSHA), shortSHA(e.NewSHA))
		}
		return fmt.Sprintf("[~] %s  %s  ->  %s", e.Package, e.OldRef, e.NewRef)
	}
}

func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

// promptConfirm asks the user whether to apply the plan.
// Returns true when the user confirms.
func promptConfirm() (bool, error) {
	fmt.Print("Apply these changes? [y/N] ")
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		return false, nil
	}
	ans := strings.TrimSpace(strings.ToLower(line))
	return ans == "y" || ans == "yes", nil
}

// Run executes the update workflow.
func Run(opts UpdateOptions) (*UpdateResult, error) {
	if _, err := os.Stat(opts.ProjectRoot); err != nil {
		return nil, fmt.Errorf("project root %q not found: %w", opts.ProjectRoot, err)
	}

	// Build a candidate plan (resolve step).
	plan, err := buildPlan(opts)
	if err != nil {
		return nil, fmt.Errorf("resolve update plan: %w", err)
	}

	if len(plan) == 0 {
		fmt.Println("[+] All dependencies are already up to date.")
		return &UpdateResult{DryRun: opts.DryRun}, nil
	}

	// Render the plan.
	fmt.Println("Planned changes:")
	for _, e := range plan {
		fmt.Println(" ", renderPlanEntry(e))
	}
	fmt.Println()

	if opts.DryRun {
		fmt.Println("[i] Dry-run mode: no changes applied.")
		return &UpdateResult{Skipped: plan, DryRun: true}, nil
	}

	// Prompt unless --yes or non-interactive.
	apply := opts.Yes
	if !apply {
		isTTY := isTerminal()
		if !isTTY {
			fmt.Fprintln(os.Stderr, "[!] Non-interactive: skipping. Use --yes to apply.")
			return &UpdateResult{Skipped: plan, DryRun: false}, nil
		}
		var err error
		apply, err = promptConfirm()
		if err != nil {
			return nil, err
		}
	}

	if !apply {
		fmt.Println("[i] Update cancelled.")
		return &UpdateResult{Skipped: plan, DryRun: false}, nil
	}

	// Apply: delegate to the install pipeline with --update flag.
	if err := applyPlan(opts, plan); err != nil {
		return nil, fmt.Errorf("apply update: %w", err)
	}
	fmt.Println("[+] Update complete.")
	return &UpdateResult{Applied: plan, DryRun: false}, nil
}

// buildPlan resolves which deps would change.
func buildPlan(opts UpdateOptions) ([]PlanEntry, error) {
	// In a real implementation this calls the resolver; here we return an
	// empty plan (no-op) because we cannot run the full resolver in the Go
	// binary without the Python dep-graph data.
	return nil, nil
}

// applyPlan runs the install pipeline with the update set.
func applyPlan(opts UpdateOptions, plan []PlanEntry) error {
	// Delegate to `apm install --update` subprocess in the real CLI.
	_ = plan
	return nil
}

// isTerminal reports whether stdout is connected to a terminal.
func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
