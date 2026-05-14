// Package wfparser parses workflow definition files with YAML frontmatter.
package wfparser

import (
"bufio"
"os"
"strings"
)

// WorkflowDefinition holds parsed workflow data.
type WorkflowDefinition struct {
Name            string
FilePath        string
Description     string
Author          string
MCPDependencies []string
InputParameters []string
LLMModel        string
Content         string
}

// Validate returns validation errors for the workflow.
func (w *WorkflowDefinition) Validate() []string {
var errs []string
if w.Description == "" {
errs = append(errs, "Missing 'description' in frontmatter")
}
return errs
}

// ParseWorkflowFile parses a workflow file with YAML frontmatter.
func ParseWorkflowFile(filePath string) (*WorkflowDefinition, error) {
data, err := os.ReadFile(filePath)
if err != nil {
return nil, err
}
meta, content := splitFrontmatter(string(data))
name := workflowName(filePath)
w := &WorkflowDefinition{
Name:     name,
FilePath: filePath,
Content:  content,
}
parseFrontmatter(meta, w)
return w, nil
}

func workflowName(filePath string) string {
parts := strings.Split(filePath, string(os.PathSeparator))
base := parts[len(parts)-1]
base = strings.TrimSuffix(base, ".prompt.md")
base = strings.TrimSuffix(base, ".md")
return base
}

func splitFrontmatter(content string) (meta, body string) {
if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
return "", content
}
rest := content[4:]
end := strings.Index(rest, "\n---")
if end < 0 {
return "", content
}
return rest[:end], rest[end+4:]
}

func parseFrontmatter(meta string, w *WorkflowDefinition) {
scanner := bufio.NewScanner(strings.NewReader(meta))
var inMCP, inInput bool
for scanner.Scan() {
line := scanner.Text()
trimmed := strings.TrimSpace(line)
if trimmed == "" {
inMCP = false
inInput = false
continue
}
if kv := parseKV(trimmed); kv[0] != "" {
inMCP = false
inInput = false
switch kv[0] {
case "description":
w.Description = kv[1]
case "author":
w.Author = kv[1]
case "llm":
w.LLMModel = kv[1]
case "mcp":
if kv[1] == "" {
inMCP = true
}
case "input":
if kv[1] == "" {
inInput = true
}
}
} else if strings.HasPrefix(line, "  - ") || strings.HasPrefix(line, "- ") {
val := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "")
if inMCP {
w.MCPDependencies = append(w.MCPDependencies, val)
} else if inInput {
w.InputParameters = append(w.InputParameters, val)
}
}
}
}

func parseKV(line string) [2]string {
idx := strings.Index(line, ":")
if idx < 0 {
return [2]string{}
}
key := strings.TrimSpace(line[:idx])
val := strings.TrimSpace(line[idx+1:])
return [2]string{key, val}
}
