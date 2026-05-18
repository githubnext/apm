package listcmd

import (
	"testing"
)

func TestParseScripts_Empty(t *testing.T) {
	got := parseScripts("")
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestParseScripts_NoScriptsSection(t *testing.T) {
	content := "name: myapp\nversion: 1.0\n"
	got := parseScripts(content)
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestParseScripts_SimpleScript(t *testing.T) {
	content := "name: myapp\nscripts:\n  start: echo hello\n  build: make build\n"
	got := parseScripts(content)
	if got["start"] != "echo hello" {
		t.Errorf("expected 'echo hello', got %q", got["start"])
	}
	if got["build"] != "make build" {
		t.Errorf("expected 'make build', got %q", got["build"])
	}
}

func TestParseScripts_QuotedCommand(t *testing.T) {
	content := "scripts:\n  start: \"codex run main.prompt.md\"\n"
	got := parseScripts(content)
	if got["start"] != "codex run main.prompt.md" {
		t.Errorf("expected 'codex run main.prompt.md', got %q", got["start"])
	}
}

func TestParseScripts_SingleQuotedCommand(t *testing.T) {
	content := "scripts:\n  lint: 'ruff check src/'\n"
	got := parseScripts(content)
	if got["lint"] != "ruff check src/" {
		t.Errorf("expected 'ruff check src/', got %q", got["lint"])
	}
}

func TestParseScripts_CommentLines(t *testing.T) {
	content := "# top comment\nscripts:\n  # skip this\n  test: pytest\n"
	got := parseScripts(content)
	if got["test"] != "pytest" {
		t.Errorf("expected 'pytest', got %q", got["test"])
	}
	if _, ok := got["# skip this"]; ok {
		t.Error("should not have parsed comment as a key")
	}
}

func TestParseScripts_BlockEndsOnNewTopLevel(t *testing.T) {
	content := "scripts:\n  run: go run .\nother:\n  key: val\n"
	got := parseScripts(content)
	if len(got) != 1 {
		t.Errorf("expected 1 script, got %d: %v", len(got), got)
	}
	if got["run"] != "go run ." {
		t.Errorf("expected 'go run .', got %q", got["run"])
	}
}

func TestParseScripts_MultipleScripts(t *testing.T) {
	content := "scripts:\n  build: make\n  test: go test ./...\n  lint: golangci-lint run\n"
	got := parseScripts(content)
	if len(got) != 3 {
		t.Errorf("expected 3 scripts, got %d", len(got))
	}
}

func TestParseScripts_ColonInCommand(t *testing.T) {
	content := "scripts:\n  serve: http://localhost:8080\n"
	got := parseScripts(content)
	// "serve" is the key; value should be "http://localhost:8080"
	if got["serve"] != "http://localhost:8080" {
		t.Errorf("unexpected value: %q", got["serve"])
	}
}
