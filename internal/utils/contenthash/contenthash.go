// Package contenthash provides deterministic SHA-256 content hashing for
// package integrity verification.
package contenthash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const (
	// MarkerFilename is the cache-pin marker excluded from package hashes.
	MarkerFilename = ".apm-pin"
)

var excludedDirs = map[string]bool{
	".git":        true,
	"__pycache__": true,
}

// emptyHash is the well-known hash for an empty or missing package.
var emptyHash = "sha256:" + func() string {
	h := sha256.Sum256([]byte{})
	return fmt.Sprintf("%x", h)
}()

// ComputePackageHash computes a deterministic SHA-256 hash of a package's
// file tree. The hash is computed over sorted file paths and their contents,
// making it independent of filesystem ordering and metadata.
//
// Returns a hash string in format "sha256:<hex_digest>".
func ComputePackageHash(packagePath string) (string, error) {
	info, err := os.Lstat(packagePath)
	if err != nil || !info.IsDir() {
		return emptyHash, nil
	}

	var relFiles []string
	err = filepath.WalkDir(packagePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip symlinks
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		rel, relErr := filepath.Rel(packagePath, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}
		// Skip excluded directories
		parts := splitPath(rel)
		for _, part := range parts {
			if excludedDirs[part] {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if d.IsDir() {
			return nil
		}
		// Exclude root-level marker files
		if len(parts) == 1 && parts[0] == MarkerFilename {
			return nil
		}
		relFiles = append(relFiles, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("contenthash: walking %s: %w", packagePath, err)
	}

	if len(relFiles) == 0 {
		return emptyHash, nil
	}

	sort.Strings(relFiles)

	h := sha256.New()
	for _, rel := range relFiles {
		h.Write([]byte(rel))
		f, openErr := os.Open(filepath.Join(packagePath, filepath.FromSlash(rel)))
		if openErr != nil {
			return "", fmt.Errorf("contenthash: opening %s: %w", rel, openErr)
		}
		_, copyErr := io.Copy(h, f)
		f.Close()
		if copyErr != nil {
			return "", fmt.Errorf("contenthash: reading %s: %w", rel, copyErr)
		}
	}

	return fmt.Sprintf("sha256:%x", h.Sum(nil)), nil
}

// ComputeFileHash computes SHA-256 of a single file's contents.
// Returns "sha256:<hex_digest>". Returns the empty-content hash when the
// path does not exist or is not a regular file.
func ComputeFileHash(filePath string) (string, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return emptyHash, nil
	}
	if !info.Mode().IsRegular() {
		return emptyHash, nil
	}
	f, err := os.Open(filePath)
	if err != nil {
		return emptyHash, nil
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", fmt.Errorf("contenthash: reading %s: %w", filePath, err)
	}
	return fmt.Sprintf("sha256:%x", h.Sum(nil)), nil
}

// VerifyPackageHash verifies a package's content matches the expected hash.
// Returns true if hash matches.
func VerifyPackageHash(packagePath, expectedHash string) (bool, error) {
	actual, err := ComputePackageHash(packagePath)
	if err != nil {
		return false, err
	}
	return actual == expectedHash, nil
}

// splitPath splits a slash-separated relative path into its components.
func splitPath(p string) []string {
	s := filepath.ToSlash(p)
	var parts []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '/' {
			if seg := s[start:i]; seg != "" && seg != "." {
				parts = append(parts, seg)
			}
			start = i + 1
		}
	}
	return parts
}
