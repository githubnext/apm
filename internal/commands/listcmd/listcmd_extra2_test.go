package listcmd

import (
	"testing"
)

func TestScript_NameAndCommand(t *testing.T) {
	s := Script{Name: "build", Command: "go build ./..."}
	if s.Name != "build" {
		t.Errorf("unexpected Name: %q", s.Name)
	}
	if s.Command != "go build ./..." {
		t.Errorf("unexpected Command: %q", s.Command)
	}
}

func TestScript_ZeroValue(t *testing.T) {
	var s Script
	if s.Name != "" || s.Command != "" {
		t.Error("zero-value Script should have empty Name and Command")
	}
}

func TestParseScripts_NoScriptsSectionDeps(t *testing.T) {
	content := "dependencies:\n  - owner/pkg\n"
	result := parseScripts(content)
	if len(result) != 0 {
		t.Errorf("expected no scripts, got %v", result)
	}
}

func TestParseScripts_SingleScriptEntry(t *testing.T) {
	content := "scripts:\n  test: go test ./...\n"
	result := parseScripts(content)
	if result["test"] != "go test ./..." {
		t.Errorf("unexpected test command: %q", result["test"])
	}
}

func TestParseScripts_MultipleScriptEntries(t *testing.T) {
	content := "scripts:\n  build: go build\n  lint: golint\n  test: go test\n"
	result := parseScripts(content)
	if len(result) != 3 {
		t.Errorf("expected 3 scripts, got %d: %v", len(result), result)
	}
}

func TestParseScripts_EmptyYAML(t *testing.T) {
	result := parseScripts("")
	if len(result) != 0 {
		t.Errorf("expected no scripts from empty YAML, got %v", result)
	}
}

func TestParseScripts_SingleQuotedCommandExtra2(t *testing.T) {
	content := "scripts:\n  greet: 'echo hello'\n"
	result := parseScripts(content)
	if result["greet"] != "echo hello" {
		t.Errorf("expected 'echo hello', got %q", result["greet"])
	}
}

func TestParseScripts_DoubleQuotedCommand(t *testing.T) {
	content := "scripts:\n  greet: \"echo world\"\n"
	result := parseScripts(content)
	if result["greet"] != "echo world" {
		t.Errorf("expected 'echo world', got %q", result["greet"])
	}
}

func TestParseScripts_ReturnType(t *testing.T) {
	content := "scripts:\n  foo: bar\n"
	result := parseScripts(content)
	if result == nil {
		t.Error("parseScripts should return non-nil map")
	}
}

func TestParseScripts_StopsAtNextTopLevelKey(t *testing.T) {
	content := "scripts:\n  a: cmd-a\ndependencies:\n  - owner/dep\n"
	result := parseScripts(content)
	if _, hasDep := result["dependencies"]; hasDep {
		t.Error("should not include keys from section after scripts")
	}
	if result["a"] != "cmd-a" {
		t.Errorf("expected a=cmd-a, got %q", result["a"])
	}
}

func TestParseScripts_ScriptWithColonInValue(t *testing.T) {
	content := "scripts:\n  url: http://example.com\n"
	result := parseScripts(content)
	if result["url"] == "" {
		t.Error("expected non-empty command for url key")
	}
}

func TestScript_Slice(t *testing.T) {
	scripts := []Script{
		{Name: "a", Command: "cmd-a"},
		{Name: "b", Command: "cmd-b"},
	}
	if len(scripts) != 2 {
		t.Errorf("expected 2 scripts, got %d", len(scripts))
	}
	if scripts[0].Name != "a" {
		t.Errorf("expected name 'a', got %q", scripts[0].Name)
	}
}
