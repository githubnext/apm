// apmyml.go provides a minimal apm.yml parser for the Go CLI rewrite.
// Only the fields needed for read-only CLI commands are parsed.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ApmProject holds the parsed apm.yml structure.
type ApmProject struct {
	Name        string
	Version     string
	Description string
	Author      string
	Targets     []string
	Scripts     map[string]string
	Deps        []ApmDep
	MCPDeps     []ApmDep
	Marketplaces []ApmMarketplace
	PolicyDeny  []string
}

// ApmDep is a single dependency entry (owner/repo or owner/repo@ref).
type ApmDep struct {
	Package string
	Ref     string
}

// ApmMarketplace is a registered marketplace source.
type ApmMarketplace struct {
	Name string
	URL  string
}

// findApmYML walks up from dir looking for apm.yml.
func findApmYML(dir string) (string, error) {
	current := dir
	for {
		candidate := filepath.Join(current, "apm.yml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("apm.yml not found (searched from %s)", dir)
}

// parseApmYML does a line-by-line best-effort parse of apm.yml.
// It handles simple YAML scalars and list entries only.
func parseApmYML(path string) (*ApmProject, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := &ApmProject{Scripts: map[string]string{}}
	scanner := bufio.NewScanner(f)

	var section string
	var depSection string // "apm" or "mcp"
	var policySubSection string // "dependencies.deny" etc

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Top-level key detection.
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "-") {
			if idx := strings.Index(trimmed, ":"); idx >= 0 {
				key := strings.TrimSpace(trimmed[:idx])
				val := strings.TrimSpace(trimmed[idx+1:])
				switch key {
				case "name":
					p.Name = unquote(val)
					section = ""
				case "version":
					p.Version = unquote(val)
					section = ""
				case "description":
					p.Description = unquote(val)
					section = ""
				case "author":
					p.Author = unquote(val)
					section = ""
				case "targets":
					section = "targets"
					if val != "" && val != "[]" {
						p.Targets = parseInlineList(val)
					}
				case "scripts":
					section = "scripts"
				case "dependencies":
					section = "dependencies"
					depSection = ""
				case "marketplace":
					section = "marketplace"
				case "policy":
					section = "policy"
				default:
					section = key
				}
				continue
			}
		}

		// Section-specific parsing.
		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		switch section {
		case "targets":
			if strings.HasPrefix(trimmed, "-") {
				val := strings.TrimSpace(trimmed[1:])
				if val != "" {
					p.Targets = append(p.Targets, unquote(val))
				}
			}
		case "scripts":
			if idx := strings.Index(trimmed, ":"); idx >= 0 {
				key := strings.TrimSpace(trimmed[:idx])
				val := strings.TrimSpace(trimmed[idx+1:])
				if key != "" && !strings.HasPrefix(key, "-") {
					p.Scripts[key] = unquote(val)
				}
			}
		case "dependencies":
			if indent == 2 || indent == 0 {
				if strings.HasSuffix(trimmed, ":") {
					depSection = strings.TrimSuffix(trimmed, ":")
				}
			}
			if strings.HasPrefix(trimmed, "-") {
				val := strings.TrimSpace(trimmed[1:])
				if val != "" {
					dep := parseDep(unquote(val))
					switch depSection {
					case "apm":
						p.Deps = append(p.Deps, dep)
					case "mcp":
						p.MCPDeps = append(p.MCPDeps, dep)
					}
				}
			}
		case "marketplace":
			// Parse marketplace entries (name: URL or - name: url)
			if idx := strings.Index(trimmed, ":"); idx >= 0 {
				key := strings.TrimSpace(trimmed[:idx])
				val := strings.TrimSpace(trimmed[idx+1:])
				if key != "" && val != "" && !strings.HasPrefix(key, "-") {
					p.Marketplaces = append(p.Marketplaces, ApmMarketplace{Name: key, URL: unquote(val)})
				}
			}
		case "policy":
			if trimmed == "deny:" {
				policySubSection = "deny"
			} else if trimmed == "dependencies:" {
				policySubSection = "dependencies"
			} else if policySubSection == "deny" && strings.HasPrefix(trimmed, "- ") {
				p.PolicyDeny = append(p.PolicyDeny, unquote(strings.TrimSpace(trimmed[2:])))
			} else if !strings.HasPrefix(trimmed, "- ") && !strings.Contains(trimmed, ":") {
				policySubSection = ""
			}
		}
	}
	return p, scanner.Err()
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}

func parseInlineList(s string) []string {
	s = strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(s), "]"), "[")
	var out []string
	for _, part := range strings.Split(s, ",") {
		v := strings.TrimSpace(part)
		if v != "" {
			out = append(out, unquote(v))
		}
	}
	return out
}

func parseDep(s string) ApmDep {
	parts := strings.SplitN(s, "@", 2)
	if len(parts) == 2 {
		return ApmDep{Package: parts[0], Ref: parts[1]}
	}
	return ApmDep{Package: s}
}
