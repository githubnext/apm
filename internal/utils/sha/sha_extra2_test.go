package sha_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/sha"
)

func TestFormatShortSHA_UppercaseHexInput(t *testing.T) {
	input := "ABCDEF1234567890"
	got := sha.FormatShortSHA(input)
	if got != "ABCDEF12" {
		t.Errorf("got %q, want %q", got, "ABCDEF12")
	}
}

func TestFormatShortSHA_MixedCaseHex(t *testing.T) {
	input := "aAbBcCdDeEfF0123"
	got := sha.FormatShortSHA(input)
	if got != "aAbBcCdD" {
		t.Errorf("got %q, want %q", got, "aAbBcCdD")
	}
}

func TestFormatShortSHA_SentinelCached_LowerCase(t *testing.T) {
	got := sha.FormatShortSHA("cached")
	if got != "" {
		t.Errorf("expected empty for 'cached', got %q", got)
	}
}

func TestFormatShortSHA_SentinelCached_UpperCase(t *testing.T) {
	got := sha.FormatShortSHA("CACHED")
	if got != "" {
		t.Errorf("expected empty for 'CACHED', got %q", got)
	}
}

func TestFormatShortSHA_SentinelUnknown(t *testing.T) {
	got := sha.FormatShortSHA("unknown")
	if got != "" {
		t.Errorf("expected empty for 'unknown', got %q", got)
	}
}

func TestFormatShortSHA_WhitespaceOnly(t *testing.T) {
	got := sha.FormatShortSHA("   ")
	if got != "" {
		t.Errorf("expected empty for whitespace-only, got %q", got)
	}
}

func TestFormatShortSHA_LeadingTrailingWhitespace(t *testing.T) {
	got := sha.FormatShortSHA("  abcdef1234567890  ")
	if got != "abcdef12" {
		t.Errorf("got %q, want %q", got, "abcdef12")
	}
}

func TestFormatShortSHA_NonHexCharacter(t *testing.T) {
	got := sha.FormatShortSHA("ghijklmn")
	if got != "" {
		t.Errorf("expected empty for non-hex chars, got %q", got)
	}
}

func TestFormatShortSHA_ContainsHyphen(t *testing.T) {
	got := sha.FormatShortSHA("abc-def1")
	if got != "" {
		t.Errorf("expected empty for input with hyphen, got %q", got)
	}
}

func TestFormatShortSHA_ReturnsFirst8Chars(t *testing.T) {
	input := "fedcba9876543210fedcba9876543210fedcba98"
	got := sha.FormatShortSHA(input)
	if len(got) != 8 {
		t.Errorf("expected 8 chars, got %d: %q", len(got), got)
	}
	if !strings.HasPrefix(input, got) {
		t.Errorf("returned %q is not a prefix of input %q", got, input)
	}
}

func TestFormatShortSHA_AllZerosInput(t *testing.T) {
	got := sha.FormatShortSHA("0000000000000000")
	if got != "00000000" {
		t.Errorf("got %q, want %q", got, "00000000")
	}
}

func TestFormatShortSHA_AllFs(t *testing.T) {
	got := sha.FormatShortSHA("ffffffffffffffff")
	if got != "ffffffff" {
		t.Errorf("got %q, want %q", got, "ffffffff")
	}
}
