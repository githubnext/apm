package subprocenv_test

import (
	"strings"
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

func TestExternalProcessEnvReturnsCopy(t *testing.T) {
	base := map[string]string{"KEY": "value"}
	env := subprocenv.ExternalProcessEnv(base)
	// Mutating the returned map should not affect the original.
	env["KEY"] = "modified"
	if base["KEY"] != "value" {
		t.Error("ExternalProcessEnv should return an independent copy")
	}
}

func TestExternalProcessEnvPreservesNonLibraryVars(t *testing.T) {
	base := map[string]string{
		"HOME":   "/home/user",
		"PATH":   "/usr/bin:/bin",
		"EDITOR": "vim",
	}
	env := subprocenv.ExternalProcessEnv(base)
	if env["HOME"] != "/home/user" {
		t.Errorf("HOME should be preserved, got %q", env["HOME"])
	}
	if env["PATH"] != "/usr/bin:/bin" {
		t.Errorf("PATH should be preserved, got %q", env["PATH"])
	}
	if env["EDITOR"] != "vim" {
		t.Errorf("EDITOR should be preserved, got %q", env["EDITOR"])
	}
}

func TestMapToSliceEmpty(t *testing.T) {
	slice := subprocenv.MapToSlice(map[string]string{})
	if len(slice) != 0 {
		t.Errorf("expected empty slice, got %v", slice)
	}
}

func TestMapToSliceFormatting(t *testing.T) {
	env := map[string]string{"MYKEY": "myval"}
	slice := subprocenv.MapToSlice(env)
	if len(slice) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(slice))
	}
	if slice[0] != "MYKEY=myval" {
		t.Errorf("expected MYKEY=myval, got %q", slice[0])
	}
}

func TestMapToSliceContainsEquals(t *testing.T) {
	env := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}
	for _, kv := range subprocenv.MapToSlice(env) {
		if !strings.Contains(kv, "=") {
			t.Errorf("MapToSlice entry %q missing '='", kv)
		}
	}
}

func TestMapToSliceEmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY": ""}
	slice := subprocenv.MapToSlice(env)
	if len(slice) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(slice))
	}
	if slice[0] != "EMPTY=" {
		t.Errorf("expected EMPTY=, got %q", slice[0])
	}
}

func TestExternalProcessEnvLibraryPathsNotModifiedWhenNotFrozen(t *testing.T) {
	// When not frozen, all keys are preserved including the ORIG variants.
	base := map[string]string{
		"LD_LIBRARY_PATH":            "/bundled",
		"LD_LIBRARY_PATH_ORIG":       "/system",
		"DYLD_LIBRARY_PATH":          "/bundled-mac",
		"DYLD_LIBRARY_PATH_ORIG":     "/system-mac",
		"DYLD_FRAMEWORK_PATH":        "/bundled-fw",
		"DYLD_FRAMEWORK_PATH_ORIG":   "/system-fw",
	}
	env := subprocenv.ExternalProcessEnv(base)
	// Without freeze, all keys are present unchanged.
	if env["LD_LIBRARY_PATH"] != "/bundled" {
		t.Errorf("LD_LIBRARY_PATH: got %q", env["LD_LIBRARY_PATH"])
	}
	if env["LD_LIBRARY_PATH_ORIG"] != "/system" {
		t.Errorf("LD_LIBRARY_PATH_ORIG: got %q", env["LD_LIBRARY_PATH_ORIG"])
	}
}
