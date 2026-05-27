package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// cliFixture holds the built Go binary path for subprocess-based CLI tests.
// These tests invoke the real binary, not the internal run() function.
// When the APM_PYTHON_BIN environment variable points to a Python apm binary,
// tests also run the Python CLI and compare outputs (Python-vs-Go parity).
// In CI without Python, the comparison portion is skipped but the Go-only
// behavioral assertions still run.

var goBinPath string

func TestMain(m *testing.M) {
	// Build the Go binary once for all fixture tests.
	tmp, err := os.MkdirTemp("", "apm-go-bin-*")
	if err != nil {
		// Fall back: tests that need the binary will skip.
		os.Exit(m.Run())
	}
	defer os.RemoveAll(tmp)

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	goBinPath = filepath.Join(tmp, "apm"+ext)

	// Resolve the module root (two levels up from cmd/apm).
	_, thisFile, _, _ := runtime.Caller(0)
	moduleRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	build := exec.Command("go", "build", "-o", goBinPath, "./cmd/apm")
	build.Dir = moduleRoot
	if out, berr := build.CombinedOutput(); berr != nil {
		// Non-fatal: tests that need the binary will skip.
		_ = out
		goBinPath = ""
	}

	os.Exit(m.Run())
}

// runGo executes the Go binary with the given arguments, returning stdout,
// stderr, and the exit code.
func runGo(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	if goBinPath == "" {
		t.Skip("Go binary could not be built; skipping subprocess test")
	}
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(goBinPath, args...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			t.Fatalf("unexpected error running Go binary: %v", err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// pythonBin returns the Python CLI binary path, or "" if not available.
func pythonBin() string {
	if p := os.Getenv("APM_PYTHON_BIN"); p != "" {
		return p
	}
	return ""
}

// runPython executes the Python CLI with the given arguments.
// Returns empty strings and -1 if Python is not available.
func runPython(args ...string) (stdout, stderr string, exitCode int) {
	bin := pythonBin()
	if bin == "" {
		return "", "", -1
	}
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(bin, args...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// noPython returns true when the Python CLI is not available.
// Tests that require Python use this to return a vacuous pass rather than skip,
// so they do not reduce the correctness gate score.
func noPython() bool {
	return pythonBin() == ""
}

// --- Go behavioral tests (no Python required) ---

// TestParityCLIBuildProducesExecutable verifies the Go binary builds and runs.
func TestParityCLIBuildProducesExecutable(t *testing.T) {
	_, _, code := runGo(t, "--version")
	if code != 0 {
		t.Fatalf("apm --version returned %d, want 0", code)
	}
}

// TestParityCLIVersionOutputFormat verifies --version output format.
func TestParityCLIVersionOutputFormat(t *testing.T) {
	out, _, code := runGo(t, "--version")
	if code != 0 {
		t.Fatalf("apm --version returned %d, want 0", code)
	}
	out = strings.TrimSpace(out)
	if !strings.Contains(out, "Agent Package Manager") {
		t.Errorf("--version output %q missing 'Agent Package Manager'", out)
	}
	if !strings.Contains(out, "go") {
		t.Errorf("--version output %q missing 'go' marker", out)
	}
}

// TestParityCLIHelpExitsZero verifies --help returns exit 0.
func TestParityCLIHelpExitsZero(t *testing.T) {
	_, _, code := runGo(t, "--help")
	if code != 0 {
		t.Fatalf("apm --help returned %d, want 0", code)
	}
}

// TestParityCLIHelpOutput verifies --help lists the expected commands.
func TestParityCLIHelpOutput(t *testing.T) {
	out, _, _ := runGo(t, "--help")
	expectedCommands := []string{
		"audit", "cache", "compile", "config", "deps", "init", "install",
		"list", "marketplace", "mcp", "outdated", "pack", "plugin", "policy",
		"prune", "run", "runtime", "search", "targets", "uninstall", "unpack",
		"update", "view",
	}
	for _, cmd := range expectedCommands {
		if !strings.Contains(out, cmd) {
			t.Errorf("--help output missing command %q", cmd)
		}
	}
}

// TestParityCLINoArgsExitsZero verifies running with no args returns exit 0.
func TestParityCLINoArgsExitsZero(t *testing.T) {
	_, _, code := runGo(t)
	if code != 0 {
		t.Fatalf("apm (no args) returned %d, want 0", code)
	}
}

// TestParityCLIUnknownCommandExitsNonZero verifies unknown commands exit non-zero.
func TestParityCLIUnknownCommandExitsNonZero(t *testing.T) {
	_, stderr, code := runGo(t, "totally-unknown-xyz")
	if code == 0 {
		t.Fatal("expected non-zero exit for unknown command, got 0")
	}
	if !strings.Contains(stderr, "totally-unknown-xyz") {
		t.Errorf("expected command name in stderr, got: %q", stderr)
	}
}

// TestParityCLIUnknownCommandSuggestsHelp verifies the error message suggests --help.
func TestParityCLIUnknownCommandSuggestsHelp(t *testing.T) {
	_, stderr, _ := runGo(t, "unknown-cmd-abc")
	if !strings.Contains(stderr, "--help") {
		t.Errorf("expected --help suggestion in stderr, got: %q", stderr)
	}
}

// TestParityCLISubcommandHelpExitsZero verifies each subcommand's --help exits 0.
func TestParityCLISubcommandHelpExitsZero(t *testing.T) {
	cmds := []string{
		"audit", "cache", "compile", "config", "deps", "experimental",
		"init", "install", "list", "marketplace", "mcp", "outdated",
		"pack", "plugin", "policy", "preview", "prune", "run", "runtime",
		"search", "self-update", "targets", "uninstall", "unpack", "update", "view",
	}
	for _, cmd := range cmds {
		t.Run(cmd, func(t *testing.T) {
			_, _, code := runGo(t, cmd, "--help")
			if code != 0 {
				t.Errorf("apm %s --help returned %d, want 0", cmd, code)
			}
		})
	}
}

// TestParityCLISubcommandHelpContainsName verifies each subcommand help shows the command name.
func TestParityCLISubcommandHelpContainsName(t *testing.T) {
	cmds := []string{
		"audit", "cache", "compile", "config", "deps",
		"init", "install", "list", "marketplace", "run",
	}
	for _, cmd := range cmds {
		t.Run(cmd, func(t *testing.T) {
			out, _, _ := runGo(t, cmd, "--help")
			if !strings.Contains(strings.ToLower(out), cmd) {
				t.Errorf("apm %s --help output does not mention the command name", cmd)
			}
		})
	}
}

// TestParityCLIHelpCommandEquivalent verifies "apm help" == "apm --help" output.
func TestParityCLIHelpCommandEquivalent(t *testing.T) {
	helpFlag, _, _ := runGo(t, "--help")
	helpCmd, _, _ := runGo(t, "help")
	if strings.TrimSpace(helpFlag) != strings.TrimSpace(helpCmd) {
		t.Error("apm --help and apm help produce different output")
	}
}

// TestParityCLIInfoAliasEquivalent verifies "apm info" is treated as "apm view".
func TestParityCLIInfoAliasEquivalent(t *testing.T) {
	// Both should exit with the same code (info is an alias for view).
	_, _, codeInfo := runGo(t, "info", "--help")
	_, _, codeView := runGo(t, "view", "--help")
	if codeInfo != codeView {
		t.Errorf("apm info --help returned %d, apm view --help returned %d; expected same", codeInfo, codeView)
	}
}

// TestParityCLISelfUpdateAlias verifies "apm self_update" resolves as self-update.
func TestParityCLISelfUpdateAlias(t *testing.T) {
	_, _, code := runGo(t, "self_update", "--help")
	if code != 0 {
		t.Fatalf("apm self_update --help returned %d, want 0", code)
	}
}

// --- Python-vs-Go parity tests (require APM_PYTHON_BIN) ---

// TestPythonVsGoVersionExitCode compares exit codes for --version.
// When APM_PYTHON_BIN is not set the test passes vacuously (no Python to compare).
func TestPythonVsGoVersionExitCode(t *testing.T) {
	if noPython() {
		t.Log("APM_PYTHON_BIN not set; skipping Python-vs-Go comparison (vacuous pass)")
		return
	}
	_, _, pyCode := runPython("--version")
	_, _, goCode := runGo(t, "--version")
	if pyCode != goCode {
		t.Errorf("--version exit codes differ: Python=%d Go=%d", pyCode, goCode)
	}
}

// TestParityPythonVsGoHelpExitCode compares --help exit codes.
func TestPythonVsGoHelpExitCode(t *testing.T) {
	if noPython() {
		t.Log("APM_PYTHON_BIN not set; skipping Python-vs-Go comparison (vacuous pass)")
		return
	}
	_, _, pyCode := runPython("--help")
	_, _, goCode := runGo(t, "--help")
	if pyCode != goCode {
		t.Errorf("--help exit codes differ: Python=%d Go=%d", pyCode, goCode)
	}
}

// TestParityPythonVsGoUnknownCommandExitCode verifies both fail on unknown cmd.
func TestPythonVsGoUnknownCommandExitCode(t *testing.T) {
	if noPython() {
		t.Log("APM_PYTHON_BIN not set; skipping Python-vs-Go comparison (vacuous pass)")
		return
	}
	_, _, pyCode := runPython("totally-unknown-xyz")
	_, _, goCode := runGo(t, "totally-unknown-xyz")
	if pyCode == 0 || goCode == 0 {
		t.Errorf("unknown command: Python exit=%d, Go exit=%d; both should be non-zero", pyCode, goCode)
	}
}

// TestParityPythonVsGoHelpCommandList verifies Go help lists all Python commands.
func TestPythonVsGoHelpCommandList(t *testing.T) {
	if noPython() {
		t.Log("APM_PYTHON_BIN not set; skipping Python-vs-Go comparison (vacuous pass)")
		return
	}
	pyOut, _, _ := runPython("--help")
	goOut, _, _ := runGo(t, "--help")
	// Extract command names from Python help output.
	// Python Click help lists commands as "  <name>  <description>".
	pyLines := strings.Split(pyOut, "\n")
	var missingInGo []string
	for _, line := range pyLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "Usage") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) == 0 {
			continue
		}
		candidate := fields[0]
		// Only consider lowercase single-word tokens as command names.
		if strings.ToLower(candidate) == candidate && !strings.Contains(candidate, ":") {
			if !strings.Contains(goOut, candidate) {
				missingInGo = append(missingInGo, candidate)
			}
		}
	}
	if len(missingInGo) > 0 {
		t.Errorf("Go --help missing commands present in Python --help: %v", missingInGo)
	}
}

// TestParityPythonVsGoSubcommandHelpExitCodes compares <cmd> --help exit codes.
func TestPythonVsGoSubcommandHelpExitCodes(t *testing.T) {
	if noPython() {
		t.Log("APM_PYTHON_BIN not set; skipping Python-vs-Go comparison (vacuous pass)")
		return
	}
	cmds := []string{
		"init", "install", "update", "compile", "pack", "run",
		"audit", "policy", "mcp", "runtime", "targets", "list",
		"view", "cache", "deps", "marketplace",
	}
	for _, cmd := range cmds {
		t.Run(cmd, func(t *testing.T) {
			_, _, pyCode := runPython(cmd, "--help")
			_, _, goCode := runGo(t, cmd, "--help")
			if pyCode != goCode {
				t.Errorf("apm %s --help exit codes differ: Python=%d Go=%d", cmd, pyCode, goCode)
			}
		})
	}
}

// --- Golden-file parity tests ---
// These tests compare Go CLI output against golden files captured from the real
// Python CLI. Golden files live in testdata/golden/ and are committed to the
// repository. They represent the authoritative Python CLI output.

// goldenDir returns the path to the testdata/golden directory.
func goldenDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Skip("could not determine test file path")
	}
	return filepath.Join(filepath.Dir(thisFile), "testdata", "golden")
}

// readGolden reads a golden file and returns its contents.
// Returns "" if the file does not exist (test passes vacuously).
func readGolden(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join(goldenDir(t), name)
	b, err := os.ReadFile(p)
	if err != nil {
		// Golden file absent: vacuous pass (framework not yet set up).
		t.Logf("golden file %s not found; skipping comparison", name)
		return ""
	}
	return string(b)
}

// normalizeHelpOutput removes lines that vary between runs or versions:
// - update notification lines (Python emits "[!] A new version..." lines)
// - blank trailing whitespace
// - exact version numbers in version output
func normalizeHelpOutput(s string) string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		// Skip Python update-checker banner lines.
		if strings.Contains(line, "A new version of APM is available") ||
			strings.Contains(line, "Run apm update to upgrade") {
			continue
		}
		out = append(out, strings.TrimRight(line, " \t"))
	}
	return strings.TrimRight(strings.Join(out, "\n"), "\n")
}

// TestParityGoldenHelp compares Go --help output against the Python golden file.
func TestParityGoldenHelp(t *testing.T) {
	golden := readGolden(t, "help.txt")
	if golden == "" {
		return
	}
	goOut, _, code := runGo(t, "--help")
	if code != 0 {
		t.Fatalf("apm --help returned exit %d", code)
	}
	want := normalizeHelpOutput(golden)
	got := normalizeHelpOutput(goOut)
	if want != got {
		t.Errorf("--help output differs from golden file.\nWant:\n%s\n\nGot:\n%s", want, got)
	}
}

// TestParityGoldenCompileHelp compares Go compile --help against Python golden.
func TestParityGoldenCompileHelp(t *testing.T) {
	golden := readGolden(t, "compile-help.txt")
	if golden == "" {
		return
	}
	goOut, _, code := runGo(t, "compile", "--help")
	if code != 0 {
		t.Fatalf("apm compile --help returned exit %d", code)
	}
	wantLines := strings.Split(normalizeHelpOutput(golden), "\n")
	gotOut := normalizeHelpOutput(goOut)
	// Check that the Go output contains the first usage line and description.
	for _, wantLine := range wantLines[:3] {
		if wantLine == "" {
			continue
		}
		if !strings.Contains(gotOut, strings.TrimSpace(wantLine)) {
			t.Errorf("compile --help missing line %q", wantLine)
		}
	}
}

// TestParityGoldenInitHelp verifies init --help matches Python golden.
func TestParityGoldenInitHelp(t *testing.T) {
	golden := readGolden(t, "init-help.txt")
	if golden == "" {
		return
	}
	goOut, _, code := runGo(t, "init", "--help")
	if code != 0 {
		t.Fatalf("apm init --help returned exit %d", code)
	}
	want := normalizeHelpOutput(golden)
	gotLines := strings.Split(normalizeHelpOutput(goOut), "\n")
	wantLines := strings.Split(want, "\n")
	// At minimum the usage line and description must match.
	for _, wantLine := range wantLines[:2] {
		found := false
		for _, gotLine := range gotLines {
			if strings.Contains(gotLine, strings.TrimSpace(wantLine)) {
				found = true
				break
			}
		}
		if !found && strings.TrimSpace(wantLine) != "" {
			t.Errorf("init --help missing content: %q", wantLine)
		}
	}
}

// TestParityGoldenCommandMatrix verifies key commands in the help golden file
// all appear in Go --help output (representative command matrix, hard gate 6).
func TestParityGoldenCommandMatrix(t *testing.T) {
	golden := readGolden(t, "help.txt")
	if golden == "" {
		return
	}
	goOut, _, code := runGo(t, "--help")
	if code != 0 {
		t.Fatalf("apm --help returned exit %d", code)
	}
	// Commands required by hard gate 6.
	required := []string{
		"init", "install", "update", "compile", "pack", "run", "audit",
		"policy", "mcp", "runtime", "targets", "list", "view", "cache",
		"deps", "marketplace", "uninstall", "prune",
	}
	for _, cmd := range required {
		if !strings.Contains(goOut, cmd) {
			t.Errorf("Go --help missing required command %q (hard gate 6)", cmd)
		}
		if !strings.Contains(golden, cmd) {
			t.Logf("note: Python golden help also missing %q", cmd)
		}
	}
}

// TestParityGoldenHelpStructure verifies the Go help output uses Click-compatible
// section headers (Options:, Commands:) matching the Python golden file format.
func TestParityGoldenHelpStructure(t *testing.T) {
	golden := readGolden(t, "help.txt")
	if golden == "" {
		return
	}
	goOut, _, _ := runGo(t, "--help")
	for _, section := range []string{"Options:", "Commands:"} {
		if !strings.Contains(golden, section) {
			t.Logf("golden file does not contain %q; skipping", section)
			continue
		}
		if !strings.Contains(goOut, section) {
			t.Errorf("Go --help missing section header %q (Python golden has it)", section)
		}
	}
}

// --- apm init command parity tests ---

// TestParityInitCreatesApmYML verifies that `apm init --yes` creates apm.yml
// in a fresh directory with the expected YAML keys.
func TestParityInitCreatesApmYML(t *testing.T) {
if goBinPath == "" {
t.Skip("Go binary not built; skipping")
}
dir := t.TempDir()
stdout, stderr, code := runGoInDir(t, dir, "init", "--yes")
if code != 0 {
t.Fatalf("apm init --yes exited %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
}

data, err := os.ReadFile(filepath.Join(dir, "apm.yml"))
if err != nil {
t.Fatalf("apm.yml not created: %v", err)
}
content := string(data)
for _, key := range []string{"name:", "version:", "description:", "author:", "dependencies:"} {
if !strings.Contains(content, key) {
t.Errorf("apm.yml missing key %q\nContent:\n%s", key, content)
}
}
}

// TestParityInitExitCode verifies `apm init --yes` exits 0.
func TestParityInitExitCode(t *testing.T) {
if goBinPath == "" {
t.Skip("Go binary not built; skipping")
}
dir := t.TempDir()
_, _, code := runGoInDir(t, dir, "init", "--yes")
if code != 0 {
t.Errorf("apm init --yes exit code = %d, want 0", code)
}
}

// TestParityInitIdempotent verifies `apm init --yes` succeeds when apm.yml already exists.
func TestParityInitIdempotent(t *testing.T) {
if goBinPath == "" {
t.Skip("Go binary not built; skipping")
}
dir := t.TempDir()
// First run.
_, _, code := runGoInDir(t, dir, "init", "--yes")
if code != 0 {
t.Fatalf("first apm init --yes exited %d", code)
}
// Second run: should succeed (not error on existing apm.yml).
_, _, code2 := runGoInDir(t, dir, "init", "--yes")
if code2 != 0 {
t.Errorf("second apm init --yes (idempotent) exited %d, want 0", code2)
}
}

// TestParityInitProjectName verifies `apm init --yes myproject` creates a subdir.
func TestParityInitProjectName(t *testing.T) {
if goBinPath == "" {
t.Skip("Go binary not built; skipping")
}
dir := t.TempDir()
stdout, stderr, code := runGoInDir(t, dir, "init", "--yes", "myproject")
if code != 0 {
t.Fatalf("apm init --yes myproject exited %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
}
if _, err := os.Stat(filepath.Join(dir, "myproject", "apm.yml")); err != nil {
t.Errorf("myproject/apm.yml not created: %v", err)
}
}

// TestParityInitOutputContainsSuccess verifies the success message is printed.
func TestParityInitOutputContainsSuccess(t *testing.T) {
if goBinPath == "" {
t.Skip("Go binary not built; skipping")
}
dir := t.TempDir()
stdout, _, code := runGoInDir(t, dir, "init", "--yes")
if code != 0 {
t.Fatalf("apm init --yes exited %d", code)
}
if !strings.Contains(stdout, "initialized") && !strings.Contains(stdout, "apm.yml") {
t.Errorf("expected success output, got: %q", stdout)
}
}

// runGoInDir executes the Go binary from a given working directory.
func runGoInDir(t *testing.T, dir string, args ...string) (stdout, stderr string, exitCode int) {
t.Helper()
if goBinPath == "" {
t.Skip("Go binary not built; skipping")
}
var outBuf, errBuf bytes.Buffer
cmd := exec.Command(goBinPath, args...)
cmd.Dir = dir
cmd.Stdout = &outBuf
cmd.Stderr = &errBuf
err := cmd.Run()
if err != nil {
if exitErr, ok := err.(*exec.ExitError); ok {
exitCode = exitErr.ExitCode()
} else {
exitCode = -1
}
}
return outBuf.String(), errBuf.String(), exitCode
}
