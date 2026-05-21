package listcmd

import (
	"strings"
	"testing"
)

func TestScript_NameOnly_Extra3(t *testing.T) {
	s := Script{Name: "build"}
	if s.Name != "build" {
		t.Errorf("Name = %q, want build", s.Name)
	}
	if s.Command != "" {
		t.Errorf("Command should be empty, got %q", s.Command)
	}
}

func TestParseScripts_NestedCommand_Extra3(t *testing.T) {
	content := `
scripts:
  deploy: echo deploying
`
	result := parseScripts(content)
	if result["deploy"] != "echo deploying" {
		t.Errorf("deploy = %q, want echo deploying", result["deploy"])
	}
}

func TestParseScripts_MultipleScriptKeys_Extra3(t *testing.T) {
	content := `
scripts:
  lint: ruff check .
  fmt: ruff format .
  test: pytest
`
	result := parseScripts(content)
	if len(result) != 3 {
		t.Errorf("len = %d, want 3", len(result))
	}
	if result["lint"] != "ruff check ." {
		t.Errorf("lint = %q", result["lint"])
	}
	if result["test"] != "pytest" {
		t.Errorf("test = %q", result["test"])
	}
}

func TestParseScripts_HyphenatedKey_Extra3(t *testing.T) {
	content := `
scripts:
  check-types: mypy src/
`
	result := parseScripts(content)
	if result["check-types"] != "mypy src/" {
		t.Errorf("check-types = %q", result["check-types"])
	}
}

func TestParseScripts_UnderscoreKey_Extra3(t *testing.T) {
	content := `
scripts:
  run_all: make all
`
	result := parseScripts(content)
	if result["run_all"] != "make all" {
		t.Errorf("run_all = %q", result["run_all"])
	}
}

func TestParseScripts_CommandWithEqualSign_Extra3(t *testing.T) {
	content := `
scripts:
  env: export FOO=bar
`
	result := parseScripts(content)
	if !strings.Contains(result["env"], "FOO") {
		t.Errorf("env command = %q should contain FOO", result["env"])
	}
}

func TestParseScripts_LeadingSpaces_Extra3(t *testing.T) {
	content := "scripts:\n  build:   go build ./...\n"
	result := parseScripts(content)
	v := result["build"]
	if v == "" {
		t.Error("build should not be empty")
	}
	if !strings.Contains(v, "go build") {
		t.Errorf("build command %q should contain go build", v)
	}
}

func TestParseScripts_EmptyScriptsSection_Extra3(t *testing.T) {
	content := "scripts:\ndependencies:\n  foo: bar\n"
	result := parseScripts(content)
	if len(result) != 0 {
		t.Errorf("empty scripts section should give empty map, got %v", result)
	}
}

func TestScript_CommandVariants_Extra3(t *testing.T) {
	scripts := []Script{
		{Name: "a", Command: "go test ./..."},
		{Name: "b", Command: "npm run build"},
		{Name: "c", Command: "python -m pytest"},
	}
	for _, s := range scripts {
		if s.Name == "" {
			t.Error("Name should not be empty")
		}
		if s.Command == "" {
			t.Error("Command should not be empty")
		}
	}
}

func TestParseScripts_ReturnIsMap_Extra3(t *testing.T) {
	content := "scripts:\n  x: y\n"
	result := parseScripts(content)
	if _, ok := result["x"]; !ok {
		t.Error("expected key x in result map")
	}
}
