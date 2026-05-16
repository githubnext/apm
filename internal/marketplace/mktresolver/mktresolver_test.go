package mktresolver_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mktresolver"
)

func TestParseMarketplaceRef_Valid(t *testing.T) {
	ref := mktresolver.ParseMarketplaceRef("my-plugin@my-market")
	if ref == nil {
		t.Fatal("expected parsed ref, got nil")
	}
	if ref.Name != "my-plugin" {
		t.Fatalf("Name mismatch: %q", ref.Name)
	}
	if ref.Marketplace != "my-market" {
		t.Fatalf("Marketplace mismatch: %q", ref.Marketplace)
	}
	if ref.Ref != "" {
		t.Fatalf("Ref should be empty, got %q", ref.Ref)
	}
}

func TestParseMarketplaceRef_WithRef(t *testing.T) {
	ref := mktresolver.ParseMarketplaceRef("my-plugin@my-market#v1.2.0")
	if ref == nil {
		t.Fatal("expected parsed ref, got nil")
	}
	if ref.Ref != "v1.2.0" {
		t.Fatalf("Ref mismatch: %q", ref.Ref)
	}
}

func TestParseMarketplaceRef_Invalid(t *testing.T) {
	if mktresolver.ParseMarketplaceRef("owner/repo") != nil {
		t.Fatal("expected nil for non-marketplace ref")
	}
	if mktresolver.ParseMarketplaceRef("") != nil {
		t.Fatal("expected nil for empty string")
	}
}

func TestIsMarketplaceRef(t *testing.T) {
	if !mktresolver.IsMarketplaceRef("plugin@market") {
		t.Fatal("expected true")
	}
	if mktresolver.IsMarketplaceRef("owner/repo") {
		t.Fatal("expected false for owner/repo")
	}
}

func TestIsSemverRange(t *testing.T) {
	if !mktresolver.IsSemverRange("^1.2.3") {
		t.Fatal("^ should be semver range")
	}
	if !mktresolver.IsSemverRange(">=1.0.0") {
		t.Fatal(">= should be semver range")
	}
	if mktresolver.IsSemverRange("v1.2.3") {
		t.Fatal("plain version should not be semver range")
	}
}

func TestNormalizeOwnerRepoSlug(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"Owner/Repo", "owner/repo"},
		{"owner/repo.git", "owner/repo"},
		{"owner/repo/", "owner/repo"},
	}
	for _, tt := range tests {
		got := mktresolver.NormalizeOwnerRepoSlug(tt.in)
		if got != tt.out {
			t.Errorf("NormalizeOwnerRepoSlug(%q) = %q, want %q", tt.in, got, tt.out)
		}
	}
}

func TestMarketplaceProjectSlug(t *testing.T) {
	got := mktresolver.MarketplaceProjectSlug("Owner", "Repo")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoFieldForMatch_HTTPS(t *testing.T) {
	got := mktresolver.NormalizeRepoFieldForMatch("https://github.com/owner/repo", "github.com")
	if got != "owner/repo" {
		t.Fatalf("expected 'owner/repo', got %q", got)
	}
}

func TestNormalizeRepoFieldForMatch_WrongHost(t *testing.T) {
	got := mktresolver.NormalizeRepoFieldForMatch("https://gitlab.com/owner/repo", "github.com")
	if got != "" {
		t.Fatalf("expected empty for wrong host, got %q", got)
	}
}
