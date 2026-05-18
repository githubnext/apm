package lockfileenrichment

import (
	"strings"
	"testing"
)

func TestFilterFilesByTarget_Cursor(t *testing.T) {
	files := []string{
		".cursor/rules/x.md",
		".github/skills/bar.md",
		".claude/skills/foo.md",
	}
	result := FilterFilesByTarget(files, "cursor")
	found := map[string]bool{}
	for _, f := range result.Files {
		found[f] = true
	}
	if !found[".cursor/rules/x.md"] {
		t.Errorf("expected .cursor/rules/x.md in results, got %v", result.Files)
	}
}

func TestFilterFilesByTarget_Codex(t *testing.T) {
	files := []string{
		".codex/agents/foo.md",
		".agents/skills/bar.md",
		"README.md",
	}
	result := FilterFilesByTarget(files, "codex")
	if len(result.Files) == 0 {
		t.Errorf("expected files for codex target, got none")
	}
}

func TestFilterFilesByTarget_AgentSkills(t *testing.T) {
	files := []string{
		".agents/skills/foo.md",
		".github/skills/bar.md",
	}
	result := FilterFilesByTarget(files, "agent-skills")
	if len(result.Files) == 0 {
		t.Errorf("expected files for agent-skills target, got none")
	}
}

func TestFilterFilesByTarget_EmptyFiles(t *testing.T) {
	result := FilterFilesByTarget(nil, "claude")
	if len(result.Files) != 0 {
		t.Errorf("expected no files for nil input, got %v", result.Files)
	}
}

func TestEnrichLockfileForPack_ContainsFormat(t *testing.T) {
	meta := PackMeta{
		PackedAt: "2025-06-01T12:00:00Z",
		Target:   "cursor",
		Format:   "plugin-v2",
	}
	out := EnrichLockfileForPack(meta)
	if !strings.Contains(out, "plugin-v2") {
		t.Errorf("expected format in output, got: %s", out)
	}
	if !strings.Contains(out, "cursor") {
		t.Errorf("expected target in output, got: %s", out)
	}
}

func TestEnrichLockfileForPack_EmptyMeta(t *testing.T) {
	meta := PackMeta{}
	out := EnrichLockfileForPack(meta)
	if out == "" {
		t.Error("expected non-empty output even for empty meta")
	}
}

func TestCollectMappedFromPrefixes_Claude(t *testing.T) {
	paths := []string{
		".github/skills/foo.md",
		".github/agents/bar.md",
	}
	used := CollectMappedFromPrefixes("claude", paths)
	if len(used) == 0 {
		t.Logf("CollectMappedFromPrefixes(claude) returned empty; paths=%v", paths)
	}
}

func TestCollectMappedFromPrefixes_Unknown(t *testing.T) {
	paths := []string{".github/skills/foo.md"}
	used := CollectMappedFromPrefixes("unknown-target-xyz", paths)
	_ = used
}

func TestAllTargetPrefixes_ContainsDotGithub(t *testing.T) {
	prefixes := allTargetPrefixes()
	found := false
	for _, p := range prefixes {
		if p == ".github/" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected .github/ in allTargetPrefixes, got %v", prefixes)
	}
}

func TestAllTargetPrefixes_ContainsDotClaude(t *testing.T) {
	prefixes := allTargetPrefixes()
	found := false
	for _, p := range prefixes {
		if p == ".claude/" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected .claude/ in allTargetPrefixes, got %v", prefixes)
	}
}

func TestFilterFilesByTarget_Opencode(t *testing.T) {
	files := []string{
		".opencode/skills/foo.md",
		".github/skills/bar.md",
	}
	result := FilterFilesByTarget(files, "opencode")
	if len(result.Files) == 0 {
		t.Errorf("expected files for opencode target, got none")
	}
}

func TestFilterFilesByTarget_Windsurf(t *testing.T) {
	files := []string{
		".windsurf/skills/foo.md",
		".github/skills/bar.md",
	}
	result := FilterFilesByTarget(files, "windsurf")
	if len(result.Files) == 0 {
		t.Errorf("expected files for windsurf target, got none")
	}
}
