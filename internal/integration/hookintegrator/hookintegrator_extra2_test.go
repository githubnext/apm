package hookintegrator

import (
	"testing"
)

func TestHookIntegrationResult_HooksIntegrated_Zero(t *testing.T) {
	r := &HookIntegrationResult{}
	if r.HooksIntegrated() != 0 {
		t.Errorf("HooksIntegrated() = %d, want 0", r.HooksIntegrated())
	}
}

func TestHookIntegrationResult_HooksIntegrated_NonZero(t *testing.T) {
	r := &HookIntegrationResult{FilesIntegrated: 3}
	if r.HooksIntegrated() != 3 {
		t.Errorf("HooksIntegrated() = %d, want 3", r.HooksIntegrated())
	}
}

func TestNew_ReturnsNonNil(t *testing.T) {
	hi := New()
	if hi == nil {
		t.Fatal("New() returned nil")
	}
}

func TestFindHookFiles_EmptyDir(t *testing.T) {
	hi := New()
	dir := t.TempDir()
	files := hi.FindHookFiles(dir)
	if len(files) != 0 {
		t.Errorf("expected no hook files, got %v", files)
	}
}

func TestToSlice_Nil(t *testing.T) {
	result := toSlice(nil)
	if len(result) != 0 {
		t.Errorf("toSlice(nil) = %v, want []", result)
	}
}

func TestToSlice_AlreadySlice(t *testing.T) {
	in := []interface{}{"a", "b", "c"}
	out := toSlice(in)
	if len(out) != 3 {
		t.Errorf("toSlice(slice) len = %d, want 3", len(out))
	}
}

func TestToSlice_SingleElement(t *testing.T) {
	// toSlice only returns non-nil for []interface{} - single elements return nil
	out := toSlice("hello")
	if out != nil {
		t.Errorf("toSlice(string) should return nil, got %v", out)
	}
}

func TestDeepCopyMap_Basic_Extra2(t *testing.T) {
	m := map[string]interface{}{"k1": "v1", "k2": "v2"}
	c := deepCopyMap(m)
	if c["k1"] != "v1" || c["k2"] != "v2" {
		t.Errorf("deepCopyMap = %v", c)
	}
}

func TestDeepCopyMap_Independence_Extra2(t *testing.T) {
	m := map[string]interface{}{"key": "original"}
	c := deepCopyMap(m)
	c["key"] = "modified"
	if m["key"] != "original" {
		t.Error("deepCopyMap should be independent copy")
	}
}

func TestDeepCopyMap_NestedMap(t *testing.T) {
	m := map[string]interface{}{
		"nested": map[string]interface{}{"inner": "value"},
	}
	c := deepCopyMap(m)
	inner, ok := c["nested"].(map[string]interface{})
	if !ok {
		t.Fatal("nested key should be a map")
	}
	if inner["inner"] != "value" {
		t.Errorf("inner[inner] = %v", inner["inner"])
	}
}

func TestFilterHookFilesForTarget_Gemini_Extra2(t *testing.T) {
	files := []string{
		"/pkg/hooks/gemini-hooks.json",
		"/pkg/hooks/copilot-hooks.json",
	}
	got := filterHookFilesForTarget(files, "gemini")
	if len(got) != 1 || got[0] != "/pkg/hooks/gemini-hooks.json" {
		t.Errorf("filterHookFilesForTarget gemini = %v", got)
	}
}

func TestFilterHookFilesForTarget_NoMatch(t *testing.T) {
	files := []string{"/pkg/hooks/copilot-hooks.json"}
	got := filterHookFilesForTarget(files, "cursor")
	if len(got) != 0 {
		t.Errorf("expected no match, got %v", got)
	}
}

func TestPortableRelpath_SameDir(t *testing.T) {
	rel := portableRelpath("/a/b/c", "/a/b/c")
	if rel != "." && rel != "" {
		t.Errorf("portableRelpath same dir = %q", rel)
	}
}

func TestPortableRelpath_Child(t *testing.T) {
	rel := portableRelpath("/a/b/c/d", "/a/b")
	if rel != "c/d" {
		t.Errorf("portableRelpath child = %q, want c/d", rel)
	}
}
