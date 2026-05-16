package mktvalidator_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mktvalidator"
)

func TestValidatePluginSchema_Valid(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "my-plugin", Source: "owner/repo"},
	}
	r := mktvalidator.ValidatePluginSchema(plugins)
	if !r.Passed {
		t.Fatalf("expected passed, got errors: %v", r.Errors)
	}
	if r.CheckName != "plugin_schema" {
		t.Fatalf("CheckName mismatch: %q", r.CheckName)
	}
}

func TestValidatePluginSchema_EmptyName(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "", Source: "owner/repo"},
	}
	r := mktvalidator.ValidatePluginSchema(plugins)
	if r.Passed {
		t.Fatal("expected failure for empty name")
	}
	if len(r.Errors) == 0 {
		t.Fatal("expected at least one error")
	}
}

func TestValidatePluginSchema_EmptySource(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "my-plugin", Source: ""},
	}
	r := mktvalidator.ValidatePluginSchema(plugins)
	if r.Passed {
		t.Fatal("expected failure for empty source")
	}
}

func TestValidateNoDuplicateNames_NoDups(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "a", Source: "o/a"},
		{Name: "b", Source: "o/b"},
	}
	r := mktvalidator.ValidateNoDuplicateNames(plugins)
	if !r.Passed {
		t.Fatalf("expected passed, got errors: %v", r.Errors)
	}
}

func TestValidateNoDuplicateNames_WithDup(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "dup", Source: "o/a"},
		{Name: "dup", Source: "o/b"},
	}
	r := mktvalidator.ValidateNoDuplicateNames(plugins)
	if r.Passed {
		t.Fatal("expected failure for duplicate name")
	}
}

func TestValidateMarketplace_AllPass(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "p1", Source: "o/p1"},
		{Name: "p2", Source: "o/p2"},
	}
	results := mktvalidator.ValidateMarketplace(plugins)
	for _, r := range results {
		if !r.Passed {
			t.Errorf("check %q failed: %v", r.CheckName, r.Errors)
		}
	}
}
