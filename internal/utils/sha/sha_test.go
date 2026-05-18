// Package sha_test tests the SHA short-form helper.
package sha_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/sha"
)

func TestFormatShortSHA(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"cached", ""},
		{"unknown", ""},
		{"CACHED", ""},
		{"abc123", ""},           // too short
		{"abc12345", "abc12345"}, // exactly 8 hex chars
		{"abc123456789abcd", "abc12345"},
		{"xyz12345", ""},         // non-hex char
		{"  abc12345  ", "abc12345"}, // trims whitespace
	}
	for _, tt := range tests {
		got := sha.FormatShortSHA(tt.input)
		if got != tt.want {
			t.Errorf("FormatShortSHA(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatShortSHA_AllHexChars(t *testing.T) {
	// All valid hex digits should be accepted.
	validHexSHA := "0123456789abcdef"
	got := sha.FormatShortSHA(validHexSHA)
	if got != "01234567" {
		t.Errorf("FormatShortSHA(%q) = %q, want 01234567", validHexSHA, got)
	}
}

func TestFormatShortSHA_UppercaseHex(t *testing.T) {
	input := "ABCDEF1234567890"
	got := sha.FormatShortSHA(input)
	if got != "ABCDEF12" {
		t.Errorf("FormatShortSHA(%q) = %q, want ABCDEF12", input, got)
	}
}

func TestFormatShortSHA_MixedCase(t *testing.T) {
	input := "aAbBcCdDeEfF0011"
	got := sha.FormatShortSHA(input)
	if got != "aAbBcCdD" {
		t.Errorf("FormatShortSHA(%q) = %q, want aAbBcCdD", input, got)
	}
}

func TestFormatShortSHA_SentinelLowercase(t *testing.T) {
	for _, s := range []string{"cached", "unknown"} {
		got := sha.FormatShortSHA(s)
		if got != "" {
			t.Errorf("FormatShortSHA(%q) = %q, want empty (sentinel)", s, got)
		}
	}
}

func TestFormatShortSHA_SentinelMixedCase(t *testing.T) {
	for _, s := range []string{"CACHED", "UNKNOWN", "Cached", "Unknown"} {
		got := sha.FormatShortSHA(s)
		if got != "" {
			t.Errorf("FormatShortSHA(%q) = %q, want empty (sentinel case-insensitive)", s, got)
		}
	}
}

func TestFormatShortSHA_TooShort(t *testing.T) {
	for _, s := range []string{"a", "ab", "abc", "abcd", "abcde", "abcdef", "abcdefg"} {
		got := sha.FormatShortSHA(s)
		if got != "" {
			t.Errorf("FormatShortSHA(%q) = %q, want empty (too short)", s, got)
		}
	}
}

func TestFormatShortSHA_NonHexChars(t *testing.T) {
	for _, s := range []string{
		"ghijklmn", // g-n are invalid hex
		"xyz12345",
		"!@#$%^&*",
		"12345678!",
		"hello123",
	} {
		got := sha.FormatShortSHA(s)
		if got != "" {
			t.Errorf("FormatShortSHA(%q) = %q, want empty (invalid hex)", s, got)
		}
	}
}

func TestFormatShortSHA_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"  abc12345  ", "abc12345"},
		{"\tabc12345\n", "abc12345"},
		{"abc12345", "abc12345"},
		{"  ", ""},
	}
	for _, tc := range tests {
		got := sha.FormatShortSHA(tc.input)
		if got != tc.want {
			t.Errorf("FormatShortSHA(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFormatShortSHA_ExactlyEightChars(t *testing.T) {
	input := "deadbeef"
	got := sha.FormatShortSHA(input)
	if got != "deadbeef" {
		t.Errorf("FormatShortSHA(%q) = %q, want deadbeef", input, got)
	}
}

func TestFormatShortSHA_TruncatesLongSHA(t *testing.T) {
	full := "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4"
	got := sha.FormatShortSHA(full)
	if len(got) != 8 {
		t.Errorf("FormatShortSHA long SHA: len = %d, want 8", len(got))
	}
	if got != full[:8] {
		t.Errorf("FormatShortSHA long SHA: %q, want %q", got, full[:8])
	}
}
