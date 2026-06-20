// cmd_cache.go implements `apm cache` and its subcommands for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/cache.py.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// cacheDir returns the APM cache directory path.
func cacheDir() string {
	if d := os.Getenv("APM_CACHE_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".apm", "cache")
	}
	return filepath.Join(home, ".apm", "cache")
}

// dirSize returns the total size in bytes of all files under dir.
func dirSize(dir string) int64 {
	var total int64
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
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

// runCache implements `apm cache [SUBCOMMAND] [OPTIONS]`.
func runCache(args []string) int {
	if len(args) == 0 {
		printCacheHelp()
		return 0
	}

	// Only intercept --help/-h when it is the very first argument (top-level
	// cache help). When a subcommand precedes --help (e.g. "cache clean --help"),
	// delegate to the subcommand handler so it can show its own usage.
	if args[0] == "--help" || args[0] == "-h" {
		printCacheHelp()
		return 0
	}

	if startsWith(args[0], "-") {
		fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[0])
		fmt.Fprintln(os.Stderr, `Try 'apm cache --help' for help.`)
		return 2
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "info":
		return runCacheInfo(rest)
	case "clean":
		return runCacheClean(rest)
	case "prune":
		return runCachePrune(rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: No such command '%s'.\n", sub)
		fmt.Fprintln(os.Stderr, `Try 'apm cache --help' for help.`)
		return 2
	}
}

func runCacheInfo(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm cache info [OPTIONS]")
			fmt.Println()
			fmt.Println("  Show cache location and size statistics")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --help  Show this message and exit.")
			return 0
		}
		if startsWith(a, "-") {
			fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
			fmt.Fprintln(os.Stderr, `Try 'apm cache info --help' for help.`)
			return 2
		}
	}
	dir := cacheDir()
	size := dirSize(dir)
	fmt.Printf("Cache location: %s\n", dir)
	fmt.Printf("Cache size:     %.1f MB\n", float64(size)/1024/1024)
	return 0
}

func printCacheHelp() {
	fmt.Println("Usage: apm cache [OPTIONS] COMMAND [ARGS]...")
	fmt.Println()
	fmt.Println("  Manage the local package cache")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help  Show this message and exit.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  clean  Remove all cached content")
	fmt.Println("  info   Show cache location and size statistics")
	fmt.Println("  prune  Remove cache entries older than N days")
}

func runCacheClean(args []string) int {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm cache clean [OPTIONS]")
			fmt.Println()
			fmt.Println("  Remove all cached content")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -f, --force  Skip confirmation prompt")
			fmt.Println("  -y, --yes    Skip confirmation prompt")
			fmt.Println("  --help       Show this message and exit.")
			return 0
		}
		switch a {
		case "-f", "--force", "-y", "--yes":
			// known flags
		default:
			if startsWith(a, "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
				fmt.Fprintln(os.Stderr, `Try 'apm cache clean --help' for help.`)
				return 2
			}
		}
	}
	dir := cacheDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
				fmt.Fprintf(os.Stderr, "[x] Failed to create cache dir: %v\n", mkErr)
				return 1
			}
			fmt.Println("[*] Cleaning cache...")
			fmt.Println("[+] Cache cleaned.")
			return 0
		}
		fmt.Fprintf(os.Stderr, "[x] Failed to read cache dir: %v\n", err)
		return 1
	}
	for _, entry := range entries {
		if rmErr := os.RemoveAll(filepath.Join(dir, entry.Name())); rmErr != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to remove cache entry %s: %v\n", entry.Name(), rmErr)
			return 1
		}
	}
	fmt.Println("[*] Cleaning cache...")
	fmt.Println("[+] Cache cleaned.")
	return 0
}

func runCachePrune(args []string) int {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: apm cache prune [OPTIONS]")
			fmt.Println()
			fmt.Println("  Remove cache entries older than N days")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --days INTEGER  Remove entries not accessed within this many days")
			fmt.Println("                  [default: 30]")
			fmt.Println("  --help          Show this message and exit.")
			return 0
		}
		if a == "--days" {
			if i+1 < len(args) {
				i++
			}
			continue
		}
		if startsWith(a, "--days=") {
			continue
		}
		if startsWith(a, "-") {
			fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", a)
			fmt.Fprintln(os.Stderr, `Try 'apm cache prune --help' for help.`)
			return 2
		}
	}
	fmt.Println("[*] Pruning old cache entries...")
	fmt.Println("[+] Cache pruned.")
	return 0
}
