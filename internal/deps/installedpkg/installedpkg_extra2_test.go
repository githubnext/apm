package installedpkg

import (
	"testing"
)

func TestInstalledPackage_ZeroValue(t *testing.T) {
	var pkg InstalledPackage
	if pkg.DepRefURL != "" {
		t.Error("zero-value DepRefURL should be empty")
	}
	if pkg.Depth != 0 {
		t.Error("zero-value Depth should be 0")
	}
	if pkg.IsDev {
		t.Error("zero-value IsDev should be false")
	}
}

func TestInstalledPackage_DepRefURL_Set(t *testing.T) {
	pkg := InstalledPackage{DepRefURL: "https://github.com/owner/repo"}
	if pkg.DepRefURL != "https://github.com/owner/repo" {
		t.Errorf("unexpected DepRefURL: %q", pkg.DepRefURL)
	}
}

func TestInstalledPackage_ResolvedCommit_Set(t *testing.T) {
	pkg := InstalledPackage{ResolvedCommit: "abc123def456"}
	if pkg.ResolvedCommit != "abc123def456" {
		t.Errorf("unexpected ResolvedCommit: %q", pkg.ResolvedCommit)
	}
}

func TestInstalledPackage_DevPackage(t *testing.T) {
	pkg := InstalledPackage{IsDev: true, DepRefURL: "https://github.com/dev/pkg"}
	if !pkg.IsDev {
		t.Error("expected IsDev=true")
	}
}

func TestInstalledPackage_RegistryFields(t *testing.T) {
	pkg := InstalledPackage{
		RegistryHost:   "github.example.com",
		RegistryPrefix: "myorg",
	}
	if pkg.RegistryHost != "github.example.com" {
		t.Errorf("unexpected RegistryHost: %q", pkg.RegistryHost)
	}
	if pkg.RegistryPrefix != "myorg" {
		t.Errorf("unexpected RegistryPrefix: %q", pkg.RegistryPrefix)
	}
}

func TestInstalledPackage_ResolvedByField(t *testing.T) {
	pkg := InstalledPackage{ResolvedBy: "lockfile"}
	if pkg.ResolvedBy != "lockfile" {
		t.Errorf("unexpected ResolvedBy: %q", pkg.ResolvedBy)
	}
}

func TestInstalledPackage_DepthValues(t *testing.T) {
	for _, d := range []int{0, 1, 5, 100} {
		pkg := InstalledPackage{Depth: d}
		if pkg.Depth != d {
			t.Errorf("expected Depth=%d, got %d", d, pkg.Depth)
		}
	}
}

func TestInstalledPackage_EmptyRegistryFields(t *testing.T) {
	pkg := InstalledPackage{DepRefURL: "https://github.com/a/b"}
	if pkg.RegistryHost != "" || pkg.RegistryPrefix != "" {
		t.Error("expected empty registry fields by default")
	}
}

func TestInstalledPackage_FullyPopulated(t *testing.T) {
	pkg := InstalledPackage{
		DepRefURL:      "https://github.com/org/repo",
		ResolvedCommit: "deadbeef1234",
		Depth:          3,
		ResolvedBy:     "direct",
		IsDev:          false,
		RegistryHost:   "github.com",
		RegistryPrefix: "org",
	}
	if pkg.DepRefURL == "" {
		t.Error("DepRefURL should not be empty")
	}
	if pkg.ResolvedCommit == "" {
		t.Error("ResolvedCommit should not be empty")
	}
	if pkg.Depth != 3 {
		t.Errorf("expected Depth=3, got %d", pkg.Depth)
	}
}

func TestInstalledPackage_MapKey(t *testing.T) {
	m := map[string]InstalledPackage{}
	m["pkg-a"] = InstalledPackage{DepRefURL: "https://github.com/a/a"}
	m["pkg-b"] = InstalledPackage{DepRefURL: "https://github.com/b/b", IsDev: true}
	if len(m) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m))
	}
	if !m["pkg-b"].IsDev {
		t.Error("pkg-b should be dev")
	}
}
