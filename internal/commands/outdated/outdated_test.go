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
