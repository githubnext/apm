package sha

import (
"testing"
)

func TestFormatShortSHA_ExactlyEightHexChars(t *testing.T) {
result := FormatShortSHA("abcdef12")
if result != "abcdef12" {
t.Errorf("expected 'abcdef12', got %q", result)
}
}

func TestFormatShortSHA_LongHexReturnsFirst8(t *testing.T) {
result := FormatShortSHA("1234567890abcdef")
if result != "12345678" {
t.Errorf("expected '12345678', got %q", result)
}
}

func TestFormatShortSHA_AllUpperCase(t *testing.T) {
result := FormatShortSHA("ABCDEF1234567890")
if result != "ABCDEF12" {
t.Errorf("expected 'ABCDEF12', got %q", result)
}
}

func TestFormatShortSHA_SevenChars_Empty(t *testing.T) {
result := FormatShortSHA("1234567")
if result != "" {
t.Errorf("expected '', got %q", result)
}
}

func TestFormatShortSHA_SentinelCachedMixed(t *testing.T) {
result := FormatShortSHA("Cached")
if result != "" {
t.Errorf("expected '' for 'Cached', got %q", result)
}
}

func TestFormatShortSHA_GCharNotHex(t *testing.T) {
result := FormatShortSHA("abcdefg1")
if result != "" {
t.Errorf("expected '' for non-hex char, got %q", result)
}
}

func TestFormatShortSHA_OnlySpaces(t *testing.T) {
result := FormatShortSHA("        ")
if result != "" {
t.Errorf("expected '' for spaces only, got %q", result)
}
}
