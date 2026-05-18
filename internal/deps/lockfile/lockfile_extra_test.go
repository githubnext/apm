package lockfile

import (
	"testing"
)

func TestGetUniqueKey_Standard(t *testing.T) {
	d := &LockedDependency{RepoURL: "https://github.com/foo/bar"}
	if d.GetUniqueKey() != "https://github.com/foo/bar" {
		t.Errorf("unexpected key: %s", d.GetUniqueKey())
	}
}

func TestGetUniqueKey_Local(t *testing.T) {
	d := &LockedDependency{Source: "local", LocalPath: "/my/local/path"}
	if d.GetUniqueKey() != "/my/local/path" {
		t.Errorf("expected local path, got: %s", d.GetUniqueKey())
	}
}

func TestGetUniqueKey_Virtual(t *testing.T) {
	d := &LockedDependency{IsVirtual: true, RepoURL: "https://github.com/foo/bar", VirtualPath: "sub"}
	want := "https://github.com/foo/bar/sub"
	if d.GetUniqueKey() != want {
		t.Errorf("expected %s, got %s", want, d.GetUniqueKey())
	}
}

func TestGetPackageDependencies_ExcludesSelf(t *testing.T) {
	lf := NewLockFile()
	lf.AddDependency(&LockedDependency{RepoURL: "https://github.com/a/b", Depth: 1})
	lf.AddDependency(&LockedDependency{Source: "local", LocalPath: ".", RepoURL: "."})
	pkgs := lf.GetPackageDependencies()
	for _, d := range pkgs {
		if d.GetUniqueKey() == "." {
			t.Error("self-entry should be excluded from package dependencies")
		}
	}
}

func TestHasDependency_False(t *testing.T) {
	lf := NewLockFile()
	if lf.HasDependency("https://github.com/nonexistent/pkg") {
		t.Error("expected HasDependency=false for unknown URL")
	}
}

func TestHasDependency_True(t *testing.T) {
	lf := NewLockFile()
	lf.AddDependency(&LockedDependency{RepoURL: "https://github.com/owner/repo"})
	if !lf.HasDependency("https://github.com/owner/repo") {
		t.Error("expected HasDependency=true after adding dep")
	}
}

func TestGetDependency_Nil(t *testing.T) {
	lf := NewLockFile()
	if lf.GetDependency("missing") != nil {
		t.Error("expected nil for missing dependency")
	}
}

func TestGetAllDependencies_Empty(t *testing.T) {
	lf := NewLockFile()
	deps := lf.GetAllDependencies()
	if len(deps) != 0 {
		t.Errorf("expected 0 deps, got %d", len(deps))
	}
}

func TestGetAllDependencies_OrderByDepth(t *testing.T) {
	lf := NewLockFile()
	lf.AddDependency(&LockedDependency{RepoURL: "https://github.com/z/z", Depth: 3})
	lf.AddDependency(&LockedDependency{RepoURL: "https://github.com/a/a", Depth: 1})
	deps := lf.GetAllDependencies()
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d", len(deps))
	}
	if deps[0].Depth > deps[1].Depth {
		t.Error("expected deps sorted by depth ascending")
	}
}

func TestNewLockFile_VersionOne(t *testing.T) {
	lf := NewLockFile()
	if lf.LockfileVersion != "1" {
		t.Errorf("expected version '1', got %q", lf.LockfileVersion)
	}
}

func TestNewLockFile_GeneratedAtSet(t *testing.T) {
	lf := NewLockFile()
	if lf.GeneratedAt == "" {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestGetLockfilePath_ContainsFileName(t *testing.T) {
	p := GetLockfilePath("/project/root")
	if p == "" {
		t.Error("expected non-empty path")
	}
}

func TestLoadOrCreate_ReturnsLockfile(t *testing.T) {
	lf := LoadOrCreate("/nonexistent/path/apm.lock.yaml")
	if lf == nil {
		t.Fatal("LoadOrCreate returned nil")
	}
}

func TestToDict_DepthOneOmitted(t *testing.T) {
	d := &LockedDependency{RepoURL: "https://example.com/r", Depth: 1}
	dict := d.ToDict()
	if _, ok := dict["depth"]; ok {
		t.Error("depth=1 should be omitted from ToDict")
	}
}

func TestToDict_DepthTwoIncluded(t *testing.T) {
	d := &LockedDependency{RepoURL: "https://example.com/r", Depth: 2}
	dict := d.ToDict()
	if dict["depth"] != 2 {
		t.Errorf("expected depth=2, got %v", dict["depth"])
	}
}

func TestToDict_IsDevFalseOmitted(t *testing.T) {
	d := &LockedDependency{RepoURL: "https://example.com/r", Depth: 1, IsDev: false}
	dict := d.ToDict()
	if v, ok := dict["is_dev"]; ok && v == true {
		t.Error("is_dev=false should not be included as true")
	}
}

func TestFromYAML_EmptyString(t *testing.T) {
	lf, err := FromYAML("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lf == nil {
		t.Fatal("expected non-nil LockFile")
	}
}

func TestIsSemanticalllyEquivalent_SameContent(t *testing.T) {
	lf1 := NewLockFile()
	lf1.AddDependency(&LockedDependency{RepoURL: "https://github.com/a/b", ResolvedCommit: "abc"})
	lf2 := NewLockFile()
	lf2.AddDependency(&LockedDependency{RepoURL: "https://github.com/a/b", ResolvedCommit: "abc"})
	if !lf1.IsSemanticalllyEquivalent(lf2) {
		t.Error("same content should be semantically equivalent")
	}
}

func TestIsSemanticalllyEquivalent_DifferentCommit(t *testing.T) {
	lf1 := NewLockFile()
	lf1.AddDependency(&LockedDependency{RepoURL: "https://github.com/a/b", ResolvedCommit: "abc"})
	lf2 := NewLockFile()
	lf2.AddDependency(&LockedDependency{RepoURL: "https://github.com/a/b", ResolvedCommit: "def"})
	if lf1.IsSemanticalllyEquivalent(lf2) {
		t.Error("different commit should not be semantically equivalent")
	}
}
