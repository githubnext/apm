// Package versionpins provides a ref pin cache for marketplace plugin immutability checks.
//
// Mirrors src/apm_cli/marketplace/version_pins.py.
//
// Records plugin-to-ref mappings per marketplace, keyed on the plugin's declared
// "version" field from the standard marketplace spec. When the same
// (marketplace, plugin, version) triple resolves to a different ref, a warning
// is emitted -- this may indicate a ref-swap attack.
//
// The pin file lives at ~/.apm/cache/marketplace/version-pins.json.
// All functions are fail-open: filesystem or JSON errors are silently ignored.
package versionpins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const pinsFilename = "version-pins.json"

// pinsPath returns the full path to the version-pins JSON file.
// If pinsDir is empty, the default ~/.apm/cache/marketplace/ is used.
func pinsPath(pinsDir string) string {
	if pinsDir != "" {
		return filepath.Join(pinsDir, pinsFilename)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".apm", "cache", "marketplace", pinsFilename)
}

// pinKey builds the canonical dict key for a marketplace/plugin/version triple.
func pinKey(marketplaceName, pluginName, version string) string {
	base := fmt.Sprintf("%s/%s", strings.ToLower(marketplaceName), strings.ToLower(pluginName))
	if version != "" {
		return fmt.Sprintf("%s/%s", base, strings.ToLower(version))
	}
	return base
}

// LoadRefPins loads the ref-pins file from disk.
// Returns an empty map when the file is missing or contains invalid JSON.
// Never returns an error.
func LoadRefPins(pinsDir string) map[string]string {
	path := pinsPath(pinsDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return map[string]string{}
	}
	result := make(map[string]string, len(raw))
	for k, v := range raw {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result
}

// SaveRefPins persists pins to disk atomically using a temp file + os.Rename.
// Errors are silently ignored (advisory system).
func SaveRefPins(pins map[string]string, pinsDir string) {
	path := pinsPath(pinsDir)
	tmpPath := path + ".tmp"

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return
	}

	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return
	}
	_ = os.Rename(tmpPath, path)
}

// CheckRefPin checks whether ref matches the previously-recorded pin.
//
// Returns the previously pinned ref if it differs from ref (possible ref swap).
// Returns empty string if this is the first time seeing the plugin/version or the
// ref matches.
func CheckRefPin(marketplaceName, pluginName, ref, version, pinsDir string) string {
	pins := LoadRefPins(pinsDir)
	key := pinKey(marketplaceName, pluginName, version)
	previous, ok := pins[key]
	if !ok || previous == "" {
		return ""
	}
	if previous == ref {
		return ""
	}
	return previous
}

// RecordRefPin stores a plugin-to-ref mapping in the pin cache.
// Overwrites any existing pin for the same plugin/version.
func RecordRefPin(marketplaceName, pluginName, ref, version, pinsDir string) {
	pins := LoadRefPins(pinsDir)
	key := pinKey(marketplaceName, pluginName, version)
	pins[key] = ref
	SaveRefPins(pins, pinsDir)
}
