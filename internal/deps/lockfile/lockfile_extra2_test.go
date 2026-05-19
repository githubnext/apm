package lockfile

import (
	"testing"
)

func TestAddDependency_MultipleEntries(t *testing.T) {
	lf := NewLockFile()
	dep1 := &LockedDependency{RepoURL: "github.com/a/b", ResolvedCommit: "abc123"}
	dep2 := &LockedDependency{RepoURL: "github.com/c/d", ResolvedCommit: "def456"}
	lf.AddDependency(dep1)
	lf.AddDependency(dep2)
	all := lf.GetAllDependencies()
	if len(all) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(all))
	}
}

func TestHasDependency_AfterAdd(t *testing.T) {
	lf := NewLockFile()
	dep := &LockedDependency{RepoURL: "github.com/x/y", ResolvedCommit: "sha1"}
	lf.AddDependency(dep)
	key := dep.GetUniqueKey()
	if !lf.HasDependency(key) {
		t.Errorf("expected HasDependency true for %q", key)
	}
}

func TestGetDependency_Returns(t *testing.T) {
	lf := NewLockFile()
	dep := &LockedDependency{RepoURL: "github.com/x/y", ResolvedCommit: "sha1"}
	lf.AddDependency(dep)
	key := dep.GetUniqueKey()
	got := lf.GetDependency(key)
	if got == nil {
		t.Fatal("expected non-nil dependency")
	}
	if got.RepoURL != dep.RepoURL {
		t.Errorf("expected %q, got %q", dep.RepoURL, got.RepoURL)
	}
}

func TestGetAllDependencies_OrderedByDepth(t *testing.T) {
	lf := NewLockFile()
	lf.AddDependency(&LockedDependency{RepoURL: "github.com/a/a", Depth: 2})
	lf.AddDependency(&LockedDependency{RepoURL: "github.com/b/b", Depth: 1})
	all := lf.GetAllDependencies()
	if all[0].Depth > all[1].Depth {
		t.Error("expected dependencies sorted by depth ascending")
	}
}

func TestLockedDependency_IsDevField(t *testing.T) {
	dep := &LockedDependency{RepoURL: "github.com/d/d", IsDev: true}
	if !dep.IsDev {
		t.Error("expected IsDev true")
	}
}

func TestLockedDependency_ToDict_RepoURL(t *testing.T) {
	dep := &LockedDependency{RepoURL: "github.com/a/b", ResolvedCommit: "abc", Depth: 1}
	d := dep.ToDict()
	if d["repo_url"] != "github.com/a/b" {
		t.Errorf("expected repo_url=github.com/a/b, got %v", d["repo_url"])
	}
}

func TestLockedDependency_ToDict_OmitsDepthOne(t *testing.T) {
	dep := &LockedDependency{RepoURL: "github.com/a/b", Depth: 1}
	d := dep.ToDict()
	if _, ok := d["depth"]; ok {
		t.Error("expected depth omitted for depth=1")
	}
}

func TestLockedDependency_ToDict_IncludesDepthTwo(t *testing.T) {
	dep := &LockedDependency{RepoURL: "github.com/a/b", Depth: 2}
	d := dep.ToDict()
	if _, ok := d["depth"]; !ok {
		t.Error("expected depth present for depth=2")
	}
}

func TestGetPackageDependencies_ExcludesSelf2(t *testing.T) {
	lf := NewLockFile()
	lf.AddDependency(&LockedDependency{RepoURL: "github.com/a/b", LocalPath: "."})
	lf.AddDependency(&LockedDependency{RepoURL: "github.com/c/d"})
	pkgDeps := lf.GetPackageDependencies()
	for _, d := range pkgDeps {
		if d.LocalPath == "." {
			t.Error("GetPackageDependencies should exclude self (LocalPath='.')")
		}
	}
}

func TestIsSemanticalllyEquivalent_DifferentDepCount(t *testing.T) {
	lf1 := NewLockFile()
	lf2 := NewLockFile()
	lf1.AddDependency(&LockedDependency{RepoURL: "github.com/a/b", ResolvedCommit: "sha1"})
	if lf1.IsSemanticalllyEquivalent(lf2) {
		t.Error("expected non-equivalent when dep counts differ")
	}
}

func TestFromYAML_BasicParsing(t *testing.T) {
	yaml := "version: 1\ngenerated_at: 2024-01-01T00:00:00Z\ndependencies:\n"
	lf, err := FromYAML(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf == nil {
		t.Fatal("expected non-nil lockfile")
	}
}

func TestFromYAML_EmptyDeps(t *testing.T) {
	yaml := ""
	lf, err := FromYAML(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.GetAllDependencies()) != 0 {
		t.Error("expected no dependencies from empty YAML")
	}
}
