package request_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/request"
)

func TestDefaultInstallRequest(t *testing.T) {
	r := request.DefaultInstallRequest()
	if r.ParallelDownloads != 4 {
		t.Errorf("ParallelDownloads: got %d, want 4", r.ParallelDownloads)
	}
	if r.UpdateRefs {
		t.Error("UpdateRefs should default to false")
	}
	if r.Force {
		t.Error("Force should default to false")
	}
	if r.AllowInsecure {
		t.Error("AllowInsecure should default to false")
	}
	if r.NoPolicy {
		t.Error("NoPolicy should default to false")
	}
}

func TestInstallRequestFields(t *testing.T) {
	r := request.InstallRequest{
		ApmPackagePath:    "/some/path",
		UpdateRefs:        true,
		Verbose:           true,
		OnlyPackages:      []string{"pkg1", "pkg2"},
		Force:             true,
		ParallelDownloads: 8,
		Target:            "vscode",
		AllowInsecure:     true,
		AllowInsecureHosts: []string{"example.com"},
		NoPolicy:          true,
		SkillSubset:       []string{"skill1"},
		SkillSubsetFromCLI: true,
		LegacySkillPaths:  true,
		Frozen:            true,
		ProtocolPref:      "https",
	}
	if r.ApmPackagePath != "/some/path" {
		t.Errorf("ApmPackagePath mismatch")
	}
	if r.ParallelDownloads != 8 {
		t.Errorf("ParallelDownloads: got %d, want 8", r.ParallelDownloads)
	}
	if len(r.OnlyPackages) != 2 {
		t.Errorf("OnlyPackages length: got %d, want 2", len(r.OnlyPackages))
	}
	if r.ProtocolPref != "https" {
		t.Errorf("ProtocolPref: got %q, want %q", r.ProtocolPref, "https")
	}
}
