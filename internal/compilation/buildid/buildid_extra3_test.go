package buildid

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestStabilizeBuildID_MultilinePre(t *testing.T) {
	content := "line1\nline2\n" + compilationconst.BuildIDPlaceholder + "\nline3\n"
	result := StabilizeBuildID(content)
	if strings.Contains(result, compilationconst.BuildIDPlaceholder) {
		t.Fatal("placeholder should have been replaced")
	}
	if !strings.Contains(result, "line1") || !strings.Contains(result, "line3") {
		t.Fatal("surrounding lines should be preserved")
	}
}

func TestStabilizeBuildID_PlaceholderAtEnd(t *testing.T) {
	content := "first line\n" + compilationconst.BuildIDPlaceholder
	result := StabilizeBuildID(content)
	if strings.Contains(result, compilationconst.BuildIDPlaceholder) {
		t.Fatal("placeholder should be replaced even at end of content")
	}
}

func TestStabilizeBuildID_MultiplePlaceholders_OnlyFirst(t *testing.T) {
	p := compilationconst.BuildIDPlaceholder
	content := p + "\n" + p + "\n"
	result := StabilizeBuildID(content)
	count := strings.Count(result, p)
	if count != 1 {
		t.Fatalf("expected 1 remaining placeholder, got %d", count)
	}
}

func TestStabilizeBuildID_OutputContainsBuildIDComment(t *testing.T) {
	content := compilationconst.BuildIDPlaceholder + "\n"
	result := StabilizeBuildID(content)
	if !strings.Contains(result, "<!-- Build ID:") {
		t.Fatalf("expected '<!-- Build ID:' in output, got %q", result)
	}
}

func TestStabilizeBuildID_HashIs12HexChars(t *testing.T) {
	content := compilationconst.BuildIDPlaceholder + "\n"
	result := StabilizeBuildID(content)
	start := strings.Index(result, "<!-- Build ID: ")
	if start < 0 {
		t.Fatal("no build id comment found")
	}
	id := result[start+len("<!-- Build ID: "):]
	id = strings.TrimSuffix(id, " -->")
	id = strings.TrimSuffix(id, " -->\n")
	end := strings.Index(id, " ")
	if end > 0 {
		id = id[:end]
	}
	id = strings.TrimRight(id, " \n>-")
	if len(id) < 12 {
		t.Fatalf("expected hash of at least 12 chars, got %q (len=%d)", id, len(id))
	}
}

func TestStabilizeBuildID_ContentBeforePlaceholderPreservedVariant(t *testing.T) {
	content := "header\n" + compilationconst.BuildIDPlaceholder + "\n"
	result := StabilizeBuildID(content)
	if !strings.Contains(result, "header") {
		t.Fatal("header line should be preserved")
	}
}

func TestStabilizeBuildID_NoPlaceholderReturnsSame(t *testing.T) {
	content := "no placeholder here\nanother line\n"
	result := StabilizeBuildID(content)
	if result != content {
		t.Fatalf("expected unchanged content, got %q", result)
	}
}

func TestStabilizeBuildID_SameSeedSameHash(t *testing.T) {
	content := "abc\n" + compilationconst.BuildIDPlaceholder + "\ndef\n"
	r1 := StabilizeBuildID(content)
	r2 := StabilizeBuildID(content)
	if r1 != r2 {
		t.Fatal("same input should produce same hash")
	}
}

func TestStabilizeBuildID_DiffSeedDiffHash(t *testing.T) {
	p := compilationconst.BuildIDPlaceholder
	r1 := StabilizeBuildID("aaa\n" + p + "\n")
	r2 := StabilizeBuildID("bbb\n" + p + "\n")
	if r1 == r2 {
		t.Fatal("different seeds should produce different hashes")
	}
}
