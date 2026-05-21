package mktvalidator_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mktvalidator"
)

func TestValidatePluginSchema_Empty(t *testing.T) {
	r := mktvalidator.ValidatePluginSchema(nil)
	if !r.Passed {
		t.Error("nil list should pass schema check")
	}
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors, got %v", r.Errors)
	}
}

func TestValidatePluginSchema_BothEmpty(t *testing.T) {
	plugins := []mktvalidator.Plugin{{Name: "", Source: ""}}
	r := mktvalidator.ValidatePluginSchema(plugins)
	if r.Passed {
		t.Error("both empty fields should fail")
	}
	if len(r.Errors) < 2 {
		t.Errorf("expected at least 2 errors, got %d: %v", len(r.Errors), r.Errors)
	}
}

func TestValidatePluginSchema_ManyErrors(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "", Source: ""},
		{Name: "ok", Source: ""},
		{Name: "", Source: "owner/repo"},
	}
	r := mktvalidator.ValidatePluginSchema(plugins)
	if r.Passed {
		t.Error("expected failure")
	}
	if len(r.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d", len(r.Errors))
	}
}

func TestValidateNoDuplicateNames_EmptyList(t *testing.T) {
	r := mktvalidator.ValidateNoDuplicateNames(nil)
	if !r.Passed {
		t.Error("nil list should pass dedup check")
	}
}

func TestValidateNoDuplicateNames_Single(t *testing.T) {
	r := mktvalidator.ValidateNoDuplicateNames([]mktvalidator.Plugin{{Name: "only", Source: "o/r"}})
	if !r.Passed {
		t.Error("single plugin should pass dedup check")
	}
}

func TestValidateNoDuplicateNames_ManyDups(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "a", Source: "o/a"},
		{Name: "b", Source: "o/b"},
		{Name: "a", Source: "o/a2"},
		{Name: "b", Source: "o/b2"},
	}
	r := mktvalidator.ValidateNoDuplicateNames(plugins)
	if r.Passed {
		t.Error("expected failure for duplicate names")
	}
	if len(r.Errors) < 2 {
		t.Errorf("expected >= 2 errors, got %d", len(r.Errors))
	}
}

func TestValidateMarketplace_BothPass(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "plugin-a", Source: "owner/a"},
		{Name: "plugin-b", Source: "owner/b"},
	}
	results := mktvalidator.ValidateMarketplace(plugins)
	if len(results) == 0 {
		t.Fatal("expected results from ValidateMarketplace")
	}
	for _, r := range results {
		if !r.Passed {
			t.Errorf("check %q failed unexpectedly: %v", r.CheckName, r.Errors)
		}
	}
}

func TestValidateMarketplace_DuplicatesDetected(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "dup", Source: "o/a"},
		{Name: "dup", Source: "o/b"},
	}
	results := mktvalidator.ValidateMarketplace(plugins)
	found := false
	for _, r := range results {
		if r.CheckName == "no_duplicate_names" && !r.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected dedup check to fail")
	}
}

func TestValidateMarketplace_SchemaFailure(t *testing.T) {
	plugins := []mktvalidator.Plugin{{Name: "", Source: "o/r"}}
	results := mktvalidator.ValidateMarketplace(plugins)
	found := false
	for _, r := range results {
		if r.CheckName == "plugin_schema" && !r.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected schema check to fail for empty name")
	}
}

func TestValidationResult_Warnings(t *testing.T) {
	r := mktvalidator.ValidationResult{
		CheckName: "custom",
		Passed:    true,
		Warnings:  []string{"warn1", "warn2"},
	}
	if len(r.Warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(r.Warnings))
	}
}

func TestPlugin_Fields(t *testing.T) {
	p := mktvalidator.Plugin{Name: "myplugin", Source: "org/myplugin"}
	if p.Name != "myplugin" {
		t.Errorf("Name = %q", p.Name)
	}
	if p.Source != "org/myplugin" {
		t.Errorf("Source = %q", p.Source)
	}
}
