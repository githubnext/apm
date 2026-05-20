package normalization

import (
	"bytes"
	"strings"
	"testing"
)

func TestStripBuildID_SingleLine(t *testing.T) {
	input := []byte("<!-- Build ID: abc123 -->\nrest of content\n")
	got := StripBuildID(input)
	if bytes.Contains(got, []byte("Build ID")) {
		t.Errorf("StripBuildID should remove Build ID comment, got: %q", got)
	}
	if !bytes.Contains(got, []byte("rest of content")) {
		t.Errorf("StripBuildID should keep rest of content, got: %q", got)
	}
}

func TestStripBuildID_CaseInsensitiveHeader(t *testing.T) {
	input := []byte("<!-- build id: ABCDEF -->\ncontent\n")
	got := StripBuildID(input)
	if bytes.Contains(got, []byte("build id")) {
		t.Errorf("StripBuildID should be case-insensitive, got: %q", got)
	}
}

func TestStripBuildID_NoComment(t *testing.T) {
	input := []byte("no build id comment here\n")
	got := StripBuildID(input)
	if !bytes.Equal(got, input) {
		t.Errorf("StripBuildID with no comment should return unchanged, got: %q", got)
	}
}

func TestStripBuildID_EmptyInput(t *testing.T) {
	got := StripBuildID([]byte{})
	if len(got) != 0 {
		t.Errorf("StripBuildID(empty) = %q, want empty", got)
	}
}

func TestNormalizeLineEndings_MixedCRLFAndLF(t *testing.T) {
	input := []byte("line1\r\nline2\nline3\r\n")
	got := NormalizeLineEndings(input)
	want := []byte("line1\nline2\nline3\n")
	if !bytes.Equal(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeLineEndings_AllCRLF(t *testing.T) {
	input := []byte("a\r\nb\r\nc\r\n")
	got := NormalizeLineEndings(input)
	if bytes.Contains(got, []byte("\r\n")) {
		t.Errorf("expected no CRLF after normalization, got: %q", got)
	}
}

func TestNormalizeLineEndings_OnlyLF(t *testing.T) {
	input := []byte("a\nb\nc\n")
	got := NormalizeLineEndings(input)
	if !bytes.Equal(got, input) {
		t.Errorf("LF-only input should be unchanged, got: %q", got)
	}
}

func TestNormalizeLineEndings_EmptyBytes(t *testing.T) {
	got := NormalizeLineEndings([]byte{})
	if len(got) != 0 {
		t.Errorf("NormalizeLineEndings(empty) = %q, want empty", got)
	}
}

func TestNormalize_AllTransformations(t *testing.T) {
	bom := []byte{0xef, 0xbb, 0xbf}
	input := append(bom, []byte("<!-- Build ID: deadbeef -->\r\ncontent here\r\n")...)
	got := Normalize(input)
	if bytes.HasPrefix(got, bom) {
		t.Error("BOM should be stripped")
	}
	if bytes.Contains(got, []byte("Build ID")) {
		t.Error("Build ID comment should be stripped")
	}
	if bytes.Contains(got, []byte("\r\n")) {
		t.Error("CRLF should be normalized")
	}
	if !bytes.Contains(got, []byte("content here")) {
		t.Error("content should be preserved")
	}
}

func TestNormalize_NilInput(t *testing.T) {
	got := Normalize(nil)
	if len(got) != 0 {
		t.Errorf("Normalize(nil) = %q, want empty", got)
	}
}

func TestStripBuildID_MultilineDocument(t *testing.T) {
	lines := []string{
		"# Header",
		"<!-- Build ID: cafebabe -->",
		"## Section",
		"Some content",
	}
	input := []byte(strings.Join(lines, "\n") + "\n")
	got := StripBuildID(input)
	if bytes.Contains(got, []byte("Build ID")) {
		t.Errorf("Build ID should be stripped: %q", got)
	}
	if !bytes.Contains(got, []byte("# Header")) {
		t.Errorf("Header should be preserved: %q", got)
	}
	if !bytes.Contains(got, []byte("Some content")) {
		t.Errorf("Content should be preserved: %q", got)
	}
}
