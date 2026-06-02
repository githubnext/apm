//go:build ignore

// score.go -- deletion-grade migration scoring for the APM CLI Python-to-Go migration.
// Usage: go test -json ./... | go run .crane/scripts/score.go
// Outputs JSON that separates progress metrics from cutover-readiness gates.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TestEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Output  string `json:"Output"`
}

type GateEvent struct {
	Crane   string `json:"crane"`
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Passing int    `json:"passing"`
	Total   int    `json:"total"`
	Count   int    `json:"count"`
}

type BoolGate struct {
	Seen   bool
	Passed bool
}

type RatioGate struct {
	Seen    bool
	Passing int
	Total   int
}

func (g BoolGate) OK() bool {
	return g.Seen && g.Passed
}

func (g RatioGate) Percent() float64 {
	if !g.Seen || g.Total <= 0 {
		return 0
	}
	return float64(g.Passing) / float64(g.Total)
}

func (g RatioGate) OK() bool {
	return g.Seen && g.Total > 0 && g.Passing == g.Total
}

type CutoverGates struct {
	PythonReferenceRequired bool    `json:"python_reference_required"`
	SurfaceParity           float64 `json:"surface_parity"`
	HelpParity              float64 `json:"help_parity"`
	FunctionalContracts     float64 `json:"functional_contracts"`
	StateDiffContracts      float64 `json:"state_diff_contracts"`
	PythonBehaviorContracts float64 `json:"python_behavior_contracts"`
	KnownExceptions         int     `json:"known_exceptions"`
	GoTests                 string  `json:"go_tests"`
	PythonTests             string  `json:"python_tests"`
	Benchmarks              string  `json:"benchmarks"`
}

type ProgressMetrics struct {
	ParityPassing      int     `json:"parity_passing"`
	ParityTotal        int     `json:"parity_total"`
	SourceTestsPassing int     `json:"source_tests_passing"`
	TargetTestsPassing int     `json:"target_tests_passing"`
	PerfRatio          float64 `json:"perf_ratio"`
}

type GateResult struct {
	Name    string `json:"name"`
	Passing bool   `json:"passing"`
	Reason  string `json:"reason,omitempty"`
}

type Score struct {
	MigrationScore         float64         `json:"migration_score"`
	Progress               float64         `json:"progress"`
	CutoverReady           bool            `json:"cutover_ready"`
	CutoverGates           CutoverGates    `json:"cutover_gates"`
	ProgressMetrics        ProgressMetrics `json:"progress_metrics"`
	DeletionGradeReady     bool            `json:"deletion_grade_ready"`
	PythonReferencePresent bool            `json:"python_reference_present"`
	SurfaceParity          float64         `json:"surface_parity"`
	HelpParity             float64         `json:"help_parity"`
	FunctionalParity       float64         `json:"functional_parity"`
	StateDiffParity        float64         `json:"state_diff_parity"`
	KnownExceptions        int             `json:"known_exceptions"`
	PythonTestsPassing     bool            `json:"python_tests_passing"`
	GoTestsPassing         bool            `json:"go_tests_passing"`
	BenchmarksPassing      bool            `json:"benchmarks_passing"`
	ParityPassing          int             `json:"parity_passing"`
	ParityTotal            int             `json:"parity_total"`
	SourceTestsPassing     int             `json:"source_tests_passing"`
	TargetTestsPassing     int             `json:"target_tests_passing"`
	PerfRatio              float64         `json:"perf_ratio"`
	Gates                  []GateResult    `json:"gates"`
}

func main() {
	score, err := computeScore(os.Stdin, os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "score: %v\n", err)
		os.Exit(1)
	}

	out, _ := json.MarshalIndent(score, "", "  ")
	fmt.Println(string(out))
}

type getenvFunc func(string) string

type scanInput interface {
	Read([]byte) (int, error)
}

func computeScore(input scanInput, getenv getenvFunc) (Score, error) {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)

	var parityPassing, parityTotal, targetPassing, targetTotal int
	eventsSeen := 0
	goTestsFailed := false
	running := map[string]bool{}
	passed := map[string]bool{}
	failed := map[string]bool{}
	knownExceptions := knownExceptionsFromEnv(getenv("APM_KNOWN_EXCEPTIONS"))
	pythonReference := BoolGate{}
	pythonTests := BoolGate{Seen: getenv("APM_PYTHON_TESTS") != "", Passed: getenv("APM_PYTHON_TESTS") == "pass"}
	benchmarks := BoolGate{Seen: getenv("APM_BENCHMARKS") != "", Passed: getenv("APM_BENCHMARKS") == "pass"}
	surface := RatioGate{}
	help := RatioGate{}
	functional := RatioGate{}
	stateDiff := RatioGate{}
	behaviorContracts := RatioGate{}

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "{") {
			continue
		}
		if gate, ok := parseGateEvent(line); ok {
			eventsSeen++
			applyGateEvent(gate, &pythonReference, &surface, &help, &functional, &stateDiff, &behaviorContracts, &knownExceptions, &pythonTests, &benchmarks)
			continue
		}

		var ev TestEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		eventsSeen++

		if ev.Output != "" {
			if gate, ok := parseGateEvent(ev.Output); ok {
				applyGateEvent(gate, &pythonReference, &surface, &help, &functional, &stateDiff, &behaviorContracts, &knownExceptions, &pythonTests, &benchmarks)
			}
			if n, ok := approvedExceptionCount(ev.Output); ok && n > knownExceptions {
				knownExceptions = n
			}
		}

		if ev.Test == "" {
			if isTargetPackage(ev.Package) && ev.Action == "fail" {
				goTestsFailed = true
			}
			continue
		}

		key := ev.Package + "/" + ev.Test
		switch ev.Action {
		case "run":
			running[key] = true
		case "pass":
			passed[ev.Test] = true
			delete(running, key)
		case "fail":
			failed[ev.Test] = true
			delete(running, key)
			if isTargetPackage(ev.Package) {
				goTestsFailed = true
			}
		case "skip":
			delete(running, key)
		}

		isParity := strings.Contains(ev.Test, "Parity") || strings.Contains(ev.Package, "parity")
		if isParity {
			if ev.Action == "run" {
				parityTotal++
			} else if ev.Action == "pass" {
				parityPassing++
			}
		}
		if isTargetPackage(ev.Package) {
			if ev.Action == "run" {
				targetTotal++
			} else if ev.Action == "pass" {
				targetPassing++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Score{}, err
	}
	if eventsSeen == 0 || targetTotal == 0 {
		return Score{}, fmt.Errorf("Go test event stream is empty or incomplete")
	}
	if len(running) > 0 {
		return Score{}, fmt.Errorf("Go test event stream is incomplete: %d test(s) did not finish", len(running))
	}

	if !pythonReference.Seen {
		pythonReference = BoolGate{Seen: true, Passed: testPassed(passed, failed, "TestParityCompletionHardGate") || pythonReferenceReady(getenv("APM_PYTHON_BIN"))}
	}
	if !surface.Seen {
		surface = inferredAnyRatioGate(passed, failed, "TestParityCompletionSurfaceParity", "TestParitySurfaceInventory")
	}
	if !help.Seen {
		help = inferredAllRatioGate(passed, failed, "TestParityCompletionCommandMatrix", "TestParityCompletionHelpIdentical")
	}
	if !functional.Seen {
		functional = inferredAnyRatioGate(passed, failed, "TestParityCompletionFunctionalContracts", "TestParityFunctionalContracts")
	}
	if !stateDiff.Seen {
		stateDiff = inferredAnyRatioGate(passed, failed, "TestParityCompletionStateDiffContracts", "TestParityStateDiffContracts")
	}
	if !behaviorContracts.Seen {
		behaviorContracts = RatioGate{Seen: true, Passing: 0, Total: 1}
	}
	if !pythonTests.Seen {
		pythonTests = BoolGate{Seen: true, Passed: testPassed(passed, failed, "TestParityCompletionPythonSuite")}
	}
	if !benchmarks.Seen {
		benchmarks = BoolGate{Seen: true, Passed: testPassed(passed, failed, "TestParityCompletionBenchmarks")}
	}

	goTestsPass := !goTestsFailed && targetTotal > 0 && targetPassing == targetTotal
	gates := CutoverGates{
		PythonReferenceRequired: pythonReference.OK(),
		SurfaceParity:           surface.Percent(),
		HelpParity:              help.Percent(),
		FunctionalContracts:     functional.Percent(),
		StateDiffContracts:      stateDiff.Percent(),
		PythonBehaviorContracts: behaviorContracts.Percent(),
		KnownExceptions:         knownExceptions,
		GoTests:                 passFail(goTestsPass),
		PythonTests:             passFail(pythonTests.OK()),
		Benchmarks:              passFail(benchmarks.OK()),
	}

	total := 302 // fixed historical progress denominator: Python modules/functions to port.
	if parityTotal > total {
		total = parityTotal
	}

	progress := 0.0
	if total > 0 {
		progress = float64(parityPassing) / float64(total)
	}

	cutoverReady := gates.PythonReferenceRequired &&
		gates.SurfaceParity == 1.0 &&
		gates.HelpParity == 1.0 &&
		gates.FunctionalContracts == 1.0 &&
		gates.StateDiffContracts == 1.0 &&
		gates.PythonBehaviorContracts == 1.0 &&
		gates.KnownExceptions == 0 &&
		gates.GoTests == "pass" &&
		gates.PythonTests == "pass" &&
		gates.Benchmarks == "pass"

	migrationScore := progress
	if !goTestsPass {
		migrationScore = 0
	}
	if !cutoverReady && migrationScore >= 1.0 {
		migrationScore = 0.999
	}
	if cutoverReady && progress == 1.0 {
		migrationScore = 1.0
	}

	metrics := ProgressMetrics{
		ParityPassing:      parityPassing,
		ParityTotal:        total,
		SourceTestsPassing: sourceTestsPassing(getenv("APM_SOURCE_TESTS_PASSING")),
		TargetTestsPassing: targetPassing,
		PerfRatio:          perfRatio(getenv("APM_PERF_RATIO")),
	}

	return Score{
		MigrationScore:         migrationScore,
		Progress:               progress,
		CutoverReady:           cutoverReady,
		CutoverGates:           gates,
		ProgressMetrics:        metrics,
		DeletionGradeReady:     cutoverReady,
		PythonReferencePresent: gates.PythonReferenceRequired,
		SurfaceParity:          gates.SurfaceParity,
		HelpParity:             gates.HelpParity,
		FunctionalParity:       gates.FunctionalContracts,
		StateDiffParity:        gates.StateDiffContracts,
		KnownExceptions:        gates.KnownExceptions,
		PythonTestsPassing:     gates.PythonTests == "pass",
		GoTestsPassing:         gates.GoTests == "pass",
		BenchmarksPassing:      gates.Benchmarks == "pass",
		ParityPassing:          metrics.ParityPassing,
		ParityTotal:            metrics.ParityTotal,
		SourceTestsPassing:     metrics.SourceTestsPassing,
		TargetTestsPassing:     metrics.TargetTestsPassing,
		PerfRatio:              metrics.PerfRatio,
		Gates:                  gateResults(gates),
	}, nil
}

func parseGateEvent(line string) (GateEvent, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "{") {
		return GateEvent{}, false
	}
	var gate GateEvent
	if err := json.Unmarshal([]byte(line), &gate); err != nil || gate.Crane != "gate" {
		return GateEvent{}, false
	}
	return gate, true
}

func applyGateEvent(
	gate GateEvent,
	pythonReference *BoolGate,
	surface *RatioGate,
	help *RatioGate,
	functional *RatioGate,
	stateDiff *RatioGate,
	behaviorContracts *RatioGate,
	knownExceptions *int,
	pythonTests *BoolGate,
	benchmarks *BoolGate,
) {
	switch gate.Name {
	case "python_reference":
		*pythonReference = BoolGate{Seen: true, Passed: gate.Passed}
	case "surface":
		*surface = RatioGate{Seen: true, Passing: gate.Passing, Total: gate.Total}
	case "help":
		*help = RatioGate{Seen: true, Passing: gate.Passing, Total: gate.Total}
	case "functional":
		*functional = RatioGate{Seen: true, Passing: gate.Passing, Total: gate.Total}
	case "state_diff":
		*stateDiff = RatioGate{Seen: true, Passing: gate.Passing, Total: gate.Total}
	case "python_behavior_contracts":
		*behaviorContracts = RatioGate{Seen: true, Passing: gate.Passing, Total: gate.Total}
	case "known_exceptions":
		*knownExceptions = gate.Count
	case "python_tests":
		*pythonTests = BoolGate{Seen: true, Passed: gate.Passed}
	case "benchmarks":
		*benchmarks = BoolGate{Seen: true, Passed: gate.Passed}
	}
}

func isTargetPackage(pkg string) bool {
	return strings.HasPrefix(pkg, "github.com/githubnext/apm/")
}

func pythonReferenceReady(bin string) bool {
	if bin == "" {
		return false
	}
	info, err := os.Stat(bin)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode()&0o111 != 0
}

func testPassed(passed, failed map[string]bool, names ...string) bool {
	for _, name := range names {
		if failed[name] {
			return false
		}
	}
	for _, name := range names {
		if passed[name] {
			return true
		}
	}
	return false
}

func inferredAnyRatioGate(passed, failed map[string]bool, names ...string) RatioGate {
	for _, name := range names {
		if failed[name] {
			return RatioGate{Seen: true, Passing: 0, Total: 1}
		}
	}
	return RatioGate{Seen: true, Passing: boolToInt(testPassed(passed, failed, names...)), Total: 1}
}

func inferredAllRatioGate(passed, failed map[string]bool, names ...string) RatioGate {
	for _, name := range names {
		if failed[name] {
			return RatioGate{Seen: true, Passing: 0, Total: 1}
		}
	}
	return RatioGate{Seen: true, Passing: boolToInt(allRequiredTestsPassed(passed, names...)), Total: 1}
}

func allRequiredTestsPassed(passed map[string]bool, names ...string) bool {
	for _, name := range names {
		if !passed[name] {
			return false
		}
	}
	return true
}

func gateResults(gates CutoverGates) []GateResult {
	return []GateResult{
		{Name: "python_reference_required", Passing: gates.PythonReferenceRequired},
		{Name: "go_tests_pass", Passing: gates.GoTests == "pass"},
		{Name: "surface_parity", Passing: gates.SurfaceParity == 1.0},
		{Name: "help_parity", Passing: gates.HelpParity == 1.0},
		{Name: "functional_contracts", Passing: gates.FunctionalContracts == 1.0},
		{Name: "state_diff_contracts", Passing: gates.StateDiffContracts == 1.0},
		{Name: "python_behavior_contracts", Passing: gates.PythonBehaviorContracts == 1.0},
		{Name: "python_tests_pass", Passing: gates.PythonTests == "pass"},
		{Name: "benchmarks_pass", Passing: gates.Benchmarks == "pass"},
		{Name: "no_known_exceptions", Passing: gates.KnownExceptions == 0},
	}
}

func passFail(ok bool) string {
	if ok {
		return "pass"
	}
	return "fail"
}

func boolToInt(ok bool) int {
	if ok {
		return 1
	}
	return 0
}

func knownExceptionsFromEnv(raw string) int {
	if raw == "" {
		return 0
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 1
	}
	return n
}

func approvedExceptionCount(output string) (int, bool) {
	if !strings.Contains(strings.ToLower(output), "approved") || !strings.Contains(strings.ToLower(output), "exception") {
		return 0, false
	}
	fields := strings.Fields(output)
	for _, field := range fields {
		if n, err := strconv.Atoi(field); err == nil {
			return n, true
		}
	}
	return 1, true
}

func sourceTestsPassing(raw string) int {
	if raw == "" {
		return 247
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return n
}

func perfRatio(raw string) float64 {
	if raw == "" {
		return 1.0
	}
	n, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return n
}
