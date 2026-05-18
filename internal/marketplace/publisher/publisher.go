// Package publisher implements the MarketplacePublisher: update consumer repos
// with new marketplace package versions.
//
// The publisher reads the local marketplace.yml, computes a deterministic
// branch name and commit message, clones each consumer repo, updates its
// apm.yml, and pushes a feature branch.
//
// Design notes:
//   - Byte integrity: publisher NEVER regenerates marketplace.json; only copies it.
//   - Token redaction: stderr from git subprocesses is redacted.
//   - Atomic writes: state files and apm.yml updates use write-tmp + rename.
//   - Error isolation: failures in one target never abort other targets.
//
// Migrated from: src/apm_cli/marketplace/publisher.py
package publisher

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// -------------------------------------------------------------------
// Data types
// -------------------------------------------------------------------

// PublishStatus records the outcome for one consumer repo.
type PublishStatus string

const (
	StatusSuccess PublishStatus = "success"
	StatusSkipped PublishStatus = "skipped"
	StatusFailed  PublishStatus = "failed"
)

// ConsumerUpdate describes the update to push to a single consumer.
type ConsumerUpdate struct {
	Repo         string
	BranchName   string
	CommitMsg    string
	APMYMLPatch  string // new content for the consumer's apm.yml
	PackageName  string
	NewVersion   string
	OldVersion   string
}

// PublishResult is the outcome for one consumer repo.
type PublishResult struct {
	Repo    string
	Status  PublishStatus
	Branch  string
	Error   error
	Skipped bool
	Reason  string
}

// PublishReport summarises a publish run.
type PublishReport struct {
	Results   []PublishResult
	StartedAt time.Time
	Duration  time.Duration
	Errors    []error
}

// OK returns true when all results succeeded or were skipped.
func (r *PublishReport) OK() bool {
	for _, res := range r.Results {
		if res.Status == StatusFailed {
			return false
		}
	}
	return true
}

// PublishOptions controls a publish run.
type PublishOptions struct {
	Concurrency int
	DryRun      bool
	Force       bool
	Token       string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() PublishOptions {
	return PublishOptions{Concurrency: 4}
}

// -------------------------------------------------------------------
// Marketplace types (minimal)
// -------------------------------------------------------------------

// MarketplaceYML holds the parsed marketplace.yml for a source repo.
type MarketplaceYML struct {
	Name       string            `yaml:"name" json:"name"`
	Version    string            `yaml:"version" json:"version"`
	Consumers  []string          `yaml:"consumers" json:"consumers"`
	Packages   map[string]string `yaml:"packages" json:"packages"`
}

// -------------------------------------------------------------------
// MarketplacePublisher
// -------------------------------------------------------------------

// MarketplacePublisher pushes version bumps to consumer repositories.
type MarketplacePublisher struct {
	sourceDir string
	yml       *MarketplaceYML
	token     string
	mu        sync.Mutex
}

// New constructs a MarketplacePublisher for sourceDir.
func New(sourceDir, token string) (*MarketplacePublisher, error) {
	if sourceDir == "" {
		sourceDir = "."
	}
	abs, err := filepath.Abs(sourceDir)
	if err != nil {
		return nil, err
	}
	p := &MarketplacePublisher{sourceDir: abs, token: token}
	if err := p.loadYML(); err != nil {
		return nil, err
	}
	return p, nil
}

// loadYML reads marketplace.yml from sourceDir.
func (p *MarketplacePublisher) loadYML() error {
	path := filepath.Join(p.sourceDir, "marketplace.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("marketplace.yml not found at %s", path)
		}
		return err
	}
	// Minimal line-based YAML parse for name/version/consumers.
	yml := &MarketplaceYML{Packages: make(map[string]string)}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			yml.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		} else if strings.HasPrefix(line, "version:") {
			yml.Version = strings.TrimSpace(strings.TrimPrefix(line, "version:"))
		} else if strings.HasPrefix(line, "- ") && yml.Name != "" {
			yml.Consumers = append(yml.Consumers, strings.TrimPrefix(line, "- "))
		}
	}
	p.yml = yml
	return nil
}

// -------------------------------------------------------------------
// Publish
// -------------------------------------------------------------------

// Publish runs the publish loop for all consumers.
func (p *MarketplacePublisher) Publish(opts PublishOptions) (*PublishReport, error) {
	if p.yml == nil {
		return nil, errors.New("marketplace.yml not loaded")
	}
	if len(p.yml.Consumers) == 0 {
		return &PublishReport{StartedAt: time.Now()}, nil
	}

	t0 := time.Now()
	report := &PublishReport{StartedAt: t0}

	sem := make(chan struct{}, max(opts.Concurrency, 1))
	var wg sync.WaitGroup
	resultCh := make(chan PublishResult, len(p.yml.Consumers))

	for _, consumer := range p.yml.Consumers {
		wg.Add(1)
		sem <- struct{}{}
		go func(repo string) {
			defer wg.Done()
			defer func() { <-sem }()
			result := p.publishToConsumer(repo, opts)
			resultCh <- result
		}(consumer)
	}

	wg.Wait()
	close(resultCh)

	for r := range resultCh {
		report.Results = append(report.Results, r)
		if r.Error != nil {
			report.Errors = append(report.Errors, r.Error)
		}
	}
	sort.Slice(report.Results, func(i, j int) bool {
		return report.Results[i].Repo < report.Results[j].Repo
	})
	report.Duration = time.Since(t0)
	return report, nil
}

func (p *MarketplacePublisher) publishToConsumer(repo string, opts PublishOptions) PublishResult {
	result := PublishResult{Repo: repo, Status: StatusFailed}

	branch := p.branchName(repo)
	result.Branch = branch

	tmpDir, err := os.MkdirTemp("", "apm-publish-*")
	if err != nil {
		result.Error = err
		return result
	}
	defer os.RemoveAll(tmpDir)

	repoURL := p.buildRepoURL(repo)
	if err := p.gitClone(repoURL, tmpDir); err != nil {
		result.Error = fmt.Errorf("clone %s: %w", repo, err)
		return result
	}

	apmYMLPath := filepath.Join(tmpDir, "apm.yml")
	updated, oldVer, err := p.patchAPMYML(apmYMLPath)
	if err != nil {
		result.Error = fmt.Errorf("patch apm.yml: %w", err)
		return result
	}
	if !updated {
		result.Status = StatusSkipped
		result.Skipped = true
		result.Reason = "already up to date"
		return result
	}

	if opts.DryRun {
		result.Status = StatusSuccess
		result.Reason = "dry run"
		return result
	}

	commitMsg := fmt.Sprintf("chore: update %s to %s (was %s)\n\nPublished by APM marketplace publisher.",
		p.yml.Name, p.yml.Version, oldVer)

	if err := p.gitCommitAndPush(tmpDir, branch, commitMsg); err != nil {
		result.Error = fmt.Errorf("push to %s: %w", repo, err)
		return result
	}

	result.Status = StatusSuccess
	return result
}

func (p *MarketplacePublisher) branchName(repo string) string {
	h := sha256.Sum256([]byte(repo + p.yml.Name + p.yml.Version))
	return fmt.Sprintf("apm/marketplace-update-%s-%s-%x", p.yml.Name, p.yml.Version, h[:4])
}

func (p *MarketplacePublisher) buildRepoURL(repo string) string {
	if p.token != "" {
		return fmt.Sprintf("https://x-access-token:%s@github.com/%s.git", p.token, repo)
	}
	return fmt.Sprintf("https://github.com/%s.git", repo)
}

func (p *MarketplacePublisher) gitClone(url, dir string) error {
	cmd := exec.Command("git", "clone", "--depth=1", url, dir)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, redactToken(string(out), p.token))
	}
	return nil
}

func (p *MarketplacePublisher) patchAPMYML(path string) (updated bool, oldVer string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, "", nil
		}
		return false, "", err
	}

	// Look for the package name and update its version.
	re := regexp.MustCompile(`(?m)^(\s*` + regexp.QuoteMeta(p.yml.Name) + `\s*:\s*)(.+)$`)
	current := string(data)
	m := re.FindStringSubmatch(current)
	if m == nil {
		return false, "", nil // package not referenced
	}
	oldVer = strings.TrimSpace(m[2])
	if oldVer == p.yml.Version {
		return false, oldVer, nil // already at target version
	}

	patched := re.ReplaceAllString(current, "${1}"+p.yml.Version)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(patched), 0o644); err != nil {
		return false, oldVer, err
	}
	return true, oldVer, os.Rename(tmp, path)
}

func (p *MarketplacePublisher) gitCommitAndPush(dir, branch, msg string) error {
	run := func(args ...string) error {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git %s: %w: %s", args[0], err, redactToken(string(out), p.token))
		}
		return nil
	}

	if err := run("checkout", "-b", branch); err != nil {
		return err
	}
	if err := run("add", "apm.yml"); err != nil {
		return err
	}
	if err := run("commit", "-m", msg); err != nil {
		return err
	}
	return run("push", "origin", branch)
}

// -------------------------------------------------------------------
// Marketplace JSON copy (byte-integrity guarantee)
// -------------------------------------------------------------------

// CopyMarketplaceJSON copies marketplace.json verbatim from src to dst.
// It NEVER regenerates or modifies the file content.
func CopyMarketplaceJSON(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	tmp := dst + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	out.Close()
	return os.Rename(tmp, dst)
}

// -------------------------------------------------------------------
// State file helpers
// -------------------------------------------------------------------

// PublishState records the last publish run's outcomes for idempotency.
type PublishState struct {
	LastPublished time.Time             `json:"last_published"`
	Consumers     map[string]string     `json:"consumers"` // repo -> branch pushed
	Version       string                `json:"version"`
}

// LoadPublishState reads the publish state file.
func LoadPublishState(path string) (*PublishState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &PublishState{Consumers: make(map[string]string)}, nil
		}
		return nil, err
	}
	var state PublishState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	if state.Consumers == nil {
		state.Consumers = make(map[string]string)
	}
	return &state, nil
}

// SavePublishState writes the publish state atomically.
func SavePublishState(path string, state *PublishState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// -------------------------------------------------------------------
// Token redaction
// -------------------------------------------------------------------

func redactToken(s, token string) string {
	if token == "" {
		return s
	}
	return strings.ReplaceAll(s, token, "[REDACTED]")
}

// -------------------------------------------------------------------
// Semver helpers
// -------------------------------------------------------------------

// BumpPatch increments the patch component of a semver string.
func BumpPatch(version string) (string, error) {
	re := regexp.MustCompile(`^(v?)(\d+)\.(\d+)\.(\d+)(.*)$`)
	m := re.FindStringSubmatch(version)
	if m == nil {
		return "", fmt.Errorf("invalid semver: %q", version)
	}
	var patch int
	fmt.Sscanf(m[4], "%d", &patch)
	return fmt.Sprintf("%s%s.%s.%d%s", m[1], m[2], m[3], patch+1, m[5]), nil
}

// -------------------------------------------------------------------
// Tag rendering
// -------------------------------------------------------------------

// RenderTag replaces {version} in a tag pattern template.
func RenderTag(pattern, version string) string {
	return strings.ReplaceAll(pattern, "{version}", version)
}

// -------------------------------------------------------------------
// Report rendering
// -------------------------------------------------------------------

// RenderReport returns a human-readable summary of a publish report.
func RenderReport(r *PublishReport) string {
	if r == nil {
		return ""
	}
	var sb strings.Builder
	for _, res := range r.Results {
		switch res.Status {
		case StatusSuccess:
			sb.WriteString(fmt.Sprintf("[+] %s -> %s\n", res.Repo, res.Branch))
		case StatusSkipped:
			sb.WriteString(fmt.Sprintf("[i] %s skipped: %s\n", res.Repo, res.Reason))
		case StatusFailed:
			sb.WriteString(fmt.Sprintf("[x] %s failed: %v\n", res.Repo, res.Error))
		}
	}
	return sb.String()
}

// -------------------------------------------------------------------
// Go 1.21+ max helper (stdlib max was added in 1.21)
// -------------------------------------------------------------------

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
