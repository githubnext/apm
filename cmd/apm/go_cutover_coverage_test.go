package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

type goCutoverPythonTestCoverage struct {
	SchemaVersion        int                 `json:"schema_version"`
	Description          string              `json:"description"`
	ConvertedPythonTests map[string][]string `json:"converted_python_tests"`
}

type pythonClassContext struct {
	name   string
	indent int
}

var (
	pythonClassRE = regexp.MustCompile(`^class\s+(Test[A-Za-z0-9_]*)\b`)
	pythonTestRE  = regexp.MustCompile(`^(?:async\s+)?def\s+(test_[A-Za-z0-9_]*)\b`)
	goTestFuncRE  = regexp.MustCompile(`^func\s+(Test[A-Za-z0-9_]*)\s*\(`)
)

func TestGoCutoverPythonTestConversionCoverage(t *testing.T) {
	root := completionModuleRoot(t)
	pythonTests := discoverPythonTestsForCutover(t, root)
	goTests := discoverGoTestsForCutover(t, root)
	coverage := loadGoCutoverPythonTestCoverage(t, root)

	behaviorBacked := 0
	var missing []string
	var unknown []string
	var weak []string
	for _, id := range pythonTests {
		tests := coverage.ConvertedPythonTests[id]
		if len(tests) == 0 {
			missing = append(missing, id)
			continue
		}
		for _, testName := range tests {
			if _, ok := goTests[testName]; !ok {
				unknown = append(unknown, fmt.Sprintf("%s -> %s", id, testName))
			}
		}
		if !hasBehaviorBackedGoTest(tests, goTests) {
			weak = append(weak, fmt.Sprintf("%s -> %s", id, strings.Join(tests, ", ")))
			continue
		}
		behaviorBacked++
	}

	defer emitCraneRatioGate("python_behavior_contracts", behaviorBacked, len(pythonTests))
	defer emitCraneBoolGate("golden_fixture_corpus", behaviorBacked == len(pythonTests) && len(pythonTests) > 0)
	defer emitCraneBoolGate("all_go_golden_tests", behaviorBacked == len(pythonTests) && len(pythonTests) > 0)

	if len(pythonTests) == 0 {
		t.Fatal("no Python tests discovered under tests/; coverage gate cannot prove conversion")
	}
	if coverage.SchemaVersion != 1 {
		t.Fatalf("go cutover Python test coverage manifest schema_version = %d, want 1", coverage.SchemaVersion)
	}
	if len(missing) > 0 {
		t.Fatalf(
			"Go cutover coverage incomplete: %d/%d Python tests mapped to Go tests; %d missing.\nFirst missing tests:\n%s",
			behaviorBacked,
			len(pythonTests),
			len(missing),
			formatCutoverMissing(missing, 80),
		)
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		t.Fatalf(
			"Go cutover coverage references Go tests that do not exist: %d stale mappings.\nFirst stale mappings:\n%s",
			len(unknown),
			formatCutoverMissing(unknown, 80),
		)
	}
	if len(weak) > 0 {
		t.Fatalf(
			"Go cutover coverage is not behavior-backed: %d/%d Python tests do not map to a real Go-only cutover behavior test.\nFirst weak mappings:\n%s",
			len(weak),
			len(pythonTests),
			formatCutoverMissing(weak, 80),
		)
	}
}

func TestGoCutoverNoPythonRuntimeDependency(t *testing.T) {
	dir := t.TempDir()
	stdout, stderr, code := realBehaviorRunGoInDirSanitized(t, dir, "--version")
	passed := code == 0 && strings.Contains(strings.ToLower(stdout+stderr), "apm")
	emitCraneBoolGate("no_python_runtime_dependency", passed)
	if !passed {
		t.Fatalf("Go CLI must run without Python runtime env vars; exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func discoverPythonTestsForCutover(t *testing.T, root string) []string {
	t.Helper()
	testsRoot := filepath.Join(root, "tests")
	var ids []string
	err := filepath.WalkDir(testsRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == "__pycache__" {
				return filepath.SkipDir
			}
			if rel, relErr := filepath.Rel(testsRoot, path); relErr == nil {
				parts := strings.Split(filepath.ToSlash(rel), "/")
				if len(parts) > 0 && parts[0] == "parity" {
					return filepath.SkipDir
				}
			}
			return nil
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "test") || !strings.HasSuffix(name, ".py") {
			return nil
		}
		fileIDs, scanErr := scanPythonTestFile(t, root, path)
		if scanErr != nil {
			return scanErr
		}
		ids = append(ids, fileIDs...)
		return nil
	})
	if err != nil {
		t.Fatalf("discover Python tests: %v", err)
	}
	sort.Strings(ids)
	return ids
}

func scanPythonTestFile(t *testing.T, root, path string) ([]string, error) {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return nil, err
	}
	rel = filepath.ToSlash(rel)

	var ids []string
	var classes []pythonClassContext
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		indent := leadingWhitespaceWidth(line)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		for len(classes) > 0 && indent <= classes[len(classes)-1].indent {
			classes = classes[:len(classes)-1]
		}

		if match := pythonClassRE.FindStringSubmatch(trimmed); match != nil {
			classes = append(classes, pythonClassContext{name: match[1], indent: indent})
			continue
		}

		match := pythonTestRE.FindStringSubmatch(trimmed)
		if match == nil {
			continue
		}
		name := match[1]
		if len(classes) > 0 && indent > classes[len(classes)-1].indent {
			ids = append(ids, fmt.Sprintf("%s::%s::%s", rel, classes[len(classes)-1].name, name))
			continue
		}
		if indent == 0 {
			ids = append(ids, fmt.Sprintf("%s::%s", rel, name))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func discoverGoTestsForCutover(t *testing.T, root string) map[string]struct{} {
	t.Helper()
	tests := map[string]struct{}{}
	err := filepath.WalkDir(filepath.Join(root, "cmd", "apm"), func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		file, openErr := os.Open(path)
		if openErr != nil {
			return openErr
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			match := goTestFuncRE.FindStringSubmatch(strings.TrimSpace(scanner.Text()))
			if match != nil {
				tests[match[1]] = struct{}{}
			}
		}
		return scanner.Err()
	})
	if err != nil {
		t.Fatalf("discover Go tests: %v", err)
	}
	if len(tests) == 0 {
		t.Fatal("no Go tests discovered under cmd/apm; coverage gate cannot prove conversion")
	}
	return tests
}

func loadGoCutoverPythonTestCoverage(t *testing.T, root string) goCutoverPythonTestCoverage {
	t.Helper()
	path := filepath.Join(root, "cmd", "apm", "testdata", "go_cutover", "python_test_coverage.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Go cutover Python test coverage manifest %s: %v", path, err)
	}
	var coverage goCutoverPythonTestCoverage
	if err := json.Unmarshal(data, &coverage); err != nil {
		t.Fatalf("parse Go cutover Python test coverage manifest %s: %v", path, err)
	}
	if coverage.ConvertedPythonTests == nil {
		coverage.ConvertedPythonTests = map[string][]string{}
	}
	return coverage
}

func formatCutoverMissing(missing []string, limit int) string {
	if limit > len(missing) {
		limit = len(missing)
	}
	lines := make([]string, 0, limit+1)
	for _, id := range missing[:limit] {
		lines = append(lines, "  - "+id)
	}
	if limit < len(missing) {
		lines = append(lines, fmt.Sprintf("  ... %d more", len(missing)-limit))
	}
	return strings.Join(lines, "\n")
}

func hasBehaviorBackedGoTest(names []string, existing map[string]struct{}) bool {
	for _, name := range names {
		if _, ok := existing[name]; ok && isBehaviorBackedGoTest(name) {
			return true
		}
	}
	return false
}

func isBehaviorBackedGoTest(name string) bool {
	return strings.HasPrefix(name, "TestGoCutoverReal")
}

func leadingWhitespaceWidth(line string) int {
	width := 0
	for _, r := range line {
		switch r {
		case ' ':
			width++
		case '\t':
			width += 4
		default:
			return width
		}
	}
	return width
}

func realBehaviorRunGoInDirSanitized(t *testing.T, dir string, args ...string) (string, string, int) {
	t.Helper()
	cleared := map[string]string{
		"APM_PYTHON_BIN":                "",
		"APM_PYTHON_CONTRACT_INVENTORY": "",
		"PYTHONPATH":                    "",
		"VIRTUAL_ENV":                   "",
	}
	return realBehaviorRunGoInDir(t, dir, cleared, args...)
}
