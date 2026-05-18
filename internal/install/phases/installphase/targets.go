// Package installphase provides Go implementations of install pipeline phases.
// Migrated from src/apm_cli/install/phases/targets.py
package installphase

import (
	"os"
	"path/filepath"
	"strings"
)

// Target represents an integration target (e.g. "claude", "vscode").
type Target struct {
	// Name is the canonical target name.
	Name string
	// ConfigDir is the configuration directory for this target, if known.
	ConfigDir string
}

// TargetDetectionResult holds the outcome of target detection.
type TargetDetectionResult struct {
	// Targets is the ordered list of detected or user-specified targets.
	Targets []Target
	// Provenance describes how the targets were determined.
	Provenance string
	// Integrators maps primitive type names to integrator instances.
	// Values are opaque (interface{}) because integrators are Go implementations
	// of the various BaseIntegrator subclasses.
	Integrators map[string]interface{}
}

// TargetSource indicates where the target list came from.
type TargetSource int

const (
	TargetSourceCLI     TargetSource = iota // --target flag
	TargetSourceYAML                        // targets: in apm.yml
	TargetSourceEnv                         // APM_TARGET env var
	TargetSourceDetect                      // auto-detection
)

// ParseTargetsField parses a targets/target YAML field (string or string list).
// Returns nil when neither key is present.
func ParseTargetsField(data map[string]interface{}) []string {
	var raw interface{}
	if v, ok := data["targets"]; ok {
		raw = v
	} else if v, ok := data["target"]; ok {
		raw = v
	} else {
		return nil
	}
	switch v := raw.(type) {
	case string:
		var result []string
		for _, t := range strings.Split(v, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				result = append(result, t)
			}
		}
		return result
	case []interface{}:
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				s = strings.TrimSpace(s)
				if s != "" {
					result = append(result, s)
				}
			}
		}
		return result
	}
	return nil
}

// ReadYAMLTargets reads the targets/target field from an apm.yml file.
// Returns nil when neither key is present, or on any error.
func ReadYAMLTargets(apmYMLPath string) []string {
	path := filepath.Join(apmYMLPath, "apm.yml")
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	// Minimal YAML parse: scan for "targets:" or "target:" key
	parsed := parseSimpleYAMLMap(string(data))
	return ParseTargetsField(parsed)
}

// parseSimpleYAMLMap is a line-scanner for simple flat YAML maps (no nesting).
func parseSimpleYAMLMap(content string) map[string]interface{} {
	result := map[string]interface{}{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		result[key] = val
	}
	return result
}

// KnownTargets is the canonical set of supported integration targets.
var KnownTargets = map[string]bool{
	"claude":       true,
	"vscode":       true,
	"windsurf":     true,
	"cursor":       true,
	"opencode":     true,
	"codex":        true,
	"copilot":      true,
	"all":          true,
}

// ValidateTargets checks that all specified targets are known.
// Returns a list of unknown target names.
func ValidateTargets(targets []string) []string {
	var unknown []string
	for _, t := range targets {
		if !KnownTargets[strings.ToLower(t)] {
			unknown = append(unknown, t)
		}
	}
	return unknown
}

// ExpandAllTarget replaces "all" with the full list of known targets (except "all").
func ExpandAllTarget(targets []string) []string {
	for _, t := range targets {
		if strings.ToLower(t) == "all" {
			var all []string
			for name := range KnownTargets {
				if name != "all" {
					all = append(all, name)
				}
			}
			return all
		}
	}
	return targets
}

// FormatProvenance returns a human-readable description of target provenance.
func FormatProvenance(source TargetSource, value string) string {
	switch source {
	case TargetSourceCLI:
		return "from --target flag: " + value
	case TargetSourceYAML:
		return "from apm.yml targets field: " + value
	case TargetSourceEnv:
		return "from APM_TARGET environment variable: " + value
	case TargetSourceDetect:
		return "auto-detected: " + value
	default:
		return value
	}
}

// DetectTargetsFromEnv reads the APM_TARGET env var.
// Returns nil when the variable is unset or empty.
func DetectTargetsFromEnv() []string {
	val := os.Getenv("APM_TARGET")
	if val == "" {
		return nil
	}
	var targets []string
	for _, t := range strings.Split(val, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			targets = append(targets, t)
		}
	}
	return targets
}
