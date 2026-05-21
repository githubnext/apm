package outdated

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSemverParts_Basic(t *testing.T) {
	parts := semverParts("v1.2.3")
	if len(parts) != 3 || parts[0] != 1 || parts[1] != 2 || parts[2] != 3 {
		t.Errorf("unexpected parts: %v", parts)
	}
}

func TestSemverParts_NoV(t *testing.T) {
	parts := semverParts("2.10.5")
	if len(parts) != 3 || parts[0] != 2 || parts[1] != 10 || parts[2] != 5 {
		t.Errorf("unexpected parts: %v", parts)
	}
}

func TestSemverParts_Zeros(t *testing.T) {
	parts := semverParts("v0.0.0")
	if len(parts) != 3 || parts[0] != 0 || parts[1] != 0 || parts[2] != 0 {
		t.Errorf("unexpected parts: %v", parts)
	}
}

func TestCompareSemver_Equal(t *testing.T) {
	if compareSemver("v3.0.0", "v3.0.0") != 0 {
		t.Error("equal versions should return 0")
	}
}

func TestCompareSemver_MajorDiff(t *testing.T) {
	if compareSemver("v3.0.0", "v1.9.9") != 1 {
		t.Error("higher major should return 1")
	}
	if compareSemver("v1.0.0", "v3.0.0") != -1 {
		t.Error("lower major should return -1")
	}
}

func TestCompareSemver_MinorDiffExtra(t *testing.T) {
	if compareSemver("v1.2.0", "v1.1.9") != 1 {
		t.Error("higher minor should return 1")
	}
}

func TestTruncate_Short(t *testing.T) {
	s := truncate("hello", 10)
	if s != "hello" {
		t.Errorf("short string should not be truncated: %q", s)
	}
}

func TestTruncate_Exact(t *testing.T) {
	s := truncate("hello", 5)
	if s != "hello" {
		t.Errorf("exact length should not be truncated: %q", s)
	}
}

func TestTruncate_Long(t *testing.T) {
	s := truncate("hello world!", 8)
	if len(s) != 8 {
		t.Errorf("truncated string should be exactly n chars, got %d: %q", len(s), s)
	}
	if s[5:] != "..." {
		t.Errorf("truncated string should end with ...: %q", s)
	}
}

func TestIsTagRef_Prerelease(t *testing.T) {
	// v1.0.0-rc1 starts with v1.0.0 so should match
	if !isTagRef("v1.0.0-rc1") {
		t.Error("prerelease tag should be recognized as tag ref")
	}
}

func TestIsTagRef_ShortSHA(t *testing.T) {
	if isTagRef("abc1234") {
		t.Error("short SHA should not be tag ref")
	}
}

func TestParseLockFile_Missing(t *testing.T) {
	_, err := ParseLockFile("/nonexistent/path/apm.lock.yaml")
	if err == nil {
		t.Error("expected error for missing lock file")
	}
}

func TestParseLockFile_Empty(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "apm.lock.yaml")
	if err := os.WriteFile(f, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseLockFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(lf.Entries))
	}
}

func TestParseLockFile_WithEntry(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "apm.lock.yaml")
	// top-level key is the package name; indented keys are its fields
	content := "myorg/mypkg:\n  ref: main\n  commit: abc1234567890abcdef\n"
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err := ParseLockFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(lf.Entries))
	}
	if lf.Entries[0].Name != "myorg/mypkg" {
		t.Errorf("unexpected name: %q", lf.Entries[0].Name)
	}
	if lf.Entries[0].LockedRef != "main" {
		t.Errorf("unexpected ref: %q", lf.Entries[0].LockedRef)
	}
}

func TestOutdatedRow_Fields(t *testing.T) {
	row := OutdatedRow{
		Package: "org/pkg",
		Current: "v1.0.0",
		Latest:  "v2.0.0",
		Status:  "outdated",
		Source:  "apm.yml",
	}
	if row.Package != "org/pkg" {
		t.Errorf("unexpected Package: %q", row.Package)
	}
	if row.Status != "outdated" {
		t.Errorf("unexpected Status: %q", row.Status)
	}
}

func TestCheckResult_Fields(t *testing.T) {
	r := &CheckResult{
		Rows: []OutdatedRow{
			{Package: "a/b", Status: "ok"},
		},
		ErrorCount: 0,
	}
	if len(r.Rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(r.Rows))
	}
	if r.ErrorCount != 0 {
		t.Errorf("expected 0 errors, got %d", r.ErrorCount)
	}
}

func TestLatestSemverTag_NoSemver(t *testing.T) {
	refs := []RemoteRef{
		{Name: "main", IsTag: false},
		{Name: "feature/x", IsTag: false},
	}
	got := latestSemverTag(refs)
	if got != "" {
		t.Errorf("expected empty string for no semver tags, got %q", got)
	}
}

func TestLatestSemverTag_Single(t *testing.T) {
	refs := []RemoteRef{{Name: "v1.0.0", IsTag: true}}
	got := latestSemverTag(refs)
	if got != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %q", got)
	}
}
