package core

import (
	"sort"
	"strings"
)

// joinStrings joins a slice of strings with sep.
func joinStrings(ss []string, sep string) string {
	return strings.Join(ss, sep)
}

// splitCSV splits a comma-separated string, trimming whitespace.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

// sortStrings sorts a slice of strings in place.
func sortStrings(ss []string) {
	sort.Strings(ss)
}

// sortedKeys returns the keys of a map[string]bool sorted.
func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// stripBracketNoise removes leading/trailing []'" and space characters.
func stripBracketNoise(s string) string {
	return strings.Trim(s, "[]'\" ")
}
