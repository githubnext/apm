// Package compilation provides compilation pipeline types and utilities for APM.
package compilation

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// BuildIDPlaceholder is the sentinel string inserted into compiled outputs.
const BuildIDPlaceholder = "<!-- Build ID: __BUILD_ID__ -->"

// ConstitutionMarkerBegin marks the start of an injected constitution block.
const ConstitutionMarkerBegin = "<!-- SPEC-KIT CONSTITUTION: BEGIN -->"

// ConstitutionMarkerEnd marks the end of an injected constitution block.
const ConstitutionMarkerEnd = "<!-- SPEC-KIT CONSTITUTION: END -->"

// ConstitutionRelativePath is the repo-root-relative path to the constitution file.
const ConstitutionRelativePath = ".specify/memory/constitution.md"

// StabilizeBuildID replaces BuildIDPlaceholder in content with a deterministic
// 12-char SHA256 hash computed over the content with the placeholder line removed.
//
// Idempotent: returns content unchanged if no placeholder is present.
// Preserves a trailing newline when the input has one.
func StabilizeBuildID(content string) string {
	lines := strings.Split(content, "\n")

	// Preserve trailing newline: splitlines leaves an empty last element
	// if content ends with "\n". We track it but exclude from hash input.
	trailingNL := strings.HasSuffix(content, "\n")

	idx := -1
	for i, line := range lines {
		if line == BuildIDPlaceholder {
			idx = i
			break
		}
	}
	if idx == -1 {
		return content
	}

	// Build hash input from all lines except the placeholder.
	hashLines := make([]string, 0, len(lines)-1)
	for i, line := range lines {
		if i != idx {
			hashLines = append(hashLines, line)
		}
	}
	// Remove the empty trailing element that split adds for "\n"-terminated content.
	if trailingNL && len(hashLines) > 0 && hashLines[len(hashLines)-1] == "" {
		hashLines = hashLines[:len(hashLines)-1]
	}
	h := sha256.Sum256([]byte(strings.Join(hashLines, "\n")))
	buildID := fmt.Sprintf("%x", h)[:12]

	lines[idx] = fmt.Sprintf("<!-- Build ID: %s -->", buildID)

	result := strings.Join(lines, "\n")
	if trailingNL && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result
}
