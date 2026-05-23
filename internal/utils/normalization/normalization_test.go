// Package normalization_test provides parity tests for normalization helpers.
// Mirrors the behaviour of src/apm_cli/utils/normalization.py.
package normalization_test

import (
	"bytes"
	"testing"

	"github.com/githubnext/apm/internal/utils/normalization"
)

// TestParityNormalizationStripBOM verifies BOM stripping matches Python.
func TestParityNormalizationStripBOM(t *testing.T) {
	bom := normalization.BOM
	cases := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{"no bom", []byte("hello"), []byte("hello")},
		{"bom prefix", append(append([]byte{}, bom...), []byte("hello")...), []byte("hello")},
		{"bom only", bom, []byte{}},
		{"bom in middle not stripped", []byte("hel\xef\xbb\xbflo"), []byte("hel\xef\xbb\xbflo")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalization.StripBOM(tc.input)
			if !bytes.Equal(got, tc.expected) {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

// TestParityNormalizationCRLF verifies CRLF-to-LF conversion matches Python.
func TestParityNormalizationCRLF(t *testing.T) {
	cases := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{"no crlf", []byte("hello\nworld\n"), []byte("hello\nworld\n")},
		{"crlf", []byte("hello\r\nworld\r\n"), []byte("hello\nworld\n")},
		{"bare cr preserved", []byte("hello\rworld"), []byte("hello\rworld")},
		{"mixed", []byte("a\r\nb\nc\r\n"), []byte("a\nb\nc\n")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalization.NormalizeLineEndings(tc.input)
			if !bytes.Equal(got, tc.expected) {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

// TestParityNormalizationStripBuildID verifies Build ID header stripping.
func TestParityNormalizationStripBuildID(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"no build id", "hello world", "hello world"},
		{"build id", "<!-- Build ID: abc123 -->\nhello", "hello"},
		{"build id uppercase", "<!-- BUILD ID: abc123 -->\nhello", "hello"},
		{"build id no newline", "<!-- Build ID: abc123 -->hello", "hello"},
		{"build id with spaces", "<!--  Build ID:  abc123  -->\nhello", "hello"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(normalization.StripBuildID([]byte(tc.input)))
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

// TestParityNormalizationNormalize verifies composite Normalize function.
func TestParityNormalizationNormalize(t *testing.T) {
	bom := normalization.BOM
	// Build a payload with BOM + Build ID + CRLF
	input := append(append([]byte{}, bom...), []byte("<!-- Build ID: deadbeef -->\r\nhello\r\nworld\r\n")...)
	expected := []byte("hello\nworld\n")
	got := normalization.Normalize(input)
	if !bytes.Equal(got, expected) {
		t.Errorf("got %q, want %q", got, expected)
	}
}
