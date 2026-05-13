// Package contenthash provides deterministic SHA-256 content hashing for package integrity verification.
package contenthash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// MarkerFilename is the cache-pin marker excluded from content hashes.
const MarkerFilename = ".apm-pin"

// emptyHash is returned for empty or missing packages.
var emptyHash = "sha256:" + fmt.Sprintf("%x", sha256.Sum256([]byte{}))

// excludedDirs are not relevant to package content.
var excludedDirs = map[string]bool{
	".git":        true,
	"__pycache__": true,
}

// ComputePackageHash computes a deterministic SHA-256 hash of a package's file tree.
//
// The hash is computed over sorted file paths and their contents, making it
// independent of filesystem ordering and metadata.
func ComputePackageHash(packagePath string) (string, error) {
	info, err := os.Stat(packagePath)
	if err != nil || !info.IsDir() {
		return emptyHash, nil
	}

	type fileEntry struct {
		rel  string
		full string
	}
	var files []fileEntry

	err = filepath.WalkDir(packagePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// Skip symlinks
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		rel, relErr := filepath.Rel(packagePath, path)
		if relErr != nil {
			return nil
		}
		// Skip excluded directories
		parts := splitPath(rel)
		for _, part := range parts[:max(len(parts)-1, 0)] {
			if excludedDirs[part] {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if d.IsDir() {
			if excludedDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		// Skip excluded root files
		if len(parts) == 1 && parts[0] == MarkerFilename {
			return nil
		}
		files = append(files, fileEntry{rel: filepath.ToSlash(rel), full: path})
		return nil
	})
	if err != nil {
		return emptyHash, err
	}

	if len(files) == 0 {
		return emptyHash, nil
	}

	sort.Slice(files, func(i, j int) bool { return files[i].rel < files[j].rel })

	h := sha256.New()
	for _, f := range files {
		h.Write([]byte(f.rel))
		data, readErr := os.ReadFile(f.full)
		if readErr != nil {
			return emptyHash, readErr
		}
		h.Write(data)
	}

	return fmt.Sprintf("sha256:%x", h.Sum(nil)), nil
}

// ComputeFileHash computes SHA-256 of a single file's contents.
func ComputeFileHash(filePath string) (string, error) {
	info, err := os.Lstat(filePath)
	if err != nil || !info.Mode().IsRegular() {
		return emptyHash, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return emptyHash, nil
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return emptyHash, err
	}
	return fmt.Sprintf("sha256:%x", h.Sum(nil)), nil
}

// VerifyPackageHash verifies a package's content matches the expected hash.
func VerifyPackageHash(packagePath, expectedHash string) (bool, error) {
	actual, err := ComputePackageHash(packagePath)
	if err != nil {
		return false, err
	}
	return actual == expectedHash, nil
}

func splitPath(p string) []string {
	return filepath.SplitList(filepath.ToSlash(p))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
