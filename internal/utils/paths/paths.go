// Package paths provides cross-platform path utilities for APM CLI.
// Mirrors src/apm_cli/utils/paths.py (portable_relpath function).
package paths

import (
	"path/filepath"
	"strings"
)

// PortableRelpath returns a forward-slash relative path from base to path,
// resolving both sides first. When path is not under base (or resolution
// fails), falls back to the resolved absolute path.
//
// Mirrors Python's portable_relpath() from utils/paths.py.
func PortableRelpath(path, base string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return toSlash(path)
	}
	absBase, err := filepath.Abs(base)
	if err != nil {
		return toSlash(absPath)
	}
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return toSlash(absPath)
	}
	return toSlash(rel)
}

func toSlash(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}
