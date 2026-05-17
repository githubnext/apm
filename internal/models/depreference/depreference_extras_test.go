package depreference

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// VirtualType
// ---------------------------------------------------------------------------

func TestVirtualType_not_virtual(t *testing.T) {
	d := &DependencyReference{IsVirtual: false}
	if d.VirtualType() != -1 {
		t.Errorf("expected -1 for non-virtual, got %v", d.VirtualType())
	}
}

func TestVirtualType_virtual_no_path(t *testing.T) {
	d := &DependencyReference{IsVirtual: true, VirtualPath: ""}
	if d.VirtualType() != -1 {
		t.Errorf("expected -1 for virtual with empty path, got %v", d.VirtualType())
	}
}

func TestVirtualType_virtual_file(t *testing.T) {
	d := &DependencyReference{
		IsVirtual:   true,
		VirtualPath: "prompts/code-review.prompt.md",
	}
	if d.VirtualType() != VirtualPackageFile {
		t.Errorf("expected VirtualPackageFile, got %v", d.VirtualType())
	}
}

func TestVirtualType_virtual_subdir(t *testing.T) {
	d := &DependencyReference{
		IsVirtual:   true,
		VirtualPath: "collections/project-planning",
	}
	if d.VirtualType() != VirtualPackageSubdirectory {
		t.Errorf("expected VirtualPackageSubdirectory, got %v", d.VirtualType())
	}
}

// ---------------------------------------------------------------------------
// GetVirtualPackageName
// ---------------------------------------------------------------------------

func TestGetVirtualPackageName_non_virtual(t *testing.T) {
	d := &DependencyReference{
		RepoURL:   "owner/my-plugin",
		IsVirtual: false,
	}
	got := d.GetVirtualPackageName()
	if got != "my-plugin" {
		t.Errorf("expected my-plugin, got %q", got)
	}
}

func TestGetVirtualPackageName_virtual_file(t *testing.T) {
	d := &DependencyReference{
		RepoURL:     "owner/my-repo",
		IsVirtual:   true,
		VirtualPath: "prompts/code-review.prompt.md",
	}
	got := d.GetVirtualPackageName()
	if !strings.HasPrefix(got, "my-repo-") {
		t.Errorf("expected my-repo-* prefix, got %q", got)
	}
	if strings.HasSuffix(got, ".md") {
		t.Errorf("extension should be stripped, got %q", got)
	}
}

func TestGetVirtualPackageName_virtual_subdir(t *testing.T) {
	d := &DependencyReference{
		RepoURL:     "owner/my-repo",
		IsVirtual:   true,
		VirtualPath: "collections/project-planning",
	}
	got := d.GetVirtualPackageName()
	if got != "my-repo-project-planning" {
		t.Errorf("expected my-repo-project-planning, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// GetIdentity
// ---------------------------------------------------------------------------

func TestGetIdentity_local(t *testing.T) {
	d := &DependencyReference{IsLocal: true, LocalPath: "/home/user/my-plugin"}
	got := d.GetIdentity()
	if got != "/home/user/my-plugin" {
		t.Errorf("expected local path, got %q", got)
	}
}

func TestGetIdentity_github_default_host(t *testing.T) {
	d, err := Parse("owner/repo#main")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got := d.GetIdentity()
	if !strings.Contains(got, "owner/repo") {
		t.Errorf("identity should contain owner/repo, got %q", got)
	}
}

func TestGetIdentity_virtual_path(t *testing.T) {
	d := &DependencyReference{
		RepoURL:     "owner/repo",
		IsVirtual:   true,
		VirtualPath: "prompts/review.md",
	}
	got := d.GetIdentity()
	if !strings.Contains(got, "owner/repo") || !strings.Contains(got, "prompts/review.md") {
		t.Errorf("identity missing components, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// ToCloneURL
// ---------------------------------------------------------------------------

func TestToCloneURL_equals_github_url(t *testing.T) {
	d, err := Parse("owner/repo#main")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if d.ToCloneURL() != d.ToGitHubURL() {
		t.Errorf("ToCloneURL != ToGitHubURL: %q vs %q", d.ToCloneURL(), d.ToGitHubURL())
	}
}

func TestToCloneURL_has_https_scheme(t *testing.T) {
	d, err := Parse("owner/repo")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	url := d.ToCloneURL()
	if !strings.HasPrefix(url, "https://") {
		t.Errorf("expected https:// prefix, got %q", url)
	}
}

// ---------------------------------------------------------------------------
// GetDisplayName
// ---------------------------------------------------------------------------

func TestGetDisplayName_with_alias(t *testing.T) {
	d := &DependencyReference{
		RepoURL: "owner/repo",
		Alias:   "my-alias",
	}
	got := d.GetDisplayName()
	if got != "my-alias" {
		t.Errorf("expected alias, got %q", got)
	}
}

func TestGetDisplayName_local(t *testing.T) {
	d := &DependencyReference{
		IsLocal:   true,
		LocalPath: "/path/to/plugin",
	}
	got := d.GetDisplayName()
	if got != "/path/to/plugin" {
		t.Errorf("expected local path, got %q", got)
	}
}

func TestGetDisplayName_virtual(t *testing.T) {
	d := &DependencyReference{
		RepoURL:     "owner/repo",
		IsVirtual:   true,
		VirtualPath: "prompts/review.prompt.md",
	}
	got := d.GetDisplayName()
	// Should return the virtual package name
	if got == "owner/repo" {
		t.Errorf("virtual dep should not return raw repoURL, got %q", got)
	}
}

func TestGetDisplayName_fallback_repo(t *testing.T) {
	d := &DependencyReference{RepoURL: "owner/my-pkg"}
	got := d.GetDisplayName()
	if got != "owner/my-pkg" {
		t.Errorf("expected repoURL fallback, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

func TestString_local(t *testing.T) {
	d := &DependencyReference{IsLocal: true, LocalPath: "./my-plugin"}
	got := d.String()
	if got != "./my-plugin" {
		t.Errorf("expected local path, got %q", got)
	}
}

func TestString_simple(t *testing.T) {
	d, err := Parse("owner/repo#v1.0.0")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got := d.String()
	if !strings.Contains(got, "owner/repo") {
		t.Errorf("expected owner/repo in string, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// IsArtifactory
// ---------------------------------------------------------------------------

func TestIsArtifactory_false(t *testing.T) {
	d := &DependencyReference{RepoURL: "owner/repo"}
	if d.IsArtifactory() {
		t.Error("expected not artifactory")
	}
}

func TestIsArtifactory_true(t *testing.T) {
	d := &DependencyReference{
		RepoURL:           "owner/repo",
		ArtifactoryPrefix: "artifactory.example.com/vcs",
	}
	if !d.IsArtifactory() {
		t.Error("expected artifactory")
	}
}

// ---------------------------------------------------------------------------
// IsLocalPath edge cases
// ---------------------------------------------------------------------------

func TestIsLocalPath_tilde(t *testing.T) {
	if !IsLocalPath("~/my-plugin") {
		t.Error("expected tilde path to be local")
	}
}

func TestIsLocalPath_double_slash(t *testing.T) {
	// double-slash is NOT local (it is a protocol-relative URL)
	if IsLocalPath("//github.com/owner/repo") {
		t.Error("expected double-slash to NOT be local")
	}
}

func TestIsLocalPath_parent(t *testing.T) {
	if !IsLocalPath("../sibling-plugin") {
		t.Error("expected ../ to be local")
	}
}

// ---------------------------------------------------------------------------
// GetUniqueKey
// ---------------------------------------------------------------------------

func TestGetUniqueKey_local(t *testing.T) {
	d := &DependencyReference{IsLocal: true, LocalPath: "/absolute/path"}
	if d.GetUniqueKey() != "/absolute/path" {
		t.Errorf("unexpected key: %q", d.GetUniqueKey())
	}
}

func TestGetUniqueKey_virtual(t *testing.T) {
	d := &DependencyReference{
		RepoURL:     "owner/repo",
		IsVirtual:   true,
		VirtualPath: "prompts/foo.md",
	}
	got := d.GetUniqueKey()
	if got != "owner/repo/prompts/foo.md" {
		t.Errorf("expected owner/repo/prompts/foo.md, got %q", got)
	}
}

func TestGetUniqueKey_plain(t *testing.T) {
	d := &DependencyReference{RepoURL: "owner/repo"}
	if d.GetUniqueKey() != "owner/repo" {
		t.Errorf("expected owner/repo, got %q", d.GetUniqueKey())
	}
}

// ---------------------------------------------------------------------------
// GetCanonicalDependencyString
// ---------------------------------------------------------------------------

func TestGetCanonicalDependencyString_delegates_to_key(t *testing.T) {
	d := &DependencyReference{RepoURL: "owner/repo"}
	if d.GetCanonicalDependencyString() != d.GetUniqueKey() {
		t.Error("GetCanonicalDependencyString should equal GetUniqueKey")
	}
}
