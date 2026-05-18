package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetadataFromDict_required(t *testing.T) {
	data := map[string]interface{}{
		"id":      "test-id",
		"name":    "Test Plugin",
		"version": "1.0.0",
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.ID != "test-id" || m.Name != "Test Plugin" || m.Version != "1.0.0" {
		t.Errorf("unexpected fields: %+v", m)
	}
}

func TestMetadataFromDict_missingID(t *testing.T) {
	data := map[string]interface{}{"name": "n", "version": "1"}
	_, err := MetadataFromDict(data)
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestMetadataFromDict_missingName(t *testing.T) {
	data := map[string]interface{}{"id": "x", "version": "1"}
	_, err := MetadataFromDict(data)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestMetadataFromDict_missingVersion(t *testing.T) {
	data := map[string]interface{}{"id": "x", "name": "y"}
	_, err := MetadataFromDict(data)
	if err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestMetadataFromDict_optional(t *testing.T) {
	data := map[string]interface{}{
		"id":           "id1",
		"name":         "Plugin A",
		"version":      "2.0.0",
		"description":  "A plugin",
		"author":       "Alice",
		"repository":   "https://github.com/a/b",
		"homepage":     "https://example.com",
		"license":      "MIT",
		"tags":         []interface{}{"tag1", "tag2"},
		"dependencies": []interface{}{"dep1"},
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Tags) != 2 || m.Tags[0] != "tag1" {
		t.Errorf("unexpected tags: %v", m.Tags)
	}
	if len(m.Dependencies) != 1 || m.Dependencies[0] != "dep1" {
		t.Errorf("unexpected deps: %v", m.Dependencies)
	}
	if m.License != "MIT" {
		t.Errorf("unexpected license: %s", m.License)
	}
}

func TestToDict_nilSlices(t *testing.T) {
	m := &PluginMetadata{ID: "x", Name: "n", Version: "1"}
	d := m.ToDict()
	tags, ok := d["tags"].([]string)
	if !ok {
		t.Fatalf("tags not []string: %T", d["tags"])
	}
	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}

func TestFromPath_noMetadata(t *testing.T) {
	dir := t.TempDir()
	_, err := FromPath(dir)
	if err == nil {
		t.Fatal("expected error for missing plugin.json")
	}
}

func TestFromPath_validPlugin(t *testing.T) {
	dir := t.TempDir()
	pluginJSON := `{"id":"p1","name":"Plugin1","version":"0.1.0"}`
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(pluginJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := FromPath(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Metadata.ID != "p1" {
		t.Errorf("unexpected ID: %s", p.Metadata.ID)
	}
	if p.Path != dir {
		t.Errorf("unexpected path: %s", p.Path)
	}
}

func TestFromPath_githubSubdir(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, ".github", "plugin")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	pluginJSON := `{"id":"p2","name":"Plugin2","version":"0.2.0"}`
	if err := os.WriteFile(filepath.Join(subdir, "plugin.json"), []byte(pluginJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := FromPath(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Metadata.ID != "p2" {
		t.Errorf("unexpected ID: %s", p.Metadata.ID)
	}
}

func TestFromPath_withSkills(t *testing.T) {
	dir := t.TempDir()
	pluginJSON := `{"id":"p3","name":"Plugin3","version":"0.3.0"}`
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(pluginJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	skillDir := filepath.Join(dir, "skills", "myskill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := FromPath(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(p.Skills))
	}
}
