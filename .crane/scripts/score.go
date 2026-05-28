//go:build ignore

// score.go -- deletion-grade migration scoring for the APM CLI Python-to-Go migration.
//
// Usage:
//   APM_PYTHON_BIN=/path/to/apm go test -count=1 -json ./... | go run .crane/scripts/score.go
//
// This script implements the deletion-grade framework from issue #96.
// migration_score = 1.0 only when ALL of the following gates pass:
//
//   Gate 1 -- python_reference_required: APM_PYTHON_BIN must be set and valid.
//             TestParityCompletionHardGate must PASS. A missing or invalid Python
//             binary is a hard failure -- never a warning or vacuous pass.
//
//   Gate 2 -- go_tests_pass: every Go test in the module must pass. A single
//             failing non-parity test voids the gate.
//
//   Gate 3 -- help_parity: TestParityCompletionCommandMatrix must pass.
//             Every public command must respond to --help with exit 0.
//
//   Gate 4 -- version_parity: TestParityCompletionVersionEquivalent must pass.
//
//   Gate 5 -- init_parity: TestParityCompletionInitParity must pass.
//             The init command must produce apm.yml equivalent to Python.
//
//   Gate 6 -- error_parity: TestParityCompletionErrorParity must pass.
//             Unknown commands must produce matching non-zero exit codes.
//
//   Gate 7 -- no_known_exceptions: the test output must not contain any
//             "approved exception" log line. Final cutover requires zero exceptions.
//
// If Gate 1 fails, migration_score is forced to 0.0 regardless of other gates.
// Empty or all-skipped test streams also force migration_score to 0.0.
//
// The progress field shows the fraction of deletion-grade gates passing
// (even when migration_score is 0 due to Gate 1 failure).

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TestEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Output  string `json:"Output"`
}

// GateResult tracks the pass/fail state of a single deletion-grade gate.
type GateResult struct {
	Name    string `json:"name"`
	Passing bool   `json:"passing"`
	Reason  string `json:"reason,omitempty"`
}

type Score struct {
	MigrationScore float64      `json:"migration_score"`
	Progress       float64      `json:"progress"`
	ParityPassing  int          `json:"parity_passing"`
	ParityTotal    int          `json:"parity_total"`
	GoTestsTotal   int          `json:"go_tests_total"`
	GoTestsPassing int          `json:"go_tests_passing"`
	Gates          []GateResult `json:"gates"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)

	// Deletion-grade gate test names (exact or prefix).
	const (
		gateHardGate      = "TestParityCompletionHardGate"
		gateCmdMatrix     = "TestParityCompletionCommandMatrix"
		gateVersionParity = "TestParityCompletionVersionEquivalent"
		gateInitParity    = "TestParityCompletionInitParity"
		gateErrorParity   = "TestParityCompletionErrorParity"
	)

	// Track per-test pass/fail.
	testPassed := map[string]bool{}
	testFailed := map[string]bool{}
	var totalTests, passingTests int
	knownExceptionsFound := false
	anyEvents := false

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var ev TestEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}

		anyEvents = true

		// Scan output lines for approved-exception markers.
		// Tests log "APPROVED-EXCEPTION:" via t.Logf; final cutover requires zero.
		if ev.Action == "output" && ev.Output != "" {
			if strings.Contains(ev.Output, "APPROVED-EXCEPTION") {
				knownExceptionsFound = true
			}
		}

		if ev.Test == "" {
			continue
		}

		switch ev.Action {
		case "run":
			totalTests++
		case "pass":
			passingTests++
			testPassed[ev.Test] = true
		case "fail":
			testFailed[ev.Test] = true
		}
	}

	// Gate 1: python_reference_required
	gate1 := GateResult{Name: "python_reference_required"}
	if !anyEvents {
		gate1.Passing = false
		gate1.Reason = "empty test stream -- no test events received"
	} else if testFailed[gateHardGate] {
		gate1.Passing = false
		gate1.Reason = "TestParityCompletionHardGate failed -- APM_PYTHON_BIN missing or invalid"
	} else if testPassed[gateHardGate] {
		gate1.Passing = true
	} else {
		gate1.Passing = false
		gate1.Reason = "TestParityCompletionHardGate not found in test stream"
	}

	// Gate 2: go_tests_pass
	gate2 := GateResult{Name: "go_tests_pass"}
	if totalTests == 0 {
		gate2.Passing = false
		gate2.Reason = "no tests ran"
	} else if passingTests == totalTests {
		gate2.Passing = true
	} else {
		gate2.Passing = false
		gate2.Reason = fmt.Sprintf("%d/%d tests passing", passingTests, totalTests)
	}

	// Gate 3: help_parity (command matrix)
	gate3 := GateResult{Name: "help_parity"}
	if testPassed[gateCmdMatrix] && !testFailed[gateCmdMatrix] {
		gate3.Passing = true
	} else if testFailed[gateCmdMatrix] {
		gate3.Passing = false
		gate3.Reason = "TestParityCompletionCommandMatrix failed"
	} else {
		gate3.Passing = false
		gate3.Reason = "TestParityCompletionCommandMatrix not found"
	}

	// Gate 4: version_parity
	gate4 := GateResult{Name: "version_parity"}
	if testPassed[gateVersionParity] && !testFailed[gateVersionParity] {
		gate4.Passing = true
	} else if testFailed[gateVersionParity] {
		gate4.Passing = false
		gate4.Reason = "TestParityCompletionVersionEquivalent failed"
	} else {
		gate4.Passing = false
		gate4.Reason = "TestParityCompletionVersionEquivalent not found"
	}

	// Gate 5: init_parity
	gate5 := GateResult{Name: "init_parity"}
	if testPassed[gateInitParity] && !testFailed[gateInitParity] {
		gate5.Passing = true
	} else if testFailed[gateInitParity] {
		gate5.Passing = false
		gate5.Reason = "TestParityCompletionInitParity failed"
	} else {
		gate5.Passing = false
		gate5.Reason = "TestParityCompletionInitParity not found"
	}

	// Gate 6: error_parity
	gate6 := GateResult{Name: "error_parity"}
	if testPassed[gateErrorParity] && !testFailed[gateErrorParity] {
		gate6.Passing = true
	} else if testFailed[gateErrorParity] {
		gate6.Passing = false
		gate6.Reason = "TestParityCompletionErrorParity failed"
	} else {
		gate6.Passing = false
		gate6.Reason = "TestParityCompletionErrorParity not found"
	}

	// Gate 7: no_known_exceptions
	gate7 := GateResult{Name: "no_known_exceptions"}
	if knownExceptionsFound {
		gate7.Passing = false
		gate7.Reason = "output contains 'approved exception' -- all exceptions must be resolved for cutover"
	} else {
		gate7.Passing = true
	}

	gates := []GateResult{gate1, gate2, gate3, gate4, gate5, gate6, gate7}

	// Count parity tests (any test with "Parity" in name from cmd/apm).
	parityPassing, parityTotal := 0, 0
	for name, passed := range testPassed {
		if strings.Contains(name, "Parity") {
			parityTotal++
			if passed {
				parityPassing++
			}
		}
	}
	for name := range testFailed {
		if strings.Contains(name, "Parity") && !testPassed[name] {
			parityTotal++
		}
	}

	// Compute migration score.
	gatesPassing := 0
	for _, g := range gates {
		if g.Passing {
			gatesPassing++
		}
	}
	progress := float64(gatesPassing) / float64(len(gates))

	var migrationScore float64
	if !gate1.Passing {
		// Hard gate: Python reference missing forces score to 0.
		migrationScore = 0.0
	} else {
		// All gates must pass for score 1.0; partial credit by gate fraction.
		migrationScore = progress
	}

	score := Score{
		MigrationScore: migrationScore,
		Progress:       progress,
		ParityPassing:  parityPassing,
		ParityTotal:    parityTotal,
		GoTestsTotal:   totalTests,
		GoTestsPassing: passingTests,
		Gates:          gates,
	}

	out, _ := json.MarshalIndent(score, "", "  ")
	fmt.Println(string(out))
}
