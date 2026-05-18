package mktresolver_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mktresolver"
)

func TestParseMarketplaceRef_NoAt(t *testing.T) {
	if mktresolver.ParseMarketplaceRef("noplugin") != nil {
		t.Error("expected nil for spec with no @")
	}
}

func TestParseMarketplaceRef_AtOnly(t *testing.T) {
	// "@market" has no plugin name -- should be nil
	if mktresolver.ParseMarketplaceRef("@market") != nil {
		t.Error("expected nil for spec with empty plugin name")
	}
}

func TestParseMarketplaceRef_EmptyRef(t *testing.T) {
	ref := mktresolver.ParseMarketplaceRef("plugin@market#")
	// An empty fragment after # may or may not parse; just assert no panic.
	_ = ref
}

func TestNormalizeOwnerRepoSlug_TrailingSlash(t *testing.T) {
	got := mktresolver.NormalizeOwnerRepoSlug("Owner/Repo/")
	if got != "owner/repo" {
		t.Errorf("NormalizeOwnerRepoSlug with trailing slash: got %q", got)
	}
}

func TestNormalizeOwnerRepoSlug_GitSuffix(t *testing.T) {
	got := mktresolver.NormalizeOwnerRepoSlug("OWNER/REPO.git")
	if got != "owner/repo" {
		t.Errorf("NormalizeOwnerRepoSlug .git: got %q", got)
	}
}

func TestNormalizeOwnerRepoSlug_AlreadyLower(t *testing.T) {
	got := mktresolver.NormalizeOwnerRepoSlug("owner/repo")
	if got != "owner/repo" {
		t.Errorf("NormalizeOwnerRepoSlug already lower: got %q", got)
	}
}

func TestMarketplaceProjectSlug_Spaces(t *testing.T) {
	got := mktresolver.MarketplaceProjectSlug("Owner", "Repo")
	if got != "owner/repo" {
		t.Errorf("MarketplaceProjectSlug: got %q", got)
	}
}

func TestIsSemverRange_TildeOperator(t *testing.T) {
	if !mktresolver.IsSemverRange("~1.2.0") {
		t.Error("~ should be a semver range indicator")
	}
}

func TestIsSemverRange_ExactVersion(t *testing.T) {
	if mktresolver.IsSemverRange("1.2.3") {
		t.Error("plain version should not be a semver range")
	}
}

func TestIsSemverRange_EmptyString(t *testing.T) {
	if mktresolver.IsSemverRange("") {
		t.Error("empty string should not be a semver range")
	}
}

func TestNormalizeRepoFieldForMatch_SCP(t *testing.T) {
	// git@ SCP style -- no URL scheme, should not match https prefix stripping
	got := mktresolver.NormalizeRepoFieldForMatch("git@github.com:owner/repo.git", "github.com")
	// May return empty or partial -- just assert no panic
	_ = got
}

func TestNormalizeRepoFieldForMatch_SSHScheme(t *testing.T) {
	got := mktresolver.NormalizeRepoFieldForMatch("ssh://github.com/owner/repo", "github.com")
	if got != "owner/repo" {
		t.Errorf("ssh:// scheme: got %q, want owner/repo", got)
	}
}

func TestNormalizeRepoFieldForMatch_HTTPScheme(t *testing.T) {
	got := mktresolver.NormalizeRepoFieldForMatch("http://github.com/owner/repo", "github.com")
	if got != "owner/repo" {
		t.Errorf("http:// scheme: got %q, want owner/repo", got)
	}
}

func TestGitSourceToCanonical_NoRef(t *testing.T) {
	src := map[string]interface{}{"repo": "Owner/Repo"}
	got := mktresolver.GitSourceToCanonical(src)
	if got != "owner/repo" {
		t.Errorf("GitSourceToCanonical no ref: got %q", got)
	}
}

func TestGitSourceToCanonical_WithRef(t *testing.T) {
	src := map[string]interface{}{"repo": "owner/repo", "ref": "v1.0"}
	got := mktresolver.GitSourceToCanonical(src)
	if got != "owner/repo#v1.0" {
		t.Errorf("GitSourceToCanonical with ref: got %q", got)
	}
}

func TestGitSourceToCanonical_WithVersion(t *testing.T) {
	src := map[string]interface{}{"repo": "owner/repo", "version": "2.0"}
	got := mktresolver.GitSourceToCanonical(src)
	if got != "owner/repo#2.0" {
		t.Errorf("GitSourceToCanonical with version: got %q", got)
	}
}

func TestURLSourceToCanonical_Basic(t *testing.T) {
	src := map[string]interface{}{"url": "https://example.com/plugin.zip"}
	got := mktresolver.URLSourceToCanonical(src)
	if got != "https://example.com/plugin.zip" {
		t.Errorf("URLSourceToCanonical: got %q", got)
	}
}

func TestURLSourceToCanonical_Empty(t *testing.T) {
	src := map[string]interface{}{}
	got := mktresolver.URLSourceToCanonical(src)
	if got != "" {
		t.Errorf("URLSourceToCanonical empty: got %q", got)
	}
}

func TestClassifyPluginSource_GitHub(t *testing.T) {
	src := map[string]interface{}{"github": map[string]interface{}{"repo": "owner/repo"}}
	if mktresolver.ClassifyPluginSource(src) != mktresolver.PluginSourceGitHub {
		t.Error("expected PluginSourceGitHub")
	}
}

func TestClassifyPluginSource_URL(t *testing.T) {
	src := map[string]interface{}{"url": "https://example.com/x.zip"}
	if mktresolver.ClassifyPluginSource(src) != mktresolver.PluginSourceURL {
		t.Error("expected PluginSourceURL")
	}
}

func TestMarketplaceHostNeedsExplicitGitPath_GitHub(t *testing.T) {
	if mktresolver.MarketplaceHostNeedsExplicitGitPath("github.com") {
		t.Error("github.com should not need explicit git path")
	}
}

func TestMarketplaceHostNeedsExplicitGitPath_GHE(t *testing.T) {
	if mktresolver.MarketplaceHostNeedsExplicitGitPath("acme.ghe.com") {
		t.Error("GHE host should not need explicit git path")
	}
}

func TestMarketplaceHostNeedsExplicitGitPath_GitLab(t *testing.T) {
	if !mktresolver.MarketplaceHostNeedsExplicitGitPath("gitlab.com") {
		t.Error("gitlab.com should need explicit git path")
	}
}

func TestIsMarketplaceRef_WithRef(t *testing.T) {
	if !mktresolver.IsMarketplaceRef("myplugin@mymkt#v1.0") {
		t.Error("plugin@mkt#ref should be a marketplace ref")
	}
}

func TestIsMarketplaceRef_OnlyAtSign(t *testing.T) {
	if mktresolver.IsMarketplaceRef("@") {
		t.Error("bare @ should not match")
	}
}
