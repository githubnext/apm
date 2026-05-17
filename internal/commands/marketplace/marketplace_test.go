package marketplace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidAlias(t *testing.T) {
	valid := []string{"foo", "my-pkg", "pkg.name", "Pkg_123", "a", "x1", "A-B_C.D"}
	for _, v := range valid {
		if !IsValidAlias(v) {
			t.Errorf("IsValidAlias(%q) = false, want true", v)
		}
	}
	invalid := []string{"", "has space", "has/slash", "has@at", "has#hash", "bad!", "a b"}
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

func TestMarketplaceConfigStruct(t *testing.T) {
	c := MarketplaceConfig{
		Alias:   "core",
		URL:     "https://example.com/marketplace",
		Branch:  "stable",
		Default: true,
	}
	if c.Alias != "core" {
		t.Errorf("unexpected alias %q", c.Alias)
	}
	if !c.Default {
		t.Error("expected Default true")
	}
}

func makeTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func writeMarketplaces(t *testing.T, root string, entries []MarketplaceEntry) {
	t.Helper()
	cfgPath := filepath.Join(root, ".apm", "marketplaces.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgs := make([]MarketplaceConfig, len(entries))
	for i, e := range entries {
		cfgs[i] = MarketplaceConfig{Alias: e.Alias, URL: e.URL, Branch: e.Branch, Default: e.Default}
	}
	data, _ := json.MarshalIndent(cfgs, "", "  ")
	if err := os.WriteFile(cfgPath, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestAddNewMarketplace(t *testing.T) {
	root := makeTestDir(t)
	result, err := Add(AddOptions{
		ProjectRoot: root,
		Alias:       "acme",
		URL:         "https://example.com/acme",
		Branch:      "main",
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.Alias != "acme" {
		t.Errorf("got alias %q, want %q", result.Alias, "acme")
	}
	if !result.Created {
		t.Error("expected Created true")
	}
	if result.URL != "https://example.com/acme" {
		t.Errorf("got URL %q", result.URL)
	}
}

func TestAddInvalidAlias(t *testing.T) {
	root := makeTestDir(t)
	_, err := Add(AddOptions{
		ProjectRoot: root,
		Alias:       "bad alias!",
		URL:         "https://example.com",
	})
	if err == nil {
		t.Error("expected error for invalid alias")
	}
	if !strings.Contains(err.Error(), "invalid marketplace alias") {
		t.Errorf("error should mention invalid alias, got: %v", err)
	}
}

func TestAddMissingURL(t *testing.T) {
	root := makeTestDir(t)
	_, err := Add(AddOptions{
		ProjectRoot: root,
		Alias:       "good-alias",
		URL:         "",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestAddDuplicate(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "acme", URL: "https://example.com/acme"},
	})
	_, err := Add(AddOptions{
		ProjectRoot: root,
		Alias:       "acme",
		URL:         "https://other.com",
	})
	if err == nil {
		t.Error("expected error for duplicate alias")
	}
	if !strings.Contains(err.Error(), "already registered") {
		t.Errorf("error should mention already registered: %v", err)
	}
}

func TestAddForceOverwrite(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "acme", URL: "https://old.com"},
	})
	result, err := Add(AddOptions{
		ProjectRoot: root,
		Alias:       "acme",
		URL:         "https://new.com",
		Force:       true,
	})
	if err != nil {
		t.Fatalf("Add with force: %v", err)
	}
	if result.URL != "https://new.com" {
		t.Errorf("expected overwritten URL, got %q", result.URL)
	}
}

func TestAddMultipleMarketplaces(t *testing.T) {
	root := makeTestDir(t)
	for _, alias := range []string{"alpha", "beta", "gamma"} {
		_, err := Add(AddOptions{
			ProjectRoot: root,
			Alias:       alias,
			URL:         "https://example.com/" + alias,
		})
		if err != nil {
			t.Fatalf("Add %q: %v", alias, err)
		}
	}
	res, err := List(ListOptions{ProjectRoot: root})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Errorf("got %d entries, want 3", len(res.Entries))
	}
}

func TestAddWithSetDefault(t *testing.T) {
	root := makeTestDir(t)
	result, err := Add(AddOptions{
		ProjectRoot: root,
		Alias:       "default-mp",
		URL:         "https://example.com",
		SetDefault:  true,
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.Alias != "default-mp" {
		t.Errorf("unexpected alias %q", result.Alias)
	}
}

func TestAddWithBranch(t *testing.T) {
	root := makeTestDir(t)
	result, err := Add(AddOptions{
		ProjectRoot: root,
		Alias:       "versioned",
		URL:         "https://example.com/mp",
		Branch:      "v2",
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.Branch != "v2" {
		t.Errorf("got branch %q, want %q", result.Branch, "v2")
	}
}

func TestRemoveExisting(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "alpha", URL: "https://a.com"},
		{Alias: "beta", URL: "https://b.com"},
	})
	if err := Remove(RemoveOptions{ProjectRoot: root, Alias: "alpha"}); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	res, _ := List(ListOptions{ProjectRoot: root})
	if len(res.Entries) != 1 {
		t.Errorf("got %d entries after remove, want 1", len(res.Entries))
	}
	if res.Entries[0].Alias != "beta" {
		t.Errorf("unexpected remaining alias %q", res.Entries[0].Alias)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "alpha", URL: "https://a.com"},
	})
	err := Remove(RemoveOptions{ProjectRoot: root, Alias: "nonexistent"})
	if err == nil {
		t.Error("expected error removing nonexistent marketplace")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention not found: %v", err)
	}
}

func TestRemoveNotExistFile(t *testing.T) {
	root := makeTestDir(t)
	err := Remove(RemoveOptions{ProjectRoot: root, Alias: "any"})
	if err == nil {
		t.Error("expected error when config file missing")
	}
}

func TestListEmpty(t *testing.T) {
	root := makeTestDir(t)
	res, err := List(ListOptions{ProjectRoot: root})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(res.Entries) != 0 {
		t.Errorf("got %d entries, want 0", len(res.Entries))
	}
}

func TestListMultiple(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "a", URL: "https://a.com"},
		{Alias: "b", URL: "https://b.com", Default: true},
	})
	res, err := List(ListOptions{ProjectRoot: root})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("got %d entries, want 2", len(res.Entries))
	}
}

func TestListPreservesFields(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "mp1", URL: "https://mp1.com", Branch: "dev", Default: true},
	})
	res, _ := List(ListOptions{ProjectRoot: root})
	if len(res.Entries) != 1 {
		t.Fatal("expected 1 entry")
	}
	e := res.Entries[0]
	if e.Alias != "mp1" || e.URL != "https://mp1.com" || e.Branch != "dev" || !e.Default {
		t.Errorf("fields not preserved: %+v", e)
	}
}

func TestValidateValid(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "ok", URL: "https://example.com/mp"},
	})
	result, err := Validate(ValidateOptions{ProjectRoot: root, Alias: "ok"})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !result.Valid {
		t.Errorf("expected Valid true, got errors: %v", result.Errors)
	}
}

func TestValidateHTTPURL(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "noscheme", URL: "example.com/mp"},
	})
	result, err := Validate(ValidateOptions{ProjectRoot: root, Alias: "noscheme"})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	hasURLWarning := false
	for _, e := range result.Errors {
		if strings.Contains(e, "https://") {
			hasURLWarning = true
		}
	}
	if !hasURLWarning {
		t.Errorf("expected URL warning for URL without https://, errors: %v", result.Errors)
	}
}

func TestValidateNonExistent(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "existing", URL: "https://example.com"},
	})
	_, err := Validate(ValidateOptions{ProjectRoot: root, Alias: "missing"})
	if err == nil {
		t.Error("expected error for missing alias")
	}
}

func TestBrowseReturnsEmpty(t *testing.T) {
	result, err := Browse(BrowseOptions{Alias: "mp", Query: "test"})
	if err != nil {
		t.Fatalf("Browse: %v", err)
	}
	if result == nil {
		t.Fatal("Browse returned nil result")
	}
}

func TestUpdateReturnsNil(t *testing.T) {
	if err := Update(UpdateOptions{}); err != nil {
		t.Errorf("Update returned error: %v", err)
	}
}

func TestInitCreatesManifest(t *testing.T) {
	root := makeTestDir(t)
	err := Init(InitOptions{
		ProjectRoot: root,
		Name:        "my-package",
		Description: "A test package",
		Author:      "test-author",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	manifestPath := filepath.Join(root, "my-package", "marketplace.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse manifest JSON: %v", err)
	}
	if manifest["name"] != "my-package" {
		t.Errorf("unexpected name %v", manifest["name"])
	}
	if manifest["version"] != "0.1.0" {
		t.Errorf("unexpected version %v", manifest["version"])
	}
}

func TestInitCustomOutputDir(t *testing.T) {
	root := makeTestDir(t)
	outDir := filepath.Join(root, "custom-output")
	err := Init(InitOptions{
		ProjectRoot: root,
		Name:        "custom-pkg",
		OutputDir:   outDir,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	manifestPath := filepath.Join(outDir, "marketplace.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Errorf("manifest not created at custom dir: %v", err)
	}
}

func TestInitMissingName(t *testing.T) {
	root := makeTestDir(t)
	err := Init(InitOptions{ProjectRoot: root, Name: ""})
	if err == nil {
		t.Error("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("error should mention name required: %v", err)
	}
}

func TestInitManifestFields(t *testing.T) {
	root := makeTestDir(t)
	err := Init(InitOptions{
		ProjectRoot: root,
		Name:        "field-test",
		Description: "desc",
		Author:      "author",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(root, "field-test", "marketplace.json"))
	var manifest map[string]any
	_ = json.Unmarshal(data, &manifest)
	if manifest["description"] != "desc" {
		t.Errorf("description field missing or wrong: %v", manifest["description"])
	}
	if manifest["author"] != "author" {
		t.Errorf("author field missing or wrong: %v", manifest["author"])
	}
}

func TestCheckValidManifest(t *testing.T) {
	root := makeTestDir(t)
	manifest := map[string]any{"name": "test-pkg", "version": "1.0.0"}
	data, _ := json.MarshalIndent(manifest, "", "  ")
	if err := os.WriteFile(filepath.Join(root, "marketplace.json"), data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	result, err := Check(CheckOptions{ProjectRoot: root})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !result.Valid {
		t.Errorf("expected Valid true, got issues: %v", result.Issues)
	}
}

func TestCheckMissingFields(t *testing.T) {
	root := makeTestDir(t)
	manifest := map[string]any{"name": "pkg-no-version"}
	data, _ := json.MarshalIndent(manifest, "", "  ")
	_ = os.WriteFile(filepath.Join(root, "marketplace.json"), data, 0o644)
	result, err := Check(CheckOptions{ProjectRoot: root})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if result.Valid {
		t.Error("expected invalid manifest")
	}
	if len(result.Issues) == 0 {
		t.Error("expected issues for missing fields")
	}
}

func TestCheckInvalidJSON(t *testing.T) {
	root := makeTestDir(t)
	_ = os.WriteFile(filepath.Join(root, "marketplace.json"), []byte("not json"), 0o644)
	result, err := Check(CheckOptions{ProjectRoot: root})
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if result.Valid {
		t.Error("expected invalid for bad JSON")
	}
}

func TestCheckMissingFile(t *testing.T) {
	root := makeTestDir(t)
	_, err := Check(CheckOptions{ProjectRoot: root})
	if err == nil {
		t.Error("expected error when marketplace.json missing")
	}
}

func TestMigrateNoOp(t *testing.T) {
	if err := Migrate(MigrateOptions{}); err != nil {
		t.Errorf("Migrate returned error: %v", err)
	}
}

func TestOutdatedEmpty(t *testing.T) {
	result, err := Outdated(OutdatedOptions{})
	if err != nil {
		t.Fatalf("Outdated: %v", err)
	}
	if len(result.Packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(result.Packages))
	}
}

func TestDoctorEmpty(t *testing.T) {
	result, err := Doctor(DoctorOptions{})
	if err != nil {
		t.Fatalf("Doctor: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPublishNoOp(t *testing.T) {
	if err := Publish(PublishOptions{}); err != nil {
		t.Errorf("Publish returned error: %v", err)
	}
}

func TestPackageEmpty(t *testing.T) {
	result, err := Package(PackageOptions{})
	if err != nil {
		t.Fatalf("Package: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestSearchEmpty(t *testing.T) {
	result, err := Search(SearchOptions{Query: "test"})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(result.Packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(result.Packages))
	}
}

func TestAddRemoveRoundTrip(t *testing.T) {
	root := makeTestDir(t)
	aliases := []string{"mp1", "mp2", "mp3"}
	for _, a := range aliases {
		_, err := Add(AddOptions{ProjectRoot: root, Alias: a, URL: "https://example.com/" + a})
		if err != nil {
			t.Fatalf("Add %s: %v", a, err)
		}
	}
	if err := Remove(RemoveOptions{ProjectRoot: root, Alias: "mp2"}); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	res, _ := List(ListOptions{ProjectRoot: root})
	if len(res.Entries) != 2 {
		t.Errorf("got %d entries, want 2", len(res.Entries))
	}
	for _, e := range res.Entries {
		if e.Alias == "mp2" {
			t.Error("mp2 should have been removed")
		}
	}
}

func TestListJSONOption(t *testing.T) {
	root := makeTestDir(t)
	writeMarketplaces(t, root, []MarketplaceEntry{
		{Alias: "x", URL: "https://x.com"},
	})
	res, err := List(ListOptions{ProjectRoot: root, JSON: true})
	if err != nil {
		t.Fatalf("List with JSON: %v", err)
	}
	if len(res.Entries) != 1 {
		t.Errorf("got %d entries", len(res.Entries))
	}
}

func TestValidateNoConfigFile(t *testing.T) {
	root := makeTestDir(t)
	_, err := Validate(ValidateOptions{ProjectRoot: root, Alias: "any"})
	if err == nil {
		t.Error("expected error when no config file")
	}
}

func TestPackageSummaryStruct(t *testing.T) {
	ps := PackageSummary{
		Name:        "pkg",
		Version:     "1.0",
		Description: "test",
		Stars:       42,
	}
	if ps.Name != "pkg" || ps.Stars != 42 {
		t.Errorf("unexpected PackageSummary: %+v", ps)
	}
}

func TestOutdatedPackageStruct(t *testing.T) {
	op := OutdatedPackage{
		Name:           "foo",
		CurrentVersion: "1.0",
		LatestVersion:  "1.1",
	}
	if op.Name != "foo" || op.LatestVersion != "1.1" {
		t.Errorf("unexpected OutdatedPackage: %+v", op)
	}
}

func TestAddResultStruct(t *testing.T) {
	r := AddResult{Alias: "a", URL: "u", Branch: "b", Created: true}
	if !r.Created {
		t.Error("expected Created true")
	}
}

func TestListResultStruct(t *testing.T) {
	r := ListResult{
		Entries: []MarketplaceEntry{{Alias: "x", URL: "u"}},
	}
	if len(r.Entries) != 1 {
		t.Error("expected 1 entry")
	}
}

func TestValidateResultStruct(t *testing.T) {
	r := ValidateResult{Alias: "a", Valid: true}
	if !r.Valid {
		t.Error("expected Valid true")
	}
}

func TestCheckResultStruct(t *testing.T) {
	r := CheckResult{Issues: []string{"err"}, Valid: false}
	if r.Valid || len(r.Issues) != 1 {
		t.Error("unexpected CheckResult")
	}
}

func TestDoctorResultStruct(t *testing.T) {
	r := DoctorResult{Issues: []string{"warn"}, Fixed: []string{"fix"}}
	if len(r.Issues) != 1 || len(r.Fixed) != 1 {
		t.Errorf("unexpected DoctorResult: %+v", r)
	}
}

func TestBrowseResultStruct(t *testing.T) {
	r := BrowseResult{Packages: []PackageSummary{{Name: "p", Version: "1.0"}}}
	if len(r.Packages) != 1 {
		t.Error("expected 1 package")
	}
}
