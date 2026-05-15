package semver

import (
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
		major   int
		minor   int
		patch   int
		pre     string
	}{
		{"1.2.3", false, 1, 2, 3, ""},
		{"0.0.1", false, 0, 0, 1, ""},
		{"10.20.30", false, 10, 20, 30, ""},
		{"1.2.3-alpha.1", false, 1, 2, 3, "alpha.1"},
		{"1.2.3-beta+build.1", false, 1, 2, 3, "beta"},
		{"invalid", true, 0, 0, 0, ""},
		{"1.2", true, 0, 0, 0, ""},
		{"", true, 0, 0, 0, ""},
	}
	for _, tc := range cases {
		v, err := Parse(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("Parse(%q): expected error", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if v.Major != tc.major || v.Minor != tc.minor || v.Patch != tc.patch {
			t.Errorf("Parse(%q): got %d.%d.%d, want %d.%d.%d", tc.input,
				v.Major, v.Minor, v.Patch, tc.major, tc.minor, tc.patch)
		}
		if v.Prerelease != tc.pre {
			t.Errorf("Parse(%q) prerelease: got %q, want %q", tc.input, v.Prerelease, tc.pre)
		}
	}
}

func TestCompare(t *testing.T) {
	mustParse := func(s string) SemVer {
		v, err := Parse(s)
		if err != nil {
			t.Fatalf("Parse(%q): %v", s, err)
		}
		return v
	}
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0", "1.0.0-alpha", 1},
		{"1.0.0-alpha", "1.0.0-beta", -1},
		{"1.0.0-beta", "1.0.0-alpha", 1},
	}
	for _, tc := range cases {
		a, b := mustParse(tc.a), mustParse(tc.b)
		got := a.Compare(b)
		if got != tc.want {
			t.Errorf("Compare(%q, %q): got %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSatisfiesRange(t *testing.T) {
	mustParse := func(s string) SemVer {
		v, err := Parse(s)
		if err != nil {
			t.Fatalf("Parse(%q): %v", s, err)
		}
		return v
	}
	cases := []struct {
		version string
		rangeS  string
		want    bool
	}{
		{"1.2.3", "*", true},
		{"1.2.3", "1.2.3", true},
		{"1.2.4", "1.2.3", false},
		{"1.2.3", "^1.0.0", true},
		{"1.2.3", "^2.0.0", false},
		{"1.2.3", "~1.2.0", true},
		{"1.3.0", "~1.2.0", false},
		{"1.2.3", ">=1.2.0", true},
		{"1.1.9", ">=1.2.0", false},
		{"1.2.3", ">1.2.2", true},
		{"1.2.3", ">1.2.3", false},
		{"1.2.3", "<=1.2.3", true},
		{"1.2.4", "<=1.2.3", false},
		{"1.2.2", "<1.2.3", true},
		{"1.2.3", "<1.2.3", false},
		{"1.2.3", "1.2.x", true},
		{"1.3.0", "1.2.x", true},  // 1.2.x matches any same-major version
		{"2.0.0", "1.2.x", false}, // different major does not match
	}
	for _, tc := range cases {
		v := mustParse(tc.version)
		got := SatisfiesRange(v, tc.rangeS)
		if got != tc.want {
			t.Errorf("SatisfiesRange(%q, %q): got %v, want %v", tc.version, tc.rangeS, got, tc.want)
		}
	}
}

func TestSatisfiesRangeEmpty(t *testing.T) {
	v, _ := Parse("1.0.0")
	if !SatisfiesRange(v, "") {
		t.Error("empty range should match everything")
	}
}

func TestSatisfiesRangeAnd(t *testing.T) {
	v, _ := Parse("1.5.0")
	if !SatisfiesRange(v, ">=1.0.0 <=2.0.0") {
		t.Error("1.5.0 should satisfy >=1.0.0 <=2.0.0")
	}
	if SatisfiesRange(v, ">=1.0.0 <=1.4.9") {
		t.Error("1.5.0 should not satisfy >=1.0.0 <=1.4.9")
	}
}
