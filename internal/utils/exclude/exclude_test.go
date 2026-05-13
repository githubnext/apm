package exclude_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/exclude"
)

func TestValidateExcludePatterns_nil(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns(nil)
	if err != nil || len(out) != 0 {
		t.Errorf("nil input: got %v %v", out, err)
	}
}

func TestValidateExcludePatterns_normal(t *testing.T) {
	patterns := []string{"docs/**", "build/", "*.log"}
	out, err := exclude.ValidateExcludePatterns(patterns)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3, got %d", len(out))
	}
}

func TestValidateExcludePatterns_tooManyStars(t *testing.T) {
	pattern := "a/**/b/**/c/**/d/**/e/**/f/**"
	_, err := exclude.ValidateExcludePatterns([]string{pattern})
	if err == nil {
		t.Error("expected error for too many ** segments")
	}
}

func TestValidateExcludePatterns_collapsesConsecutiveStars(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{"a/**/**/b"})
	if err != nil {
		t.Fatal(err)
	}
	if out[0] != "a/**/b" {
		t.Errorf("expected collapsed pattern, got %s", out[0])
	}
}

func TestShouldExclude_basic(t *testing.T) {
	cases := []struct {
		rel     string
		pattern string
		want    bool
	}{
		{"docs/foo.md", "docs/**", true},
		{"src/main.go", "docs/**", false},
		{"build/out.bin", "build/", true},
		{"log.txt", "*.log", false},
		{"foo.log", "*.log", true},
		{"a/b/c.go", "a/**/c.go", true},
		{"a/x/y/c.go", "a/**/c.go", true},
		{"a/b/d.go", "a/**/c.go", false},
	}
	for _, tc := range cases {
		got := exclude.ShouldExclude("/base/"+tc.rel, "/base", []string{tc.pattern})
		if got != tc.want {
			t.Errorf("ShouldExclude(%q, %q): got %v, want %v", tc.rel, tc.pattern, got, tc.want)
		}
	}
}

func TestShouldExclude_noPatterns(t *testing.T) {
	if exclude.ShouldExclude("/base/file.go", "/base", nil) {
		t.Error("nil patterns should never exclude")
	}
}
