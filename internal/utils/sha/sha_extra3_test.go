package sha

import (
	"strings"
	"testing"
)

func TestFormatShortSHA_ValidExactly8(t *testing.T) {
	got := FormatShortSHA("abcdef01")
	if got != "abcdef01" {
		t.Errorf("expected abcdef01, got %q", got)
	}
}

func TestFormatShortSHA_ValidLongerThan8(t *testing.T) {
	got := FormatShortSHA("abcdef0123456789")
	if got != "abcdef01" {
		t.Errorf("expected abcdef01, got %q", got)
	}
}

func TestFormatShortSHA_SpacesTrimmed(t *testing.T) {
	got := FormatShortSHA("  abcdef01  ")
	if got != "abcdef01" {
		t.Errorf("expected abcdef01, got %q", got)
	}
}

func TestFormatShortSHA_7Chars(t *testing.T) {
	if FormatShortSHA("abcdef0") != "" {
		t.Error("7-char hex string should return empty")
	}
}

func TestFormatShortSHA_8ZeroChars(t *testing.T) {
	got := FormatShortSHA("00000000")
	if got != "00000000" {
		t.Errorf("expected 00000000, got %q", got)
	}
}

func TestFormatShortSHA_SentinelCachedVariant(t *testing.T) {
	for _, s := range []string{"cached", "CACHED", "Cached"} {
		if FormatShortSHA(s) != "" {
			t.Errorf("sentinel %q should return empty", s)
		}
	}
}

func TestFormatShortSHA_SentinelUnknownVariant(t *testing.T) {
	for _, s := range []string{"unknown", "UNKNOWN", "Unknown"} {
		if FormatShortSHA(s) != "" {
			t.Errorf("sentinel %q should return empty", s)
		}
	}
}

func TestFormatShortSHA_NonHexCharsReturnEmpty(t *testing.T) {
	cases := []string{"zzzzzzzz", "abcdefgg", "12345678z", "--------"}
	for _, c := range cases {
		if FormatShortSHA(c) != "" {
			t.Errorf("non-hex %q should return empty", c)
		}
	}
}

func TestFormatShortSHA_UpperAndLowerMixed(t *testing.T) {
	got := FormatShortSHA("AbCdEf01")
	if len(got) != 8 {
		t.Errorf("expected 8-char result, got %q", got)
	}
}

func TestFormatShortSHA_ReturnsFirst8Chars_Long(t *testing.T) {
	input := strings.Repeat("a", 40)
	got := FormatShortSHA(input)
	if got != "aaaaaaaa" {
		t.Errorf("expected aaaaaaaa, got %q", got)
	}
}

func TestFormatShortSHA_PureDigits(t *testing.T) {
	got := FormatShortSHA("12345678")
	if got != "12345678" {
		t.Errorf("expected 12345678, got %q", got)
	}
}
