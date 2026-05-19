package subprocenv_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/utils/subprocenv"
)

func TestMapToSlice_SingleEntry(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	out := subprocenv.MapToSlice(env)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0] != "FOO=bar" {
		t.Errorf("expected 'FOO=bar', got %q", out[0])
	}
}

func TestMapToSlice_MultipleEntries(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := subprocenv.MapToSlice(env)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	for _, entry := range out {
		if !strings.Contains(entry, "=") {
			t.Errorf("entry %q missing '='", entry)
		}
	}
}

func TestMapToSlice_EmptyValueFormatted(t *testing.T) {
	env := map[string]string{"EMPTY": ""}
	out := subprocenv.MapToSlice(env)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0] != "EMPTY=" {
		t.Errorf("expected 'EMPTY=', got %q", out[0])
	}
}

func TestExternalProcessEnv_PreservesRegularKey(t *testing.T) {
	base := map[string]string{"MY_KEY": "my_val"}
	out := subprocenv.ExternalProcessEnv(base)
	if out["MY_KEY"] != "my_val" {
		t.Errorf("expected MY_KEY=my_val, got %q", out["MY_KEY"])
	}
}

func TestExternalProcessEnv_ReturnsCopy_NotSameMap(t *testing.T) {
	base := map[string]string{"K": "v"}
	out := subprocenv.ExternalProcessEnv(base)
	out["K"] = "modified"
	if base["K"] != "v" {
		t.Error("modifying output should not affect input base map")
	}
}

func TestExternalProcessEnv_EmptyBaseReturnsEmptyMap(t *testing.T) {
	base := map[string]string{}
	out := subprocenv.ExternalProcessEnv(base)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestMapToSlice_EmptyMap(t *testing.T) {
	out := subprocenv.MapToSlice(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %v", out)
	}
}

func TestMapToSlice_ValueContainsEquals(t *testing.T) {
	env := map[string]string{"PATH": "/a:/b=/c"}
	out := subprocenv.MapToSlice(env)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0] != "PATH=/a:/b=/c" {
		t.Errorf("value containing '=' should be preserved, got %q", out[0])
	}
}

func TestExternalProcessEnv_MultipleKeys(t *testing.T) {
	base := map[string]string{"A": "alpha", "B": "beta", "C": "gamma"}
	out := subprocenv.ExternalProcessEnv(base)
	if len(out) < 3 {
		t.Errorf("expected at least 3 keys, got %d", len(out))
	}
}

func TestMapToSlice_KeyContainsUnderscore(t *testing.T) {
	env := map[string]string{"MY_VAR": "hello"}
	out := subprocenv.MapToSlice(env)
	found := false
	for _, e := range out {
		if e == "MY_VAR=hello" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'MY_VAR=hello' in output, got %v", out)
	}
}
