package normalization

import (
	"bytes"
	"strings"
	"testing"
)

func TestStripBuildID(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"no build id here", "no build id here"},
		{"<!-- Build ID: abc123 -->\nrest", "rest"},
		{"<!-- Build ID: DEADBEEF -->\n", ""},
		{"before\n<!-- Build ID: 1a2b3c -->\nafter", "before\nafter"},
		{"<!-- build id: abc123 -->\n", ""},
	}
	for _, c := range cases {
		got := string(StripBuildID([]byte(c.in)))
		if got != c.want {
			t.Errorf("StripBuildID(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeLineEndings(t *testing.T) {
	in := []byte("line1\r\nline2\r\nline3")
	want := []byte("line1\nline2\nline3")
	got := NormalizeLineEndings(in)
	if !bytes.Equal(got, want) {
		t.Errorf("NormalizeLineEndings(%q) = %q, want %q", in, got, want)
	}
	// Already LF
	lf := []byte("a\nb\n")
	if !bytes.Equal(NormalizeLineEndings(lf), lf) {
		t.Error("NormalizeLineEndings should not alter LF-only content")
	}
}

func TestStripBOM(t *testing.T) {
	withBOM := append([]byte{0xef, 0xbb, 0xbf}, []byte("content")...)
	got := StripBOM(withBOM)
	if !bytes.Equal(got, []byte("content")) {
		t.Errorf("StripBOM should remove BOM, got %q", got)
	}
	noBOM := []byte("no bom")
	if !bytes.Equal(StripBOM(noBOM), noBOM) {
		t.Error("StripBOM should not alter content without BOM")
	}
}

func TestNormalize(t *testing.T) {
	bom := []byte{0xef, 0xbb, 0xbf}
	input := append(bom, []byte("<!-- Build ID: abc123 -->\r\ncontent\r\n")...)
	got := string(Normalize(input))
	want := "content\n"
	if got != want {
		t.Errorf("Normalize(%q) = %q, want %q", input, got, want)
	}
}

func TestStripBuildID_multipleHeaders(t *testing.T) {
	input := []byte("<!-- Build ID: aaa111 -->\n<!-- Build ID: bbb222 -->\nbody\n")
	got := string(StripBuildID(input))
	if strings.Contains(got, "Build ID") {
		t.Errorf("StripBuildID should remove all headers, got %q", got)
	}
	if got != "body\n" {
		t.Errorf("StripBuildID result = %q, want %q", got, "body\n")
	}
}

func TestStripBuildID_noMatch(t *testing.T) {
	input := []byte("no build id header\n")
	got := StripBuildID(input)
	if !bytes.Equal(got, input) {
		t.Errorf("StripBuildID should not alter content without header")
	}
}

func TestNormalizeLineEndings_empty(t *testing.T) {
	if !bytes.Equal(NormalizeLineEndings(nil), []byte(nil)) && !bytes.Equal(NormalizeLineEndings([]byte{}), []byte{}) {
		// Either result is acceptable; just ensure no panic.
	}
}

func TestNormalizeLineEndings_mixedEndings(t *testing.T) {
	in := []byte("line1\r\nline2\nline3\r\n")
	want := []byte("line1\nline2\nline3\n")
	got := NormalizeLineEndings(in)
	if !bytes.Equal(got, want) {
		t.Errorf("NormalizeLineEndings(%q) = %q, want %q", in, got, want)
	}
}

func TestStripBOM_noBOM(t *testing.T) {
	input := []byte("already clean")
	if !bytes.Equal(StripBOM(input), input) {
		t.Error("StripBOM should return identical slice when no BOM")
	}
}

func TestStripBOM_empty(t *testing.T) {
	if !bytes.Equal(StripBOM([]byte{}), []byte{}) {
		t.Error("StripBOM on empty slice should return empty")
	}
}

func TestNormalize_idempotent(t *testing.T) {
	input := []byte("clean content\n")
	once := Normalize(input)
	twice := Normalize(once)
	if !bytes.Equal(once, twice) {
		t.Errorf("Normalize should be idempotent: %q vs %q", once, twice)
	}
}
