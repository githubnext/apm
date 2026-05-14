// Package plugin provides data models for APM plugin management.
//
// Mirrors src/apm_cli/models/plugin.py.
package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PluginMetadata holds metadata for a plugin.
type PluginMetadata struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Author       string   `json:"author"`
	Repository   string   `json:"repository,omitempty"`
	Homepage     string   `json:"homepage,omitempty"`
	License      string   `json:"license,omitempty"`
	Tags         []string `json:"tags"`
	Dependencies []string `json:"dependencies"`
}

// ToDict converts metadata to a map for JSON serialisation.
func (m *PluginMetadata) ToDict() map[string]interface{} {
	tags := m.Tags
	if tags == nil {
		tags = []string{}
	}
	deps := m.Dependencies
	if deps == nil {
		deps = []string{}
	}
	return map[string]interface{}{
		"id":           m.ID,
		"name":         m.Name,
		"version":      m.Version,
		"description":  m.Description,
		"author":       m.Author,
		"repository":   m.Repository,
		"homepage":     m.Homepage,
		"license":      m.License,
		"tags":         tags,
		"dependencies": deps,
	}
}

// MetadataFromDict creates PluginMetadata from a JSON-decoded map.
func MetadataFromDict(data map[string]interface{}) (*PluginMetadata, error) {
	getString := func(key string) (string, bool) {
		v, ok := data[key]
		if !ok || v == nil {
			return "", false
		}
		s, ok := v.(string)
		return s, ok
	}
	getStrings := func(key string) []string {
		v, ok := data[key]
		if !ok || v == nil {
			return nil
		}
		raw, ok := v.([]interface{})
		if !ok {
			return nil
		}
		out := make([]string, 0, len(raw))
		for _, item := range raw {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}

	id, ok := getString("id")
	if !ok {
		return nil, fmt.Errorf("missing required field: id")
	}
	name, ok := getString("name")
	if !ok {
		return nil, fmt.Errorf("missing required field: name")
	}
	version, ok := getString("version")
	if !ok {
		return nil, fmt.Errorf("missing required field: version")
	}
	description, _ := getString("description")
	author, _ := getString("author")
	repository, _ := getString("repository")
	homepage, _ := getString("homepage")
	license, _ := getString("license")

	return &PluginMetadata{
		ID:           id,
		Name:         name,
		Version:      version,
		Description:  description,
		Author:       author,
		Repository:   repository,
		Homepage:     homepage,
		License:      license,
		Tags:         getStrings("tags"),
		Dependencies: getStrings("dependencies"),
	}, nil
}

// Plugin represents an installed plugin.
type Plugin struct {
	Metadata *PluginMetadata
	Path     string
	Commands []string
	Agents   []string
	Hooks    []string
	Skills   []string
}

// findPluginJSON locates the plugin.json file under pluginPath.
// It checks root, .github/plugin/, and .claude-plugin/ in order.
func findPluginJSON(pluginPath string) string {
	candidates := []string{
		filepath.Join(pluginPath, "plugin.json"),
		filepath.Join(pluginPath, ".github", "plugin", "plugin.json"),
		filepath.Join(pluginPath, ".claude-plugin", "plugin.json"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// globRec walks root for files matching the given extension (e.g. ".py").
func globRec(root, ext string) []string {
	var out []string
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ext {
			out = append(out, path)
		}
		return nil
	})
	return out
}

// globRecSuffix walks root for files whose base name has the given suffix.
func globRecSuffix(root, suffix string) []string {
	var out []string
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := filepath.Base(path)
		if len(name) >= len(suffix) && name[len(name)-len(suffix):] == suffix {
			out = append(out, path)
		}
		return nil
	})
	return out
}

// FromPath loads a Plugin from its installation directory.
//
// Plugin structure: plugin.json can be in root, .github/plugin/, or
// .claude-plugin/. Primitives are always at the repository root.
func FromPath(pluginPath string) (*Plugin, error) {
	metaFile := findPluginJSON(pluginPath)
	if metaFile == "" {
		return nil, fmt.Errorf("plugin metadata not found in any expected location: %s", pluginPath)
	}

	raw, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, fmt.Errorf("reading plugin.json: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("invalid plugin.json: %w", err)
	}

	meta, err := MetadataFromDict(data)
	if err != nil {
		return nil, err
	}

	// Discover components at repo root.
	commandsDir := filepath.Join(pluginPath, "commands")
	var commands []string
	if _, e := os.Stat(commandsDir); e == nil {
		commands = globRec(commandsDir, ".py")
	}

	agentsDir := filepath.Join(pluginPath, "agents")
	var agents []string
	if _, e := os.Stat(agentsDir); e == nil {
		agents = globRecSuffix(agentsDir, ".md")
	}

	hooksDir := filepath.Join(pluginPath, "hooks")
	var hooks []string
	if _, e := os.Stat(hooksDir); e == nil {
		hooks = globRec(hooksDir, ".py")
	}

	// Skills: each subdirectory must contain SKILL.md.
	skillsDir := filepath.Join(pluginPath, "skills")
	var skills []string
	if entries, e := os.ReadDir(skillsDir); e == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
			if _, se := os.Stat(skillFile); se == nil {
				skills = append(skills, skillFile)
			}
		}
	}

	return &Plugin{
		Metadata: meta,
		Path:     pluginPath,
		Commands: commands,
		Agents:   agents,
		Hooks:    hooks,
		Skills:   skills,
	}, nil
}
