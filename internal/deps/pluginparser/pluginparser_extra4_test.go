package pluginparser

import (
"encoding/json"
"os"
"path/filepath"
"strings"
"testing"
)

func TestPluginManifest_EmptyAgents_Extra4(t *testing.T) {
m := PluginManifest{}
if len(m.Agents) != 0 {
t.Errorf("expected 0 agents, got %d", len(m.Agents))
}
}

func TestPluginManifest_EmptySkills_Extra4(t *testing.T) {
m := PluginManifest{}
if len(m.Skills) != 0 {
t.Errorf("expected 0 skills, got %d", len(m.Skills))
}
}

func TestPluginManifest_EmptyCommands_Extra4(t *testing.T) {
m := PluginManifest{}
if len(m.Commands) != 0 {
t.Errorf("expected 0 commands, got %d", len(m.Commands))
}
}

func TestPluginManifest_NameSet_Extra4(t *testing.T) {
m := PluginManifest{Name: "my-plugin"}
if m.Name != "my-plugin" {
t.Errorf("expected name my-plugin, got %s", m.Name)
}
}

func TestPluginManifest_MultipleSkills_Extra4(t *testing.T) {
m := PluginManifest{
Name:   "multi",
Skills: []string{"skill-a", "skill-b", "skill-c"},
}
if len(m.Skills) != 3 {
t.Fatalf("expected 3 skills, got %d", len(m.Skills))
}
if m.Skills[0] != "skill-a" {
t.Errorf("expected skill-a, got %s", m.Skills[0])
}
}

func TestMCPServerConfig_CommandField_Extra4(t *testing.T) {
cfg := MCPServerConfig{Command: "npx"}
if cfg.Command != "npx" {
t.Errorf("expected npx, got %s", cfg.Command)
}
}

func TestMCPServerConfig_URLField_Extra4(t *testing.T) {
cfg := MCPServerConfig{URL: "https://example.com/mcp"}
if cfg.URL != "https://example.com/mcp" {
t.Errorf("unexpected URL: %s", cfg.URL)
}
}

func TestMCPServerConfig_EmptyEnv_Extra4(t *testing.T) {
cfg := MCPServerConfig{}
if len(cfg.Env) != 0 {
t.Errorf("expected empty env, got %v", cfg.Env)
}
}

func TestMCPServerConfig_WithMultipleEnv_Extra4(t *testing.T) {
cfg := MCPServerConfig{
Env: map[string]string{
"KEY1": "val1",
"KEY2": "val2",
},
}
if cfg.Env["KEY1"] != "val1" {
t.Errorf("expected val1, got %s", cfg.Env["KEY1"])
}
if cfg.Env["KEY2"] != "val2" {
t.Errorf("expected val2, got %s", cfg.Env["KEY2"])
}
}

func TestMCPServerConfig_ArgsList_Extra4(t *testing.T) {
cfg := MCPServerConfig{Args: []string{"-p", "my-server"}}
if len(cfg.Args) != 2 {
t.Errorf("expected 2 args, got %d", len(cfg.Args))
}
}

func TestMCPDepEntry_URLField_Extra4(t *testing.T) {
e := MCPDepEntry{URL: "https://mcp.example.com"}
if e.URL != "https://mcp.example.com" {
t.Errorf("unexpected URL: %s", e.URL)
}
}

func TestMCPDepEntry_HeadersEmpty_Extra4(t *testing.T) {
e := MCPDepEntry{}
if len(e.Headers) != 0 {
t.Errorf("expected empty headers, got %v", e.Headers)
}
}

func TestMCPDepEntry_RegistryFalse_Extra4(t *testing.T) {
e := MCPDepEntry{}
if e.Registry {
t.Error("expected Registry to be false by default")
}
}

func TestParsePluginManifest_AgentsOnly_Extra4(t *testing.T) {
dir := t.TempDir()
pluginJSON := filepath.Join(dir, "plugin.json")
data := map[string]interface{}{
"name":   "testpkg",
"agents": []string{"agent-one", "agent-two"},
}
b, _ := json.Marshal(data)
if err := os.WriteFile(pluginJSON, b, 0o644); err != nil {
t.Fatal(err)
}
m, err := ParsePluginManifest(pluginJSON)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(m.Agents) != 2 {
t.Errorf("expected 2 agents, got %d", len(m.Agents))
}
}

func TestParsePluginManifest_SkillsOnly_Extra4(t *testing.T) {
dir := t.TempDir()
pluginJSON := filepath.Join(dir, "plugin.json")
data := map[string]interface{}{
"name":   "testpkg",
"skills": []string{"skill-one"},
}
b, _ := json.Marshal(data)
if err := os.WriteFile(pluginJSON, b, 0o644); err != nil {
t.Fatal(err)
}
m, err := ParsePluginManifest(pluginJSON)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(m.Skills) != 1 {
t.Errorf("expected 1 skill, got %d", len(m.Skills))
}
}

func TestSynthesizeApmYML_WithCommands_Extra4(t *testing.T) {
dir := t.TempDir()
m := &PluginManifest{
Name:     "cmd-pkg",
Commands: []string{"run", "build"},
}
apmYMLPath, err := SynthesizeApmYMLFromPlugin(dir, m)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !strings.HasSuffix(apmYMLPath, "apm.yml") {
t.Errorf("expected path ending in apm.yml, got: %s", apmYMLPath)
}
}
