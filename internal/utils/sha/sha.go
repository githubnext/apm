// Package sha provides short SHA formatting helpers.
// Mirrors src/apm_cli/utils/short_sha.py.
package sha

import "strings"

// sentinels are values that collapse to empty string.
var sentinels = map[string]bool{
	"cached":  true,
	"unknown": true,
}

// FormatShortSHA returns an 8-char short SHA or "" for invalid inputs.
// Rules:
//   - nil / non-string -> ""
//   - sentinel strings ("cached", "unknown") -> ""
//   - shorter than 8 chars -> ""
//   - contains non-hex characters -> ""
//   - otherwise: first 8 chars
func FormatShortSHA(value string) string {
	candidate := strings.TrimSpace(value)
	if candidate == "" {
		return ""
	}
	if sentinels[strings.ToLower(candidate)] {
		return ""
	}
	if len(candidate) < 8 {
		return ""
	}
	for _, ch := range candidate {
		if !isHex(ch) {
			return ""
		}
	}
	return candidate[:8]
}

func isHex(ch rune) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}
