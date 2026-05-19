package cleanuphelper

import "testing"

func TestValidateDeployPath_MultipleValidPrefixes(t *testing.T) {
	prefixes := []string{".github/", ".cursor/", ".claude/"}
	cases := []struct {
		path  string
		valid bool
	}{
		{".github/copilot-instructions.md", true},
		{".cursor/rules.md", true},
		{".claude/settings.json", true},
		{"src/main.go", false},
		{"../etc/passwd", false},
	}
	for _, tc := range cases {
		got := ValidateDeployPath(tc.path, "/project", prefixes)
		if got != tc.valid {
			t.Errorf("ValidateDeployPath(%q): got %v, want %v", tc.path, got, tc.valid)
		}
	}
}

func TestValidateDeployPath_EmptyPrefixes(t *testing.T) {
	// No prefixes -> nothing is valid (even safe paths).
	got := ValidateDeployPath(".github/file.md", "/project", []string{})
	if got {
		t.Error("expected false with empty prefixes")
	}
}

func TestValidateDeployPath_NilPrefixes(t *testing.T) {
	got := ValidateDeployPath(".github/file.md", "/project", nil)
	if got {
		t.Error("expected false with nil prefixes")
	}
}

func TestValidateDeployPath_DotDotHidden(t *testing.T) {
	// Traversal hidden inside a longer path.
	got := ValidateDeployPath(".github/../etc/passwd", "/project", []string{".github/"})
	if got {
		t.Error("expected false for path with '..' component")
	}
}

func TestDiagnosticCollector_MultipleWarnings(t *testing.T) {
	dc := &DiagnosticCollector{}
	dc.Warn("pkg-a", "warning one")
	dc.Warn("pkg-b", "warning two")
	dc.Warn("pkg-a", "warning three")
	if len(dc.Warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d", len(dc.Warnings))
	}
	if dc.Warnings[0].Package != "pkg-a" || dc.Warnings[0].Message != "warning one" {
		t.Error("first warning mismatch")
	}
	if dc.Warnings[1].Package != "pkg-b" {
		t.Error("second warning package mismatch")
	}
	if dc.Warnings[2].Message != "warning three" {
		t.Error("third warning message mismatch")
	}
}

func TestDiagnosticCollector_ZeroValue(t *testing.T) {
	var dc DiagnosticCollector
	if len(dc.Warnings) != 0 {
		t.Error("zero-value DiagnosticCollector should have no warnings")
	}
}

func TestDiagnostic_Fields(t *testing.T) {
	d := Diagnostic{Package: "my-pkg", Message: "some msg"}
	if d.Package != "my-pkg" || d.Message != "some msg" {
		t.Error("Diagnostic fields not set correctly")
	}
}

func TestCleanupResult_ZeroValue(t *testing.T) {
	var r CleanupResult
	if len(r.Deleted) != 0 || len(r.Failed) != 0 ||
		len(r.SkippedUserEdit) != 0 || len(r.SkippedUnmanaged) != 0 {
		t.Error("zero-value CleanupResult should have empty slices")
	}
}

func TestValidateDeployPath_CoworkURIAlwaysRejected(t *testing.T) {
	// cowork:// URIs are rejected regardless of prefixes.
	got := ValidateDeployPath("cowork://some/path", "/project", []string{"cowork://"})
	if got {
		t.Error("cowork:// URI should always be rejected by ValidateDeployPath")
	}
}
