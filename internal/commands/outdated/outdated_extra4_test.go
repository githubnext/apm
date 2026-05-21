package outdated

import "testing"

func TestOutdatedRow_AllFields_Extra4(t *testing.T) {
	row := OutdatedRow{
		Package:   "mypkg",
		Current:   "v1.0.0",
		Latest:    "v2.0.0",
		Status:    "outdated",
		ExtraTags: []string{"v1.5.0"},
		Source:    "github",
	}
	if row.Package != "mypkg" {
		t.Errorf("Package = %q", row.Package)
	}
	if row.Current != "v1.0.0" {
		t.Errorf("Current = %q", row.Current)
	}
	if row.Latest != "v2.0.0" {
		t.Errorf("Latest = %q", row.Latest)
	}
	if row.Status != "outdated" {
		t.Errorf("Status = %q", row.Status)
	}
	if len(row.ExtraTags) != 1 {
		t.Errorf("ExtraTags len = %d", len(row.ExtraTags))
	}
	if row.Source != "github" {
		t.Errorf("Source = %q", row.Source)
	}
}

func TestOutdatedRow_ZeroValue_Extra4(t *testing.T) {
	var row OutdatedRow
	if row.Package != "" {
		t.Errorf("zero Package = %q", row.Package)
	}
	if row.ExtraTags != nil {
		t.Errorf("zero ExtraTags should be nil")
	}
}

func TestLockEntry_Fields_Extra4(t *testing.T) {
	e := LockEntry{
		Name:            "toolpkg",
		LockedRef:       "v3.1.0",
		LockedCommit:    "deadbeef",
		Source:          "npm",
		MarketplaceName: "toolpkg-npm",
	}
	if e.Name != "toolpkg" {
		t.Errorf("Name = %q", e.Name)
	}
	if e.LockedRef != "v3.1.0" {
		t.Errorf("LockedRef = %q", e.LockedRef)
	}
	if e.LockedCommit != "deadbeef" {
		t.Errorf("LockedCommit = %q", e.LockedCommit)
	}
	if e.MarketplaceName != "toolpkg-npm" {
		t.Errorf("MarketplaceName = %q", e.MarketplaceName)
	}
}

func TestLockEntry_ZeroValue_Extra4(t *testing.T) {
	var e LockEntry
	if e.Name != "" {
		t.Errorf("zero Name = %q", e.Name)
	}
}

func TestLockFile_EmptyEntries_Extra4(t *testing.T) {
	lf := LockFile{Entries: []LockEntry{}}
	if len(lf.Entries) != 0 {
		t.Errorf("len(Entries) = %d", len(lf.Entries))
	}
}

func TestLockFile_MultipleEntries_Extra4(t *testing.T) {
	lf := LockFile{Entries: []LockEntry{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
	}}
	if len(lf.Entries) != 3 {
		t.Errorf("len(Entries) = %d, want 3", len(lf.Entries))
	}
}

func TestCheckOptions_ZeroValue_Extra4(t *testing.T) {
	var opts CheckOptions
	if opts.ProjectRoot != "" {
		t.Errorf("zero ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Verbose {
		t.Error("zero Verbose should be false")
	}
}

func TestCheckOptions_Fields_Extra4(t *testing.T) {
	opts := CheckOptions{
		ProjectRoot: "/repo",
		Verbose:     true,
		Format:      "json",
		NoFetch:     true,
	}
	if opts.ProjectRoot != "/repo" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if !opts.Verbose {
		t.Error("Verbose should be true")
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q", opts.Format)
	}
	if !opts.NoFetch {
		t.Error("NoFetch should be true")
	}
}

func TestCheckResult_ZeroValue_Extra4(t *testing.T) {
	var r CheckResult
	if r.Rows != nil {
		t.Errorf("zero Rows should be nil")
	}
}
