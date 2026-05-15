// Package listcmd implements the "apm list" command, which prints available
// scripts from the project apm.yml file.
package listcmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Script represents a named runnable script from apm.yml.
type Script struct {
	Name    string
	Command string
}

// parseScripts extracts the scripts section from apm.yml content using a simple line scanner.
func parseScripts(content string) map[string]string {
	scripts := make(map[string]string)
	inScripts := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		// Detect top-level "scripts:" block
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			if strings.HasPrefix(trimmed, "scripts:") {
				inScripts = true
				continue
			}
			if inScripts {
				break // left the scripts block
			}
			continue
		}
		if !inScripts {
			continue
		}
		// Inside scripts block: "  name: command"
		if idx := strings.Index(trimmed, ":"); idx > 0 {
			name := strings.TrimSpace(trimmed[:idx])
			cmd := strings.TrimSpace(trimmed[idx+1:])
			// Strip surrounding quotes
			if len(cmd) >= 2 && ((cmd[0] == '"' && cmd[len(cmd)-1] == '"') ||
				(cmd[0] == '\'' && cmd[len(cmd)-1] == '\'')) {
				cmd = cmd[1 : len(cmd)-1]
			}
			scripts[name] = cmd
		}
	}
	return scripts
}

// ListScripts reads apm.yml from the current directory and returns all scripts.
func ListScripts() ([]Script, error) {
	raw, err := os.ReadFile("apm.yml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("listcmd: read apm.yml: %w", err)
	}

	scripts := parseScripts(string(raw))
	var result []Script
	for name, cmd := range scripts {
		result = append(result, Script{Name: name, Command: cmd})
	}
	return result, nil
}

// Run executes the list command, printing available scripts.
func Run() error {
	scripts, err := ListScripts()
	if err != nil {
		return err
	}

	if len(scripts) == 0 {
		fmt.Println("[!] No scripts found.")
		fmt.Println()
		fmt.Println("    Add scripts to your apm.yml file, for example:")
		fmt.Println("    scripts:")
		fmt.Println(`      start: "codex run main.prompt.md"`)
		return nil
	}

	hasStart := false
	for _, s := range scripts {
		if s.Name == "start" {
			hasStart = true
			break
		}
	}

	fmt.Println()
	fmt.Printf("  %-20s  %s\n", "SCRIPT", "COMMAND")
	fmt.Printf("  %-20s  %s\n", "------", "-------")
	for _, s := range scripts {
		marker := "  "
		if s.Name == "start" {
			marker = ">>"
		}
		fmt.Printf("  %s %-18s  %s\n", marker, s.Name, s.Command)
	}

	if hasStart {
		fmt.Println()
		fmt.Println("  [i] >> = default script (runs when no script name specified)")
	}
	return nil
}
