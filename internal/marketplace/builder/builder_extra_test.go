package builder

import (
	"encoding/json"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// DefaultBuildOptions
// ---------------------------------------------------------------------------

func TestDefaultBuildOptions_Concurrency(t *testing.T) {
	opts := DefaultBuildOptions()
	if opts.Concurrency != 8 {
		t.Errorf("DefaultBuildOptions().Concurrency = %d, want 8", opts.Concurrency)
	}
}

func TestDefaultBuildOptions_Timeout(t *testing.T) {
	opts := DefaultBuildOptions()
	if opts.TimeoutSeconds != 10.0 {
		t.Errorf("DefaultBuildOptions().TimeoutSeconds = %f, want 10.0", opts.TimeoutSeconds)
	}
}

func TestDefaultBuildOptions_FlagsOff(t *testing.T) {
	opts := DefaultBuildOptions()
	if opts.IncludePrerelease {
		t.Error("IncludePrerelease should default to false")
	}
	if opts.AllowHead {
		t.Error("AllowHead should default to false")
	}
	if opts.ContinueOnError {
		t.Error("ContinueOnError should default to false")
	}
	if opts.Offline {
		t.Error("Offline should default to false")
	}
	if opts.DryRun {
		t.Error("DryRun should default to false")
	}
}

// ---------------------------------------------------------------------------
// ResolveResult.OK
// ---------------------------------------------------------------------------

func TestResolveResult_OK_empty(t *testing.T) {
	r := ResolveResult{}
	if !r.OK() {
		t.Error("empty ResolveResult should be OK")
	}
}

func TestResolveResult_OK_withEntries(t *testing.T) {
	r := ResolveResult{
		Entries: []ResolvedPackage{{Name: "pkg"}},
	}
	if !r.OK() {
		t.Error("ResolveResult with only entries should be OK")
	}
}

func TestResolveResult_OK_withErrors(t *testing.T) {
	r := ResolveResult{
		Errors: [][2]string{{"pkg", "some error"}},
	}
	if r.OK() {
		t.Error("ResolveResult with errors should not be OK")
	}
}

func TestResolveResult_OK_multipleErrors(t *testing.T) {
	r := ResolveResult{
		Errors: [][2]string{
			{"pkg1", "err1"},
			{"pkg2", "err2"},
		},
	}
	if r.OK() {
		t.Error("ResolveResult with multiple errors should not be OK")
	}
}

// ---------------------------------------------------------------------------
// stripRefPrefix
// ---------------------------------------------------------------------------

func TestStripRefPrefix_tag(t *testing.T) {
	got := stripRefPrefix("refs/tags/v1.2.3")
	if got != "v1.2.3" {
		t.Errorf("stripRefPrefix(refs/tags/v1.2.3) = %q, want %q", got, "v1.2.3")
	}
}

func TestStripRefPrefix_head(t *testing.T) {
	got := stripRefPrefix("refs/heads/main")
	if got != "main" {
		t.Errorf("stripRefPrefix(refs/heads/main) = %q, want %q", got, "main")
	}
}

func TestStripRefPrefix_plain(t *testing.T) {
	got := stripRefPrefix("v1.0.0")
	if got != "v1.0.0" {
		t.Errorf("stripRefPrefix(v1.0.0) = %q, want %q", got, "v1.0.0")
	}
}

func TestStripRefPrefix_empty(t *testing.T) {
	got := stripRefPrefix("")
	if got != "" {
		t.Errorf("stripRefPrefix('') = %q, want empty", got)
	}
}

func TestStripRefPrefix_otherRefs(t *testing.T) {
	got := stripRefPrefix("refs/pull/42/head")
	if got != "refs/pull/42/head" {
		t.Errorf("stripRefPrefix(refs/pull/42/head) = %q, want unchanged", got)
	}
}

// ---------------------------------------------------------------------------
// Error types
// ---------------------------------------------------------------------------

func TestBuildError_Error(t *testing.T) {
	e := &BuildError{Msg: "build failed", Package: "mypkg"}
	if e.Error() != "build failed" {
		t.Errorf("BuildError.Error() = %q, want %q", e.Error(), "build failed")
	}
}

func TestHeadNotAllowedError_Error(t *testing.T) {
	e := &HeadNotAllowedError{Package: "mypkg", Ref: "main"}
	msg := e.Error()
	if !strings.Contains(msg, "mypkg") || !strings.Contains(msg, "main") {
		t.Errorf("HeadNotAllowedError.Error() missing pkg/ref: %q", msg)
	}
}

func TestRefNotFoundError_Error(t *testing.T) {
	e := &RefNotFoundError{Package: "mypkg", Ref: "v9.9.9", OwnerRepo: "owner/repo"}
	msg := e.Error()
	if !strings.Contains(msg, "mypkg") || !strings.Contains(msg, "v9.9.9") || !strings.Contains(msg, "owner/repo") {
		t.Errorf("RefNotFoundError.Error() missing details: %q", msg)
	}
}

func TestNoMatchingVersionError_Error(t *testing.T) {
	e := &NoMatchingVersionError{Package: "mypkg", VersionRange: "^2.0.0", Detail: "no tags"}
	msg := e.Error()
	if !strings.Contains(msg, "mypkg") || !strings.Contains(msg, "^2.0.0") {
		t.Errorf("NoMatchingVersionError.Error() missing details: %q", msg)
	}
}

// ---------------------------------------------------------------------------
// extractPluginSHAs
// ---------------------------------------------------------------------------

func TestExtractPluginSHAs_empty(t *testing.T) {
	data := map[string]interface{}{}
	shas := extractPluginSHAs(data)
	if len(shas) != 0 {
		t.Errorf("expected empty, got %v", shas)
	}
}

func TestExtractPluginSHAs_stringSource(t *testing.T) {
	data := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{
				"name":   "myplugin",
				"source": "abc123sha",
			},
		},
	}
	shas := extractPluginSHAs(data)
	if shas["myplugin"] != "abc123sha" {
		t.Errorf("expected abc123sha, got %q", shas["myplugin"])
	}
}

func TestExtractPluginSHAs_mapSourceSha(t *testing.T) {
	data := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{
				"name":   "myplugin",
				"source": map[string]interface{}{"sha": "deadbeef"},
			},
		},
	}
	shas := extractPluginSHAs(data)
	if shas["myplugin"] != "deadbeef" {
		t.Errorf("expected deadbeef, got %q", shas["myplugin"])
	}
}

func TestExtractPluginSHAs_mapSourceCommit(t *testing.T) {
	data := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{
				"name":   "myplugin",
				"source": map[string]interface{}{"commit": "cafebabe"},
			},
		},
	}
	shas := extractPluginSHAs(data)
	if shas["myplugin"] != "cafebabe" {
		t.Errorf("expected cafebabe, got %q", shas["myplugin"])
	}
}

func TestExtractPluginSHAs_multiplePlugins(t *testing.T) {
	data := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{"name": "p1", "source": "sha1"},
			map[string]interface{}{"name": "p2", "source": "sha2"},
			map[string]interface{}{"name": "p3", "source": "sha3"},
		},
	}
	shas := extractPluginSHAs(data)
	if len(shas) != 3 {
		t.Errorf("expected 3 entries, got %d", len(shas))
	}
	if shas["p1"] != "sha1" || shas["p2"] != "sha2" || shas["p3"] != "sha3" {
		t.Errorf("unexpected shas: %v", shas)
	}
}

// ---------------------------------------------------------------------------
// computeDiff
// ---------------------------------------------------------------------------

func TestComputeDiff_nilOld(t *testing.T) {
	newJSON := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{"name": "p1", "source": "sha1"},
			map[string]interface{}{"name": "p2", "source": "sha2"},
		},
	}
	unchanged, added, updated, removed := computeDiff(nil, newJSON)
	if unchanged != 0 || added != 2 || updated != 0 || removed != 0 {
		t.Errorf("computeDiff(nil,...) = %d,%d,%d,%d want 0,2,0,0", unchanged, added, updated, removed)
	}
}

func TestComputeDiff_allUnchanged(t *testing.T) {
	j := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{"name": "p1", "source": "sha1"},
		},
	}
	unchanged, added, updated, removed := computeDiff(j, j)
	if unchanged != 1 || added != 0 || updated != 0 || removed != 0 {
		t.Errorf("computeDiff(same,same) = %d,%d,%d,%d want 1,0,0,0", unchanged, added, updated, removed)
	}
}

func TestComputeDiff_updatedPlugin(t *testing.T) {
	oldJSON := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{"name": "p1", "source": "oldsha"},
		},
	}
	newJSON := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{"name": "p1", "source": "newsha"},
		},
	}
	unchanged, added, updated, removed := computeDiff(oldJSON, newJSON)
	if unchanged != 0 || added != 0 || updated != 1 || removed != 0 {
		t.Errorf("computeDiff(updated) = %d,%d,%d,%d want 0,0,1,0", unchanged, added, updated, removed)
	}
}

func TestComputeDiff_removedPlugin(t *testing.T) {
	oldJSON := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{"name": "p1", "source": "sha1"},
			map[string]interface{}{"name": "p2", "source": "sha2"},
		},
	}
	newJSON := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{"name": "p1", "source": "sha1"},
		},
	}
	unchanged, added, updated, removed := computeDiff(oldJSON, newJSON)
	if unchanged != 1 || added != 0 || updated != 0 || removed != 1 {
		t.Errorf("computeDiff(removed) = %d,%d,%d,%d want 1,0,0,1", unchanged, added, updated, removed)
	}
}

// ---------------------------------------------------------------------------
// serializeJSON
// ---------------------------------------------------------------------------

func TestSerializeJSON_basic(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	b, err := serializeJSON(data)
	if err != nil {
		t.Fatalf("serializeJSON error: %v", err)
	}
	if len(b) == 0 {
		t.Error("serializeJSON returned empty bytes")
	}
	// Should end with newline
	if b[len(b)-1] != '\n' {
		t.Error("serializeJSON should end with newline")
	}
	// Should be valid JSON
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Errorf("serializeJSON output is not valid JSON: %v", err)
	}
}

func TestSerializeJSON_empty(t *testing.T) {
	b, err := serializeJSON(map[string]interface{}{})
	if err != nil {
		t.Fatalf("serializeJSON empty error: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty output for empty map")
	}
}

// ---------------------------------------------------------------------------
// isDisplayVersion (additional cases)
// ---------------------------------------------------------------------------

func TestIsDisplayVersion_sha(t *testing.T) {
	if !isDisplayVersion("abc1234567890") {
		t.Error("SHA should be treated as display version")
	}
}

func TestIsDisplayVersion_v_prefix(t *testing.T) {
	cases := []struct {
		v    string
		want bool
	}{
		{"v1.2.3", true},
		{"V1.0.0", true},
		{"v0.0.1-alpha.1", true},
	}
	for _, c := range cases {
		got := isDisplayVersion(c.v)
		if got != c.want {
			t.Errorf("isDisplayVersion(%q) = %v, want %v", c.v, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// subtractPluginRoot (additional cases)
// ---------------------------------------------------------------------------

func TestSubtractPluginRoot_nested(t *testing.T) {
	got, err := subtractPluginRoot("plugins/root/sub/dir/file.txt", "plugins/root")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "./sub/dir/file.txt" {
		t.Errorf("got %q, want ./sub/dir/file.txt", got)
	}
}

func TestSubtractPluginRoot_mismatch(t *testing.T) {
	_, err := subtractPluginRoot("other/path/file.txt", "plugins/root")
	if err == nil {
		t.Error("expected error for non-matching prefix")
	}
}
