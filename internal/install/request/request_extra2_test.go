package request

import (
	"testing"
)

func TestInstallRequest_ZeroValue(t *testing.T) {
	var r InstallRequest
	if r.Force {
		t.Error("Force should default to false")
	}
	if r.Frozen {
		t.Error("Frozen should default to false")
	}
	if r.Verbose {
		t.Error("Verbose should default to false")
	}
	if r.AllowInsecure {
		t.Error("AllowInsecure should default to false")
	}
}

func TestDefaultInstallRequest_Parallel(t *testing.T) {
	r := DefaultInstallRequest()
	if r.ParallelDownloads != 4 {
		t.Errorf("expected ParallelDownloads=4, got %d", r.ParallelDownloads)
	}
}

func TestDefaultInstallRequest_NotFrozen(t *testing.T) {
	r := DefaultInstallRequest()
	if r.Frozen {
		t.Error("default should not be frozen")
	}
}

func TestDefaultInstallRequest_NoPolicy(t *testing.T) {
	r := DefaultInstallRequest()
	if r.NoPolicy {
		t.Error("default NoPolicy should be false")
	}
}

func TestInstallRequest_SkillSubset(t *testing.T) {
	r := InstallRequest{SkillSubset: []string{"skill-a", "skill-b"}}
	if len(r.SkillSubset) != 2 {
		t.Errorf("expected 2 skills, got %d", len(r.SkillSubset))
	}
}

func TestInstallRequest_ProtocolPref(t *testing.T) {
	r := InstallRequest{ProtocolPref: "ssh"}
	if r.ProtocolPref != "ssh" {
		t.Errorf("ProtocolPref = %q, want ssh", r.ProtocolPref)
	}
}

func TestInstallRequest_OnlyPackages(t *testing.T) {
	r := InstallRequest{OnlyPackages: []string{"owner/repo"}}
	if len(r.OnlyPackages) != 1 {
		t.Errorf("expected 1 package, got %d", len(r.OnlyPackages))
	}
}

func TestInstallRequest_UpdateRefs(t *testing.T) {
	r := InstallRequest{UpdateRefs: true}
	if !r.UpdateRefs {
		t.Error("UpdateRefs should be true")
	}
}
