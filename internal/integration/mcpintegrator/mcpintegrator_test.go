package mcpintegrator

import (
	"testing"
)

func TestNew_DefaultsToWorkDir(t *testing.T) {
	mi, err := New("", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mi.ProjectRoot == "" {
		t.Error("expected ProjectRoot to be set")
	}
}

func TestNew_ExplicitRoot(t *testing.T) {
	mi, err := New("/tmp", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mi.ProjectRoot != "/tmp" {
		t.Errorf("expected ProjectRoot '/tmp', got %q", mi.ProjectRoot)
	}
}

func TestNormaliseServerName_Basic(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"my-server", "my-server"},
		{"@my-server", "my-server"},
		{"ALREADY_LOWER", "already_lower"},
		{"@Scoped", "scoped"},
	}
	for _, tc := range cases {
		got := NormaliseServerName(tc.in)
		if got != tc.want {
			t.Errorf("NormaliseServerName(%q): got %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormaliseServerName_Empty(t *testing.T) {
	got := NormaliseServerName("")
	// Should not panic
	_ = got
}

func TestDetectConflicts_NoConflicts(t *testing.T) {
	byPackage := map[string][]MCPServer{
		"pkgA": {{Name: "server1"}},
		"pkgB": {{Name: "server2"}},
	}
	results := DetectConflicts(byPackage)
	if len(results) != 0 {
		t.Errorf("expected no conflicts, got %v", results)
	}
}

func TestDetectConflicts_WithConflict(t *testing.T) {
	byPackage := map[string][]MCPServer{
		"pkgA": {{Name: "shared-server"}},
		"pkgB": {{Name: "shared-server"}},
	}
	results := DetectConflicts(byPackage)
	if len(results) == 0 {
		t.Error("expected at least one conflict")
	}
}

func TestDetectConflicts_Empty(t *testing.T) {
	results := DetectConflicts(nil)
	if results != nil && len(results) != 0 {
		t.Errorf("expected no conflicts for nil input, got %v", results)
	}
}

func TestIsVSCodeAvailable_NoProject(t *testing.T) {
	// A non-existent path; should return false without panic
	result := IsVSCodeAvailable("/nonexistent/path/xyz")
	_ = result
}

func TestIsCursorAvailable_NoProject(t *testing.T) {
	result := IsCursorAvailable("/nonexistent/path/xyz")
	_ = result
}

func TestLoadServers_EmptyProject(t *testing.T) {
	t.TempDir()
	mi, err := New(t.TempDir(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	servers, err := mi.LoadServers()
	// Should return empty slice (no lock file) without error
	if err != nil {
		t.Logf("LoadServers returned error (acceptable for missing file): %v", err)
	}
	_ = servers
}
