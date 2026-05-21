package outdated

import (
	"testing"
)

func TestOutdatedRowZeroValue_Extra3(t *testing.T) {
	var row OutdatedRow
	if row.Package != "" {
		t.Errorf("expected empty Package, got %q", row.Package)
	}
}

func TestCheckOptionsDefaults_Extra3(t *testing.T) {
	opts := CheckOptions{}
	if opts.Verbose {
		t.Error("Verbose should default to false")
	}
	if opts.Format != "" {
		t.Errorf("Format should default to empty string, got %q", opts.Format)
	}
}

func TestCheckResultZeroValue_Extra3(t *testing.T) {
	var result CheckResult
	if result.Rows != nil {
		t.Error("Rows should be nil by default")
	}
}

func TestLockEntryZeroValue_Extra3(t *testing.T) {
	var e LockEntry
	if e.Name != "" {
		t.Errorf("Name should be empty, got %q", e.Name)
	}
	if e.LockedRef != "" {
		t.Errorf("LockedRef should be empty, got %q", e.LockedRef)
	}
}

func TestLockFileZeroValue_Extra3(t *testing.T) {
	var lf LockFile
	if lf.Entries != nil {
		t.Error("Entries should be nil by default")
	}
}

func TestRemoteRefZeroValue_Extra3(t *testing.T) {
	var ref RemoteRef
	if ref.Name != "" {
		t.Errorf("Name should be empty, got %q", ref.Name)
	}
	if ref.Commit != "" {
		t.Errorf("Commit should be empty, got %q", ref.Commit)
	}
}

func TestIsTagRef_AnnotatedTag_Extra3(t *testing.T) {
	cases := []struct {
		ref  string
		want bool
	}{
		{"v1.0.0", true},
		{"v2.3.4", true},
		{"main", false},
		{"abc123", false},
		{"", false},
	}
	for _, c := range cases {
		got := isTagRef(c.ref)
		if got != c.want {
			t.Errorf("isTagRef(%q) = %v, want %v", c.ref, got, c.want)
		}
	}
}

func TestStripV_WithDoubleV_Extra3(t *testing.T) {
	got := stripV("vv1.0")
	if got != "v1.0" {
		t.Errorf("stripV(\"vv1.0\") = %q, want v1.0", got)
	}
}

func TestCompareSemver_PatchOnly_Extra3(t *testing.T) {
	if compareSemver("1.0.1", "1.0.0") <= 0 {
		t.Error("1.0.1 should be greater than 1.0.0")
	}
}

func TestCompareSemver_LessMinor_Extra3(t *testing.T) {
	if compareSemver("1.0.0", "1.1.0") >= 0 {
		t.Error("1.0.0 should be less than 1.1.0")
	}
}

func TestTruncate_AtBoundary_Extra3(t *testing.T) {
	s := "hello"
	got := truncate(s, 5)
	if got != s {
		t.Errorf("truncate(%q, 5) = %q, want %q", s, got, s)
	}
}
