package subprocenv_test

import (
	"testing"

	"github.com/githubnext/apm/internal/utils/subprocenv"
)

func TestExternalProcessEnvNoFreeze(t *testing.T) {
	// Without _MEIPASS set, ExternalProcessEnv returns a clean copy of the base.
	base := map[string]string{
		"LD_LIBRARY_PATH":      "/bundled/lib",
		"LD_LIBRARY_PATH_ORIG": "/usr/lib",
		"HOME":                 "/home/user",
	}
	env := subprocenv.ExternalProcessEnv(base)
	// When not frozen the map is returned as-is (no restoration).
	if env["LD_LIBRARY_PATH"] != "/bundled/lib" {
		t.Errorf("expected /bundled/lib, got %s", env["LD_LIBRARY_PATH"])
	}
}

func TestMapToSlice(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	slice := subprocenv.MapToSlice(env)
	if len(slice) != 2 {
		t.Errorf("expected 2 entries, got %d", len(slice))
	}
	seen := map[string]bool{}
	for _, s := range slice {
		seen[s] = true
	}
	if !seen["FOO=bar"] || !seen["BAZ=qux"] {
		t.Errorf("missing expected entries: %v", slice)
	}
}

func TestExternalProcessEnvNilBase(t *testing.T) {
	// With nil base we get the real process env -- just verify no panic.
	env := subprocenv.ExternalProcessEnv(nil)
	if env == nil {
		t.Error("expected non-nil env")
	}
}
