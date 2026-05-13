// Package paths provides cross-platform path utilities for APM CLI.
// Migrated from src/apm_cli/utils/paths.py
package paths

import (
	"path/filepath"
	"strings"
)

// PortableRelpath returns a forward-slash relative path, resolving both
// sides first. When path is not under base (or resolution fails), falls
// back to an absolute POSIX path.
func PortableRelpath(path, base string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return toForwardSlash(path)
	}
	absBase, err := filepath.Abs(base)
	if err != nil {
		return toForwardSlash(absPath)
	}
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return toForwardSlash(absPath)
	}
	return toForwardSlash(rel)
}

func toForwardSlash(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}
