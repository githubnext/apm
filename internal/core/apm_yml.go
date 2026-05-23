package core

import "fmt"

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

// ParseTargetsField parses the targets/target field from a raw apm.yml data
// map. Returns a canonical list of target names. An empty list means neither
// key was present (caller should fall through to auto-detect).
func ParseTargetsField(yamlData map[string]interface{}) ([]string, error) {
	_, hasTargets := yamlData["targets"]
	_, hasTarget := yamlData["target"]

	if hasTargets && hasTarget {
		return nil, NewConflictingTargetsError()
	}

	if hasTargets {
		raw := yamlData["targets"]
		switch v := raw.(type) {
		case nil:
			return nil, NewEmptyTargetsListError()
		case []interface{}:
			if len(v) == 0 {
				return nil, NewEmptyTargetsListError()
			}
			tokens := make([]string, 0, len(v))
			for _, item := range v {
				t := fmt.Sprintf("%v", item)
				if t != "" {
					tokens = append(tokens, t)
				}
			}
			if err := validateCanonical(tokens); err != nil {
				return nil, err
			}
			return tokens, nil
		default:
			// Single value under targets key
			tokens := []string{fmt.Sprintf("%v", v)}
			if err := validateCanonical(tokens); err != nil {
				return nil, err
			}
			return tokens, nil
		}
	}

	if hasTarget {
		raw := yamlData["target"]
		if raw == nil {
			return []string{}, nil
		}
		switch v := raw.(type) {
		case []interface{}:
			tokens := make([]string, 0, len(v))
			for _, item := range v {
				t := fmt.Sprintf("%v", item)
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
		default:
			rawStr := fmt.Sprintf("%v", v)
			if rawStr == "" {
				return []string{}, nil
			}
			// CSV sugar: "claude,copilot" -> ["claude", "copilot"]
			tokens := splitCSV(rawStr)
			if len(tokens) == 0 {
				return []string{}, nil
			}
			if err := validateCanonical(tokens); err != nil {
				return nil, err
			}
			return tokens, nil
		}
	}

	return []string{}, nil
}

// validateCanonical checks every token is in CanonicalTargets.
func validateCanonical(tokens []string) error {
	valid := sortedKeys(CanonicalTargets)
	for _, t := range tokens {
		if !CanonicalTargets[t] {
			return NewUnknownTargetError(t, valid)
		}
	}
	return nil
}
