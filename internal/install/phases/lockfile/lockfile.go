// Package lockfile assembles and persists the apm.lock.yaml from install
// artefacts. Mirrors src/apm_cli/install/phases/lockfile.py.
package lockfile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"
)

// DeployedFileHash computes the SHA-256 hash of a single deployed file.
// Returns "sha256:<hex>" or empty string on error.
func DeployedFileHash(absPath string) string {
	f, err := os.Open(absPath)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

// ComputeDeployedHashes hashes currently on-disk deployed files for provenance.
// projectRoot is the absolute path to the project directory.
// relPaths is a slice of paths relative to projectRoot.
// Returns map[relPath]"sha256:<hex>"; symlinks and unreadable paths are omitted.
func ComputeDeployedHashes(projectRoot string, relPaths []string) map[string]string {
	out := make(map[string]string, len(relPaths))
	for _, rel := range relPaths {
		if rel == "" {
			continue
		}
		full := projectRoot + "/" + rel
		info, err := os.Lstat(full)
		if err != nil {
			continue
		}
		// Skip symlinks and non-regular files.
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			continue
		}
		if h := DeployedFileHash(full); h != "" {
			out[rel] = h
		}
	}
	return out
}

// LockfileEntry holds the minimum metadata for one locked dependency as
// needed by the LockfileBuilder logic.
type LockfileEntry struct {
	DepKey          string
	RepoURL         string
	DeployedFiles   []string
	DeployedHashes  map[string]string
	ContentHash     string
	PackageType     string
	ResolvedRef     string
	ResolvedCommit  string
	SkillSubset     []string
	MarketplaceProvenance map[string]string
}

// WriteIfChanged writes newContent to path only when the on-disk content
// differs, to avoid unnecessary mtime churn.
func WriteIfChanged(path string, newContent []byte) (changed bool, err error) {
	existing, rerr := os.ReadFile(path)
	if rerr == nil && string(existing) == string(newContent) {
		return false, nil
	}
	tmp, err := os.CreateTemp("", "apm-lock-*")
	if err != nil {
		return false, fmt.Errorf("lockfile temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err = tmp.Write(newContent); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return false, fmt.Errorf("lockfile write: %w", err)
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmpName)
		return false, fmt.Errorf("lockfile close: %w", err)
	}
	if err = os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return false, fmt.Errorf("lockfile rename: %w", err)
	}
	return true, nil
}

// SortedDeployedFiles returns a deterministically sorted copy of the
// deployed files list for lockfile serialisation.
func SortedDeployedFiles(files []string) []string {
	cp := make([]string, len(files))
	copy(cp, files)
	sort.Strings(cp)
	return cp
}
