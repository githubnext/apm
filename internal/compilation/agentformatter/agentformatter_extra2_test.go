package agentformatter_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/agentformatter"
)

func TestRenderGeminiStub_WithSlashInPath(t *testing.T) {
	got := agentformatter.RenderGeminiStub("docs/AGENTS.md", "1.2.3")
	if !strings.Contains(got, "docs/AGENTS.md") {
		t.Errorf("stub should contain path, got %q", got)
	}
}

func TestRenderGeminiStub_VersionInComment(t *testing.T) {
	got := agentformatter.RenderGeminiStub("AGENTS.md", "v9.8.7")
	if !strings.Contains(got, "v9.8.7") {
		t.Errorf("stub should contain version v9.8.7, got %q", got)
	}
}

func TestRenderGeminiStub_DefaultPathFallback(t *testing.T) {
	got := agentformatter.RenderGeminiStub("", "1.0")
	if !strings.Contains(got, "@AGENTS.md") {
		t.Errorf("empty path should fallback to AGENTS.md, got %q", got)
	}
}

func TestRenderGeminiStub_HasBuildIDPlaceholder(t *testing.T) {
	got := agentformatter.RenderGeminiStub("AGENTS.md", "1.0")
	if !strings.Contains(got, "__BUILD_ID__") {
		t.Errorf("stub should contain __BUILD_ID__ placeholder, got %q", got)
	}
}

func TestRenderClaudeHeader_IsHTMLComment2(t *testing.T) {
	hdr := agentformatter.RenderClaudeHeader()
	if !strings.HasPrefix(hdr, "<!--") {
		t.Errorf("header should start with <!--, got %q", hdr)
	}
	if !strings.Contains(hdr, "-->") {
		t.Errorf("header should close HTML comment, got %q", hdr)
	}
}

func TestSummarizeClaudeResult_SuccessCountInMessage(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{
		Success:    true,
		Placements: make([]agentformatter.ClaudePlacement, 5),
	}
	msg := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(msg, "5") {
		t.Errorf("success summary should mention count 5, got %q", msg)
	}
}

func TestSummarizeClaudeResult_FailureContainsError(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{
		Success: false,
		Errors:  []string{"file not found"},
	}
	msg := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(msg, "file not found") {
		t.Errorf("failure summary should include error, got %q", msg)
	}
}

func TestClaudePlacement_FieldsZeroValue(t *testing.T) {
	var p agentformatter.ClaudePlacement
	if p.ClaudePath != "" {
		t.Error("ClaudePath should be empty by default")
	}
	if p.InstructionFiles != nil {
		t.Error("InstructionFiles should be nil by default")
	}
}

func TestClaudeCompilationResult_SuccessDefault(t *testing.T) {
	var r agentformatter.ClaudeCompilationResult
	if r.Success {
		t.Error("Success should be false by default")
	}
}

func TestGeminiPlacement_ZeroValue(t *testing.T) {
	var gp agentformatter.GeminiPlacement
	if gp.GeminiPath != "" {
		t.Error("GeminiPath should be empty by default")
	}
}

func TestGeminiCompilationResult_ZeroValue(t *testing.T) {
	var gr agentformatter.GeminiCompilationResult
	if gr.Success {
		t.Error("Success should be false by default")
	}
	if gr.Stats != nil {
		t.Error("Stats should be nil by default")
	}
}

func TestClaudePlacement_DependenciesField(t *testing.T) {
	p := agentformatter.ClaudePlacement{
		Dependencies: []string{"dep1", "dep2"},
	}
	if len(p.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(p.Dependencies))
	}
}

func TestGeminiCompilationResult_WarningsField(t *testing.T) {
	gr := agentformatter.GeminiCompilationResult{
		Success:  true,
		Warnings: []string{"warn1"},
	}
	if len(gr.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(gr.Warnings))
	}
}

func TestSummarizeClaudeResult_ZeroPlacementsSuccess(t *testing.T) {
	r := &agentformatter.ClaudeCompilationResult{Success: true}
	msg := agentformatter.SummarizeClaudeResult(r)
	if !strings.Contains(msg, "[+]") {
		t.Errorf("success summary should have [+] prefix, got %q", msg)
	}
}
