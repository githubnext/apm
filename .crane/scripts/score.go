//go:build ignore

// score.go -- migration scoring script for the APM CLI Python-to-Go migration.
// Usage: go test -json ./... | go run .crane/scripts/score.go
// Outputs JSON with migration_score and progress metrics.
//
// Scoring formula:
//   migration_score = (parity_passing / parity_total) * correctness_gate
//   correctness_gate = 1.0 if all target tests pass, 0.0 otherwise
//
// NOTE: This script must NOT be modified after milestone 1 is accepted.

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

type Score struct {
	MigrationScore     float64 `json:"migration_score"`
	Progress           float64 `json:"progress"`
	ParityPassing      int     `json:"parity_passing"`
	ParityTotal        int     `json:"parity_total"`
	SourceTestsPassing int     `json:"source_tests_passing"`
	TargetTestsPassing int     `json:"target_tests_passing"`
	PerfRatio          float64 `json:"perf_ratio"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var parityPassing, parityTotal, targetPassing, targetTotal int

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var ev TestEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		if ev.Test == "" {
			continue
		}

		isParity := strings.Contains(ev.Test, "Parity") || strings.Contains(ev.Package, "parity")
		isTarget := strings.HasPrefix(ev.Package, "github.com/githubnext/apm/")

		if isParity {
			if ev.Action == "run" {
				parityTotal++
			} else if ev.Action == "pass" {
				parityPassing++
			}
		}
		if isTarget {
			if ev.Action == "run" {
				targetTotal++
			} else if ev.Action == "pass" {
				targetPassing++
			}
		}
	}

	correctnessGate := 1.0
	if targetTotal > 0 && targetPassing < targetTotal {
		correctnessGate = 0.0
	}

	total := 302 // fixed: total Python modules/functions to port
	if parityTotal > total {
		total = parityTotal
	}

	var migrationScore float64
	if total > 0 {
		migrationScore = (float64(parityPassing) / float64(total)) * correctnessGate
	}

	progress := float64(parityPassing) / float64(total)

	score := Score{
		MigrationScore:     migrationScore,
		Progress:           progress,
		ParityPassing:      parityPassing,
		ParityTotal:        total,
		SourceTestsPassing: 247, // stable Python baseline
		TargetTestsPassing: targetPassing,
		PerfRatio:          1.0,
	}

	out, _ := json.MarshalIndent(score, "", "  ")
	fmt.Println(string(out))
}
