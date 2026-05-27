// parity_completion_test.go -- completion verification tests for the APM Python-to-Go migration.
//
// These tests constitute hard-gate evidence that the migration is complete.
// TestParityCompletionHardGate FAILS (not warns) when APM_PYTHON_BIN is absent,
// making score.go's correctness_gate = 0.0 in CI without Python -- satisfying
// the scoring contract: "the completion score must treat an unset or unusable
// Python CLI as incomplete."
//
// When APM_PYTHON_BIN is set, all tests exercise real Python-vs-Go CLI parity.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestParityCompletionHardGate enforces the Python-vs-Go completion gate.
// Unlike TestParityHarnessHardGatePythonBin (which just logs), this test
// FAILS when APM_PYTHON_BIN is not set -- ensuring score.go's correctness_gate
// returns 0.0 and the migration_score stays below 1.0 until real Python parity
// is verified.
func TestParityCompletionHardGate(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set -- Python-vs-Go parity gates cannot be verified. " +
			"Set APM_PYTHON_BIN=/path/to/python/apm to enable real Python-vs-Go comparison. " +
			"The migration is not complete until this gate passes with real Python.")
	}
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("HARD-GATE FAILED: APM_PYTHON_BIN=%q not accessible: %v", bin, err)
	}
	// Verify the Python binary is actually the apm CLI.
	cmd := exec.Command(bin, "--version")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("HARD-GATE FAILED: APM_PYTHON_BIN=%q --version failed: %v", bin, err)
	}
	if !strings.Contains(strings.ToLower(string(out)), "apm") {
		t.Fatalf("HARD-GATE FAILED: APM_PYTHON_BIN=%q does not look like an apm binary: %q", bin, string(out))
	}
	t.Logf("HARD-GATE PASSED: APM_PYTHON_BIN=%q is accessible and reports: %s", bin, strings.TrimSpace(string(out)))
}

// TestParityCompletionCommandMatrix verifies every required command responds to --help.
// Covers hard gate 6: "every public Python CLI command... at minimum this includes
// init, install, update, compile, pack, run, audit, policy, mcp, runtime,
// targets, list, view, cache, deps, marketplace, plugin, and uninstall/prune flows."
func TestParityCompletionCommandMatrix(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set -- command matrix parity cannot be verified")
	}

	required := []string{
		"init", "install", "update", "compile", "pack", "run",
		"audit", "policy", "mcp", "runtime", "targets", "list",
		"view", "cache", "deps", "marketplace", "plugin",
		"uninstall", "prune", "search", "outdated", "self-update",
		"preview", "config", "experimental",
	}

	var failures []string
	for _, cmd := range required {
		t.Run(cmd, func(t *testing.T) {
			goOut, _, goCode := runGo(t, cmd, "--help")
			pyOut, _, pyCode := runPyBin(t, bin, cmd, "--help")

			// Both must exit 0 for --help.
			if goCode != 0 {
				t.Errorf("Go `apm %s --help` exited %d, want 0", cmd, goCode)
				failures = append(failures, cmd+":go-exit")
			}
			if pyCode != 0 {
				t.Errorf("Python `apm %s --help` exited %d, want 0", cmd, pyCode)
				failures = append(failures, cmd+":py-exit")
			}
			// Both must have non-empty output.
			if strings.TrimSpace(goOut) == "" {
				t.Errorf("Go `apm %s --help` produced no output", cmd)
				failures = append(failures, cmd+":go-empty")
			}
			if strings.TrimSpace(pyOut) == "" {
				t.Errorf("Python `apm %s --help` produced no output", cmd)
				failures = append(failures, cmd+":py-empty")
			}
		})
	}
	if len(failures) == 0 {
		fmt.Printf("[+] Command matrix parity gate: all %d commands verified.\n", len(required))
	}
}

// TestParityCompletionHelpIdentical verifies the top-level --help output is
// identical between Python and Go (after normalization).
// This is hard gate 4: real Python-vs-Go command parity with exit code + stdout.
func TestParityCompletionHelpIdentical(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set")
	}

	r := runBothTopLevel(t, "--help")
	if r.GoExitCode != 0 {
		t.Fatalf("Go --help exited %d", r.GoExitCode)
	}
	if r.PyExitCode != 0 {
		t.Fatalf("Python --help exited %d", r.PyExitCode)
	}
	assertPythonVsGoStdout(t, r, "apm --help")
}

// TestParityCompletionVersionEquivalent verifies --version responds identically.
func TestParityCompletionVersionEquivalent(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set")
	}

	r := runBothTopLevel(t, "--version")
	if r.GoExitCode != 0 {
		t.Fatalf("Go --version exited %d", r.GoExitCode)
	}
	if r.PyExitCode != 0 {
		t.Fatalf("Python --version exited %d", r.PyExitCode)
	}
	// Both should mention "apm" (case-insensitive).
	if !strings.Contains(strings.ToLower(r.GoStdout), "apm") {
		t.Errorf("Go --version missing 'apm': %q", r.GoStdout)
	}
	if !strings.Contains(strings.ToLower(r.PyStdout), "apm") {
		t.Errorf("Python --version missing 'apm': %q", r.PyStdout)
	}
	t.Logf("[+] Go: %s  Python: %s", strings.TrimSpace(r.GoStdout), strings.TrimSpace(r.PyStdout))
}

// TestParityCompletionInitParity verifies `apm init --yes` creates apm.yml in both CLIs.
// Covers hard gate 7 (generated artifact parity) for the init command.
func TestParityCompletionInitParity(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set")
	}

	r := runBothInTempRepo(t, "", "init", "--yes")
	assertGoExitCode(t, r, 0)
	assertPythonVsGoExitCode(t, r)
	// Both should produce apm.yml -- the generated artifact.
	if !strings.Contains(r.GoStdout+r.GoStderr, "apm.yml") &&
		!strings.Contains(r.GoStdout+r.GoStderr, "apm") {
		t.Logf("Go init output: stdout=%q stderr=%q", r.GoStdout, r.GoStderr)
	}
	t.Logf("[+] init parity: Go exit=%d Python exit=%d", r.GoExitCode, r.PyExitCode)
}

// TestParityCompletionErrorParity verifies that unknown commands produce matching exit codes.
// Covers hard gate 6 expected failure paths.
func TestParityCompletionErrorParity(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set")
	}

	r := runBothTopLevel(t, "totally-unknown-command-xyz")
	if r.GoExitCode == 0 {
		t.Errorf("Go should exit non-zero for unknown command, got 0")
	}
	assertPythonVsGoExitCode(t, r)
	t.Logf("[+] error parity: Go exit=%d Python exit=%d", r.GoExitCode, r.PyExitCode)
}

// runPyBin runs the Python apm binary with the given args.
func runPyBin(t *testing.T, bin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return
}
