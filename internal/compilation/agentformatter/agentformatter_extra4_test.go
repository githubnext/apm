package agentformatter_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/agentformatter"
)

func TestClaudePlacement_ZeroValue_Extra4(t *testing.T) {
	var p agentformatter.ClaudePlacement
	if p.ClaudePath != "" {
		t.Fatal("expected empty ClaudePath")
	}
}

func TestClaudeCompilationResult_ZeroValue_Extra4(t *testing.T) {
	var r agentformatter.ClaudeCompilationResult
	if r.Success {
		t.Fatal("expected Success=false")
	}
}

func TestGeminiPlacement_ZeroValue_Extra4(t *testing.T) {
	var p agentformatter.GeminiPlacement
	if p.GeminiPath != "" {
		t.Fatal("expected empty GeminiPath")
	}
}

func TestGeminiCompilationResult_ZeroValue_Extra4(t *testing.T) {
	var r agentformatter.GeminiCompilationResult
	if r.Success {
		t.Fatal("expected Success=false")
	}
}

func TestRenderGeminiStub_ContainsAt_Extra4(t *testing.T) {
	out := agentformatter.RenderGeminiStub("AGENTS.md", "1.0.0")
	if !strings.Contains(out, "@") {
		t.Fatal("expected @ in output")
	}
}

func TestRenderGeminiStub_ContainsPath_Extra4(t *testing.T) {
	out := agentformatter.RenderGeminiStub("some/path.md", "2.0.0")
	if !strings.Contains(out, "some/path.md") {
		t.Fatalf("expected path in output, got: %s", out)
	}
}

func TestSummarizeClaudeResult_Failure_Extra4(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{
		Success: false,
		Errors:  []string{"missing input"},
	}
	s := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(s, "failed") {
		t.Fatalf("expected 'failed' in summary, got: %s", s)
	}
}

func TestSummarizeClaudeResult_Success_Extra4(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{
		Success: true,
		Placements: []agentformatter.ClaudePlacement{
			{ClaudePath: "CLAUDE.md"},
		},
	}
	s := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(s, "compiled") {
		t.Fatalf("expected 'compiled' in summary, got: %s", s)
	}
}

func TestRenderClaudeHeader_NotEmpty_Extra4(t *testing.T) {
	h := agentformatter.RenderClaudeHeader()
	if h == "" {
		t.Fatal("expected non-empty claude header")
	}
}

func TestRenderGeminiStub_ForwardSlash_Extra4(t *testing.T) {
	out := agentformatter.RenderGeminiStub("dir/sub/AGENTS.md", "0.1.0")
	if !strings.Contains(out, "dir/sub/AGENTS.md") {
		t.Fatalf("expected forward-slash path in output, got: %s", out)
	}
}
