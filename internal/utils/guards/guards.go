// Package guards provides a read-only project-tree guard for drift detection.
//
// When apm audit runs the install pipeline against a scratch directory to
// compute drift, the working tree must remain untouched. ReadOnlyProjectGuard
// takes a stat snapshot of every protected path on entry and asserts no
// mutation occurred on exit. Any divergence returns a ProtectedPathMutationError.
//
// This is a defense-in-depth check: the primary mechanism is redirecting all
// writes via project_root=scratch_root. The guard catches accidental
// direct-path writes that bypass the redirection.
package guards

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ProtectedPathMutationError is returned when a path under guard was mutated.
type ProtectedPathMutationError struct {
	Violations []string
}

func (e *ProtectedPathMutationError) Error() string {
	return "Drift replay mutated protected project paths:\n  - " +
		strings.Join(e.Violations, "\n  - ")
}

type fileInfo struct {
	mtimeNs int64
	size    int64
	exists  bool
}

func statFile(path string) fileInfo {
	fi, err := os.Stat(path)
	if err != nil {
		return fileInfo{exists: false}
	}
	return fileInfo{
		mtimeNs: fi.ModTime().UnixNano(),
		size:    fi.Size(),
		exists:  true,
	}
}

// walkProtected enumerates every regular file under each root (recursive).
// Missing roots are silently dropped. Symlinks are not followed.
func walkProtected(roots []string) []string {
	var files []string
	for _, root := range roots {
		fi, err := os.Lstat(root)
		if err != nil {
			continue
		}
		if fi.Mode().IsRegular() {
			files = append(files, root)
			continue
		}
		if fi.IsDir() {
			_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.Type()&os.ModeSymlink != 0 {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				if d.Type().IsRegular() {
					files = append(files, path)
				}
				return nil
			})
		}
	}
	return files
}

// ReadOnlyProjectGuard snapshots protected paths and asserts no mutation.
//
// Usage:
//
//	g := NewReadOnlyProjectGuard(projectRoot, []string{".apm", "apm.lock.yaml", ".github"})
//	if err := g.Enter(); err != nil { ... }
//	runReplay(...)
//	if err := g.Exit(nil); err != nil { ... }
type ReadOnlyProjectGuard struct {
	projectRoot      string
	protectedRoots   []string
	snapshot         map[string]fileInfo
}

// NewReadOnlyProjectGuard creates a new guard.
func NewReadOnlyProjectGuard(projectRoot string, protectedSubpaths []string) *ReadOnlyProjectGuard {
	abs, _ := filepath.Abs(projectRoot)
	roots := make([]string, len(protectedSubpaths))
	for i, sp := range protectedSubpaths {
		roots[i] = filepath.Join(abs, sp)
	}
	return &ReadOnlyProjectGuard{
		projectRoot:    abs,
		protectedRoots: roots,
		snapshot:       make(map[string]fileInfo),
	}
}

// Enter takes the initial snapshot of protected paths.
func (g *ReadOnlyProjectGuard) Enter() error {
	files := walkProtected(g.protectedRoots)
	for _, f := range files {
		g.snapshot[f] = statFile(f)
	}
	return nil
}

// Exit checks for mutations. Pass the original error (if any) so that
// ProtectedPathMutationError is only surfaced when no other error is
// propagating (mirrors Python's __exit__ exc_type handling).
func (g *ReadOnlyProjectGuard) Exit(origErr error) error {
	currentFiles := walkProtected(g.protectedRoots)
	currentSet := make(map[string]struct{}, len(currentFiles))
	for _, f := range currentFiles {
		currentSet[f] = struct{}{}
	}

	var violations []string

	// Newly-appeared files under protected roots are violations.
	snapshotSet := make(map[string]struct{}, len(g.snapshot))
	for path := range g.snapshot {
		snapshotSet[path] = struct{}{}
	}
	for path := range currentSet {
		if _, seen := snapshotSet[path]; !seen {
			violations = append(violations, fmt.Sprintf("created: %s", path))
		}
	}

	// Snapshotted files that vanished or changed.
	for path, prev := range g.snapshot {
		cur := statFile(path)
		if !prev.exists && !cur.exists {
			continue // missing -> still missing: fine
		}
		if !prev.exists && cur.exists {
			violations = append(violations, fmt.Sprintf("created: %s", path))
		} else if prev.exists && !cur.exists {
			violations = append(violations, fmt.Sprintf("deleted: %s", path))
		} else if prev.mtimeNs != cur.mtimeNs || prev.size != cur.size {
			violations = append(violations, fmt.Sprintf("modified: %s", path))
		}
	}

	if len(violations) > 0 && origErr == nil {
		sort.Strings(violations)
		return &ProtectedPathMutationError{Violations: violations}
	}
	return nil
}
