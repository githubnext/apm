package agentscompiler

import (
	"errors"
	"testing"
)

var errSentinel = errors.New("test error")

func TestDefaultConfig_OutputPath(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.OutputPath == "" {
		t.Fatal("expected non-empty OutputPath in DefaultConfig")
	}
}

func TestDefaultConfig_TargetIsAll(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Target != TargetAll {
		t.Fatalf("expected target 'all', got %q", cfg.Target)
	}
}

func TestDefaultConfig_StrategyDistributed(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Strategy != StrategyDistributed {
		t.Fatalf("expected strategy 'distributed', got %q", cfg.Strategy)
	}
}

func TestDefaultConfig_ResolveLinksTrue(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.ResolveLinks {
		t.Fatal("expected ResolveLinks to be true in DefaultConfig")
	}
}

func TestCompilationResult_OKFalseWhenError(t *testing.T) {
	r := CompilationResult{Error: errSentinel}
	if r.OK() {
		t.Fatal("OK should return false when Error is set")
	}
}

func TestCompilationResult_OKTrueWhenNoError(t *testing.T) {
	r := CompilationResult{}
	if !r.OK() {
		t.Fatal("OK should return true when Error is nil")
	}
}

func TestMergedResult_OKNoResults2(t *testing.T) {
	m := MergedResult{}
	if !m.OK() {
		t.Fatal("empty MergedResult should be OK")
	}
}

func TestTargetConstants_NonEmpty(t *testing.T) {
	for _, t2 := range []CompileTargetType{TargetAll, TargetAgents, TargetCopilot, TargetClaude, TargetGemini} {
		if t2 == "" {
			t.Fatal("target constant should not be empty")
		}
	}
}

func TestStrategyConstants_NonEmpty(t *testing.T) {
	for _, s := range []CompilationStrategy{StrategyDistributed, StrategySingleFile} {
		if s == "" {
			t.Fatal("strategy constant should not be empty")
		}
	}
}

func TestNew_ReturnNotNil(t *testing.T) {
	a := New(".")
	if a == nil {
		t.Fatal("New should return non-nil AgentsCompiler")
	}
}

func TestBuildIDPlaceholder_NonEmpty(t *testing.T) {
	if BuildIDPlaceholder == "" {
		t.Fatal("BuildIDPlaceholder should not be empty")
	}
}

func TestCopilotRootGeneratedMarker_NonEmpty(t *testing.T) {
	if CopilotRootGeneratedMarker == "" {
		t.Fatal("CopilotRootGeneratedMarker should not be empty")
	}
}

func TestCompilationConfig_ZeroValue(t *testing.T) {
	var cfg CompilationConfig
	if cfg.OutputPath != "" || cfg.Target != "" {
		t.Fatal("zero-value CompilationConfig should have empty fields")
	}
}
