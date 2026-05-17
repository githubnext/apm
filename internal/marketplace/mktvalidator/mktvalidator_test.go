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

func TestValidatePluginSchema_MultipleErrors(t *testing.T) {
plugins := []mktvalidator.Plugin{
{Name: "", Source: ""},
{Name: "ok", Source: ""},
}
r := mktvalidator.ValidatePluginSchema(plugins)
if r.Passed {
t.Fatal("expected failure")
}
if len(r.Errors) < 2 {
t.Errorf("expected at least 2 errors, got %d", len(r.Errors))
}
}

func TestValidateNoDuplicateNames_MultipleDups(t *testing.T) {
plugins := []mktvalidator.Plugin{
{Name: "x", Source: "o/1"},
{Name: "x", Source: "o/2"},
{Name: "x", Source: "o/3"},
}
r := mktvalidator.ValidateNoDuplicateNames(plugins)
if r.Passed {
t.Fatal("expected failure for duplicate name")
}
}

func TestValidateNoDuplicateNames_Empty(t *testing.T) {
r := mktvalidator.ValidateNoDuplicateNames(nil)
if !r.Passed {
t.Fatal("empty list should pass")
}
}

func TestValidateMarketplace_CheckCount(t *testing.T) {
results := mktvalidator.ValidateMarketplace(nil)
if len(results) < 2 {
t.Errorf("expected at least 2 checks, got %d", len(results))
}
}

func TestValidateMarketplace_CheckNames(t *testing.T) {
results := mktvalidator.ValidateMarketplace([]mktvalidator.Plugin{{Name: "p", Source: "o/p"}})
names := map[string]bool{}
for _, r := range results {
names[r.CheckName] = true
}
if !names["plugin_schema"] {
t.Error("expected plugin_schema check")
}
if !names["no_duplicate_names"] {
t.Error("expected no_duplicate_names check")
}
}

func TestValidationResult_Fields(t *testing.T) {
r := mktvalidator.ValidationResult{
CheckName: "test_check",
Passed:    true,
Warnings:  []string{"w1"},
Errors:    nil,
}
if r.CheckName != "test_check" {
t.Errorf("unexpected check name: %s", r.CheckName)
}
if !r.Passed {
t.Error("expected Passed=true")
}
if len(r.Warnings) != 1 || r.Warnings[0] != "w1" {
t.Errorf("unexpected warnings: %v", r.Warnings)
}
}
