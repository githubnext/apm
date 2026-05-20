package mktvalidator_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/mktvalidator"
)

func TestValidatePluginSchema_AllValid(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "alpha", Source: "https://example.com/alpha"},
		{Name: "beta", Source: "https://example.com/beta"},
	}
	r := mktvalidator.ValidatePluginSchema(plugins)
	if !r.Passed {
		t.Errorf("expected pass, got errors: %v", r.Errors)
	}
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors, got %v", r.Errors)
	}
}

func TestValidatePluginSchema_OnlySourceMissing(t *testing.T) {
	plugins := []mktvalidator.Plugin{{Name: "p", Source: ""}}
	r := mktvalidator.ValidatePluginSchema(plugins)
	if r.Passed {
		t.Error("expected fail")
	}
	if len(r.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(r.Errors))
	}
}

func TestValidatePluginSchema_CheckName(t *testing.T) {
	r := mktvalidator.ValidatePluginSchema(nil)
	if r.CheckName != "plugin_schema" {
		t.Errorf("unexpected check name: %q", r.CheckName)
	}
}

func TestValidateNoDuplicateNames_CheckName(t *testing.T) {
	r := mktvalidator.ValidateNoDuplicateNames(nil)
	if r.CheckName != "no_duplicate_names" {
		t.Errorf("unexpected check name: %q", r.CheckName)
	}
}

func TestValidateNoDuplicateNames_ThreeDups(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "x", Source: "a"},
		{Name: "x", Source: "b"},
		{Name: "x", Source: "c"},
	}
	r := mktvalidator.ValidateNoDuplicateNames(plugins)
	if r.Passed {
		t.Error("expected fail due to duplicates")
	}
}

func TestValidateMarketplace_ResultCount(t *testing.T) {
	results := mktvalidator.ValidateMarketplace(nil)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestValidateMarketplace_NilPluginsPass(t *testing.T) {
	results := mktvalidator.ValidateMarketplace(nil)
	for _, r := range results {
		if !r.Passed {
			t.Errorf("check %q failed for nil plugins", r.CheckName)
		}
	}
}

func TestValidationResult_CheckNamePreserved(t *testing.T) {
	r := mktvalidator.ValidatePluginSchema([]mktvalidator.Plugin{{Name: "n", Source: "s"}})
	if !strings.Contains(r.CheckName, "schema") {
		t.Errorf("check name should contain 'schema', got %q", r.CheckName)
	}
}

func TestPlugin_ZeroValue(t *testing.T) {
	var p mktvalidator.Plugin
	if p.Name != "" || p.Source != "" {
		t.Error("zero value should have empty fields")
	}
}

func TestValidateNoDuplicateNames_ErrorMentionsName(t *testing.T) {
	plugins := []mktvalidator.Plugin{
		{Name: "dup", Source: "a"},
		{Name: "dup", Source: "b"},
	}
	r := mktvalidator.ValidateNoDuplicateNames(plugins)
	if len(r.Errors) == 0 {
		t.Fatal("expected error")
	}
	if !strings.Contains(r.Errors[0], "dup") {
		t.Errorf("error should mention 'dup': %q", r.Errors[0])
	}
}
