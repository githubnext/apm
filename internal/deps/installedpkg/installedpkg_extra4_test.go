package installedpkg

import "testing"

func TestInstalledPackage_DepRefURLField_Extra4(t *testing.T) {
p := InstalledPackage{DepRefURL: "https://github.com/org/repo"}
if p.DepRefURL != "https://github.com/org/repo" {
t.Errorf("unexpected DepRefURL: %s", p.DepRefURL)
}
}

func TestInstalledPackage_RegistryHostField_Extra4(t *testing.T) {
p := InstalledPackage{RegistryHost: "registry.example.com"}
if p.RegistryHost != "registry.example.com" {
t.Errorf("unexpected RegistryHost: %s", p.RegistryHost)
}
}

func TestInstalledPackage_RegistryPrefixField_Extra4(t *testing.T) {
p := InstalledPackage{RegistryPrefix: "myorg/"}
if p.RegistryPrefix != "myorg/" {
t.Errorf("unexpected RegistryPrefix: %s", p.RegistryPrefix)
}
}

func TestInstalledPackage_ResolvedCommit_Extra4(t *testing.T) {
p := InstalledPackage{ResolvedCommit: "abc1234567890abcdef"}
if p.ResolvedCommit != "abc1234567890abcdef" {
t.Errorf("unexpected ResolvedCommit: %s", p.ResolvedCommit)
}
}

func TestInstalledPackage_DepthOne_Extra4(t *testing.T) {
p := InstalledPackage{Depth: 1}
if p.Depth != 1 {
t.Errorf("unexpected Depth: %d", p.Depth)
}
}

func TestInstalledPackage_DepthFive_Extra4(t *testing.T) {
p := InstalledPackage{Depth: 5}
if p.Depth != 5 {
t.Errorf("unexpected Depth: %d", p.Depth)
}
}

func TestInstalledPackage_IsDevTrue_Extra4(t *testing.T) {
p := InstalledPackage{IsDev: true}
if !p.IsDev {
t.Error("expected IsDev to be true")
}
}

func TestInstalledPackage_RegistryHostEmpty_Extra4(t *testing.T) {
p := InstalledPackage{}
if p.RegistryHost != "" {
t.Errorf("expected empty registry host, got %s", p.RegistryHost)
}
}

func TestInstalledPackage_ResolvedByAPI_Extra4(t *testing.T) {
p := InstalledPackage{ResolvedBy: "api"}
if p.ResolvedBy != "api" {
t.Errorf("unexpected ResolvedBy: %s", p.ResolvedBy)
}
}

func TestInstalledPackage_ResolvedByCached_Extra4(t *testing.T) {
p := InstalledPackage{ResolvedBy: "cached"}
if p.ResolvedBy != "cached" {
t.Errorf("unexpected ResolvedBy: %s", p.ResolvedBy)
}
}

func TestInstalledPackage_AllFieldsEmpty_Extra4(t *testing.T) {
p := InstalledPackage{}
if p.DepRefURL != "" || p.ResolvedCommit != "" || p.Depth != 0 || p.IsDev {
t.Error("expected all fields at zero value")
}
}
