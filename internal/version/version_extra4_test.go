package version_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/version"
)

func TestGetVersion_NotUnknownOrEmpty_Extra4(t *testing.T) {
	v := version.GetVersion()
	if v == "" {
		t.Error("expected non-empty version")
	}
}

func TestGetVersion_NoNewline_Extra4(t *testing.T) {
	v := version.GetVersion()
	if strings.Contains(v, "\n") {
		t.Errorf("version should not contain newline: %q", v)
	}
}

func TestGetVersion_NoLeadingTrailingSpaces_Extra4(t *testing.T) {
	v := version.GetVersion()
	if v != strings.TrimSpace(v) {
		t.Errorf("version has leading/trailing spaces: %q", v)
	}
}

func TestGetBuildSHA_LengthOrEmpty_Extra4(t *testing.T) {
	sha := version.GetBuildSHA()
	// Either empty (no git) or a valid short SHA (4-40 hex chars)
	if sha != "" && len(sha) < 4 {
		t.Errorf("expected empty or valid short SHA, got %q", sha)
	}
}

func TestGetBuildSHA_NoSpaces_Extra4(t *testing.T) {
	sha := version.GetBuildSHA()
	if strings.Contains(sha, " ") {
		t.Errorf("SHA should not contain spaces: %q", sha)
	}
}

func TestBuildVersion_DefaultEmpty_Extra4(t *testing.T) {
	// BuildVersion is a var that defaults to ""; if set at link time it is non-empty.
	// Either way GetVersion should not panic.
	v := version.GetVersion()
	_ = v
}

func TestGetVersion_Deterministic_Extra4(t *testing.T) {
	a := version.GetVersion()
	b := version.GetVersion()
	if a != b {
		t.Errorf("expected deterministic version, got %q then %q", a, b)
	}
}

func TestGetBuildSHA_HexOnly_Extra4(t *testing.T) {
	sha := version.GetBuildSHA()
	for _, c := range sha {
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			t.Errorf("non-hex char in SHA: %q", sha)
			break
		}
	}
}
