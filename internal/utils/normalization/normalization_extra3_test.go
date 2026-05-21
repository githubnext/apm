package normalization

import (
	"bytes"
	"testing"
)

func TestStripBuildID_SingleOccurrence(t *testing.T) {
	in := []byte("<!-- Build ID: abc123def456 -->\nhello\n")
	out := StripBuildID(in)
	if bytes.Contains(out, []byte("Build ID")) {
		t.Errorf("Build ID header should be removed, got: %q", out)
	}
	if !bytes.Contains(out, []byte("hello")) {
		t.Error("non-header content should be preserved")
	}
}

func TestStripBuildID_TwoOccurrences(t *testing.T) {
	in := []byte("<!-- Build ID: aaa -->\n<!-- Build ID: bbb -->\ncontent\n")
	out := StripBuildID(in)
	if bytes.Contains(out, []byte("Build ID")) {
		t.Errorf("all Build ID headers should be removed, got: %q", out)
	}
}

func TestStripBuildID_NilInput(t *testing.T) {
	out := StripBuildID(nil)
	if out != nil && len(out) != 0 {
		t.Errorf("nil input should produce nil/empty, got: %q", out)
	}
}

func TestNormalizeLineEndings_CRLFtoLF(t *testing.T) {
	in := []byte("a\r\nb\r\nc\r\n")
	out := NormalizeLineEndings(in)
	if bytes.Contains(out, []byte("\r\n")) {
		t.Error("CRLF should be replaced")
	}
	if !bytes.Equal(out, []byte("a\nb\nc\n")) {
		t.Errorf("unexpected result: %q", out)
	}
}

func TestNormalizeLineEndings_NoChange(t *testing.T) {
	in := []byte("a\nb\nc\n")
	out := NormalizeLineEndings(in)
	if !bytes.Equal(out, in) {
		t.Errorf("LF-only content should be unchanged, got: %q", out)
	}
}

func TestNormalizeLineEndings_NilInput(t *testing.T) {
	out := NormalizeLineEndings(nil)
	if len(out) != 0 {
		t.Errorf("nil input should produce empty, got: %q", out)
	}
}

func TestStripBOM_WithBOM_Removed(t *testing.T) {
	bom := []byte{0xef, 0xbb, 0xbf}
	in := append(bom, []byte("content")...)
	out := StripBOM(in)
	if bytes.HasPrefix(out, bom) {
		t.Error("BOM should be stripped")
	}
	if !bytes.Equal(out, []byte("content")) {
		t.Errorf("unexpected result: %q", out)
	}
}

func TestStripBOM_WithoutBOM_Unchanged(t *testing.T) {
	in := []byte("no bom here")
	out := StripBOM(in)
	if !bytes.Equal(out, in) {
		t.Errorf("content without BOM should be unchanged, got: %q", out)
	}
}

func TestNormalize_AppliesAll(t *testing.T) {
	bom := []byte{0xef, 0xbb, 0xbf}
	in := append(bom, []byte("<!-- Build ID: abc -->\r\nhello\r\n")...)
	out := Normalize(in)
	if bytes.HasPrefix(out, bom) {
		t.Error("BOM should be stripped")
	}
	if bytes.Contains(out, []byte("\r\n")) {
		t.Error("CRLF should be normalized")
	}
	if bytes.Contains(out, []byte("Build ID")) {
		t.Error("Build ID should be stripped")
	}
	if !bytes.Contains(out, []byte("hello")) {
		t.Error("content should be preserved")
	}
}

func TestNormalize_Idempotent(t *testing.T) {
	in := []byte("clean content\n")
	first := Normalize(in)
	second := Normalize(first)
	if !bytes.Equal(first, second) {
		t.Error("Normalize should be idempotent")
	}
}
