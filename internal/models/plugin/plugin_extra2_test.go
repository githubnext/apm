package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPluginMetadata_ZeroValue_Extra2(t *testing.T) {
	var m PluginMetadata
	if m.ID != "" || m.Name != "" || m.Version != "" {
		t.Error("zero-value PluginMetadata should have empty fields")
	}
}

func TestPluginMetadata_Fields_Extra2(t *testing.T) {
	m := PluginMetadata{
		ID:          "myorg/myplugin",
		Name:        "My Plugin",
		Version:     "2.0.0",
		Description: "A test plugin",
		Author:      "Alice",
		Repository:  "https://github.com/myorg/myplugin",
		License:     "MIT",
		Tags:        []string{"testing", "utils"},
		Dependencies: []string{"dep1", "dep2"},
	}
	if m.ID != "myorg/myplugin" {
		t.Errorf("ID = %q", m.ID)
	}
	if m.Version != "2.0.0" {
		t.Errorf("Version = %q", m.Version)
	}
	if len(m.Tags) != 2 {
		t.Errorf("Tags len = %d, want 2", len(m.Tags))
	}
}

func TestToDict_HasRequiredKeys_Extra2(t *testing.T) {
	m := PluginMetadata{ID: "a/b", Name: "B", Version: "1.0.0"}
	d := m.ToDict()
	if d == nil {
		t.Fatal("ToDict returned nil")
	}
	for _, key := range []string{"id", "name", "version"} {
		if _, ok := d[key]; !ok {
			t.Errorf("ToDict missing key %q", key)
		}
	}
}

func TestToDict_TagsNotNil_Extra2(t *testing.T) {
	m := PluginMetadata{ID: "x/y", Name: "Y", Version: "0.1.0"}
	d := m.ToDict()
	tags, ok := d["tags"]
	if !ok {
		t.Error("expected tags key in ToDict")
	}
	if tags == nil {
		t.Error("tags should not be nil")
	}
}

func TestMetadataFromDict_Valid_Extra2(t *testing.T) {
	data := map[string]interface{}{
		"id":      "org/plugin",
		"name":    "Plugin",
		"version": "1.0.0",
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("MetadataFromDict error: %v", err)
	}
	if m.ID != "org/plugin" {
		t.Errorf("ID = %q", m.ID)
	}
}

func TestMetadataFromDict_MissingID_Extra2(t *testing.T) {
	data := map[string]interface{}{"name": "Plugin", "version": "1.0.0"}
	_, err := MetadataFromDict(data)
	if err == nil {
		t.Error("expected error for missing ID")
	}
}

func TestMetadataFromDict_MissingName_Extra2(t *testing.T) {
	data := map[string]interface{}{"id": "org/plugin", "version": "1.0.0"}
	_, err := MetadataFromDict(data)
	if err == nil {
		t.Error("expected error for missing Name")
	}
}

func TestMetadataFromDict_MissingVersion_Extra2(t *testing.T) {
	data := map[string]interface{}{"id": "org/plugin", "name": "Plugin"}
	_, err := MetadataFromDict(data)
	if err == nil {
		t.Error("expected error for missing Version")
	}
}

func TestPlugin_ZeroValue_Extra2(t *testing.T) {
	var p Plugin
	if p.Path != "" || len(p.Commands) != 0 || len(p.Agents) != 0 {
		t.Error("zero-value Plugin should have empty fields")
	}
}

func TestPlugin_Fields_Extra2(t *testing.T) {
	m := &PluginMetadata{ID: "x/y", Name: "Y", Version: "0.1.0"}
	p := Plugin{
		Metadata: m,
		Path:     "/some/path",
		Commands: []string{"cmd1", "cmd2"},
		Agents:   []string{"agent1"},
		Skills:   []string{"skill1", "skill2"},
	}
	if p.Path != "/some/path" {
		t.Errorf("Path = %q", p.Path)
	}
	if len(p.Commands) != 2 {
		t.Errorf("Commands len = %d", len(p.Commands))
	}
	if len(p.Skills) != 2 {
		t.Errorf("Skills len = %d", len(p.Skills))
	}
}

func TestFromPath_ValidPlugin_Extra2(t *testing.T) {
	dir := t.TempDir()
	metadata := map[string]interface{}{
		"id": "org/myplugin", "name": "MyPlugin", "version": "1.0.0",
	}
	data, _ := json.Marshal(metadata)
	_ = os.WriteFile(filepath.Join(dir, "plugin.json"), data, 0o644)

	p, err := FromPath(dir)
	if err != nil {
		t.Fatalf("FromPath error: %v", err)
	}
	if p.Metadata.ID != "org/myplugin" {
		t.Errorf("Metadata.ID = %q", p.Metadata.ID)
	}
	if p.Path != dir {
		t.Errorf("Path = %q, want %q", p.Path, dir)
	}
}

func TestFromPath_MissingPluginJSON_Extra2(t *testing.T) {
	dir := t.TempDir()
	_, err := FromPath(dir)
	if err == nil {
		t.Error("expected error when plugin.json is absent")
	}
}

func TestMetadataFromDict_Tags_Extra2(t *testing.T) {
	data := map[string]interface{}{
		"id":      "x/y",
		"name":    "Y",
		"version": "0.1",
		"tags":    []interface{}{"a", "b"},
	}
	m, err := MetadataFromDict(data)
	if err != nil {
		t.Fatalf("MetadataFromDict error: %v", err)
	}
	if len(m.Tags) != 2 {
		t.Errorf("Tags len = %d, want 2", len(m.Tags))
	}
}
