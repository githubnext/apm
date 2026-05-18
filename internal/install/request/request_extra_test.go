package request_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/request"
)

func TestInstallRequest_AllowProtocolFallback_nil(t *testing.T) {
	r := request.DefaultInstallRequest()
	if r.AllowProtocolFallback != nil {
		t.Error("AllowProtocolFallback should be nil by default")
	}
}

func TestInstallRequest_AllowProtocolFallback_set(t *testing.T) {
	b := true
	r := request.InstallRequest{AllowProtocolFallback: &b}
	if r.AllowProtocolFallback == nil || !*r.AllowProtocolFallback {
		t.Error("AllowProtocolFallback should be true")
	}
}

func TestInstallRequest_AllowProtocolFallback_false(t *testing.T) {
	b := false
	r := request.InstallRequest{AllowProtocolFallback: &b}
	if r.AllowProtocolFallback == nil || *r.AllowProtocolFallback {
		t.Error("AllowProtocolFallback should be false")
	}
}

func TestInstallRequest_SkillSubset_single(t *testing.T) {
	r := request.InstallRequest{
		SkillSubset:       []string{"core"},
		SkillSubsetFromCLI: true,
	}
	if !r.SkillSubsetFromCLI {
		t.Error("expected SkillSubsetFromCLI=true")
	}
	if r.SkillSubset[0] != "core" {
		t.Errorf("expected 'core', got %s", r.SkillSubset[0])
	}
}

func TestInstallRequest_Verbose(t *testing.T) {
	r := request.InstallRequest{Verbose: true}
	if !r.Verbose {
		t.Error("expected Verbose=true")
	}
}

func TestInstallRequest_Target(t *testing.T) {
	r := request.InstallRequest{Target: "claude"}
	if r.Target != "claude" {
		t.Errorf("expected Target=claude, got %s", r.Target)
	}
}

func TestInstallRequest_EmptyTarget(t *testing.T) {
	r := request.DefaultInstallRequest()
	if r.Target != "" {
		t.Errorf("default Target should be empty, got %s", r.Target)
	}
}

func TestInstallRequest_AllowInsecureHosts_multiple(t *testing.T) {
	r := request.InstallRequest{
		AllowInsecure:      true,
		AllowInsecureHosts: []string{"host1.example.com", "host2.example.com", "192.168.1.1"},
	}
	if len(r.AllowInsecureHosts) != 3 {
		t.Errorf("expected 3 insecure hosts, got %d", len(r.AllowInsecureHosts))
	}
}
