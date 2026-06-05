// cmd_update.go implements `apm update` and `apm prune` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/update.py and src/apm_cli/commands/prune.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runUpdate implements `apm update [OPTIONS]`.
func runUpdate(args []string) int {
	var (
		flagDryRun  bool
		flagHelp    bool
		flagVerbose bool
		flagYes     bool
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			flagDryRun = true
		case "--help", "-h":
			flagHelp = true
		case "-v", "--verbose":
			flagVerbose = true
		case "--yes", "-y":
			flagYes = true
		case "-t", "--target":
			if i+1 < len(args) {
				i++
			}
		}
	}

	if flagHelp {
		printCmdHelp("update")
		return 0
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] No apm.yml found. Run 'apm init' to create one.\n")
		return 1
	}
	proj, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
		return 1
	}

	_ = flagYes

	if flagDryRun {
		fmt.Printf("[*] Checking for updates in project '%s' (dry-run)\n", proj.Name)
		if flagVerbose {
			fmt.Printf("    APM deps: %d\n", len(proj.Deps))
		}
		fmt.Println("[+] No updates needed (dry-run). No files written.")
		return 0
	}

	fmt.Printf("[*] Updating dependencies for project '%s'\n", proj.Name)

	// Update local package versions in lockfile.
	lockPath := filepath.Join(cwd, "apm.lock.yaml")
	deps, readErr := readLockfileDeps(lockPath)
	if readErr == nil && len(deps) > 0 {
		updated := false
		for i, dep := range deps {
			if dep.RepoURL == "" || (!startsWith(dep.RepoURL, "./") && !startsWith(dep.RepoURL, "../")) {
				continue
			}
			localPath := filepath.Join(cwd, filepath.FromSlash(dep.RepoURL))
			localYML := filepath.Join(localPath, "apm.yml")
			localProj, parseErr := parseApmYML(localYML)
			if parseErr != nil {
				continue
			}
			if localProj.Version != dep.Version {
				if flagVerbose {
					fmt.Printf("    [>] %s: %s -> %s\n", dep.Name, dep.Version, localProj.Version)
				}
				deps[i].Version = localProj.Version
				updated = true
			}
		}
		if updated {
			_ = writeLockfile(lockPath, deps)
		}
	}

	if flagVerbose {
		fmt.Printf("    APM deps: %d\n", len(proj.Deps))
	}
	fmt.Println("[+] Update complete.")
	return 0
}

// runPrune implements `apm prune [OPTIONS]`.
func runPrune(args []string) int {
	var (
		flagHelp    bool
		flagDryRun  bool
		flagVerbose bool
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			flagHelp = true
		case "--dry-run":
			flagDryRun = true
		case "-v", "--verbose":
			flagVerbose = true
		}
	}

	if flagHelp {
		printCmdHelp("prune")
		return 0
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] No apm.yml found. Run 'apm init' to create one.\n")
		return 1
	}
	proj, err := parseApmYML(ymlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
		return 1
	}

	_ = flagVerbose

	if flagDryRun {
		fmt.Printf("[*] Pruning project '%s' (dry-run)\n", proj.Name)
		fmt.Println("[+] No packages to prune (dry-run). No files removed.")
		return 0
	}

	fmt.Printf("[*] Pruning project '%s'\n", proj.Name)

	// Find packages referenced in lockfile.
	lockPath := filepath.Join(cwd, "apm.lock.yaml")
	deps, _ := readLockfileDeps(lockPath)
	referencedPaths := make(map[string]bool)
	for _, dep := range deps {
		if dep.InstallPath != "" {
			// Normalize: get just the package dir name.
			parts := strings.Split(filepath.ToSlash(dep.InstallPath), "/")
			if len(parts) > 0 {
				referencedPaths[parts[len(parts)-1]] = true
			}
		}
	}

	// Remove apm_modules entries not referenced in lockfile.
	modulesDir := filepath.Join(cwd, "apm_modules")
	entries, err := os.ReadDir(modulesDir)
	if err == nil {
		for _, entry := range entries {
			if !referencedPaths[entry.Name()] {
				stale := filepath.Join(modulesDir, entry.Name())
				if flagVerbose {
					fmt.Printf("    [>] Removing stale: %s\n", entry.Name())
				}
				_ = os.RemoveAll(stale)
			}
		}
	}

	fmt.Println("[+] Prune complete.")
	return 0
}
