package marketplace

import "testing"

func TestIsValidAlias(t *testing.T) {
	valid := []string{"foo", "my-pkg", "pkg.name", "Pkg_123", "a"}
	for _, v := range valid {
		if !IsValidAlias(v) {
			t.Errorf("IsValidAlias(%q) = false, want true", v)
		}
	}
	invalid := []string{"", "has space", "has/slash", "has@at", "has#hash"}
	for _, v := range invalid {
		if IsValidAlias(v) {
			t.Errorf("IsValidAlias(%q) = true, want false", v)
		}
	}
}

func TestMarketplaceEntryStruct(t *testing.T) {
	e := MarketplaceEntry{
		Alias:  "mypkg",
		URL:    "github.com/owner/repo",
		Branch: "main",
	}
	if e.Alias != "mypkg" {
		t.Errorf("unexpected alias %q", e.Alias)
	}
	if e.Default {
		t.Error("expected Default false")
	}
}
