package exclude_test

import (
	"strings"
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

func TestValidateExcludePatterns_exactlyMaxStars(t *testing.T) {
// 5 ** segments should be valid (at max limit)
pattern := "a/**/b/**/c/**/d/**/e/**"
out, err := exclude.ValidateExcludePatterns([]string{pattern})
if err != nil {
t.Errorf("expected no error for exactly max ** segments, got %v", err)
}
if len(out) != 1 {
t.Errorf("expected 1 output, got %d", len(out))
}
}

func TestValidateExcludePatterns_backslashNormalized(t *testing.T) {
// Windows-style backslashes should be normalized to forward slashes
pattern := `docs\**\*.md`
out, err := exclude.ValidateExcludePatterns([]string{pattern})
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(out) != 1 || !strings.Contains(out[0], "/") {
t.Errorf("expected normalized pattern with forward slashes, got %q", out[0])
}
}

func TestShouldExclude_multiplePatternsFirstMatch(t *testing.T) {
patterns := []string{"build/**", "dist/**"}
if !exclude.ShouldExclude("/base/build/out.bin", "/base", patterns) {
t.Error("build/out.bin should be excluded by build/**")
}
if !exclude.ShouldExclude("/base/dist/bundle.js", "/base", patterns) {
t.Error("dist/bundle.js should be excluded by dist/**")
}
if exclude.ShouldExclude("/base/src/main.go", "/base", patterns) {
t.Error("src/main.go should not be excluded")
}
}

func TestShouldExclude_exactFilePattern(t *testing.T) {
patterns := []string{"README.md"}
if !exclude.ShouldExclude("/base/README.md", "/base", patterns) {
t.Error("README.md should be excluded by README.md pattern")
}
if exclude.ShouldExclude("/base/docs/README.md", "/base", patterns) {
t.Error("docs/README.md should not be excluded by top-level pattern")
}
}

func TestShouldExclude_emptyPatterns(t *testing.T) {
if exclude.ShouldExclude("/base/anything", "/base", []string{}) {
t.Error("empty patterns should not exclude")
}
}

func TestMaxDoubleStarSegments_value(t *testing.T) {
if exclude.MaxDoubleStarSegments <= 0 {
t.Errorf("MaxDoubleStarSegments should be positive, got %d", exclude.MaxDoubleStarSegments)
}
}
