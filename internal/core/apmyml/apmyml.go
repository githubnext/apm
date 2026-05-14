// Package apmyml provides a schema parser for the targets/target field in
// apm.yml.
//
// Mirrors src/apm_cli/core/apm_yml.py.
//
// Rules:
//   - 'targets: [a, b]'   -> ["a", "b"]  (canonical, plural)
//   - 'target: a'         -> ["a"]        (singular sugar)
//   - 'target: "a,b"'     -> ["a", "b"]  (CSV sugar)
//   - 'target: [a, b]'    -> ["a", "b"]  (list sugar under singular key)
//   - both present        -> error
//   - neither present     -> []           (empty = auto-detect upstream)
package apmyml

import (
	"fmt"
	"sort"
	"strings"
)

// CanonicalTargets is the set of target names accepted by APM.
var CanonicalTargets = map[string]bool{
	"claude":       true,
	"copilot":      true,
	"cursor":       true,
	"opencode":     true,
	"codex":        true,
	"gemini":       true,
	"windsurf":     true,
	"agent-skills": true,
}

// ConflictingTargetsError is returned when both 'targets' and 'target' are
// present in an apm.yml.
type ConflictingTargetsError struct {
	Message string
}

func (e *ConflictingTargetsError) Error() string {
	return e.Message
}

// EmptyTargetsListError is returned when 'targets:' is present but empty.
type EmptyTargetsListError struct {
	Message string
}

func (e *EmptyTargetsListError) Error() string {
	return e.Message
}

// UnknownTargetError is returned when a target token is not in CanonicalTargets.
type UnknownTargetError struct {
	Token   string
	Message string
}

func (e *UnknownTargetError) Error() string {
	return e.Message
}

// sortedTargets returns the canonical targets in sorted order for error messages.
func sortedTargets() []string {
	out := make([]string, 0, len(CanonicalTargets))
	for t := range CanonicalTargets {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// validateCanonical checks every token is in CanonicalTargets.
func validateCanonical(tokens []string) error {
	for _, token := range tokens {
		if !CanonicalTargets[token] {
			known := sortedTargets()
			msg := fmt.Sprintf(
				"[x] Unknown target %q\n\nSupported targets: %s\n\nRun 'apm targets' to list all.",
				token, strings.Join(known, ", "),
			)
			return &UnknownTargetError{Token: token, Message: msg}
		}
	}
	return nil
}

// ParseTargetsField parses the targets/target field from raw apm.yml data.
//
// data is expected to be a map[string]interface{} decoded from YAML.
// Returns a canonical list of target names. An empty slice means neither key
// was present (caller should fall through to auto-detect).
func ParseTargetsField(data map[string]interface{}) ([]string, error) {
	_, hasTargets := data["targets"]
	_, hasTarget := data["target"]

	if hasTargets && hasTarget {
		msg := "[x] Both 'targets' and 'target' keys found in apm.yml\n\n" +
			"Use only 'targets' (canonical) or 'target' (sugar), not both.\n\n" +
			"Fix with:\n\n  apm init        # regenerate apm.yml\n"
		return nil, &ConflictingTargetsError{Message: msg}
	}

	if hasTargets {
		raw := data["targets"]
		if raw == nil {
			return nil, &EmptyTargetsListError{
				Message: "[x] 'targets:' in apm.yml is empty\n\nThe targets list must contain at least one target.\n",
			}
		}
		rawList, ok := raw.([]interface{})
		if !ok {
			// Single value under targets: key.
			token := strings.TrimSpace(fmt.Sprintf("%v", raw))
			if err := validateCanonical([]string{token}); err != nil {
				return nil, err
			}
			return []string{token}, nil
		}
		if len(rawList) == 0 {
			return nil, &EmptyTargetsListError{
				Message: "[x] 'targets:' in apm.yml is empty\n\nThe targets list must contain at least one target.\n",
			}
		}
		var tokens []string
		for _, item := range rawList {
			t := strings.TrimSpace(fmt.Sprintf("%v", item))
			if t != "" {
				tokens = append(tokens, t)
			}
		}
		if err := validateCanonical(tokens); err != nil {
			return nil, err
		}
		return tokens, nil
	}

	if hasTarget {
		raw := data["target"]
		if raw == nil {
			return []string{}, nil
		}
		// List sugar: 'target: [claude, copilot]'
		if rawList, ok := raw.([]interface{}); ok {
			var tokens []string
			for _, item := range rawList {
				t := strings.TrimSpace(fmt.Sprintf("%v", item))
				if t != "" {
					tokens = append(tokens, t)
				}
			}
			if len(tokens) == 0 {
				return []string{}, nil
			}
			if err := validateCanonical(tokens); err != nil {
				return nil, err
			}
			return tokens, nil
		}
		rawStr := strings.TrimSpace(fmt.Sprintf("%v", raw))
		if rawStr == "" {
			return []string{}, nil
		}
		// CSV sugar: "claude,copilot"
		parts := strings.Split(rawStr, ",")
		var tokens []string
		for _, p := range parts {
			t := strings.TrimSpace(p)
			if t != "" {
				tokens = append(tokens, t)
			}
		}
		if err := validateCanonical(tokens); err != nil {
			return nil, err
		}
		return tokens, nil
	}

	// Neither key present.
	return []string{}, nil
}
