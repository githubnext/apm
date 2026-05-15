// Package gitcache implements a persistent content-addressable git cache.
//
// Two-tier structure:
//   - git/db_v1/<shard>/    -- bare git repositories
//   - git/checkouts_v1/<shard>/<sha>/  -- per-SHA working copies
//
// Cache keys are derived from normalized repository URLs.
// Checkouts are keyed by resolved SHA, never by mutable ref strings.
package gitcache

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/githubnext/apm/internal/cache/cachepaths"
	"github.com/githubnext/apm/internal/cache/integrity"
	"github.com/githubnext/apm/internal/cache/locking"
	"github.com/githubnext/apm/internal/cache/urlnormalize"
)

// getGitDBPath returns the git bare-repo database directory.
func getGitDBPath(cacheRoot string) string {
	return filepath.Join(cacheRoot, cachepaths.GitDBBucket)
}

// getGitCheckoutsPath returns the git working-copy checkouts directory.
func getGitCheckoutsPath(cacheRoot string) string {
	return filepath.Join(cacheRoot, cachepaths.GitCheckoutsBucket)
}

// fullSHARe matches 40-char hex SHA strings.
var fullSHARe = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)

// CacheStats holds aggregate statistics about the git cache.
type CacheStats struct {
	DBCount        int
	CheckoutCount  int
	TotalSizeBytes int64
}

// GitCache is a content-addressable git cache with integrity verification.
type GitCache struct {
	cacheRoot     string
	dbRoot        string
	checkoutsRoot string
	refresh       bool
}

// New creates a GitCache rooted at cacheRoot. If refresh is true,
// integrity is revalidated on every access.
func New(cacheRoot string, refresh bool) (*GitCache, error) {
	dbRoot := getGitDBPath(cacheRoot)
	checkoutsRoot := getGitCheckoutsPath(cacheRoot)

	for _, dir := range []string{dbRoot, checkoutsRoot} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, fmt.Errorf("gitcache: mkdir %s: %w", dir, err)
		}
	}

	locking.CleanupIncomplete(dbRoot)
	locking.CleanupIncomplete(checkoutsRoot)

	return &GitCache{
		cacheRoot:     cacheRoot,
		dbRoot:        dbRoot,
		checkoutsRoot: checkoutsRoot,
		refresh:       refresh,
	}, nil
}

// GetCheckout returns the path to a cached working-tree checkout for url@ref.
// If lockedSHA is non-empty it is used directly without ls-remote resolution.
// env is passed to all git subprocesses.
func (c *GitCache) GetCheckout(url, ref, lockedSHA string, env []string) (string, error) {
	shardKey := urlnormalize.CacheKey(url)
	sha, err := c.resolveSHA(url, ref, lockedSHA, env)
	if err != nil {
		return "", err
	}

	checkoutDir := filepath.Join(c.checkoutsRoot, shardKey, sha)

	if !c.refresh {
		if fi, err := os.Stat(checkoutDir); err == nil && fi.IsDir() {
			if integrity.VerifyCheckout(checkoutDir, sha) {
				return checkoutDir, nil
			}
			// Integrity failure -- evict
			_ = os.RemoveAll(checkoutDir)
		}
	}

	if err := c.ensureBareRepo(url, shardKey, sha, env); err != nil {
		return "", err
	}
	return c.createCheckout(url, shardKey, sha, env)
}

func (c *GitCache) resolveSHA(url, ref, lockedSHA string, env []string) (string, error) {
	if lockedSHA != "" && fullSHARe.MatchString(lockedSHA) {
		return lockedSHA, nil
	}
	if ref != "" && fullSHARe.MatchString(ref) {
		return ref, nil
	}
	return c.lsRemoteSHA(url, ref, env)
}

func (c *GitCache) lsRemoteSHA(url, ref string, env []string) (string, error) {
	args := []string{"ls-remote", url}
	if ref != "" {
		args = append(args, ref)
	}
	cmd := exec.Command("git", args...)
	cmd.Env = mergeEnv(os.Environ(), env)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("gitcache: ls-remote %s %s: %w", sanitizeURL(url), ref, err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 1 && fullSHARe.MatchString(parts[0]) {
			return parts[0], nil
		}
	}
	return "", fmt.Errorf("gitcache: cannot resolve ref %q in %s", ref, sanitizeURL(url))
}

func (c *GitCache) ensureBareRepo(url, shardKey, sha string, env []string) error {
	bareDir := filepath.Join(c.dbRoot, shardKey)
	if fi, err := os.Stat(bareDir); err == nil && fi.IsDir() {
		// Verify the sha is present
		cmd := exec.Command("git", "-C", bareDir, "cat-file", "-e", sha+"^{commit}")
		cmd.Env = mergeEnv(os.Environ(), env)
		if cmd.Run() == nil {
			return nil
		}
		// Fetch to get the missing sha
		fetch := exec.Command("git", "-C", bareDir, "fetch", "--quiet", "origin")
		fetch.Env = mergeEnv(os.Environ(), env)
		_ = fetch.Run()
		return nil
	}

	staged := locking.StagePath(bareDir)
	if err := os.MkdirAll(filepath.Dir(staged), 0o700); err != nil {
		return fmt.Errorf("gitcache: staged dir: %w", err)
	}
	clone := exec.Command("git", "clone", "--bare", "--quiet", url, staged)
	clone.Env = mergeEnv(os.Environ(), env)
	if out, err := clone.CombinedOutput(); err != nil {
		_ = os.RemoveAll(staged)
		return fmt.Errorf("gitcache: bare clone %s: %w\n%s", sanitizeURL(url), err, out)
	}
	// Redact remote URL
	exec.Command("git", "-C", staged, "remote", "set-url", "origin", "redacted").Run() //nolint:errcheck

	lock := locking.NewShardLock(bareDir, 0)
	_, err := locking.AtomicLand(staged, bareDir, lock)
	return err
}

func (c *GitCache) createCheckout(url, shardKey, sha string, env []string) (string, error) {
	bareDir := filepath.Join(c.dbRoot, shardKey)
	checkoutDir := filepath.Join(c.checkoutsRoot, shardKey, sha)

	staged := locking.StagePath(checkoutDir)
	if err := os.MkdirAll(filepath.Dir(staged), 0o700); err != nil {
		return "", fmt.Errorf("gitcache: checkout stage dir: %w", err)
	}
	clone := exec.Command("git", "clone", "--quiet", "--local", bareDir, staged)
	clone.Env = mergeEnv(os.Environ(), env)
	if out, err := clone.CombinedOutput(); err != nil {
		_ = os.RemoveAll(staged)
		return "", fmt.Errorf("gitcache: clone from bare %s: %w\n%s", sanitizeURL(url), err, out)
	}
	checkout := exec.Command("git", "-C", staged, "checkout", "--quiet", sha)
	checkout.Env = mergeEnv(os.Environ(), env)
	if out, err := checkout.CombinedOutput(); err != nil {
		_ = os.RemoveAll(staged)
		return "", fmt.Errorf("gitcache: checkout %s: %w\n%s", sha[:12], err, out)
	}

	lock := locking.NewShardLock(checkoutDir, 0)
	if _, err := locking.AtomicLand(staged, checkoutDir, lock); err != nil {
		return "", err
	}
	return checkoutDir, nil
}

// EvictCheckout removes a cached checkout directory.
func (c *GitCache) EvictCheckout(checkoutDir string) {
	_ = os.RemoveAll(checkoutDir)
}

// GetCacheStats returns aggregate statistics.
func (c *GitCache) GetCacheStats() CacheStats {
	var stats CacheStats
	stats.DBCount = countDirs(c.dbRoot)
	stats.CheckoutCount = countDirs(c.checkoutsRoot)
	stats.TotalSizeBytes = dirSize(c.dbRoot) + dirSize(c.checkoutsRoot)
	return stats
}

// CleanAll removes all git cache content.
func (c *GitCache) CleanAll() {
	_ = os.RemoveAll(c.dbRoot)
	_ = os.RemoveAll(c.checkoutsRoot)
	_ = os.MkdirAll(c.dbRoot, 0o700)
	_ = os.MkdirAll(c.checkoutsRoot, 0o700)
}

// Prune removes checkout entries not accessed within maxAgeDays.
func (c *GitCache) Prune(maxAgeDays int) int {
	cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
	count := 0
	_ = filepath.WalkDir(c.checkoutsRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || path == c.checkoutsRoot {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().Before(cutoff) {
			_ = os.RemoveAll(path)
			count++
		}
		return nil
	})
	return count
}

// sanitizeURL strips credentials from a URL for logging.
func sanitizeURL(url string) string {
	if idx := strings.Index(url, "@"); idx != -1 {
		if proto := strings.Index(url, "://"); proto != -1 && proto < idx {
			return url[:proto+3] + "***@" + url[idx+1:]
		}
	}
	return url
}

func mergeEnv(base, extra []string) []string {
	if len(extra) == 0 {
		return base
	}
	return append(base, extra...)
}

func countDirs(root string) int {
	n := 0
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && d.IsDir() && path != root {
			n++
		}
		return nil
	})
	return n
}

func dirSize(root string) int64 {
	var total int64
	_ = filepath.WalkDir(root, func(_ string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err == nil {
			total += info.Size()
		}
		return nil
	})
	return total
}
