package coworkpaths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsCoworkPath_ValidVariants(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"cowork://skills/foo", true},
		{"cowork://skills/foo/bar", true},
		{"/local/path", false},
		{"./relative", false},
		{"https://github.com/repo", false},
	}
	for _, tc := range cases {
		got := IsCoworkPath(tc.path)
		if got != tc.want {
			t.Errorf("IsCoworkPath(%q): got %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestToLockfilePath_DeepPath(t *testing.T) {
	root := t.TempDir()
	deep := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	lp, err := ToLockfilePath(deep, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lp == "" {
		t.Error("expected non-empty lockfile path")
	}
	if !IsCoworkPath(lp) {
		t.Errorf("expected cowork path, got %q", lp)
	}
}

func TestFromLockfilePath_NestedPath(t *testing.T) {
	root := t.TempDir()
	lp := "cowork://skills/nested/plugin"
	got, err := FromLockfilePath(lp, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(root, "nested", "plugin")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRoundTrip_SingleSegment(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "mysub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	lp, err := ToLockfilePath(sub, root)
	if err != nil {
		t.Fatal(err)
	}
	back, err := FromLockfilePath(lp, root)
	if err != nil {
		t.Fatal(err)
	}
	if back != sub {
		t.Errorf("round-trip: want %q, got %q", sub, back)
	}
}

func TestResolveCoworkSkillsDir_NoEnv(t *testing.T) {
	t.Setenv("APM_COPILOT_COWORK_SKILLS_DIR", "")
	// Without env, it may succeed or fail depending on system config
	_, _ = ResolveCoworkSkillsDir()
}

func TestCoworkResolutionError_Error(t *testing.T) {
	err := &CoworkResolutionError{Msg: "something went wrong"}
	if err.Error() != "something went wrong" {
		t.Errorf("unexpected error: %q", err.Error())
	}
}

func TestFromLockfilePath_MissingScheme(t *testing.T) {
	_, err := FromLockfilePath("skills/foo", "/root")
	if err == nil {
		t.Fatal("expected error for path missing cowork:// scheme")
	}
}

func TestToLockfilePath_RootItself(t *testing.T) {
	root := t.TempDir()
	// The root itself -- should produce a path at the root level
	lp, err := ToLockfilePath(root, root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = lp
}
