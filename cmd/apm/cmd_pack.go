// cmd_pack.go implements `apm pack` and `apm unpack` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/pack.py.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// runPack implements `apm pack [OPTIONS]`.
func runPack(args []string) int {
	var (
		flagDryRun bool
		flagHelp   bool
		flagJSON   bool
		output     string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			flagDryRun = true
		case "--json":
			flagJSON = true
		case "--help", "-h":
			flagHelp = true
		case "-o", "--output":
			if i+1 < len(args) {
				i++
				output = args[i]
			}
		default:
			if startsWith(args[i], "--output=") {
				output = args[i][9:]
			} else if startsWith(args[i], "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[i])
				fmt.Fprintln(os.Stderr, `Try 'apm pack --help' for help.`)
				return 2
			}
		}
	}

	if flagHelp {
		printCmdHelp("pack")
		return 0
	}

	if output == "" {
		output = "./build"
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

	if flagDryRun {
		if flagJSON {
			fmt.Printf(`{"project":%q,"output":%q,"dry_run":true,"artifacts":[]}`, proj.Name, output)
			fmt.Println()
		} else {
			fmt.Printf("[*] Packing project '%s' (dry-run)\n", proj.Name)
			fmt.Printf("    Output: %s\n", output)
			fmt.Printf("    Dependencies: %d APM, %d MCP\n", len(proj.Deps), len(proj.MCPDeps))
			fmt.Println("[+] Dry-run complete. No files written.")
		}
		return 0
	}

	if flagJSON {
		fmt.Printf(`{"project":%q,"output":%q,"artifacts":[%q]}`, proj.Name, output, filepath.Join(output, proj.Name+"-"+proj.Version+".apm"))
		fmt.Println()
	} else {
		fmt.Printf("[*] Packing project '%s'\n", proj.Name)
		fmt.Printf("    Output: %s\n", output)
	}

	// Create output directory and write a bundle manifest.
	outDir := filepath.Join(cwd, output)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to create output dir: %v\n", err)
		return 1
	}
	bundleName := proj.Name + "-" + proj.Version + ".apm"
	bundlePath := filepath.Join(outDir, bundleName)
	bundleContent := "name: " + proj.Name + "\nversion: " + proj.Version + "\n"
	if err := os.WriteFile(bundlePath, []byte(bundleContent), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to write bundle: %v\n", err)
		return 1
	}

	if !flagJSON {
		fmt.Println("[+] Pack complete.")
	}
	return 0
}

// runUnpack implements `apm unpack [OPTIONS]`.
func runUnpack(args []string) int {
	var (
		flagHelp bool
		bundle   string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			flagHelp = true
		default:
			if startsWith(args[i], "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[i])
				fmt.Fprintln(os.Stderr, "Try 'apm unpack --help' for help.")
				return 2
			}
			if bundle == "" {
				bundle = args[i]
			}
		}
	}

	if flagHelp {
		printCmdHelp("unpack")
		return 0
	}

	if bundle == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing BUNDLE argument.")
		fmt.Fprintln(os.Stderr, `Try 'apm unpack --help' for help.`)
		return 2
	}

	if _, err := os.Stat(bundle); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Bundle not found: %s\n", bundle)
		return 1
	}

	fmt.Printf("[*] Unpacking bundle: %s\n", bundle)

	cwd, _ := os.Getwd()
	bundleAbs := bundle
	if !filepath.IsAbs(bundle) {
		bundleAbs = filepath.Join(cwd, bundle)
	}

	if err := copyDirTree(bundleAbs, cwd); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to unpack: %v\n", err)
		return 1
	}

	fmt.Println("[+] Unpack complete.")
	return 0
}
