// Package normalization provides bytes-in / bytes-out content normalization
// helpers. Mirrors src/apm_cli/utils/normalization.py.
//
// Used by drift-detection to compare deployed file bytes against the replay
// scratch tree without flagging legitimate, deterministic differences:
//   - Line-ending differences (CRLF vs LF)
//   - UTF-8 BOMs at the start of the file
//   - APM <!-- Build ID: <sha> --> headers re-stamped on every recompile
package normalization

import (
	"bytes"
	"regexp"
)

// BOM is the UTF-8 byte order mark.
var BOM = []byte{0xef, 0xbb, 0xbf}

// buildIDPattern matches APM <!-- Build ID: <sha> --> headers.
var buildIDPattern = regexp.MustCompile(`(?i)<!--\s*Build ID:\s*[a-f0-9]+\s*-->\s*\n?`)

// StripBuildID removes APM <!-- Build ID: <sha> --> headers wherever they appear.
func StripBuildID(content []byte) []byte {
	return buildIDPattern.ReplaceAll(content, nil)
}

// NormalizeLineEndings converts CRLF to LF; leaves bare CR alone.
func NormalizeLineEndings(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
}

// StripBOM drops a UTF-8 BOM at the start of the file (only at offset 0).
func StripBOM(content []byte) []byte {
	if bytes.HasPrefix(content, BOM) {
		return content[len(BOM):]
	}
	return content
}

// Normalize applies all drift-tolerant normalizations to a file's bytes.
func Normalize(content []byte) []byte {
	return StripBuildID(NormalizeLineEndings(StripBOM(content)))
}
