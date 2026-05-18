package outdated

import "testing"

func TestIsTagRef(t *testing.T) {
	valid := []string{"v1.0.0", "v2.3.4", "1.0.0", "0.1.2"}
	for _, v := range valid {
		if !isTagRef(v) {
			t.Errorf("isTagRef(%q) = false, want true", v)
		}
	}
	invalid := []string{"main", "abc123", "feature/x", ""}
	for _, v := range invalid {
		if isTagRef(v) {
			t.Errorf("isTagRef(%q) = true, want false", v)
		}
	}
}

func TestStripV(t *testing.T) {
	tests := []struct{ in, want string }{
		{"v1.0.0", "1.0.0"},
		{"1.0.0", "1.0.0"},
		{"vfoo", "foo"},
		{"", ""},
	}
	for _, tc := range tests {
		got := stripV(tc.in)
		if got != tc.want {
			t.Errorf("stripV(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"v2.0.0", "v1.0.0", 1},
		{"v1.0.0", "v2.0.0", -1},
		{"v1.0.0", "v1.0.0", 0},
		{"v1.2.3", "v1.2.2", 1},
		{"v1.0.0", "v1.0.1", -1},
	}
	for _, tc := range tests {
		got := compareSemver(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("compareSemver(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestLatestSemverTag(t *testing.T) {
	refs := []RemoteRef{
		{Name: "main", IsTag: false},
		{Name: "v1.0.0", IsTag: true},
		{Name: "v2.0.0", IsTag: true},
		{Name: "v1.5.0", IsTag: true},
	}
	got := latestSemverTag(refs)
	if got != "v2.0.0" {
		t.Errorf("latestSemverTag = %q, want %q", got, "v2.0.0")
	}
}

func TestLatestSemverTagEmpty(t *testing.T) {
	refs := []RemoteRef{{Name: "main", IsTag: false}}
	got := latestSemverTag(refs)
	if got != "" {
		t.Errorf("latestSemverTag (no tags) = %q, want empty", got)
	}
}

func TestCompareSemver_PatchDiff(t *testing.T) {
	if compareSemver("v1.0.2", "v1.0.1") != 1 {
		t.Error("expected 1.0.2 > 1.0.1")
	}
	if compareSemver("v1.0.0", "v1.0.3") != -1 {
		t.Error("expected 1.0.0 < 1.0.3")
	}
}

func TestCompareSemver_MinorDiff(t *testing.T) {
	if compareSemver("v1.3.0", "v1.2.9") != 1 {
		t.Error("expected 1.3.0 > 1.2.9")
	}
	if compareSemver("v1.1.0", "v1.2.0") != -1 {
		t.Error("expected 1.1.0 < 1.2.0")
	}
}

func TestIsTagRef_SemverVariants(t *testing.T) {
	valid := []string{"v0.0.1", "v10.0.0", "v1.2.3", "0.0.0", "100.200.300"}
	for _, v := range valid {
		if !isTagRef(v) {
			t.Errorf("isTagRef(%q) should be true", v)
		}
	}
}

func TestStripV_NoPrefix(t *testing.T) {
	cases := []struct{ in, want string }{
		{"1.0.0", "1.0.0"},
		{"abc", "abc"},
		{"v", ""},
	}
	for _, tc := range cases {
		if got := stripV(tc.in); got != tc.want {
			t.Errorf("stripV(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		s    string
		n    int
		want string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"abcde", 5, "abcde"},
		{"abcdef", 6, "abcdef"},
		{"abcdefg", 6, "abc..."},
	}
	for _, c := range cases {
		got := truncate(c.s, c.n)
		if got != c.want {
			t.Errorf("truncate(%q,%d)=%q want %q", c.s, c.n, got, c.want)
		}
	}
}

func TestOutdatedRowFields(t *testing.T) {
	row := OutdatedRow{
		Package: "owner/repo",
		Current: "v1.0.0",
		Latest:  "v2.0.0",
		Status:  "outdated",
		Source:  "github.com",
	}
	if row.Package != "owner/repo" {
		t.Errorf("unexpected Package: %q", row.Package)
	}
	if row.Status != "outdated" {
		t.Errorf("unexpected Status: %q", row.Status)
	}
}

func TestRemoteRefFields(t *testing.T) {
	r := RemoteRef{Name: "v1.2.3", IsTag: true, Commit: "abc123"}
	if !r.IsTag {
		t.Error("expected IsTag=true")
	}
	if r.Commit != "abc123" {
		t.Errorf("unexpected Commit: %q", r.Commit)
	}
}
