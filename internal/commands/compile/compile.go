// Package compile implements the "apm compile" command.
//
// Compiles APM primitives (instructions, contexts, chatmodes) from .apm/
// directories into a single AGENTS.md constitution file and writes a
// build-ID-stamped output. Supports single-file mode, directory mode, and
// watch mode.
//
// Migrated from: src/apm_cli/commands/compile/cli.py
package compile

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// CompileOptions configures a single compilation run.
type CompileOptions struct {
	ProjectRoot string
	Output      string
	DryRun      bool
	Watch       bool
	Force       bool
	Strict      bool
	Verbose     bool
}

// CompileStats holds counters accumulated during compilation.
type CompileStats struct {
	Instructions int
	Contexts     int
	Chatmodes    int
	Primitives   int
	Warnings     []string
}

// CompileResult is returned by Compile.
type CompileResult struct {
	OutputPath    string
	ConstitutionHash string
	Status        string
	Stats         CompileStats
	DryRun        bool
}

// Compile discovers and compiles APM primitives in the project.
func Compile(opts CompileOptions) (*CompileResult, error) {
	projectRoot := opts.ProjectRoot
	if projectRoot == "" {
		var err error
		projectRoot, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getting cwd: %w", err)
		}
	}

	apmDir := filepath.Join(projectRoot, ".apm")
	if _, err := os.Stat(apmDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no .apm directory found at %s", projectRoot)
	}

	stats, sections, err := discoverPrimitives(apmDir, opts.Strict)
	if err != nil {
		return nil, fmt.Errorf("discovering primitives: %w", err)
	}

	constitution := buildConstitution(sections)
	hash := computeHash(constitution)

	outputPath := opts.Output
	if outputPath == "" {
		outputPath = filepath.Join(projectRoot, "AGENTS.md")
	}

	status := "unchanged"
	if opts.Force || !fileMatchesContent(outputPath, constitution) {
		status = "updated"
		if !opts.DryRun {
			if err := writeAtomic(outputPath, []byte(constitution)); err != nil {
				return nil, fmt.Errorf("writing %s: %w", outputPath, err)
			}
		}
	}

	return &CompileResult{
		OutputPath:       outputPath,
		ConstitutionHash: hash,
		Status:           status,
		Stats:            *stats,
		DryRun:           opts.DryRun,
	}, nil
}

// WatchOptions configures the watch mode.
type WatchOptions struct {
	CompileOptions
	Interval time.Duration
}

// Watch runs Compile in a loop, recompiling when .apm/ files change.
func Watch(opts WatchOptions, done <-chan struct{}) error {
	interval := opts.Interval
	if interval == 0 {
		interval = 500 * time.Millisecond
	}

	var lastHash string
	for {
		select {
		case <-done:
			return nil
		case <-time.After(interval):
		}

		result, err := Compile(opts.CompileOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[x] compile error: %v\n", err)
			continue
		}
		if result.ConstitutionHash != lastHash {
			lastHash = result.ConstitutionHash
			fmt.Printf("[*] recompiled: %s (hash %s)\n", result.OutputPath, result.ConstitutionHash[:8])
		}
	}
}

// PrimitiveSection holds the content and metadata for a discovered primitive.
type PrimitiveSection struct {
	Kind    string // "instruction", "context", "chatmode"
	Path    string
	Content string
	Title   string
}

// discoverPrimitives scans the .apm directory and returns accumulated stats
// and section content for the constitution.
func discoverPrimitives(apmDir string, strict bool) (*CompileStats, []PrimitiveSection, error) {
	var stats CompileStats
	var sections []PrimitiveSection

	err := filepath.WalkDir(apmDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		kind := ""
		switch {
		case strings.HasSuffix(path, ".instructions.md"):
			kind = "instruction"
		case strings.HasSuffix(path, ".context.md"):
			kind = "context"
		case strings.HasSuffix(path, ".chatmode.md"):
			kind = "chatmode"
		default:
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			if strict {
				return fmt.Errorf("reading %s: %w", path, err)
			}
			stats.Warnings = append(stats.Warnings, fmt.Sprintf("could not read %s: %v", path, err))
			return nil
		}

		title := extractTitle(string(content), filepath.Base(path))
		sections = append(sections, PrimitiveSection{
			Kind:    kind,
			Path:    path,
			Content: string(content),
			Title:   title,
		})

		switch kind {
		case "instruction":
			stats.Instructions++
		case "context":
			stats.Contexts++
		case "chatmode":
			stats.Chatmodes++
		}
		stats.Primitives++
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// Stable ordering: instructions first, then contexts, then chatmodes.
	sort.Slice(sections, func(i, j int) bool {
		kindOrder := map[string]int{"instruction": 0, "context": 1, "chatmode": 2}
		ki, kj := kindOrder[sections[i].Kind], kindOrder[sections[j].Kind]
		if ki != kj {
			return ki < kj
		}
		return sections[i].Path < sections[j].Path
	})

	return &stats, sections, nil
}

// buildConstitution concatenates the discovered sections into a single markdown document.
func buildConstitution(sections []PrimitiveSection) string {
	if len(sections) == 0 {
		return "# APM Constitution\n\n*(No primitives found.)*\n"
	}

	var sb strings.Builder
	sb.WriteString("# APM Constitution\n\n")
	sb.WriteString(fmt.Sprintf("*Generated by apm compile on %s.*\n\n", time.Now().UTC().Format("2006-01-02 15:04 UTC")))
	sb.WriteString("---\n\n")

	for _, s := range sections {
		sb.WriteString(fmt.Sprintf("<!-- primitive: %s type: %s -->\n\n", s.Title, s.Kind))
		sb.WriteString(s.Content)
		if !strings.HasSuffix(s.Content, "\n") {
			sb.WriteByte('\n')
		}
		sb.WriteString("\n---\n\n")
	}
	return sb.String()
}

// computeHash returns a short SHA-256 hex digest of content.
func computeHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", h[:8])
}

// fileMatchesContent returns true when the file at path has the same bytes as content.
func fileMatchesContent(path, content string) bool {
	existing, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return string(existing) == content
}

// writeAtomic writes data to path via a temp file + rename.
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".agents-tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

// extractTitle pulls the first heading from markdown content, falling back to filename.
func extractTitle(content, filename string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
