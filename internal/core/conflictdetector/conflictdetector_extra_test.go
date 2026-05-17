package conflictdetector_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/conflictdetector"
)

func TestGetExistingServerConfigs_Nil(t *testing.T) {
	d := conflictdetector.New(nil, nil, nil)
	cfg := d.GetExistingServerConfigs()
	if cfg == nil {
		t.Error("expected non-nil empty map for nil GetExistingServers")
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty map, got %v", cfg)
	}
}

func TestGetExistingServerConfigs_NonEmpty(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{
			"server-a": {"type": "command"},
			"server-b": {"type": "url"},
		}
	}, nil, nil)
	cfg := d.GetExistingServerConfigs()
	if len(cfg) != 2 {
		t.Errorf("expected 2 configs, got %d", len(cfg))
	}
}

func TestCheckServerExists_MultipleServersNoMatch(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{
			"alpha": {},
			"beta":  {},
		}
	}, nil, nil)
	res := d.CheckServerExists("gamma")
	if res.Exists {
		t.Error("expected no conflict for 'gamma'")
	}
}

func TestCheckServerExists_EmptyRef(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{"": {}}
	}, nil, nil)
	// Empty ref should match empty key (edge case).
	res := d.CheckServerExists("")
	_ = res // just verifying no panic
}

func TestGetCanonicalServerName_SlashSeparated(t *testing.T) {
	d := conflictdetector.New(nil, nil, nil)
	name := d.GetCanonicalServerName("org/repo/tool")
	if name != "tool" {
		t.Errorf("expected 'tool', got %q", name)
	}
}

func TestGetCanonicalServerName_NilResolver(t *testing.T) {
	d := conflictdetector.New(nil, nil, nil)
	name := d.GetCanonicalServerName("just-name")
	if name != "just-name" {
		t.Errorf("expected 'just-name', got %q", name)
	}
}

func TestFindConflicts_Multiple(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{
			"server1": {},
			"server2": {},
		}
	}, nil, nil)
	// Neither server1 nor server2 matches "new-server" by canonical name.
	conflicts := d.FindConflicts("owner/new-server")
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", conflicts)
	}
}

func TestFindConflicts_ExactMatch(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{
			"my-server": {},
		}
	}, nil, nil)
	conflicts := d.FindConflicts("owner/my-server")
	if len(conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d: %v", len(conflicts), conflicts)
	}
	if conflicts[0] != "my-server" {
		t.Errorf("expected 'my-server', got %q", conflicts[0])
	}
}

func TestCheckServerExists_CustomNameResolver(t *testing.T) {
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{
			"resolved-name": {},
		}
	}, func(ref string) (string, error) {
		return "resolved-name", nil
	}, nil)
	res := d.CheckServerExists("any/ref")
	if !res.Exists {
		t.Error("expected conflict with custom name resolver")
	}
	if res.ConflictName != "resolved-name" {
		t.Errorf("expected 'resolved-name', got %q", res.ConflictName)
	}
}
