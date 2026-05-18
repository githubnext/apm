// Package unpacker extracts and verifies APM bundles.
//
// Migrated from src/apm_cli/bundle/unpacker.py
package unpacker

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UnpackResult holds the result of an unpack operation.
type UnpackResult struct {
	ExtractedDir     string
	Files            []string
	Verified         bool
	DependencyFiles  map[string][]string
	SkippedCount     int
	SecurityWarnings int
	SecurityCritical int
	PackMeta         map[string]interface{}
}

// LockEntry represents a single dependency entry from a bundle lockfile.
type LockEntry struct {
	Name          string
	Version       string
	DeployedFiles []string
}

// BundleLockfile holds parsed bundle lockfile data.
type BundleLockfile struct {
	Dependencies []LockEntry
	PackMeta     map[string]interface{}
	RawData      map[string]interface{}
}

// ParseBundleLockfile parses an apm.lock.yaml or legacy apm.lock file.
func ParseBundleLockfile(lockfilePath string) (*BundleLockfile, error) {
	data, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, fmt.Errorf("reading lockfile: %w", err)
	}

	lf := &BundleLockfile{
		PackMeta: map[string]interface{}{},
		RawData:  map[string]interface{}{},
	}

	// Simple YAML parser for the fields we need: dependencies[].deployed_files and pack:
	var currentDep *LockEntry
	inDependencies := false
	inDeployedFiles := false
	inPack := false
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Top-level section detection
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			inDependencies = strings.HasPrefix(trimmed, "dependencies:")
			inPack = strings.HasPrefix(trimmed, "pack:")
			inDeployedFiles = false
			if inPack {
				currentDep = nil
			}
			continue
		}

		if inPack {
			// Parse pack: sub-fields
			if strings.HasPrefix(strings.TrimSpace(line), "format:") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					lf.PackMeta["format"] = strings.TrimSpace(parts[1])
				}
			} else if strings.HasPrefix(strings.TrimSpace(line), "target:") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					lf.PackMeta["target"] = strings.TrimSpace(parts[1])
				}
			}
			continue
		}

		if !inDependencies {
			continue
		}

		// Detect new dependency entry (2-space indent starting with -)
		if strings.HasPrefix(line, "  - ") || (strings.HasPrefix(line, "- ") && !strings.HasPrefix(line, "  ")) {
			if currentDep != nil {
				lf.Dependencies = append(lf.Dependencies, *currentDep)
			}
			currentDep = &LockEntry{}
			inDeployedFiles = false
			rest := strings.TrimPrefix(strings.TrimPrefix(line, "  - "), "- ")
			if strings.HasPrefix(rest, "name:") {
				currentDep.Name = strings.TrimSpace(strings.TrimPrefix(rest, "name:"))
			}
			continue
		}

		if currentDep == nil {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		if indent >= 4 && strings.HasPrefix(trimmed, "name:") {
			currentDep.Name = strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
			inDeployedFiles = false
		} else if indent >= 4 && strings.HasPrefix(trimmed, "version:") {
			currentDep.Version = strings.TrimSpace(strings.TrimPrefix(trimmed, "version:"))
			inDeployedFiles = false
		} else if indent >= 4 && strings.HasPrefix(trimmed, "deployed_files:") {
			inDeployedFiles = true
		} else if inDeployedFiles && strings.HasPrefix(trimmed, "- ") {
			f := strings.TrimPrefix(trimmed, "- ")
			f = strings.Trim(f, `"'`)
			currentDep.DeployedFiles = append(currentDep.DeployedFiles, f)
		} else if inDeployedFiles && indent < 6 {
			inDeployedFiles = false
		}
	}

	if currentDep != nil {
		lf.Dependencies = append(lf.Dependencies, *currentDep)
	}

	return lf, nil
}

// UnpackBundle extracts and applies an APM bundle to a project directory.
func UnpackBundle(bundlePath, outputDir string, skipVerify, dryRun bool) (*UnpackResult, error) {
	sourceDir, tempDir, err := prepareSourceDir(bundlePath)
	if err != nil {
		return nil, err
	}
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}

	// Find lockfile
	lockfilePath := filepath.Join(sourceDir, "apm.lock.yaml")
	if _, err2 := os.Stat(lockfilePath); os.IsNotExist(err2) {
		legacyPath := filepath.Join(sourceDir, "apm.lock")
		if _, err3 := os.Stat(legacyPath); err3 == nil {
			lockfilePath = legacyPath
		}
	}

	lf, err := ParseBundleLockfile(lockfilePath)
	if err != nil {
		return nil, fmt.Errorf("lockfile missing from bundle: %w", err)
	}

	// Collect deployed_files per dependency
	depFileMap := map[string][]string{}
	seen := map[string]bool{}
	var uniqueFiles []string

	for _, dep := range lf.Dependencies {
		key := dep.Name
		if dep.Version != "" {
			key = dep.Name + "@" + dep.Version
		}
		var depFiles []string
		for _, f := range dep.DeployedFiles {
			depFiles = append(depFiles, f)
			if !seen[f] {
				seen[f] = true
				uniqueFiles = append(uniqueFiles, f)
			}
		}
		if len(depFiles) > 0 {
			depFileMap[key] = depFiles
		}
	}

	// Verify completeness
	verified := true
	if !skipVerify {
		var missing []string
		for _, f := range uniqueFiles {
			if _, err2 := os.Stat(filepath.Join(sourceDir, f)); os.IsNotExist(err2) {
				missing = append(missing, f)
			}
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("bundle verification failed -- missing files: %s",
				strings.Join(missing, ", "))
		}
	} else {
		verified = false
	}

	if dryRun {
		return &UnpackResult{
			ExtractedDir:    bundlePath,
			Files:           uniqueFiles,
			Verified:        verified,
			DependencyFiles: depFileMap,
			PackMeta:        lf.PackMeta,
		}, nil
	}

	// Copy files to output_dir (additive, no deletes)
	skipped := 0
	outputAbs, _ := filepath.Abs(outputDir)

	for _, relPath := range uniqueFiles {
		// Guard against path traversal
		p := filepath.Clean(relPath)
		if filepath.IsAbs(p) || strings.Contains(p, "..") {
			return nil, fmt.Errorf("refusing unsafe path from bundle lockfile: %q", relPath)
		}

		dest := filepath.Join(outputDir, relPath)
		destAbs, _ := filepath.Abs(dest)
		if !strings.HasPrefix(destAbs, outputAbs+string(os.PathSeparator)) && destAbs != outputAbs {
			return nil, fmt.Errorf("refusing path escaping output directory: %q", relPath)
		}

		src := filepath.Join(sourceDir, relPath)
		fi, err2 := os.Lstat(src)
		if err2 != nil {
			skipped++
			continue
		}

		// Skip symlinks
		if fi.Mode()&os.ModeSymlink != 0 {
			skipped++
			continue
		}

		if fi.IsDir() {
			if err3 := copyDir(src, dest); err3 != nil {
				return nil, err3
			}
		} else {
			if err3 := os.MkdirAll(filepath.Dir(dest), 0o755); err3 != nil {
				return nil, err3
			}
			if err3 := copyFile(src, dest); err3 != nil {
				return nil, err3
			}
		}
	}

	return &UnpackResult{
		ExtractedDir:    bundlePath,
		Files:           uniqueFiles,
		Verified:        verified,
		DependencyFiles: depFileMap,
		SkippedCount:    skipped,
		PackMeta:        lf.PackMeta,
	}, nil
}

// prepareSourceDir returns the source directory for a bundle.
// For .tar.gz archives, it extracts to a temp dir and returns the inner dir.
func prepareSourceDir(bundlePath string) (sourceDir, tempDir string, err error) {
	fi, err := os.Stat(bundlePath)
	if err != nil {
		return "", "", fmt.Errorf("bundle not found: %w", err)
	}

	if fi.IsDir() {
		return bundlePath, "", nil
	}

	if !strings.HasSuffix(bundlePath, ".tar.gz") {
		return "", "", fmt.Errorf("unsupported bundle format: %s", bundlePath)
	}

	tmp, err := os.MkdirTemp("", "apm-unpack-")
	if err != nil {
		return "", "", fmt.Errorf("creating temp dir: %w", err)
	}

	if err := extractTarGz(bundlePath, tmp); err != nil {
		os.RemoveAll(tmp)
		return "", "", err
	}

	// Locate inner directory
	entries, err := os.ReadDir(tmp)
	if err != nil {
		os.RemoveAll(tmp)
		return "", "", err
	}

	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(tmp, entries[0].Name()), tmp, nil
	}
	return tmp, tmp, nil
}

// extractTarGz extracts a .tar.gz archive to destDir.
func extractTarGz(src, destDir string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Security: reject path traversal and symlinks
		if filepath.IsAbs(hdr.Name) || strings.Contains(hdr.Name, "..") {
			return fmt.Errorf("refusing path-traversal entry: %s", hdr.Name)
		}
		if hdr.Typeflag == tar.TypeSymlink || hdr.Typeflag == tar.TypeLink {
			return fmt.Errorf("refusing symlink/hardlink entry: %s", hdr.Name)
		}

		dest := filepath.Join(destDir, hdr.Name)
		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return err
		}
		out.Close()
	}
	return nil
}

// copyFile copies src to dst (no symlink follow).
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// copyDir recursively copies src directory to dst.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(dest, info.Mode())
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil // skip symlinks
		}
		return copyFile(path, dest)
	})
}
