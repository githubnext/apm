package installedpkg

import "testing"

func TestInstalledPackage_AllFieldsSet_Extra3(t *testing.T) {
	p := InstalledPackage{
		DepRefURL:      "owner/repo",
		ResolvedCommit: "abc1234",
		Depth:          2,
		ResolvedBy:     "github",
		IsDev:          true,
		RegistryHost:   "ghcr.io",
		RegistryPrefix: "prefix/",
	}
	if p.DepRefURL != "owner/repo" {
		t.Errorf("DepRefURL = %q, want owner/repo", p.DepRefURL)
	}
	if p.Depth != 2 {
		t.Errorf("Depth = %d, want 2", p.Depth)
	}
	if !p.IsDev {
		t.Error("IsDev should be true")
	}
	if p.RegistryHost != "ghcr.io" {
		t.Errorf("RegistryHost = %q, want ghcr.io", p.RegistryHost)
	}
}

func TestInstalledPackage_DepthZero_Extra3(t *testing.T) {
	p := InstalledPackage{Depth: 0}
	if p.Depth != 0 {
		t.Errorf("Depth = %d, want 0", p.Depth)
	}
}

func TestInstalledPackage_EmptyRegistryPrefix_Extra3(t *testing.T) {
	p := InstalledPackage{RegistryHost: "ghcr.io", RegistryPrefix: ""}
	if p.RegistryPrefix != "" {
		t.Error("RegistryPrefix should be empty")
	}
}

func TestInstalledPackage_NotDev_Extra3(t *testing.T) {
	p := InstalledPackage{IsDev: false}
	if p.IsDev {
		t.Error("IsDev should be false")
	}
}

func TestInstalledPackage_LargeDepth_Extra3(t *testing.T) {
	p := InstalledPackage{Depth: 100}
	if p.Depth != 100 {
		t.Errorf("Depth = %d, want 100", p.Depth)
	}
}

func TestInstalledPackage_ResolvedByEmpty_Extra3(t *testing.T) {
	p := InstalledPackage{}
	if p.ResolvedBy != "" {
		t.Errorf("ResolvedBy = %q, want empty", p.ResolvedBy)
	}
}

func TestInstalledPackage_RegistryPrefixSet_Extra3(t *testing.T) {
	p := InstalledPackage{RegistryPrefix: "my/prefix/"}
	if p.RegistryPrefix != "my/prefix/" {
		t.Errorf("RegistryPrefix = %q", p.RegistryPrefix)
	}
}

func TestInstalledPackage_CommitShort_Extra3(t *testing.T) {
	p := InstalledPackage{ResolvedCommit: "abc"}
	if p.ResolvedCommit != "abc" {
		t.Errorf("ResolvedCommit = %q, want abc", p.ResolvedCommit)
	}
}

func TestInstalledPackage_MultiplePackages_Extra3(t *testing.T) {
	pkgs := []InstalledPackage{
		{DepRefURL: "owner/a", Depth: 0},
		{DepRefURL: "owner/b", Depth: 1},
		{DepRefURL: "owner/c", Depth: 2, IsDev: true},
	}
	if len(pkgs) != 3 {
		t.Fatalf("expected 3 packages")
	}
	if pkgs[2].IsDev != true {
		t.Error("pkgs[2] should be dev")
	}
}

func TestInstalledPackage_DepRefURLVariants_Extra3(t *testing.T) {
	variants := []string{
		"owner/repo",
		"owner/repo/subpath",
		"https://github.com/owner/repo",
	}
	for _, url := range variants {
		p := InstalledPackage{DepRefURL: url}
		if p.DepRefURL != url {
			t.Errorf("DepRefURL = %q, want %q", p.DepRefURL, url)
		}
	}
}
