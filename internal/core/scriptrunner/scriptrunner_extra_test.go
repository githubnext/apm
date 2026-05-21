package scriptrunner

import (
	"testing"
)

func TestRuntimeKindConstants(t *testing.T) {
	kinds := []RuntimeKind{RuntimeCopilot, RuntimeCodex, RuntimeLLM, RuntimeGemini, RuntimeUnknown}
	for _, k := range kinds {
		if k == "" {
			t.Error("RuntimeKind constant should not be empty")
		}
	}
	// all must be distinct
	seen := map[RuntimeKind]bool{}
	for _, k := range kinds {
		if seen[k] {
			t.Errorf("duplicate RuntimeKind %q", k)
		}
		seen[k] = true
	}
}

func TestScriptRunnerNewFalse(t *testing.T) {
	s := New(false)
	if s == nil {
		t.Fatal("New returned nil")
	}
	if s.Compiler == nil {
		t.Error("Compiler should not be nil")
	}
	if s.UseColor {
		t.Error("UseColor should be false")
	}
}

func TestScriptRunnerNewWithColor(t *testing.T) {
	s := New(true)
	if !s.UseColor {
		t.Error("UseColor should be true")
	}
}

func TestDetectRuntimeLLMVariants(t *testing.T) {
	cases := []string{"llm run something", "llm prompt exec"}
	for _, c := range cases {
		if detectRuntime(c) != RuntimeLLM {
			t.Errorf("expected LLM for %q, got %q", c, detectRuntime(c))
		}
	}
}

func TestSubstituteParametersEmptyParams(t *testing.T) {
	result := substituteParameters("no vars here", nil)
	if result != "no vars here" {
		t.Errorf("got %q", result)
	}
}

func TestSubstituteParametersRepeatedVar(t *testing.T) {
	result := substituteParameters("${input:x} and ${input:x} again", map[string]string{"x": "hello"})
	if result != "hello and hello again" {
		t.Errorf("got %q", result)
	}
}

func TestCopyEnvNil(t *testing.T) {
	cp := copyEnv(nil)
	if cp == nil {
		t.Error("copyEnv(nil) should return non-nil map")
	}
}

func TestSplitArgsBasicExtra(t *testing.T) {
	args := splitArgs("cmd arg1 arg2")
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d: %v", len(args), args)
	}
}

func TestBuildCollisionError(t *testing.T) {
	msg := buildCollisionError("my-script", []string{"a/my-script.prompt.md", "b/my-script.prompt.md"})
	if msg == "" {
		t.Error("expected non-empty collision error")
	}
}

func TestGenerateRuntimeCommandCodex(t *testing.T) {
	cmd := generateRuntimeCommand(RuntimeCodex, "my.prompt.md")
	if cmd == "" {
		t.Error("expected non-empty runtime command for codex")
	}
}

func TestGenerateRuntimeCommandLLM(t *testing.T) {
	cmd := generateRuntimeCommand(RuntimeLLM, "my.prompt.md")
	if cmd == "" {
		t.Error("expected non-empty runtime command for llm")
	}
}

func TestGenerateRuntimeCommandGemini(t *testing.T) {
	cmd := generateRuntimeCommand(RuntimeGemini, "my.prompt.md")
	if cmd == "" {
		t.Error("expected non-empty runtime command for gemini")
	}
}
