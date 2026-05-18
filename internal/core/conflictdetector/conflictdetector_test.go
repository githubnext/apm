package conflictdetector_test

import (
	"testing"

	"github.com/githubnext/apm/internal/core/conflictdetector"
)

func noServers() map[string]conflictdetector.ServerConfig {
	return map[string]conflictdetector.ServerConfig{}
}

func withServer(name string, cfg conflictdetector.ServerConfig) func() map[string]conflictdetector.ServerConfig {
	return func() map[string]conflictdetector.ServerConfig {
		return map[string]conflictdetector.ServerConfig{name: cfg}
	}
}

func TestCheckServerExists_NotFound(t *testing.T) {
	d := conflictdetector.New(noServers, nil, nil)
	result := d.CheckServerExists("github.com/owner/myserver")
	if result.Exists {
		t.Error("expected no conflict")
	}
}

func TestCheckServerExists_ByCanonicalName(t *testing.T) {
	d := conflictdetector.New(
		withServer("myserver", conflictdetector.ServerConfig{}),
		nil, nil,
	)
	result := d.CheckServerExists("github.com/owner/myserver")
	if !result.Exists {
		t.Error("expected conflict by canonical name")
	}
	if result.ConflictName != "myserver" {
		t.Errorf("expected conflict name 'myserver', got %q", result.ConflictName)
	}
}

func TestCheckServerExists_ByUUID(t *testing.T) {
	existing := withServer("some-server", conflictdetector.ServerConfig{"id": "uuid-123"})
	findFn := func(ref string) (map[string]interface{}, error) {
		return map[string]interface{}{"id": "uuid-123"}, nil
	}
	d := conflictdetector.New(existing, nil, findFn)
	result := d.CheckServerExists("any/ref")
	if !result.Exists {
		t.Error("expected UUID-based conflict detection")
	}
	if result.ConflictUUID != "uuid-123" {
		t.Errorf("expected UUID 'uuid-123', got %q", result.ConflictUUID)
	}
}

func TestGetCanonicalServerName_FallbackLastComponent(t *testing.T) {
	d := conflictdetector.New(noServers, nil, nil)
	name := d.GetCanonicalServerName("github.com/owner/myserver")
	if name != "myserver" {
		t.Errorf("expected 'myserver', got %q", name)
	}
}

func TestGetCanonicalServerName_CustomResolver(t *testing.T) {
	d := conflictdetector.New(noServers, func(ref string) (string, error) {
		return "custom-name", nil
	}, nil)
	name := d.GetCanonicalServerName("anything")
	if name != "custom-name" {
		t.Errorf("expected 'custom-name', got %q", name)
	}
}

func TestFindConflicts_None(t *testing.T) {
	d := conflictdetector.New(noServers, nil, nil)
	conflicts := d.FindConflicts("owner/newserver")
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", conflicts)
	}
}

func TestFindConflicts_Found(t *testing.T) {
	d := conflictdetector.New(
		withServer("newserver", conflictdetector.ServerConfig{}),
		nil, nil,
	)
	conflicts := d.FindConflicts("owner/newserver")
	if len(conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(conflicts))
	}
}
