package buildid_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/buildid"
	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestStabilizeBuildID_PlaceholderLastLine(t *testing.T) {
	content := "header\nbody\n" + compilationconst.BuildIDPlaceholder + "\n"
	got := buildid.StabilizeBuildID(content)
	if strings.Contains(got, compilationconst.BuildIDPlaceholder) {
		t.Error("placeholder on last line should be replaced")
	}
	if !strings.Contains(got, "<!-- Build ID:") {
		t.Error("expected Build ID comment")
	}
}

func TestStabilizeBuildID_PlaceholderMiddleLine(t *testing.T) {
	content := "top\n" + compilationconst.BuildIDPlaceholder + "\nbottom\n"
	got := buildid.StabilizeBuildID(content)
	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if lines[1] != "<!-- Build ID: "+lines[1][len("<!-- Build ID: "):len(lines[1])-4]+" -->" {
		// just verify it starts with the prefix
		if !strings.HasPrefix(lines[1], "<!-- Build ID:") {
			t.Errorf("middle line should be replaced: %q", lines[1])
		}
	}
}

func TestStabilizeBuildID_IdempotentStable(t *testing.T) {
	content := "a\n" + compilationconst.BuildIDPlaceholder + "\nb\n"
	once := buildid.StabilizeBuildID(content)
	twice := buildid.StabilizeBuildID(once)
	if once != twice {
		t.Errorf("second call changed output: %q vs %q", once, twice)
	}
}

func TestStabilizeBuildID_OnlyFirstPlaceholderReplaced(t *testing.T) {
	content := compilationconst.BuildIDPlaceholder + "\n" + compilationconst.BuildIDPlaceholder + "\n"
	got := buildid.StabilizeBuildID(content)
	// First placeholder replaced; second one stays because the loop stops at first
	count := strings.Count(got, compilationconst.BuildIDPlaceholder)
	if count != 1 {
		t.Errorf("expected 1 remaining placeholder, got %d in: %q", count, got)
	}
}

func TestStabilizeBuildID_CommentFormat(t *testing.T) {
	content := compilationconst.BuildIDPlaceholder + "\n"
	got := buildid.StabilizeBuildID(content)
	got = strings.TrimSuffix(got, "\n")
	if !strings.HasPrefix(got, "<!-- Build ID: ") {
		t.Errorf("output should start with '<!-- Build ID: ', got %q", got)
	}
	if !strings.HasSuffix(got, " -->") {
		t.Errorf("output should end with ' -->', got %q", got)
	}
}

func TestStabilizeBuildID_ContentAroundPlaceholderPreserved(t *testing.T) {
	content := "before\n" + compilationconst.BuildIDPlaceholder + "\nafter\n"
	got := buildid.StabilizeBuildID(content)
	if !strings.Contains(got, "before\n") {
		t.Error("before content should be preserved")
	}
	if !strings.Contains(got, "\nafter\n") {
		t.Error("after content should be preserved")
	}
}

func TestStabilizeBuildID_TwoLinesNoBuildID(t *testing.T) {
	content := "line1\nline2"
	got := buildid.StabilizeBuildID(content)
	if got != content {
		t.Errorf("no placeholder: content should be unchanged, got %q", got)
	}
}

func TestStabilizeBuildID_SingleLineNoBuildID(t *testing.T) {
	content := "just a line"
	got := buildid.StabilizeBuildID(content)
	if got != content {
		t.Errorf("expected unchanged, got %q", got)
	}
}
