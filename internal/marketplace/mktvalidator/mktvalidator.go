// Package mktvalidator provides marketplace manifest validation.
package mktvalidator

// Plugin is a minimal plugin record for validation.
type Plugin struct {
Name   string
Source string
}

// ValidationResult holds the result of a single validation check.
type ValidationResult struct {
CheckName string
Passed    bool
Warnings  []string
Errors    []string
}

// ValidatePluginSchema checks that all plugins have required fields.
func ValidatePluginSchema(plugins []Plugin) ValidationResult {
r := ValidationResult{CheckName: "plugin_schema", Passed: true}
for _, p := range plugins {
if p.Name == "" {
r.Errors = append(r.Errors, "Plugin entry has empty name")
r.Passed = false
}
if p.Source == "" {
r.Errors = append(r.Errors, "Plugin '"+p.Name+"' has empty source")
r.Passed = false
}
}
return r
}

// ValidateNoDuplicateNames checks for duplicate plugin names.
func ValidateNoDuplicateNames(plugins []Plugin) ValidationResult {
r := ValidationResult{CheckName: "no_duplicate_names", Passed: true}
seen := map[string]bool{}
for _, p := range plugins {
if seen[p.Name] {
r.Errors = append(r.Errors, "Duplicate plugin name: "+p.Name)
r.Passed = false
}
seen[p.Name] = true
}
return r
}

// ValidateMarketplace runs all validation checks on a list of plugins.
func ValidateMarketplace(plugins []Plugin) []ValidationResult {
return []ValidationResult{
ValidatePluginSchema(plugins),
ValidateNoDuplicateNames(plugins),
}
}
