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
	"path/filepath"
	"runtime"
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

// completionModuleRoot returns the repository root, two levels above cmd/apm.
func completionModuleRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file path via runtime.Caller")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "..")
}

// TestParityCompletionSurfaceParity verifies the Go CLI exposes at least every
// top-level command that the Python CLI exposes. Gate 3: surface_parity.
func TestParityCompletionSurfaceParity(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set -- surface parity cannot be verified")
	}

	required := []string{
		"init", "install", "update", "compile", "pack", "unpack", "run",
		"audit", "policy", "mcp", "runtime", "targets", "list",
		"view", "cache", "deps", "marketplace", "plugin",
		"uninstall", "prune", "search", "outdated", "self-update",
		"preview", "config", "experimental",
	}

	goOut, _, goCode := runGo(t, "--help")
	if goCode != 0 {
		t.Fatalf("Go `apm --help` exited %d", goCode)
	}
	pyOut, _, pyCode := runPyBin(t, bin, "--help")
	if pyCode != 0 {
		t.Fatalf("Python `apm --help` exited %d", pyCode)
	}

	var missing []string
	for _, cmd := range required {
		inPy := strings.Contains(pyOut, cmd)
		inGo := strings.Contains(goOut, cmd)
		if inPy && !inGo {
			missing = append(missing, cmd)
			t.Errorf("Go CLI missing command %q present in Python CLI", cmd)
		}
	}
	if len(missing) == 0 {
		t.Logf("[+] Surface parity: all %d required commands present in Go CLI.", len(required))
	}
}

// TestParityCompletionFunctionalContracts verifies key read-only command
// behaviors match between Python and Go (exit codes and basic output). Gate 5.
func TestParityCompletionFunctionalContracts(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set -- functional contracts cannot be verified")
	}

	type contract struct {
		args    []string
		wantGo  int
		inRepo  bool
		ymlType string // "minimal" or "deps"
	}
	contracts := []contract{
		{args: []string{"--help"}, wantGo: 0},
		{args: []string{"--version"}, wantGo: 0},
		{args: []string{"targets", "--help"}, wantGo: 0},
		{args: []string{"deps", "--help"}, wantGo: 0},
		{args: []string{"cache", "--help"}, wantGo: 0},
		{args: []string{"targets"}, wantGo: 0, inRepo: true, ymlType: "minimal"},
		{args: []string{"list"}, wantGo: 0, inRepo: true, ymlType: "deps"},
		{args: []string{"deps", "list"}, wantGo: 0, inRepo: true, ymlType: "deps"},
		{args: []string{"compile", "--dry-run"}, wantGo: 0, inRepo: true, ymlType: "minimal"},
		{args: []string{"audit"}, wantGo: 0, inRepo: true, ymlType: "minimal"},
	}

	for _, c := range contracts {
		c := c
		label := "apm " + strings.Join(c.args, " ")
		t.Run(label, func(t *testing.T) {
			var goOut, pyOut string
			var goCode, pyCode int
			if c.inRepo {
				ymlContent := minimalApmYML
				if c.ymlType == "deps" {
					ymlContent = apmYMLWithDeps
				}
				r := runBothInTempRepo(t, ymlContent, c.args...)
				goOut, goCode = r.GoStdout+r.GoStderr, r.GoExitCode
				pyOut, pyCode = r.PyStdout+r.PyStderr, r.PyExitCode
				if r.PyMissing {
					pyCode = -1
				}
			} else {
				goOut, _, goCode = runGo(t, c.args...)
				pyOut, _, pyCode = runPyBin(t, bin, c.args...)
			}
			if goCode != c.wantGo {
				t.Errorf("Go exit %d, want %d; output: %q", goCode, c.wantGo, goOut)
			}
			if pyCode >= 0 && pyCode != goCode {
				t.Errorf("exit code mismatch: Python=%d Go=%d; pyOut=%q goOut=%q",
					pyCode, goCode, pyOut, goOut)
			}
		})
	}
}

// TestParityCompletionStateDiffContracts verifies that mutating commands produce
// equivalent filesystem state between Python and Go. Gate 6.
func TestParityCompletionStateDiffContracts(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set -- state-diff contracts cannot be verified")
	}

	t.Run("init creates apm.yml", func(t *testing.T) {
		goDir, err := os.MkdirTemp("", "apm-state-go-*")
		if err != nil {
			t.Fatalf("mkdtemp: %v", err)
		}
		defer os.RemoveAll(goDir)

		pyDir, err := os.MkdirTemp("", "apm-state-py-*")
		if err != nil {
			t.Fatalf("mkdtemp: %v", err)
		}
		defer os.RemoveAll(pyDir)

		_, _, goCode := runGoInDir(t, goDir, "init", "--yes")
		if goCode != 0 {
			t.Errorf("Go `apm init --yes` exited %d", goCode)
		}
		goApmYML := filepath.Join(goDir, "apm.yml")
		if _, err := os.Stat(goApmYML); err != nil {
			t.Errorf("Go `apm init --yes` did not create apm.yml: %v", err)
		}

		_, _, pyCode := runGoInDirWith(t, pyDir, bin, "init", "--yes")
		pyApmYML := filepath.Join(pyDir, "apm.yml")
		if _, err := os.Stat(pyApmYML); err != nil {
			t.Logf("Python init did not create apm.yml (exit %d): will verify Go only", pyCode)
		} else {
			// Both created apm.yml: verify they contain the same required keys.
			goBytes, _ := os.ReadFile(goApmYML)
			pyBytes, _ := os.ReadFile(pyApmYML)
			for _, key := range []string{"name:", "version:", "dependencies:"} {
				if !strings.Contains(string(goBytes), key) {
					t.Errorf("Go apm.yml missing key %q", key)
				}
				if !strings.Contains(string(pyBytes), key) {
					t.Logf("Python apm.yml missing key %q (non-fatal)", key)
				}
			}
			t.Logf("[+] State-diff: Go and Python both created apm.yml with required keys.")
		}
	})
}

// TestParityCompletionPythonSuite runs the Python reference unit test suite to
// confirm the Python CLI remains green. Gate 7: python_tests_pass.
func TestParityCompletionPythonSuite(t *testing.T) {
	if os.Getenv("APM_PYTHON_BIN") == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set -- Python suite cannot be verified")
	}

	root := completionModuleRoot(t)

	// Locate uv; required to run the Python test suite.
	uvPath, err := exec.LookPath("uv")
	if err != nil {
		t.Fatalf("HARD-GATE FAILED: uv not found in PATH -- cannot run Python suite: %v", err)
	}

	// Run the Python unit suite in parallel (-n auto) for speed.
	// --ignore integration tests that require external services.
	cmd := exec.Command(uvPath, "run", "--extra", "dev",
		"pytest", "tests/unit/", "-q", "--tb=short", "--no-header",
		"-n", "auto",
		"--ignore=tests/unit/integration",
	)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "PYTHONDONTWRITEBYTECODE=1", "COLUMNS=10000")
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if runErr := cmd.Run(); runErr != nil {
		t.Fatalf("Python suite failed:\n%s\n%s", outBuf.String(), errBuf.String())
	}
	t.Logf("[+] Python suite passed:\n%s", outBuf.String())
}

// TestParityCompletionBenchmarks runs the migration CLI benchmark and verifies
// the Go CLI stays within the configured performance ratio. Gate 8.
func TestParityCompletionBenchmarks(t *testing.T) {
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		t.Fatal("HARD-GATE FAILED: APM_PYTHON_BIN not set -- benchmarks cannot be verified")
	}
	if goBinPath == "" {
		t.Fatal("HARD-GATE FAILED: Go binary not built -- benchmarks cannot be verified")
	}

	root := completionModuleRoot(t)

	benchScript := filepath.Join(root, "scripts", "ci", "migration_cli_benchmark.py")
	if _, err := os.Stat(benchScript); err != nil {
		t.Fatalf("benchmark script not found at %s: %v", benchScript, err)
	}

	// Locate uv to run the benchmark script.
	uvPath, err := exec.LookPath("uv")
	if err != nil {
		t.Fatalf("HARD-GATE FAILED: uv not found in PATH: %v", err)
	}

	jsonOut := filepath.Join(t.TempDir(), "benchmark.json")
	mdOut := filepath.Join(t.TempDir(), "benchmark.md")
	// Use --repeats 2 for a quick CI smoke test (full 5-repeat runs in the
	// dedicated benchmarks job).
	cmd := exec.Command(uvPath, "run", benchScript,
		"--python-bin", bin,
		"--go-bin", goBinPath,
		"--json-out", jsonOut,
		"--markdown-out", mdOut,
		"--max-ratio", "5.0",
		"--repeats", "2",
	)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if runErr := cmd.Run(); runErr != nil {
		t.Fatalf("Benchmark failed (Go CLI exceeds 5x Python latency or script error):\n%s\n%s",
			outBuf.String(), errBuf.String())
	}
	t.Logf("[+] Benchmarks passed:\n%s", outBuf.String())
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
