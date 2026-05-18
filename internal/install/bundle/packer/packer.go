// Package packer creates self-contained APM bundles from the resolved dependency tree.
//
// Migrated from src/apm_cli/bundle/packer.py
package packer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PackResult holds the result of a pack operation.
type PackResult struct {
	BundlePath      string
	Files           []string
	LockfileEnriched bool
	MappedCount     int
	PathMappings    map[string]string
}

// PackOptions configures a pack operation.
type PackOptions struct {
	// ProjectRoot is the root of the project containing apm.lock.yaml.
	ProjectRoot string
	// OutputDir is the directory where the bundle will be created.
	OutputDir string
	// Format is the bundle format: "apm" (default) or "plugin".
	Format string
	// Target is the target filter: "copilot", "claude", "all", or comma-separated list.
	// Empty means auto-detect.
	Target string
	// Archive creates a .tar.gz and removes the directory.
	Archive bool
	// DryRun resolves the file list but writes nothing to disk.
	DryRun bool
	// Force overwrites on collision.
	Force bool
}

// DeployedFile represents a file to be included in the bundle.
type DeployedFile struct {
	SourcePath string // absolute path on disk
	BundlePath string // relative path in bundle
}

// BundleDependency represents a dependency entry with its deployed files.
type BundleDependency struct {
	Name          string
	Version       string
	DeployedFiles []string
}

// PackBundle creates a self-contained bundle from installed APM dependencies.
// It reads deployed files from project_root and copies them into a bundle directory.
func PackBundle(opts PackOptions) (*PackResult, error) {
	if opts.Format == "" {
		opts.Format = "apm"
	}

	// Find and read lockfile
	lockfilePath := findLockfile(opts.ProjectRoot)
	if lockfilePath == "" {
		return nil, fmt.Errorf("apm.lock.yaml not found -- run 'apm install' first")
	}

	deps, err := readDeployedFiles(lockfilePath)
	if err != nil {
		return nil, fmt.Errorf("reading lockfile: %w", err)
	}

	// Resolve target
	target := opts.Target
	if target == "" {
		target = detectTarget(opts.ProjectRoot)
	}

	// Filter files by target with cross-target mapping
	type filteredDep struct {
		dep     BundleDependency
		files   []string
		mapping map[string]string
	}

	var allFiles []string
	seenFiles := map[string]bool{}
	allMappings := map[string]string{}
	var filtered []filteredDep

	for _, dep := range deps {
		files, mappings := filterFilesByTarget(dep.DeployedFiles, target)
		for f, orig := range mappings {
			allMappings[f] = orig
		}
		for _, f := range files {
			if !seenFiles[f] {
				seenFiles[f] = true
				allFiles = append(allFiles, f)
			}
		}
		filtered = append(filtered, filteredDep{dep: dep, files: files, mapping: mappings})
	}

	// Verify files exist on disk (skip local-content files)
	var missing []string
	for _, f := range allFiles {
		src := filepath.Join(opts.ProjectRoot, f)
		if _, err2 := os.Stat(src); os.IsNotExist(err2) {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 && !opts.Force {
		return nil, fmt.Errorf("bundle verification failed -- %d files missing from disk", len(missing))
	}

	if opts.DryRun {
		return &PackResult{
			BundlePath:      opts.OutputDir,
			Files:           allFiles,
			LockfileEnriched: true,
			MappedCount:     len(allMappings),
			PathMappings:    allMappings,
		}, nil
	}

	// Create bundle directory
	bundleDir := filepath.Join(opts.OutputDir, "bundle")
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating bundle dir: %w", err)
	}

	// Copy files into bundle
	for _, f := range allFiles {
		// Map bundle path back to disk path
		diskRel := f
		if orig, ok := allMappings[f]; ok {
			diskRel = orig
		}
		src := filepath.Join(opts.ProjectRoot, diskRel)
		dst := filepath.Join(bundleDir, f)

		fi, err2 := os.Stat(src)
		if err2 != nil {
			continue // skip missing files (already checked above with non-force)
		}

		if err3 := os.MkdirAll(filepath.Dir(dst), 0o755); err3 != nil {
			return nil, err3
		}

		if fi.IsDir() {
			if err3 := copyDirContents(src, dst); err3 != nil {
				return nil, err3
			}
		} else {
			if err3 := copyFile(src, dst); err3 != nil {
				return nil, err3
			}
		}
	}

	// Copy the lockfile into the bundle
	lockfileDst := filepath.Join(bundleDir, "apm.lock.yaml")
	_ = copyFile(lockfilePath, lockfileDst)

	bundlePath := bundleDir
	if opts.Archive {
		archivePath := bundleDir + ".tar.gz"
		if err := createTarGz(bundleDir, archivePath); err != nil {
			return nil, fmt.Errorf("creating archive: %w", err)
		}
		os.RemoveAll(bundleDir)
		bundlePath = archivePath
	}

	return &PackResult{
		BundlePath:      bundlePath,
		Files:           allFiles,
		LockfileEnriched: true,
		MappedCount:     len(allMappings),
		PathMappings:    allMappings,
	}, nil
}

// findLockfile finds the lockfile in project_root.
func findLockfile(projectRoot string) string {
	candidates := []string{
		filepath.Join(projectRoot, "apm.lock.yaml"),
		filepath.Join(projectRoot, "apm.lock"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// detectTarget auto-detects the target from the project structure.
func detectTarget(projectRoot string) string {
	dirs := map[string]string{
		".github":   "copilot",
		".claude":   "claude",
		".cursor":   "cursor",
		".windsurf": "windsurf",
		".agents":   "agent-skills",
	}
	for dir, target := range dirs {
		if _, err := os.Stat(filepath.Join(projectRoot, dir)); err == nil {
			return target
		}
	}
	return "all"
}

// readDeployedFiles parses deployed_files from a lockfile.
func readDeployedFiles(lockfilePath string) ([]BundleDependency, error) {
	data, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, err
	}

	var deps []BundleDependency
	var current *BundleDependency
	inDeps := false
	inDeployedFiles := false

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			inDeps = strings.HasPrefix(trimmed, "dependencies:")
			inDeployedFiles = false
			if current != nil {
				deps = append(deps, *current)
				current = nil
			}
			continue
		}

		if !inDeps {
			continue
		}

		if strings.HasPrefix(line, "  - ") {
			if current != nil {
				deps = append(deps, *current)
			}
			current = &BundleDependency{}
			inDeployedFiles = false
			rest := strings.TrimPrefix(line, "  - ")
			if strings.HasPrefix(rest, "name:") {
				current.Name = strings.TrimSpace(strings.TrimPrefix(rest, "name:"))
			}
			continue
		}

		if current == nil {
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "deployed_files:") {
			inDeployedFiles = true
		} else if inDeployedFiles && strings.HasPrefix(trimmed, "- ") {
			f := strings.TrimPrefix(trimmed, "- ")
			f = strings.Trim(f, `"'`)
			current.DeployedFiles = append(current.DeployedFiles, f)
		} else if inDeployedFiles {
			inDeployedFiles = false
		}
	}
	if current != nil {
		deps = append(deps, *current)
	}
	return deps, nil
}

// knownTargetPrefixes maps target names to effective pack prefixes.
var knownTargetPrefixes = map[string][]string{
	"copilot":      {".github/"},
	"vscode":       {".github/"},
	"claude":       {".claude/"},
	"cursor":       {".cursor/"},
	"opencode":     {".opencode/"},
	"codex":        {".codex/", ".agents/"},
	"windsurf":     {".windsurf/"},
	"agent-skills": {".agents/"},
}

var crossTargetMaps = map[string]map[string]string{
	"claude": {
		".github/skills/": ".claude/skills/",
		".github/agents/": ".claude/agents/",
	},
	"vscode": {
		".claude/skills/": ".github/skills/",
		".claude/agents/": ".github/agents/",
	},
	"copilot": {
		".claude/skills/": ".github/skills/",
		".claude/agents/": ".github/agents/",
	},
	"cursor": {
		".github/skills/": ".cursor/skills/",
		".github/agents/": ".cursor/agents/",
	},
	"opencode": {
		".github/skills/": ".opencode/skills/",
		".github/agents/": ".opencode/agents/",
	},
	"codex": {
		".github/skills/": ".agents/skills/",
		".github/agents/": ".codex/agents/",
	},
	"windsurf": {
		".github/skills/": ".windsurf/skills/",
		".github/agents/": ".windsurf/skills/",
	},
}

func filterFilesByTarget(files []string, target string) ([]string, map[string]string) {
	targets := strings.Split(target, ",")
	var prefixes []string
	seen := map[string]bool{}
	crossMap := map[string]string{}

	for _, t := range targets {
		t = strings.TrimSpace(t)
		ps := knownTargetPrefixes[t]
		if t == "all" || len(ps) == 0 {
			// union all
			for _, tps := range knownTargetPrefixes {
				for _, p := range tps {
					if !seen[p] {
						seen[p] = true
						prefixes = append(prefixes, p)
					}
				}
			}
		} else {
			for _, p := range ps {
				if !seen[p] {
					seen[p] = true
					prefixes = append(prefixes, p)
				}
			}
		}
		for k, v := range crossTargetMaps[t] {
			crossMap[k] = v
		}
	}

	var direct []string
	directSet := map[string]bool{}
	for _, f := range files {
		for _, p := range prefixes {
			if strings.HasPrefix(f, p) {
				direct = append(direct, f)
				directSet[f] = true
				break
			}
		}
	}

	mappings := map[string]string{}
	for _, f := range files {
		if directSet[f] {
			continue
		}
		for src, dst := range crossMap {
			if strings.HasPrefix(f, src) {
				mapped := dst + f[len(src):]
				if !directSet[mapped] {
					direct = append(direct, mapped)
					directSet[mapped] = true
					mappings[mapped] = f
				}
				break
			}
		}
	}
	return direct, mappings
}

// createTarGz creates a .tar.gz archive from a directory.
func createTarGz(srcDir, archivePath string) error {
	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = rel
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tw, file)
			return err
		}
		return nil
	})
}

// copyFile copies src to dst.
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

// copyDirContents recursively copies contents of src into dst.
func copyDirContents(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		dest := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(dest, info.Mode())
		}
		return copyFile(path, dest)
	})
}
