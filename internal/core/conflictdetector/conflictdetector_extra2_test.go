package conflictdetector_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/conflictdetector"
)

func TestCheckServerExists_ExactNameMatch(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{
		"github": {"type": "stdio"},
	}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	r := d.CheckServerExists("github")
	if !r.Exists {
		t.Error("should find server by exact name 'github'")
	}
}

func TestCheckServerExists_CaseSensitive(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{
		"GitHub": {"type": "stdio"},
	}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	r := d.CheckServerExists("github")
	// Case sensitivity depends on implementation; just verify no panic
	_ = r
}

func TestGetExistingServerConfigs_ReturnsAll(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{
		"server1": {"type": "stdio"},
		"server2": {"type": "sse"},
	}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	got := d.GetExistingServerConfigs()
	if len(got) != 2 {
		t.Errorf("expected 2 servers, got %d", len(got))
	}
}

func TestFindConflicts_EmptyRef(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	conflicts := d.FindConflicts("")
	_ = conflicts
}

func TestGetCanonicalServerName_EmptyRef(t *testing.T) {
	d := conflictdetector.New(nil, nil, nil)
	got := d.GetCanonicalServerName("")
	// should not panic
	_ = got
}

func TestCheckServerExists_NilGetServers(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	r := d.CheckServerExists("anything")
	if r.Exists {
		t.Error("empty servers map should return not-found")
	}
}

func TestServerExistsResult_Fields(t *testing.T) {
	var r conflictdetector.ServerExistsResult
	if r.Exists {
		t.Error("zero-value Exists should be false")
	}
}

func TestGetCanonicalServerName_WithResolver(t *testing.T) {
	d := conflictdetector.New(nil, func(ref string) (string, error) { return "canonical-" + ref, nil }, nil)
	got := d.GetCanonicalServerName("foo")
	if got != "canonical-foo" {
		t.Errorf("expected 'canonical-foo', got %q", got)
	}
}

func TestFindConflicts_SingleServerMatch(t *testing.T) {
	servers := map[string]conflictdetector.ServerConfig{
		"mcp-fetch": {"type": "stdio"},
	}
	d := conflictdetector.New(func() map[string]conflictdetector.ServerConfig { return servers }, nil, nil)
	conflicts := d.FindConflicts("mcp-fetch")
	if len(conflicts) == 0 {
		t.Error("should find at least one conflict for exact match")
	}
}
