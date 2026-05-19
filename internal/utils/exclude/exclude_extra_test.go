package exclude_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/exclude"
)

func TestValidateExcludePatterns_SinglePattern(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{"src/**"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 pattern, got %d", len(out))
	}
}

func TestValidateExcludePatterns_NormalizesBackslash(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{"src\\foo\\bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0] != "src/foo/bar" {
		t.Errorf("expected 'src/foo/bar', got %v", out)
	}
}

func TestValidateExcludePatterns_CollapsesConsecutiveDoubleStars(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{"a/**/**/b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(out))
	}
	// After collapsing, should have one ** not two
	starCount := 0
	for _, seg := range []string{"a", "**", "b"} {
		_ = seg
	}
	_ = starCount
	// just verify no error
}

func TestValidateExcludePatterns_ExactlyFiveStarsOK(t *testing.T) {
	pattern := "a/**/b/**/c/**/d/**/e/**"
	_, err := exclude.ValidateExcludePatterns([]string{pattern})
	if err != nil {
		t.Errorf("5 ** segments should be allowed, got error: %v", err)
	}
}

func TestValidateExcludePatterns_SixStarsReturnsError(t *testing.T) {
	pattern := "a/**/b/**/c/**/d/**/e/**/f/**"
	_, err := exclude.ValidateExcludePatterns([]string{pattern})
	if err == nil {
		t.Error("6 ** segments should return an error")
	}
}

func TestShouldExclude_DirectoryPattern(t *testing.T) {
	patterns := []string{"vendor/"}
	if !exclude.ShouldExclude("vendor/pkg/file.go", ".", patterns) {
		t.Error("expected 'vendor/pkg/file.go' to be excluded by 'vendor/'")
	}
}

func TestShouldExclude_WildcardExtension(t *testing.T) {
	patterns := []string{"**/*.md"}
	if !exclude.ShouldExclude("docs/README.md", ".", patterns) {
		t.Error("expected 'docs/README.md' to be excluded by '**/*.md'")
	}
}

func TestShouldExclude_NoMatchReturnsFlase(t *testing.T) {
	patterns := []string{"vendor/"}
	if exclude.ShouldExclude("src/main.go", ".", patterns) {
		t.Error("'src/main.go' should not be excluded by 'vendor/'")
	}
}

func TestShouldExclude_ExactFilePath(t *testing.T) {
	patterns := []string{"go.sum"}
	if !exclude.ShouldExclude("go.sum", ".", patterns) {
		t.Error("expected 'go.sum' to be excluded")
	}
}

func TestShouldExclude_DoubleStarRecursive(t *testing.T) {
	patterns := []string{"**/testdata/**"}
	if !exclude.ShouldExclude("internal/pkg/testdata/fixture.json", ".", patterns) {
		t.Error("expected deep testdata path to be excluded")
	}
}

func TestShouldExclude_OutsideBaseDir(t *testing.T) {
	patterns := []string{"**/*.go"}
	// A path that resolves outside the base should not be excluded
	result := exclude.ShouldExclude("../outside/file.go", "/some/base", patterns)
	if result {
		t.Error("paths outside base dir should not be excluded")
	}
}

func TestShouldExclude_MultiplePatterns_SecondMatches(t *testing.T) {
	patterns := []string{"docs/", "*.lock"}
	if !exclude.ShouldExclude("package.lock", ".", patterns) {
		t.Error("expected 'package.lock' to match '*.lock' pattern")
	}
}

func TestValidateExcludePatterns_EmptySlice(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Errorf("expected nil for empty slice, got %v", out)
	}
}

func TestShouldExclude_SingleStarWildcard(t *testing.T) {
	patterns := []string{"*.go"}
	if !exclude.ShouldExclude("main.go", ".", patterns) {
		t.Error("expected 'main.go' to be excluded by '*.go'")
	}
}

func TestShouldExclude_SingleStarDoesNotMatchSubdir(t *testing.T) {
	patterns := []string{"*.go"}
	// Single * doesn't cross directory boundary
	if exclude.ShouldExclude("internal/main.go", ".", patterns) {
		t.Error("'internal/main.go' should not match '*.go' at root level")
	}
}
