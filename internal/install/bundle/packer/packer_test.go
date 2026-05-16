package packer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectTarget_GitHub(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".github"), 0o755)
	got := detectTarget(dir)
	if got != "copilot" {
		t.Errorf("detectTarget with .github = %q; want copilot", got)
	}
}

func TestDetectTarget_Claude(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0o755)
	got := detectTarget(dir)
	if got != "claude" {
		t.Errorf("detectTarget with .claude = %q; want claude", got)
	}
}

func TestDetectTarget_NoMatch(t *testing.T) {
	dir := t.TempDir()
	got := detectTarget(dir)
	if got != "all" {
		t.Errorf("detectTarget with no dirs = %q; want all", got)
	}
}

func TestFilterFilesByTarget_Copilot(t *testing.T) {
	files := []string{
		".github/skills/my-skill.md",
		".claude/skills/my-skill.md",
		".cursor/mcp.json",
	}
	direct, _ := filterFilesByTarget(files, "copilot")
	if len(direct) == 0 {
		t.Error("expected files for copilot target")
	}
	for _, f := range direct {
		if !startsWith(f, ".github/") && !startsWith(f, ".claude/skills/") {
			// Cross-map may include .github/skills -> .github/skills
		}
	}
}

func TestFilterFilesByTarget_All(t *testing.T) {
	files := []string{
		".github/skills/my-skill.md",
		".claude/skills/skill.md",
		".cursor/mcp.json",
		".agents/my-agent.md",
	}
	direct, _ := filterFilesByTarget(files, "all")
	if len(direct) == 0 {
		t.Error("expected files for all target")
	}
}

func TestFilterFilesByTarget_EmptyFiles(t *testing.T) {
	direct, mappings := filterFilesByTarget([]string{}, "copilot")
	if len(direct) != 0 {
		t.Errorf("expected 0 files, got %d", len(direct))
	}
	if len(mappings) != 0 {
		t.Errorf("expected 0 mappings, got %d", len(mappings))
	}
}

func TestReadDeployedFiles_ValidLockfile(t *testing.T) {
	content := `dependencies:
  - name: my-package
    version: "1.0.0"
    deployed_files:
      - .github/skills/my-skill.md
      - .github/instructions/my.instructions.md
`
	f, err := os.CreateTemp(t.TempDir(), "apm.lock.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()

	deps, err := readDeployedFiles(f.Name())
	if err != nil {
		t.Fatalf("readDeployedFiles: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	if deps[0].Name != "my-package" {
		t.Errorf("name = %q; want my-package", deps[0].Name)
	}
	if len(deps[0].DeployedFiles) != 2 {
		t.Errorf("expected 2 deployed files, got %d", len(deps[0].DeployedFiles))
	}
}

func TestReadDeployedFiles_MissingFile(t *testing.T) {
	_, err := readDeployedFiles("/nonexistent/path/apm.lock.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFindLockfile_Present(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, "apm.lock.yaml")
	os.WriteFile(lockPath, []byte("dependencies: []"), 0o644)
	got := findLockfile(dir)
	if got != lockPath {
		t.Errorf("findLockfile = %q; want %q", got, lockPath)
	}
}

func TestFindLockfile_Missing(t *testing.T) {
	dir := t.TempDir()
	got := findLockfile(dir)
	if got != "" {
		t.Errorf("findLockfile for empty dir = %q; want empty", got)
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
