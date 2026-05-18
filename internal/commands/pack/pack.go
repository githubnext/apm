// Package pack implements the "apm pack" and "apm unpack" commands.
//
// pack: produces distributable artifacts (bundle, .tar.gz archive, or
//   marketplace plugin manifest) from the project's apm.yml.
// unpack: extracts a previously-packed bundle.
//
// Migrated from: src/apm_cli/commands/pack.py
package pack

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Format selects the bundle output format.
type Format string

const (
	FormatPlugin Format = "plugin"
	FormatAPM    Format = "apm"
)

// PackOptions configures a pack run.
type PackOptions struct {
	ProjectRoot        string
	Format             Format
	Archive            bool
	OutputDir          string
	Offline            bool
	DryRun             bool
	MarketplaceOutput  string
	Verbose            bool
}

// PackResult records what was produced.
type PackResult struct {
	OutputPaths []string
	DryRun      bool
}

// Run executes the pack command.
func Run(opts PackOptions) (*PackResult, error) {
	if opts.OutputDir == "" {
		opts.OutputDir = filepath.Join(opts.ProjectRoot, "build")
	}
	if opts.Format == "" {
		opts.Format = FormatPlugin
	}

	if !opts.DryRun {
		if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
			return nil, fmt.Errorf("create output dir: %w", err)
		}
	}

	// Determine what the project contains.
	hasDeps := fileExists(filepath.Join(opts.ProjectRoot, "apm.yml"))
	if !hasDeps {
		return nil, fmt.Errorf("no apm.yml found in %s", opts.ProjectRoot)
	}

	var outputs []string

	// Build bundle.
	bundlePath, err := buildBundle(opts)
	if err != nil {
		return nil, fmt.Errorf("build bundle: %w", err)
	}
	if bundlePath != "" {
		outputs = append(outputs, bundlePath)
	}

	if opts.DryRun {
		fmt.Println("[i] Dry-run: no files written.")
		return &PackResult{DryRun: true}, nil
	}

	fmt.Printf("[+] Pack complete: %s\n", strings.Join(outputs, ", "))
	return &PackResult{OutputPaths: outputs}, nil
}

// buildBundle assembles the package contents and optionally archives them.
func buildBundle(opts PackOptions) (string, error) {
	projectName := filepath.Base(opts.ProjectRoot)
	if opts.DryRun {
		fmt.Printf("[i] Would write bundle for %s to %s\n", projectName, opts.OutputDir)
		return "", nil
	}

	var outputPath string
	if opts.Archive {
		outputPath = filepath.Join(opts.OutputDir, projectName+".tar.gz")
		if err := createTarGZ(opts.ProjectRoot, outputPath); err != nil {
			return "", err
		}
	} else {
		outputPath = filepath.Join(opts.OutputDir, projectName)
		if err := copyDir(opts.ProjectRoot, outputPath); err != nil {
			return "", err
		}
	}
	return outputPath, nil
}

// UnpackOptions configures an unpack run.
type UnpackOptions struct {
	BundlePath  string
	DestDir     string
	ProjectRoot string
	DryRun      bool
	Verbose     bool
}

// RunUnpack extracts a bundle.
func RunUnpack(opts UnpackOptions) error {
	if opts.BundlePath == "" {
		return fmt.Errorf("bundle path is required")
	}
	if opts.DestDir == "" {
		opts.DestDir = opts.ProjectRoot
	}
	if opts.DryRun {
		fmt.Printf("[i] Would unpack %s to %s\n", opts.BundlePath, opts.DestDir)
		return nil
	}

	if strings.HasSuffix(opts.BundlePath, ".tar.gz") || strings.HasSuffix(opts.BundlePath, ".tgz") {
		return extractTarGZ(opts.BundlePath, opts.DestDir)
	}
	// Directory bundle: copy
	return copyDir(opts.BundlePath, opts.DestDir)
}

// --- helpers ---

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createTarGZ(src, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		if rel == "." {
			return nil
		}
		info, err := d.Info()
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
		if d.IsDir() {
			return nil
		}
		sf, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sf.Close()
		_, err = io.Copy(tw, sf)
		return err
	})
}

func extractTarGZ(src, dest string) error {
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
		target := filepath.Join(dest, filepath.Clean(hdr.Name))
		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		df, err := os.Create(target)
		if err != nil {
			return err
		}
		if _, err := io.Copy(df, tr); err != nil {
			df.Close()
			return err
		}
		df.Close()
	}
	return nil
}

func copyDir(src, dest string) error {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		if rel == "." {
			return nil
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dest string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}
