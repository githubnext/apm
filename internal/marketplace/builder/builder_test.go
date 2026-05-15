package builder

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// isDisplayVersion
// ---------------------------------------------------------------------------

func TestIsDisplayVersionSimple(t *testing.T) {
	cases := []struct {
		v    string
		want bool
	}{
		{"1.2.3", true},
		{"v1.0.0", true},
		{"1.2.3-beta", true},
		{"", false},
		{"^1.0.0", false},
		{"~1.0.0", false},
		{">1.0.0", false},
		{"<1.0.0", false},
		{">=1.0.0", false},
		{"1.x", false},
		{"1.2.*", false},
		{"1 2 3", false},
	}
	for _, c := range cases {
		got := isDisplayVersion(c.v)
		if got != c.want {
			t.Errorf("isDisplayVersion(%q) = %v, want %v", c.v, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// subtractPluginRoot
// ---------------------------------------------------------------------------

func TestSubtractPluginRoot(t *testing.T) {
	cases := []struct {
		src, root, want string
		wantErr         bool
	}{
		{"./plugins/my-plugin/file.json", "./plugins/my-plugin", "./file.json", false},
		{"plugins/my-plugin/sub/file.json", "plugins/my-plugin", "./sub/file.json", false},
		{"other/path/file.json", "plugins/my-plugin", "", true},
		{"./plugins/my-plugin", "./plugins/my-plugin", "", true}, // yields empty
	}
	for _, c := range cases {
		got, err := subtractPluginRoot(c.src, c.root)
		if c.wantErr {
			if err == nil {
				t.Errorf("subtractPluginRoot(%q, %q): expected error, got %q", c.src, c.root, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("subtractPluginRoot(%q, %q): unexpected error: %v", c.src, c.root, err)
			continue
		}
		if got != c.want {
			t.Errorf("subtractPluginRoot(%q, %q) = %q, want %q", c.src, c.root, got, c.want)
		}
	}
}

func TestSubtractPluginRootTraversal(t *testing.T) {
	_, err := subtractPluginRoot("plugins/my-plugin/../../etc/passwd", "plugins/my-plugin")
	if err == nil {
		t.Error("expected error for path traversal")
	}
}

// ---------------------------------------------------------------------------
// Error types
// ---------------------------------------------------------------------------

func TestBuildErrorMessage(t *testing.T) {
	e := &BuildError{Msg: "something went wrong", Package: "pkg-a"}
	if e.Error() != "something went wrong" {
		t.Errorf("unexpected: %q", e.Error())
	}
}

func TestHeadNotAllowedError(t *testing.T) {
	e := &HeadNotAllowedError{Package: "pkg", Ref: "main"}
	msg := e.Error()
	if !strings.Contains(msg, "pkg") || !strings.Contains(msg, "main") {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestRefNotFoundError(t *testing.T) {
	e := &RefNotFoundError{Package: "pkg", Ref: "v1.2.3", OwnerRepo: "owner/repo"}
	msg := e.Error()
	if !strings.Contains(msg, "pkg") || !strings.Contains(msg, "v1.2.3") || !strings.Contains(msg, "owner/repo") {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestNoMatchingVersionError(t *testing.T) {
	e := &NoMatchingVersionError{Package: "pkg", VersionRange: "^1.0.0", Detail: "no tags"}
	msg := e.Error()
	if !strings.Contains(msg, "^1.0.0") || !strings.Contains(msg, "no tags") {
		t.Errorf("unexpected message: %q", msg)
	}
}

// ---------------------------------------------------------------------------
// DefaultBuildOptions
// ---------------------------------------------------------------------------

func TestDefaultBuildOptions(t *testing.T) {
	opts := DefaultBuildOptions()
	if opts.Concurrency <= 0 {
		t.Errorf("expected positive Concurrency, got %d", opts.Concurrency)
	}
	if opts.DryRun {
		t.Error("DryRun should default to false")
	}
}

// ---------------------------------------------------------------------------
// ResolveResult.OK
// ---------------------------------------------------------------------------

func TestResolveResultOK(t *testing.T) {
	ok := ResolveResult{Entries: []ResolvedPackage{{}}, Errors: nil}
	if !ok.OK() {
		t.Error("expected OK")
	}
	notOk := ResolveResult{Errors: [][2]string{{"pkg", "failed"}}}
	if notOk.OK() {
		t.Error("expected not OK")
	}
}

// ---------------------------------------------------------------------------
// stripRefPrefix
// ---------------------------------------------------------------------------

func TestStripRefPrefix(t *testing.T) {
	cases := []struct{ in, out string }{
		{"refs/tags/v1.2.3", "v1.2.3"},
		{"refs/heads/main", "main"},
		{"v1.0.0", "v1.0.0"},
		{"", ""},
	}
	for _, c := range cases {
		got := stripRefPrefix(c.in)
		if got != c.out {
			t.Errorf("stripRefPrefix(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}
