// Package configcmd implements the "apm config" command group.
//
// Sub-commands:
//   - apm config          -- show current configuration
//   - apm config get KEY  -- get a configuration value
//   - apm config set KEY VALUE -- set a configuration value
//
// Corresponds to src/apm_cli/commands/config.py.
package configcmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Known configuration keys with their canonical display names.
var configKeyDisplayNames = map[string]string{
	"auto_integrate":            "auto-integrate",
	"temp_dir":                  "temp-dir",
	"copilot_cowork_skills_dir": "copilot-cowork-skills-dir",
}

// booleanTrueValues is the set of strings that mean true.
var booleanTrueValues = map[string]bool{"true": true, "1": true, "yes": true}

// booleanFalseValues is the set of strings that mean false.
var booleanFalseValues = map[string]bool{"false": true, "0": true, "no": true}

// ParseBoolValue parses a CLI boolean string.
func ParseBoolValue(value string) (bool, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if booleanTrueValues[normalized] {
		return true, nil
	}
	if booleanFalseValues[normalized] {
		return false, nil
	}
	return false, fmt.Errorf("invalid value %q; use 'true' or 'false'", value)
}

// APMConfig represents key fields parsed from apm.yml.
type APMConfig struct {
	Name        string
	Version     string
	Entrypoint  string
	MCPDepCount int
}

// parseAPMYML extracts known fields from apm.yml using a simple line scanner.
func parseAPMYML(content string) APMConfig {
	var cfg APMConfig
	inDeps := false
	inMCP := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		if indent == 0 {
			inDeps = false
			inMCP = false
		}

		kv := func(key string) string {
			if strings.HasPrefix(trimmed, key+":") {
				val := strings.TrimSpace(trimmed[len(key)+1:])
				// Strip quotes
				if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') ||
					(val[0] == '\'' && val[len(val)-1] == '\'')) {
					val = val[1 : len(val)-1]
				}
				return val
			}
			return ""
		}

		switch {
		case kv("name") != "" && indent == 0:
			cfg.Name = kv("name")
		case kv("version") != "" && indent == 0:
			cfg.Version = kv("version")
		case kv("entrypoint") != "" && indent == 0:
			cfg.Entrypoint = kv("entrypoint")
		case strings.HasPrefix(trimmed, "dependencies:") && indent == 0:
			inDeps = true
		case inDeps && strings.HasPrefix(trimmed, "mcp:"):
			inMCP = true
		case inMCP && strings.HasPrefix(trimmed, "-"):
			cfg.MCPDepCount++
		}
	}
	return cfg
}

// LoadAPMConfig reads apm.yml from the current directory.
func LoadAPMConfig() (*APMConfig, error) {
	raw, err := os.ReadFile("apm.yml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("configcmd: read apm.yml: %w", err)
	}
	cfg := parseAPMYML(string(raw))
	return &cfg, nil
}

// UserConfig holds persistent APM CLI settings stored in ~/.config/apm/config.json.
type UserConfig struct {
	AutoIntegrate bool   `json:"auto_integrate"`
	TempDir       string `json:"temp_dir"`
}

func userConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "apm", "config.json"), nil
}

// LoadUserConfig reads the user-level APM config.
func LoadUserConfig() (*UserConfig, error) {
	path, err := userConfigPath()
	if err != nil {
		return &UserConfig{}, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return &UserConfig{}, nil //nolint:nilerr // default config
	}
	var cfg UserConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return &UserConfig{}, nil //nolint:nilerr
	}
	return &cfg, nil
}

// SaveUserConfig persists the user-level APM config.
func SaveUserConfig(cfg *UserConfig) error {
	path, err := userConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o600)
}

// GetAutoIntegrate returns the current auto-integrate setting.
func GetAutoIntegrate() (bool, error) {
	cfg, err := LoadUserConfig()
	if err != nil {
		return false, err
	}
	return cfg.AutoIntegrate, nil
}

// SetAutoIntegrate sets the auto-integrate setting.
func SetAutoIntegrate(value bool) error {
	cfg, err := LoadUserConfig()
	if err != nil {
		return err
	}
	cfg.AutoIntegrate = value
	return SaveUserConfig(cfg)
}

// RunShow prints the current configuration.
func RunShow() error {
	cfg, err := LoadAPMConfig()
	if err != nil {
		return fmt.Errorf("[x] Error reading apm.yml: %w", err)
	}

	userCfg, _ := LoadUserConfig()

	fmt.Println()
	fmt.Printf("  %-16s  %-24s  %s\n", "CATEGORY", "SETTING", "VALUE")
	fmt.Printf("  %-16s  %-24s  %s\n", "--------", "-------", "-----")

	if cfg != nil {
		fmt.Printf("  %-16s  %-24s  %s\n", "Project", "Name", cfg.Name)
		fmt.Printf("  %-16s  %-24s  %s\n", "", "Version", cfg.Version)
		fmt.Printf("  %-16s  %-24s  %s\n", "", "Entrypoint", cfg.Entrypoint)
		fmt.Printf("  %-16s  %-24s  %d\n", "", "MCP Dependencies", cfg.MCPDepCount)
	}

	if userCfg != nil {
		fmt.Printf("  %-16s  %-24s  %v\n", "CLI", "auto-integrate", userCfg.AutoIntegrate)
		if userCfg.TempDir != "" {
			fmt.Printf("  %-16s  %-24s  %s\n", "", "temp-dir", userCfg.TempDir)
		}
	}

	fmt.Println()
	return nil
}

// RunGet prints the value for a configuration key.
func RunGet(key string) error {
	userCfg, err := LoadUserConfig()
	if err != nil {
		return err
	}
	switch key {
	case "auto-integrate":
		fmt.Println(userCfg.AutoIntegrate)
	case "temp-dir":
		fmt.Println(userCfg.TempDir)
	default:
		return fmt.Errorf("[x] Unknown config key %q. Valid keys: auto-integrate, temp-dir", key)
	}
	return nil
}

// RunSet sets a configuration key to value.
func RunSet(key, value string) error {
	userCfg, err := LoadUserConfig()
	if err != nil {
		return err
	}
	switch key {
	case "auto-integrate":
		b, err := ParseBoolValue(value)
		if err != nil {
			return err
		}
		userCfg.AutoIntegrate = b
		if err := SaveUserConfig(userCfg); err != nil {
			return err
		}
		fmt.Printf("[+] auto-integrate = %v\n", b)
	case "temp-dir":
		userCfg.TempDir = value
		if err := SaveUserConfig(userCfg); err != nil {
			return err
		}
		fmt.Printf("[+] temp-dir = %s\n", value)
	default:
		return fmt.Errorf("[x] Unknown config key %q. Valid keys: auto-integrate, temp-dir", key)
	}
	return nil
}

// ValidConfigKeys returns the list of valid configuration key names.
func ValidConfigKeys() []string {
	return []string{"auto-integrate", "temp-dir"}
}

// DisplayName returns the human-readable name for a config key.
func DisplayName(key string) string {
	if name, ok := configKeyDisplayNames[key]; ok {
		return name
	}
	return key
}
