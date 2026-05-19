package scope

import "testing"

func TestInstallScope_String_Variants(t *testing.T) {
	tests := []struct {
		scope InstallScope
		want  string
	}{
		{ScopeProject, "project"},
		{ScopeUser, "user"},
		{InstallScope(99), "project"}, // unknown falls back to "project"
	}
	for _, tc := range tests {
		got := tc.scope.String()
		if got != tc.want {
			t.Errorf("InstallScope(%d).String() = %q, want %q", tc.scope, got, tc.want)
		}
	}
}

func TestParseScope_CaseInsensitive(t *testing.T) {
	cases := []struct {
		input string
		scope InstallScope
		ok    bool
	}{
		{"USER", ScopeUser, true},
		{"User", ScopeUser, true},
		{"PROJECT", ScopeProject, true},
		{"Project", ScopeProject, true},
		{"invalid", ScopeProject, false},
		{"", ScopeProject, false},
	}
	for _, tc := range cases {
		got, ok := ParseScope(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseScope(%q): ok=%v, want %v", tc.input, ok, tc.ok)
		}
		if got != tc.scope {
			t.Errorf("ParseScope(%q): scope=%v, want %v", tc.input, got, tc.scope)
		}
	}
}

func TestGetDeployRoot_ProjectCWD(t *testing.T) {
	root, err := GetDeployRoot(ScopeProject)
	if err != nil {
		t.Fatalf("GetDeployRoot(project): %v", err)
	}
	if root == "" {
		t.Error("GetDeployRoot(project) returned empty string")
	}
}

func TestGetAPMDir_UserHasAPMSuffix(t *testing.T) {
	dir, err := GetAPMDir(ScopeUser)
	if err != nil {
		t.Fatalf("GetAPMDir(user): %v", err)
	}
	// Must end with ".apm"
	if len(dir) < 4 || dir[len(dir)-4:] != ".apm" {
		t.Errorf("GetAPMDir(user) = %q, expected to end with '.apm'", dir)
	}
}

func TestGetAPMDir_ProjectNotEmpty(t *testing.T) {
	dir, err := GetAPMDir(ScopeProject)
	if err != nil {
		t.Fatalf("GetAPMDir(project): %v", err)
	}
	if dir == "" {
		t.Error("GetAPMDir(project) returned empty string")
	}
}

func TestScopeIota_Values(t *testing.T) {
	if ScopeProject != 0 {
		t.Error("ScopeProject should be 0")
	}
	if ScopeUser != 1 {
		t.Error("ScopeUser should be 1")
	}
}

func TestInstallScope_Distinctness(t *testing.T) {
	if ScopeProject == ScopeUser {
		t.Error("ScopeProject and ScopeUser must be distinct")
	}
}

func TestGetLockfileDir_MatchesAPMDir(t *testing.T) {
	apmDir, err := GetAPMDir(ScopeProject)
	if err != nil {
		t.Fatalf("GetAPMDir: %v", err)
	}
	lockDir, err := GetLockfileDir(ScopeProject)
	if err != nil {
		t.Fatalf("GetLockfileDir: %v", err)
	}
	if lockDir != apmDir {
		t.Errorf("GetLockfileDir should equal GetAPMDir: got %q vs %q", lockDir, apmDir)
	}
}
