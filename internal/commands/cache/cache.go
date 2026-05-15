// Package cachecmd implements CLI commands for cache management (apm cache info|clean|prune).
package cachecmd

import (
	"fmt"
	"os"

	"github.com/githubnext/apm/internal/cache/cachepaths"
	"github.com/githubnext/apm/internal/cache/gitcache"
	"github.com/githubnext/apm/internal/cache/httpcache"
)

// CacheInfo prints cache location and size statistics.
func CacheInfo() error {
	root, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		return fmt.Errorf("[x] Cannot resolve cache root: %w", err)
	}

	fmt.Printf("[i] Cache root: %s\n", root)

	gc, err := gitcache.New(root, false)
	if err != nil {
		return fmt.Errorf("[x] Cannot open git cache: %w", err)
	}
	gitStats := gc.GetCacheStats()

	hc, err := httpcache.New(root)
	if err != nil {
		return fmt.Errorf("[x] Cannot open http cache: %w", err)
	}
	httpStats := hc.GetStats()

	totalBytes := gitStats.TotalSizeBytes + httpStats.TotalSizeBytes

	fmt.Println()
	fmt.Printf("  [*] Git repositories (db):    %d\n", gitStats.DBCount)
	fmt.Printf("  [*] Git checkouts:            %d\n", gitStats.CheckoutCount)
	fmt.Printf("  [*] HTTP cache entries:       %d\n", httpStats.EntryCount)
	fmt.Println()
	fmt.Printf("  [*] Total size:               %s\n", formatSize(totalBytes))
	fmt.Printf("      Git:                      %s\n", formatSize(gitStats.TotalSizeBytes))
	fmt.Printf("      HTTP:                     %s\n", formatSize(httpStats.TotalSizeBytes))
	return nil
}

// CacheClean removes all cached content after optional confirmation.
func CacheClean(force bool) error {
	root, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		return fmt.Errorf("[x] Cannot resolve cache root: %w", err)
	}

	if !force {
		fmt.Printf("Remove all cache content in %s? [y/N] ", root)
		var answer string
		fmt.Fscan(os.Stdin, &answer)
		if answer != "y" && answer != "Y" {
			fmt.Println("[i] Aborted.")
			return nil
		}
	}

	fmt.Println("[*] Cleaning cache...")

	gc, err := gitcache.New(root, false)
	if err != nil {
		return fmt.Errorf("[x] Cannot open git cache: %w", err)
	}
	gc.CleanAll()

	hc, err := httpcache.New(root)
	if err != nil {
		return fmt.Errorf("[x] Cannot open http cache: %w", err)
	}
	hc.CleanAll()

	fmt.Println("[+] Cache cleaned.")
	return nil
}

// CachePrune removes cache entries older than maxAgeDays days.
func CachePrune(maxAgeDays int) error {
	root, err := cachepaths.GetCacheRoot(false)
	if err != nil {
		return fmt.Errorf("[x] Cannot resolve cache root: %w", err)
	}

	fmt.Printf("[i] Pruning entries older than %d days...\n", maxAgeDays)

	gc, err := gitcache.New(root, false)
	if err != nil {
		return fmt.Errorf("[x] Cannot open git cache: %w", err)
	}
	pruned := gc.Prune(maxAgeDays)

	fmt.Printf("[+] Pruned %d checkout(s).\n", pruned)
	return nil
}

func formatSize(b int64) string {
	switch {
	case b < 1024:
		return fmt.Sprintf("%d B", b)
	case b < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	case b < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	default:
		return fmt.Sprintf("%.1f GB", float64(b)/(1024*1024*1024))
	}
}
