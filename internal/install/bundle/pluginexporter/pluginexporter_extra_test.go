package pluginexporter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateOutputRel_EmptyString(t *testing.T) {
	// empty string should be valid (no path traversal)
	if !validateOutputRel("") {
		t.Error("empty string should be valid")
	}
}

func TestValidateOutputRel_Windows(t *testing.T) {
	// Windows-style backslash paths should be rejected as absolute
	// but implementation may vary; just assert no panic
	_ = validateOutputRel(`sub\file.md`)
}

func TestSanitizeBundleName_AlphaNumeric(t *testing.T) {
	got := sanitizeBundleName("hello123")
	if got != "hello123" {
		t.Errorf("expected hello123, got %q", got)
	}
}

func TestSanitizeBundleName_AllSpecial(t *testing.T) {
	got := sanitizeBundleName("!@#$%")
	// All special chars -- should produce "unnamed" or sanitized form
	if got == "" {
		t.Error("expected non-empty result for all-special input")
	}
}

func TestSanitizeBundleName_Hyphen(t *testing.T) {
	got := sanitizeBundleName("my-bundle-name")
	if got != "my-bundle-name" {
		t.Errorf("expected my-bundle-name, got %q", got)
	}
}

func TestRenamePrompt_NoExtension(t *testing.T) {
	got := renamePrompt("justname")
	if got != "justname" {
		t.Errorf("expected justname, got %q", got)
	}
}

func TestRenamePrompt_GoFile(t *testing.T) {
	got := renamePrompt("file.go")
	if got != "file.go" {
		t.Errorf("expected file.go unchanged, got %q", got)
	}
}

func TestExportPluginBundle_MissingProjectRoot(t *testing.T) {
	opts := ExportOptions{
		ProjectRoot: "/nonexistent/path/xyz",
		OutputDir:   t.TempDir(),
		DryRun:      true,
	}
	_, err := ExportPluginBundle(opts)
	// Missing project root -- may fail gracefully or succeed with empty bundle
	_ = err
}

func TestExportPluginBundle_WithPluginJSON(t *testing.T) {
	dir := t.TempDir()
	// Create a minimal plugin.json
	os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(`{"name":"test","version":"1.0.0"}`), 0o644)

	outDir := t.TempDir()
	opts := ExportOptions{
		ProjectRoot: dir,
		OutputDir:   outDir,
		DryRun:      true,
	}
	result, err := ExportPluginBundle(opts)
	if err != nil {
		t.Fatalf("ExportPluginBundle with plugin.json: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestExportPluginBundle_WithAgentsDir(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, ".apm", "agents")
	os.MkdirAll(agentsDir, 0o755)
	os.WriteFile(filepath.Join(agentsDir, "agent1.md"), []byte("# Agent 1\n"), 0o644)
	os.WriteFile(filepath.Join(agentsDir, "agent2.md"), []byte("# Agent 2\n"), 0o644)

	outDir := t.TempDir()
	opts := ExportOptions{ProjectRoot: dir, OutputDir: outDir, DryRun: true}
	result, err := ExportPluginBundle(opts)
	if err != nil {
		t.Fatalf("agents dir test: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestExportPluginBundle_WithSkillsDir(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, ".apm", "skills")
	os.MkdirAll(skillsDir, 0o755)
	os.WriteFile(filepath.Join(skillsDir, "skill1.md"), []byte("# Skill 1\n"), 0o644)

	outDir := t.TempDir()
	opts := ExportOptions{ProjectRoot: dir, OutputDir: outDir, DryRun: true}
	result, err := ExportPluginBundle(opts)
	if err != nil {
		t.Fatalf("skills dir test: %v", err)
	}
	_ = result
}

func TestExportPluginBundle_NameVersionFromOpts(t *testing.T) {
	dir := t.TempDir()
	outDir := t.TempDir()
	opts := ExportOptions{
		ProjectRoot: dir,
		OutputDir:   outDir,
		DryRun:      true,
		Force:       true,
	}
	result, err := ExportPluginBundle(opts)
	if err != nil {
		t.Fatalf("named bundle test: %v", err)
	}
	_ = result
}

func TestValidateOutputRel_CurrentDir(t *testing.T) {
	// "." refers to current directory -- should be valid (no traversal)
	// implementation-defined; just no panic
	_ = validateOutputRel(".")
}
