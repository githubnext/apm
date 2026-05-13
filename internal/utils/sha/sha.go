// Package sha provides short-form SHA helpers for user-facing output.
// Migrated from src/apm_cli/utils/short_sha.py
package sha

import "strings"

var sentinels = map[string]struct{}{
	"cached":  {},
	"unknown": {},
}

// FormatShortSHA returns an 8-char short SHA or "" for invalid inputs.
// Non-string inputs (empty string) and sentinel values collapse to "".
// Strings shorter than 8 chars or containing non-hex characters return "".
func FormatShortSHA(value string) string {
	candidate := strings.TrimSpace(value)
	if candidate == "" {
		return ""
	}
	if _, isSentinel := sentinels[strings.ToLower(candidate)]; isSentinel {
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

func isHex(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}
