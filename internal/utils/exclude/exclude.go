// Package exclude provides glob-style pattern matching for compilation and primitive discovery.
//
// Supports ** (recursive directory) patterns. Used to filter paths against
// compilation.exclude patterns from apm.yml.
package exclude

import (
	"errors"
	"path/filepath"
	"strings"
)

// MaxDoubleStarSegments is the maximum number of ** segments allowed in a pattern.
const MaxDoubleStarSegments = 5

// ValidateExcludePatterns validates and normalises exclude patterns.
func ValidateExcludePatterns(patterns []string) ([]string, error) {
	if len(patterns) == 0 {
		return nil, nil
	}
	validated := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		normalized := strings.ReplaceAll(pattern, `\`, "/")
		parts := strings.Split(normalized, "/")
		// Collapse consecutive ** segments
		collapsed := make([]string, 0, len(parts))
		for _, p := range parts {
			if p == "**" && len(collapsed) > 0 && collapsed[len(collapsed)-1] == "**" {
				continue
			}
			collapsed = append(collapsed, p)
		}
		normalized = strings.Join(collapsed, "/")
		count := 0
		for _, p := range collapsed {
			if p == "**" {
				count++
			}
		}
		if count > MaxDoubleStarSegments {
			return nil, errors.New("exclude pattern '" + pattern + "' has too many ** segments (max 5)")
		}
		validated = append(validated, normalized)
	}
	return validated, nil
}

// ShouldExclude checks whether a file path should be excluded.
func ShouldExclude(filePath, baseDir string, excludePatterns []string) bool {
	if len(excludePatterns) == 0 {
		return false
	}
	absFile, err := filepath.Abs(filePath)
	if err != nil {
		absFile = filePath
	}
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		absBase = baseDir
	}
	rel, err := filepath.Rel(absBase, absFile)
	if err != nil {
		return false
	}
	relStr := filepath.ToSlash(rel)
	for _, pattern := range excludePatterns {
		if matchesPattern(relStr, pattern) {
			return true
		}
	}
	return false
}

func matchesPattern(relPath, pattern string) bool {
	if strings.Contains(pattern, "**") {
		pathParts := strings.Split(relPath, "/")
		patternParts := strings.Split(pattern, "/")
		return matchGlobRecursive(pathParts, patternParts)
	}
	if matched, _ := filepath.Match(pattern, relPath); matched {
		return true
	}
	if strings.HasSuffix(pattern, "/") {
		if strings.HasPrefix(relPath, pattern) || relPath == strings.TrimSuffix(pattern, "/") {
			return true
		}
	} else if strings.HasPrefix(relPath, pattern+"/") || relPath == pattern {
		return true
	}
	return false
}

func matchGlobRecursive(pathParts, patternParts []string) bool {
	// Strip trailing empty parts
	for len(patternParts) > 0 && patternParts[len(patternParts)-1] == "" {
		patternParts = patternParts[:len(patternParts)-1]
	}
	pi, xi := 0, 0
	// Fast iterative path for leading non-** segments
	for pi < len(patternParts) && xi < len(pathParts) {
		part := patternParts[pi]
		if part == "**" {
			break
		}
		matched, _ := filepath.Match(part, pathParts[xi])
		if !matched {
			return false
		}
		pi++
		xi++
	}
	if pi == len(patternParts) {
		return xi == len(pathParts)
	}
	return matchDoubleStar(pathParts[xi:], patternParts[pi:])
}

func matchDoubleStar(pathParts, patternParts []string) bool {
	if len(patternParts) == 0 {
		return len(pathParts) == 0
	}
	if len(pathParts) == 0 {
		for _, p := range patternParts {
			if p != "**" && p != "" {
				return false
			}
		}
		return true
	}
	part := patternParts[0]
	if part == "**" {
		if matchDoubleStar(pathParts, patternParts[1:]) {
			return true
		}
		return matchDoubleStar(pathParts[1:], patternParts)
	}
	matched, _ := filepath.Match(part, pathParts[0])
	if matched {
		return matchDoubleStar(pathParts[1:], patternParts[1:])
	}
	return false
}
