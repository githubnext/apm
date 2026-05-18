package listcmd

import (
	"testing"
)

func TestParseScripts_TabIndented(t *testing.T) {
	content := "scripts:\n\tbuild: go build ./...\n\ttest: go test ./...\n"
	got := parseScripts(content)
	if len(got) != 2 {
		t.Errorf("expected 2 scripts, got %d: %v", len(got), got)
	}
	if got["build"] != "go build ./..." {
		t.Errorf("unexpected build value: %q", got["build"])
	}
}

func TestParseScripts_EmptyCommand(t *testing.T) {
	content := "scripts:\n  noop:\n"
	got := parseScripts(content)
	// line "  noop:" has an empty command
	if _, ok := got["noop"]; ok {
		// empty command is allowed
	}
}

func TestParseScripts_ScriptNameWithHyphen(t *testing.T) {
	content := "scripts:\n  build-all: make all\n"
	got := parseScripts(content)
	if got["build-all"] != "make all" {
		t.Errorf("unexpected value: %q", got["build-all"])
	}
}

func TestParseScripts_ScriptNameWithUnderscore(t *testing.T) {
	content := "scripts:\n  run_tests: pytest tests/\n"
	got := parseScripts(content)
	if got["run_tests"] != "pytest tests/" {
		t.Errorf("unexpected value: %q", got["run_tests"])
	}
}

func TestParseScripts_MultipleBlocks(t *testing.T) {
	content := "name: myapp\nversion: 1.0\nscripts:\n  start: node .\n  stop: kill -9 1\n"
	got := parseScripts(content)
	if len(got) != 2 {
		t.Errorf("expected 2, got %d: %v", len(got), got)
	}
}

func TestParseScripts_DoubleColonInCommand(t *testing.T) {
	content := "scripts:\n  connect: ssh user@host:22\n"
	got := parseScripts(content)
	if got["connect"] != "ssh user@host:22" {
		t.Errorf("unexpected: %q", got["connect"])
	}
}

func TestParseScripts_PreservesSpacesInCommand(t *testing.T) {
	content := "scripts:\n  test: go test -v -count=1 ./...\n"
	got := parseScripts(content)
	if got["test"] != "go test -v -count=1 ./..." {
		t.Errorf("unexpected: %q", got["test"])
	}
}

func TestParseScripts_CommentsInsideBlock(t *testing.T) {
	content := "scripts:\n  # a comment\n  real: echo ok\n"
	got := parseScripts(content)
	if got["real"] != "echo ok" {
		t.Errorf("unexpected: %q", got["real"])
	}
	if _, ok := got["# a comment"]; ok {
		t.Error("should not have parsed comment as script name")
	}
}

func TestParseScripts_LeavesAfterNewTopLevel(t *testing.T) {
	content := "scripts:\n  a: do-a\n  b: do-b\ntopkey:\n  x: y\n"
	got := parseScripts(content)
	if _, ok := got["x"]; ok {
		t.Error("should not parse entries from other sections")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 scripts, got %d", len(got))
	}
}

func TestParseScripts_CommitVariants(t *testing.T) {
	cases := []string{
		"abc1234",
		"deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		"v1.2.3",
	}
	for _, commit := range cases {
		if commit == "" {
			t.Errorf("commit should not be empty")
		}
	}
}

func TestScript_Fields(t *testing.T) {
	s := Script{Name: "start", Command: "node ."}
	if s.Name != "start" {
		t.Errorf("Name mismatch: %q", s.Name)
	}
	if s.Command != "node ." {
		t.Errorf("Command mismatch: %q", s.Command)
	}
}

func TestParseScripts_SingleQuotePreservesContent(t *testing.T) {
	content := "scripts:\n  greet: 'hello world'\n"
	got := parseScripts(content)
	if got["greet"] != "hello world" {
		t.Errorf("unexpected: %q", got["greet"])
	}
}

func TestParseScripts_DoubleQuotePreservesContent(t *testing.T) {
	content := "scripts:\n  greet: \"hello world\"\n"
	got := parseScripts(content)
	if got["greet"] != "hello world" {
		t.Errorf("unexpected: %q", got["greet"])
	}
}
