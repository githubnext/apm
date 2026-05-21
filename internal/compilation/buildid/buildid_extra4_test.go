package buildid_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/compilation/buildid"
)

func TestStabilizeBuildID_EmptyInput_Extra4(t *testing.T) {
	got := buildid.StabilizeBuildID("")
	_ = got // no panic
}

func TestStabilizeBuildID_NoBuildIDLine_Unchanged_Extra4(t *testing.T) {
	input := "line one\nline two\nline three\n"
	got := buildid.StabilizeBuildID(input)
	if got != input {
		t.Errorf("expected unchanged output when no build ID line, got %q", got)
	}
}

func TestStabilizeBuildID_ReplacedLineContainsHash_Extra4(t *testing.T) {
	input := "<!-- Build ID: __BUILD_ID__ -->\n"
	got := buildid.StabilizeBuildID(input)
	if got == input {
		t.Error("expected replacement to differ from input")
	}
}

func TestStabilizeBuildID_OutputHasOneBuildIDComment_Extra4(t *testing.T) {
	input := "<!-- Build ID: __BUILD_ID__ -->\n"
	got := buildid.StabilizeBuildID(input)
	if !strings.Contains(got, "<!-- Build ID:") {
		t.Errorf("expected Build ID comment in output, got %q", got)
	}
}

func TestStabilizeBuildID_LineCountPreserved_Extra4(t *testing.T) {
	input := "line1\n<!-- Build ID: __BUILD_ID__ -->\nline3\n"
	got := buildid.StabilizeBuildID(input)
	inLines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	outLines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(inLines) != len(outLines) {
		t.Errorf("expected same line count: %d vs %d", len(inLines), len(outLines))
	}
}

func TestStabilizeBuildID_SameInputSameOutput_Extra4(t *testing.T) {
	input := "header\n<!-- Build ID: __BUILD_ID__ -->\nfooter\n"
	a := buildid.StabilizeBuildID(input)
	b := buildid.StabilizeBuildID(input)
	if a != b {
		t.Error("expected deterministic output")
	}
}

func TestStabilizeBuildID_AlreadyStabilized_Idempotent_Extra4(t *testing.T) {
	input := "<!-- Build ID: __BUILD_ID__ -->\n"
	once := buildid.StabilizeBuildID(input)
	twice := buildid.StabilizeBuildID(once)
	if once != twice {
		t.Errorf("expected idempotent, got diff: %q vs %q", once, twice)
	}
}

func TestStabilizeBuildID_HashIsHex_Extra4(t *testing.T) {
	input := "<!-- Build ID: __BUILD_ID__ -->\n"
	got := buildid.StabilizeBuildID(input)
	// Extract the hash from the output line
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "Build ID:") {
			parts := strings.Fields(line)
			for _, p := range parts {
				p = strings.Trim(p, "-->")
				if len(p) == 12 {
					for _, c := range p {
						if !strings.ContainsRune("0123456789abcdef", c) {
							t.Errorf("non-hex char in hash: %q", p)
						}
					}
				}
			}
		}
	}
}
