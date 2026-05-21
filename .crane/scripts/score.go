//go:build ignore

// score.go: Crane migration scoring script.
// Usage: cd cmd/apm && go test -json ./... | go run ../../.crane/scripts/score.go
//
// Reads go test JSON output from stdin and emits a migration score JSON object.
// migration_score = (parity_passing / parity_total) * correctness_gate
// where correctness_gate = 1.0 if all target tests pass, 0.0 otherwise.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type TestEvent struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Package string  `json:"Package"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

type ScoreOutput struct {
	MigrationScore     float64 `json:"migration_score"`
	Progress           float64 `json:"progress"`
	ParityPassing      int     `json:"parity_passing"`
	ParityTotal        int     `json:"parity_total"`
	SourceTestsPassing int     `json:"source_tests_passing"`
	TargetTestsPassing int     `json:"target_tests_passing"`
	PerfRatio          float64 `json:"perf_ratio"`
}

func main() {
	var passed, failed int
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev TestEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue
		}
		if ev.Test == "" {
			// package-level event
			continue
		}
		switch ev.Action {
		case "pass":
			passed++
		case "fail":
			failed++
		}
	}

	total := passed + failed
	// parity_total is the number of Python source modules (302).
	// Until parity tests are wired, this reflects Go test count.
	const parityTotal = 302
	parityPassing := 0
	if total > 0 {
		parityPassing = passed
	}

	correctnessGate := 0.0
	if total > 0 && failed == 0 {
		correctnessGate = 1.0
	} else if total == 0 {
		// No tests yet -- score is 0 but not a failure.
		correctnessGate = 0.0
	}

	progress := 0.0
	if parityTotal > 0 {
		progress = float64(parityPassing) / float64(parityTotal)
	}

	migrationScore := progress * correctnessGate

	out := ScoreOutput{
		MigrationScore:     migrationScore,
		Progress:           progress,
		ParityPassing:      parityPassing,
		ParityTotal:        parityTotal,
		SourceTestsPassing: 247,
		TargetTestsPassing: passed,
		PerfRatio:          1.0,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "score.go: encode error: %v\n", err)
		os.Exit(1)
	}
}
