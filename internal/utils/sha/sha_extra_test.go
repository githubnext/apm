package sha_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/sha"
)

func TestFormatShortSHA_ValidFullSHA(t *testing.T) {
	// A typical 40-char commit SHA
	input := "a1b2c3d4e5f67890abcdef1234567890abcdef12"
	got := sha.FormatShortSHA(input)
	if got != "a1b2c3d4" {
		t.Errorf("got %q, want %q", got, "a1b2c3d4")
	}
}

func TestFormatShortSHA_ExactlyEight_AllLower(t *testing.T) {
	input := "abcdef12"
	got := sha.FormatShortSHA(input)
	if got != "abcdef12" {
		t.Errorf("got %q, want %q", got, "abcdef12")
	}
}

func TestFormatShortSHA_SevenChars_TooShort(t *testing.T) {
	input := "abcdef1"
	got := sha.FormatShortSHA(input)
	if got != "" {
		t.Errorf("7-char input should return empty, got %q", got)
	}
}

func TestFormatShortSHA_SentinelCached_Empty(t *testing.T) {
	got := sha.FormatShortSHA("cached")
	if got != "" {
		t.Errorf("'cached' sentinel should return empty, got %q", got)
	}
}

func TestFormatShortSHA_SentinelUnknown_Empty(t *testing.T) {
	got := sha.FormatShortSHA("unknown")
	if got != "" {
		t.Errorf("'unknown' sentinel should return empty, got %q", got)
	}
}

func TestFormatShortSHA_SentinelCachedUppercase(t *testing.T) {
	got := sha.FormatShortSHA("CACHED")
	if got != "" {
		t.Errorf("'CACHED' sentinel (case-insensitive) should return empty, got %q", got)
	}
}

func TestFormatShortSHA_SentinelUnknownMixed(t *testing.T) {
	got := sha.FormatShortSHA("UnKnOwN")
	if got != "" {
		t.Errorf("'UnKnOwN' sentinel should return empty, got %q", got)
	}
}

func TestFormatShortSHA_OnlySpaces(t *testing.T) {
	got := sha.FormatShortSHA("        ")
	if got != "" {
		t.Errorf("whitespace-only should return empty, got %q", got)
	}
}

func TestFormatShortSHA_HasGChar(t *testing.T) {
	// 'g' is not hex
	got := sha.FormatShortSHA("gabcdef1234567")
	if got != "" {
		t.Errorf("non-hex char 'g' should return empty, got %q", got)
	}
}

func TestFormatShortSHA_HasHyphen(t *testing.T) {
	got := sha.FormatShortSHA("abcdef12-56789abc")
	if got != "" {
		t.Errorf("hyphen is not hex; should return empty, got %q", got)
	}
}

func TestFormatShortSHA_LeadingSpaceTrimmed(t *testing.T) {
	got := sha.FormatShortSHA("  abcdef1234567890  ")
	if got != "abcdef12" {
		t.Errorf("leading/trailing spaces should be trimmed; got %q, want %q", got, "abcdef12")
	}
}

func TestFormatShortSHA_AllZeros(t *testing.T) {
	got := sha.FormatShortSHA("0000000000000000")
	if got != "00000000" {
		t.Errorf("got %q, want %q", got, "00000000")
	}
}

func TestFormatShortSHA_AllUpperHex(t *testing.T) {
	got := sha.FormatShortSHA("ABCDEF1234567890")
	if got != "ABCDEF12" {
		t.Errorf("got %q, want %q", got, "ABCDEF12")
	}
}

func TestFormatShortSHA_NineChars(t *testing.T) {
	got := sha.FormatShortSHA("abcdef123")
	if got != "abcdef12" {
		t.Errorf("9-char hex should return first 8; got %q, want %q", got, "abcdef12")
	}
}

func TestFormatShortSHA_Deterministic(t *testing.T) {
	input := "deadbeefcafe1234"
	r1 := sha.FormatShortSHA(input)
	r2 := sha.FormatShortSHA(input)
	if r1 != r2 {
		t.Errorf("non-deterministic: %q vs %q", r1, r2)
	}
}

func TestFormatShortSHA_EmptyString(t *testing.T) {
	got := sha.FormatShortSHA("")
	if got != "" {
		t.Errorf("empty string should return empty, got %q", got)
	}
}
