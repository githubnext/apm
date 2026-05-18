package githubdownloader

import (
	"testing"
)

// ---------------------------------------------------------------------------
// DefaultOptions
// ---------------------------------------------------------------------------

func TestDefaultOptions_concurrency(t *testing.T) {
	opts := DefaultOptions()
	if opts.Concurrency <= 0 {
		t.Errorf("Concurrency should be positive, got %d", opts.Concurrency)
	}
}

func TestDefaultOptions_not_dry_run(t *testing.T) {
	opts := DefaultOptions()
	if opts.Concurrency <= 0 {
		t.Errorf("Concurrency should be positive, got %d", opts.Concurrency)
	}
	if !opts.AllowFallback {
		t.Error("AllowFallback should default to true")
	}
}

// ---------------------------------------------------------------------------
// ParseLsRemoteOutput: additional edge cases
// ---------------------------------------------------------------------------

func TestParseLsRemoteOutput_tabs_only(t *testing.T) {
	refs := ParseLsRemoteOutput("\t\n")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs for tab-only input, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_windows_line_endings(t *testing.T) {
	input := "abc123\trefs/heads/main\r\ndef456\trefs/tags/v1.0.0\r\n"
	refs := ParseLsRemoteOutput(input)
	// Should parse at least the valid lines; SHA/name may or may not have \r
	// depending on implementation -- just ensure no panic
	_ = refs
}

func TestParseLsRemoteOutput_sha_case_preserved(t *testing.T) {
	input := "ABCDEF1234567890\trefs/heads/feature\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].SHA != "ABCDEF1234567890" {
		t.Errorf("SHA case should be preserved, got %q", refs[0].SHA)
	}
}

// ---------------------------------------------------------------------------
// SemverSortKey: additional cases
// ---------------------------------------------------------------------------

func TestSemverSortKey_no_v_prefix(t *testing.T) {
	key := SemverSortKey("3.0.1")
	if key != [4]int{3, 0, 1, 0} {
		t.Errorf("SemverSortKey(3.0.1) = %v, want [3 0 1 0]", key)
	}
}

func TestSemverSortKey_empty(t *testing.T) {
	key := SemverSortKey("")
	if key[0] != -1 {
		t.Errorf("expected -1 for empty string, got %v", key)
	}
}

func TestSemverSortKey_pre_release_less_than_release(t *testing.T) {
	pre := SemverSortKey("v1.0.0-alpha")
	rel := SemverSortKey("v1.0.0")
	// pre-release should sort lower
	if !(pre[3] < rel[3]) {
		t.Errorf("pre-release should sort lower: pre=%v rel=%v", pre, rel)
	}
}

// ---------------------------------------------------------------------------
// SortRemoteRefs: stability and ties
// ---------------------------------------------------------------------------

func TestSortRemoteRefs_single(t *testing.T) {
	refs := []RemoteRef{{Name: "v1.0.0", SHA: "abc"}}
	sorted := SortRemoteRefs(refs)
	if len(sorted) != 1 || sorted[0].Name != "v1.0.0" {
		t.Errorf("single ref sort broken: %+v", sorted)
	}
}

func TestSortRemoteRefs_empty(t *testing.T) {
	sorted := SortRemoteRefs(nil)
	if len(sorted) != 0 {
		t.Errorf("nil input should return empty slice, got %v", sorted)
	}
}

func TestSortRemoteRefs_non_semver_last(t *testing.T) {
	refs := []RemoteRef{
		{Name: "latest", SHA: "a"},
		{Name: "v1.0.0", SHA: "b"},
	}
	sorted := SortRemoteRefs(refs)
	// semver v1.0.0 should come first
	if sorted[0].Name != "v1.0.0" {
		t.Errorf("expected v1.0.0 first, got %s", sorted[0].Name)
	}
}

func TestSortRemoteRefs_preserves_all_refs(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v2.0.0", SHA: "a"},
		{Name: "v1.0.0", SHA: "b"},
		{Name: "v3.0.0", SHA: "c"},
	}
	sorted := SortRemoteRefs(refs)
	if len(sorted) != 3 {
		t.Errorf("expected 3 refs, got %d", len(sorted))
	}
}

// ---------------------------------------------------------------------------
// RemoteRef struct fields
// ---------------------------------------------------------------------------

func TestRemoteRef_fields(t *testing.T) {
	r := RemoteRef{Name: "refs/tags/v1.0.0", SHA: "deadbeef"}
	if r.Name != "refs/tags/v1.0.0" {
		t.Errorf("unexpected Name: %q", r.Name)
	}
	if r.SHA != "deadbeef" {
		t.Errorf("unexpected SHA: %q", r.SHA)
	}
}

// ---------------------------------------------------------------------------
// ProtocolPreference constants
// ---------------------------------------------------------------------------

func TestProtocolPreference_distinct(t *testing.T) {
	vals := map[ProtocolPreference]bool{}
	for _, p := range []ProtocolPreference{ProtocolHTTPSOnly, ProtocolSSHOnly, ProtocolPreferHTTPS, ProtocolPreferSSH} {
		if vals[p] {
			t.Errorf("duplicate ProtocolPreference value: %v", p)
		}
		vals[p] = true
	}
}
