package marketplace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidAlias_AllowsDots(t *testing.T) {
	if !IsValidAlias("my.marketplace") {
		t.Error("expected dot to be valid in alias")
	}
}

func TestIsValidAlias_AllowsUnderscore(t *testing.T) {
	if !IsValidAlias("my_market") {
		t.Error("expected underscore to be valid")
	}
}

func TestIsValidAlias_RejectsSlash(t *testing.T) {
	if IsValidAlias("my/market") {
		t.Error("slash should be invalid")
	}
}

func TestIsValidAlias_RejectsSpace(t *testing.T) {
	if IsValidAlias("my market") {
		t.Error("space should be invalid")
	}
}

func TestIsValidAlias_SingleChar(t *testing.T) {
	if !IsValidAlias("a") {
		t.Error("single char should be valid")
	}
}

func TestIsValidAlias_Numeric(t *testing.T) {
	if !IsValidAlias("123") {
		t.Error("numeric alias should be valid")
	}
}

func TestAddOptionsFields(t *testing.T) {
	o := AddOptions{
		ProjectRoot: "/tmp",
		Alias:       "test",
		URL:         "https://example.com",
		Branch:      "main",
		SetDefault:  true,
		Force:       false,
	}
	if o.ProjectRoot != "/tmp" {
		t.Errorf("ProjectRoot = %q", o.ProjectRoot)
	}
	if !o.SetDefault {
		t.Error("SetDefault should be true")
	}
	if o.Force {
		t.Error("Force should be false")
	}
}

func TestAddResultFields(t *testing.T) {
	r := AddResult{Alias: "prod", URL: "https://p.example.com", Branch: "v1", Created: true}
	if !r.Created {
		t.Error("Created should be true")
	}
	if r.Alias != "prod" {
		t.Errorf("Alias = %q", r.Alias)
	}
	if r.Branch != "v1" {
		t.Errorf("Branch = %q", r.Branch)
	}
}

func TestListOptionsFields(t *testing.T) {
	o := ListOptions{ProjectRoot: "/repo"}
	if o.ProjectRoot != "/repo" {
		t.Errorf("ProjectRoot = %q", o.ProjectRoot)
	}
}

func TestListResultZero(t *testing.T) {
	r := ListResult{}
	if len(r.Entries) != 0 {
		t.Error("expected empty entries")
	}
}

func TestValidateOptionsFields(t *testing.T) {
	o := ValidateOptions{ProjectRoot: "/repo", Alias: "prod"}
	if o.Alias != "prod" {
		t.Errorf("Alias = %q", o.Alias)
	}
}

func TestValidateResultZero(t *testing.T) {
	r := ValidateResult{}
	if r.Valid {
		t.Error("zero ValidateResult should have Valid=false")
	}
}

func TestBrowseOptionsFields(t *testing.T) {
	o := BrowseOptions{ProjectRoot: "/r", Alias: "a"}
	if o.ProjectRoot != "/r" || o.Alias != "a" {
		t.Errorf("unexpected fields: %+v", o)
	}
}

func TestBrowseResultNotNil(t *testing.T) {
	r := &BrowseResult{}
	if r == nil {
		t.Fatal("nil BrowseResult")
	}
}

func TestPackageSummaryFields(t *testing.T) {
	p := PackageSummary{Name: "mypkg", Version: "1.0.0"}
	if p.Name != "mypkg" {
		t.Errorf("Name = %q", p.Name)
	}
	if p.Version != "1.0.0" {
		t.Errorf("Version = %q", p.Version)
	}
}

func TestMarketplaceConfigFieldsExtra2(t *testing.T) {
	mc := MarketplaceConfig{Alias: "local", URL: "https://local.test", Branch: "dev", Default: true}
	if mc.Alias != "local" {
		t.Errorf("Alias = %q", mc.Alias)
	}
	if !mc.Default {
		t.Error("Default should be true")
	}
}

func TestAdd_MissingProjectRoot(t *testing.T) {
	_, err := Add(AddOptions{Alias: "x", URL: "https://x.com"})
	if err == nil {
		t.Error("expected error for missing ProjectRoot")
	}
}

func TestRemoveOptions_Fields(t *testing.T) {
	o := RemoveOptions{ProjectRoot: "/r", Alias: "a"}
	if o.Alias != "a" {
		t.Errorf("Alias = %q", o.Alias)
	}
}

func TestAddAndRemoveRoundtrip(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".apm")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	res, err := Add(AddOptions{
		ProjectRoot: dir,
		Alias:       "ci",
		URL:         "https://ci.example.com",
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if res.Alias != "ci" {
		t.Errorf("Alias = %q", res.Alias)
	}
	if err := Remove(RemoveOptions{ProjectRoot: dir, Alias: "ci"}); err != nil {
		t.Fatalf("Remove: %v", err)
	}
}

func TestIsValidAlias_Hyphen(t *testing.T) {
	if !IsValidAlias("my-market") {
		t.Error("hyphen should be valid")
	}
}

func TestIsValidAlias_Empty(t *testing.T) {
	if IsValidAlias("") {
		t.Error("empty should be invalid")
	}
}

func TestIsValidAlias_SpecialCharsInvalid(t *testing.T) {
	for _, bad := range []string{"@foo", "foo!", "foo#bar", "foo$"} {
		if IsValidAlias(bad) {
			t.Errorf("expected %q to be invalid", bad)
		}
	}
}

func TestMarketplaceEntryFields_Extra2(t *testing.T) {
	e := MarketplaceEntry{Alias: "ent", URL: "https://ent.com", Branch: "stable", Default: false}
	if e.Alias != "ent" || e.Branch != "stable" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestList_MissingProjectRoot(t *testing.T) {
	// List with empty ProjectRoot should not panic
	_, _ = List(ListOptions{})
}

func TestListResult_WithEntries(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".apm")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Add(AddOptions{ProjectRoot: dir, Alias: "a1", URL: "https://a1.com"}); err != nil {
		t.Fatal(err)
	}
	if _, err := Add(AddOptions{ProjectRoot: dir, Alias: "a2", URL: "https://a2.com"}); err != nil {
		t.Fatal(err)
	}
	lr, err := List(ListOptions{ProjectRoot: dir})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(lr.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(lr.Entries))
	}
}

func TestOutdatedResultFields_Extra2(t *testing.T) {
	o := OutdatedResult{Packages: []OutdatedPackage{{Name: "foo", CurrentVersion: "1.0.0"}}}
	if len(o.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(o.Packages))
	}
}

func TestDoctorResultFields_Extra2(t *testing.T) {
	d := DoctorResult{Issues: []string{"issue1"}, Fixed: []string{"fix1"}}
	if len(d.Issues) != 1 {
		t.Error("Issues should have 1 entry")
	}
	if len(d.Fixed) != 1 {
		t.Error("Fixed should have 1 entry")
	}
}

func TestUpdateOptions_Fields(t *testing.T) {
	o := UpdateOptions{ProjectRoot: "/r"}
	if o.ProjectRoot != "/r" {
		t.Errorf("ProjectRoot = %q", o.ProjectRoot)
	}
}
