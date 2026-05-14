// Package skilltransformer converts SKILL.md primitives to platform-native formats.
// Mirrors src/apm_cli/integration/skill_transformer.py.
package skilltransformer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Skill holds the minimal fields from primitives.Skill needed by this package.
type Skill struct {
	Name        string
	Description string
	Content     string
	Source      string
}

var (
	reCamel       = regexp.MustCompile(`([a-z])([A-Z])`)
	reInvalidChar = regexp.MustCompile(`[^a-z0-9-]`)
	reConsecHyph  = regexp.MustCompile(`-+`)
)

// ToHyphenCase converts a name to hyphen-case for file naming.
// Handles underscores, spaces, and camelCase.
func ToHyphenCase(name string) string {
	result := strings.ReplaceAll(name, "_", "-")
	result = strings.ReplaceAll(result, " ", "-")
	result = reCamel.ReplaceAllString(result, "$1-$2")
	result = strings.ToLower(result)
	result = reInvalidChar.ReplaceAllString(result, "")
	result = reConsecHyph.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	return result
}

// SkillTransformer transforms SKILL.md to platform-native formats.
type SkillTransformer struct{}

// TransformToAgent transforms SKILL.md -> .github/agents/{name}.agent.md for VSCode.
// Returns the path where the file would be written. If dryRun is true, no file is created.
func (t *SkillTransformer) TransformToAgent(skill Skill, outputDir string, dryRun bool) (string, error) {
	content := t.generateAgentContent(skill)
	agentName := ToHyphenCase(skill.Name)
	agentPath := filepath.Join(outputDir, ".github", "agents", fmt.Sprintf("%s.agent.md", agentName))
	if dryRun {
		return agentPath, nil
	}
	if err := os.MkdirAll(filepath.Dir(agentPath), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(agentPath, []byte(content), 0o644); err != nil {
		return "", err
	}
	return agentPath, nil
}

// generateAgentContent builds the agent.md content with frontmatter.
func (t *SkillTransformer) generateAgentContent(skill Skill) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("name: %s\n", skill.Name))
	sb.WriteString(fmt.Sprintf("description: %s\n", skill.Description))
	sb.WriteString("---\n\n")
	if skill.Source != "" && skill.Source != "local" {
		sb.WriteString(fmt.Sprintf("<!-- Source: %s -->\n\n", skill.Source))
	}
	sb.WriteString(skill.Content)
	return sb.String()
}

// GetAgentName returns the hyphen-case agent filename for a skill.
func (t *SkillTransformer) GetAgentName(skill Skill) string {
	return ToHyphenCase(skill.Name)
}
