package normalization

import "testing"

func TestStripBuildID_NoBuildID_Extra4(t *testing.T) {
	input := []byte("# Hello\nsome content\n")
	out := StripBuildID(input)
	if string(out) != string(input) {
		t.Errorf("StripBuildID changed content without build ID: %q", out)
	}
}

func TestStripBuildID_WithBuildID_Extra4(t *testing.T) {
	input := []byte("<!-- Build ID: abc123 -->\n# Hello\n")
	out := StripBuildID(input)
	if len(out) >= len(input) {
		t.Errorf("StripBuildID did not shorten content: got %d bytes, want < %d", len(out), len(input))
	}
}

func TestNormalizeLineEndings_Mixed_Extra4(t *testing.T) {
	input := []byte("line1\r\nline2\r\nline3\n")
	out := NormalizeLineEndings(input)
	if string(out) != "line1\nline2\nline3\n" {
		t.Errorf("NormalizeLineEndings = %q", out)
	}
}

func TestNormalizeLineEndings_NoChange_Extra4(t *testing.T) {
	input := []byte("line1\nline2\n")
	out := NormalizeLineEndings(input)
	if string(out) != string(input) {
		t.Errorf("NormalizeLineEndings changed clean content")
	}
}

func TestStripBOM_WithBOM_Extra4(t *testing.T) {
	bom := []byte{0xef, 0xbb, 0xbf}
	input := append(bom, []byte("content")...)
	out := StripBOM(input)
	if string(out) != "content" {
		t.Errorf("StripBOM result = %q, want content", out)
	}
}

func TestStripBOM_WithoutBOM_Extra4(t *testing.T) {
	input := []byte("content")
	out := StripBOM(input)
	if string(out) != "content" {
		t.Errorf("StripBOM should not change content without BOM")
	}
}

func TestNormalize_Empty_Extra4(t *testing.T) {
	out := Normalize([]byte{})
	if len(out) != 0 {
		t.Errorf("Normalize(empty) = %q", out)
	}
}

func TestNormalize_AllTransforms_Extra4(t *testing.T) {
	bom := []byte{0xef, 0xbb, 0xbf}
	input := append(bom, []byte("<!-- Build ID: deadbeef -->\r\ntext\r\n")...)
	out := Normalize(input)
	if len(out) == 0 {
		t.Error("Normalize should return non-empty result")
	}
	for _, b := range out {
		if b == '\r' {
			t.Error("Normalize should remove all CR characters")
		}
	}
}

func TestStripBuildID_CaseInsensitive_Extra4(t *testing.T) {
	input := []byte("<!-- build id: abc123 -->\n# Title\n")
	out := StripBuildID(input)
	if len(out) >= len(input) {
		t.Errorf("case-insensitive strip failed: still %d bytes", len(out))
	}
}
