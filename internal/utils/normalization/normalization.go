// Package normalization provides bytes-in / bytes-out content normalization helpers.
// Migrated from src/apm_cli/utils/normalization.py
package normalization

import (
	"bytes"
	"regexp"
)

var (
	buildIDPattern = regexp.MustCompile(`(?i)<!--\s*Build ID:\s*[a-f0-9]+\s*-->\s*\n?`)
	bom            = []byte{0xef, 0xbb, 0xbf}
)

// StripBuildID removes APM <!-- Build ID: <sha> --> headers.
func StripBuildID(content []byte) []byte {
	return buildIDPattern.ReplaceAll(content, nil)
}

// NormalizeLineEndings converts CRLF to LF.
func NormalizeLineEndings(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
}

// StripBOM drops a UTF-8 BOM at the start of the file.
func StripBOM(content []byte) []byte {
	if bytes.HasPrefix(content, bom) {
		return content[len(bom):]
	}
	return content
}

// Normalize applies all drift-tolerant normalizations to a file's bytes.
func Normalize(content []byte) []byte {
	return StripBuildID(NormalizeLineEndings(StripBOM(content)))
}
