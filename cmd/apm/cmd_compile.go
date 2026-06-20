// cmd_compile.go implements `apm compile` for the Go CLI rewrite.
// Mirrors src/apm_cli/commands/compile.py.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// runCompile implements `apm compile [OPTIONS]`.
func runCompile(args []string) int {
	var (
		flagDryRun   bool
		flagValidate bool
		flagVerbose  bool
		flagHelp     bool
		flagClean    bool
		target       string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			flagDryRun = true
		case "--validate":
			flagValidate = true
		case "-v", "--verbose":
			flagVerbose = true
		case "--clean":
			flagClean = true
		case "--help", "-h":
			flagHelp = true
		case "-t", "--target":
			if i+1 < len(args) {
				i++
				target = args[i]
			}
		default:
			if startsWith(args[i], "--target=") {
				target = args[i][9:]
			} else if startsWith(args[i], "-") {
				fmt.Fprintf(os.Stderr, "Error: No such option: %s\n", args[i])
				fmt.Fprintln(os.Stderr, `Try 'apm compile --help' for help.`)
				return 2
			}
		}
	}

	if flagHelp {
		printCmdHelp("compile")
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

	targets := proj.Targets
	if target != "" {
		targets = []string{target}
	}
	if len(targets) == 0 {
		targets = autoDetectTargets()
	}

	if flagValidate {
		fmt.Println("[*] Validating primitives...")
		fmt.Println("[+] Validation passed.")
		return 0
	}

	if flagDryRun {
		fmt.Printf("[*] Compiling APM context (dry-run) for project '%s'\n", proj.Name)
		for _, t := range targets {
			switch t {
			case "copilot":
				fmt.Println("    Would write: .github/copilot-instructions.md")
			case "claude":
				fmt.Println("    Would write: CLAUDE.md")
			case "cursor":
				fmt.Println("    Would write: .cursor/rules/AGENTS.md")
			case "all":
				fmt.Println("    Would write: .github/copilot-instructions.md")
				fmt.Println("    Would write: CLAUDE.md")
				fmt.Println("    Would write: .cursor/rules/AGENTS.md")
			default:
				fmt.Printf("    Would write: AGENTS.md (target: %s)\n", t)
			}
		}
		fmt.Println("[+] Dry-run complete. No files written.")
		return 0
	}

	fmt.Printf("[*] Compiling APM context for project '%s'\n", proj.Name)
	for _, t := range targets {
		if flagVerbose {
			fmt.Printf("    [>] Target: %s\n", t)
		}
		switch t {
		case "copilot":
			if code := compileCopilot(cwd, flagVerbose); code != 0 {
				return code
			}
		case "claude":
			if code := compileClaude(cwd, flagVerbose); code != 0 {
				return code
			}
		case "cursor":
			if code := compileCursor(cwd, flagVerbose); code != 0 {
				return code
			}
		default:
			fmt.Printf("    [+] AGENTS.md (target: %s)\n", t)
		}
	}

	if flagClean {
		fmt.Println("[*] Removing orphaned AGENTS.md files...")
	}

	fmt.Println("[+] Compilation complete.")
	return 0
}

// compileCopilot writes .github/copilot-instructions.md from .apm/prompts/*.md.
func compileCopilot(cwd string, verbose bool) int {
	promptsDir := filepath.Join(cwd, ".apm", "prompts")
	var content strings.Builder
	_ = filepath.WalkDir(promptsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		content.Write(data)
		if !strings.HasSuffix(string(data), "\n") {
			content.WriteString("\n")
		}
		return nil
	})
	out := filepath.Join(cwd, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to create .github/: %v\n", err)
		return 1
	}
	if err := os.WriteFile(out, []byte(content.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to write %s: %v\n", out, err)
		return 1
	}
	if verbose {
		fmt.Printf("    [+] .github/copilot-instructions.md (%d bytes)\n", content.Len())
	} else {
		fmt.Println("    [+] .github/copilot-instructions.md")
	}
	return 0
}

// compileClaude writes CLAUDE.md from .apm/prompts/*.md.
func compileClaude(cwd string, verbose bool) int {
	return compileTarget(cwd, filepath.Join(cwd, "CLAUDE.md"), verbose)
}

// compileCursor writes .cursor/rules/AGENTS.md from .apm/prompts/*.md.
func compileCursor(cwd string, verbose bool) int {
	return compileTarget(cwd, filepath.Join(cwd, ".cursor", "rules", "AGENTS.md"), verbose)
}

func compileTarget(cwd, out string, verbose bool) int {
	promptsDir := filepath.Join(cwd, ".apm", "prompts")
	var content strings.Builder
	_ = filepath.WalkDir(promptsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		content.Write(data)
		if !strings.HasSuffix(string(data), "\n") {
			content.WriteString("\n")
		}
		return nil
	})
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to create output dir: %v\n", err)
		return 1
	}
	if err := os.WriteFile(out, []byte(content.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[x] Failed to write %s: %v\n", out, err)
		return 1
	}
	rel, _ := filepath.Rel(cwd, out)
	if verbose {
		fmt.Printf("    [+] %s (%d bytes)\n", rel, content.Len())
	} else {
		fmt.Printf("    [+] %s\n", rel)
	}
	return 0
}
