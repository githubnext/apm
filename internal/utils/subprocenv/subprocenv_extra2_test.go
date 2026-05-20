package subprocenv_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/subprocenv"
)

func TestIsWindows_ReturnsBool(t *testing.T) {
	// Just ensure no panic and returns a bool.
	_ = subprocenv.IsWindows()
}

func TestMapToSlice_NilMap(t *testing.T) {
	// nil map should produce empty slice without panic.
	out := subprocenv.MapToSlice(nil)
	if out == nil {
		t.Error("expected non-nil slice for nil map")
	}
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(out))
	}
}

func TestExternalProcessEnv_OrigRestored_WhenFrozen(t *testing.T) {
	// When _MEIPASS is set the ORIG sibling should be restored.
	t.Setenv("_MEIPASS", "/some/tmpdir")
	base := map[string]string{
		"LD_LIBRARY_PATH":      "/bundled",
		"LD_LIBRARY_PATH_ORIG": "/usr/lib",
		"HOME":                 "/home/user",
	}
	out := subprocenv.ExternalProcessEnv(base)
	if got := out["LD_LIBRARY_PATH"]; got != "/usr/lib" {
		t.Errorf("LD_LIBRARY_PATH should be restored to ORIG, got %q", got)
	}
	if _, present := out["LD_LIBRARY_PATH_ORIG"]; present {
		t.Error("LD_LIBRARY_PATH_ORIG should be stripped from the output map")
	}
}

func TestExternalProcessEnv_VarRemoved_WhenNoOrig_AndFrozen(t *testing.T) {
	// When frozen but no ORIG exists the managed var should be removed.
	t.Setenv("_MEIPASS", "/some/tmpdir")
	base := map[string]string{
		"LD_LIBRARY_PATH": "/bundled",
		"HOME":            "/home/user",
	}
	out := subprocenv.ExternalProcessEnv(base)
	if _, present := out["LD_LIBRARY_PATH"]; present {
		t.Error("LD_LIBRARY_PATH should be removed when no ORIG and process is frozen")
	}
	if out["HOME"] != "/home/user" {
		t.Error("unrelated key HOME should be preserved")
	}
}

func TestExternalProcessEnv_NoMutation_WhenNotFrozen(t *testing.T) {
	// Without _MEIPASS the map is returned unchanged.
	base := map[string]string{
		"LD_LIBRARY_PATH":      "/bundled",
		"LD_LIBRARY_PATH_ORIG": "/sys",
		"FOO":                  "bar",
	}
	out := subprocenv.ExternalProcessEnv(base)
	if out["LD_LIBRARY_PATH"] != "/bundled" {
		t.Errorf("LD_LIBRARY_PATH should not be modified when not frozen")
	}
	if out["LD_LIBRARY_PATH_ORIG"] != "/sys" {
		t.Errorf("LD_LIBRARY_PATH_ORIG should be present when not frozen")
	}
}
