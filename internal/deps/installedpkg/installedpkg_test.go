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
