package installedpkg_test

import (
	"testing"

	"github.com/githubnext/apm/internal/deps/installedpkg"
)

func TestInstalledPackage_Fields(t *testing.T) {
	pkg := installedpkg.InstalledPackage{
		DepRefURL:      "https://github.com/owner/repo",
		ResolvedCommit: "abc1234",
		Depth:          1,
		ResolvedBy:     "direct",
		IsDev:          false,
		RegistryHost:   "",
		RegistryPrefix: "",
	}
	if pkg.DepRefURL != "https://github.com/owner/repo" {
		t.Fatalf("DepRefURL mismatch")
	}
	if pkg.ResolvedCommit != "abc1234" {
		t.Fatalf("ResolvedCommit mismatch")
	}
	if pkg.Depth != 1 {
		t.Fatalf("Depth mismatch")
	}
	if pkg.IsDev {
		t.Fatal("IsDev should be false")
	}
}

func TestInstalledPackage_DevPackage(t *testing.T) {
	pkg := installedpkg.InstalledPackage{
		DepRefURL: "https://github.com/owner/dev-tool",
		IsDev:     true,
	}
	if !pkg.IsDev {
		t.Fatal("IsDev should be true")
	}
}

func TestInstalledPackage_RegistryFields(t *testing.T) {
	pkg := installedpkg.InstalledPackage{
		RegistryHost:   "artifactory.example.com",
		RegistryPrefix: "npm-proxy",
	}
	if pkg.RegistryHost != "artifactory.example.com" {
		t.Fatalf("RegistryHost mismatch")
	}
	if pkg.RegistryPrefix != "npm-proxy" {
		t.Fatalf("RegistryPrefix mismatch")
	}
}

func TestInstalledPackage_ZeroValue(t *testing.T) {
	var pkg installedpkg.InstalledPackage
	if pkg.Depth != 0 {
		t.Errorf("zero Depth should be 0, got %d", pkg.Depth)
	}
	if pkg.IsDev {
		t.Error("zero IsDev should be false")
	}
	if pkg.DepRefURL != "" {
		t.Errorf("zero DepRefURL should be empty, got %q", pkg.DepRefURL)
	}
}

func TestInstalledPackage_DepthLevels(t *testing.T) {
	for _, depth := range []int{0, 1, 2, 5, 10} {
		pkg := installedpkg.InstalledPackage{Depth: depth}
		if pkg.Depth != depth {
			t.Errorf("Depth: got %d, want %d", pkg.Depth, depth)
		}
	}
}

func TestInstalledPackage_ResolvedBy(t *testing.T) {
	cases := []string{"direct", "transitive", "dev", "peer"}
	for _, by := range cases {
		pkg := installedpkg.InstalledPackage{ResolvedBy: by}
		if pkg.ResolvedBy != by {
			t.Errorf("ResolvedBy: got %q, want %q", pkg.ResolvedBy, by)
		}
	}
}

func TestInstalledPackage_CommitFormats(t *testing.T) {
	commits := []string{"abc1234", "deadbeef12345678", "0000000"}
	for _, c := range commits {
		pkg := installedpkg.InstalledPackage{ResolvedCommit: c}
		if pkg.ResolvedCommit != c {
			t.Errorf("ResolvedCommit: got %q, want %q", pkg.ResolvedCommit, c)
		}
	}
}
