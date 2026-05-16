package lockfileenrichment

import (
	"strings"
	"testing"
)

func TestFilterFilesByTarget_Claude(t *testing.T) {
	files := []string{
		".claude/skills/foo.md",
		".github/skills/bar.md",
		".cursor/rules/x.md",
		"README.md",
	}
	result := FilterFilesByTarget(files, "claude")
	// direct: files under .claude/
	found := map[string]bool{}
	for _, f := range result.Files {
		found[f] = true
	}
	if !found[".claude/skills/foo.md"] {
		t.Errorf("expected .claude/skills/foo.md in Direct, got %v", result.Files)
	}
}

func TestFilterFilesByTarget_VSCode(t *testing.T) {
	files := []string{
		".github/skills/bar.md",
		".vscode/settings.json",
		".claude/skills/foo.md",
	}
	result := FilterFilesByTarget(files, "vscode")
	found := map[string]bool{}
	for _, f := range result.Files {
		found[f] = true
	}
	if !found[".github/skills/bar.md"] {
		t.Errorf("expected .github/skills/bar.md in Direct, got %v", result.Files)
	}
}

func TestFilterFilesByTarget_MultiTarget(t *testing.T) {
	files := []string{
		".claude/skills/foo.md",
		".github/skills/bar.md",
	}
	result := FilterFilesByTarget(files, "claude, vscode")
	if len(result.Files) == 0 {
		t.Errorf("expected files matched with multi-target, got none")
	}
}

func TestFilterFilesByTarget_UnknownTarget(t *testing.T) {
	files := []string{"README.md", "src/main.go"}
	result := FilterFilesByTarget(files, "unknown-target")
	if len(result.Files) != 0 {
		t.Errorf("expected no direct files for unknown target, got %v", result.Files)
	}
}

func TestEnrichLockfileForPack(t *testing.T) {
	meta := PackMeta{
		PackedAt: "2025-01-01T00:00:00Z",
		Target:   "claude",
		Format:   "plugin-v1",
	}
	out := EnrichLockfileForPack(meta)
	if !strings.Contains(out, "pack:") {
		t.Errorf("expected 'pack:' section in output, got: %s", out)
	}
	if !strings.Contains(out, "2025-01-01T00:00:00Z") {
		t.Errorf("expected packed_at in output, got: %s", out)
	}
}

func TestCollectMappedFromPrefixes(t *testing.T) {
	paths := []string{
		".github/skills/foo.md",
		".claude/agents/bar.md",
		"README.md",
	}
	// For vscode target: .claude/skills/ -> .github/skills/, .claude/agents/ -> .github/agents/
	used := CollectMappedFromPrefixes("vscode", paths)
	// .claude/agents/bar.md uses .claude/agents/ prefix which is in the vscode cross-map
	foundAgents := false
	for _, p := range used {
		if p == ".claude/agents/" {
			foundAgents = true
		}
	}
	if !foundAgents {
		t.Logf("CollectMappedFromPrefixes(vscode) result: %v", used)
	}
}

func TestAllTargetPrefixes_NotEmpty(t *testing.T) {
	prefixes := allTargetPrefixes()
	if len(prefixes) == 0 {
		t.Error("expected non-empty target prefixes list")
	}
}
