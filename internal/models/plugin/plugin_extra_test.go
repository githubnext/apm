package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetadataFromDict_AllFields(t *testing.T) {
	data := map[string]interface{}{
		"id":           "full-id",
		"name":         "Full Plugin",
		"version":      "3.0.0",
		"description":  "complete",
		"author":       "Bob",
		"repository":   "https://github.com/b/c",
		"homepage":     "https://home.example.com",
		"license":      "Apache-2.0",
		"tags":         []interface{}{"t1", "t2", "t3"},
		"dependencies": []interface{}{"dep-a", "dep-b"},
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Description != "complete" {
		t.Errorf("Description = %q", m.Description)
	}
	if m.Author != "Bob" {
		t.Errorf("Author = %q", m.Author)
	}
	if m.Repository != "https://github.com/b/c" {
		t.Errorf("Repository = %q", m.Repository)
	}
	if m.Homepage != "https://home.example.com" {
		t.Errorf("Homepage = %q", m.Homepage)
	}
	if m.License != "Apache-2.0" {
		t.Errorf("License = %q", m.License)
	}
	if len(m.Tags) != 3 {
		t.Errorf("Tags len = %d", len(m.Tags))
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("Dependencies len = %d", len(m.Dependencies))
	}
}

func TestMetadataFromDict_EmptyTags(t *testing.T) {
	data := map[string]interface{}{
		"id":      "e",
		"name":    "E",
		"version": "1",
		"tags":    []interface{}{},
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Tags) != 0 {
		t.Errorf("expected empty tags, got %v", m.Tags)
	}
}

func TestToDict_HasAllKeys(t *testing.T) {
	m := &PluginMetadata{
		ID:           "k",
		Name:         "K",
		Version:      "9",
		Description:  "desc",
		Author:       "auth",
		Repository:   "repo",
		Homepage:     "home",
		License:      "MIT",
		Tags:         []string{"a"},
		Dependencies: []string{"b"},
	}
	d := m.ToDict()
	for _, key := range []string{"id", "name", "version", "description", "author", "repository", "homepage", "license", "tags", "dependencies"} {
		if _, ok := d[key]; !ok {
			t.Errorf("ToDict missing key %q", key)
		}
	}
}

func TestToDict_TagsNotNilWhenSet(t *testing.T) {
	m := &PluginMetadata{
		ID: "t", Name: "T", Version: "1",
		Tags: []string{"x", "y"},
	}
	d := m.ToDict()
	tags, ok := d["tags"].([]string)
	if !ok {
		t.Fatalf("tags not []string: %T", d["tags"])
	}
	if len(tags) != 2 || tags[0] != "x" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestFromPath_MissingPluginJSON(t *testing.T) {
	dir := t.TempDir()
	_, err := FromPath(dir)
	if err == nil {
		t.Error("expected error for missing plugin.json")
	}
}

func TestFromPath_ClaudePluginSubdir(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, ".claude-plugin")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `{"id":"claude1","name":"ClaudePlugin","version":"0.0.1"}`
	if err := os.WriteFile(filepath.Join(subdir, "plugin.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := FromPath(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Metadata.ID != "claude1" {
		t.Errorf("unexpected ID: %s", p.Metadata.ID)
	}
}

func TestFromPath_WithAgents(t *testing.T) {
	dir := t.TempDir()
	content := `{"id":"ag","name":"AgPlugin","version":"1.0.0"}`
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	agentDir := filepath.Join(dir, "agents", "myagent")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(agentDir, "AGENT.md"), []byte("agent"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := FromPath(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(p.Agents))
	}
}

func TestFromPath_PathIsSet(t *testing.T) {
	dir := t.TempDir()
	content := `{"id":"pp","name":"PP","version":"1"}`
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := FromPath(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Path != dir {
		t.Errorf("expected path %s, got %s", dir, p.Path)
	}
}

func TestMetadataFromDict_IDWithSpaces(t *testing.T) {
	data := map[string]interface{}{
		"id":      "id with spaces",
		"name":    "Name",
		"version": "1",
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.ID != "id with spaces" {
		t.Errorf("unexpected ID: %q", m.ID)
	}
}

func TestMetadataFromDict_NullTags(t *testing.T) {
	data := map[string]interface{}{
		"id":      "x",
		"name":    "N",
		"version": "2",
		"tags":    nil,
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Tags != nil {
		t.Errorf("expected nil tags, got %v", m.Tags)
	}
}
