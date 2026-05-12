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
