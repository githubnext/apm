package compilation_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation"
)

// TestParityBuildIDConstants verifies constant values match Python source.
func TestParityBuildIDConstants(t *testing.T) {
	if compilation.BuildIDPlaceholder != "<!-- Build ID: __BUILD_ID__ -->" {
		t.Fatalf("unexpected BuildIDPlaceholder: %s", compilation.BuildIDPlaceholder)
	}
	if compilation.ConstitutionRelativePath != ".specify/memory/constitution.md" {
		t.Fatalf("unexpected path: %s", compilation.ConstitutionRelativePath)
	}
}

// TestParityConstitutionMarkers verifies constitution marker constants.
func TestParityConstitutionMarkers(t *testing.T) {
	if compilation.ConstitutionMarkerBegin != "<!-- SPEC-KIT CONSTITUTION: BEGIN -->" {
		t.Fatalf("unexpected begin marker: %s", compilation.ConstitutionMarkerBegin)
	}
	if compilation.ConstitutionMarkerEnd != "<!-- SPEC-KIT CONSTITUTION: END -->" {
		t.Fatalf("unexpected end marker: %s", compilation.ConstitutionMarkerEnd)
	}
}

// TestParityStabilizeBuildIDNoPlaceholder returns input unchanged when no placeholder.
func TestParityStabilizeBuildIDNoPlaceholder(t *testing.T) {
	content := "# My doc\nSome content.\n"
	got := compilation.StabilizeBuildID(content)
	if got != content {
		t.Fatalf("expected unchanged content, got: %s", got)
	}
}

// TestParityStabilizeBuildIDReplacesPlaceholder verifies placeholder is replaced.
func TestParityStabilizeBuildIDReplacesPlaceholder(t *testing.T) {
	content := "# My doc\n" + compilation.BuildIDPlaceholder + "\nSome content.\n"
	got := compilation.StabilizeBuildID(content)
	if strings.Contains(got, "__BUILD_ID__") {
		t.Fatal("placeholder was not replaced")
	}
	if !strings.Contains(got, "<!-- Build ID:") {
		t.Fatalf("expected Build ID comment in output: %s", got)
	}
}

// TestParityStabilizeBuildIDDeterministic verifies same input produces same hash.
func TestParityStabilizeBuildIDDeterministic(t *testing.T) {
	content := "line1\n" + compilation.BuildIDPlaceholder + "\nline3\n"
	r1 := compilation.StabilizeBuildID(content)
	r2 := compilation.StabilizeBuildID(content)
	if r1 != r2 {
		t.Fatalf("not deterministic: %q vs %q", r1, r2)
	}
}

// TestParityStabilizeBuildIDHashLength verifies hash is 12 chars.
func TestParityStabilizeBuildIDHashLength(t *testing.T) {
	content := "hello\n" + compilation.BuildIDPlaceholder + "\nworld\n"
	got := compilation.StabilizeBuildID(content)
	// Extract the hash value from "<!-- Build ID: <hash> -->"
	const prefix = "<!-- Build ID: "
	const suffix = " -->"
	idx := strings.Index(got, prefix)
	if idx == -1 {
		t.Fatalf("no Build ID comment found in: %s", got)
	}
	start := idx + len(prefix)
	end := strings.Index(got[start:], suffix)
	if end == -1 {
		t.Fatalf("malformed Build ID comment: %s", got)
	}
	hash := got[start : start+end]
	if len(hash) != 12 {
		t.Fatalf("expected 12-char hash, got %d chars: %s", len(hash), hash)
	}
}

// TestParityStabilizeBuildIDIdempotent verifies second call is idempotent.
func TestParityStabilizeBuildIDIdempotent(t *testing.T) {
	content := "abc\n" + compilation.BuildIDPlaceholder + "\ndef\n"
	once := compilation.StabilizeBuildID(content)
	twice := compilation.StabilizeBuildID(once)
	if once != twice {
		t.Fatalf("not idempotent: %q vs %q", once, twice)
	}
}

// TestParityStabilizeBuildIDTrailingNewline verifies trailing newline preservation.
func TestParityStabilizeBuildIDTrailingNewline(t *testing.T) {
	content := "# doc\n" + compilation.BuildIDPlaceholder + "\ncontent\n"
	got := compilation.StabilizeBuildID(content)
	if !strings.HasSuffix(got, "\n") {
		t.Fatal("trailing newline was lost")
	}
}

// TestParityStabilizeBuildIDNoTrailingNewline verifies non-trailing-newline input.
func TestParityStabilizeBuildIDNoTrailingNewline(t *testing.T) {
	content := "# doc\n" + compilation.BuildIDPlaceholder
	got := compilation.StabilizeBuildID(content)
	if strings.HasSuffix(got, "\n") {
		t.Fatal("unexpected trailing newline added")
	}
}

// TestParityStabilizeBuildIDDifferentContentDifferentHash verifies hash changes.
func TestParityStabilizeBuildIDDifferentContentDifferentHash(t *testing.T) {
	c1 := "aaa\n" + compilation.BuildIDPlaceholder + "\nbbb\n"
	c2 := "xxx\n" + compilation.BuildIDPlaceholder + "\nyyy\n"
	r1 := compilation.StabilizeBuildID(c1)
	r2 := compilation.StabilizeBuildID(c2)
	if r1 == r2 {
		t.Fatal("different content should produce different hash")
	}
}

// TestParityStabilizeBuildIDEmptyContent handles edge case of just placeholder.
func TestParityStabilizeBuildIDEmptyContent(t *testing.T) {
	got := compilation.StabilizeBuildID(compilation.BuildIDPlaceholder)
	if strings.Contains(got, "__BUILD_ID__") {
		t.Fatal("placeholder should be replaced even for minimal content")
	}
}
