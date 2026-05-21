package conflictdetector

import (
	"fmt"
	"testing"
)

func TestCheckServerExists_NotFound_Extra4(t *testing.T) {
	d := New(
		func() map[string]ServerConfig { return map[string]ServerConfig{} },
		nil,
		nil,
	)
	result := d.CheckServerExists("owner/pkg")
	if result.Exists {
		t.Error("expected not found for empty servers")
	}
}

func TestCheckServerExists_FoundByCanonical_Extra4(t *testing.T) {
	d := New(
		func() map[string]ServerConfig {
			return map[string]ServerConfig{"owner/pkg": {}}
		},
		nil,
		nil,
	)
	result := d.CheckServerExists("owner/pkg")
	if !result.Exists {
		t.Error("expected found for exact canonical name")
	}
	if result.ConflictName != "owner/pkg" {
		t.Errorf("expected owner/pkg, got %s", result.ConflictName)
	}
}

func TestCheckServerExists_WithResolver_Extra4(t *testing.T) {
	d := New(
		func() map[string]ServerConfig {
			return map[string]ServerConfig{"pkg": {}}
		},
		func(ref string) (string, error) { return ref, nil },
		nil,
	)
	result := d.CheckServerExists("pkg")
	if !result.Exists {
		t.Error("expected found via canonical resolver")
	}
}

func TestCheckServerExists_ResolverError_Extra4(t *testing.T) {
	d := New(
		func() map[string]ServerConfig { return map[string]ServerConfig{} },
		func(ref string) (string, error) { return "", fmt.Errorf("err") },
		nil,
	)
	result := d.CheckServerExists("anything")
	if result.Exists {
		t.Error("expected not found when resolver errors and servers empty")
	}
}

func TestCheckServerExists_UUIDMatch_Extra4(t *testing.T) {
	d := New(
		func() map[string]ServerConfig {
			return map[string]ServerConfig{
				"myserver": {"id": "uuid-123"},
			}
		},
		nil,
		func(ref string) (map[string]interface{}, error) {
			return map[string]interface{}{"id": "uuid-123"}, nil
		},
	)
	result := d.CheckServerExists("owner/pkg")
	if !result.Exists {
		t.Error("expected found by UUID")
	}
	if result.ConflictUUID != "uuid-123" {
		t.Errorf("expected uuid-123, got %s", result.ConflictUUID)
	}
}

func TestCheckServerExists_UUIDNoMatch_Extra4(t *testing.T) {
	d := New(
		func() map[string]ServerConfig {
			return map[string]ServerConfig{
				"myserver": {"id": "uuid-xyz"},
			}
		},
		nil,
		func(ref string) (map[string]interface{}, error) {
			return map[string]interface{}{"id": "uuid-abc"}, nil
		},
	)
	result := d.CheckServerExists("owner/pkg")
	if result.Exists {
		t.Error("expected not found when UUIDs don't match and canonical doesn't exist")
	}
}

func TestNew_NotNil_Extra4(t *testing.T) {
	d := New(nil, nil, nil)
	if d == nil {
		t.Error("expected non-nil detector")
	}
}

func TestCheckServerExists_AliasMatch_Extra4(t *testing.T) {
	d := New(
		func() map[string]ServerConfig {
			return map[string]ServerConfig{"alias": {}}
		},
		func(ref string) (string, error) { return "alias", nil },
		nil,
	)
	result := d.CheckServerExists("owner/alias")
	if !result.Exists {
		t.Error("expected alias match via resolver")
	}
}
