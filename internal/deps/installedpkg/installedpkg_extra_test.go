package installedpkg_test

import (
	"testing"

	"github.com/githubnext/apm/internal/deps/installedpkg"
)

func TestInstalledPackage_DepthVariants(t *testing.T) {
	for _, depth := range []int{0, 1, 2, 5, 10} {
		pkg := installedpkg.InstalledPackage{Depth: depth}
		if pkg.Depth != depth {
			t.Errorf("Depth %d mismatch: got %d", depth, pkg.Depth)
		}
	}
}

func TestInstalledPackage_ResolvedByVariants(t *testing.T) {
	cases := []string{"direct", "transitive", "pinned", ""}
	for _, v := range cases {
		pkg := installedpkg.InstalledPackage{ResolvedBy: v}
		if pkg.ResolvedBy != v {
			t.Errorf("ResolvedBy %q mismatch: got %q", v, pkg.ResolvedBy)
		}
	}
}

func TestInstalledPackage_CommitVariants(t *testing.T) {
	cases := []string{
		"abc1234",
		"deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		"v1.2.3",
		"",
	}
	for _, commit := range cases {
		pkg := installedpkg.InstalledPackage{ResolvedCommit: commit}
		if pkg.ResolvedCommit != commit {
			t.Errorf("ResolvedCommit %q mismatch", commit)
		}
	}
}

func TestInstalledPackage_AllFieldsSet(t *testing.T) {
	pkg := installedpkg.InstalledPackage{
		DepRefURL:      "https://github.com/org/repo",
		ResolvedCommit: "abc123",
		Depth:          3,
		ResolvedBy:     "transitive",
		IsDev:          true,
		RegistryHost:   "registry.example.com",
		RegistryPrefix: "my-prefix",
	}
	if pkg.DepRefURL != "https://github.com/org/repo" {
		t.Error("DepRefURL mismatch")
	}
	if pkg.ResolvedCommit != "abc123" {
		t.Error("ResolvedCommit mismatch")
	}
	if pkg.Depth != 3 {
		t.Error("Depth mismatch")
	}
	if pkg.ResolvedBy != "transitive" {
		t.Error("ResolvedBy mismatch")
	}
	if !pkg.IsDev {
		t.Error("IsDev should be true")
	}
	if pkg.RegistryHost != "registry.example.com" {
		t.Error("RegistryHost mismatch")
	}
	if pkg.RegistryPrefix != "my-prefix" {
		t.Error("RegistryPrefix mismatch")
	}
}

func TestInstalledPackage_Slice(t *testing.T) {
	pkgs := []installedpkg.InstalledPackage{
		{DepRefURL: "https://github.com/a/b", Depth: 0},
		{DepRefURL: "https://github.com/c/d", Depth: 1, IsDev: true},
	}
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 pkgs, got %d", len(pkgs))
	}
	if pkgs[1].IsDev != true {
		t.Error("second pkg IsDev should be true")
	}
}

func TestInstalledPackage_IsDevFalseByDefault(t *testing.T) {
	pkg := installedpkg.InstalledPackage{DepRefURL: "https://github.com/x/y"}
	if pkg.IsDev {
		t.Error("IsDev should default to false")
	}
}

func TestInstalledPackage_GHERegistryHost(t *testing.T) {
	pkg := installedpkg.InstalledPackage{
		DepRefURL:    "https://ghe.company.com/org/pkg",
		RegistryHost: "ghe.company.com",
	}
	if pkg.RegistryHost != "ghe.company.com" {
		t.Errorf("RegistryHost mismatch: %q", pkg.RegistryHost)
	}
}

func TestInstalledPackage_EmptyRegistryIsValid(t *testing.T) {
	pkg := installedpkg.InstalledPackage{
		DepRefURL:      "https://github.com/a/b",
		RegistryHost:   "",
		RegistryPrefix: "",
	}
	if pkg.RegistryHost != "" {
		t.Errorf("expected empty RegistryHost, got %q", pkg.RegistryHost)
	}
}
