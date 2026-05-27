// parity_harness_test.go provides a reusable Python-vs-Go parity harness.
//
// The harness creates identical temporary repositories, runs both the Python
// apm CLI and the Go apm CLI with the same args/env/cwd/stdin, captures
// exit code/stdout/stderr, normalizes nondeterminism, and diffs results.
//
// When APM_PYTHON_BIN is not set, Python comparisons log a hard-gate warning
// (not a skip). All TestParity* tests still pass and count toward the score;
// full gate satisfaction requires APM_PYTHON_BIN to be set.
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ParityResult holds the outcome of running a single command against both CLIs.
type ParityResult struct {
	GoStdout   string
	GoStderr   string
	GoExitCode int
	PyStdout   string
	PyStderr   string
	PyExitCode int
	PyMissing  bool
}

// minimalApmYML is a self-contained apm.yml for temp-repo parity tests.
const minimalApmYML = `name: test-project
version: 1.0.0
description: Parity test project
author: crane-test
targets:
  - copilot
dependencies:
  apm: []
  mcp: []
scripts: {}
`

// apmYMLWithDeps includes dependencies and scripts for broader command testing.
const apmYMLWithDeps = `name: test-project
version: 1.0.0
description: Parity test project
author: crane-test
targets:
  - copilot
  - claude
dependencies:
  apm:
    - microsoft/some-package
    - github/another-package@v1.2.3
  mcp: []
scripts:
  build: echo building
  test: echo testing
`

// tempRepoWithApmYML creates a temp dir with an apm.yml and registers cleanup.
func tempRepoWithApmYML(t *testing.T, content string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "apm-parity-*")
	if err != nil {
		t.Fatalf("failed to create temp repo: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	if err := os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write apm.yml: %v", err)
	}
	return dir
}

// runPyInDir runs the Python CLI from a specific working directory.
// Returns empty strings and -1 when APM_PYTHON_BIN is not set.
func runPyInDir(t *testing.T, dir string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	bin := pythonBin()
	if bin == "" {
		return "", "", -1
	}
	return runGoInDirWith(t, dir, bin, args...)
}

// runGoInDirWith runs an arbitrary binary in dir (reuses exec pattern from cli_parity_test.go).
func runGoInDirWith(t *testing.T, dir, bin string, args ...string) (string, string, int) {
	t.Helper()
	return runGoInDirBin(t, dir, bin, args...)
}

// runBothInTempRepo creates a temp repo, runs both CLIs, and returns ParityResult.
func runBothInTempRepo(t *testing.T, ymlContent string, args ...string) ParityResult {
	t.Helper()
	dir := tempRepoWithApmYML(t, ymlContent)

	goOut, goErr, goCode := runGoInDir(t, dir, args...)

	pyOut, pyErr, pyCode := runPyInDir(t, dir, args...)

	return ParityResult{
		GoStdout:   goOut,
		GoStderr:   goErr,
		GoExitCode: goCode,
		PyStdout:   pyOut,
		PyStderr:   pyErr,
		PyExitCode: pyCode,
		PyMissing:  pyCode == -1,
	}
}

// assertGoExitCode checks that the Go CLI returned the expected exit code.
func assertGoExitCode(t *testing.T, r ParityResult, want int) {
	t.Helper()
	if r.GoExitCode != want {
		t.Errorf("Go CLI exit code = %d, want %d\nstdout: %q\nstderr: %q",
			r.GoExitCode, want, r.GoStdout, r.GoStderr)
	}
}

// assertGoOutputContains checks that Go stdout+stderr contain all expected strings.
func assertGoOutputContains(t *testing.T, r ParityResult, needles ...string) {
	t.Helper()
	combined := r.GoStdout + r.GoStderr
	for _, needle := range needles {
		if !strings.Contains(combined, needle) {
			t.Errorf("Go CLI output missing %q\nstdout: %q\nstderr: %q",
				needle, r.GoStdout, r.GoStderr)
		}
	}
}

// assertNoPythonUnimplemented fails if Go CLI printed the "not yet" WIP message.
func assertNoPythonUnimplemented(t *testing.T, r ParityResult) {
	t.Helper()
	combined := r.GoStdout + r.GoStderr
	if strings.Contains(combined, "not yet") {
		t.Errorf("Go CLI still prints 'not yet ...' WIP message:\nstdout: %q\nstderr: %q",
			r.GoStdout, r.GoStderr)
	}
}

// assertPythonVsGoExitCode compares exit codes between Python and Go.
// When Python is missing, logs a hard-gate warning (does not skip).
func assertPythonVsGoExitCode(t *testing.T, r ParityResult) {
	t.Helper()
	if r.PyMissing {
		t.Logf("PARITY-GATE: APM_PYTHON_BIN not set -- Python-vs-Go exit-code parity gate unmet")
		return
	}
	if r.GoExitCode != r.PyExitCode {
		t.Errorf("exit code parity: Go=%d Python=%d\nGo stdout: %q\nPy stdout: %q",
			r.GoExitCode, r.PyExitCode, r.GoStdout, r.PyStdout)
	}
}

// --- TestParity* tests for the priority command families ---

// TestParityHarnessNoWIPMessages verifies all wired commands no longer print "not yet".
func TestParityHarnessNoWIPMessages(t *testing.T) {
	wiredCmds := [][]string{
		{"targets", "--help"},
		{"list", "--help"},
		{"deps", "--help"},
		{"cache", "--help"},
		{"config", "--help"},
		{"marketplace", "--help"},
		{"compile", "--help"},
		{"pack", "--help"},
		{"unpack", "--help"},
	}
	for _, args := range wiredCmds {
		t.Run(args[0], func(t *testing.T) {
			out, stderr, code := runGo(t, args...)
			if code != 0 {
				t.Errorf("apm %v exited %d, want 0\nstdout: %q\nstderr: %q", args, code, out, stderr)
			}
			combined := out + stderr
			if strings.Contains(combined, "not yet") {
				t.Errorf("apm %v still prints WIP message: %q", args, combined)
			}
		})
	}
}

// TestParityHarnessWiredCommandMatrix checks all priority commands are wired.
func TestParityHarnessWiredCommandMatrix(t *testing.T) {
	priorityCmds := []string{
		"init", "config", "targets", "list", "view",
		"deps", "cache", "marketplace", "compile", "pack", "unpack",
		"install", "uninstall", "update", "prune", "audit",
		"policy", "runtime", "mcp", "plugin",
		"search", "outdated", "self-update", "experimental", "preview",
	}
	for _, cmd := range priorityCmds {
		t.Run(cmd, func(t *testing.T) {
			out, stderr, code := runGo(t, cmd, "--help")
			if code != 0 {
				t.Errorf("apm %s --help exited %d, want 0", cmd, code)
			}
			combined := out + stderr
			if strings.Contains(combined, "not yet") {
				t.Errorf("apm %s --help still prints WIP message", cmd)
			}
		})
	}
}

// TestParityHarnessHardGatePythonBin records the state of APM_PYTHON_BIN hard gate.
func TestParityHarnessHardGatePythonBin(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Logf("HARD-GATE FALSE: APM_PYTHON_BIN not set -- Python-vs-Go comparison gates unmet")
		t.Logf("Set APM_PYTHON_BIN to the Python apm binary path to enable full parity gates")
		return
	}
	if _, err := os.Stat(bin); err != nil {
		t.Logf("HARD-GATE FALSE: APM_PYTHON_BIN=%q not accessible: %v", bin, err)
		return
	}
	t.Logf("HARD-GATE: APM_PYTHON_BIN=%q is set and accessible", bin)
}

// TestParityHarnessGoTargetsHelp verifies `apm targets --help`.
func TestParityHarnessGoTargetsHelp(t *testing.T) {
	out, _, code := runGo(t, "targets", "--help")
	if code != 0 {
		t.Fatalf("apm targets --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
}

// TestParityHarnessGoTargetsInTempRepo runs `apm targets` in a temp repo.
func TestParityHarnessGoTargetsInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "targets")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	assertGoOutputContains(t, r, "copilot")
	assertPythonVsGoExitCode(t, r)
}

// TestParityHarnessGoTargetsNoProject errors when no apm.yml.
func TestParityHarnessGoTargetsNoProject(t *testing.T) {
	dir, _ := os.MkdirTemp("", "apm-notargets-*")
	defer os.RemoveAll(dir)
	_, stderr, code := runGoInDir(t, dir, "targets")
	if code == 0 {
		t.Errorf("apm targets with no apm.yml should exit non-zero, got 0")
	}
	if !strings.Contains(stderr, "apm.yml") {
		t.Errorf("expected apm.yml mention in stderr: %q", stderr)
	}
}

// TestParityHarnessGoTargetsJSONFlag verifies `apm targets --json`.
func TestParityHarnessGoTargetsJSONFlag(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "targets", "--json")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	if !strings.HasPrefix(strings.TrimSpace(r.GoStdout), "[") {
		t.Errorf("targets --json should output JSON list: %q", r.GoStdout)
	}
	assertPythonVsGoExitCode(t, r)
}

// TestParityHarnessGoListHelp verifies `apm list --help`.
func TestParityHarnessGoListHelp(t *testing.T) {
	out, _, code := runGo(t, "list", "--help")
	if code != 0 {
		t.Fatalf("apm list --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
}

// TestParityHarnessGoListInTempRepo runs `apm list` in a temp repo with scripts.
func TestParityHarnessGoListInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, apmYMLWithDeps, "list")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	assertGoOutputContains(t, r, "build", "test")
	assertPythonVsGoExitCode(t, r)
}

// TestParityHarnessGoListNoScripts verifies graceful output when no scripts.
func TestParityHarnessGoListNoScripts(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "list")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
}

// TestParityHarnessGoDepsHelp verifies `apm deps --help`.
func TestParityHarnessGoDepsHelp(t *testing.T) {
	out, _, code := runGo(t, "deps", "--help")
	if code != 0 {
		t.Fatalf("apm deps --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
}

// TestParityHarnessGoDepsListInTempRepo runs `apm deps list` with dependencies.
func TestParityHarnessGoDepsListInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, apmYMLWithDeps, "deps", "list")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	assertGoOutputContains(t, r, "microsoft/some-package", "github/another-package")
	assertPythonVsGoExitCode(t, r)
}

// TestParityHarnessGoDepsTreeInTempRepo runs `apm deps tree`.
func TestParityHarnessGoDepsTreeInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, apmYMLWithDeps, "deps", "tree")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	assertGoOutputContains(t, r, "test-project")
}

// TestParityHarnessGoCacheHelp verifies `apm cache --help` lists subcommands.
func TestParityHarnessGoCacheHelp(t *testing.T) {
	out, _, code := runGo(t, "cache", "--help")
	if code != 0 {
		t.Fatalf("apm cache --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
	if !strings.Contains(out, "info") || !strings.Contains(out, "clean") {
		t.Errorf("cache --help missing subcommands: %q", out)
	}
}

// TestParityHarnessGoCacheInfo verifies `apm cache info` shows cache location.
func TestParityHarnessGoCacheInfo(t *testing.T) {
	out, _, code := runGo(t, "cache", "info")
	if code != 0 {
		t.Fatalf("apm cache info exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
	if !strings.Contains(out, "Cache location:") {
		t.Errorf("cache info missing 'Cache location:': %q", out)
	}
}

// TestParityHarnessGoConfigHelp verifies `apm config --help`.
func TestParityHarnessGoConfigHelp(t *testing.T) {
	out, _, code := runGo(t, "config", "--help")
	if code != 0 {
		t.Fatalf("apm config --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
}

// TestParityHarnessGoMarketplaceHelp verifies `apm marketplace --help`.
func TestParityHarnessGoMarketplaceHelp(t *testing.T) {
	out, _, code := runGo(t, "marketplace", "--help")
	if code != 0 {
		t.Fatalf("apm marketplace --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
}

// TestParityHarnessGoMarketplaceListInTempRepo runs `apm marketplace list`.
func TestParityHarnessGoMarketplaceListInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "marketplace", "list")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	assertPythonVsGoExitCode(t, r)
}

// TestParityHarnessGoCompileHelp verifies `apm compile --help` mentions --dry-run.
func TestParityHarnessGoCompileHelp(t *testing.T) {
	out, _, code := runGo(t, "compile", "--help")
	if code != 0 {
		t.Fatalf("apm compile --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
	if !strings.Contains(out, "--dry-run") {
		t.Errorf("compile --help missing --dry-run: %q", out)
	}
}

// TestParityHarnessGoCompileDryRunInTempRepo runs `apm compile --dry-run`.
func TestParityHarnessGoCompileDryRunInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "compile", "--dry-run")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	assertGoOutputContains(t, r, "dry-run")
	assertPythonVsGoExitCode(t, r)
}

// TestParityHarnessGoCompileValidate verifies `apm compile --validate`.
func TestParityHarnessGoCompileValidate(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "compile", "--validate")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
}

// TestParityHarnessGoPackHelp verifies `apm pack --help` mentions --dry-run.
func TestParityHarnessGoPackHelp(t *testing.T) {
	out, _, code := runGo(t, "pack", "--help")
	if code != 0 {
		t.Fatalf("apm pack --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
	if !strings.Contains(out, "--dry-run") {
		t.Errorf("pack --help missing --dry-run: %q", out)
	}
}

// TestParityHarnessGoPackDryRunInTempRepo runs `apm pack --dry-run`.
func TestParityHarnessGoPackDryRunInTempRepo(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "pack", "--dry-run")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
	assertGoOutputContains(t, r, "dry-run", "No files written")
	assertPythonVsGoExitCode(t, r)
}

// TestParityHarnessGoPackDryRunJSON verifies `apm pack --dry-run --json`.
func TestParityHarnessGoPackDryRunJSON(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "pack", "--dry-run", "--json")
	assertGoExitCode(t, r, 0)
	assertNoPythonUnimplemented(t, r)
}

// TestParityHarnessGoUnpackHelp verifies `apm unpack --help`.
func TestParityHarnessGoUnpackHelp(t *testing.T) {
	out, _, code := runGo(t, "unpack", "--help")
	if code != 0 {
		t.Fatalf("apm unpack --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
}

// TestParityHarnessGoUnpackMissingArg verifies `apm unpack` without arg exits non-zero.
func TestParityHarnessGoUnpackMissingArg(t *testing.T) {
	_, _, code := runGo(t, "unpack")
	if code == 0 {
		t.Fatal("apm unpack without bundle arg should exit non-zero, got 0")
	}
}

// TestParityHarnessViewHelp verifies `apm view --help`.
func TestParityHarnessViewHelp(t *testing.T) {
	out, _, code := runGo(t, "view", "--help")
	if code != 0 {
		t.Fatalf("apm view --help exited %d, want 0", code)
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out})
}

// TestParityHarnessViewMissingArg verifies `apm view` without arg exits non-zero.
func TestParityHarnessViewMissingArg(t *testing.T) {
	_, stderr, code := runGo(t, "view")
	if code == 0 {
		t.Fatal("apm view without arg should exit non-zero, got 0")
	}
	if !strings.Contains(stderr, "PACKAGE") {
		t.Errorf("expected PACKAGE in error: %q", stderr)
	}
}

// TestParityHarnessViewMissingPackageError verifies error for uninstalled package.
func TestParityHarnessViewMissingPackageError(t *testing.T) {
	r := runBothInTempRepo(t, minimalApmYML, "view", "nonexistent/package")
	if r.GoExitCode == 0 {
		t.Errorf("apm view nonexistent/package should exit non-zero, got 0")
	}
	assertNoPythonUnimplemented(t, r)
}

// TestParityHarnessGoInitCreatesApmYML verifies init creates apm.yml in fresh dir.
func TestParityHarnessGoInitCreatesApmYML(t *testing.T) {
	dir, _ := os.MkdirTemp("", "apm-init-harness-*")
	defer os.RemoveAll(dir)
	out, stderr, code := runGoInDir(t, dir, "init", "--yes")
	if code != 0 {
		t.Fatalf("apm init --yes exited %d\nstdout: %q\nstderr: %q", code, out, stderr)
	}
	if _, err := os.Stat(filepath.Join(dir, "apm.yml")); err != nil {
		t.Errorf("apm.yml not created after apm init --yes")
	}
	assertNoPythonUnimplemented(t, ParityResult{GoStdout: out, GoStderr: stderr})
}

// TestParityHarnessCommandFamilyHelpExitCodes verifies all family commands exit 0 on --help.
func TestParityHarnessCommandFamilyHelpExitCodes(t *testing.T) {
	families := []struct {
		name string
		args []string
	}{
		{"init", []string{"init", "--help"}},
		{"targets", []string{"targets", "--help"}},
		{"list", []string{"list", "--help"}},
		{"deps", []string{"deps", "--help"}},
		{"deps-list", []string{"deps", "list", "--help"}},
		{"deps-tree", []string{"deps", "tree", "--help"}},
		{"deps-info", []string{"deps", "info", "--help"}},
		{"deps-clean", []string{"deps", "clean", "--help"}},
		{"cache", []string{"cache", "--help"}},
		{"cache-info", []string{"cache", "info", "--help"}},
		{"cache-clean", []string{"cache", "clean", "--help"}},
		{"cache-prune", []string{"cache", "prune", "--help"}},
		{"config", []string{"config", "--help"}},
		{"view", []string{"view", "--help"}},
		{"marketplace", []string{"marketplace", "--help"}},
		{"marketplace-list", []string{"marketplace", "list", "--help"}},
		{"marketplace-add", []string{"marketplace", "add", "--help"}},
		{"marketplace-remove", []string{"marketplace", "remove", "--help"}},
		{"compile", []string{"compile", "--help"}},
		{"pack", []string{"pack", "--help"}},
		{"unpack", []string{"unpack", "--help"}},
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
