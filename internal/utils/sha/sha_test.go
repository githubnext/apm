// Package sha_test provides parity tests for short SHA formatting.
// Mirrors src/apm_cli/utils/short_sha.py.
package sha_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/sha"
)

// TestParitySHAFormatShortSHA verifies FormatShortSHA matches the Python
// format_short_sha implementation.
func TestParitySHAFormatShortSHA(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid sha40", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2", "a1b2c3d4"},
		{"valid sha8", "a1b2c3d4", "a1b2c3d4"},
		{"valid sha64", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2a1b2c3d4e5f6a1b2c3d4e5f6", "a1b2c3d4"},
		{"empty", "", ""},
		{"whitespace only", "   ", ""},
		{"sentinel cached", "cached", ""},
		{"sentinel unknown", "unknown", ""},
		{"sentinel cached upper", "CACHED", ""},
		{"too short", "abc123", ""},
		{"exactly 7", "abcdef1", ""},
		{"non hex chars", "xyz12345", ""},
		{"uppercase valid", "ABCDEF12", "ABCDEF12"},
		{"mixed case valid", "aAbBcCdD", "aAbBcCdD"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sha.FormatShortSHA(tc.input)
			if got != tc.expected {
				t.Errorf("FormatShortSHA(%q): got %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}
