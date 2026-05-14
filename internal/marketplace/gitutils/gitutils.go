// Package gitutils provides shared git-related utilities for marketplace modules.
// Migrated from src/apm_cli/marketplace/_git_utils.py
package gitutils

import "regexp"

// tokenRE matches auth tokens in git URLs.
// Covers: https://TOKEN@host, http://TOKEN@host, and ?token=VALUE query params.
var tokenRE = regexp.MustCompile(`https?://[^@\s]*@|([?&])token=[^\s&]*`)

// RedactToken replaces auth tokens in text with redacted placeholders.
func RedactToken(text string) string {
	return tokenRE.ReplaceAllStringFunc(text, func(m string) string {
		for _, r := range m {
			if r == '@' {
				return "https://***@"
			}
		}
		// query-param match: preserve the leading ? or &
		if len(m) > 0 && (m[0] == '?' || m[0] == '&') {
			return string(m[0]) + "token=***"
		}
		return m
	})
}
