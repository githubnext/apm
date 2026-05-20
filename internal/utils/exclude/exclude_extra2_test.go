package exclude_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/exclude"
)

func TestValidateExcludePatterns_EmptySliceReturnsNil(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Errorf("expected nil for empty slice, got %v", out)
	}
}

func TestValidateExcludePatterns_NilSlice(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Errorf("expected nil for nil input, got %v", out)
	}
}

func TestValidateExcludePatterns_MultiplePatterns(t *testing.T) {
	patterns := []string{"a/b/c", "d/e/**", "*.go"}
	out, err := exclude.ValidateExcludePatterns(patterns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 patterns, got %d", len(out))
	}
}

func TestValidateExcludePatterns_RejectsExcessiveDoubleStars(t *testing.T) {
	// 6 non-consecutive ** segments exceed the MaxDoubleStarSegments=5 limit
	pattern := "a/**/b/**/c/**/d/**/e/**/f/**"
	_, err := exclude.ValidateExcludePatterns([]string{pattern})
	if err == nil {
		t.Error("expected error for pattern with 6 ** segments")
	}
	if !strings.Contains(err.Error(), "**") {
		t.Errorf("error should mention **: %v", err)
	}
}

func TestValidateExcludePatterns_ExactlyMaxDoubleStars(t *testing.T) {
	// Exactly MaxDoubleStarSegments (5) ** segments should be valid
	pattern := "**/**/**/**/**"
	out, err := exclude.ValidateExcludePatterns([]string{pattern})
	if err != nil {
		t.Fatalf("unexpected error for exactly max ** segments: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 pattern, got %d", len(out))
	}
}

func TestValidateExcludePatterns_SimpleWildcard(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{"*.py"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0] != "*.py" {
		t.Errorf("expected [*.py], got %v", out)
	}
}

func TestValidateExcludePatterns_PreservesNonGlobPaths(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{"src/apm_cli/utils"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0] != "src/apm_cli/utils" {
		t.Errorf("expected unmodified path, got %v", out)
	}
}

func TestValidateExcludePatterns_MixedValidAndTripleStarOK(t *testing.T) {
	out, err := exclude.ValidateExcludePatterns([]string{"src/**", "tests/**", "docs/**"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 patterns, got %d", len(out))
	}
}
