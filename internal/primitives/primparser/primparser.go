// Package primparser parses APM primitive definition files.
// Migrated from src/apm_cli/primitives/parser.py.
package primparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/githubnext/apm/internal/primitives/primmodels"
)

// ParseSkillFile parses a SKILL.md file and returns a Skill primitive.
// source is an optional identifier like "local" or "dependency:pkg".
func ParseSkillFile(filePath string, source string) (*primmodels.Skill, error) {
	meta, content, err := parseFrontmatter(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SKILL.md file %s: %w", filePath, err)
	}

	name := meta["name"]
	if name == "" {
		// Derive from parent directory name.
		name = filepath.Base(filepath.Dir(filePath))
	}

	return &primmodels.Skill{
		Name:        name,
		FilePath:    filePath,
		Description: meta["description"],
		Content:     content,
		Source:      source,
	}, nil
}

// ParsePrimitiveFile parses a primitive file (.chatmode.md, .instructions.md,
// .context.md, .memory.md) and returns the appropriate Primitive.
func ParsePrimitiveFile(filePath string, source string) (primmodels.Primitive, error) {
	meta, content, err := parseFrontmatter(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse primitive file %s: %w", filePath, err)
	}

	name := extractPrimitiveName(filePath)
	base := filepath.Base(filePath)

	switch {
	case strings.HasSuffix(base, ".chatmode.md") || strings.HasSuffix(base, ".agent.md"):
		return parseChatmode(name, filePath, meta, content, source), nil
	case strings.HasSuffix(base, ".instructions.md"):
		return parseInstruction(name, filePath, meta, content, source), nil
	case strings.HasSuffix(base, ".context.md") || strings.HasSuffix(base, ".memory.md") || isContextFile(filePath):
		return parseContext(name, filePath, meta, content, source), nil
	default:
		return nil, fmt.Errorf("unknown primitive file type: %s", filePath)
	}
}

// ValidatePrimitive returns a list of validation errors for the primitive.
func ValidatePrimitive(p primmodels.Primitive) []string {
	return p.Validate()
}

func parseChatmode(name, filePath string, meta map[string]string, content, source string) *primmodels.Chatmode {
	return &primmodels.Chatmode{
		Name:        name,
		FilePath:    filePath,
		Description: meta["description"],
		ApplyTo:     meta["applyTo"],
		Content:     content,
		Author:      meta["author"],
		Version:     meta["version"],
		Source:      source,
	}
}

func parseInstruction(name, filePath string, meta map[string]string, content, source string) *primmodels.Instruction {
	return &primmodels.Instruction{
		Name:        name,
		FilePath:    filePath,
		Description: meta["description"],
		ApplyTo:     meta["applyTo"],
		Content:     content,
		Author:      meta["author"],
		Version:     meta["version"],
		Source:      source,
	}
}

func parseContext(name, filePath string, meta map[string]string, content, source string) *primmodels.Context {
	return &primmodels.Context{
		Name:        name,
		FilePath:    filePath,
		Content:     content,
		Description: meta["description"],
		Author:      meta["author"],
		Version:     meta["version"],
		Source:      source,
	}
}

// extractPrimitiveName derives the primitive name from the file path following
// APM naming conventions.
func extractPrimitiveName(filePath string) string {
	abs, _ := filepath.Abs(filePath)
	parts := strings.Split(filepath.ToSlash(abs), "/")

	// Check for structured directories (.apm/ or .github/)
	subDirs := map[string]bool{
		"chatmodes": true, "instructions": true,
		"context": true, "memory": true, "agents": true,
	}
	for i, p := range parts {
		if (p == ".apm" || p == ".github") && i+2 < len(parts) && subDirs[parts[i+1]] {
			return stripPrimExt(filepath.Base(filePath))
		}
	}

	return stripPrimExt(filepath.Base(filePath))
}

func stripPrimExt(basename string) string {
	suffixes := []string{
		".chatmode.md", ".instructions.md", ".context.md",
		".memory.md", ".agent.md",
	}
	for _, s := range suffixes {
		if strings.HasSuffix(basename, s) {
			return strings.TrimSuffix(basename, s)
		}
	}
	if strings.HasSuffix(basename, ".md") {
		return strings.TrimSuffix(basename, ".md")
	}
	ext := filepath.Ext(basename)
	return strings.TrimSuffix(basename, ext)
}

// isContextFile returns true for files directly under .apm/memory/ or .github/memory/.
func isContextFile(filePath string) bool {
	dir := filepath.Base(filepath.Dir(filePath))
	parent := filepath.Base(filepath.Dir(filepath.Dir(filePath)))
	if dir != "memory" {
		return false
	}
	return parent == ".apm" || parent == ".github"
}

// parseFrontmatter reads a file and splits YAML frontmatter (--- ... ---) from
// the body. Returns the parsed key/value pairs and the body content.
// Only flat key: value pairs are supported (no nesting or lists).
func parseFrontmatter(filePath string) (map[string]string, string, error) {
	f, err := os.Open(filePath) // #nosec G304
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, "", err
	}

	meta := map[string]string{}
	if len(lines) == 0 {
		return meta, "", nil
	}

	// Check for leading frontmatter delimiter.
	if strings.TrimSpace(lines[0]) != "---" {
		return meta, strings.Join(lines, "\n"), nil
	}

	// Find closing delimiter.
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		// No closing delimiter -- treat entire file as content.
		return meta, strings.Join(lines, "\n"), nil
	}

	// Parse frontmatter block.
	for _, line := range lines[1:end] {
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		// Strip surrounding quotes.
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		meta[key] = val
	}

	content := strings.Join(lines[end+1:], "\n")
	return meta, strings.TrimLeft(content, "\n"), nil
}
