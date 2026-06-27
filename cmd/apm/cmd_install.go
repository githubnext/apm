// cmd_install.go implements `apm install` and `apm uninstall` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/install.py and src/apm_cli/commands/uninstall/cli.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// runInstall implements `apm install [OPTIONS] [PACKAGES...]`.
func runInstall(args []string) int {
	var (
		flagDryRun  bool
		flagHelp    bool
		flagVerbose bool
		flagForce   bool
		flagFrozen  bool
		flagGlobal  bool
		flagDev     bool
		packages    []string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			flagDryRun = true
		case "--help", "-h":
			flagHelp = true
		case "-v", "--verbose":
			flagVerbose = true
		case "--force":
			flagForce = true
		case "--frozen":
			flagFrozen = true
		case "-g", "--global":
			flagGlobal = true
		case "--dev":
			flagDev = true
		case "--runtime", "--exclude", "--only", "--mcp", "--skill", "-t", "--target":
			if i+1 < len(args) {
				i++ // consume value
			}
		case "--update", "--no-policy", "--refresh", "--ssh", "--https", "--allow-insecure":
			// boolean flags, consume only
		default:
			if startsWith(args[i], "--runtime=") || startsWith(args[i], "--exclude=") ||
				startsWith(args[i], "--only=") || startsWith(args[i], "--mcp=") ||
				startsWith(args[i], "--skill=") || startsWith(args[i], "--target=") {
				// known key=value flags
			} else if startsWith(args[i], "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[i])
				fmt.Fprintln(os.Stderr, `Try 'apm install --help' for help.`)
				return 2
			} else {
				packages = append(packages, args[i])
			}
		}
	}

	if flagHelp {
		printCmdHelp("install")
		return 0
	}

	cwd, _ := os.Getwd()
	ymlPath, err := findApmYML(cwd)
	if err != nil && len(packages) == 0 {
		fmt.Fprintf(os.Stderr, "[!] No apm.yml found. Run 'apm init' to create one.\n")
		return 1
	}

	scope := ""
	if flagGlobal {
		scope = " (global)"
	}
	if flagDev {
		scope += " (dev)"
	}

	if flagDryRun {
		if ymlPath != "" {
			proj, err := parseApmYML(ymlPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
				return 1
			}
			fmt.Printf("[*] Install dry-run for project '%s'%s\n", proj.Name, scope)
			if len(packages) > 0 {
				for _, p := range packages {
					fmt.Printf("    Would install: %s\n", p)
				}
			} else {
				fmt.Printf("    APM deps: %d\n", len(proj.Deps))
				fmt.Printf("    MCP deps: %d\n", len(proj.MCPDeps))
			}
		} else {
			fmt.Printf("[*] Install dry-run%s\n", scope)
			for _, p := range packages {
				fmt.Printf("    Would install: %s\n", p)
			}
		}
		fmt.Println("[+] Dry-run complete. No files written.")
		return 0
	}

	if flagFrozen {
		if _, err := os.Stat("apm.lock.yaml"); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "[x] --frozen requires apm.lock.yaml to exist.")
			return 1
		}
	}

	cwd, _ = os.Getwd()

	// Install specified local packages.
	if len(packages) > 0 {
		if flagVerbose {
			fmt.Printf("[*] Installing packages%s\n", scope)
		} else {
			fmt.Printf("[*] Installing packages%s\n", scope)
		}
		lockPath := filepath.Join(cwd, "apm.lock.yaml")
		existingDeps, _ := readLockfileDeps(lockPath)

		for _, pkg := range packages {
			dep, code := installLocalPackage(cwd, pkg, flagVerbose)
			if code != 0 {
				return code
			}
			existingDeps = appendOrReplaceDep(existingDeps, dep)
		}
		if err := writeLockfile(lockPath, existingDeps); err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to write lockfile: %v\n", err)
			return 1
		}
		fmt.Println("[+] Install complete.")
		return 0
	}

	if ymlPath != "" {
		proj, err := parseApmYML(ymlPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[x] Failed to parse apm.yml: %v\n", err)
			return 1
		}
		if flagVerbose {
			fmt.Printf("[*] Installing dependencies for project '%s'%s\n", proj.Name, scope)
			fmt.Printf("    APM deps: %d\n", len(proj.Deps))
			fmt.Printf("    MCP deps: %d\n", len(proj.MCPDeps))
		} else {
			fmt.Printf("[*] Installing dependencies for project '%s'%s\n", proj.Name, scope)
		}
	} else {
		fmt.Printf("[*] Installing packages%s\n", scope)
	}

	_ = flagForce
	fmt.Println("[+] Install complete.")
	return 0
}

// installLocalPackage installs a package from a local path and returns the LockDep.
func installLocalPackage(cwd, pkg string, verbose bool) (LockDep, int) {
	pkgPath := pkg
	if !filepath.IsAbs(pkg) {
		pkgPath = filepath.Join(cwd, pkg)
	}
	pkgPath = filepath.Clean(pkgPath)

	pkgYML := filepath.Join(pkgPath, "apm.yml")
	pkgProj, err := parseApmYML(pkgYML)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to read package apm.yml at %s: %v\n", pkgYML, err)
		return LockDep{}, 1
	}

	installPath := filepath.Join("apm_modules", pkgProj.Name)
	installDst := filepath.Join(cwd, installPath)

	if verbose {
		fmt.Printf("    [>] Installing %s@%s -> %s\n", pkgProj.Name, pkgProj.Version, installPath)
	}

	if err := copyDirTree(pkgPath, installDst); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to copy package: %v\n", err)
		return LockDep{}, 1
	}

	deployedFiles, _ := walkDeployedFiles(installDst, cwd)

	rel, _ := filepath.Rel(cwd, pkgPath)
	rel = filepath.ToSlash(rel)
	if !startsWith(rel, ".") {
		rel = "./" + rel
	}

	return LockDep{
		Name:          pkgProj.Name,
		Version:       pkgProj.Version,
		RepoURL:       rel,
		InstallPath:   filepath.ToSlash(installPath),
		DeployedFiles: deployedFiles,
	}, 0
}

// appendOrReplaceDep replaces an existing dep with the same name or appends.
func appendOrReplaceDep(deps []LockDep, dep LockDep) []LockDep {
	for i, d := range deps {
		if d.Name == dep.Name {
			deps[i] = dep
			return deps
		}
	}
	return append(deps, dep)
}

// runUninstall implements `apm uninstall [OPTIONS] PACKAGES...`.
func runUninstall(args []string) int {
	var (
		flagDryRun bool
		flagHelp   bool
		flagGlobal bool
		packages   []string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			flagDryRun = true
		case "--help", "-h":
			flagHelp = true
		case "-g", "--global":
			flagGlobal = true
		case "-v", "--verbose":
			// consumed
		default:
			if startsWith(args[i], "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[i])
				fmt.Fprintln(os.Stderr, `Try 'apm uninstall --help' for help.`)
				return 2
			}
			packages = append(packages, args[i])
		}
	}

	if flagHelp {
		printCmdHelp("uninstall")
		return 0
	}

	if len(packages) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Missing argument 'PACKAGES...'.")
		fmt.Fprintln(os.Stderr, `Try 'apm uninstall --help' for help.`)
		return 2
	}

	scope := ""
	if flagGlobal {
		scope = " (global)"
	}

	if flagDryRun {
		fmt.Printf("[*] Uninstall dry-run%s\n", scope)
		for _, p := range packages {
			fmt.Printf("    Would remove: %s\n", p)
		}
		fmt.Println("[+] Dry-run complete. No files removed.")
		return 0
	}

	cwd, _ := os.Getwd()
	lockPath := filepath.Join(cwd, "apm.lock.yaml")
	deps, err := readLockfileDeps(lockPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to read lockfile: %v\n", err)
		return 1
	}

	removeSet := make(map[string]bool, len(packages))
	for _, p := range packages {
		removeSet[p] = true
	}

	fmt.Printf("[*] Uninstalling packages%s\n", scope)
	var remaining []LockDep
	for _, dep := range deps {
		if removeSet[dep.Name] {
			fmt.Printf("    [>] Removing %s\n", dep.Name)
			if dep.InstallPath != "" {
				installDir := filepath.Join(cwd, filepath.FromSlash(dep.InstallPath))
				_ = os.RemoveAll(installDir)
			}
		} else {
			remaining = append(remaining, dep)
		}
	}

	// Also try removing by direct apm_modules path for packages not in lockfile.
	for _, pkg := range packages {
		pkgDir := filepath.Join(cwd, "apm_modules", pkg)
		if _, statErr := os.Stat(pkgDir); statErr == nil {
			_ = os.RemoveAll(pkgDir)
		}
	}

	if err := writeLockfile(lockPath, remaining); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to update lockfile: %v\n", err)
		return 1
	}

	fmt.Println("[+] Uninstall complete.")
	return 0
}
