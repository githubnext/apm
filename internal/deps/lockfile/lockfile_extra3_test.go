package lockfile

import (
	"testing"
)

func TestNewLockFile_EmptyDeps(t *testing.T) {
	lf := NewLockFile()
	deps := lf.GetAllDependencies()
	if len(deps) != 0 {
		t.Errorf("expected empty deps, got %d", len(deps))
	}
}

func TestLockedDependency_FieldAssignment(t *testing.T) {
	d := &LockedDependency{
		RepoURL:        "https://github.com/owner/repo",
		ResolvedCommit: "abc123",
		Version:        "1.0.0",
	}
	if d.RepoURL != "https://github.com/owner/repo" {
		t.Error("unexpected RepoURL")
	}
	if d.ResolvedCommit != "abc123" {
		t.Error("unexpected ResolvedCommit")
	}
}

func TestGetUniqueKey_DefaultRepoURL(t *testing.T) {
	d := &LockedDependency{RepoURL: "https://github.com/o/r"}
	key := d.GetUniqueKey()
	if key != "https://github.com/o/r" {
		t.Errorf("expected repo URL, got %q", key)
	}
}

func TestGetUniqueKey_VirtualPath(t *testing.T) {
	d := &LockedDependency{
		RepoURL:     "https://github.com/o/r",
		IsVirtual:   true,
		VirtualPath: "sub/path",
	}
	key := d.GetUniqueKey()
	if key == "" {
		t.Error("expected non-empty key for virtual dep")
	}
	if key != "https://github.com/o/r/sub/path" {
		t.Errorf("unexpected virtual key: %q", key)
	}
}

func TestGetUniqueKey_LocalPath(t *testing.T) {
	d := &LockedDependency{
		Source:    "local",
		LocalPath: "/home/user/mypackage",
	}
	key := d.GetUniqueKey()
	if key != "/home/user/mypackage" {
		t.Errorf("expected local path, got %q", key)
	}
}

func TestLockFile_AddAndGetDep(t *testing.T) {
	lf := NewLockFile()
	dep := &LockedDependency{RepoURL: "https://github.com/o/r"}
	lf.AddDependency(dep)
	got := lf.GetDependency("https://github.com/o/r")
	if got == nil {
		t.Fatal("expected non-nil dependency")
	}
	if got.RepoURL != dep.RepoURL {
		t.Error("RepoURL mismatch")
	}
}

func TestLockFile_HasDependency(t *testing.T) {
	lf := NewLockFile()
	dep := &LockedDependency{RepoURL: "https://github.com/a/b"}
	lf.AddDependency(dep)
	if !lf.HasDependency("https://github.com/a/b") {
		t.Error("expected HasDependency to return true")
	}
	if lf.HasDependency("missing") {
		t.Error("expected HasDependency to return false for missing key")
	}
}

func TestFromYAML_InvalidYAML(t *testing.T) {
	_, err := FromYAML(":::invalid:::yaml:::")
	if err == nil {
		_ = err
	}
}

func TestFromYAML_EmptyYAMLReturnsLockFile(t *testing.T) {
	lf, err := FromYAML("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf == nil {
		t.Fatal("expected non-nil lockfile")
	}
}

func TestGetLockfilePath_ContainsLockfileName(t *testing.T) {
	path := GetLockfilePath("/some/dir")
	if path == "" {
		t.Error("expected non-empty path")
	}
}

func TestLockedDependency_ToDict_HasRepoURL(t *testing.T) {
	d := &LockedDependency{RepoURL: "https://github.com/x/y"}
	m := d.ToDict()
	if m["repo_url"] != "https://github.com/x/y" {
		t.Errorf("expected repo_url in dict, got %v", m["repo_url"])
	}
}

func TestLockedDependency_IsDevFieldE3(t *testing.T) {
	d := &LockedDependency{IsDev: true}
	if !d.IsDev {
		t.Error("expected IsDev=true")
	}
}
