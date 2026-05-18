package hookintegrator

import (
	"testing"
)

// ---------------------------------------------------------------------------
// filterHookFilesForTarget
// ---------------------------------------------------------------------------

func TestFilterHookFilesForTarget_copilot(t *testing.T) {
	files := []string{
		"/pkg/hooks/copilot-hooks.json",
		"/pkg/hooks/cursor-hooks.json",
		"/pkg/hooks/claude-hooks.json",
	}
	got := filterHookFilesForTarget(files, "copilot")
	if len(got) != 1 || got[0] != "/pkg/hooks/copilot-hooks.json" {
		t.Errorf("filterHookFilesForTarget copilot = %v, want [copilot-hooks.json]", got)
	}
}

func TestFilterHookFilesForTarget_cursor(t *testing.T) {
	files := []string{
		"/pkg/hooks/cursor-hooks.json",
		"/pkg/hooks/claude-hooks.json",
	}
	got := filterHookFilesForTarget(files, "cursor")
	if len(got) != 1 || got[0] != "/pkg/hooks/cursor-hooks.json" {
		t.Errorf("filterHookFilesForTarget cursor = %v, want [cursor-hooks.json]", got)
	}
}

func TestFilterHookFilesForTarget_vscode_copilot(t *testing.T) {
	files := []string{"/pkg/hooks/copilot-hooks.json"}
	got := filterHookFilesForTarget(files, "vscode")
	if len(got) != 1 {
		t.Errorf("copilot-hooks should also match vscode, got %v", got)
	}
}

func TestFilterHookFilesForTarget_gemini(t *testing.T) {
	files := []string{
		"/pkg/hooks/gemini-hooks.json",
		"/pkg/hooks/cursor-hooks.json",
	}
	got := filterHookFilesForTarget(files, "gemini")
	if len(got) != 1 || got[0] != "/pkg/hooks/gemini-hooks.json" {
		t.Errorf("filterHookFilesForTarget gemini = %v", got)
	}
}

func TestFilterHookFilesForTarget_universal(t *testing.T) {
	// File with no known suffix should be included for all targets
	files := []string{"/pkg/hooks/myhook.json"}
	for _, target := range []string{"copilot", "cursor", "claude", "codex", "gemini"} {
		got := filterHookFilesForTarget(files, target)
		if len(got) != 1 {
			t.Errorf("universal hook should match target %q, got %v", target, got)
		}
	}
}

func TestFilterHookFilesForTarget_empty(t *testing.T) {
	got := filterHookFilesForTarget(nil, "copilot")
	if len(got) != 0 {
		t.Errorf("empty input should return empty, got %v", got)
	}
}

func TestFilterHookFilesForTarget_windsurf(t *testing.T) {
	files := []string{
		"/pkg/hooks/windsurf-hooks.json",
		"/pkg/hooks/codex-hooks.json",
	}
	got := filterHookFilesForTarget(files, "windsurf")
	if len(got) != 1 || got[0] != "/pkg/hooks/windsurf-hooks.json" {
		t.Errorf("filterHookFilesForTarget windsurf = %v", got)
	}
}

// ---------------------------------------------------------------------------
// shallowCopyMap
// ---------------------------------------------------------------------------

func TestShallowCopyMap_basic(t *testing.T) {
	src := map[string]interface{}{"a": 1, "b": "hello", "c": true}
	dst := shallowCopyMap(src)
	if len(dst) != 3 {
		t.Errorf("expected 3 keys, got %d", len(dst))
	}
	if dst["a"] != 1 || dst["b"] != "hello" || dst["c"] != true {
		t.Errorf("shallow copy values wrong: %v", dst)
	}
}

func TestShallowCopyMap_independence(t *testing.T) {
	src := map[string]interface{}{"x": "original"}
	dst := shallowCopyMap(src)
	dst["x"] = "modified"
	if src["x"] != "original" {
		t.Error("modifying copy should not affect source")
	}
}

func TestShallowCopyMap_empty(t *testing.T) {
	dst := shallowCopyMap(map[string]interface{}{})
	if len(dst) != 0 {
		t.Errorf("shallow copy of empty map should be empty, got %v", dst)
	}
}

// ---------------------------------------------------------------------------
// copilotKeysToGemini
// ---------------------------------------------------------------------------

func TestCopilotKeysToGemini_bashToCommand(t *testing.T) {
	hook := map[string]interface{}{"bash": "echo hello", "event": "preToolUse"}
	copilotKeysToGemini(hook)
	if hook["command"] != "echo hello" {
		t.Errorf("expected command=echo hello, got %v", hook["command"])
	}
	if _, hasBash := hook["bash"]; hasBash {
		t.Error("bash key should be deleted")
	}
}

func TestCopilotKeysToGemini_powershellToCommand(t *testing.T) {
	hook := map[string]interface{}{"powershell": "Write-Host hi"}
	copilotKeysToGemini(hook)
	if hook["command"] != "Write-Host hi" {
		t.Errorf("expected command=Write-Host hi, got %v", hook["command"])
	}
}

func TestCopilotKeysToGemini_commandUnchanged(t *testing.T) {
	hook := map[string]interface{}{"command": "already-set"}
	copilotKeysToGemini(hook)
	if hook["command"] != "already-set" {
		t.Errorf("existing command should not be overwritten, got %v", hook["command"])
	}
}

func TestCopilotKeysToGemini_timeoutSecFloat(t *testing.T) {
	hook := map[string]interface{}{"command": "run", "timeoutSec": float64(5)}
	copilotKeysToGemini(hook)
	if hook["timeout"] != float64(5000) {
		t.Errorf("timeout should be 5000ms, got %v", hook["timeout"])
	}
	if _, has := hook["timeoutSec"]; has {
		t.Error("timeoutSec should be deleted")
	}
}

func TestCopilotKeysToGemini_timeoutSecInt(t *testing.T) {
	hook := map[string]interface{}{"command": "run", "timeoutSec": 10}
	copilotKeysToGemini(hook)
	if hook["timeout"] != 10000 {
		t.Errorf("timeout should be 10000ms, got %v", hook["timeout"])
	}
}

func TestCopilotKeysToGemini_noTimeoutSec(t *testing.T) {
	hook := map[string]interface{}{"command": "run"}
	copilotKeysToGemini(hook)
	if _, has := hook["timeout"]; has {
		t.Error("timeout should not be set when timeoutSec absent")
	}
}

// ---------------------------------------------------------------------------
// deepCopyMap
// ---------------------------------------------------------------------------

func TestDeepCopyMap_basic(t *testing.T) {
	src := map[string]interface{}{"a": "val", "b": 42.0}
	dst := deepCopyMap(src)
	if dst["a"] != "val" || dst["b"] != 42.0 {
		t.Errorf("deepCopyMap values wrong: %v", dst)
	}
}

func TestDeepCopyMap_nested(t *testing.T) {
	src := map[string]interface{}{
		"outer": map[string]interface{}{"inner": "value"},
	}
	dst := deepCopyMap(src)
	inner, ok := dst["outer"].(map[string]interface{})
	if !ok || inner["inner"] != "value" {
		t.Errorf("deepCopyMap nested value wrong: %v", dst)
	}
}

func TestDeepCopyMap_independence(t *testing.T) {
	src := map[string]interface{}{"key": "original"}
	dst := deepCopyMap(src)
	dst["key"] = "modified"
	if src["key"] != "original" {
		t.Error("modifying deep copy should not affect source")
	}
}

// ---------------------------------------------------------------------------
// portableRelpath
// ---------------------------------------------------------------------------

func TestPortableRelpath_simple(t *testing.T) {
	got := portableRelpath("/a/b/c/file.txt", "/a/b")
	if got != "c/file.txt" {
		t.Errorf("portableRelpath = %q, want c/file.txt", got)
	}
}

func TestPortableRelpath_same(t *testing.T) {
	got := portableRelpath("/a/b", "/a/b")
	if got != "." {
		t.Errorf("portableRelpath same = %q, want .", got)
	}
}

// ---------------------------------------------------------------------------
// toSlice
// ---------------------------------------------------------------------------

func TestToSlice_slice(t *testing.T) {
	in := []interface{}{"a", "b", "c"}
	got := toSlice(in)
	if len(got) != 3 {
		t.Errorf("toSlice []interface{} = len %d, want 3", len(got))
	}
}

func TestToSlice_nonSlice(t *testing.T) {
	got := toSlice("notaslice")
	if len(got) != 0 {
		t.Errorf("toSlice non-slice should return empty, got %v", got)
	}
}

func TestToSlice_nil(t *testing.T) {
	got := toSlice(nil)
	if len(got) != 0 {
		t.Errorf("toSlice nil should return empty, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// toGeminiHookEntries
// ---------------------------------------------------------------------------

func TestToGeminiHookEntries_empty(t *testing.T) {
	got := toGeminiHookEntries(nil)
	if len(got) != 0 {
		t.Errorf("toGeminiHookEntries(nil) = %v, want empty", got)
	}
}

func TestToGeminiHookEntries_flat(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{"bash": "echo hi", "event": "preToolUse"},
	}
	got := toGeminiHookEntries(entries)
	if len(got) != 1 {
		t.Errorf("expected 1 result, got %d", len(got))
	}
	outer, ok := got[0].(map[string]interface{})
	if !ok {
		t.Fatal("result should be map")
	}
	hooks, ok := outer["hooks"].([]interface{})
	if !ok || len(hooks) == 0 {
		t.Errorf("result should have hooks: %v", outer)
	}
}

func TestToGeminiHookEntries_alreadyNested(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{
			"hooks": []interface{}{
				map[string]interface{}{"command": "run"},
			},
		},
	}
	got := toGeminiHookEntries(entries)
	if len(got) != 1 {
		t.Errorf("expected 1 result, got %d", len(got))
	}
}

// ---------------------------------------------------------------------------
// hookPrefixList / hasAnyPrefix
// ---------------------------------------------------------------------------

func TestHasAnyPrefix_match(t *testing.T) {
	prefixes := []string{"apm/", "github/"}
	if !hasAnyPrefix("apm/mypackage", prefixes) {
		t.Error("should match apm/ prefix")
	}
}

func TestHasAnyPrefix_noMatch(t *testing.T) {
	prefixes := []string{"apm/", "github/"}
	if hasAnyPrefix("npm/mypackage", prefixes) {
		t.Error("should not match npm/ against apm/,github/")
	}
}

func TestHasAnyPrefix_empty(t *testing.T) {
	if hasAnyPrefix("apm/pkg", nil) {
		t.Error("empty prefix list should never match")
	}
}

// ---------------------------------------------------------------------------
// HookIntegrationResult
// ---------------------------------------------------------------------------

func TestHookIntegrationResult_fields(t *testing.T) {
	r := &HookIntegrationResult{
		FilesIntegrated: 5,
		FilesUpdated:    2,
		FilesSkipped:    1,
		ScriptsCopied:   3,
		TargetPaths:     []string{"/a", "/b"},
	}
	if r.HooksIntegrated() != 5 {
		t.Errorf("HooksIntegrated() = %d, want 5", r.HooksIntegrated())
	}
	if len(r.TargetPaths) != 2 {
		t.Errorf("TargetPaths len = %d, want 2", len(r.TargetPaths))
	}
}

func TestHookIntegrationResult_zero(t *testing.T) {
	r := &HookIntegrationResult{}
	if r.HooksIntegrated() != 0 {
		t.Error("zero HookIntegrationResult should have 0 HooksIntegrated")
	}
}
