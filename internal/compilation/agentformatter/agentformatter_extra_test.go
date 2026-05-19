package agentformatter_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/agentformatter"
)

func TestRenderGeminiStub_VersionVariants(t *testing.T) {
	versions := []string{"0.0.1", "1.0.0", "2.3.4-beta", "10.20.30"}
	for _, v := range versions {
		stub := agentformatter.RenderGeminiStub("AGENTS.md", v)
		if !strings.Contains(stub, v) {
			t.Errorf("stub for version %q should contain version string", v)
		}
	}
}

func TestRenderGeminiStub_PathNormalization(t *testing.T) {
	// paths with subdirectories
	stub := agentformatter.RenderGeminiStub("docs/sub/AGENTS.md", "1.0.0")
	if !strings.Contains(stub, "docs/sub/AGENTS.md") {
		t.Errorf("stub should preserve path, got: %s", stub)
	}
}

func TestRenderGeminiStub_AlwaysEndsWithNewline(t *testing.T) {
	stub := agentformatter.RenderGeminiStub("AGENTS.md", "1.0.0")
	if !strings.HasSuffix(stub, "\n") {
		t.Error("stub should end with newline")
	}
}

func TestRenderGeminiStub_StartsWithHTMLComment(t *testing.T) {
	stub := agentformatter.RenderGeminiStub("AGENTS.md", "1.0.0")
	if !strings.HasPrefix(stub, "<!--") {
		t.Errorf("stub should start with <!-- comment, got: %q", stub[:min(20, len(stub))])
	}
}

func TestRenderGeminiStub_ContainsAtSymbol(t *testing.T) {
	stub := agentformatter.RenderGeminiStub("AGENTS.md", "1.0.0")
	if !strings.Contains(stub, "@") {
		t.Error("stub should contain @ path reference")
	}
}

func TestRenderClaudeHeader_EndsWithNewline(t *testing.T) {
	header := agentformatter.RenderClaudeHeader()
	if !strings.HasSuffix(header, "\n") {
		t.Error("claude header should end with newline")
	}
}

func TestRenderClaudeHeader_HasContent(t *testing.T) {
	header := agentformatter.RenderClaudeHeader()
	if header == "" {
		t.Error("claude header should not be empty")
	}
}

func TestRenderClaudeHeader_IsHTMLComment(t *testing.T) {
	header := agentformatter.RenderClaudeHeader()
	if !strings.HasPrefix(header, "<!--") {
		t.Errorf("expected HTML comment, got %q", header)
	}
}

func TestSummarizeClaudeResult_ZeroPlacements(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{
		Success:    true,
		Placements: []agentformatter.ClaudePlacement{},
	}
	summary := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(summary, "[+]") {
		t.Error("zero-placement success should still show [+]")
	}
	if !strings.Contains(summary, "0") {
		t.Error("summary should mention 0 placements")
	}
}

func TestSummarizeClaudeResult_ThreeErrors(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{
		Success: false,
		Errors:  []string{"err1", "err2", "err3"},
	}
	summary := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(summary, "[x]") {
		t.Error("failure should show [x]")
	}
	// At minimum should contain first error
	if !strings.Contains(summary, "err1") {
		t.Error("summary should contain error messages")
	}
}

func TestSummarizeClaudeResult_SuccessWithManyPlacements(t *testing.T) {
	placements := make([]agentformatter.ClaudePlacement, 10)
	r := &agentformatter.ClaudeCompilationResult{
		Success:    true,
		Placements: placements,
	}
	summary := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(summary, "10") {
		t.Errorf("summary should mention 10 placements, got: %s", summary)
	}
}

func TestClaudeCompilationResult_SuccessField(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{Success: true}
	if !r.Success {
		t.Error("expected Success=true")
	}
	r2 := &agentformatter.ClaudeCompilationResult{Success: false}
	if r2.Success {
		t.Error("expected Success=false")
	}
}

func TestClaudeCompilationResult_WarningsField(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{
		Success:  true,
		Warnings: []string{"warn1", "warn2"},
	}
	if len(r.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(r.Warnings))
	}
}

func TestGeminiCompilationResult_AllFields(t *testing.T) {
	r := &agentformatter.GeminiCompilationResult{
		Success:  true,
		Warnings: []string{"w1"},
		Errors:   []string{},
		Stats:    map[string]float64{"count": 3.0},
	}
	if !r.Success {
		t.Error("expected Success=true")
	}
	if r.Stats["count"] != 3.0 {
		t.Error("expected Stats count=3.0")
	}
}

func TestGeminiPlacement_Fields(t *testing.T) {
	gp := agentformatter.GeminiPlacement{
		GeminiPath:       "/path/to/GEMINI.md",
		InstructionFiles: []string{"a.md", "b.md"},
	}
	if gp.GeminiPath != "/path/to/GEMINI.md" {
		t.Error("GeminiPath not set")
	}
	if len(gp.InstructionFiles) != 2 {
		t.Error("expected 2 InstructionFiles")
	}
}

func TestClaudePlacement_AllFields(t *testing.T) {
	cp := agentformatter.ClaudePlacement{
		ClaudePath:       ".claude/CLAUDE.md",
		InstructionFiles: []string{"a.md"},
		AgentFiles:       []string{"agent.md"},
		Dependencies:     []string{"dep1"},
		CoveragePatterns: []string{"src/**"},
		SourceAttribution: map[string]string{"a.md": "source"},
	}
	if cp.ClaudePath != ".claude/CLAUDE.md" {
		t.Error("ClaudePath not set")
	}
	if len(cp.AgentFiles) != 1 {
		t.Error("AgentFiles not set")
	}
	if cp.SourceAttribution["a.md"] != "source" {
		t.Error("SourceAttribution not set")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
