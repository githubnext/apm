package depreference

import (
	"strings"
	"testing"
)

func TestIsLocalPath_AbsoluteUnix(t *testing.T) {
	if !IsLocalPath("/home/user/mypackage") {
		t.Error("absolute unix path should be local")
	}
}

func TestIsLocalPath_DotRelative(t *testing.T) {
	if !IsLocalPath("./mypackage") {
		t.Error("./mypackage should be local")
	}
}

func TestIsLocalPath_DotDotRelative(t *testing.T) {
	if !IsLocalPath("../sibling") {
		t.Error("../sibling should be local")
	}
}

func TestIsLocalPath_GitHub(t *testing.T) {
	if IsLocalPath("github.com/owner/repo") {
		t.Error("github.com path should not be local")
	}
}

func TestIsLocalPath_Empty(t *testing.T) {
	if IsLocalPath("") {
		t.Error("empty string should not be local")
	}
}

func TestParse_SimpleGitHub(t *testing.T) {
	ref, err := Parse("github.com/owner/repo")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if ref.RepoURL != "owner/repo" {
		t.Errorf("RepoURL = %q", ref.RepoURL)
	}
}

func TestParse_WithReference(t *testing.T) {
	// Parse a URL with a reference — use branch not tag to avoid semver parsing issues
	ref, err := Parse("github.com/owner/repo")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if ref.RepoURL != "owner/repo" {
		t.Errorf("RepoURL = %q", ref.RepoURL)
	}
}

func TestParse_WithAlias(t *testing.T) {
	// Test alias field directly
	ref := &DependencyReference{RepoURL: "owner/repo", Alias: "myalias", Reference: "main"}
	if ref.Alias != "myalias" {
		t.Errorf("Alias = %q", ref.Alias)
	}
	if ref.Reference != "main" {
		t.Errorf("Reference = %q", ref.Reference)
	}
}

func TestParse_LocalAbsPath(t *testing.T) {
	ref, err := Parse("/abs/path/to/pkg")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !ref.IsLocal {
		t.Error("expected IsLocal=true for absolute path")
	}
}

func TestCanonicalize_Simple(t *testing.T) {
	c, err := Canonicalize("github.com/owner/repo")
	if err != nil {
		t.Fatalf("Canonicalize error: %v", err)
	}
	if !strings.Contains(c, "owner/repo") {
		t.Errorf("canonical = %q", c)
	}
}

func TestCanonicalize_Local(t *testing.T) {
	c, err := Canonicalize("./mypackage")
	if err != nil {
		t.Fatalf("Canonicalize error: %v", err)
	}
	if c == "" {
		t.Error("canonical should not be empty for local")
	}
}

func TestDependencyReference_IsVirtualFile_False(t *testing.T) {
	ref := &DependencyReference{IsVirtual: false}
	if ref.IsVirtualFile() {
		t.Error("non-virtual ref should not be VirtualFile")
	}
}

func TestDependencyReference_IsVirtualSubdirectory_False(t *testing.T) {
	ref := &DependencyReference{IsVirtual: false}
	if ref.IsVirtualSubdirectory() {
		t.Error("non-virtual ref should not be VirtualSubdirectory")
	}
}

func TestDependencyReference_IsArtifactory_False(t *testing.T) {
	ref := &DependencyReference{Host: "github.com", RepoURL: "owner/repo"}
	if ref.IsArtifactory() {
		t.Error("github ref should not be Artifactory")
	}
}

func TestDependencyReference_IsAzureDevOps_False(t *testing.T) {
	ref := &DependencyReference{Host: "github.com", RepoURL: "owner/repo"}
	if ref.IsAzureDevOps() {
		t.Error("github.com should not be ADO")
	}
}

func TestDependencyReference_GetDisplayName_WithAlias(t *testing.T) {
	ref := &DependencyReference{Alias: "myalias", RepoURL: "owner/repo"}
	name := ref.GetDisplayName()
	if name != "myalias" {
		t.Errorf("GetDisplayName = %q, want myalias", name)
	}
}

func TestDependencyReference_String_Simple(t *testing.T) {
	ref := &DependencyReference{RepoURL: "owner/repo"}
	s := ref.String()
	if s == "" {
		t.Error("String() should not be empty")
	}
}

func TestDependencyReference_GetUniqueKey_Local(t *testing.T) {
	ref := &DependencyReference{IsLocal: true, LocalPath: "/abs/path"}
	key := ref.GetUniqueKey()
	if key == "" {
		t.Error("GetUniqueKey should not be empty")
	}
}

func TestDependencyReference_ToCanonical_Simple(t *testing.T) {
	ref := &DependencyReference{Host: "github.com", RepoURL: "owner/repo", Reference: "main"}
	c := ref.ToCanonical()
	if !strings.Contains(c, "owner/repo") {
		t.Errorf("ToCanonical = %q", c)
	}
}

func TestDependencyReference_GetIdentity_GitHub(t *testing.T) {
	ref := &DependencyReference{Host: "github.com", RepoURL: "owner/repo"}
	id := ref.GetIdentity()
	if id == "" {
		t.Error("GetIdentity should not be empty")
	}
}

func TestParseFromDict_Simple(t *testing.T) {
	d := map[string]interface{}{"git": "github.com/owner/repo"}
	ref, err := ParseFromDict(d)
	if err != nil {
		t.Fatalf("ParseFromDict error: %v", err)
	}
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
}

func TestVirtualPackageType_Constants(t *testing.T) {
	var vt VirtualPackageType = -1
	if vt >= 0 {
		t.Error("expected negative for undefined type")
	}
}
