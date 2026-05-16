package normalization

import (
	"bytes"
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
