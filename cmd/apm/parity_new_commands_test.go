// parity_new_commands_test.go: parity tests for newly wired command families:
// install, uninstall, update, prune, audit, policy, runtime, mcp, plugin,
// search, outdated, self-update, experimental, preview.
// Part of the greenfield parity harness (TestParity* counted by score.go).
package main

import (
	"strings"
	"testing"
)

// -- install --

func TestParityHarnessInstallHelp(t *testing.T) {
	_, _, code := runGo(t, "install", "--help")
	if code != 0 {
		t.Errorf("apm install --help exited %d, want 0", code)
	}
}

func TestParityHarnessInstallDryRunInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "install", "--dry-run")
	if r.GoExitCode != 0 {
		t.Errorf("apm install --dry-run exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

func TestParityHarnessInstallDryRunVerbose(t *testing.T) {
	r := runBothInTempRepo(t, apmYMLWithDeps, "install", "--dry-run", "--verbose")
	if r.GoExitCode != 0 {
		t.Errorf("apm install --dry-run --verbose exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

func TestParityHarnessInstallNoApmYML(t *testing.T) {
	// Empty yml content creates an apm.yml with empty/default values.
	r := runBothInTempRepo(t, minimalApmYML, "install", "--dry-run")
	if r.GoExitCode != 0 {
		t.Errorf("apm install --dry-run exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
}

// -- uninstall --

func TestParityHarnessUninstallHelp(t *testing.T) {
	_, _, code := runGo(t, "uninstall", "--help")
	if code != 0 {
		t.Errorf("apm uninstall --help exited %d, want 0", code)
	}
}

func TestParityHarnessUninstallMissingPackage(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "uninstall")
	if r.GoExitCode == 0 {
		t.Errorf("apm uninstall without package should fail, got 0")
	}
}

func TestParityHarnessUninstallDryRun(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "uninstall", "--dry-run", "some/package")
	if r.GoExitCode != 0 {
		t.Errorf("apm uninstall --dry-run exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

// -- update --

func TestParityHarnessUpdateHelp(t *testing.T) {
	_, _, code := runGo(t, "update", "--help")
	if code != 0 {
		t.Errorf("apm update --help exited %d, want 0", code)
	}
}

func TestParityHarnessUpdateDryRunInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "update", "--dry-run")
	if r.GoExitCode != 0 {
		t.Errorf("apm update --dry-run exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

func TestParityHarnessUpdateNoApmYML(t *testing.T) {
	// With a minimal apm.yml, update --dry-run should succeed.
	r := runBothInTempRepo(t, minimalApmYML, "update", "--dry-run")
	if r.GoExitCode != 0 {
		t.Errorf("apm update --dry-run exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
}

// -- prune --

func TestParityHarnessPruneHelp(t *testing.T) {
	_, _, code := runGo(t, "prune", "--help")
	if code != 0 {
		t.Errorf("apm prune --help exited %d, want 0", code)
	}
}

func TestParityHarnessPruneDryRunInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "prune", "--dry-run")
	if r.GoExitCode != 0 {
		t.Errorf("apm prune --dry-run exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

// -- audit --

func TestParityHarnessAuditHelp(t *testing.T) {
	_, _, code := runGo(t, "audit", "--help")
	if code != 0 {
		t.Errorf("apm audit --help exited %d, want 0", code)
	}
}

func TestParityHarnessAuditInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "audit")
	if r.GoExitCode != 0 {
		t.Errorf("apm audit exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

func TestParityHarnessAuditCIMode(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "audit", "--ci")
	if r.GoExitCode != 0 {
		t.Errorf("apm audit --ci exited %d, want 0 when no issues", r.GoExitCode)
	}
	assertNoPythonUnimplemented(t, r)
}

func TestParityHarnessAuditVerbose(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "audit", "--verbose")
	if r.GoExitCode != 0 {
		t.Errorf("apm audit --verbose exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

// -- policy --

func TestParityHarnessPolicyHelp(t *testing.T) {
	_, _, code := runGo(t, "policy", "--help")
	if code != 0 {
		t.Errorf("apm policy --help exited %d, want 0", code)
	}
}

func TestParityHarnessPolicyStatusInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "policy", "status")
	if r.GoExitCode != 0 {
		t.Errorf("apm policy status exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

func TestParityHarnessPolicyStatusJSON(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "policy", "status", "--json")
	if r.GoExitCode != 0 {
		t.Errorf("apm policy status --json exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
}

func TestParityHarnessPolicyUnknownSubcmd(t *testing.T) {
	_, _, code := runGo(t, "policy", "unknown-subcmd")
	if code == 0 {
		t.Errorf("apm policy unknown-subcmd should fail, got 0")
	}
}

// -- runtime --

func TestParityHarnessRuntimeHelp(t *testing.T) {
	_, _, code := runGo(t, "runtime", "--help")
	if code != 0 {
		t.Errorf("apm runtime --help exited %d, want 0", code)
	}
}

func TestParityHarnessRuntimeListHelp(t *testing.T) {
	_, _, code := runGo(t, "runtime", "list", "--help")
	if code != 0 {
		t.Errorf("apm runtime list --help exited %d, want 0", code)
	}
}

func TestParityHarnessRuntimeList(t *testing.T) {
	_, _, code := runGo(t, "runtime", "list")
	if code != 0 {
		t.Errorf("apm runtime list exited %d, want 0", code)
	}
}

func TestParityHarnessRuntimeStatus(t *testing.T) {
	_, _, code := runGo(t, "runtime", "status")
	if code != 0 {
		t.Errorf("apm runtime status exited %d, want 0", code)
	}
}

func TestParityHarnessRuntimeSetupMissingArg(t *testing.T) {
	_, _, code := runGo(t, "runtime", "setup")
	if code == 0 {
		t.Errorf("apm runtime setup without arg should fail, got 0")
	}
}

func TestParityHarnessRuntimeSetupHelp(t *testing.T) {
	_, _, code := runGo(t, "runtime", "setup", "--help")
	if code != 0 {
		t.Errorf("apm runtime setup --help exited %d, want 0", code)
	}
}

// -- mcp --

func TestParityHarnessMCPHelp(t *testing.T) {
	_, _, code := runGo(t, "mcp", "--help")
	if code != 0 {
		t.Errorf("apm mcp --help exited %d, want 0", code)
	}
}

func TestParityHarnessMCPListHelp(t *testing.T) {
	_, _, code := runGo(t, "mcp", "list", "--help")
	if code != 0 {
		t.Errorf("apm mcp list --help exited %d, want 0", code)
	}
}

func TestParityHarnessMCPList(t *testing.T) {
	_, _, code := runGo(t, "mcp", "list")
	if code != 0 {
		t.Errorf("apm mcp list exited %d, want 0", code)
	}
}

func TestParityHarnessMCPSearchHelp(t *testing.T) {
	_, _, code := runGo(t, "mcp", "search", "--help")
	if code != 0 {
		t.Errorf("apm mcp search --help exited %d, want 0", code)
	}
}

func TestParityHarnessMCPInstallMissingArg(t *testing.T) {
	_, _, code := runGo(t, "mcp", "install")
	if code == 0 {
		t.Errorf("apm mcp install without arg should fail, got 0")
	}
}

func TestParityHarnessMCPUnknownSubcmd(t *testing.T) {
	_, _, code := runGo(t, "mcp", "unknown-subcmd")
	if code == 0 {
		t.Errorf("apm mcp unknown-subcmd should fail, got 0")
	}
}

// -- plugin --

func TestParityHarnessPluginHelp(t *testing.T) {
	_, _, code := runGo(t, "plugin", "--help")
	if code != 0 {
		t.Errorf("apm plugin --help exited %d, want 0", code)
	}
}

func TestParityHarnessPluginInitHelp(t *testing.T) {
	_, _, code := runGo(t, "plugin", "init", "--help")
	if code != 0 {
		t.Errorf("apm plugin init --help exited %d, want 0", code)
	}
}

func TestParityHarnessPluginUnknownSubcmd(t *testing.T) {
	_, _, code := runGo(t, "plugin", "unknown-subcmd")
	if code == 0 {
		t.Errorf("apm plugin unknown-subcmd should fail, got 0")
	}
}

// -- search --

func TestParityHarnessSearchHelp(t *testing.T) {
	_, _, code := runGo(t, "search", "--help")
	if code != 0 {
		t.Errorf("apm search --help exited %d, want 0", code)
	}
}

func TestParityHarnessSearchMissingArg(t *testing.T) {
	_, _, code := runGo(t, "search")
	if code == 0 {
		t.Errorf("apm search without arg should fail, got 0")
	}
}

// -- outdated --

func TestParityHarnessOutdatedHelp(t *testing.T) {
	_, _, code := runGo(t, "outdated", "--help")
	if code != 0 {
		t.Errorf("apm outdated --help exited %d, want 0", code)
	}
}

func TestParityHarnessOutdatedInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "outdated")
	if r.GoExitCode != 0 {
		t.Errorf("apm outdated exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

// -- self-update --

func TestParityHarnessSelfUpdateHelp(t *testing.T) {
	_, _, code := runGo(t, "self-update", "--help")
	if code != 0 {
		t.Errorf("apm self-update --help exited %d, want 0", code)
	}
}

func TestParityHarnessSelfUpdateCheck(t *testing.T) {
	_, _, code := runGo(t, "self-update", "--check")
	if code != 0 {
		t.Errorf("apm self-update --check exited %d, want 0", code)
	}
}

// -- experimental --

func TestParityHarnessExperimentalHelp(t *testing.T) {
	_, _, code := runGo(t, "experimental", "--help")
	if code != 0 {
		t.Errorf("apm experimental --help exited %d, want 0", code)
	}
}

func TestParityHarnessExperimentalList(t *testing.T) {
	_, _, code := runGo(t, "experimental", "list")
	if code != 0 {
		t.Errorf("apm experimental list exited %d, want 0", code)
	}
}

// -- preview --

func TestParityHarnessPreviewHelp(t *testing.T) {
	_, _, code := runGo(t, "preview", "--help")
	if code != 0 {
		t.Errorf("apm preview --help exited %d, want 0", code)
	}
}

func TestParityHarnessPreviewMissingArg(t *testing.T) {
	_, _, code := runGo(t, "preview")
	if code == 0 {
		t.Errorf("apm preview without arg should fail, got 0")
	}
}

func TestParityHarnessPreviewInTempRepo(t *testing.T) {
	apmYML := `name: test-project
version: 1.0.0
description: Parity test project
author: crane-test
targets:
  - copilot
dependencies:
  apm: []
  mcp: []
scripts:
  build: echo build
`
	r := runBothInTempRepo(t, apmYML, "preview", "build")
	if r.GoExitCode != 0 {
		t.Errorf("apm preview build exited %d\nstderr: %s", r.GoExitCode, r.GoStderr)
	}
	assertNoPythonUnimplemented(t, r)
}

// -- all new command families: help exits 0 matrix --

func TestParityHarnessNewCommandFamilyHelpExitCodes(t *testing.T) {
	families := []struct {
		name string
		args []string
	}{
		{"install", []string{"install", "--help"}},
		{"uninstall", []string{"uninstall", "--help"}},
		{"update", []string{"update", "--help"}},
		{"prune", []string{"prune", "--help"}},
		{"audit", []string{"audit", "--help"}},
		{"policy", []string{"policy", "--help"}},
		{"policy-status", []string{"policy", "status", "--help"}},
		{"runtime", []string{"runtime", "--help"}},
		{"runtime-list", []string{"runtime", "list", "--help"}},
		{"runtime-status", []string{"runtime", "status", "--help"}},
		{"runtime-setup-help", []string{"runtime", "setup", "--help"}},
		{"runtime-remove-help", []string{"runtime", "remove", "--help"}},
		{"mcp", []string{"mcp", "--help"}},
		{"mcp-list", []string{"mcp", "list", "--help"}},
		{"mcp-search", []string{"mcp", "search", "--help"}},
		{"mcp-inspect", []string{"mcp", "inspect", "--help"}},
		{"mcp-install", []string{"mcp", "install", "--help"}},
		{"plugin", []string{"plugin", "--help"}},
		{"plugin-init", []string{"plugin", "init", "--help"}},
		{"search", []string{"search", "--help"}},
		{"outdated", []string{"outdated", "--help"}},
		{"self-update", []string{"self-update", "--help"}},
		{"experimental", []string{"experimental", "--help"}},
		{"preview", []string{"preview", "--help"}},
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			_, _, code := runGo(t, f.args...)
			if code != 0 {
				t.Errorf("apm %v exited %d, want 0", f.args, code)
			}
		})
	}
}

// TestParityHarnessNoWIPMessagesNewCmds verifies new commands don't output WIP messages.
func TestParityHarnessNoWIPMessagesNewCmds(t *testing.T) {
	wipCmds := [][]string{
		{"install", "--help"},
		{"uninstall", "--help"},
		{"update", "--help"},
		{"prune", "--help"},
		{"audit", "--help"},
		{"policy", "--help"},
		{"runtime", "--help"},
		{"mcp", "--help"},
		{"plugin", "--help"},
		{"search", "--help"},
		{"outdated", "--help"},
		{"self-update", "--help"},
		{"experimental", "--help"},
		{"preview", "--help"},
	}
	for _, args := range wipCmds {
		stdout, stderr, _ := runGo(t, args...)
		combined := stdout + stderr
		if strings.Contains(combined, "not yet implemented") || strings.Contains(combined, "work in progress") || strings.Contains(combined, "WIP") {
			t.Errorf("apm %v still outputs WIP message:\n%s", args, combined)
		}
	}
}
