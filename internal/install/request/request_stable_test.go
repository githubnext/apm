package request_test

import (
"testing"

"github.com/githubnext/apm/internal/install/request"
)

func TestDefaultInstallRequest_ParallelDownloads(t *testing.T) {
r := request.DefaultInstallRequest()
if r.ParallelDownloads != 4 {
t.Errorf("expected default ParallelDownloads=4, got %d", r.ParallelDownloads)
}
}

func TestInstallRequest_UpdateRefs_default(t *testing.T) {
r := request.DefaultInstallRequest()
if r.UpdateRefs {
t.Error("UpdateRefs should default to false")
}
}

func TestInstallRequest_Force_set(t *testing.T) {
r := request.InstallRequest{Force: true}
if !r.Force {
t.Error("expected Force=true")
}
}

func TestInstallRequest_NoPolicy_default(t *testing.T) {
r := request.DefaultInstallRequest()
if r.NoPolicy {
t.Error("NoPolicy should default to false")
}
}

func TestInstallRequest_NoPolicy_set(t *testing.T) {
r := request.InstallRequest{NoPolicy: true}
if !r.NoPolicy {
t.Error("expected NoPolicy=true")
}
}

func TestInstallRequest_LegacySkillPaths(t *testing.T) {
r := request.InstallRequest{LegacySkillPaths: true}
if !r.LegacySkillPaths {
t.Error("expected LegacySkillPaths=true")
}
}

func TestInstallRequest_Frozen_default(t *testing.T) {
r := request.DefaultInstallRequest()
if r.Frozen {
t.Error("Frozen should default to false")
}
}

func TestInstallRequest_Frozen_set(t *testing.T) {
r := request.InstallRequest{Frozen: true}
if !r.Frozen {
t.Error("expected Frozen=true")
}
}

func TestInstallRequest_ProtocolPref_empty(t *testing.T) {
r := request.DefaultInstallRequest()
if r.ProtocolPref != "" {
t.Errorf("expected empty ProtocolPref, got %q", r.ProtocolPref)
}
}

func TestInstallRequest_ProtocolPref_https(t *testing.T) {
r := request.InstallRequest{ProtocolPref: "https"}
if r.ProtocolPref != "https" {
t.Errorf("expected ProtocolPref=https, got %q", r.ProtocolPref)
}
}

func TestInstallRequest_ProtocolPref_ssh(t *testing.T) {
r := request.InstallRequest{ProtocolPref: "ssh"}
if r.ProtocolPref != "ssh" {
t.Errorf("expected ProtocolPref=ssh, got %q", r.ProtocolPref)
}
}

func TestInstallRequest_ApmPackagePath_set(t *testing.T) {
r := request.InstallRequest{ApmPackagePath: "/path/to/apm.yml"}
if r.ApmPackagePath != "/path/to/apm.yml" {
t.Errorf("expected ApmPackagePath=/path/to/apm.yml, got %q", r.ApmPackagePath)
}
}

func TestInstallRequest_OnlyPackages_multiple(t *testing.T) {
r := request.InstallRequest{
OnlyPackages: []string{"pkg1", "pkg2", "pkg3"},
}
if len(r.OnlyPackages) != 3 {
t.Errorf("expected 3 OnlyPackages, got %d", len(r.OnlyPackages))
}
}

func TestInstallRequest_OnlyPackages_empty_default(t *testing.T) {
r := request.DefaultInstallRequest()
if len(r.OnlyPackages) != 0 {
t.Errorf("expected empty OnlyPackages, got %d items", len(r.OnlyPackages))
}
}

func TestInstallRequest_SkillSubset_multiple(t *testing.T) {
r := request.InstallRequest{
SkillSubset:        []string{"core", "extras", "experimental"},
SkillSubsetFromCLI: true,
}
if len(r.SkillSubset) != 3 {
t.Errorf("expected 3 skill subsets, got %d", len(r.SkillSubset))
}
if !r.SkillSubsetFromCLI {
t.Error("expected SkillSubsetFromCLI=true")
}
}

func TestInstallRequest_AllowInsecure_default(t *testing.T) {
r := request.DefaultInstallRequest()
if r.AllowInsecure {
t.Error("AllowInsecure should default to false")
}
}

func TestInstallRequest_AllowInsecure_set(t *testing.T) {
r := request.InstallRequest{AllowInsecure: true}
if !r.AllowInsecure {
t.Error("expected AllowInsecure=true")
}
}

func TestInstallRequest_AllowProtocolFallback_truePtr(t *testing.T) {
b := true
r := request.InstallRequest{AllowProtocolFallback: &b}
if r.AllowProtocolFallback == nil {
t.Fatal("AllowProtocolFallback should not be nil")
}
if !*r.AllowProtocolFallback {
t.Error("expected *AllowProtocolFallback=true")
}
}

func TestInstallRequest_AllowInsecureHosts_single(t *testing.T) {
r := request.InstallRequest{
AllowInsecureHosts: []string{"internal.corp.example.com"},
}
if len(r.AllowInsecureHosts) != 1 {
t.Errorf("expected 1 insecure host, got %d", len(r.AllowInsecureHosts))
}
if r.AllowInsecureHosts[0] != "internal.corp.example.com" {
t.Errorf("unexpected host: %s", r.AllowInsecureHosts[0])
}
}
