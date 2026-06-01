package main

import (
	"encoding/json"
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
	if os.Getenv("APM_PYTHON_CONTRACT_INVENTORY") == "" {
		t.Skip("set APM_PYTHON_CONTRACT_INVENTORY to run option-coverage checks (migration CI only)")
	}
	inv := loadPythonBehaviorInventory(t, false)
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
			help := normalizeContractHelp(goOut + goErr)
			for _, param := range command.Params {
				if param.Type != "Option" {
					continue
				}
				opts := append([]string{}, param.Opts...)
				opts = append(opts, param.SecondaryOpts...)
				for _, opt := range opts {
					if opt == "" {
						continue
					}
					if !strings.Contains(help, opt) {
						t.Logf("TRACKING: %s help missing Python option %s (migration in progress)", command.ID, opt)
					}
				}
			}
		})
	}
}

func TestParityCompletionPythonBehaviorContracts(t *testing.T) {
	inventoryPath := os.Getenv("APM_PYTHON_CONTRACT_INVENTORY")
	if inventoryPath == "" {
		t.Skip("set APM_PYTHON_CONTRACT_INVENTORY to enforce the behavior-contracts coverage gate (migration CI only)")
	}

	root := completionModuleRoot(t)
	python := pythonInterpreterForContracts(t, true)

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
		t.Fatalf("Python behavior contracts are not fully covered:\n%s", string(out))
	}
}
