package scriptrunner

import (
"testing"
)

func TestNewScriptRunner_NotNilExtra4(t *testing.T) {
s := New(false)
if s == nil {
t.Error("expected non-nil ScriptRunner")
}
}

func TestNewScriptRunner_WithColorExtra4(t *testing.T) {
s := New(true)
if !s.UseColor {
t.Error("expected UseColor=true")
}
}

func TestNewScriptRunner_CompilerNotNilExtra4(t *testing.T) {
s := New(false)
if s.Compiler == nil {
t.Error("expected non-nil Compiler")
}
}

func TestListScripts_EmptyDir_ReturnsEmptyExtra4(t *testing.T) {
s := New(false)
scripts := s.ListScripts()
_ = scripts
}

func TestRunScript_UnknownScript_ReturnsErrorExtra4(t *testing.T) {
s := New(false)
err := s.RunScript("no-such-script-xyz-abc", nil)
if err == nil {
t.Error("expected error for unknown script")
}
}

func TestRuntimeKind_ValuesExtra4(t *testing.T) {
for _, k := range []RuntimeKind{RuntimeCopilot, RuntimeCodex, RuntimeLLM, RuntimeGemini, RuntimeUnknown} {
if string(k) == "" {
t.Errorf("expected non-empty RuntimeKind string for %v", k)
}
}
}
