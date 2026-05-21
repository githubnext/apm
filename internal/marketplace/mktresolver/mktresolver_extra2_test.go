package mktresolver

import (
	"testing"
)

func TestParseMarketplaceRef_ValidSpec(t *testing.T) {
	ref := ParseMarketplaceRef("myplugin@mymkt")
	if ref == nil {
		t.Fatal("expected non-nil result")
	}
	if ref.Name != "myplugin" {
		t.Errorf("Name = %q, want myplugin", ref.Name)
	}
	if ref.Marketplace != "mymkt" {
		t.Errorf("Marketplace = %q, want mymkt", ref.Marketplace)
	}
	if ref.Ref != "" {
		t.Errorf("Ref should be empty, got %q", ref.Ref)
	}
}

func TestParseMarketplaceRef_WithRef(t *testing.T) {
	ref := ParseMarketplaceRef("plugin@mkt#v1.2")
	if ref == nil {
		t.Fatal("expected non-nil result")
	}
	if ref.Ref != "v1.2" {
		t.Errorf("Ref = %q, want v1.2", ref.Ref)
	}
}

func TestParseMarketplaceRef_NoAtSign(t *testing.T) {
	ref := ParseMarketplaceRef("justname")
	if ref != nil {
		t.Error("expected nil for non-marketplace spec")
	}
}

func TestIsMarketplaceRef_True(t *testing.T) {
	if !IsMarketplaceRef("tool@market") {
		t.Error("expected true for tool@market")
	}
}

func TestIsMarketplaceRef_False(t *testing.T) {
	if IsMarketplaceRef("owner/repo") {
		t.Error("expected false for owner/repo")
	}
}

func TestIsSemverRange_CaretOperator(t *testing.T) {
	if !IsSemverRange("^1.0.0") {
		t.Error("expected true for caret range")
	}
}

func TestIsSemverRange_PlainTag(t *testing.T) {
	if IsSemverRange("v1.0.0") {
		t.Error("expected false for plain tag")
	}
}

func TestNormalizeOwnerRepoSlug_UpperCase(t *testing.T) {
	got := NormalizeOwnerRepoSlug("Owner/Repo")
	want := "owner/repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMarketplaceProjectSlug_Basic(t *testing.T) {
	slug := MarketplaceProjectSlug("owner", "repo")
	if slug == "" {
		t.Error("expected non-empty slug")
	}
}

func TestMarketplaceHostNeedsExplicitGitPath_GitLab(t *testing.T) {
	if !MarketplaceHostNeedsExplicitGitPath("gitlab.com") {
		t.Error("expected gitlab.com to need explicit git path")
	}
}
