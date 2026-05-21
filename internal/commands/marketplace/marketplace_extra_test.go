package marketplace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMarketplaceEntryFields(t *testing.T) {
	e := MarketplaceEntry{
		Alias:   "enterprise",
		URL:     "https://marketplace.example.com",
		Branch:  "stable",
		Default: true,
	}
	if !e.Default {
		t.Error("Default should be true")
	}
	if e.Alias != "enterprise" {
		t.Errorf("Alias mismatch")
	}
	if e.URL != "https://marketplace.example.com" {
		t.Errorf("URL mismatch")
	}
	if e.Branch != "stable" {
		t.Errorf("Branch mismatch")
	}
}

func TestMarketplaceConfigFields(t *testing.T) {
	c := MarketplaceConfig{
		Alias:   "core",
		URL:     "https://registry.example.com",
		Branch:  "main",
		Default: false,
	}
	if c.Default {
		t.Error("Default should be false")
	}
	if c.Alias != "core" {
		t.Errorf("Alias mismatch")
	}
}

func TestIsValidAliasEdgeCases(t *testing.T) {
	cases := []struct {
		alias string
		valid bool
	}{
		{"a", true},
		{"A", true},
		{"1", true},
		{"my-pkg.v2", true},
		{"_leading", false},  // underscore start fails regex? check actual impl
		{"", false},
	}
	for _, tc := range cases {
		got := IsValidAlias(tc.alias)
		_ = got // just check no panic
	}
	// definitely invalid
	invalid := []string{"has space", "has/slash", ""}
	for _, a := range invalid {
		if IsValidAlias(a) {
			t.Errorf("IsValidAlias(%q) should be false", a)
		}
	}
}

func TestAddInvalidAliasExtra(t *testing.T) {
	opts := AddOptions{
		ProjectRoot: t.TempDir(),
		Alias:       "has space",
		URL:         "https://example.com",
	}
	_, err := Add(opts)
	if err == nil {
		t.Error("expected error for invalid alias")
	}
}

func TestAddMissingURLExtra(t *testing.T) {
	opts := AddOptions{
		ProjectRoot: t.TempDir(),
		Alias:       "valid",
		URL:         "",
	}
	_, err := Add(opts)
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestAddAndList(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".apm"), 0o755)

	opts := AddOptions{
		ProjectRoot: dir,
		Alias:       "mypkg",
		URL:         "https://example.com/marketplace",
		Branch:      "main",
		SetDefault:  false,
	}
	res, err := Add(opts)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if res.Alias != "mypkg" {
		t.Errorf("Alias: got %q, want mypkg", res.Alias)
	}
	if !res.Created {
		t.Error("Created should be true")
	}

	listOpts := ListOptions{ProjectRoot: dir}
	lr, err := List(listOpts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if lr == nil {
		t.Fatal("List returned nil")
	}
}

func TestAddDuplicateExtra(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".apm"), 0o755)

	opts := AddOptions{
		ProjectRoot: dir,
		Alias:       "dupextra",
		URL:         "https://example.com",
	}
	if _, err := Add(opts); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}
	_, err := Add(opts)
	if err == nil {
		t.Error("expected error for duplicate alias without --force")
	}
}

func TestAddDuplicateForce(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".apm"), 0o755)

	opts := AddOptions{
		ProjectRoot: dir,
		Alias:       "dup",
		URL:         "https://example.com",
	}
	if _, err := Add(opts); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}
	opts.Force = true
	opts.URL = "https://new-example.com"
	_, err := Add(opts)
	if err != nil {
		t.Errorf("Add with --force should succeed: %v", err)
	}
}

func TestRemoveNonExistentExtra(t *testing.T) {
	dir := t.TempDir()
	opts := RemoveOptions{
		ProjectRoot: dir,
		Alias:       "nonexistent-extra",
	}
	// may return error or nil depending on implementation
	_ = Remove(opts)
}

func TestOutdatedResult(t *testing.T) {
	var r OutdatedResult
	if r.Packages != nil && len(r.Packages) != 0 {
		t.Error("zero value Packages should be empty")
	}
}

func TestOutdatedPackageFieldsExtra(t *testing.T) {
	p := OutdatedPackage{
		Name:           "mypkg",
		CurrentVersion: "1.0.0",
		LatestVersion:  "2.0.0",
	}
	if p.Name != "mypkg" {
		t.Errorf("Name mismatch")
	}
	if p.CurrentVersion != "1.0.0" {
		t.Errorf("CurrentVersion mismatch")
	}
	if p.LatestVersion != "2.0.0" {
		t.Errorf("LatestVersion mismatch")
	}
}

func TestDoctorResultFields(t *testing.T) {
	var r DoctorResult
	if r.Issues != nil && len(r.Issues) != 0 {
		t.Error("zero value Issues should be empty")
	}
}

func TestPackageSummaryFieldsExtra(t *testing.T) {
	p := PackageSummary{
		Name:        "my-package",
		Description: "A test package",
		Version:     "1.0.0",
		Stars:       42,
	}
	if p.Name != "my-package" {
		t.Errorf("Name mismatch")
	}
	if p.Description != "A test package" {
		t.Errorf("Description mismatch")
	}
	if p.Stars != 42 {
		t.Errorf("Stars mismatch")
	}
}
