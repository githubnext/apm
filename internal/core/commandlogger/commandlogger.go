// Package commandlogger provides structured CLI output infrastructure for APM commands.
//
// Mirrors src/apm_cli/core/command_logger.py.
package commandlogger

import (
	"fmt"

	"github.com/githubnext/apm/internal/utils/console"
)

// StripSourcePrefix removes the "org:" or "url:" prefix from a policy source string.
func StripSourcePrefix(source string) string {
	if source == "" {
		return ""
	}
	for _, pfx := range []string{"org:", "url:"} {
		if len(source) > len(pfx) && source[:len(pfx)] == pfx {
			return source[len(pfx):]
		}
	}
	return source
}

// CommandLogger is the base context-aware logger for all CLI commands.
// All methods delegate to console helpers -- no new output primitives.
type CommandLogger struct {
	Command  string
	Verbose  bool
	DryRun   bool
}

// NewCommandLogger creates a new CommandLogger.
func NewCommandLogger(command string, verbose, dryRun bool) *CommandLogger {
	return &CommandLogger{Command: command, Verbose: verbose, DryRun: dryRun}
}

// Start logs the start of an operation.
func (l *CommandLogger) Start(message string) {
	console.Info(message, "running")
}

// Progress logs progress during an operation.
func (l *CommandLogger) Progress(message string) {
	console.Info(message, "info")
}

// MCPLookupHeartbeat emits a single batch heartbeat before MCP registry validation.
func (l *CommandLogger) MCPLookupHeartbeat(count int) {
	if count <= 0 {
		return
	}
	noun := "servers"
	if count == 1 {
		noun = "server"
	}
	console.Info(fmt.Sprintf("Looking up %d MCP %s in registry...", count, noun), "running")
}

// Info logs static advisory/informational context.
func (l *CommandLogger) Info(message, symbol string) {
	if symbol == "" {
		symbol = "info"
	}
	console.Info(message, symbol)
}

// Success logs successful completion.
func (l *CommandLogger) Success(message string) {
	console.Success(message, "sparkles")
}

// Warning logs a warning.
func (l *CommandLogger) Warning(message string) {
	console.Warning(message, "warning")
}

// Error logs an error.
func (l *CommandLogger) Error(message string) {
	console.Error(message, "error")
}

// VerboseDetail logs a detail only when verbose mode is enabled.
func (l *CommandLogger) VerboseDetail(message string) {
	if l.Verbose {
		console.Echo(nil, message, "dim", "", false)
	}
}

// TreeItem logs a tree sub-item (continuation line) under a package block.
func (l *CommandLogger) TreeItem(message string) {
	console.Echo(nil, message, "green", "", false)
}

// BlankLine logs a blank line.
func (l *CommandLogger) BlankLine() {
	console.Echo(nil, "", "", "", false)
}

// PackageInlineWarning logs an inline warning under a package block (verbose only).
func (l *CommandLogger) PackageInlineWarning(message string) {
	if l.Verbose {
		console.Echo(nil, message, "yellow", "", false)
	}
}

// DryRunNotice logs what would happen in dry-run mode.
func (l *CommandLogger) DryRunNotice(whatWouldHappen string) {
	console.Info(fmt.Sprintf("[dry-run] %s", whatWouldHappen), "info")
}

// ShouldExecute returns false if in dry-run mode.
func (l *CommandLogger) ShouldExecute() bool {
	return !l.DryRun
}

// AuthStep logs an auth resolution step (verbose only).
func (l *CommandLogger) AuthStep(step string, success bool, detail string) {
	if !l.Verbose {
		return
	}
	msg := fmt.Sprintf("  auth: %s", step)
	if detail != "" {
		msg += fmt.Sprintf(" (%s)", detail)
	}
	symbol := "check"
	if !success {
		symbol = "error"
	}
	console.Echo(nil, msg, "dim", symbol, false)
}

// PolicyDiscoveryMiss logs a policy-discovery non-success outcome.
func (l *CommandLogger) PolicyDiscoveryMiss(outcome, source, errText, hostOrg string) {
	if errText == "" {
		errText = "unknown"
	}
	switch outcome {
	case "absent":
		if !l.Verbose {
			return
		}
		org := hostOrg
		if org == "" {
			org = StripSourcePrefix(source)
		}
		if org == "" {
			org = "this project"
		}
		console.Info(fmt.Sprintf("No org policy found for %s", org), "info")

	case "no_git_remote":
		if !l.Verbose {
			return
		}
		console.Info("Could not determine org from git remote; policy auto-discovery skipped", "info")

	case "empty":
		src := source
		if src == "" {
			src = "this project"
		}
		console.Warning(fmt.Sprintf("Org policy at %s is present but empty; no enforcement applied", src), "warning")

	case "malformed":
		console.Warning(fmt.Sprintf("Policy at %s is malformed: %s. Contact your org admin to fix the policy file.", source, errText), "warning")

	case "cache_miss_fetch_fail":
		console.Warning(fmt.Sprintf("Could not fetch org policy from %s (%s); proceeding without policy enforcement. Retry, check connectivity, or use --no-policy to bypass.", source, errText), "warning")

	case "garbage_response":
		console.Warning(fmt.Sprintf("Policy response from %s is not valid YAML (%s); proceeding without policy enforcement. Contact your org admin or use --no-policy.", source, errText), "warning")

	case "cached_stale":
		console.Warning(fmt.Sprintf("Using stale cached policy (refresh failed: %s); enforcement still applies from cached policy.", errText), "warning")

	case "hash_mismatch":
		console.Error(fmt.Sprintf("Policy hash mismatch: pinned hash does not match fetched policy (%s). Update apm.yml policy.hash or contact your org admin.", errText), "error")

	default:
		if errText != "unknown" && errText != "" {
			console.Warning(fmt.Sprintf("Policy discovery issue: %s", errText), "warning")
		}
	}
}

// PolicyViolation records a policy violation for a dependency.
func (l *CommandLogger) PolicyViolation(depRef, reason, severity, source string) {
	// Strip depRef prefix if present.
	prefix := depRef + ": "
	if len(reason) > len(prefix) && reason[:len(prefix)] == prefix {
		reason = reason[len(prefix):]
	}
	if severity == "block" {
		console.Error(fmt.Sprintf("Policy violation: %s -- %s", depRef, reason), "error")
		if source != "" {
			msg := fmt.Sprintf("  Blocked by org policy at %s -- remove `%s` from apm.yml, contact admin to update policy, or use `--no-policy` for one-off bypass", source, depRef)
			console.Echo(nil, msg, "dim", "", false)
		}
	}
}

// PolicyDisabled logs a loud warning that policy enforcement is disabled.
func (l *CommandLogger) PolicyDisabled(reason string) {
	console.Warning(fmt.Sprintf("Policy enforcement disabled by %s for this invocation. This does NOT bypass apm audit --ci. CI will still fail the PR for the same policy violation.", reason), "warning")
}

// InstallSummary logs the final install summary.
func (l *CommandLogger) InstallSummary(apmCount, mcpCount, errors, staleCleaned int, elapsedSeconds float64, hasElapsed bool) {
	var parts []string
	if apmCount > 0 {
		noun := "dependencies"
		if apmCount == 1 {
			noun = "dependency"
		}
		parts = append(parts, fmt.Sprintf("%d APM %s", apmCount, noun))
	}
	if mcpCount > 0 {
		noun := "servers"
		if mcpCount == 1 {
			noun = "server"
		}
		parts = append(parts, fmt.Sprintf("%d MCP %s", mcpCount, noun))
	}

	cleanupSuffix := ""
	if staleCleaned > 0 {
		fNoun := "files"
		if staleCleaned == 1 {
			fNoun = "file"
		}
		cleanupSuffix = fmt.Sprintf(" (%d stale %s cleaned)", staleCleaned, fNoun)
	}

	timingSuffix := ""
	if hasElapsed {
		timingSuffix = fmt.Sprintf(" in %.1fs", elapsedSeconds)
	}

	if len(parts) > 0 {
		summary := joinParts(parts)
		if errors > 0 {
			console.Warning(fmt.Sprintf("Installed %s%s%s with %d error(s).", summary, cleanupSuffix, timingSuffix, errors), "warning")
		} else {
			console.Success(fmt.Sprintf("Installed %s%s%s.", summary, cleanupSuffix, timingSuffix), "sparkles")
		}
	} else if errors > 0 {
		console.Error(fmt.Sprintf("Installation failed with %d error(s)%s.", errors, timingSuffix), "error")
	}
}

func joinParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return parts[0] + " and " + parts[1]
}

// InstallInterrupted logs a minimal elapsed-time line for interrupted installs.
func (l *CommandLogger) InstallInterrupted(elapsedSeconds float64) {
	console.Warning(fmt.Sprintf("Install interrupted after %.1fs.", elapsedSeconds), "warning")
}

// InstallLogger is the install-specific logger with validation, resolution, and download phases.
type InstallLogger struct {
	*CommandLogger
	Partial         bool
	staleCleaned    int
}

// NewInstallLogger creates a new InstallLogger.
func NewInstallLogger(verbose, dryRun, partial bool) *InstallLogger {
	return &InstallLogger{
		CommandLogger: NewCommandLogger("install", verbose, dryRun),
		Partial:       partial,
	}
}

// ValidationStart logs start of package validation.
func (l *InstallLogger) ValidationStart(count int) {
	noun := "packages"
	if count == 1 {
		noun = "package"
	}
	console.Info(fmt.Sprintf("Validating %d %s...", count, noun), "gear")
}

// ValidationPass logs a package that passed validation.
func (l *InstallLogger) ValidationPass(canonical string, alreadyPresent bool) {
	if alreadyPresent {
		console.Echo(nil, fmt.Sprintf("%s (already in apm.yml)", canonical), "dim", "check", false)
	} else {
		console.Success(canonical, "check")
	}
}

// ValidationFail logs a package that failed validation.
func (l *InstallLogger) ValidationFail(pkg, reason string) {
	console.Error(fmt.Sprintf("%s -- %s", pkg, reason), "error")
}

// ResolutionStart logs start of dependency resolution.
func (l *InstallLogger) ResolutionStart(toInstallCount, lockfileCount int) {
	if l.Partial {
		noun := "packages"
		if toInstallCount == 1 {
			noun = "package"
		}
		console.Info(fmt.Sprintf("Installing %d new %s...", toInstallCount, noun), "running")
		if lockfileCount > 0 && l.Verbose {
			console.Echo(nil, fmt.Sprintf("  (%d existing dependencies in lockfile)", lockfileCount), "dim", "", false)
		}
	} else {
		console.Info("Installing dependencies from apm.yml...", "running")
		if lockfileCount > 0 {
			console.Info(fmt.Sprintf("Using apm.lock.yaml (%d locked dependencies)", lockfileCount), "")
		}
	}
}

// NothingToInstall logs when there's nothing to install.
func (l *InstallLogger) NothingToInstall(lockfilePresent, updateMode bool) {
	if l.Partial {
		console.Info("Requested packages are already installed.", "check")
	} else {
		console.Success("All dependencies are up to date.", "check")
	}
	if lockfilePresent && !updateMode {
		console.Info("Lockfile already satisfied -- run 'apm update' to resolve latest refs.", "")
	}
}

// DownloadStart logs start of a package download.
func (l *InstallLogger) DownloadStart(depName string, cached bool) {
	if cached {
		l.VerboseDetail(fmt.Sprintf("  Using cached: %s", depName))
	} else if l.Verbose {
		console.Info(fmt.Sprintf("  Downloading: %s", depName), "download")
	}
}

// ResolvingHeartbeat emits a per-dependency progress heartbeat during BFS resolve.
func (l *InstallLogger) ResolvingHeartbeat(depName string) {
	if l.Verbose {
		console.Info(fmt.Sprintf("  Resolving: %s", depName), "running")
	}
}

// DownloadComplete logs completion of a package download.
func (l *InstallLogger) DownloadComplete(depName string) {
	l.VerboseDetail(fmt.Sprintf("  Downloaded: %s", depName))
}
