package buildid_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/buildid"
	"github.com/githubnext/apm/internal/compilation/compilationconst"
)

func TestStabilizeBuildID_NoPlaceholder(t *testing.T) {
	content := "# Some content\nno placeholder here\n"
	got := buildid.StabilizeBuildID(content)
	if got != content {
		t.Errorf("expected unchanged content, got %q", got)
	}
}

func TestStabilizeBuildID_ReplacesPlaceholder(t *testing.T) {
	content := "line one\n" + compilationconst.BuildIDPlaceholder + "\nline three\n"
	got := buildid.StabilizeBuildID(content)
	if strings.Contains(got, compilationconst.BuildIDPlaceholder) {
		t.Error("placeholder was not replaced")
	}
	if !strings.Contains(got, "<!-- Build ID:") {
		t.Errorf("expected Build ID comment, got: %q", got)
	}
}

func TestStabilizeBuildID_Idempotent(t *testing.T) {
	content := "header\n" + compilationconst.BuildIDPlaceholder + "\nfooter\n"
	first := buildid.StabilizeBuildID(content)
	second := buildid.StabilizeBuildID(first)
	if first != second {
		t.Errorf("StabilizeBuildID is not idempotent: first=%q second=%q", first, second)
	}
}

func TestStabilizeBuildID_DeterministicHash(t *testing.T) {
	content := "a\n" + compilationconst.BuildIDPlaceholder + "\nb\n"
	got1 := buildid.StabilizeBuildID(content)
	got2 := buildid.StabilizeBuildID(content)
	if got1 != got2 {
		t.Errorf("non-deterministic: %q vs %q", got1, got2)
	}
}

func TestStabilizeBuildID_PreservesTrailingNewline(t *testing.T) {
	withNL := "x\n" + compilationconst.BuildIDPlaceholder + "\ny\n"
	withoutNL := "x\n" + compilationconst.BuildIDPlaceholder + "\ny"
	gotWith := buildid.StabilizeBuildID(withNL)
	gotWithout := buildid.StabilizeBuildID(withoutNL)
	if !strings.HasSuffix(gotWith, "\n") {
		t.Errorf("trailing newline not preserved: %q", gotWith)
	}
	if strings.HasSuffix(gotWithout, "\n") {
		t.Errorf("unexpected trailing newline added: %q", gotWithout)
	}
}

func TestStabilizeBuildID_12CharHash(t *testing.T) {
	content := "data\n" + compilationconst.BuildIDPlaceholder + "\nmore\n"
	got := buildid.StabilizeBuildID(content)
	// Extract the Build ID value from "<!-- Build ID: XXXXXXXXXXXX -->"
	const prefix = "<!-- Build ID: "
	const suffix = " -->"
	idx := strings.Index(got, prefix)
	if idx < 0 {
		t.Fatalf("no Build ID comment in %q", got)
	}
	inner := got[idx+len(prefix):]
	end := strings.Index(inner, suffix)
	if end < 0 {
		t.Fatalf("malformed Build ID comment in %q", got)
	}
	hash := inner[:end]
	if len(hash) != 12 {
		t.Errorf("expected 12-char hash, got %d chars: %q", len(hash), hash)
	}
}
