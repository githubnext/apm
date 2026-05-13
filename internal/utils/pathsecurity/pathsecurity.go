// Package pathsecurity provides centralised path-security helpers for APM CLI.
//
// Every filesystem operation whose target is derived from user-controlled
// input must pass through one of these guards before touching the disk.
package pathsecurity

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// PathTraversalError is returned when a computed path escapes its expected base directory.
type PathTraversalError struct {
	msg string
}

func (e *PathTraversalError) Error() string { return e.msg }

func traversalErr(msg string) *PathTraversalError {
	return &PathTraversalError{msg: msg}
}

// ValidatePathSegments rejects path strings containing traversal sequences.
//
// Parameters:
//   - pathStr: path-like string to validate (repo URL, virtual path, etc.)
//   - context: human-readable label for error messages
//   - rejectEmpty: if true, also reject empty segments
//   - allowCurrentDir: if true, "." segments are accepted but ".." still rejected
func ValidatePathSegments(pathStr, context string, rejectEmpty, allowCurrentDir bool) error {
	reject := map[string]bool{"..": true}
	if !allowCurrentDir {
		reject["."] = true
	}
	for _, segment := range strings.Split(strings.ReplaceAll(pathStr, `\`, "/"), "/") {
		// Iteratively percent-decode each segment to catch multi-encoded traversal
		decoded := segment
		for i := 0; i < 8; i++ {
			next, err := url.PathUnescape(decoded)
			if err != nil || next == decoded {
				break
			}
			decoded = next
		}
		if reject[segment] || reject[decoded] {
			return traversalErr("Invalid " + context + " '" + pathStr + "': segment '" + segment + "' is a traversal sequence")
		}
		if rejectEmpty && segment == "" {
			return traversalErr("Invalid " + context + " '" + pathStr + "': path segments must not be empty")
		}
	}
	return nil
}

// IsPathTraversalError reports whether err is a PathTraversalError.
func IsPathTraversalError(err error) bool {
	var t *PathTraversalError
	return errors.As(err, &t)
}

// EnsurePathWithin resolves path and asserts it lives inside baseDir.
//
// Returns the resolved path on success. Raises PathTraversalError if the
// resolved path escapes baseDir.
func EnsurePathWithin(path, baseDir string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		// Fall back to Abs if EvalSymlinks fails (path may not exist yet)
		resolved, err = filepath.Abs(path)
		if err != nil {
			return "", traversalErr("Cannot resolve path '" + path + "': " + err.Error())
		}
	}
	resolvedBase, err := filepath.EvalSymlinks(baseDir)
	if err != nil {
		resolvedBase, err = filepath.Abs(baseDir)
		if err != nil {
			return "", traversalErr("Cannot resolve base dir '" + baseDir + "': " + err.Error())
		}
	}
	// Strip Windows extended-length prefix
	resolved = stripExtendedPrefix(resolved)
	resolvedBase = stripExtendedPrefix(resolvedBase)

	rel, err := filepath.Rel(resolvedBase, resolved)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", traversalErr("Path '" + path + "' resolves to '" + resolved + "' which is outside the allowed base directory '" + resolvedBase + "'")
	}
	return resolved, nil
}

func stripExtendedPrefix(p string) string {
	if strings.HasPrefix(p, `\\?\`) {
		return p[4:]
	}
	return p
}

// SafeRmtree removes path only if it resolves within baseDir.
func SafeRmtree(path, baseDir string) error {
	if _, err := EnsurePathWithin(path, baseDir); err != nil {
		return err
	}
	return os.RemoveAll(path)
}
