// Package aggregator scans workflow files for MCP dependencies.
package aggregator

import (
"bufio"
"os"
"path/filepath"
"strings"
)

// ScanWorkflowsForDependencies scans .prompt.md files for MCP dependencies.
func ScanWorkflowsForDependencies(baseDir string) (map[string]bool, error) {
if baseDir == "" {
var err error
baseDir, err = os.Getwd()
if err != nil {
return nil, err
}
}

servers := map[string]bool{}
err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
if err != nil {
return nil
}
if d.IsDir() || !strings.HasSuffix(path, ".prompt.md") {
return nil
}
if mcps, parseErr := parseMCPFromPromptFile(path); parseErr == nil {
for _, s := range mcps {
servers[s] = true
}
}
return nil
})
return servers, err
}

func parseMCPFromPromptFile(filePath string) ([]string, error) {
f, err := os.Open(filePath)
if err != nil {
return nil, err
}
defer f.Close()

var result []string
inFrontmatter := false
inMCP := false
firstLine := true
scanner := bufio.NewScanner(f)
for scanner.Scan() {
line := scanner.Text()
if firstLine {
firstLine = false
if strings.TrimSpace(line) == "---" {
inFrontmatter = true
continue
}
return nil, nil
}
if inFrontmatter {
if strings.TrimSpace(line) == "---" {
break
}
trimmed := strings.TrimSpace(line)
if strings.HasPrefix(trimmed, "mcp:") {
val := strings.TrimSpace(strings.TrimPrefix(trimmed, "mcp:"))
if val == "" {
inMCP = true
}
continue
}
if inMCP {
if strings.HasPrefix(line, "  - ") || strings.HasPrefix(line, "- ") {
val := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "")
result = append(result, val)
continue
}
inMCP = false
}
}
}
return result, scanner.Err()
}
