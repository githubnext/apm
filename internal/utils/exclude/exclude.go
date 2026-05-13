// Package exclude provides glob-style pattern matching for filtering paths
// against compilation.exclude patterns from apm.yml.
//
// Supports ** (recursive directory) wildcard matching with a bounded-recursion
// guard to prevent exponential blowup.
package exclude

import (
	"fmt"
	"path/filepath"
	"strings"
)

// MaxDoubleStarSegments is the maximum number of ** segments allowed in a
// single pattern to prevent exponential recursion blowup.
const MaxDoubleStarSegments = 5

// ValidateExcludePatterns validates and normalizes exclude patterns, rejecting
// dangerous ones. Returns the normalized patterns or an error if any pattern
// exceeds the ** segment safety limit.
func ValidateExcludePatterns(patterns []string) ([]string, error) {
	if len(patterns) == 0 {
		return nil, nil
	}
	validated := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		normalized := strings.ReplaceAll(pattern, "\\", "/")
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
		starCount := 0
		for _, p := range collapsed {
			if p == "**" {
				starCount++
			}
		}
		if starCount > MaxDoubleStarSegments {
			return nil, fmt.Errorf(
				"exclude: pattern %q has %d '**' segments (max %d); simplify the pattern",
				pattern, starCount, MaxDoubleStarSegments,
			)
		}
		validated = append(validated, normalized)
	}
	return validated, nil
}

// ShouldExclude checks whether a file path should be excluded based on the
// pre-validated patterns. baseDir is used to compute the relative path.
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
	if strings.HasPrefix(relStr, "../") {
		return false
	}
	for _, pattern := range excludePatterns {
		if matchesPattern(relStr, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a relative path string matches a single exclusion pattern.
func matchesPattern(relPathStr, pattern string) bool {
	if strings.Contains(pattern, "**") {
		pathParts := strings.Split(relPathStr, "/")
		patternParts := strings.Split(pattern, "/")
		return matchGlobRecursive(pathParts, patternParts)
	}
	ok, _ := filepath.Match(pattern, relPathStr)
	if ok {
		return true
	}
	// Directory prefix matching
	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(relPathStr, pattern) || relPathStr == strings.TrimSuffix(pattern, "/")
	}
	return strings.HasPrefix(relPathStr, pattern+"/") || relPathStr == pattern
}

// matchGlobRecursive matches path components against pattern components with ** support.
func matchGlobRecursive(pathParts, patternParts []string) bool {
	// Strip trailing empty parts
	for len(patternParts) > 0 && patternParts[len(patternParts)-1] == "" {
		patternParts = patternParts[:len(patternParts)-1]
	}

	pi, xi := 0, 0
	for pi < len(patternParts) && xi < len(pathParts) {
		part := patternParts[pi]
		if part == "**" {
			break
		}
		ok, _ := filepath.Match(part, pathParts[xi])
		if !ok {
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

// matchDoubleStar handles ** segments with bounded recursion.
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
	ok, _ := filepath.Match(part, pathParts[0])
	if ok {
		return matchDoubleStar(pathParts[1:], patternParts[1:])
	}
	return false
}
