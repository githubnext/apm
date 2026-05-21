package conflictdetector_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/conflictdetector"
)

func TestNew_NotNil(t *testing.T) {
	d := conflictdetector.New(nil, nil, nil)
	if d == nil {
		t.Error("expected non-nil detector")
	}
}

func TestGetExistingServerConfigs_Empty(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return nil }, nil, nil)
	got := d.GetExistingServerConfigs()
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}

func TestGetCanonicalServerName_ReturnsString(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{
		"mcp-github": {"type": "stdio"},
	}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	name := d.GetCanonicalServerName("mcp-github")
	if name == "" {
		t.Error("expected non-empty canonical name")
	}
}

func TestFindConflicts_NoConflicts(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{
		"github": {"type": "stdio"},
	}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	conflicts := d.FindConflicts("newserver")
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestCheckServerExists_NotFoundExtra3(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{}
	}, nil, nil)
	r := d.CheckServerExists("nonexistent")
	if r.Exists {
		t.Error("expected server not to exist")
	}
}

func TestServerExistsResult_ZeroValue(t *testing.T) {
	var r conflictdetector.ServerExistsResult
	if r.Exists {
		t.Error("expected zero value Exists=false")
	}
}
