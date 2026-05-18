package normalization

import (
	"bytes"
	"testing"
)

func TestStripBOM_WithBOM(t *testing.T) {
	input := append([]byte{0xef, 0xbb, 0xbf}, []byte("hello")...)
	got := StripBOM(input)
	if !bytes.Equal(got, []byte("hello")) {
		t.Errorf("StripBOM with BOM: got %q, want %q", got, "hello")
	}
}

func TestStripBOM_WithoutBOM(t *testing.T) {
	input := []byte("hello world")
	got := StripBOM(input)
	if !bytes.Equal(got, input) {
		t.Errorf("StripBOM without BOM: got %q, want %q", got, input)
	}
}

func TestStripBOM_Empty(t *testing.T) {
	got := StripBOM(nil)
	if len(got) != 0 {
		t.Errorf("StripBOM(nil) = %q, want empty", got)
	}
}

func TestStripBOM_OnlyBOM(t *testing.T) {
	input := []byte{0xef, 0xbb, 0xbf}
	got := StripBOM(input)
	if len(got) != 0 {
		t.Errorf("StripBOM of BOM-only: expected empty, got %q", got)
	}
}

func TestStripBOM_BOMPlusNewline(t *testing.T) {
	input := append([]byte{0xef, 0xbb, 0xbf}, []byte("\n")...)
	got := StripBOM(input)
	if !bytes.Equal(got, []byte("\n")) {
		t.Errorf("StripBOM BOM+newline: got %q, want %q", got, "\n")
	}
}

func TestNormalizeLineEndings_CRLF(t *testing.T) {
	in := []byte("a\r\nb\r\nc\r\n")
	want := []byte("a\nb\nc\n")
	got := NormalizeLineEndings(in)
	if !bytes.Equal(got, want) {
		t.Errorf("NormalizeLineEndings: got %q, want %q", got, want)
	}
}

func TestNormalizeLineEndings_LFOnly(t *testing.T) {
	in := []byte("a\nb\n")
	got := NormalizeLineEndings(in)
	if !bytes.Equal(got, in) {
		t.Errorf("NormalizeLineEndings LF-only: modified unexpectedly: %q", got)
	}
}

func TestNormalizeLineEndings_Empty(t *testing.T) {
	got := NormalizeLineEndings(nil)
	if len(got) != 0 {
		t.Errorf("NormalizeLineEndings(nil) = %q, want empty", got)
	}
}

func TestNormalizeLineEndings_NoCRLF(t *testing.T) {
	in := []byte("no line endings at all")
	got := NormalizeLineEndings(in)
	if !bytes.Equal(got, in) {
		t.Errorf("NormalizeLineEndings no CRLF: got %q, want %q", got, in)
	}
}

func TestNormalize_BOMAndCRLF(t *testing.T) {
	input := append([]byte{0xef, 0xbb, 0xbf}, []byte("line1\r\nline2\r\n")...)
	got := Normalize(input)
	want := []byte("line1\nline2\n")
	if !bytes.Equal(got, want) {
		t.Errorf("Normalize BOM+CRLF: got %q, want %q", got, want)
	}
}

func TestNormalize_WithBuildID(t *testing.T) {
	input := []byte("<!-- Build ID: abc123 -->\ncontent\n")
	got := Normalize(input)
	if bytes.Contains(got, []byte("Build ID")) {
		t.Errorf("Normalize should strip Build ID: %q", got)
	}
	if !bytes.Contains(got, []byte("content")) {
		t.Errorf("Normalize should preserve content: %q", got)
	}
}

func TestNormalize_Clean(t *testing.T) {
	input := []byte("clean content\n")
	got := Normalize(input)
	if !bytes.Equal(got, input) {
		t.Errorf("Normalize clean: got %q, want %q", got, input)
	}
}

func TestNormalize_Empty(t *testing.T) {
	got := Normalize(nil)
	if len(got) != 0 {
		t.Errorf("Normalize(nil) = %q, want empty", got)
	}
}

func TestStripBuildID_Multiple(t *testing.T) {
	input := []byte("<!-- Build ID: aaa111 -->\nline\n<!-- Build ID: bbb222 -->\nend\n")
	got := StripBuildID(input)
	if bytes.Contains(got, []byte("Build ID")) {
		t.Errorf("multiple Build IDs should be stripped: %q", got)
	}
	if !bytes.Contains(got, []byte("line")) {
		t.Errorf("content between Build IDs should be preserved: %q", got)
	}
}

func TestStripBuildID_CaseInsensitive(t *testing.T) {
	input := []byte("<!-- build id: cafebabe -->\ncontent\n")
	got := StripBuildID(input)
	if bytes.Contains(got, []byte("cafebabe")) {
		t.Errorf("case-insensitive Build ID not stripped: %q", got)
	}
}

func TestStripBuildID_NoMatch(t *testing.T) {
	input := []byte("no build id here\n")
	got := StripBuildID(input)
	if !bytes.Equal(got, input) {
		t.Errorf("StripBuildID no match: got %q, want %q", got, input)
	}
}

func TestStripBuildID_Empty(t *testing.T) {
	got := StripBuildID(nil)
	if len(got) != 0 {
		t.Errorf("StripBuildID(nil) = %q, want empty", got)
	}
}

func TestNormalize_CRLFWithBuildID(t *testing.T) {
	input := []byte("<!-- Build ID: deadbeef -->\r\ncontent\r\n")
	got := Normalize(input)
	if bytes.Contains(got, []byte("deadbeef")) {
		t.Errorf("Build ID should be stripped: %q", got)
	}
	if bytes.Contains(got, []byte("\r\n")) {
		t.Errorf("CRLF should be normalized: %q", got)
	}
}

func TestNormalize_IdempotentOnClean(t *testing.T) {
	input := []byte("line1\nline2\nline3\n")
	got1 := Normalize(input)
	got2 := Normalize(got1)
	if !bytes.Equal(got1, got2) {
		t.Errorf("Normalize not idempotent: %q vs %q", got1, got2)
	}
}
