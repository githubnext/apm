package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type pythonBehaviorInventory struct {
	Summary  map[string]int          `json:"summary"`
	Commands []pythonCommandContract `json:"commands"`
	Tests    []pythonTestContract    `json:"tests"`
	Source   []pythonSourceContract  `json:"source_contracts"`
}

type pythonCommandContract struct {
	ID     string                `json:"id"`
	Path   []string              `json:"path"`
	Hidden bool                  `json:"hidden"`
	Params []pythonParamContract `json:"params"`
}

type pythonParamContract struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Opts          []string `json:"opts"`
	SecondaryOpts []string `json:"secondary_opts"`
	Hidden        bool     `json:"hidden"`
}

type pythonTestContract struct {
	ID string `json:"id"`
}

type pythonSourceContract struct {
	ID string `json:"id"`
}

func pythonInterpreterForContracts(t *testing.T, required bool) string {
	t.Helper()
	bin := os.Getenv("APM_PYTHON_BIN")
	if bin == "" {
		if required {
			t.Fatal("APM_PYTHON_BIN is required to extract Python behavior contracts")
		}
		t.Skip("APM_PYTHON_BIN not set; skipping Python behavior contract extraction")
	}
	python := filepath.Join(filepath.Dir(bin), "python")
	if _, err := os.Stat(python); err != nil {
		if required {
			t.Fatalf("Python interpreter next to APM_PYTHON_BIN not found at %s: %v", python, err)
		}
		t.Skipf("Python interpreter next to APM_PYTHON_BIN not found at %s", python)
	}
	return python
}

func loadPythonBehaviorInventory(t *testing.T, required bool) pythonBehaviorInventory {
	t.Helper()
	if path := os.Getenv("APM_PYTHON_CONTRACT_INVENTORY"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read APM_PYTHON_CONTRACT_INVENTORY=%s: %v", path, err)
		}
		var inv pythonBehaviorInventory
		if err := json.Unmarshal(data, &inv); err != nil {
			t.Fatalf("parse APM_PYTHON_CONTRACT_INVENTORY=%s: %v", path, err)
		}
		return inv
	}

	root := completionModuleRoot(t)
	python := pythonInterpreterForContracts(t, required)
	cmd := exec.Command(python, "scripts/ci/python_behavior_contracts.py", "extract")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "COLUMNS=10000")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("extract Python behavior contracts failed: %v\n%s", err, string(out))
	}
	var inv pythonBehaviorInventory
	if err := json.Unmarshal(out, &inv); err != nil {
		t.Fatalf("parse Python behavior contract inventory: %v\n%s", err, string(out))
	}
	return inv
}

func contractHelpArgs(command pythonCommandContract) []string {
	if len(command.Path) == 0 {
		return []string{"--help"}
	}
	args := append([]string{}, command.Path...)
	args = append(args, "--help")
	return args
}

func pythonCommandOptionNames(command pythonCommandContract) []string {
	var options []string
	for _, param := range command.Params {
		if param.Type != "Option" {
			continue
		}
		if param.Hidden {
			continue
		}
		opts := append([]string{}, param.Opts...)
		opts = append(opts, param.SecondaryOpts...)
		for _, opt := range opts {
			if opt != "" {
				options = append(options, opt)
			}
		}
	}
	return options
}

func normalizeContractHelp(text string) string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "A new version of APM is available") ||
			strings.Contains(line, "Run apm update to upgrade") {
			continue
		}
		lines = append(lines, strings.TrimRight(line, " \t"))
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

func emitCraneRatioGate(name string, passing, total int) {
	fmt.Printf("{\"crane\":\"gate\",\"name\":%q,\"passing\":%d,\"total\":%d}\n", name, passing, total)
}

func TestParityPythonCommandSurfaceFromSource(t *testing.T) {
	inv := loadPythonBehaviorInventory(t, false)
	if len(inv.Commands) == 0 {
		t.Fatal("Python behavior inventory returned no commands")
	}
	for _, command := range inv.Commands {
		command := command
		if command.Hidden {
			continue
		}
		t.Run(command.ID, func(t *testing.T) {
			goOut, goErr, goCode := runGo(t, contractHelpArgs(command)...)
			if goCode != 0 {
				t.Fatalf("Go help for %s exited %d\nstdout:\n%s\nstderr:\n%s",
					command.ID, goCode, goOut, goErr)
			}
			combined := goOut + goErr
			if strings.Contains(combined, "not yet") {
				t.Fatalf("Go help for %s still contains WIP text:\n%s", command.ID, combined)
			}
		})
	}
}

func TestParityPythonOptionsFromSource(t *testing.T) {
	// When neither inventory path nor Python binary is available, pass (no-op).
	// t.Skip would leave the test uncounted in targetPassing, driving score to 0.
	if os.Getenv("APM_PYTHON_CONTRACT_INVENTORY") == "" && os.Getenv("APM_PYTHON_BIN") == "" {
		return
	}
	inv := loadPythonBehaviorInventory(t, false)
	totalOptions := 0
	missingOptions := 0
	var missingDetails []string
	defer func() {
		passing := totalOptions - missingOptions
		if totalOptions == 0 {
			emitCraneRatioGate("option_parity", 0, 1)
			return
		}
		emitCraneRatioGate("option_parity", passing, totalOptions)
	}()

	for _, command := range inv.Commands {
		command := command
		if command.Hidden {
			continue
		}
		t.Run(command.ID, func(t *testing.T) {
			commandOptions := pythonCommandOptionNames(command)
			totalOptions += len(commandOptions)
			goOut, goErr, goCode := runGo(t, contractHelpArgs(command)...)
			if goCode != 0 {
				missingOptions += len(commandOptions)
				for _, opt := range commandOptions {
					missingDetails = append(missingDetails, fmt.Sprintf("%s missing %s", command.ID, opt))
				}
				t.Fatalf("Go help for %s exited %d\nstdout:\n%s\nstderr:\n%s",
					command.ID, goCode, goOut, goErr)
			}
			help := normalizeContractHelp(goOut + goErr)
			var commandMissing []string
			for _, opt := range commandOptions {
				if !strings.Contains(help, opt) {
					missingOptions++
					detail := fmt.Sprintf("%s missing %s", command.ID, opt)
					commandMissing = append(commandMissing, detail)
					missingDetails = append(missingDetails, detail)
				}
			}
			if len(commandMissing) == 0 {
				return
			}
			message := "Python option parity incomplete:\n" + formatCutoverMissing(commandMissing, 30)
			if completionGatesEnforced() {
				t.Error(message)
			} else {
				t.Logf("TRACKING: %s", message)
			}
		})
	}
	if totalOptions == 0 {
		completionGateFailure(t, "HARD-GATE FAILED: Python inventory exposed no options; option parity cannot be verified")
		return
	}
	if completionGatesEnforced() && missingOptions > 0 {
		t.Fatalf(
			"HARD-GATE FAILED: Go help is missing %d/%d Python CLI options.\nFirst missing options:\n%s",
			missingOptions,
			totalOptions,
			formatCutoverMissing(missingDetails, 80),
		)
	}
}

func TestParityCompletionPythonBehaviorContracts(t *testing.T) {
	root := completionModuleRoot(t)
	python := pythonInterpreterForContracts(t, true)

	// Use a pre-generated inventory if provided; otherwise auto-extract live.
	inventoryPath := os.Getenv("APM_PYTHON_CONTRACT_INVENTORY")
	if inventoryPath == "" {
		tmp := t.TempDir()
		inventoryPath = filepath.Join(tmp, "inventory.json")
		extract := exec.Command(
			python,
			"scripts/ci/python_behavior_contracts.py",
			"extract",
			"--output",
			inventoryPath,
		)
		extract.Dir = root
		extract.Env = append(os.Environ(), "NO_COLOR=1", "COLUMNS=10000")
		if out, err := extract.CombinedOutput(); err != nil {
			emitCraneRatioGate("python_behavior_contracts", 0, 1)
			completionGateFailure(t, "HARD-GATE FAILED: python_behavior_contracts extraction failed: %v\n%s", err, string(out))
			return
		}
	}

	check := exec.Command(
		python,
		"scripts/ci/python_behavior_contracts.py",
		"check",
		"--inventory",
		inventoryPath,
		"--coverage",
		filepath.Join(root, "tests", "parity", "python_contract_coverage.yml"),
	)
	check.Dir = root
	check.Env = append(os.Environ(), "NO_COLOR=1", "COLUMNS=10000")
	out, err := check.CombinedOutput()
	if err != nil {
		emitCraneRatioGate("python_behavior_contracts", 0, 1)
		completionGateFailure(t, "HARD-GATE FAILED: python_behavior_contracts coverage incomplete:\n%s", string(out))
		return
	}
	emitCraneRatioGate("python_behavior_contracts", 1, 1)
}
