package buildid_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/buildid"
	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestStabilizeBuildID_EmptyString(t *testing.T) {
	got := buildid.StabilizeBuildID("")
	if got != "" {
		t.Errorf("expected empty string unchanged, got %q", got)
	}
}

func TestStabilizeBuildID_OnlyNewline(t *testing.T) {
	got := buildid.StabilizeBuildID("\n")
	if got != "\n" {
		t.Errorf("expected single newline unchanged, got %q", got)
	}
}

func TestStabilizeBuildID_PlaceholderOnlyLine(t *testing.T) {
	content := compilationconst.BuildIDPlaceholder
	got := buildid.StabilizeBuildID(content)
	if strings.Contains(got, compilationconst.BuildIDPlaceholder) {
		t.Error("placeholder should be replaced")
	}
	if !strings.Contains(got, "<!-- Build ID:") {
		t.Errorf("expected Build ID comment, got %q", got)
	}
}

func TestStabilizeBuildID_MultipleCallsSameHash(t *testing.T) {
	content := "alpha\n" + compilationconst.BuildIDPlaceholder + "\nbeta\ngamma\n"
	r1 := buildid.StabilizeBuildID(content)
	r2 := buildid.StabilizeBuildID(content)
	if r1 != r2 {
		t.Errorf("results differ: %q vs %q", r1, r2)
	}
}

func TestStabilizeBuildID_DifferentContentsGiveDifferentHashes(t *testing.T) {
	c1 := "aaa\n" + compilationconst.BuildIDPlaceholder + "\nbbb\n"
	c2 := "xxx\n" + compilationconst.BuildIDPlaceholder + "\nyyy\n"
	r1 := buildid.StabilizeBuildID(c1)
	r2 := buildid.StabilizeBuildID(c2)
	if r1 == r2 {
		t.Error("different contents should produce different hashes")
	}
}

func TestStabilizeBuildID_NoTrailingNewline(t *testing.T) {
	content := "line1\n" + compilationconst.BuildIDPlaceholder + "\nline2"
	got := buildid.StabilizeBuildID(content)
	if strings.HasSuffix(got, "\n") {
		t.Errorf("should not append trailing newline when input has none: %q", got)
	}
}

func TestStabilizeBuildID_LargeContent(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 500; i++ {
		sb.WriteString("# line content here\n")
	}
	sb.WriteString(compilationconst.BuildIDPlaceholder)
	sb.WriteString("\n")
	for i := 0; i < 500; i++ {
		sb.WriteString("# footer line\n")
	}
	got := buildid.StabilizeBuildID(sb.String())
	if strings.Contains(got, compilationconst.BuildIDPlaceholder) {
		t.Error("placeholder not replaced in large content")
	}
}

func TestStabilizeBuildID_PlaceholderInFirstLine(t *testing.T) {
content := compilationconst.BuildIDPlaceholder + "\nsome content\n"
got := buildid.StabilizeBuildID(content)
if strings.Contains(got, compilationconst.BuildIDPlaceholder) {
t.Error("placeholder in first line should be replaced")
}
if !strings.Contains(got, "<!-- Build ID:") {
t.Errorf("expected Build ID comment, got %q", got)
}
}

func TestStabilizeBuildID_NoBuildIDPlaceholder(t *testing.T) {
content := "line1\nline2\nline3\n"
got := buildid.StabilizeBuildID(content)
if got != content {
t.Errorf("content without placeholder should be unchanged, got %q", got)
}
}

func TestStabilizeBuildID_HashLength(t *testing.T) {
content := compilationconst.BuildIDPlaceholder
got := buildid.StabilizeBuildID(content)
// Extract hash from "<!-- Build ID: <hash> -->"
if !strings.Contains(got, "<!-- Build ID:") {
t.Fatal("expected Build ID comment")
}
inner := strings.TrimPrefix(got, "<!-- Build ID: ")
inner = strings.TrimSuffix(inner, " -->")
if len(inner) != 12 {
t.Errorf("expected 12-char hash, got %d chars: %q", len(inner), inner)
}
}

func TestStabilizeBuildID_TrailingNewlinePreserved(t *testing.T) {
content := "a\n" + compilationconst.BuildIDPlaceholder + "\nb\n"
got := buildid.StabilizeBuildID(content)
if !strings.HasSuffix(got, "\n") {
t.Errorf("trailing newline should be preserved, got %q", got)
}
}

func TestStabilizeBuildID_OnlyPlaceholderNoNewline(t *testing.T) {
content := compilationconst.BuildIDPlaceholder
got := buildid.StabilizeBuildID(content)
// Must not add newline since input has none
if strings.HasSuffix(got, "\n") {
t.Error("should not add trailing newline when input has none")
}
}
