// Package experimental implements the "apm experimental" command group.
//
// Provides "apm experimental list|enable|disable|reset" to manage
// opt-in feature flags stored in ~/.apm/config.json.
//
// Migrated from: src/apm_cli/commands/experimental.py
package experimental

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Flag describes one experimental feature flag.
type Flag struct {
	Name        string
	DisplayName string
	Description string
	Default     bool
}

// KnownFlags is the registry of all experimental flags.
// Mirrors core/experimental.FLAGS in Python.
var KnownFlags = []Flag{
	{
		Name:        "parallel-install",
		DisplayName: "Parallel Install",
		Description: "Download and install packages concurrently.",
		Default:     false,
	},
	{
		Name:        "incremental-compilation",
		DisplayName: "Incremental Compilation",
		Description: "Skip unchanged primitives during apm compile.",
		Default:     false,
	},
	{
		Name:        "strict-policy",
		DisplayName: "Strict Policy",
		Description: "Treat policy warnings as errors.",
		Default:     false,
	},
	{
		Name:        "mcp-auto-configure",
		DisplayName: "MCP Auto Configure",
		Description: "Automatically write MCP configs during install.",
		Default:     false,
	},
	{
		Name:        "telemetry",
		DisplayName: "Telemetry",
		Description: "Send anonymised usage data to improve APM.",
		Default:     false,
	},
}

// ---------------------------------------------------------
// Config file
// ---------------------------------------------------------

// Config holds the full ~/.apm/config.json content.
type Config struct {
	ExperimentalFlags map[string]bool `json:"experimental_flags,omitempty"`
}

// configPath returns ~/.apm/config.json.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".apm", "config.json"), nil
}

// loadConfig reads the config file; returns an empty Config on ENOENT.
func loadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{ExperimentalFlags: make(map[string]bool)}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	if cfg.ExperimentalFlags == nil {
		cfg.ExperimentalFlags = make(map[string]bool)
	}
	return cfg, nil
}

// saveConfig writes the config file atomically.
func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// ---------------------------------------------------------
// Public API
// ---------------------------------------------------------

// IsEnabled reports whether a flag is currently enabled.
func IsEnabled(name string) (bool, error) {
	cfg, err := loadConfig()
	if err != nil {
		return false, err
	}
	if v, ok := cfg.ExperimentalFlags[name]; ok {
		return v, nil
	}
	// Fall back to the flag's default.
	for _, f := range KnownFlags {
		if f.Name == name {
			return f.Default, nil
		}
	}
	return false, nil
}

// EnableFlag enables a named flag. Returns an error for unknown flags.
func EnableFlag(name string) error {
	name = NormaliseFlag(name)
	if !isKnown(name) {
		return fmt.Errorf("unknown flag %q -- run 'apm experimental list' to see available flags", name)
	}
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	cfg.ExperimentalFlags[name] = true
	if err := saveConfig(cfg); err != nil {
		return err
	}
	fmt.Printf("[+] Enabled: %s\n", name)
	return nil
}

// DisableFlag disables a named flag.
func DisableFlag(name string) error {
	name = NormaliseFlag(name)
	if !isKnown(name) {
		return fmt.Errorf("unknown flag %q", name)
	}
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	cfg.ExperimentalFlags[name] = false
	if err := saveConfig(cfg); err != nil {
		return err
	}
	fmt.Printf("[i] Disabled: %s\n", name)
	return nil
}

// ResetFlags clears all experimental flag overrides, restoring defaults.
func ResetFlags() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	cfg.ExperimentalFlags = make(map[string]bool)
	if err := saveConfig(cfg); err != nil {
		return err
	}
	fmt.Println("[+] All experimental flags reset to defaults.")
	return nil
}

// ListFlags returns all known flags with their current enabled state.
func ListFlags() ([]FlagStatus, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	var out []FlagStatus
	for _, f := range KnownFlags {
		enabled := f.Default
		if v, ok := cfg.ExperimentalFlags[f.Name]; ok {
			enabled = v
		}
		out = append(out, FlagStatus{
			Flag:       f,
			Enabled:    enabled,
			Overridden: cfg.ExperimentalFlags[f.Name] != f.Default,
		})
	}
	// Sort by name for deterministic output.
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// FlagStatus combines a Flag definition with its current runtime state.
type FlagStatus struct {
	Flag
	Enabled    bool
	Overridden bool
}

// GetOverriddenFlags returns only flags that differ from their default.
func GetOverriddenFlags() ([]FlagStatus, error) {
	all, err := ListFlags()
	if err != nil {
		return nil, err
	}
	var out []FlagStatus
	for _, f := range all {
		if f.Overridden {
			out = append(out, f)
		}
	}
	return out, nil
}

// GetMalformedFlagKeys returns keys in the config that are not valid flag names.
func GetMalformedFlagKeys() ([]string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	var bad []string
	for k := range cfg.ExperimentalFlags {
		if !isKnown(k) {
			bad = append(bad, k)
		}
	}
	sort.Strings(bad)
	return bad, nil
}

// GetStaleConfigKeys is an alias for GetMalformedFlagKeys (legacy name).
func GetStaleConfigKeys() ([]string, error) { return GetMalformedFlagKeys() }

// ValidateFlagName returns an error if name is not a known flag.
func ValidateFlagName(name string) error {
	n := NormaliseFlag(name)
	if !isKnown(n) {
		return fmt.Errorf("unknown flag %q", name)
	}
	return nil
}

// NormaliseFlag lowercases and trims a flag name.
func NormaliseFlag(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// DisplayName returns the human-readable display name for a flag.
func DisplayName(name string) string {
	for _, f := range KnownFlags {
		if f.Name == NormaliseFlag(name) {
			return f.DisplayName
		}
	}
	return name
}

// ---------------------------------------------------------
// Helpers
// ---------------------------------------------------------

func isKnown(name string) bool {
	for _, f := range KnownFlags {
		if f.Name == name {
			return true
		}
	}
	return false
}
