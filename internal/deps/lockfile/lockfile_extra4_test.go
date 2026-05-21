package lockfile

import (
	"testing"
)

func TestLockedDependency_GetUniqueKey_LocalPath(t *testing.T) {
	d := &LockedDependency{Source: "local", LocalPath: "/some/path"}
	if d.GetUniqueKey() != "/some/path" {
		t.Errorf("expected /some/path, got %s", d.GetUniqueKey())
	}
}

func TestLockedDependency_GetUniqueKey_Virtual(t *testing.T) {
	d := &LockedDependency{IsVirtual: true, VirtualPath: "vp", RepoURL: "https://gh.com/o/r"}
	want := "https://gh.com/o/r/vp"
	if d.GetUniqueKey() != want {
		t.Errorf("expected %s, got %s", want, d.GetUniqueKey())
	}
}

func TestLockedDependency_GetUniqueKey_Default(t *testing.T) {
	d := &LockedDependency{RepoURL: "https://gh.com/o/r"}
	if d.GetUniqueKey() != "https://gh.com/o/r" {
		t.Errorf("expected repo url, got %s", d.GetUniqueKey())
	}
}

func TestNewLockFile_Empty(t *testing.T) {
	lf := NewLockFile()
	if lf == nil {
		t.Fatal("expected non-nil lockfile")
	}
	if len(lf.GetAllDependencies()) != 0 {
		t.Errorf("expected empty deps")
	}
}

func TestLockFile_HasDependency_False(t *testing.T) {
	lf := NewLockFile()
	if lf.HasDependency("missing") {
		t.Error("expected false for missing key")
	}
}

func TestLockFile_AddAndGet(t *testing.T) {
	lf := NewLockFile()
	d := &LockedDependency{RepoURL: "https://gh.com/o/r"}
	lf.AddDependency(d)
	got := lf.GetDependency(d.GetUniqueKey())
	if got == nil {
		t.Fatal("expected non-nil dep")
	}
}

func TestLockFile_HasDependency_True(t *testing.T) {
	lf := NewLockFile()
	d := &LockedDependency{RepoURL: "https://gh.com/o/r2"}
	lf.AddDependency(d)
	if !lf.HasDependency(d.GetUniqueKey()) {
		t.Error("expected true for added dep")
	}
}

func TestLockFile_GetAllDependencies_Multiple(t *testing.T) {
	lf := NewLockFile()
	for _, url := range []string{"https://a.com/o/r1", "https://a.com/o/r2", "https://a.com/o/r3"} {
		lf.AddDependency(&LockedDependency{RepoURL: url})
	}
	if len(lf.GetAllDependencies()) != 3 {
		t.Errorf("expected 3 deps, got %d", len(lf.GetAllDependencies()))
	}
}

func TestLockedDependency_IsDevDefault(t *testing.T) {
	d := &LockedDependency{}
	if d.IsDev {
		t.Error("expected IsDev=false by default")
	}
}

func TestLockedDependency_IsInsecureDefault(t *testing.T) {
	d := &LockedDependency{}
	if d.IsInsecure {
		t.Error("expected IsInsecure=false by default")
	}
}

func TestLockedDependency_Fields(t *testing.T) {
	d := &LockedDependency{
		RepoURL:        "https://gh.com/o/r",
		Host:           "gh.com",
		ResolvedCommit: "abc123",
		Version:        "1.0.0",
		IsDev:          true,
	}
	if d.Host != "gh.com" {
		t.Errorf("unexpected host: %s", d.Host)
	}
	if !d.IsDev {
		t.Error("expected IsDev=true")
	}
}

func TestGetLockfilePath_NonEmpty(t *testing.T) {
	p := GetLockfilePath("/some/root")
	if p == "" {
		t.Error("expected non-empty path")
	}
}
