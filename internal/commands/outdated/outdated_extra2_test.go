package outdated

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripV_WithV(t *testing.T) {
	if got := stripV("v1.2.3"); got != "1.2.3" {
		t.Errorf("expected 1.2.3 got %q", got)
	}
}

func TestStripV_WithoutV(t *testing.T) {
	if got := stripV("1.2.3"); got != "1.2.3" {
		t.Errorf("expected 1.2.3 got %q", got)
	}
}

func TestStripV_Empty(t *testing.T) {
	if got := stripV(""); got != "" {
		t.Errorf("expected empty got %q", got)
	}
}

func TestIsTagRef_ValidSemver(t *testing.T) {
	for _, v := range []string{"v1.0.0", "1.2.3", "v10.20.30"} {
		if !isTagRef(v) {
			t.Errorf("expected %q to be tag ref", v)
		}
	}
}

func TestIsTagRef_SHA(t *testing.T) {
	if isTagRef("abc1234567890abcdef") {
		t.Error("SHA should not be tag ref")
	}
}

func TestIsTagRef_EmptyString(t *testing.T) {
	if isTagRef("") {
		t.Error("empty string should not be tag ref")
	}
}

func TestCompareSemver_PatchDiffNew(t *testing.T) {
	if got := compareSemver("1.0.2", "1.0.1"); got <= 0 {
		t.Errorf("expected positive, got %d", got)
	}
}

func TestCompareSemver_SameVersionNew(t *testing.T) {
	if got := compareSemver("3.3.4", "3.3.4"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestSemverParts_MissingPatch(t *testing.T) {
	parts := semverParts("1.2")
	if len(parts) < 2 {
		t.Errorf("expected at least 2 parts, got %v", parts)
	}
}

func TestTruncate_ExactBoundary(t *testing.T) {
	s := "abcde"
	got := truncate(s, len(s)+1)
	if got != s {
		t.Errorf("expected %q got %q", s, got)
	}
}

func TestTruncate_OneLonger(t *testing.T) {
	got := truncate("hello world", 8)
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected ellipsis, got %q", got)
	}
}

func TestLockFile_EmptyEntries(t *testing.T) {
	lf := &LockFile{}
	if len(lf.Entries) != 0 {
		t.Error("expected empty entries")
	}
}

func TestLockEntry_Fields(t *testing.T) {
	le := LockEntry{
		Name:        "owner/repo",
		LockedRef:   "v1.0.0",
		LockedCommit: "abc123",
		Source:      "github",
	}
	if le.Name != "owner/repo" {
		t.Errorf("unexpected Name: %q", le.Name)
	}
	if le.LockedRef != "v1.0.0" {
		t.Errorf("unexpected LockedRef: %q", le.LockedRef)
	}
}

func TestOutdatedRow_ExtraTags(t *testing.T) {
	row := OutdatedRow{
		Package:   "pkg",
		Current:   "v1.0.0",
		Latest:    "v2.0.0",
		Status:    "outdated",
		ExtraTags: []string{"v1.1.0", "v1.2.0"},
	}
	if len(row.ExtraTags) != 2 {
		t.Errorf("expected 2 extra tags, got %d", len(row.ExtraTags))
	}
}

func TestParseLockFile_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "apm.lock.yaml")
	if err := os.WriteFile(p, []byte("owner/repo:\n  ref: v1.0.0\n  commit: abc123\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseLockFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf == nil {
		t.Fatal("expected non-nil LockFile")
	}
}

func TestCheckOptions_ZeroValue(t *testing.T) {
	var opts CheckOptions
	_ = opts
}

func TestRemoteRef_Fields(t *testing.T) {
	rr := RemoteRef{Name: "refs/tags/v1.0.0", Commit: "deadbeef"}
	if rr.Name == "" || rr.Commit == "" {
		t.Error("fields should not be empty")
	}
}

func TestLatestSemverTag_Multiple(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v1.0.0", Commit: "a", IsTag: true},
		{Name: "v2.0.0", Commit: "b", IsTag: true},
		{Name: "v1.5.0", Commit: "c", IsTag: true},
	}
	got := latestSemverTag(refs)
	if got != "v2.0.0" {
		t.Errorf("expected v2.0.0, got %q", got)
	}
}

func TestLatestSemverTag_WithVAndWithout(t *testing.T) {
	refs := []RemoteRef{
		{Name: "1.0.0", Commit: "a", IsTag: true},
		{Name: "v1.1.0", Commit: "b", IsTag: true},
	}
	got := latestSemverTag(refs)
	if got != "v1.1.0" {
		t.Errorf("expected v1.1.0, got %q", got)
	}
}
