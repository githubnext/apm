package marketplace

import "testing"

func TestIsValidAlias_Empty_Extra3(t *testing.T) {
	if IsValidAlias("") {
		t.Error("empty string should not be a valid alias")
	}
}

func TestIsValidAlias_LongAlpha_Extra3(t *testing.T) {
	if !IsValidAlias("my-marketplace-123") {
		t.Error("alphanumeric with hyphens should be valid")
	}
}

func TestIsValidAlias_RejectsAt_Extra3(t *testing.T) {
	if IsValidAlias("bad@alias") {
		t.Error("alias with @ should be invalid")
	}
}

func TestIsValidAlias_RejectsHash_Extra3(t *testing.T) {
	if IsValidAlias("bad#alias") {
		t.Error("alias with # should be invalid")
	}
}

func TestMarketplaceConfig_ZeroValue_Extra3(t *testing.T) {
	var cfg MarketplaceConfig
	if cfg.Alias != "" {
		t.Errorf("Alias should be empty, got %q", cfg.Alias)
	}
}

func TestMarketplaceEntry_ZeroValue_Extra3(t *testing.T) {
	var e MarketplaceEntry
	if e.Alias != "" {
		t.Errorf("Alias should be empty, got %q", e.Alias)
	}
}

func TestAddOptions_Struct_Extra3(t *testing.T) {
	opts := AddOptions{
		Alias:       "my-alias",
		URL:         "https://example.com",
		ProjectRoot: "/project",
	}
	if opts.Alias != "my-alias" {
		t.Errorf("Alias = %q, want my-alias", opts.Alias)
	}
}

func TestRemoveOptions_Struct_Extra3(t *testing.T) {
	opts := RemoveOptions{Alias: "alias-to-remove", ProjectRoot: "/root"}
	if opts.Alias != "alias-to-remove" {
		t.Errorf("Alias = %q", opts.Alias)
	}
}

func TestListOptions_Struct_Extra3(t *testing.T) {
	opts := ListOptions{ProjectRoot: "/root", JSON: true}
	if !opts.JSON {
		t.Error("JSON should be true")
	}
}

func TestUpdateOptions_Struct_Extra3(t *testing.T) {
	opts := UpdateOptions{ProjectRoot: "/root"}
	if opts.ProjectRoot != "/root" {
		t.Errorf("ProjectRoot = %q, want /root", opts.ProjectRoot)
	}
}

func TestBrowseOptions_Struct_Extra3(t *testing.T) {
	opts := BrowseOptions{ProjectRoot: "/root"}
	if opts.ProjectRoot != "/root" {
		t.Errorf("ProjectRoot = %q, want /root", opts.ProjectRoot)
	}
}

func TestPackageSummary_Struct_Extra3(t *testing.T) {
	ps := PackageSummary{Name: "my-pkg", Description: "A package"}
	if ps.Name != "my-pkg" {
		t.Errorf("Name = %q, want my-pkg", ps.Name)
	}
}
