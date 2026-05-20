package scriptrunner

import (
	"testing"
)

func TestRuntimeKind_Copilot(t *testing.T) {
	if RuntimeCopilot != "copilot" {
		t.Errorf("unexpected RuntimeCopilot value: %q", RuntimeCopilot)
	}
}

func TestRuntimeKind_Codex(t *testing.T) {
	if RuntimeCodex != "codex" {
		t.Errorf("unexpected RuntimeCodex value: %q", RuntimeCodex)
	}
}

func TestRuntimeKind_LLM(t *testing.T) {
	if RuntimeLLM != "llm" {
		t.Errorf("unexpected RuntimeLLM value: %q", RuntimeLLM)
	}
}

func TestRuntimeKind_Unknown(t *testing.T) {
	if RuntimeUnknown != "unknown" {
		t.Errorf("unexpected RuntimeUnknown value: %q", RuntimeUnknown)
	}
}

func TestNew_NotNil(t *testing.T) {
	s := New(false)
	if s == nil {
		t.Error("expected non-nil ScriptRunner")
	}
}

func TestNew_UseColorField(t *testing.T) {
	s := New(true)
	if !s.UseColor {
		t.Error("expected UseColor=true")
	}
}

func TestNew_CompilerNotNil(t *testing.T) {
	s := New(false)
	if s.Compiler == nil {
		t.Error("expected non-nil Compiler")
	}
}

func TestNewPromptCompiler_NotNil(t *testing.T) {
	c := NewPromptCompiler()
	if c == nil {
		t.Error("expected non-nil PromptCompiler")
	}
}

func TestNewPromptCompiler_CompiledDir(t *testing.T) {
	c := NewPromptCompiler()
	if c.CompiledDir == "" {
		t.Error("expected non-empty CompiledDir")
	}
}
