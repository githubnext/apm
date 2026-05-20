package lockfileenrichment

import (
	"strings"
	"testing"
)

func TestFilterFilesResult_ZeroValue(t *testing.T) {
	var r FilterFilesResult
	if r.Files != nil {
		t.Error("expected nil Files")
	}
	if r.PathMappings != nil {
		t.Error("expected nil PathMappings")
	}
}

func TestFilterFilesByTarget_CopilotPrefixes(t *testing.T) {
	files := []string{".github/skills/foo.md", ".claude/agents/bar.md"}
	result := FilterFilesByTarget(files, "copilot")
	if len(result.Files) == 0 {
		t.Error("expected files for copilot target")
	}
	for _, f := range result.Files {
		if !strings.HasPrefix(f, ".github/") {
			t.Errorf("unexpected file for copilot: %q", f)
		}
	}
}

func TestFilterFilesByTarget_EmptyTarget_Empty(t *testing.T) {
	files := []string{".github/skills/foo.md"}
	result := FilterFilesByTarget(files, "")
	_ = result // unknown target returns empty or unchanged - just no panic
}

func TestFilterFilesByTarget_WindsurfTarget(t *testing.T) {
	files := []string{".windsurf/skills/foo.md", ".github/agents/bar.md"}
	result := FilterFilesByTarget(files, "windsurf")
	if len(result.Files) == 0 {
		t.Error("expected at least one file for windsurf target")
	}
}

func TestFilterFilesByTarget_NoMatch(t *testing.T) {
	files := []string{"docs/readme.md"}
	result := FilterFilesByTarget(files, "copilot")
	if len(result.Files) != 0 {
		t.Errorf("expected 0 files for non-matching docs path, got %d", len(result.Files))
	}
}

func TestPackMeta_ZeroValue(t *testing.T) {
	var m PackMeta
	if m.Format != "" {
		t.Errorf("expected empty Format")
	}
	if m.Target != "" {
		t.Errorf("expected empty Target")
	}
}

func TestPackMeta_FieldRoundtrip(t *testing.T) {
	m := PackMeta{
		Format:   "v1",
		Target:   "copilot",
		PackedAt: "2026-01-01T00:00:00Z",
		MappedFrom: []string{"claude"},
	}
	if m.Format != "v1" {
		t.Errorf("unexpected Format %q", m.Format)
	}
	if m.Target != "copilot" {
		t.Errorf("unexpected Target %q", m.Target)
	}
	if len(m.MappedFrom) != 1 {
		t.Errorf("expected 1 MappedFrom, got %d", len(m.MappedFrom))
	}
}

func TestEnrichLockfileForPack_ContainsTarget(t *testing.T) {
	m := PackMeta{Format: "v1", Target: "copilot", PackedAt: "2026-01-01T00:00:00Z"}
	result := EnrichLockfileForPack(m)
	if !strings.Contains(result, "copilot") {
		t.Errorf("expected target in output: %q", result)
	}
}

func TestEnrichLockfileForPack_AutoTimestamp(t *testing.T) {
	m := PackMeta{Format: "v1", Target: "claude"}
	result := EnrichLockfileForPack(m)
	if result == "" {
		t.Error("expected non-empty output")
	}
}

func TestCollectMappedFromPrefixes_Copilot(t *testing.T) {
	paths := []string{".github/skills/foo.md", ".github/agents/bar.md"}
	mapped := CollectMappedFromPrefixes("copilot", paths)
	_ = mapped // just verify no panic
}

func TestCollectMappedFromPrefixes_EmptyPaths(t *testing.T) {
	mapped := CollectMappedFromPrefixes("claude", nil)
	if mapped == nil {
		mapped = []string{}
	}
	if len(mapped) != 0 {
		t.Errorf("expected empty result for nil paths, got %d", len(mapped))
	}
}
