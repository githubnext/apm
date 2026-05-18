package commandintegrator

import (
	"testing"
)

func TestParseFrontmatter_WithFrontmatter(t *testing.T) {
	content := "---\ndescription: test desc\nmodel: gpt-4\n---\n\n# Body"
	meta, body := parseFrontmatter(content)
	if meta["description"] != "test desc" {
		t.Errorf("description: got %v", meta["description"])
	}
	if meta["model"] != "gpt-4" {
		t.Errorf("model: got %v", meta["model"])
	}
	if body == "" {
		t.Error("expected non-empty body")
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	content := "# Just a heading\nSome content"
	meta, body := parseFrontmatter(content)
	if len(meta) != 0 {
		t.Errorf("expected empty meta for content without frontmatter, got %v", meta)
	}
	if body != content {
		t.Errorf("body should be the original content")
	}
}

func TestParseFrontmatter_EmptyFrontmatter(t *testing.T) {
	content := "---\n---\n\nbody here"
	meta, _ := parseFrontmatter(content)
	if len(meta) != 0 {
		t.Errorf("expected empty meta for empty frontmatter block, got %v", meta)
	}
}

func TestParseFrontmatter_OnlyFrontmatter(t *testing.T) {
	content := "---\ndescription: only front\n---"
	meta, _ := parseFrontmatter(content)
	if meta["description"] != "only front" {
		t.Errorf("unexpected description: %v", meta["description"])
	}
}

func TestBuildCommandContent_WithDescription(t *testing.T) {
	meta := map[string]interface{}{
		"description": "My command description",
	}
	body := "Do something useful"
	result := buildCommandContent(meta, body)
	if result == "" {
		t.Error("expected non-empty result")
	}
	if result[:3] != "---" {
		t.Errorf("expected result to start with ---, got %q", result[:3])
	}
}

func TestBuildCommandContent_EmptyMeta(t *testing.T) {
	result := buildCommandContent(map[string]interface{}{}, "body text")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestBuildCommandContent_WithAllowedTools(t *testing.T) {
	meta := map[string]interface{}{
		"description":   "Test",
		"allowed-tools": "bash,python",
	}
	result := buildCommandContent(meta, "# Body")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestExtractInputNamesListOfMaps(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{"name": "arg1"},
		map[string]interface{}{"name": "arg2"},
		"arg3",
	}
	valid, rejected := extractInputNames(input)
	if len(rejected) != 0 {
		t.Errorf("unexpected rejected: %v", rejected)
	}
	found := map[string]bool{}
	for _, v := range valid {
		found[v] = true
	}
	if !found["arg1"] || !found["arg2"] || !found["arg3"] {
		t.Errorf("expected arg1,arg2,arg3 in valid: %v", valid)
	}
}

func TestExtractInputNamesInvalidInSlice(t *testing.T) {
	input := []interface{}{"valid-name", "1invalid", "another-valid"}
	valid, rejected := extractInputNames(input)
	found := map[string]bool{}
	for _, v := range valid {
		found[v] = true
	}
	if !found["valid-name"] || !found["another-valid"] {
		t.Errorf("expected valid names: %v", valid)
	}
	if len(rejected) != 1 || rejected[0] != "1invalid" {
		t.Errorf("expected rejected [1invalid], got %v", rejected)
	}
}

func TestIsValidInputName_BoundaryLength(t *testing.T) {
	// 64 chars including first letter = valid
	name64 := "A" + make64chars()
	if !isValidInputName(name64) {
		t.Errorf("64-char name should be valid")
	}
	// 65 chars = invalid
	name65 := "A" + make64chars() + "x"
	if isValidInputName(name65) {
		t.Errorf("65-char name should be invalid")
	}
}

func make64chars() string {
	s := ""
	for i := 0; i < 63; i++ {
		s += "a"
	}
	return s
}

func TestPreservedCommandKeys(t *testing.T) {
	expected := []string{"description", "allowed-tools", "allowedTools", "model", "argument-hint", "argumentHint", "input"}
	for _, k := range expected {
		if !preservedCommandKeys[k] {
			t.Errorf("key %q should be in preservedCommandKeys", k)
		}
	}
}

func TestIntegrationResult_ZeroValue(t *testing.T) {
	r := IntegrationResult{}
	if r.FilesIntegrated != 0 || r.FilesUpdated != 0 || r.FilesSkipped != 0 || r.LinksResolved != 0 {
		t.Errorf("unexpected non-zero fields: %+v", r)
	}
	if len(r.TargetPaths) != 0 {
		t.Error("expected empty TargetPaths")
	}
}

func TestIntegrationResult_AllFields(t *testing.T) {
	r := IntegrationResult{
		FilesIntegrated: 3,
		FilesUpdated:    1,
		FilesSkipped:    2,
		TargetPaths:     []string{"/path/a", "/path/b"},
		LinksResolved:   5,
	}
	if r.FilesIntegrated != 3 || r.LinksResolved != 5 || len(r.TargetPaths) != 2 {
		t.Errorf("unexpected fields: %+v", r)
	}
}

func TestNew_NotNil(t *testing.T) {
	ci := New()
	if ci == nil {
		t.Error("New() returned nil")
	}
}
