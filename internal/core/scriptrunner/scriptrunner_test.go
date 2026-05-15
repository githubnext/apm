package scriptrunner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// substituteParameters
// ---------------------------------------------------------------------------

func TestSubstituteParametersBasic(t *testing.T) {
	result := substituteParameters("Hello ${input:name}!", map[string]string{"name": "world"})
	if result != "Hello world!" {
		t.Errorf("got %q", result)
	}
}

func TestSubstituteParametersMissing(t *testing.T) {
	// Missing key should leave placeholder intact.
	result := substituteParameters("Hello ${input:name}!", map[string]string{})
	if result != "Hello ${input:name}!" {
		t.Errorf("got %q", result)
	}
}

func TestSubstituteParametersMultiple(t *testing.T) {
	result := substituteParameters("${input:a} + ${input:b} = ${input:c}", map[string]string{"a": "1", "b": "2", "c": "3"})
	if result != "1 + 2 = 3" {
		t.Errorf("got %q", result)
	}
}

func TestSubstituteParametersNoPlaceholders(t *testing.T) {
	result := substituteParameters("no placeholders here", map[string]string{"x": "y"})
	if result != "no placeholders here" {
		t.Errorf("got %q", result)
	}
}

// ---------------------------------------------------------------------------
// detectRuntime
// ---------------------------------------------------------------------------

func TestDetectRuntimeCopilot(t *testing.T) {
	cases := []string{
		"gh copilot suggest something",
		"  gh copilot explain code",
	}
	for _, c := range cases {
		if detectRuntime(c) != RuntimeCopilot {
			t.Errorf("expected copilot for %q", c)
		}
	}
}

func TestDetectRuntimeCodex(t *testing.T) {
	if detectRuntime("codex run something") != RuntimeCodex {
		t.Error("expected codex")
	}
}

func TestDetectRuntimeGemini(t *testing.T) {
	if detectRuntime("gemini -p something") != RuntimeGemini {
		t.Error("expected gemini")
	}
}

func TestDetectRuntimeUnknown(t *testing.T) {
	if detectRuntime("echo hello") != RuntimeUnknown {
		t.Error("expected unknown")
	}
}

func TestDetectRuntimeLLM(t *testing.T) {
	if detectRuntime("llm prompt 'something'") != RuntimeLLM {
		t.Error("expected llm")
	}
}

// ---------------------------------------------------------------------------
// splitArgs
// ---------------------------------------------------------------------------

func TestSplitArgsSimple(t *testing.T) {
	args := splitArgs("git commit -m hello")
	if len(args) != 4 {
		t.Fatalf("expected 4, got %d: %v", len(args), args)
	}
	if args[0] != "git" || args[3] != "hello" {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestSplitArgsQuoted(t *testing.T) {
	args := splitArgs(`echo "hello world"`)
	if len(args) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(args), args)
	}
	if args[1] != "hello world" {
		t.Errorf("expected 'hello world', got %q", args[1])
	}
}

func TestSplitArgsSingleQuoted(t *testing.T) {
	args := splitArgs("echo 'hello world'")
	if len(args) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(args), args)
	}
	if args[1] != "hello world" {
		t.Errorf("expected 'hello world', got %q", args[1])
	}
}

func TestSplitArgsEmpty(t *testing.T) {
	args := splitArgs("")
	if len(args) != 0 {
		t.Errorf("expected empty, got %v", args)
	}
}

// ---------------------------------------------------------------------------
// isVirtualPackageReference
// ---------------------------------------------------------------------------

func TestIsVirtualPackageReference(t *testing.T) {
	cases := []struct {
		name   string
		result bool
	}{
		{"owner/repo/path", true},
		{"owner/repo", false},
		{"localscript", false},
		{"build", false},
	}
	for _, c := range cases {
		got := isVirtualPackageReference(c.name)
		if got != c.result {
			t.Errorf("isVirtualPackageReference(%q) = %v, want %v", c.name, got, c.result)
		}
	}
}

// ---------------------------------------------------------------------------
// isValidEnvVarName
// ---------------------------------------------------------------------------

func TestIsValidEnvVarName(t *testing.T) {
	valid := []string{"FOO", "BAR_BAZ", "_PRIV", "X1"}
	for _, v := range valid {
		if !isValidEnvVarName(v) {
			t.Errorf("expected %q to be valid", v)
		}
	}
	invalid := []string{"1INVALID", "foo-bar", "foo bar", ""}
	for _, v := range invalid {
		if isValidEnvVarName(v) {
			t.Errorf("expected %q to be invalid", v)
		}
	}
}

// ---------------------------------------------------------------------------
// parseSimpleYAML
// ---------------------------------------------------------------------------

func TestParseSimpleYAMLBasic(t *testing.T) {
	yml := "name: myapp\nversion: 1.0\n"
	result := parseSimpleYAML(yml)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["name"] != "myapp" {
		t.Errorf("got name=%q", result["name"])
	}
}

func TestParseSimpleYAMLScriptsBlock(t *testing.T) {
	yml := "scripts:\n  build: go build ./...\n  test: go test ./...\n"
	result := parseSimpleYAML(yml)
	scripts, ok := result["scripts"].(map[string]any)
	if !ok {
		t.Fatalf("expected scripts map, got %T", result["scripts"])
	}
	if scripts["build"] != "go build ./..." {
		t.Errorf("unexpected build script: %q", scripts["build"])
	}
}

func TestParseSimpleYAMLEmpty(t *testing.T) {
	result := parseSimpleYAML("")
	if result == nil {
		t.Fatal("expected non-nil")
	}
}

// ---------------------------------------------------------------------------
// unquoteYAML
// ---------------------------------------------------------------------------

func TestUnquoteYAML(t *testing.T) {
	cases := []struct{ in, out string }{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{"plain", "plain"},
		{"", ""},
	}
	for _, c := range cases {
		got := unquoteYAML(c.in)
		if got != c.out {
			t.Errorf("unquoteYAML(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

// ---------------------------------------------------------------------------
// formatScriptHeader
// ---------------------------------------------------------------------------

func TestFormatScriptHeader(t *testing.T) {
	lines := formatScriptHeader("build", map[string]string{"env": "prod"})
	if len(lines) == 0 {
		t.Error("expected at least one line")
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "build") {
		t.Errorf("expected script name in header: %q", joined)
	}
}

// ---------------------------------------------------------------------------
// copyEnv / envMapToSlice
// ---------------------------------------------------------------------------

func TestCopyEnv(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	cp := copyEnv(orig)
	cp["A"] = "99"
	if orig["A"] != "1" {
		t.Error("original should not be modified")
	}
}

func TestEnvMapToSlice(t *testing.T) {
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	slice := envMapToSlice(m)
	if len(slice) != 2 {
		t.Fatalf("expected 2, got %d", len(slice))
	}
	found := 0
	for _, s := range slice {
		if s == "FOO=bar" || s == "BAZ=qux" {
			found++
		}
	}
	if found != 2 {
		t.Errorf("unexpected slice: %v", slice)
	}
}

// ---------------------------------------------------------------------------
// PromptCompiler
// ---------------------------------------------------------------------------

func TestPromptCompilerCompile(t *testing.T) {
	dir := t.TempDir()
	promptFile := filepath.Join(dir, "hello.prompt.md")
	if err := os.WriteFile(promptFile, []byte("# Hello\n\nHello ${input:name}!"), 0o644); err != nil {
		t.Fatal(err)
	}

	compiler := &PromptCompiler{CompiledDir: filepath.Join(dir, "compiled")}
	out, err := compiler.Compile(promptFile, map[string]string{"name": "tester"})
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), "Hello tester!") {
		t.Errorf("expected substituted content, got: %s", data)
	}
}

func TestPromptCompilerCompileFrontmatterStripped(t *testing.T) {
	dir := t.TempDir()
	promptFile := filepath.Join(dir, "test.prompt.md")
	content := "---\ntitle: test\n---\nHello ${input:name}!"
	if err := os.WriteFile(promptFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	compiler := &PromptCompiler{CompiledDir: filepath.Join(dir, "compiled")}
	out, err := compiler.Compile(promptFile, map[string]string{"name": "world"})
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if strings.Contains(string(data), "title: test") {
		t.Errorf("frontmatter should be stripped: %s", data)
	}
	if !strings.Contains(string(data), "Hello world!") {
		t.Errorf("expected substituted content: %s", data)
	}
}

// ---------------------------------------------------------------------------
// New and ListScripts
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	sr := New(true)
	if sr == nil {
		t.Fatal("expected non-nil ScriptRunner")
	}
	if sr.Compiler == nil {
		t.Error("expected non-nil Compiler")
	}
}

func TestListScripts(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	_ = os.Chdir(dir)

	yml := "name: myapp\nscripts:\n  build: go build ./...\n  test: go test ./...\n"
	_ = os.WriteFile(filepath.Join(dir, "apm.yml"), []byte(yml), 0o644)

	sr := New(false)
	scripts := sr.ListScripts()
	if scripts["build"] != "go build ./..." {
		t.Errorf("unexpected scripts map: %v", scripts)
	}
}

func TestListScriptsNoFile(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	_ = os.Chdir(dir)

	sr := New(false)
	scripts := sr.ListScripts()
	if len(scripts) != 0 {
		t.Errorf("expected empty map, got %v", scripts)
	}
}

// ---------------------------------------------------------------------------
// generateRuntimeCommand
// ---------------------------------------------------------------------------

func TestGenerateRuntimeCommandCopilot(t *testing.T) {
	cmd := generateRuntimeCommand(RuntimeCopilot, "path/to/prompt.txt")
	if !strings.Contains(cmd, "copilot") && !strings.Contains(cmd, "prompt.txt") {
		t.Errorf("unexpected command: %q", cmd)
	}
}

func TestGenerateRuntimeCommandUnknown(t *testing.T) {
	cmd := generateRuntimeCommand(RuntimeUnknown, "path/to/prompt.txt")
	if cmd == "" {
		t.Error("expected non-empty fallback command")
	}
}
