package coworkpaths

import (
	"strings"
	"testing"
)

func TestIsCoworkPath_Empty(t *testing.T) {
	if IsCoworkPath("") {
		t.Error("empty string should not be a cowork path")
	}
}

func TestIsCoworkPath_Scheme(t *testing.T) {
	if !IsCoworkPath(CoworkURIScheme + "something") {
		t.Error("CoworkURIScheme prefix should be a cowork path")
	}
}

func TestCoworkURIScheme_StartsWithCowork(t *testing.T) {
	if !strings.HasPrefix(CoworkURIScheme, "cowork://") {
		t.Errorf("CoworkURIScheme should start with 'cowork://', got %q", CoworkURIScheme)
	}
}

func TestCoworkLockfilePrefix_StartsWithScheme(t *testing.T) {
	if !strings.HasPrefix(CoworkLockfilePrefix, CoworkURIScheme) {
		t.Errorf("CoworkLockfilePrefix should start with CoworkURIScheme: %q", CoworkLockfilePrefix)
	}
}

func TestCoworkLockfilePrefix_ContainsSkills(t *testing.T) {
	if !strings.Contains(CoworkLockfilePrefix, "skills") {
		t.Errorf("CoworkLockfilePrefix should contain 'skills': %q", CoworkLockfilePrefix)
	}
}

func TestToLockfilePath_Basic(t *testing.T) {
	root := t.TempDir()
	abs := root + "/some/skill"
	got, err := ToLockfilePath(abs, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(got, CoworkURIScheme) {
		t.Errorf("result should start with CoworkURIScheme: %q", got)
	}
}

func TestFromLockfilePath_RoundTrip(t *testing.T) {
	root := t.TempDir()
	abs := root + "/myfolder/skill"
	encoded, err := ToLockfilePath(abs, root)
	if err != nil {
		t.Fatalf("ToLockfilePath: %v", err)
	}
	decoded, err := FromLockfilePath(encoded, root)
	if err != nil {
		t.Fatalf("FromLockfilePath: %v", err)
	}
	// Normalize separators for comparison
	got := strings.ReplaceAll(decoded, "\\", "/")
	want := strings.ReplaceAll(abs, "\\", "/")
	if got != want {
		t.Errorf("round-trip: got %q want %q", got, want)
	}
}

func TestFromLockfilePath_NonCoworkError(t *testing.T) {
	_, err := FromLockfilePath("https://example.com/foo", t.TempDir())
	if err == nil {
		t.Error("expected error for non-cowork path")
	}
}

func TestCoworkResolutionError_ErrorString(t *testing.T) {
	e := &CoworkResolutionError{Msg: "test error msg"}
	if e.Error() != "test error msg" {
		t.Errorf("Error() = %q, want 'test error msg'", e.Error())
	}
}

func TestIsCoworkPath_LocalPath(t *testing.T) {
	if IsCoworkPath("/local/path/to/something") {
		t.Error("local absolute path should not be cowork path")
	}
}

func TestIsCoworkPath_HTTPSPath(t *testing.T) {
	if IsCoworkPath("https://github.com/owner/repo") {
		t.Error("https URL should not be cowork path")
	}
}
