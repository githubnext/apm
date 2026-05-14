// Package experimental provides a feature-flag subsystem for the APM CLI.
// Migrated from src/apm_cli/core/experimental.py
package experimental

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Flag describes a single experimental feature.
type Flag struct {
	// Name is the internal snake_case identifier.
	Name string
	// Description is a one-line summary (<= 80 chars, printable ASCII).
	Description string
	// Default is the registry default -- always false.
	Default bool
	// Hint is an optional next-step message shown after enabling.
	Hint string
}

// registry is the static map of all registered experimental flags.
var registry = map[string]Flag{
	"verbose_version": {
		Name:        "verbose_version",
		Description: "Show Python version, platform, and install path in 'apm --version'.",
		Default:     false,
		Hint:        "Run 'apm --version' to see the new output.",
	},
	"copilot_cowork": {
		Name:        "copilot_cowork",
		Description: "Enable Microsoft 365 Copilot Cowork skills deployment via OneDrive.",
		Default:     false,
		Hint: "Use '--target copilot-cowork --global' to deploy skills. " +
			"See https://microsoft.github.io/apm/integrations/copilot-cowork/",
	},
}

// Flags returns the static registry (read-only view).
func Flags() map[string]Flag {
	return registry
}

// normalizeFlagName normalizes a CLI flag name to internal snake_case.
func normalizeFlagName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, "-", "_"))
}

// DisplayName converts an internal snake_case name to kebab-case for display.
func DisplayName(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}

// configPath returns the path to ~/.apm/config.json.
func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".apm", "config.json")
}

var (
	configMu    sync.RWMutex
	configCache map[string]interface{}
)

// loadConfig reads ~/.apm/config.json, returning an empty map on failure.
func loadConfig() map[string]interface{} {
	configMu.RLock()
	if configCache != nil {
		defer configMu.RUnlock()
		return configCache
	}
	configMu.RUnlock()

	configMu.Lock()
	defer configMu.Unlock()
	if configCache != nil {
		return configCache
	}
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		configCache = map[string]interface{}{}
		return configCache
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		configCache = map[string]interface{}{}
		return configCache
	}
	configCache = cfg
	return configCache
}

// invalidateCache clears the config cache so the next load re-reads disk.
func invalidateCache() {
	configMu.Lock()
	configCache = nil
	configMu.Unlock()
}

// getExperimentalSection returns the "experimental" section from config.
func getExperimentalSection() map[string]interface{} {
	cfg := loadConfig()
	v, ok := cfg["experimental"]
	if !ok {
		return map[string]interface{}{}
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return m
}

// updateConfig merges updates into ~/.apm/config.json.
func updateConfig(updates map[string]interface{}) error {
	invalidateCache()
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Read existing
	var cfg map[string]interface{}
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, &cfg)
	}
	if cfg == nil {
		cfg = map[string]interface{}{}
	}
	for k, v := range updates {
		cfg[k] = v
	}
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(append(out, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// IsEnabled reports whether an experimental flag is currently enabled.
// Returns an error if the flag name is not registered.
func IsEnabled(name string) (bool, error) {
	if _, ok := registry[name]; !ok {
		keys := make([]string, 0, len(registry))
		for k := range registry {
			keys = append(keys, k)
		}
		return false, fmt.Errorf("unknown experimental flag: %q; registered: %s",
			name, strings.Join(keys, ", "))
	}
	experimental := getExperimentalSection()
	v, ok := experimental[name]
	if !ok {
		return registry[name].Default, nil
	}
	b, ok := v.(bool)
	if !ok {
		return registry[name].Default, nil
	}
	return b, nil
}

// ValidateFlagName validates and normalizes a flag name from CLI input.
// Returns the normalized name or an error with suggestions.
func ValidateFlagName(name string) (string, error) {
	normalized := normalizeFlagName(name)
	if _, ok := registry[normalized]; ok {
		return normalized, nil
	}
	display := DisplayName(normalized)
	// Build suggestions via simple prefix/contains matching.
	var suggestions []string
	for k := range registry {
		if strings.Contains(k, normalized) || strings.Contains(normalized, k) {
			suggestions = append(suggestions, DisplayName(k))
		}
	}
	msg := fmt.Sprintf("unknown experimental feature: %s", display)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("; did you mean: %s?", strings.Join(suggestions, ", "))
	}
	return "", fmt.Errorf("%s", msg)
}

// setFlag sets an experimental flag to a boolean value and persists it.
func setFlag(name string, value bool) (Flag, error) {
	flag, ok := registry[name]
	if !ok {
		return Flag{}, fmt.Errorf("unknown flag: %s", name)
	}
	experimental := map[string]interface{}{}
	for k, v := range getExperimentalSection() {
		experimental[k] = v
	}
	experimental[name] = value
	if err := updateConfig(map[string]interface{}{"experimental": experimental}); err != nil {
		return Flag{}, err
	}
	return flag, nil
}

// Enable enables an experimental flag and persists the change.
func Enable(name string) (Flag, error) {
	return setFlag(name, true)
}

// Disable disables an experimental flag and persists the change.
func Disable(name string) (Flag, error) {
	return setFlag(name, false)
}

// Reset resets one or all experimental flags to registry defaults.
// When name is empty, all flags are cleared. Returns the number removed.
func Reset(name string) (int, error) {
	experimental := map[string]interface{}{}
	for k, v := range getExperimentalSection() {
		experimental[k] = v
	}
	if name != "" {
		if _, ok := experimental[name]; ok {
			delete(experimental, name)
			if err := updateConfig(map[string]interface{}{"experimental": experimental}); err != nil {
				return 0, err
			}
			return 1, nil
		}
		return 0, nil
	}
	count := len(experimental)
	if count > 0 {
		if err := updateConfig(map[string]interface{}{"experimental": map[string]interface{}{}}); err != nil {
			return 0, err
		}
	}
	return count, nil
}

// GetOverriddenFlags returns flags that have user overrides in config.
func GetOverriddenFlags() map[string]bool {
	experimental := getExperimentalSection()
	out := map[string]bool{}
	for k, v := range experimental {
		if _, ok := registry[k]; !ok {
			continue
		}
		if b, ok := v.(bool); ok {
			out[k] = b
		}
	}
	return out
}

// GetStaleConfigKeys returns config keys not in the registry.
func GetStaleConfigKeys() []string {
	experimental := getExperimentalSection()
	var out []string
	for k := range experimental {
		if _, ok := registry[k]; !ok {
			out = append(out, k)
		}
	}
	return out
}

// GetMalformedFlagKeys returns registered flags with non-boolean config values.
func GetMalformedFlagKeys() []string {
	experimental := getExperimentalSection()
	var out []string
	for k, v := range experimental {
		if _, ok := registry[k]; !ok {
			continue
		}
		if _, ok := v.(bool); !ok {
			out = append(out, k)
		}
	}
	return out
}
